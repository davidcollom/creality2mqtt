package discovery

// Device represents the device information for Home Assistant
type Device struct {
	Identifiers  []string `json:"identifiers"`
	Name         string   `json:"name"`
	Manufacturer string   `json:"mf"`
	Model        string   `json:"mdl"`
}

// SensorConfig represents Home Assistant MQTT sensor discovery config
type SensorConfig struct {
	Name              string  `json:"name"`
	UniqueID          string  `json:"unique_id"`
	StateTopic        string  `json:"state_topic"`
	AvailabilityTopic string  `json:"availability_topic,omitempty"`
	PayloadAvailable  string  `json:"payload_available,omitempty"`
	PayloadNotAvail   string  `json:"payload_not_available,omitempty"`
	UnitOfMeasurement string  `json:"unit_of_measurement,omitempty"`
	DeviceClass       string  `json:"device_class,omitempty"`
	StateClass        string  `json:"state_class,omitempty"`
	Icon              string  `json:"icon,omitempty"`
	Device            *Device `json:"device"`
}

// BinarySensorConfig represents Home Assistant MQTT binary sensor discovery config
type BinarySensorConfig struct {
	Name              string  `json:"name"`
	UniqueID          string  `json:"unique_id"`
	StateTopic        string  `json:"state_topic"`
	AvailabilityTopic string  `json:"availability_topic,omitempty"`
	PayloadAvailable  string  `json:"payload_available,omitempty"`
	PayloadNotAvail   string  `json:"payload_not_available,omitempty"`
	PayloadOn         string  `json:"payload_on"`
	PayloadOff        string  `json:"payload_off"`
	DeviceClass       string  `json:"device_class,omitempty"`
	Icon              string  `json:"icon,omitempty"`
	Device            *Device `json:"device"`
}

// CameraConfig represents Home Assistant MQTT camera discovery config
type CameraConfig struct {
	Name              string  `json:"name"`
	UniqueID          string  `json:"unique_id"`
	Topic             string  `json:"topic"`
	ImageTopic        string  `json:"image_topic,omitempty"`
	AvailabilityTopic string  `json:"availability_topic,omitempty"`
	PayloadAvailable  string  `json:"payload_available,omitempty"`
	PayloadNotAvail   string  `json:"payload_not_available,omitempty"`
	Device            *Device `json:"device"`
}

// Config holds discovery configuration
type Config struct {
	DiscoveryPrefix string
	BaseTopic       string
	DeviceID        string
	DeviceName      string
	DeviceModel     string
	PrinterIP       string // IP address for camera stream
}
