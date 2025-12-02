package mapper

import (
	"fmt"

	"github.com/davidcollom/creality2mqtt/internal/types"
)

// BuildCFSBoxMessages maps Creality Filament System (CFS) boxState to MQTT topics
// Expected payload structure:
// { "boxState": { "id": 1, "state": 1, "humidity": 28.0, "temp": 23.0 } }
func BuildCFSBoxMessages(msg map[string]any, baseTopic string) []types.MqttMessage {
	out := make([]types.MqttMessage, 0, 3)

	raw, ok := msg["boxState"]
	if !ok {
		return out
	}

	m, ok := raw.(map[string]any)
	if !ok {
		return out
	}

	idAny, hasID := m["id"]
	if !hasID {
		return out
	}

	id := toInt(idAny)
	boxPrefix := fmt.Sprintf("%s/cfs/%d", baseTopic, id)

	if humAny, ok := m["humidity"]; ok {
		hum := coerceValue(humAny)
		out = append(out, types.MqttMessage{
			Topic:   fmt.Sprintf("%s/humidity", boxPrefix),
			Payload: fmt.Sprint(hum),
			Retain:  false,
		})
	}

	if tAny, ok := m["temp"]; ok {
		t := coerceValue(tAny)
		out = append(out, types.MqttMessage{
			Topic:   fmt.Sprintf("%s/temperature", boxPrefix),
			Payload: fmt.Sprint(t),
			Retain:  false,
		})
	}

	if sAny, ok := m["state"]; ok {
		s := toInt(sAny)
		out = append(out, types.MqttMessage{
			Topic:   fmt.Sprintf("%s/state", boxPrefix),
			Payload: fmt.Sprintf("%d", s),
			Retain:  false,
		})
	}

	return out
}

func toInt(v any) int {
	switch x := v.(type) {
	case int:
		return x
	case int32:
		return int(x)
	case int64:
		return int(x)
	case float32:
		return int(x)
	case float64:
		return int(x)
	default:
		return 0
	}
}
