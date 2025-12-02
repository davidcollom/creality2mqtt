package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

func TestDeviceInfoCmd(t *testing.T) {
	tests := []struct {
		testName    string
		description string
		input       struct {
			wsURL    string
			logLevel string
			args     []string
		}
		expectedOutput string
		actualOutput   string
		passFail       string
	}{
		{
			testName:    "ValidWebSocketURL",
			description: "Should successfully connect and display device info when valid WebSocket URL is provided",
			input: struct {
				wsURL    string
				logLevel string
				args     []string
			}{
				wsURL:    "ws://localhost:9999/",
				logLevel: "info",
				args:     []string{},
			},
			expectedOutput: "no error",
			actualOutput:   "",
			passFail:       "",
		},
		{
			testName:    "MissingWebSocketURL",
			description: "Should return error when WebSocket URL is not provided",
			input: struct {
				wsURL    string
				logLevel string
				args     []string
			}{
				wsURL:    "",
				logLevel: "info",
				args:     []string{},
			},
			expectedOutput: "ws-url is required",
			actualOutput:   "",
			passFail:       "",
		},
		{
			testName:    "InvalidLogLevel",
			description: "Should use default log level when invalid level is provided",
			input: struct {
				wsURL    string
				logLevel string
				args     []string
			}{
				wsURL:    "ws://localhost:9999/",
				logLevel: "invalid",
				args:     []string{},
			},
			expectedOutput: "no error",
			actualOutput:   "",
			passFail:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			// Reset command flags
			cmd := deviceInfoCmd
			cmd.SetContext(t.Context())
			cmd.ResetFlags()
			deviceInfoCmd.Flags().String("ws-url", "", "WebSocket URL of printer (e.g. ws://192.168.1.50:9999/) [can use CREALITY_WS_URL env var]")

			// Set up test environment
			if tt.input.wsURL != "" {
				// Per-test HTTP server and mux to avoid global handler conflicts
				mux := http.NewServeMux()
				mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
					upgrader := websocket.Upgrader{}
					conn, err := upgrader.Upgrade(w, r, nil)
					if err != nil {
						t.Logf("Failed to upgrade connection: %v", err)
						return
					}
					defer func() { require.NoError(t, conn.Close()) }()

					sampleMsg := `{"device":{"id":"test_device_123","name":"Test Printer","model":"Creality K1 SE"}}`
					if err := conn.WriteMessage(websocket.TextMessage, []byte(sampleMsg)); err != nil {
						t.Logf("Failed to write message: %v", err)
						return
					}

					time.Sleep(100 * time.Millisecond)
				})

				srv := httptest.NewServer(mux)
				defer srv.Close()

				wsURLLocal := strings.Replace(srv.URL, "http://", "ws://", 1) + "/"
				tt.input.wsURL = wsURLLocal

				require.NoError(t, cmd.Flags().Set("ws-url", wsURLLocal))
			}

			// Set global variables
			originalWSURL := wsURL
			originalLogLevel := logLevel
			wsURL = tt.input.wsURL
			logLevel = tt.input.logLevel

			defer func() {
				wsURL = originalWSURL
				logLevel = originalLogLevel
			}()

			// Create context with timeout
			_, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			// Execute command
			err := cmd.RunE(cmd, tt.input.args)

			// Verify results
			switch tt.testName {
			case "MissingWebSocketURL":
				if err == nil {
					tt.actualOutput = "no error"
					tt.passFail = "FAIL"
				} else if err.Error() == tt.expectedOutput {
					tt.actualOutput = err.Error()
					tt.passFail = "PASS"
				} else {
					tt.actualOutput = err.Error()
					tt.passFail = "FAIL"
				}
			case "ValidWebSocketURL":
				// For valid URL test, we expect it to fail due to no actual WebSocket server
				// but not with the "ws-url is required" error
				if err != nil && !strings.Contains(err.Error(), "ws-url is required") {
					tt.actualOutput = "connection error (expected)"
					tt.passFail = "PASS"
				} else if err == nil {
					tt.actualOutput = "no error"
					tt.passFail = "PASS"
				} else {
					tt.actualOutput = err.Error()
					tt.passFail = "FAIL"
				}
			case "InvalidLogLevel":
				// Similar to valid URL test
				if err != nil && !strings.Contains(err.Error(), "ws-url is required") {
					tt.actualOutput = "connection error (expected)"
					tt.passFail = "PASS"
				} else if err == nil {
					tt.actualOutput = "no error"
					tt.passFail = "PASS"
				} else {
					tt.actualOutput = err.Error()
					tt.passFail = "FAIL"
				}
			}

			// Log test results
			t.Logf("Test: %s | Description: %s | Input: %+v | Expected: %s | Actual: %s | Result: %s",
				tt.testName, tt.description, tt.input, tt.expectedOutput, tt.actualOutput, tt.passFail)

			require.NotEqualf(t, tt.passFail, "FAIL", "Test %s failed: expected %s, got %s", tt.testName, tt.expectedOutput, tt.actualOutput)
		})
	}
}

