package discovery

import (
	"encoding/json"
	"testing"

	require "github.com/stretchr/testify/require"
)

func TestBuildCameraSensors(t *testing.T) {
	cfg := Config{DiscoveryPrefix: "ha", BaseTopic: "bt", DeviceID: "dev", PrinterIP: "10.0.0.5"}
	device := &Device{Identifiers: []string{"dev"}, Name: "Dev"}
	msgs := BuildCameraSensors(cfg, device, "bt/availability")
	require.Equal(t, true, len(msgs) >= 2)
	// First is discovery for sensor
	var sc SensorConfig
	_ = json.Unmarshal([]byte(msgs[0].Payload), &sc)
	require.Equal(t, true, sc.UniqueID != "")
	require.Equal(t, true, sc.StateTopic != "")
}
