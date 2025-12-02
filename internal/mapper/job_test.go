package mapper

import (
	"testing"

	"github.com/davidcollom/creality2mqtt/internal/types"
)

func TestBuildJobMessages(t *testing.T) {
	tests := []struct {
		name      string
		input     map[string]any
		baseTopic string
		want      []types.MqttMessage
	}{
		{
			name: "printing state - active job",
			input: map[string]any{
				"printProgress": 50,
				"printLeftTime": 1800,
				"printJobTime":  900,
				"layer":         25,
				"TotalLayer":    100,
				"printFileName": "/usr/data/printer_data/gcodes/test.gcode",
			},
			baseTopic: "3dprinter/k1se",
			want: []types.MqttMessage{
				{Topic: "3dprinter/k1se/printing", Payload: "true", Retain: false},
				{Topic: "3dprinter/k1se/job/progress", Payload: "50", Retain: false},
				{Topic: "3dprinter/k1se/job/left_time", Payload: "1800", Retain: false},
				{Topic: "3dprinter/k1se/job/job_time", Payload: "900", Retain: false},
				{Topic: "3dprinter/k1se/job/layer/current", Payload: "25", Retain: false},
				{Topic: "3dprinter/k1se/job/layer/total", Payload: "100", Retain: false},
				{Topic: "3dprinter/k1se/job/file_name", Payload: "test.gcode", Retain: false},
			},
		},
		{
			name: "not printing - zero progress",
			input: map[string]any{
				"printProgress": 0,
				"printLeftTime": 0,
			},
			baseTopic: "3dprinter/k1se",
			want: []types.MqttMessage{
				{Topic: "3dprinter/k1se/printing", Payload: "false", Retain: false},
				{Topic: "3dprinter/k1se/job/progress", Payload: "0", Retain: false},
				{Topic: "3dprinter/k1se/job/left_time", Payload: "0", Retain: false},
			},
		},
		{
			name: "not printing - missing left time",
			input: map[string]any{
				"printProgress": 50,
			},
			baseTopic: "3dprinter/k1se",
			want: []types.MqttMessage{
				{Topic: "3dprinter/k1se/printing", Payload: "false", Retain: false},
				{Topic: "3dprinter/k1se/job/progress", Payload: "50", Retain: false},
			},
		},
		{
			name: "string numeric values",
			input: map[string]any{
				"printProgress": "75",
				"printLeftTime": "600",
				"printJobTime":  "1200",
				"layer":         "50",
				"TotalLayer":    "80",
			},
			baseTopic: "test",
			want: []types.MqttMessage{
				{Topic: "test/printing", Payload: "true", Retain: false},
				{Topic: "test/job/progress", Payload: "75", Retain: false},
				{Topic: "test/job/left_time", Payload: "600", Retain: false},
				{Topic: "test/job/job_time", Payload: "1200", Retain: false},
				{Topic: "test/job/layer/current", Payload: "50", Retain: false},
				{Topic: "test/job/layer/total", Payload: "80", Retain: false},
			},
		},
		{
			name: "float64 values",
			input: map[string]any{
				"printProgress": 85.5,
				"printLeftTime": 300.0,
				"layer":         42.0,
			},
			baseTopic: "test",
			want: []types.MqttMessage{
				{Topic: "test/printing", Payload: "true", Retain: false},
				{Topic: "test/job/progress", Payload: "85", Retain: false},
				{Topic: "test/job/left_time", Payload: "300", Retain: false},
				{Topic: "test/job/layer/current", Payload: "42", Retain: false},
			},
		},
		{
			name: "zero total layers filtered out",
			input: map[string]any{
				"printProgress": 30,
				"printLeftTime": 1200,
				"TotalLayer":    0,
			},
			baseTopic: "test",
			want: []types.MqttMessage{
				{Topic: "test/printing", Payload: "true", Retain: false},
				{Topic: "test/job/progress", Payload: "30", Retain: false},
				{Topic: "test/job/left_time", Payload: "1200", Retain: false},
			},
		},
		{
			name: "complex filename path",
			input: map[string]any{
				"printProgress": 10,
				"printLeftTime": 3600,
				"printFileName": "/usr/data/printer_data/gcodes/.Meta+2+Stock+V3 (1)_gcode.3mf/Meta+2+Stock+V3_plate_4.gcode",
			},
			baseTopic: "test",
			want: []types.MqttMessage{
				{Topic: "test/printing", Payload: "true", Retain: false},
				{Topic: "test/job/progress", Payload: "10", Retain: false},
				{Topic: "test/job/left_time", Payload: "3600", Retain: false},
				{Topic: "test/job/file_name", Payload: "Meta+2+Stock+V3_plate_4.gcode", Retain: false},
			},
		},
		{
			name: "empty filename ignored",
			input: map[string]any{
				"printProgress": 20,
				"printLeftTime": 900,
				"printFileName": "",
			},
			baseTopic: "test",
			want: []types.MqttMessage{
				{Topic: "test/printing", Payload: "true", Retain: false},
				{Topic: "test/job/progress", Payload: "20", Retain: false},
				{Topic: "test/job/left_time", Payload: "900", Retain: false},
			},
		},
		{
			name: "invalid numeric strings ignored",
			input: map[string]any{
				"printProgress": "invalid",
				"printLeftTime": "not-a-number",
			},
			baseTopic: "test",
			want: []types.MqttMessage{
				{Topic: "test/printing", Payload: "false", Retain: false},
			},
		},
		{
			name:      "empty input",
			input:     map[string]any{},
			baseTopic: "test",
			want: []types.MqttMessage{
				{Topic: "test/printing", Payload: "false", Retain: false},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildJobMessages(tt.input, tt.baseTopic)

			if len(got) != len(tt.want) {
				t.Errorf("BuildJobMessages() returned %d messages, want %d", len(got), len(tt.want))
				return
			}

			for i, want := range tt.want {
				if got[i].Topic != want.Topic {
					t.Errorf("Message %d: topic = %q, want %q", i, got[i].Topic, want.Topic)
				}
				if got[i].Payload != want.Payload {
					t.Errorf("Message %d: payload = %q, want %q", i, got[i].Payload, want.Payload)
				}
				if got[i].Retain != want.Retain {
					t.Errorf("Message %d: retain = %v, want %v", i, got[i].Retain, want.Retain)
				}
			}
		})
	}
}

