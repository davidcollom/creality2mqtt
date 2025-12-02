package discovery

import (
	"fmt"

	"github.com/davidcollom/creality2mqtt/internal/types"
)

// CleanupOldEntities generates messages to remove old/unused entities from Home Assistant.
// This function publishes empty payloads to discovery topics, which tells Home Assistant
// to delete those entity configurations.
//
// HOW TO FIND OLD ENTITY IDs:
// 1. In Home Assistant, go to Settings → Devices & Services → MQTT
// 2. Click on your Creality printer device
// 3. Look for entities that appear as "Unavailable" or have strange/duplicate names
// 4. Click on an unwanted entity and look at the "Entity ID" (e.g., sensor.k1_se_old_name)
// 5. The unique_id used in discovery is typically the part after the device prefix
//
// MQTT TOPIC PATTERN:
// Discovery topics follow this pattern:
//
//	homeassistant/{component}/{device_id}/{unique_id}/config
//
// Where:
//   - component: sensor, binary_sensor, camera, etc.
//   - device_id: your printer's device ID (e.g., "k1_se_192_168_4_87")
//   - unique_id: the entity's unique identifier (what you add to the lists below)
//
// ALTERNATIVE METHOD - MQTT Explorer:
// Use MQTT Explorer or mosquitto_sub to see all retained discovery topics:
//
//	mosquitto_sub -h localhost -t "homeassistant/#" -v
//
// Look for topics with your device_id that shouldn't exist anymore.
func CleanupOldEntities(cfg Config) []types.MqttMessage {
	messages := []types.MqttMessage{}

	// Old regular sensors that should be removed
	// Add the unique_id portion of entities you want to delete
	// Example: if entity is "sensor.k1_se_printer_online", add "printer_online"
	oldSensors := []string{
		// Add deprecated sensor unique_ids here:
		"printer_online", // If you had this before switching to LWT
		// "printer_connected",   // If you had this before switching to LWT
		// "printer_connected_2", // If you had this before switching to LWT
		"old_temp_sensor",   // Example of old naming
		"camera_stream",     // If camera was previously a sensor
		"camera_stream_url", // If you had this before camera discovery
		"last_seen",         // If you had this before camera discovery
	}

	// Remove old sensors
	for _, sensorID := range oldSensors {
		topic := fmt.Sprintf("%s/sensor/%s/%s/config", cfg.DiscoveryPrefix, cfg.DeviceID, sensorID)
		messages = append(messages, types.MqttMessage{
			Topic:   topic,
			Payload: "", // Empty payload removes the entity
			Retain:  true,
		})
	}

	// Old binary sensors that should be removed
	oldBinarySensors := []string{
		// Add deprecated binary_sensor unique_ids here:
		"light",  // Light is now a switch, not binary_sensor
		"online", // Example: old availability sensor
		"connected",
		"camera_stream",       // If camera was previously a sensor
		"printer_connected",   // If you had this before switching to LWT
		"printer_connected_2", // If you had this before switching to LWT
		// "is_printing",             // Example: old naming variant
	}

	for _, sensorID := range oldBinarySensors {
		topic := fmt.Sprintf("%s/binary_sensor/%s/%s/config", cfg.DiscoveryPrefix, cfg.DeviceID, sensorID)
		messages = append(messages, types.MqttMessage{
			Topic:   topic,
			Payload: "",
			Retain:  true,
		})
	}

	// Old cameras that should be removed
	oldCameras := []string{
		// Add deprecated camera unique_ids here:
		// "camera",                  // If you tried camera discovery before
		// "video_stream",            // Example old camera config
	}

	for _, cameraID := range oldCameras {
		topic := fmt.Sprintf("%s/camera/%s/%s/config", cfg.DiscoveryPrefix, cfg.DeviceID, cameraID)
		messages = append(messages, types.MqttMessage{
			Topic:   topic,
			Payload: "",
			Retain:  true,
		})
	}

	return messages
}
