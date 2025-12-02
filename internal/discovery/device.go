package discovery

import "strings"

// ExtractDeviceInfo tries to extract device information from the first WebSocket message
func ExtractDeviceInfo(msg map[string]any) (deviceID, deviceName, deviceModel string) {
	// Try to find device-specific fields in the message
	if v, ok := msg["deviceName"].(string); ok {
		deviceName = v
	}
	if v, ok := msg["deviceModel"].(string); ok {
		deviceModel = v
	}
	if v, ok := msg["deviceId"].(string); ok {
		deviceID = v
	}
	if v, ok := msg["device_name"].(string); ok && deviceName == "" {
		deviceName = v
	}
	if v, ok := msg["device_model"].(string); ok && deviceModel == "" {
		deviceModel = v
	}
	if v, ok := msg["device_id"].(string); ok && deviceID == "" {
		deviceID = v
	}

	// Fallback values
	if deviceID == "" {
		deviceID = "creality_printer"
	}
	if deviceName == "" {
		deviceName = "Creality Printer"
	}
	if deviceModel == "" {
		deviceModel = "K1/K1 SE/K1 Max"
	}

	// Sanitize deviceID for MQTT topics
	deviceID = strings.ToLower(deviceID)
	deviceID = strings.ReplaceAll(deviceID, " ", "_")
	deviceID = strings.ReplaceAll(deviceID, "-", "_")

	return
}
