package discovery

import (
	"github.com/charmbracelet/log"
	"github.com/davidcollom/creality2mqtt/internal/types"
)

var discoverMessages []types.MqttMessage

// GenerateDiscoveryMessages creates Home Assistant MQTT Discovery messages
func GenerateDiscoveryMessages(cfg Config) []types.MqttMessage {
	device := &Device{
		Identifiers:  []string{cfg.DeviceID},
		Name:         cfg.DeviceName,
		Manufacturer: "Creality",
		Model:        cfg.DeviceModel,
	}

	// Build availability topic (consistent with LWT)
	topics := types.NewTopicBuilder(cfg.BaseTopic, cfg.DiscoveryPrefix)
	availTopic := topics.Availability()

	// Build all sensor discovery messages
	discoverMessages = append(discoverMessages, BuildTemperatureSensors(cfg, device, availTopic)...)
	discoverMessages = append(discoverMessages, BuildStatusSensor(cfg, device, availTopic)...)
	discoverMessages = append(discoverMessages, BuildFanSensors(cfg, device, availTopic)...)
	discoverMessages = append(discoverMessages, BuildProgressSensor(cfg, device, availTopic)...)
	discoverMessages = append(discoverMessages, BuildFeedStateSensor(cfg, device, availTopic)...)

	// Build binary sensor discovery messages
	discoverMessages = append(discoverMessages, BuildPrintingSensor(cfg, device, availTopic)...)
	discoverMessages = append(discoverMessages, BuildPartFanSensor(cfg, device, availTopic)...)

	// Build switch discovery messages (bidirectional control)
	discoverMessages = append(discoverMessages, BuildLightSwitch(cfg, device, availTopic)...)

	// Build camera-related discovery messages
	discoverMessages = append(discoverMessages, BuildCameraSensors(cfg, device, availTopic)...)

	log.Info("Generated MQTT Discovery messages", "count", len(discoverMessages))
	return discoverMessages
}
