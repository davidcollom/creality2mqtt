package main

import (
	"encoding/json"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/davidcollom/creality2mqtt/internal/discovery"
	"github.com/davidcollom/creality2mqtt/internal/wsclient"
	"github.com/spf13/cobra"
)

var deviceInfoCmd = &cobra.Command{
	Use:   "device-info",
	Short: "Show device information from printer WebSocket",
	Long: `Connects to the printer's WebSocket and displays the MQTT device ID and other information.
Use this to find the device-id needed for the cleanup command.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if wsURL == "" {
			log.Error("--ws-url is required")
			return fmt.Errorf("ws-url is required")
		}

		// Set log level
		level, err := log.ParseLevel(logLevel)
		if err != nil {
			log.Warn("Invalid log level, using info", "level", logLevel)
			level = log.InfoLevel
		}
		log.SetLevel(level)

		log.Info("Connecting to printer", "url", wsURL)

		done := make(chan bool)

		handler := func(data []byte) {
			var rawMsg map[string]any
			if err := json.Unmarshal(data, &rawMsg); err != nil {
				log.Error("Failed to parse WebSocket message", "error", err)
				return
			}

			// Extract device info
			deviceID, deviceName, deviceModel := discovery.ExtractDeviceInfo(rawMsg)

			fmt.Println()
			fmt.Println("═══════════════════════════════════════════════════════════")
			fmt.Println("              DEVICE INFORMATION")
			fmt.Println("═══════════════════════════════════════════════════════════")
			fmt.Println()
			fmt.Printf("  MQTT Device ID:    %s\n", deviceID)
			fmt.Printf("  Device Name:       %s\n", deviceName)
			fmt.Printf("  Device Model:      %s\n", deviceModel)
			fmt.Println()
			fmt.Println("───────────────────────────────────────────────────────────")
			fmt.Println("  Use this device ID for cleanup:")
			fmt.Println()
			fmt.Printf("    ./creality2mqtt cleanup --device-id=\"%s\"\n", deviceID)
			fmt.Println()
			fmt.Println("═══════════════════════════════════════════════════════════")
			fmt.Println()

			close(done)
		}

		ws := wsclient.New(wsURL, handler)

		// Run WebSocket Client in background
		errCh := make(chan error, 1)
		go func() {
			if err := ws.Run(cmd.Context()); err != nil {
				errCh <- err
			}
		}()

		// Wait for first message or error
		select {
		case <-done:
			return nil
		case err := <-errCh:
			return fmt.Errorf("websocket error: %w", err)
		}
	},
}

func init() {
	deviceInfoCmd.Flags().String("ws-url", "", "WebSocket URL of printer (e.g. ws://192.168.1.50:9999/) [can use CREALITY_WS_URL env var]")
}
