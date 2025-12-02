package discovery

import (
	"encoding/json"
	"fmt"

	"github.com/davidcollom/creality2mqtt/internal/types"
)

// BuildTemperatureSensors creates temperature sensor discovery messages
func BuildTemperatureSensors(cfg Config, device *Device, availTopic string) []types.MqttMessage {
	messages := []types.MqttMessage{}

	tempSensors := []struct {
		name       string
		stateTopic string
		uniqueID   string
	}{
		{"Nozzle Temperature", fmt.Sprintf("%s/temperature/nozzle/current", cfg.BaseTopic), "nozzle_temp_current"},
		{"Nozzle Target Temperature", fmt.Sprintf("%s/temperature/nozzle/target", cfg.BaseTopic), "nozzle_temp_target"},
		{"Bed Temperature", fmt.Sprintf("%s/temperature/bed0/current", cfg.BaseTopic), "bed_temp_current"},
		{"Bed Target Temperature", fmt.Sprintf("%s/temperature/bed0/target", cfg.BaseTopic), "bed_temp_target"},
	}

	for _, ts := range tempSensors {
		configTopic := fmt.Sprintf("%s/sensor/%s/%s/config", cfg.DiscoveryPrefix, cfg.DeviceID, ts.uniqueID)
		config := SensorConfig{
			Name:              ts.name,
			UniqueID:          fmt.Sprintf("%s_%s", cfg.DeviceID, ts.uniqueID),
			StateTopic:        ts.stateTopic,
			AvailabilityTopic: availTopic,
			PayloadAvailable:  "online",
			PayloadNotAvail:   "offline",
			UnitOfMeasurement: "°C",
			DeviceClass:       "temperature",
			StateClass:        "measurement",
			Icon:              "mdi:thermometer",
			Device:            device,
		}

		payload, _ := json.Marshal(config)
		messages = append(messages, types.MqttMessage{
			Topic:   configTopic,
			Payload: string(payload),
			Retain:  true,
		})
	}

	return messages
}

// BuildFeedStateSensor creates the extruder feed state sensor
func BuildFeedStateSensor(cfg Config, device *Device, availTopic string) []types.MqttMessage {
	messages := []types.MqttMessage{}

	configTopic := fmt.Sprintf("%s/sensor/%s/feed_state/config", cfg.DiscoveryPrefix, cfg.DeviceID)
	config := SensorConfig{
		Name:              "Feed State",
		UniqueID:          fmt.Sprintf("%s_feed_state", cfg.DeviceID),
		StateTopic:        fmt.Sprintf("%s/feed_state", cfg.BaseTopic),
		AvailabilityTopic: availTopic,
		PayloadAvailable:  "online",
		PayloadNotAvail:   "offline",
		Icon:              "mdi:printer-3d-nozzle",
		Device:            device,
	}

	payload, _ := json.Marshal(config)
	messages = append(messages, types.MqttMessage{
		Topic:   configTopic,
		Payload: string(payload),
		Retain:  true,
	})

	return messages
}

// BuildStatusSensor creates printer status sensor (idle/active)
func BuildStatusSensor(cfg Config, device *Device, availTopic string) []types.MqttMessage {
	messages := []types.MqttMessage{}

	configTopic := fmt.Sprintf("%s/sensor/%s/printer_status/config", cfg.DiscoveryPrefix, cfg.DeviceID)
	config := SensorConfig{
		Name:              "Printer Status",
		UniqueID:          fmt.Sprintf("%s_printer_status", cfg.DeviceID),
		StateTopic:        fmt.Sprintf("%s/printer_status", cfg.BaseTopic),
		AvailabilityTopic: availTopic,
		PayloadAvailable:  "online",
		PayloadNotAvail:   "offline",
		Icon:              "mdi:printer-3d",
		Device:            device,
	}

	payload, _ := json.Marshal(config)
	messages = append(messages, types.MqttMessage{
		Topic:   configTopic,
		Payload: string(payload),
		Retain:  true,
	})

	return messages
}

// BuildFanSensors creates fan speed sensor discovery messages
func BuildFanSensors(cfg Config, device *Device, availTopic string) []types.MqttMessage {
	messages := []types.MqttMessage{}

	fanSensors := []struct {
		name       string
		stateTopic string
		uniqueID   string
		icon       string
	}{
		{"Model Fan Speed", fmt.Sprintf("%s/model_fan_pct", cfg.BaseTopic), "model_fan_pct", "mdi:fan"},
		{"Auxiliary Fan Speed", fmt.Sprintf("%s/auxiliary_fan_pct", cfg.BaseTopic), "auxiliary_fan_pct", "mdi:fan"},
		{"Case Fan Speed", fmt.Sprintf("%s/case_fan_pct", cfg.BaseTopic), "case_fan_pct", "mdi:fan"},
	}

	for _, fs := range fanSensors {
		configTopic := fmt.Sprintf("%s/sensor/%s/%s/config", cfg.DiscoveryPrefix, cfg.DeviceID, fs.uniqueID)
		config := SensorConfig{
			Name:              fs.name,
			UniqueID:          fmt.Sprintf("%s_%s", cfg.DeviceID, fs.uniqueID),
			StateTopic:        fs.stateTopic,
			AvailabilityTopic: availTopic,
			PayloadAvailable:  "online",
			PayloadNotAvail:   "offline",
			UnitOfMeasurement: "%",
			Icon:              fs.icon,
			StateClass:        "measurement",
			Device:            device,
		}

		payload, _ := json.Marshal(config)
		messages = append(messages, types.MqttMessage{
			Topic:   configTopic,
			Payload: string(payload),
			Retain:  true,
		})
	}

	return messages
}

// BuildProgressSensor creates the print progress sensor discovery message
func BuildProgressSensor(cfg Config, device *Device, availTopic string) []types.MqttMessage {
	progressTopic := fmt.Sprintf("%s/sensor/%s/print_progress/config", cfg.DiscoveryPrefix, cfg.DeviceID)
	progressConfig := SensorConfig{
		Name:              "Print Progress",
		UniqueID:          fmt.Sprintf("%s_print_progress", cfg.DeviceID),
		StateTopic:        fmt.Sprintf("%s/job/progress", cfg.BaseTopic),
		AvailabilityTopic: availTopic,
		PayloadAvailable:  "online",
		PayloadNotAvail:   "offline",
		UnitOfMeasurement: "%",
		Icon:              "mdi:percent",
		StateClass:        "measurement",
		Device:            device,
	}

	payload, _ := json.Marshal(progressConfig)
	return []types.MqttMessage{{
		Topic:   progressTopic,
		Payload: string(payload),
		Retain:  true,
	}}
}
