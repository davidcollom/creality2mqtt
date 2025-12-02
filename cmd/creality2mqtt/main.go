package main

import (
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	wsURL           string
	broker          string
	clientID        string
	baseTopic       string
	username        string
	password        string
	logLevel        string
	discoveryPrefix string
	deviceName      string
	mqttMinInterval time.Duration
)

// Create the rootCmd to attach everything else onto
var rootCmd = &cobra.Command{
	Use:   "creality2mqtt",
	Short: "Bridge between Creality printer WebSocket and MQTT",
	Long:  `A bridge application that connects to a Creality 3D printer via WebSocket and publishes data to an MQTT broker.`,
}

func init() {
	// Global flags shared across commands
	rootCmd.PersistentFlags().StringVar(&broker, "mqtt-broker", getEnvOrDefault("CREALITY_MQTT_BROKER", "tcp://localhost:1883"), "MQTT broker URL")
	rootCmd.PersistentFlags().StringVar(&clientID, "mqtt-client-id", getEnvOrDefault("CREALITY_MQTT_CLIENT_ID", "creality2mqtt"), "MQTT client ID")
	rootCmd.PersistentFlags().StringVar(&username, "mqtt-username", os.Getenv("CREALITY_MQTT_USERNAME"), "MQTT username")
	rootCmd.PersistentFlags().StringVar(&password, "mqtt-password", os.Getenv("CREALITY_MQTT_PASSWORD"), "MQTT password")
	rootCmd.PersistentFlags().StringVar(&discoveryPrefix, "discovery-prefix", getEnvOrDefault("CREALITY_DISCOVERY_PREFIX", "homeassistant"), "Home Assistant MQTT Discovery prefix")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", getEnvOrDefault("CREALITY_LOG_LEVEL", "info"), "Log level (debug, info, warn, error)")

	// Flags specific to the main run command
	rootCmd.PersistentFlags().StringVar(&wsURL, "ws-url", os.Getenv("CREALITY_WS_URL"), "WebSocket URL of printer (e.g. ws://192.168.1.50:9999/) [required]")
	rootCmd.PersistentFlags().StringVar(&baseTopic, "mqtt-base-topic", getEnvOrDefault("CREALITY_MQTT_BASE_TOPIC", "creality/printer"), "Base MQTT topic")
	rootCmd.PersistentFlags().StringVar(&deviceName, "device-name", os.Getenv("CREALITY_DEVICE_NAME"), "Device name override for Home Assistant")
	rootCmd.PersistentFlags().DurationVar(&mqttMinInterval, "mqtt-min-interval", getEnvOrDefaultDuration("CREALITY_MQTT_MIN_INTERVAL", 60*time.Second), "Minimum seconds between publishes per topic (0=disabled)")

	if wsURL == "" || os.Getenv("CREALITY_WS_URL") != "" {
		if err := rootCmd.MarkPersistentFlagRequired("ws-url"); err != nil {
			panic("Failed to mark ws-url flag as required: " + err.Error())
		}
	}

	// Add subcommands
	rootCmd.AddCommand(cleanupCmd)
	rootCmd.AddCommand(deviceInfoCmd)
	rootCmd.AddCommand(runCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
