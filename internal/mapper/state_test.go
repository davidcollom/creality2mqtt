package mapper

import (
	"testing"
	"time"
)

func TestBuildStateMessages(t *testing.T) {
	t.Parallel()

	baseTopic := "3dprinter/k1se"

	tests := []struct {
		name     string
		input    map[string]any
		expected map[string]string
	}{
		{
			name: "printer active and tfCard present",
			input: map[string]any{
				"printProgress": 10,
				"leftTime":      100,
				"tfCard":        1,
			},
			expected: map[string]string{
				baseTopic + "/printer_status":  "active",
				baseTopic + "/tf_card_present": "true",
			},
		},
		{
			name: "idle, no tfCard",
			input: map[string]any{
				"deviceState": 0,
			},
			expected: map[string]string{
				baseTopic + "/printer_status": "idle",
			},
		},
		{
			name: "idle with tfCard false",
			input: map[string]any{
				"tfCard": 0,
			},
			expected: map[string]string{
				baseTopic + "/printer_status":  "idle",
				baseTopic + "/tf_card_present": "false",
			},
		},
		{
			name: "no status fields, tfCard true",
			input: map[string]any{
				"tfCard": 1,
			},
			expected: map[string]string{
				baseTopic + "/printer_status":  "idle",
				baseTopic + "/tf_card_present": "true",
			},
		},
		{
			name: "active when printing heuristics",
			input: map[string]any{
				"printProgress": 50,
				"leftTime":      200,
			},
			expected: map[string]string{
				baseTopic + "/printer_status": "active",
			},
		},
		{
			name: "active via gcodeState",
			input: map[string]any{
				"gcodeState": 1,
			},
			expected: map[string]string{
				baseTopic + "/printer_status": "active",
			},
		},
		{
			name:  "empty message",
			input: map[string]any{},
			expected: map[string]string{
				baseTopic + "/printer_status": "idle",
			},
		},
		{
			name: "string values ignored",
			input: map[string]any{
				"connect":     "1",
				"deviceState": "invalid",
				"tfCard":      "true",
			},
			expected: map[string]string{
				baseTopic + "/printer_status": "idle",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset rate-limit cache to ensure first publish per test case
			printerStatusCache.mu.Lock()
			printerStatusCache.lastStatus = ""
			printerStatusCache.lastPublished = time.Time{}
			printerStatusCache.lastUpdate = time.Time{}
			printerStatusCache.mu.Unlock()

			got := BuildStateMessages(tt.input, baseTopic)
			tp := toTopicMap(got)

			// Check that we got the expected number of topics
			if len(tp) != len(tt.expected) {
				t.Errorf("expected %d topics, got %d", len(tt.expected), len(tp))
			}

			// Verify each expected topic and payload
			for topic, want := range tt.expected {
				gotVal, ok := tp[topic]
				if !ok {
					t.Errorf("expected topic %s to be present", topic)
					continue
				}
				if gotVal != want {
					t.Errorf("topic %s payload = %q, want %q", topic, gotVal, want)
				}
			}

			// Ensure no unexpected topics were generated beyond printer_status/tf_card_present
			for topic := range tp {
				if _, expected := tt.expected[topic]; !expected {
					t.Errorf("unexpected topic generated: %s", topic)
				}
			}
		})
	}
}

// func TestDeriveOnline(t *testing.T) {
// 	t.Parallel()

// 	tests := []struct {
// 		name     string
// 		input    map[string]any
// 		expected bool
// 	}{
// 		{
// 			name: "online via connect=1",
// 			input: map[string]any{
// 				"connect": 1,
// 			},
// 			expected: true,
// 		},
// 		{
// 			name: "online via deviceState=5",
// 			input: map[string]any{
// 				"deviceState": 5,
// 			},
// 			expected: true,
// 		},
// 		{
// 			name: "online via both connect and deviceState",
// 			input: map[string]any{
// 				"connect":     1,
// 				"deviceState": 2,
// 			},
// 			expected: true,
// 		},
// 		{
// 			name: "offline - connect=0, deviceState=0",
// 			input: map[string]any{
// 				"connect":     0,
// 				"deviceState": 0,
// 			},
// 			expected: false,
// 		},
// 		{
// 			name:     "offline - no fields present",
// 			input:    map[string]any{},
// 			expected: false,
// 		},
// 		{
// 			name: "offline - invalid types",
// 			input: map[string]any{
// 				"connect":     "1",
// 				"deviceState": "active",
// 			},
// 			expected: false,
// 		},
// 		{
// 			name: "connect=1 overrides deviceState=0",
// 			input: map[string]any{
// 				"connect":     1,
// 				"deviceState": 0,
// 			},
// 			expected: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := deriveOnline(tt.input)
// 			if got != tt.expected {
// 				t.Errorf("deriveOnline() = %v, want %v", got, tt.expected)
// 			}
// 		})
// 	}
// }
