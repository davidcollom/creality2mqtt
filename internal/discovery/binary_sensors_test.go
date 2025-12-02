package discovery

import (
	"encoding/json"
	"testing"

	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
)

func TestBuildPrintingSensor(t *testing.T) {
	cfg := Config{DiscoveryPrefix: "ha", BaseTopic: "bt", DeviceID: "dev"}
	device := &Device{Identifiers: []string{"dev"}, Name: "Dev"}
	msgs := BuildPrintingSensor(cfg, device, "bt/availability")
	assert.Equal(t, 1, len(msgs))
	require.Equal(t, true, msgs[0].Retain)
	var bc BinarySensorConfig
	_ = json.Unmarshal([]byte(msgs[0].Payload), &bc)
	assert.Equal(t, "bt/printing", bc.StateTopic)
	assert.Equal(t, "true", bc.PayloadOn)
	assert.Equal(t, "false", bc.PayloadOff)
}

func TestBuildPartFanSensor(t *testing.T) {
	cfg := Config{DiscoveryPrefix: "ha", BaseTopic: "bt", DeviceID: "dev"}
	device := &Device{Identifiers: []string{"dev"}, Name: "Dev"}
	msgs := BuildPartFanSensor(cfg, device, "bt/availability")
	assert.Equal(t, 1, len(msgs))
	var bc BinarySensorConfig
	_ = json.Unmarshal([]byte(msgs[0].Payload), &bc)
	assert.Equal(t, "bt/fan", bc.StateTopic)
	assert.Equal(t, "1", bc.PayloadOn)
	assert.Equal(t, "0", bc.PayloadOff)
}
