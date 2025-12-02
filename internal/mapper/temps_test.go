package mapper

import (
	"testing"

	"github.com/davidcollom/creality2mqtt/internal/types"
)

func TestBuildTempMessages(t *testing.T) {
	tests := []struct {
		name      string
		input     map[string]any
		baseTopic string
		expected  []types.MqttMessage
	}{
		{
			name: "all temperature fields present as float64",
			input: map[string]any{
				"nozzleTemp":       219.900,
				"targetNozzleTemp": 220.000,
				"bedTemp0":         60.500,
				"targetBedTemp0":   65.000,
				"boxTemp":          45.200,
			},
			baseTopic: "3dprinter/k1se",
			expected: []types.MqttMessage{
				{Topic: "3dprinter/k1se/temperature/nozzle/current", Payload: "219.900", Retain: false},
				{Topic: "3dprinter/k1se/temperature/nozzle/target", Payload: "220.000", Retain: false},
				{Topic: "3dprinter/k1se/temperature/bed0/current", Payload: "60.500", Retain: false},
				{Topic: "3dprinter/k1se/temperature/bed0/target", Payload: "65.000", Retain: false},
				{Topic: "3dprinter/k1se/temperature/box/current", Payload: "45.200", Retain: false},
			},
		},
		{
			name: "temperature fields as string numbers",
			input: map[string]any{
				"nozzleTemp":       "219.900000",
				"targetNozzleTemp": "220.000000",
				"bedTemp0":         "60.500000",
				"targetBedTemp0":   "65.000000",
				"boxTemp":          "45.200000",
			},
			baseTopic: "3dprinter/k1se",
			expected: []types.MqttMessage{
				{Topic: "3dprinter/k1se/temperature/nozzle/current", Payload: "219.900", Retain: false},
				{Topic: "3dprinter/k1se/temperature/nozzle/target", Payload: "220.000", Retain: false},
				{Topic: "3dprinter/k1se/temperature/bed0/current", Payload: "60.500", Retain: false},
				{Topic: "3dprinter/k1se/temperature/bed0/target", Payload: "65.000", Retain: false},
				{Topic: "3dprinter/k1se/temperature/box/current", Payload: "45.200", Retain: false},
			},
		},
		{
			name: "mixed integer and float types",
			input: map[string]any{
				"nozzleTemp":       int(220),
				"targetNozzleTemp": int32(225),
				"bedTemp0":         int64(60),
				"targetBedTemp0":   65.0,
				"boxTemp":          "45",
			},
			baseTopic: "test/printer",
			expected: []types.MqttMessage{
				{Topic: "test/printer/temperature/nozzle/current", Payload: "220.000", Retain: false},
				{Topic: "test/printer/temperature/nozzle/target", Payload: "225.000", Retain: false},
				{Topic: "test/printer/temperature/bed0/current", Payload: "60.000", Retain: false},
				{Topic: "test/printer/temperature/bed0/target", Payload: "65.000", Retain: false},
				{Topic: "test/printer/temperature/box/current", Payload: "45.000", Retain: false},
			},
		},
		{
			name: "partial temperature data",
			input: map[string]any{
				"nozzleTemp": 200.5,
				"bedTemp0":   55.0,
				"boxTemp":    40.2,
			},
			baseTopic: "printer/test",
			expected: []types.MqttMessage{
				{Topic: "printer/test/temperature/nozzle/current", Payload: "200.500", Retain: false},
				{Topic: "printer/test/temperature/bed0/current", Payload: "55.000", Retain: false},
				{Topic: "printer/test/temperature/box/current", Payload: "40.200", Retain: false},
			},
		},
		{
			name:      "empty input",
			input:     map[string]any{},
			baseTopic: "3dprinter/k1se",
			expected:  []types.MqttMessage{},
		},
		{
			name: "invalid temperature values ignored",
			input: map[string]any{
				"nozzleTemp":       "invalid",
				"targetNozzleTemp": nil,
				"bedTemp0":         true,
				"targetBedTemp0":   []string{"not", "a", "number"},
				"boxTemp":          50.0, // valid one
			},
			baseTopic: "3dprinter/k1se",
			expected: []types.MqttMessage{
				{Topic: "3dprinter/k1se/temperature/box/current", Payload: "50.000", Retain: false},
			},
		},
		{
			name: "zero values are valid",
			input: map[string]any{
				"nozzleTemp":       0.0,
				"targetNozzleTemp": "0",
				"bedTemp0":         int(0),
			},
			baseTopic: "3dprinter/k1se",
			expected: []types.MqttMessage{
				{Topic: "3dprinter/k1se/temperature/nozzle/current", Payload: "0.000", Retain: false},
				{Topic: "3dprinter/k1se/temperature/nozzle/target", Payload: "0.000", Retain: false},
				{Topic: "3dprinter/k1se/temperature/bed0/current", Payload: "0.000", Retain: false},
			},
		},
		{
			name: "different base topic",
			input: map[string]any{
				"nozzleTemp": 210.5,
			},
			baseTopic: "home/3dprinter/creality",
			expected: []types.MqttMessage{
				{Topic: "home/3dprinter/creality/temperature/nozzle/current", Payload: "210.500", Retain: false},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildTempMessages(tt.input, tt.baseTopic)

			if len(result) != len(tt.expected) {
				t.Errorf("expected %d messages, got %d", len(tt.expected), len(result))
				return
			}

			// Convert to maps for easier comparison since order might vary
			resultMap := make(map[string]types.MqttMessage)
			expectedMap := make(map[string]types.MqttMessage)

			for _, msg := range result {
				resultMap[msg.Topic] = msg
			}
			for _, msg := range tt.expected {
				expectedMap[msg.Topic] = msg
			}

			for topic, expectedMsg := range expectedMap {
				resultMsg, exists := resultMap[topic]
				if !exists {
					t.Errorf("expected topic %s not found in result", topic)
					continue
				}

				if resultMsg.Topic != expectedMsg.Topic {
					t.Errorf("topic mismatch: expected %s, got %s", expectedMsg.Topic, resultMsg.Topic)
				}
				if resultMsg.Payload != expectedMsg.Payload {
					t.Errorf("payload mismatch for topic %s: expected %s, got %s", topic, expectedMsg.Payload, resultMsg.Payload)
				}
				if resultMsg.Retain != expectedMsg.Retain {
					t.Errorf("retain mismatch for topic %s: expected %t, got %t", topic, expectedMsg.Retain, resultMsg.Retain)
				}
			}
		})
	}
}

