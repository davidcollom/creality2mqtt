package types

import "fmt"

// MqttMessage represents a single MQTT publish operation.
type MqttMessage struct {
	Topic   string
	Payload string
	Retain  bool
}

// TopicBuilder provides consistent topic construction across the application
type TopicBuilder struct {
	BaseTopic       string
	DiscoveryPrefix string
}

// NewTopicBuilder creates a new TopicBuilder instance
func NewTopicBuilder(baseTopic, discoveryPrefix string) *TopicBuilder {
	return &TopicBuilder{
		BaseTopic:       baseTopic,
		DiscoveryPrefix: discoveryPrefix,
	}
}

// Availability returns the availability/status topic (used for LWT)
func (tb *TopicBuilder) Availability() string {
	return fmt.Sprintf("%s/status", tb.BaseTopic)
}

// HAStatus returns the Home Assistant status topic
func (tb *TopicBuilder) HAStatus() string {
	return fmt.Sprintf("%s/status", tb.DiscoveryPrefix)
}

// LightState returns the light state topic
func (tb *TopicBuilder) LightState() string {
	return fmt.Sprintf("%s/light_sw", tb.BaseTopic)
}

// LightCommand returns the light command topic
func (tb *TopicBuilder) LightCommand() string {
	return fmt.Sprintf("%s/light_sw/set", tb.BaseTopic)
}

// CameraStreamURL returns the camera stream URL topic
func (tb *TopicBuilder) CameraStreamURL() string {
	return fmt.Sprintf("%s/camera_stream_url", tb.BaseTopic)
}

// Discovery returns a discovery config topic for a component and entity
func (tb *TopicBuilder) Discovery(component, deviceID, entityID string) string {
	return fmt.Sprintf("%s/%s/%s/%s/config", tb.DiscoveryPrefix, component, deviceID, entityID)
}

// Data returns a data topic under the base topic
func (tb *TopicBuilder) Data(subtopic string) string {
	return fmt.Sprintf("%s/%s", tb.BaseTopic, subtopic)
}
