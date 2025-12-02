package mapper

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/davidcollom/creality2mqtt/internal/types"
)

func normaliseKey(key string) string {
	key = strings.TrimSpace(key)
	key = strings.ReplaceAll(key, " ", "_")

	var out []rune
	runes := []rune(key)

	for i, r := range runes {
		if i > 0 && unicode.IsUpper(r) && unicode.IsLower(runes[i-1]) {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(r))
	}

	return string(out)
}

func coerceValue(v any) any {
	switch val := v.(type) {
	case string:
		// try to parse numeric strings
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			// format to a sensible precision for temps etc.
			return fmt.Sprintf("%.3f", f)
		}
		return val
	default:
		return val
	}
}

func DecodeAndMap(data []byte, baseTopic string) ([]types.MqttMessage, error) {
	var msg map[string]any
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return MapMessageToMqtt(msg, baseTopic), nil
}

func MapMessageToMqtt(msg map[string]any, baseTopic string) []types.MqttMessage {
	result := make([]types.MqttMessage, 0, len(msg)+16)

	// --- generic scalar → topic mapping ---
	for key, raw := range msg {
		if _, skip := noisyKeys[key]; skip {
			continue
		}

		norm := normaliseKey(key)

		switch raw.(type) {
		case map[string]any, []any:
			continue
		}

		val := coerceValue(raw)
		payload := fmt.Sprint(val)

		topic := fmt.Sprintf("%s/%s", baseTopic, norm)
		result = append(result, types.MqttMessage{
			Topic:   topic,
			Payload: payload,
			Retain:  false,
		})
	}

	// --- domain-specific derived topics ---
	result = append(result, BuildTempMessages(msg, baseTopic)...)
	result = append(result, BuildJobMessages(msg, baseTopic)...)
	result = append(result, BuildStateMessages(msg, baseTopic)...)
	result = append(result, BuildCFSBoxMessages(msg, baseTopic)...)

	return result
}
