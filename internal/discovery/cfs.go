package discovery

import (
	"encoding/json"
	"fmt"

	"github.com/davidcollom/creality2mqtt/internal/types"
)

// BuildCFSBoxSensors builds HA discovery for a CFS box id (humidity, temperature)
func BuildCFSBoxSensors(cfg Config, device *Device, availTopic string, id int) []types.MqttMessage {
	messages := []types.MqttMessage{}

	// humidity sensor
	humUID := fmt.Sprintf("cfs_%d_humidity", id)
	humCfgTopic := fmt.Sprintf("%s/sensor/%s/%s/config", cfg.DiscoveryPrefix, cfg.DeviceID, humUID)
	humCfg := SensorConfig{
		Name:              fmt.Sprintf("CFS %d Humidity", id),
		UniqueID:          fmt.Sprintf("%s_%s", cfg.DeviceID, humUID),
		StateTopic:        fmt.Sprintf("%s/cfs/%d/humidity", cfg.BaseTopic, id),
		AvailabilityTopic: availTopic,
		PayloadAvailable:  "online",
		PayloadNotAvail:   "offline",
		UnitOfMeasurement: "%",
		DeviceClass:       "humidity",
		StateClass:        "measurement",
		Icon:              "mdi:water-percent",
		Device:            device,
	}
	humPayload, _ := json.Marshal(humCfg)
	messages = append(messages, types.MqttMessage{Topic: humCfgTopic, Payload: string(humPayload), Retain: true})

	// temperature sensor
	tempUID := fmt.Sprintf("cfs_%d_temperature", id)
	tempCfgTopic := fmt.Sprintf("%s/sensor/%s/%s/config", cfg.DiscoveryPrefix, cfg.DeviceID, tempUID)
	tempCfg := SensorConfig{
		Name:              fmt.Sprintf("CFS %d Temperature", id),
		UniqueID:          fmt.Sprintf("%s_%s", cfg.DeviceID, tempUID),
		StateTopic:        fmt.Sprintf("%s/cfs/%d/temperature", cfg.BaseTopic, id),
		AvailabilityTopic: availTopic,
		PayloadAvailable:  "online",
		PayloadNotAvail:   "offline",
		UnitOfMeasurement: "°C",
		DeviceClass:       "temperature",
		StateClass:        "measurement",
		Icon:              "mdi:thermometer",
		Device:            device,
	}
	tempPayload, _ := json.Marshal(tempCfg)
	messages = append(messages, types.MqttMessage{Topic: tempCfgTopic, Payload: string(tempPayload), Retain: true})

	return messages
}
