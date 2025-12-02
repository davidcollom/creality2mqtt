package discovery

import (
	"encoding/json"
	"testing"

	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
)

func TestBuildLightSwitch(t *testing.T) {
	cfg := Config{DiscoveryPrefix: "ha", BaseTopic: "bt", DeviceID: "dev"}
	device := &Device{Identifiers: []string{"dev"}, Name: "Dev"}
	msgs := BuildLightSwitch(cfg, device, "bt/availability")
	assert.Equal(t, 1, len(msgs))
	require.Equal(t, true, msgs[0].Retain)
	var sc SwitchConfig
	_ = json.Unmarshal([]byte(msgs[0].Payload), &sc)
	require.Equal(t, true, sc.StateTopic != "")
	require.Equal(t, true, sc.CommandTopic != "")
	assert.Equal(t, "1", sc.PayloadOn)
	assert.Equal(t, "0", sc.PayloadOff)
}
