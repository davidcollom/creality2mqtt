package mapper

import (
	"fmt"
	"sync"
	"time"

	"github.com/davidcollom/creality2mqtt/internal/types"
)

var printerStatusCache struct {
	mu            sync.RWMutex
	lastStatus    string
	lastPublished time.Time
	lastUpdate    time.Time
}
var statusUpdateInterval = 10 * time.Second

// BuildStateMessages emits derived MQTT topics around device state:
//
//	<base>/printer_status    -> "idle"/"active" (based on print activity)
//	<base>/tf_card_present   -> "true"/"false" (based on tfCard)
func BuildStateMessages(msg map[string]any, baseTopic string) []types.MqttMessage {
	out := make([]types.MqttMessage, 0, 4)

	// Determine printer status with rate limiting
	if status, shouldPublish := getRateLimitedPrinterStatus(msg); shouldPublish {
		out = append(out, types.MqttMessage{
			Topic:   fmt.Sprintf("%s/printer_status", baseTopic),
			Payload: status,
			Retain:  false,
		})
	}

	if tf, ok := getInt(msg, "tfCard"); ok {
		tfPresent := tf == 1
		tfPayload := "false"
		if tfPresent {
			tfPayload = "true"
		}
		out = append(out, types.MqttMessage{
			Topic:   fmt.Sprintf("%s/tf_card_present", baseTopic),
			Payload: tfPayload,
			Retain:  false,
		})
	}
	return out
}

// getRateLimitedPrinterStatus returns status and whether it should be published
// Only publishes updates every 10 seconds or when status changes significantly
func getRateLimitedPrinterStatus(msg map[string]any) (string, bool) {
	printerStatusCache.mu.Lock()
	defer printerStatusCache.mu.Unlock()

	now := time.Now()
	currentStatus := derivePrinterStatus(msg)

	// Update last activity time if printer is active
	if currentStatus == "active" {
		printerStatusCache.lastUpdate = now
	}

	// Check if we should mark as idle due to timeout
	if currentStatus == "idle" && !printerStatusCache.lastUpdate.IsZero() {
		if now.Sub(printerStatusCache.lastUpdate) < statusUpdateInterval {
			// Still within timeout period, keep previous status if it was active
			if printerStatusCache.lastStatus == "active" {
				currentStatus = "active"
			}
		}
	}

	// Only publish if status changed or enough time has passed
	shouldPublish := false
	if currentStatus != printerStatusCache.lastStatus {
		shouldPublish = true
	} else if now.Sub(printerStatusCache.lastPublished) >= statusUpdateInterval {
		shouldPublish = true
	}

	if shouldPublish {
		printerStatusCache.lastStatus = currentStatus
		printerStatusCache.lastPublished = now
	}

	return currentStatus, shouldPublish
}

// derivePrinterStatus returns "active" when printing, "idle" otherwise
// When bridge is offline, HA will show "unavailable" via availability topic
func derivePrinterStatus(msg map[string]any) string {
	// Check if actively printing (progress > 0 and time left > 0)
	if progress, ok := getInt(msg, "printProgress"); ok && progress > 0 {
		if left, ok := getInt(msg, "leftTime"); ok && left > 0 {
			return "active"
		}
	}

	// If we have a printJobTime greater than 0 we're active
	if left, ok := getInt(msg, "printJobTime"); ok && left > 0 {
		return "active"
	}

	// If we have a print time left greater than 0 we're active
	if left, ok := getInt(msg, "printLeftTime"); ok && left > 0 {
		return "active"
	}

	// If we have a layer greater than 0 we're active
	if left, ok := getInt(msg, "layer"); ok && left > 0 {
		return "active"
	}

	// Check gcode state (printing/paused)
	if state, ok := getInt(msg, "gcodeState"); ok && state > 0 {
		return "active"
	}

	return "idle"
}
