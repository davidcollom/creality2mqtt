package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/charmbracelet/log"
	"github.com/davidcollom/creality2mqtt/internal/discovery"
	"github.com/davidcollom/creality2mqtt/internal/mapper"
	"github.com/davidcollom/creality2mqtt/internal/mqttclient"
	"github.com/davidcollom/creality2mqtt/internal/types"
	"github.com/davidcollom/creality2mqtt/internal/wsclient"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the Creality to MQTT bridge",
	Long:  `Connects to a Creality 3D printer via WebSocket and publishes data to an MQTT broker.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if wsURL == "" {
			return cmd.Help()
		}

		// Set log level
		level, err := log.ParseLevel(logLevel)
		if err != nil {
			log.Warn("Invalid log level, using info", "level", logLevel)
			level = log.InfoLevel
		}
		log.SetLevel(level)

		log.Info("Starting creality2mqtt bridge")
		log.Debug("Configuration",
			"ws_url", wsURL,
			"mqtt_broker", broker,
			"mqtt_client_id", clientID,
			"base_topic", baseTopic,
			"discovery_prefix", discoveryPrefix,
			"log_level", logLevel,
		)

		// Create topic builder for consistent topic construction
		topics := types.NewTopicBuilder(baseTopic, discoveryPrefix)

		// Set up Last Will Testament (LWT) - published when we disconnect unexpectedly
		mqttClient, err := mqttclient.New(broker, clientID, username, password, topics.Availability(), "offline")
		if err != nil {
			log.Fatal("Failed to connect to MQTT broker", "error", err)
		}
		if mqttMinInterval > 0 {
			log.Info("Enabling MQTT rate limiting", "min_interval", mqttMinInterval)
			mqttClient.SetMinInterval(mqttMinInterval)
		} else {
			log.Info("MQTT rate limiting disabled")
		}
		defer func() {
			// Publish offline status before disconnecting
			mqttClient.Publish(topics.Availability(), "offline", true)
			mqttClient.Disconnect()
		}()

		// Publish birth message (we're online)
		mqttClient.Publish(topics.Availability(), "online", true)
		log.Info("Published birth message", "topic", topics.Availability())

		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer cancel()

		// Track discovery config
		var discoveryMu sync.Mutex
		var discoCfg *discovery.Config
		var discoveryMsgs []types.MqttMessage
		// Track published CFS box discovery to avoid duplicates
		publishedCFS := map[int]bool{}

		// Helper function to publish discovery messages
		publishDiscovery := func() {
			discoveryMu.Lock()
			defer discoveryMu.Unlock()
			if discoCfg != nil && len(discoveryMsgs) > 0 {
				log.Info("Publishing MQTT Discovery messages", "count", len(discoveryMsgs))
				for _, m := range discoveryMsgs {
					log.Debug("Publishing discovery config", "topic", m.Topic)
					mqttClient.Publish(m.Topic, m.Payload, m.Retain)
				}
				log.Info("MQTT Discovery published - entities should appear in Home Assistant")
			}
		}

		// Subscribe to Home Assistant status to republish discovery on HA restart
		err = mqttClient.Subscribe(topics.HAStatus(), func(client mqtt.Client, msg mqtt.Message) {
			payload := string(msg.Payload())
			log.Debug("Home Assistant status changed", "status", payload)
			if payload == "online" {
				log.Info("Home Assistant came online, republishing discovery")
				publishDiscovery()
			}
		})
		if err != nil {
			log.Warn("Failed to subscribe to HA status", "error", err)
		}

		// Track if we've sent discovery messages
		var discoveryOnce sync.Once

		// Create WebSocket client (before handler so we can reference it)
		ws := wsclient.New(wsURL, nil)

		// Subscribe to light control command topic
		err = mqttClient.Subscribe(topics.LightCommand(), func(client mqtt.Client, msg mqtt.Message) {
			payload := string(msg.Payload())
			log.Info("Received light command", "payload", payload)

			// Convert MQTT payload to WebSocket command
			var lightValue int
			switch payload {
			case "ON", "1":
				lightValue = 1
			case "OFF", "0":
				lightValue = 0
			default:
				log.Warn("Invalid light command payload", "payload", payload)
				return
			}

			// Send command to printer via WebSocket with proper format
			wsCmd := fmt.Sprintf(`{"method":"set","params":{"lightSw":%d}}`, lightValue)
			if err := ws.SendMessage([]byte(wsCmd)); err != nil {
				log.Error("Failed to send light command to printer", "error", err)
			} else {
				log.Info("Sent light command to printer", "command", wsCmd)
			}
		})
		if err != nil {
			log.Warn("Failed to subscribe to light command topic", "error", err)
		}

		handler := func(data []byte) {
			log.Debug("Received WebSocket message", "size", len(data))

			// Parse the message to extract device info on first message
			var rawMsg map[string]any
			if err := json.Unmarshal(data, &rawMsg); err == nil {
				discoveryOnce.Do(func() {
					// Extract device info from first message
					deviceID, devName, deviceModel := discovery.ExtractDeviceInfo(rawMsg)
					if deviceName != "" {
						devName = deviceName // Use CLI override if provided
					}

					log.Info("Device detected",
						"device_id", deviceID,
						"device_name", devName,
						"device_model", deviceModel,
					)

					// Extract printer IP from WebSocket URL
					printerIP := ""
					if u, err := url.Parse(wsURL); err == nil {
						printerIP = strings.Split(u.Host, ":")[0]
					}

					// Store discovery config for later republishing
					discoveryMu.Lock()
					discoCfg = &discovery.Config{
						DiscoveryPrefix: discoveryPrefix,
						BaseTopic:       baseTopic,
						DeviceID:        deviceID,
						DeviceName:      devName,
						DeviceModel:     deviceModel,
						PrinterIP:       printerIP,
					}

					// First, cleanup old/unused entities
					cleanupMsgs := discovery.CleanupOldEntities(*discoCfg)
					if len(cleanupMsgs) > 0 {
						log.Info("Cleaning up old entities", "count", len(cleanupMsgs))
						for _, m := range cleanupMsgs {
							log.Debug("Removing old entity", "topic", m.Topic)
							mqttClient.Publish(m.Topic, m.Payload, m.Retain)
						}
					}

					// Then generate current discovery messages
					discoveryMsgs = discovery.GenerateDiscoveryMessages(*discoCfg)
					discoveryMu.Unlock()

					// Publish discovery messages
					publishDiscovery()
				})

				// Dynamic discovery for CFS box sensors when seen
				if discoCfg != nil {
					if bs, ok := rawMsg["boxState"].(map[string]any); ok {
						if idAny, ok := bs["id"]; ok {
							id := 0
							switch v := idAny.(type) {
							case float64:
								id = int(v)
							case int:
								id = v
							}
							if id > 0 {
								discoveryMu.Lock()
								if !publishedCFS[id] {
									device := &discovery.Device{
										Identifiers:  []string{discoCfg.DeviceID},
										Name:         discoCfg.DeviceName,
										Manufacturer: "Creality",
										Model:        discoCfg.DeviceModel,
									}
									topics := types.NewTopicBuilder(discoCfg.BaseTopic, discoCfg.DiscoveryPrefix)
									availTopic := topics.Availability()
									msgs := discovery.BuildCFSBoxSensors(*discoCfg, device, availTopic, id)
									for _, m := range msgs {
										log.Info("Publishing CFS discovery", "topic", m.Topic)
										mqttClient.Publish(m.Topic, m.Payload, m.Retain)
									}
									publishedCFS[id] = true
								}
								discoveryMu.Unlock()
							}
						}
					}
				}
			}

			msgs, err := mapper.DecodeAndMap(data, baseTopic)
			if err != nil {
				log.Error("Failed to decode message", "error", err)
				return
			}

			for _, m := range msgs {
				log.Debug("Publishing MQTT message", "topic", m.Topic, "payload", m.Payload)
				mqttClient.Publish(m.Topic, m.Payload, m.Retain)
			}
		}

		// Set handler on WebSocket client
		ws.SetHandler(handler)

		log.Info("Starting WebSocket connection", "url", wsURL)
		log.Info("Press Ctrl+C to stop")
		if err := ws.Run(ctx); err != nil && err != context.Canceled {
			log.Error("WebSocket error", "error", err)
		}

		log.Info("Shutting down gracefully")
		return nil
	},
}