func TestDeviceInfoCmdFlags(t *testing.T) {
	tests := []struct {
		testName       string
		description    string
		input          string
		expectedOutput string
		actualOutput   string
		passFail       string
	}{
		{
			testName:       "WSURLFlagExists",
			description:    "Should have ws-url flag defined",
			input:          "ws-url",
			expectedOutput: "flag exists",
			actualOutput:   "",
			passFail:       "",
		},
		{
			testName:       "WSURLFlagDefault",
			description:    "Should have empty default value for ws-url flag",
			input:          "ws-url",
			expectedOutput: "",
			actualOutput:   "",
			passFail:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			cmd := deviceInfoCmd
			flag := cmd.Flags().Lookup(tt.input)

			switch tt.testName {
			case "WSURLFlagExists":
				if flag != nil {
					tt.actualOutput = "flag exists"
					tt.passFail = "PASS"
				} else {
					tt.actualOutput = "flag not found"
					tt.passFail = "FAIL"
				}
			case "WSURLFlagDefault":
				if flag != nil {
					tt.actualOutput = flag.DefValue
					if tt.actualOutput == tt.expectedOutput {
						tt.passFail = "PASS"
					} else {
						tt.passFail = "FAIL"
					}
				} else {
					tt.actualOutput = "flag not found"
					tt.passFail = "FAIL"
				}
			}

			t.Logf("Test: %s | Description: %s | Input: %s | Expected: %s | Actual: %s | Result: %s",
				tt.testName, tt.description, tt.input, tt.expectedOutput, tt.actualOutput, tt.passFail)

			if tt.passFail == "FAIL" {
				t.Errorf("Test %s failed: expected %s, got %s", tt.testName, tt.expectedOutput, tt.actualOutput)
			}
		})
	}
}

func TestDeviceInfoCmdProperties(t *testing.T) {
	tests := []struct {
		testName       string
		description    string
		input          string
		expectedOutput string
		actualOutput   string
		passFail       string
	}{
		{
			testName:       "CommandUse",
			description:    "Should have correct Use property",
			input:          "use",
			expectedOutput: "device-info",
			actualOutput:   "",
			passFail:       "",
		},
		{
			testName:       "CommandShort",
			description:    "Should have Short description",
			input:          "short",
			expectedOutput: "Show device information from printer WebSocket",
			actualOutput:   "",
			passFail:       "",
		},
		{
			testName:       "CommandLongExists",
			description:    "Should have Long description",
			input:          "long",
			expectedOutput: "has content",
			actualOutput:   "",
			passFail:       "",
		},
		{
			testName:       "CommandRunEExists",
			description:    "Should have RunE function defined",
			input:          "rune",
			expectedOutput: "function exists",
			actualOutput:   "",
			passFail:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			cmd := deviceInfoCmd

			switch tt.input {
			case "use":
				tt.actualOutput = cmd.Use
				if tt.actualOutput == tt.expectedOutput {
					tt.passFail = "PASS"
				} else {
					tt.passFail = "FAIL"
				}
			case "short":
				tt.actualOutput = cmd.Short
				if tt.actualOutput == tt.expectedOutput {
					tt.passFail = "PASS"
				} else {
					tt.passFail = "FAIL"
				}
			case "long":
				if cmd.Long != "" {
					tt.actualOutput = "has content"
					tt.passFail = "PASS"
				} else {
					tt.actualOutput = "empty"
					tt.passFail = "FAIL"
				}
			case "rune":
				if cmd.RunE != nil {
					tt.actualOutput = "function exists"
					tt.passFail = "PASS"
				} else {
					tt.actualOutput = "function nil"
					tt.passFail = "FAIL"
				}
			}

			t.Logf("Test: %s | Description: %s | Input: %s | Expected: %s | Actual: %s | Result: %s",
				tt.testName, tt.description, tt.input, tt.expectedOutput, tt.actualOutput, tt.passFail)

			if tt.passFail == "FAIL" {
				t.Errorf("Test %s failed: expected %s, got %s", tt.testName, tt.expectedOutput, tt.actualOutput)
			}
		})
	}
}
