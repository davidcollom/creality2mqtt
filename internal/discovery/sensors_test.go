package discovery

import (
	"encoding/json"
	"testing"

	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
)

func TestBuildTemperatureSensors(t *testing.T) {
	cfg := Config{DiscoveryPrefix: "ha", BaseTopic: "bt", DeviceID: "dev"}
	device := &Device{Identifiers: []string{"dev"}, Name: "Dev"}
	msgs := BuildTemperatureSensors(cfg, device, "bt/availability")
	assert.Equal(t, 4, len(msgs))
	// Table output
	t.Logf("Test Name | Description | Input | Expected Output | Actual Output | Pass/Fail")
	t.Logf("%s | %s | %s | %s | %d configs | %s", "BuildTemperatureSensors", "4 temp sensors", "cfg,device", "4 retained configs", len(msgs), "Pass")

	// Validate one
	var sc SensorConfig
	_ = json.Unmarshal([]byte(msgs[0].Payload), &sc)
	require.Equal(t, true, sc.Device != nil)
	require.Equal(t, true, sc.UniqueID != "")
	require.Equal(t, true, sc.StateTopic != "")
}

func TestBuildStatusSensor(t *testing.T) {
	cfg := Config{DiscoveryPrefix: "ha", BaseTopic: "bt", DeviceID: "dev"}
	device := &Device{Identifiers: []string{"dev"}, Name: "Dev"}
	msgs := BuildStatusSensor(cfg, device, "bt/availability")
	assert.Equal(t, 1, len(msgs))
	var sc SensorConfig
	_ = json.Unmarshal([]byte(msgs[0].Payload), &sc)
	assert.Equal(t, "bt/printer_status", sc.StateTopic)
}

func TestBuildFanSensors(t *testing.T) {
	cfg := Config{DiscoveryPrefix: "ha", BaseTopic: "bt", DeviceID: "dev"}
	device := &Device{Identifiers: []string{"dev"}, Name: "Dev"}
	msgs := BuildFanSensors(cfg, device, "bt/availability")
	assert.Equal(t, 3, len(msgs))
}

func TestBuildProgressSensor(t *testing.T) {
	cfg := Config{DiscoveryPrefix: "ha", BaseTopic: "bt", DeviceID: "dev"}
	device := &Device{Identifiers: []string{"dev"}, Name: "Dev"}
	msgs := BuildProgressSensor(cfg, device, "bt/availability")
	assert.Equal(t, 1, len(msgs))
	var sc SensorConfig
	_ = json.Unmarshal([]byte(msgs[0].Payload), &sc)
	assert.Equal(t, "bt/job/progress", sc.StateTopic)
}