func TestGetInt(t *testing.T) {
	tests := []struct {
		name    string
		msg     map[string]any
		key     string
		wantVal int64
		wantOk  bool
	}{
		{
			name:    "int value",
			msg:     map[string]any{"test": 42},
			key:     "test",
			wantVal: 42,
			wantOk:  true,
		},
		{
			name:    "int32 value",
			msg:     map[string]any{"test": int32(123)},
			key:     "test",
			wantVal: 123,
			wantOk:  true,
		},
		{
			name:    "int64 value",
			msg:     map[string]any{"test": int64(456)},
			key:     "test",
			wantVal: 456,
			wantOk:  true,
		},
		{
			name:    "float64 value",
			msg:     map[string]any{"test": 78.9},
			key:     "test",
			wantVal: 78,
			wantOk:  true,
		},
		{
			name:    "string numeric",
			msg:     map[string]any{"test": "999"},
			key:     "test",
			wantVal: 999,
			wantOk:  true,
		},
		{
			name:    "string non-numeric",
			msg:     map[string]any{"test": "not-a-number"},
			key:     "test",
			wantVal: 0,
			wantOk:  false,
		},
		{
			name:    "missing key",
			msg:     map[string]any{"other": 123},
			key:     "test",
			wantVal: 0,
			wantOk:  false,
		},
		{
			name:    "unsupported type",
			msg:     map[string]any{"test": []int{1, 2, 3}},
			key:     "test",
			wantVal: 0,
			wantOk:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, gotOk := getInt(tt.msg, tt.key)
			if gotVal != tt.wantVal {
				t.Errorf("getInt() value = %v, want %v", gotVal, tt.wantVal)
			}
			if gotOk != tt.wantOk {
				t.Errorf("getInt() ok = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestSimplifyFileName(t *testing.T) {
	tests := []struct {
		name string
		full string
		want string
	}{
		{
			name: "simple filename",
			full: "test.gcode",
			want: "test.gcode",
		},
		{
			name: "full unix path",
			full: "/usr/data/printer_data/gcodes/model.gcode",
			want: "model.gcode",
		},
		{
			name: "complex nested path",
			full: "/usr/data/printer_data/gcodes/.Meta+2+Stock+V3 (1)_gcode.3mf/Meta+2+Stock+V3_plate_4.gcode",
			want: "Meta+2+Stock+V3_plate_4.gcode",
		},
		{
			name: "windows-style path",
			full: "C:\\Users\\Print\\Documents\\file.gcode",
			want: "file.gcode",
		},
		{
			name: "empty string",
			full: "",
			want: "",
		},
		{
			name: "whitespace string",
			full: "   ",
			want: "",
		},
		{
			name: "path with spaces",
			full: "/path/to/my file name.gcode",
			want: "my file name.gcode",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := simplifyFileName(tt.full)
			if got != tt.want {
				t.Errorf("simplifyFileName(%q) = %q, want %q", tt.full, got, tt.want)
			}
		})
	}
}
