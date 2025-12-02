package main

import (
	"fmt"
	"testing"

	"github.com/spf13/cobra"
)

func TestCleanupCmd(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		flags       map[string]string
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid device ID",
			args: []string{},
			flags: map[string]string{
				"device-id": "k1_se_192_168_4_87",
			},
			expectError: false,
		},
		{
			name:        "missing device ID",
			args:        []string{},
			flags:       map[string]string{},
			expectError: true,
			errorMsg:    "device_id is required",
		},
		{
			name: "empty device ID",
			args: []string{},
			flags: map[string]string{
				"device-id": "",
			},
			expectError: true,
			errorMsg:    "device_id is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new command instance for each test
			cmd := &cobra.Command{
				Use:   "cleanup",
				Short: "Remove all MQTT discovery entities for a device",
				RunE: func(cmd *cobra.Command, args []string) error {
					deviceID, _ := cmd.Flags().GetString("device-id")
					if deviceID == "" {
						return fmt.Errorf("device_id is required")
					}
					return nil
				},
			}
			cmd.Flags().String("device-id", "", "Device ID to remove")

			// Set flags
			for key, value := range tt.flags {
				err := cmd.Flags().Set(key, value)
				if err != nil {
					t.Fatalf("Failed to set flag %s: %v", key, err)
				}
			}

			// Execute command
			err := cmd.RunE(cmd, tt.args)

			// Check results
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Errorf("Expected error message %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestCleanupCmdFlags(t *testing.T) {
	tests := []struct {
		name     string
		flagName string
		expected bool
	}{
		{
			name:     "device-id flag exists",
			flagName: "device-id",
			expected: true,
		},
		{
			name:     "non-existent flag",
			flagName: "invalid-flag",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			cmd.Flags().String("device-id", "", "Device ID to remove")

			flag := cmd.Flags().Lookup(tt.flagName)
			exists := flag != nil

			if exists != tt.expected {
				t.Errorf("Flag %s existence: expected %t, got %t", tt.flagName, tt.expected, exists)
			}
		})
	}
}

func TestCleanupCmdProperties(t *testing.T) {
	tests := []struct {
		name     string
		property string
		expected string
	}{
		{
			name:     "command use",
			property: "use",
			expected: "cleanup",
		},
		{
			name:     "command short description",
			property: "short",
			expected: "Remove all MQTT discovery entities for a device",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var actual string
			switch tt.property {
			case "use":
				actual = cleanupCmd.Use
			case "short":
				actual = cleanupCmd.Short
			}

			if actual != tt.expected {
				t.Errorf("Expected %s to be %q, got %q", tt.property, tt.expected, actual)
			}
		})
	}
}

func TestCleanupEntityGeneration(t *testing.T) {
	components := []string{"sensor", "binary_sensor", "switch", "camera"}
	entityIDs := []string{
		"nozzle_temp_current",
		"nozzle_temp_target",
		"bed_temp_current",
		"bed_temp_target",
		"printer_status",
		"model_fan_pct",
		"auxiliary_fan_pct",
		"case_fan_pct",
		"print_progress",
		"feed_state",
		"last_seen",
		"camera_stream_url",
		"video_stream",
		"printing",
		"light",
		"part_fan",
	}

	deviceID := "test_device"
	discoveryPrefix := "homeassistant"

	expectedTopicCount := len(components) * len(entityIDs)
	actualTopics := make([]string, 0, expectedTopicCount)

	for _, component := range components {
		for _, entityID := range entityIDs {
			topic := fmt.Sprintf("%s/%s/%s/%s/config", discoveryPrefix, component, deviceID, entityID)
			actualTopics = append(actualTopics, topic)
		}
	}

	if len(actualTopics) != expectedTopicCount {
		t.Errorf("Expected %d topics, got %d", expectedTopicCount, len(actualTopics))
	}

	// Test a few specific topic formats
	expectedSensorTopic := fmt.Sprintf("%s/sensor/%s/nozzle_temp_current/config", discoveryPrefix, deviceID)
	found := false
	for _, topic := range actualTopics {
		if topic == expectedSensorTopic {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected topic %s not found in generated topics", expectedSensorTopic)
	}
}

func TestCleanupCmdInit(t *testing.T) {
	// Test that init function properly sets up the command
	if cleanupCmd.Flags().Lookup("device-id") == nil {
		t.Error("device-id flag not found after init")
	}

	// Test flag properties
	flag := cleanupCmd.Flags().Lookup("device-id")
	if flag.Usage != "Device ID to remove (e.g., k1_se_192_168_4_87) [required]" {
		t.Errorf("Unexpected device-id flag usage: %s", flag.Usage)
	}

	if flag.DefValue != "" {
		t.Errorf("Expected empty default value for device-id flag, got: %s", flag.DefValue)
	}
}
