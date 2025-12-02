package mapper

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/davidcollom/creality2mqtt/internal/types"
)

// BuildJobMessages emits derived MQTT topics for the current print job:
//
//	<base>/printing              -> "true"/"false"
//	<base>/job/progress          -> int percentage (string)
//	<base>/job/left_time         -> seconds remaining (string)
//	<base>/job/job_time          -> seconds elapsed (string)
//	<base>/job/layer/current     -> current layer (string)
//	<base>/job/layer/total       -> total layers (string) if non-zero
//	<base>/job/file_name         -> simplified file name (no full path)
//	<base>/feed_state            -> extruder/feed state code (101=extruding, 102=done, etc.)
func BuildJobMessages(msg map[string]any, baseTopic string) []types.MqttMessage {
	out := make([]types.MqttMessage, 0, 10)

	progress, hasProgress := getInt(msg, "printProgress")
	left, hasLeft := getInt(msg, "printLeftTime")
	jobTime, hasJobTime := getInt(msg, "printJobTime")
	layer, hasLayer := getInt(msg, "layer")
	totalLayer, hasTotalLayer := getInt(msg, "TotalLayer")

	// printing heuristic:
	// - if progress > 0 and left > 0, we treat it as printing.
	printing := hasProgress && hasLeft && progress > 0 && left > 0
	printingPayload := "false"
	if printing {
		printingPayload = "true"
	}
	out = append(out, types.MqttMessage{
		Topic:   fmt.Sprintf("%s/printing", baseTopic),
		Payload: printingPayload,
		Retain:  false,
	})

	if hasProgress {
		out = append(out, types.MqttMessage{
			Topic:   fmt.Sprintf("%s/job/progress", baseTopic),
			Payload: strconv.FormatInt(progress, 10),
			Retain:  false,
		})
	}

	if hasLeft {
		out = append(out, types.MqttMessage{
			Topic:   fmt.Sprintf("%s/job/left_time", baseTopic),
			Payload: strconv.FormatInt(left, 10),
			Retain:  false,
		})
	}

	if hasJobTime {
		out = append(out, types.MqttMessage{
			Topic:   fmt.Sprintf("%s/job/job_time", baseTopic),
			Payload: strconv.FormatInt(jobTime, 10),
			Retain:  false,
		})
	}

	if hasLayer {
		out = append(out, types.MqttMessage{
			Topic:   fmt.Sprintf("%s/job/layer/current", baseTopic),
			Payload: strconv.FormatInt(layer, 10),
			Retain:  false,
		})
	}

	if hasTotalLayer && totalLayer > 0 {
		out = append(out, types.MqttMessage{
			Topic:   fmt.Sprintf("%s/job/layer/total", baseTopic),
			Payload: strconv.FormatInt(totalLayer, 10),
			Retain:  false,
		})
	}

	// Derive a human-friendly filename from printFileName, which is often
	// something like "/usr/data/printer_data/gcodes/.../foo.gcode".
	if raw, ok := msg["printFileName"]; ok {
		if s, ok := raw.(string); ok && s != "" {
			short := simplifyFileName(s)
			out = append(out, types.MqttMessage{
				Topic:   fmt.Sprintf("%s/job/file_name", baseTopic),
				Payload: short,
				Retain:  false,
			})
		}
	}

	// Feed state (extruder activity: 101=extruding, 102=done, etc.)
	if feedState, hasFeedState := getInt(msg, "feedState"); hasFeedState {
		out = append(out, types.MqttMessage{
			Topic:   fmt.Sprintf("%s/feed_state", baseTopic),
			Payload: strconv.FormatInt(feedState, 10),
			Retain:  false,
		})
	}

	return out
}

// getInt normalises numeric-ish values to int64.
func getInt(msg map[string]any, key string) (int64, bool) {
	raw, ok := msg[key]
	if !ok {
		return 0, false
	}

	switch v := raw.(type) {
	case float64:
		return int64(v), true
	case int:
		return int64(v), true
	case int32:
		return int64(v), true
	case int64:
		return v, true
	case string:
		val, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, false
		}
		return val, true
	default:
		return 0, false
	}
}

// simplifyFileName tries to give a nice filename for HA display.
//
// For paths like ".../.Meta+2+Stock+V3 (1)_gcode.3mf/Meta+2+..._plate_4.gcode"
// we'll basically just take the last segment ("Meta+2+..._plate_4.gcode").
func simplifyFileName(full string) string {
	full = strings.TrimSpace(full)
	if full == "" {
		return full
	}

	// Normalize Windows-style separators to POSIX for filepath.Base
	full = strings.ReplaceAll(full, "\\", "/")

	// filepath.Base works fine even with mixed directory structures.
	base := filepath.Base(full)

	// If for some reason we still have weird suffixes, you can trim them here.
	return base
}
