package mapper

import (
	"fmt"
	"strconv"

	"github.com/davidcollom/creality2mqtt/internal/types"
)

// BuildTempMessages emits derived MQTT topics for temperatures.
//
// It does NOT replace the generic "nozzle_temp", "bed_temp0" etc. topics
// produced by MapMessageToMqtt – it adds more structured ones, e.g.:
//
//	<base>/temperature/nozzle/current
//	<base>/temperature/nozzle/target
//	<base>/temperature/bed0/current
//	<base>/temperature/bed0/target
//	<base>/temperature/box/current
func BuildTempMessages(msg map[string]any, baseTopic string) []types.MqttMessage {
	out := make([]types.MqttMessage, 0, 8)

	if v, ok := getFloat(msg, "nozzleTemp"); ok {
		payload := fmt.Sprintf("%.3f", v)
		out = append(out, types.MqttMessage{
			Topic:   fmt.Sprintf("%s/temperature/nozzle/current", baseTopic),
			Payload: payload,
			Retain:  false,
		})
	}

	if v, ok := getFloat(msg, "targetNozzleTemp"); ok {
		payload := fmt.Sprintf("%.3f", v)
		out = append(out, types.MqttMessage{
			Topic:   fmt.Sprintf("%s/temperature/nozzle/target", baseTopic),
			Payload: payload,
			Retain:  false,
		})
	}

	if v, ok := getFloat(msg, "bedTemp0"); ok {
		payload := fmt.Sprintf("%.3f", v)
		out = append(out, types.MqttMessage{
			Topic:   fmt.Sprintf("%s/temperature/bed0/current", baseTopic),
			Payload: payload,
			Retain:  false,
		})
	}

	if v, ok := getFloat(msg, "targetBedTemp0"); ok {
		payload := fmt.Sprintf("%.3f", v)
		out = append(out, types.MqttMessage{
			Topic:   fmt.Sprintf("%s/temperature/bed0/target", baseTopic),
			Payload: payload,
			Retain:  false,
		})
	}

	if v, ok := getFloat(msg, "boxTemp"); ok {
		payload := fmt.Sprintf("%.3f", v)
		out = append(out, types.MqttMessage{
			Topic:   fmt.Sprintf("%s/temperature/box/current", baseTopic),
			Payload: payload,
			Retain:  false,
		})
	}

	return out
}

// getFloat tries to normalise "numeric ish" values to float64.
// It is intentionally a bit defensive because the printer sometimes
// sends numeric strings ("219.900000") and sometimes numbers.
func getFloat(msg map[string]any, key string) (float64, bool) {
	raw, ok := msg[key]
	if !ok {
		return 0, false
	}

	switch v := raw.(type) {
	case float64:
		return v, true
	case int:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, false
		}
		return f, true
	default:
		return 0, false
	}
}