func TestGetFloat(t *testing.T) {
	tests := []struct {
		name     string
		msg      map[string]any
		key      string
		expected float64
		ok       bool
	}{
		{
			name:     "float64 value",
			msg:      map[string]any{"temp": 219.9},
			key:      "temp",
			expected: 219.9,
			ok:       true,
		},
		{
			name:     "int value",
			msg:      map[string]any{"temp": 220},
			key:      "temp",
			expected: 220.0,
			ok:       true,
		},
		{
			name:     "int32 value",
			msg:      map[string]any{"temp": int32(225)},
			key:      "temp",
			expected: 225.0,
			ok:       true,
		},
		{
			name:     "int64 value",
			msg:      map[string]any{"temp": int64(230)},
			key:      "temp",
			expected: 230.0,
			ok:       true,
		},
		{
			name:     "string numeric value",
			msg:      map[string]any{"temp": "219.900000"},
			key:      "temp",
			expected: 219.9,
			ok:       true,
		},
		{
			name:     "string zero",
			msg:      map[string]any{"temp": "0"},
			key:      "temp",
			expected: 0.0,
			ok:       true,
		},
		{
			name: "key not found",
			msg:  map[string]any{"other": 123},
			key:  "temp",
			ok:   false,
		},
		{
			name: "invalid string",
			msg:  map[string]any{"temp": "not_a_number"},
			key:  "temp",
			ok:   false,
		},
		{
			name: "nil value",
			msg:  map[string]any{"temp": nil},
			key:  "temp",
			ok:   false,
		},
		{
			name: "boolean value",
			msg:  map[string]any{"temp": true},
			key:  "temp",
			ok:   false,
		},
		{
			name: "slice value",
			msg:  map[string]any{"temp": []int{1, 2, 3}},
			key:  "temp",
			ok:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := getFloat(tt.msg, tt.key)

			if ok != tt.ok {
				t.Errorf("expected ok=%t, got ok=%t", tt.ok, ok)
				return
			}

			if tt.ok && result != tt.expected {
				t.Errorf("expected %f, got %f", tt.expected, result)
			}
		})
	}
}
