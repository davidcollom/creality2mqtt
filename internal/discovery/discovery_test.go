package discovery

import (
	"encoding/json"
	"testing"

	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
)

func TestGenerateDiscoveryMessages(t *testing.T) {
	cfg := Config{
		DiscoveryPrefix: "homeassistant",
		BaseTopic:       "printer/test",
		DeviceID:        "dev1",
		DeviceName:      "Test Printer",
		DeviceModel:     "K1 SE",
		PrinterIP:       "192.168.1.10",
	}

	msgs := GenerateDiscoveryMessages(cfg)
	require.Equal(t, true, len(msgs) > 0)

	// Build table output
	t.Logf("Test Name | Description | Input | Expected Output | Actual Output | Pass/Fail")
	t.Logf("%s | %s | %s | %s | %d messages | %s", "GenerateDiscoveryMessages", "Aggregates discovery payloads", "cfg",
		"+sensor/+binary_sensor/+switch/+camera", len(msgs), "Pass")

	// Quick sanity on one known config: printer status
	found := false
	for _, m := range msgs {
		if m.Topic == "homeassistant/sensor/dev1/printer_status/config" {
			found = true
			var sc SensorConfig
			_ = json.Unmarshal([]byte(m.Payload), &sc)
			assert.Equal(t, "printer/test/printer_status", sc.StateTopic)
		}
	}
	require.Equal(t, true, found)
}
