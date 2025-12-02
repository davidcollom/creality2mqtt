package discovery

import (
	"encoding/json"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/davidcollom/creality2mqtt/internal/types"
)

// BuildCameraSensors creates camera-related sensor discovery messages
func BuildCameraSensors(cfg Config, device *Device, availTopic string) []types.MqttMessage {
	if cfg.PrinterIP == "" {
		return nil
	}

	messages := []types.MqttMessage{}
	topics := types.NewTopicBuilder(cfg.BaseTopic, cfg.DiscoveryPrefix)
	streamURL := fmt.Sprintf("http://%s:8080/?action=stream", cfg.PrinterIP)

	// Camera stream URL sensor
	streamURLTopic := topics.Discovery("sensor", cfg.DeviceID, "camera_stream_url")
	streamConfig := SensorConfig{
		Name:              "Camera Stream URL",
		UniqueID:          fmt.Sprintf("%s_camera_stream_url", cfg.DeviceID),
		StateTopic:        topics.CameraStreamURL(),
		AvailabilityTopic: availTopic,
		PayloadAvailable:  "online",
		PayloadNotAvail:   "offline",
		Icon:              "mdi:video",
		Device:            device,
	}
	payload, _ := json.Marshal(streamConfig)
	messages = append(messages, types.MqttMessage{
		Topic:   streamURLTopic,
		Payload: string(payload),
		Retain:  true,
	})

	// Publish the stream URL value once
	messages = append(messages, types.MqttMessage{
		Topic:   topics.CameraStreamURL(),
		Payload: streamURL,
		Retain:  true,
	})

	log.Info("Camera stream available - manual setup required",
		"stream_url", streamURL,
		"setup", "Add to configuration.yaml → camera: - platform: mjpeg, name: Creality Camera, mjpeg_url: "+streamURL)

	// Video stream active binary sensor
	videoTopic := fmt.Sprintf("%s/binary_sensor/%s/video_stream/config", cfg.DiscoveryPrefix, cfg.DeviceID)
	videoConfig := BinarySensorConfig{
		Name:              "Camera Stream Active",
		UniqueID:          fmt.Sprintf("%s_video_stream", cfg.DeviceID),
		StateTopic:        fmt.Sprintf("%s/video", cfg.BaseTopic),
		AvailabilityTopic: availTopic,
		PayloadAvailable:  "online",
		PayloadNotAvail:   "offline",
		PayloadOn:         "1",
		PayloadOff:        "0",
		Icon:              "mdi:video",
		Device:            device,
	}
	payload, _ = json.Marshal(videoConfig)
	messages = append(messages, types.MqttMessage{
		Topic:   videoTopic,
		Payload: string(payload),
		Retain:  true,
	})

	return messages
}
