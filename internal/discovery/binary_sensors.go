package discovery

import (
	"encoding/json"
	"fmt"

	"github.com/davidcollom/creality2mqtt/internal/types"
)

// BuildPrintingSensor creates the printing binary sensor discovery message
func BuildPrintingSensor(cfg Config, device *Device, availTopic string) []types.MqttMessage {
	printingTopic := fmt.Sprintf("%s/binary_sensor/%s/printing/config", cfg.DiscoveryPrefix, cfg.DeviceID)
	printingConfig := BinarySensorConfig{
		Name:              "Printing",
		UniqueID:          fmt.Sprintf("%s_printing", cfg.DeviceID),
		StateTopic:        fmt.Sprintf("%s/printing", cfg.BaseTopic),
		AvailabilityTopic: availTopic,
		PayloadAvailable:  "online",
		PayloadNotAvail:   "offline",
		PayloadOn:         "true",
		PayloadOff:        "false",
		Icon:              "mdi:printer-3d",
		Device:            device,
	}

	payload, _ := json.Marshal(printingConfig)
	return []types.MqttMessage{{
		Topic:   printingTopic,
		Payload: string(payload),
		Retain:  true,
	}}
}

// BuildLightSensor creates the light binary sensor discovery message
func BuildLightSensor(cfg Config, device *Device, availTopic string) []types.MqttMessage {
	lightTopic := fmt.Sprintf("%s/binary_sensor/%s/light/config", cfg.DiscoveryPrefix, cfg.DeviceID)
	lightConfig := BinarySensorConfig{
		Name:              "Light",
		UniqueID:          fmt.Sprintf("%s_light", cfg.DeviceID),
		StateTopic:        fmt.Sprintf("%s/light_sw", cfg.BaseTopic),
		AvailabilityTopic: availTopic,
		PayloadAvailable:  "online",
		PayloadNotAvail:   "offline",
		PayloadOn:         "1",
		PayloadOff:        "0",
		Icon:              "mdi:lightbulb",
		Device:            device,
	}

	payload, _ := json.Marshal(lightConfig)
	return []types.MqttMessage{{
		Topic:   lightTopic,
		Payload: string(payload),
		Retain:  true,
	}}
}

// BuildPartFanSensor creates the part cooling fan binary sensor discovery message
func BuildPartFanSensor(cfg Config, device *Device, availTopic string) []types.MqttMessage {
	partFanTopic := fmt.Sprintf("%s/binary_sensor/%s/part_fan/config", cfg.DiscoveryPrefix, cfg.DeviceID)
	partFanConfig := BinarySensorConfig{
		Name:              "Part Cooling Fan",
		UniqueID:          fmt.Sprintf("%s_part_fan", cfg.DeviceID),
		StateTopic:        fmt.Sprintf("%s/fan", cfg.BaseTopic),
		AvailabilityTopic: availTopic,
		PayloadAvailable:  "online",
		PayloadNotAvail:   "offline",
		PayloadOn:         "1",
		PayloadOff:        "0",
		Icon:              "mdi:fan",
		Device:            device,
	}

	payload, _ := json.Marshal(partFanConfig)
	return []types.MqttMessage{{
		Topic:   partFanTopic,
		Payload: string(payload),
		Retain:  true,
	}}
}
