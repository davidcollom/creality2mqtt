package discovery

import (
	"encoding/json"
	"testing"

	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
)

func TestBuildCFSBoxSensors(t *testing.T) {
	cfg := Config{DiscoveryPrefix: "ha", BaseTopic: "bt", DeviceID: "dev"}
	device := &Device{Identifiers: []string{"dev"}, Name: "Dev"}
	msgs := BuildCFSBoxSensors(cfg, device, "bt/availability", 1)
	assert.Equal(t, 2, len(msgs))
	// Validate humidity
	var hum SensorConfig
	_ = json.Unmarshal([]byte(msgs[0].Payload), &hum)
	assert.Equal(t, "humidity", hum.DeviceClass)
	require.Equal(t, true, hum.StateTopic != "")
}
