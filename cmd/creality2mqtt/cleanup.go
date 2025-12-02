package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/cobra"
)

var cleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Remove all MQTT discovery entities for a device",
	Long: `Removes all MQTT discovery configurations for a Creality printer device from Home Assistant.
This will delete all entities (sensors, switches, etc.) associated with the device.
Use this to start fresh or remove orphaned entities.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		deviceID, _ := cmd.Flags().GetString("device-id")
		if deviceID == "" {
			log.Error("--device-id is required for cleanup command")
			return fmt.Errorf("device-id is required")
		}

		log.Info("Starting cleanup",
			"device_id", deviceID,
			"mqtt_broker", broker,
			"discovery_prefix", discoveryPrefix,
		)

		// Create MQTT client
		opts := mqtt.NewClientOptions()
		opts.AddBroker(broker)
		opts.SetClientID(clientID + "_cleanup")
		if username != "" {
			opts.SetUsername(username)
		}
		if password != "" {
			opts.SetPassword(password)
		}
		opts.SetCleanSession(true)
		opts.SetAutoReconnect(false)

		client := mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			log.Error("Failed to connect to MQTT broker", "error", token.Error())
			return token.Error()
		}
		defer client.Disconnect(250)

		log.Info("Connected to MQTT broker")

		// List of all entity types and their possible unique IDs
		components := []string{"sensor", "binary_sensor", "switch", "camera"}

		// All possible entity unique IDs (based on current and past implementations)
		entityIDs := []string{
			// Temperature sensors
			"nozzle_temp_current",
			"nozzle_temp_target",
			"bed_temp_current",
			"bed_temp_target",

			// Status sensors
			"printer_status",

			// Fan sensors
			"model_fan_pct",
			"auxiliary_fan_pct",
			"case_fan_pct",

			// Progress and state
			"print_progress",
			"feed_state",
			"last_seen",

			// Camera sensors
			"camera_stream_url",
			"video_stream",

			// Binary sensors
			"printing",
			"light",
			"part_fan",

			// Switches
			"light",

			// Old/deprecated entities
			"printer_online",
			"printer_connected",
			"printer_connected_2",
			"online",
			"camera_stream",
			"camera",
			"old_temp_sensor",
		}

		totalDeleted := 0

		// Delete all discovery topics for each component/entity combination
		for _, component := range components {
			for _, entityID := range entityIDs {
				topic := fmt.Sprintf("%s/%s/%s/%s/config", discoveryPrefix, component, deviceID, entityID)

				// Publish empty payload to delete the entity
				token := client.Publish(topic, 0, true, "")
				token.Wait()

				if token.Error() != nil {
					log.Error("Failed to delete entity", "component", component, "entity", entityID, "error", token.Error())
				} else {
					log.Debug("Deleted entity", "component", component, "entity", entityID, "topic", topic)
					totalDeleted++
				}
			}
		}

		// Small delay to ensure messages are sent
		time.Sleep(500 * time.Millisecond)

		log.Info("Cleanup complete", "deleted_count", totalDeleted)
		log.Info("The device should disappear from Home Assistant within a few seconds")

		return nil
	},
}

func init() {
	// Device ID flag
	cleanupCmd.Flags().String("device-id", "", "Device ID to remove (e.g., k1_se_192_168_4_87) [required]")
	if err := cleanupCmd.MarkFlagRequired("device-id"); err != nil {
		panic("Failed to mark device-id flag as required: " + err.Error())
	}
}
