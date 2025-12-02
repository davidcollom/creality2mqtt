package discovery

import (
	"encoding/json"
	"fmt"

	"github.com/davidcollom/creality2mqtt/internal/types"
)

// SwitchConfig represents a Home Assistant MQTT switch configuration
type SwitchConfig struct {
	Name              string  `json:"name"`
	UniqueID          string  `json:"unique_id"`
	StateTopic        string  `json:"state_topic"`
	CommandTopic      string  `json:"command_topic"`
	AvailabilityTopic string  `json:"availability_topic"`
	PayloadAvailable  string  `json:"payload_available"`
	PayloadNotAvail   string  `json:"payload_not_available"`
	PayloadOn         string  `json:"payload_on"`
	PayloadOff        string  `json:"payload_off"`
	StateOn           string  `json:"state_on"`
	StateOff          string  `json:"state_off"`
	Icon              string  `json:"icon,omitempty"`
	Device            *Device `json:"device"`
}

// BuildLightSwitch creates the light switch discovery message
// This allows bidirectional control - read state and send commands
func BuildLightSwitch(cfg Config, device *Device, availTopic string) []types.MqttMessage {
	topics := types.NewTopicBuilder(cfg.BaseTopic, cfg.DiscoveryPrefix)

	switchTopic := topics.Discovery("switch", cfg.DeviceID, "light")
	switchConfig := SwitchConfig{
		Name:              "Light",
		UniqueID:          fmt.Sprintf("%s_light", cfg.DeviceID),
		StateTopic:        topics.LightState(),
		CommandTopic:      topics.LightCommand(),
		AvailabilityTopic: availTopic,
		PayloadAvailable:  "online",
		PayloadNotAvail:   "offline",
		PayloadOn:         "1",
		PayloadOff:        "0",
		StateOn:           "1",
		StateOff:          "0",
		Icon:              "mdi:lightbulb",
		Device:            device,
	}

	payload, _ := json.Marshal(switchConfig)
	return []types.MqttMessage{{
		Topic:   switchTopic,
		Payload: string(payload),
		Retain:  true,
	}}
}
