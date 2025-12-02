package mapper

import (
	"testing"

	"github.com/davidcollom/creality2mqtt/internal/types"
)

func TestBuildCFSBoxMessages(t *testing.T) {
	tests := []struct {
		name      string
		input     map[string]any
		baseTopic string
		expected  []types.MqttMessage
	}{
		{
			name: "complete boxState with all fields",
			input: map[string]any{
				"boxState": map[string]any{
					"id":       1,
					"state":    1,
					"humidity": 28.0,
					"temp":     23.0,
				},
			},
			baseTopic: "3dprinter/k1se",
			expected: []types.MqttMessage{
				{Topic: "3dprinter/k1se/cfs/1/humidity", Payload: "28", Retain: false},
				{Topic: "3dprinter/k1se/cfs/1/temperature", Payload: "23", Retain: false},
				{Topic: "3dprinter/k1se/cfs/1/state", Payload: "1", Retain: false},
			},
		},
		{
			name: "boxState with different id",
			input: map[string]any{
				"boxState": map[string]any{
					"id":       2,
					"state":    0,
					"humidity": 35.5,
					"temp":     25.3,
				},
			},
			baseTopic: "printer/test",
			expected: []types.MqttMessage{
				{Topic: "printer/test/cfs/2/humidity", Payload: "35.5", Retain: false},
				{Topic: "printer/test/cfs/2/temperature", Payload: "25.3", Retain: false},
				{Topic: "printer/test/cfs/2/state", Payload: "0", Retain: false},
			},
		},
		{
			name: "partial boxState - only humidity",
			input: map[string]any{
				"boxState": map[string]any{
					"id":       1,
					"humidity": 30.0,
				},
			},
			baseTopic: "3dprinter/k1se",
			expected: []types.MqttMessage{
				{Topic: "3dprinter/k1se/cfs/1/humidity", Payload: "30", Retain: false},
			},
		},
		{
			name: "partial boxState - only temperature",
			input: map[string]any{
				"boxState": map[string]any{
					"id":   1,
					"temp": 22.5,
				},
			},
			baseTopic: "3dprinter/k1se",
			expected: []types.MqttMessage{
				{Topic: "3dprinter/k1se/cfs/1/temperature", Payload: "22.5", Retain: false},
			},
		},
		{
			name: "partial boxState - only state",
			input: map[string]any{
				"boxState": map[string]any{
					"id":    1,
					"state": 2,
				},
			},
			baseTopic: "3dprinter/k1se",
			expected: []types.MqttMessage{
				{Topic: "3dprinter/k1se/cfs/1/state", Payload: "2", Retain: false},
			},
		},
		{
			name: "boxState with integer types for numeric values",
			input: map[string]any{
				"boxState": map[string]any{
					"id":       3,
					"state":    1,
					"humidity": 28,
					"temp":     23,
				},
			},
			baseTopic: "3dprinter/k1se",
			expected: []types.MqttMessage{
				{Topic: "3dprinter/k1se/cfs/3/humidity", Payload: "28", Retain: false},
				{Topic: "3dprinter/k1se/cfs/3/temperature", Payload: "23", Retain: false},
				{Topic: "3dprinter/k1se/cfs/3/state", Payload: "1", Retain: false},
			},
		},
		{
			name: "boxState with float id",
			input: map[string]any{
				"boxState": map[string]any{
					"id":       1.0,
					"humidity": 28.0,
				},
			},
			baseTopic: "3dprinter/k1se",
			expected: []types.MqttMessage{
				{Topic: "3dprinter/k1se/cfs/1/humidity", Payload: "28", Retain: false},
			},
		},
		{
			name: "boxState with zero values",
			input: map[string]any{
				"boxState": map[string]any{
					"id":       0,
					"state":    0,
					"humidity": 0.0,
					"temp":     0.0,
				},
			},
			baseTopic: "3dprinter/k1se",
			expected: []types.MqttMessage{
				{Topic: "3dprinter/k1se/cfs/0/humidity", Payload: "0", Retain: false},
				{Topic: "3dprinter/k1se/cfs/0/temperature", Payload: "0", Retain: false},
				{Topic: "3dprinter/k1se/cfs/0/state", Payload: "0", Retain: false},
			},
		},
		{
			name:      "no boxState key",
			input:     map[string]any{"other": "value"},
			baseTopic: "3dprinter/k1se",
			expected:  []types.MqttMessage{},
		},
		{
			name:      "empty input",
			input:     map[string]any{},
			baseTopic: "3dprinter/k1se",
			expected:  []types.MqttMessage{},
		},
		{
			name: "boxState is not a map",
			input: map[string]any{
				"boxState": "invalid",
			},
			baseTopic: "3dprinter/k1se",
			expected:  []types.MqttMessage{},
		},
		{
			name: "boxState missing id",
			input: map[string]any{
				"boxState": map[string]any{
					"humidity": 28.0,
					"temp":     23.0,
				},
			},
			baseTopic: "3dprinter/k1se",
			expected:  []types.MqttMessage{},
		},
		{
			name: "boxState with nil id",
			input: map[string]any{
				"boxState": map[string]any{
					"id":       nil,
					"humidity": 28.0,
				},
			},
			baseTopic: "3dprinter/k1se",
			expected: []types.MqttMessage{
				{Topic: "3dprinter/k1se/cfs/0/humidity", Payload: "28", Retain: false},
			},
		},
		{
			name: "boxState with invalid field types",
			input: map[string]any{
				"boxState": map[string]any{
					"id":       1,
					"humidity": "invalid",
					"temp":     []int{1, 2},
					"state":    "not_number",
				},
			},
			baseTopic: "3dprinter/k1se",
			expected: []types.MqttMessage{
				{Topic: "3dprinter/k1se/cfs/1/humidity", Payload: "invalid", Retain: false},
				{Topic: "3dprinter/k1se/cfs/1/temperature", Payload: "[1 2]", Retain: false},
				{Topic: "3dprinter/k1se/cfs/1/state", Payload: "0", Retain: false},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildCFSBoxMessages(tt.input, tt.baseTopic)

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

func TestToInt(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected int
	}{
		{
			name:     "int value",
			input:    42,
			expected: 42,
		},
		{
			name:     "int32 value",
			input:    int32(100),
			expected: 100,
		},
		{
			name:     "int64 value",
			input:    int64(200),
			expected: 200,
		},
		{
			name:     "float32 value",
			input:    float32(3.14),
			expected: 3,
		},
		{
			name:     "float64 value",
			input:    float64(2.99),
			expected: 2,
		},
		{
			name:     "zero int",
			input:    0,
			expected: 0,
		},
		{
			name:     "zero float",
			input:    0.0,
			expected: 0,
		},
		{
			name:     "negative int",
			input:    -5,
			expected: -5,
		},
		{
			name:     "negative float",
			input:    -3.7,
			expected: -3,
		},
		{
			name:     "string value",
			input:    "not_a_number",
			expected: 0,
		},
		{
			name:     "nil value",
			input:    nil,
			expected: 0,
		},
		{
			name:     "boolean value",
			input:    true,
			expected: 0,
		},
		{
			name:     "slice value",
			input:    []int{1, 2, 3},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toInt(tt.input)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}
