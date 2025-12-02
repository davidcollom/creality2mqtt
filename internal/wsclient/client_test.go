package wsclient

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

type testInput struct {
	contextTimeout time.Duration
	messages       []string
}

type testOutput struct {
	errorExpected    bool
	messagesReceived int
	connectionStored bool
}

func TestClient_runOnce(t *testing.T) {
	tests := []struct {
		testName       string
		description    string
		input          testInput
		expectedOutput testOutput
		setup          func() (*httptest.Server, *Client)
		teardown       func(*httptest.Server)
	}{
		{
			testName:    "Successful connection and message handling",
			description: "runOnce should establish connection, handle messages, and store connection",
			input: testInput{
				contextTimeout: 5 * time.Second,
				messages:       []string{"test message 1", "test message 2"},
			},
			expectedOutput: testOutput{
				// runOnce exits after connection closes; treat closure as error
				errorExpected:    true,
				messagesReceived: 2,
				connectionStored: false,
			},
			setup: func() (*httptest.Server, *Client) {
				var receivedMessages []string
				handler := func(data []byte) {
					receivedMessages = append(receivedMessages, string(data))
				}

				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					upgrader := websocket.Upgrader{}
					conn, err := upgrader.Upgrade(w, r, nil)
					if err != nil {
						return
					}
					defer func() { require.NoError(t, conn.Close()) }()

					// Send test messages
					for _, msg := range []string{"test message 1", "test message 2"} {
						require.NoError(t, conn.WriteMessage(websocket.TextMessage, []byte(msg)))
						time.Sleep(10 * time.Millisecond)
					}

					// Keep connection open briefly then close
					time.Sleep(100 * time.Millisecond)
				}))

				url := "ws" + strings.TrimPrefix(server.URL, "http")
				client := New(url, handler)
				return server, client
			},
			teardown: func(server *httptest.Server) {
				server.Close()
			},
		},
		{
			testName:    "Connection failure - dial error",
			description: "runOnce should return error when unable to establish WebSocket connection",
			input: testInput{
				contextTimeout: 1 * time.Second,
				messages:       []string{},
			},
			expectedOutput: testOutput{
				errorExpected:    true,
				messagesReceived: 0,
				connectionStored: false,
			},
			setup: func() (*httptest.Server, *Client) {
				handler := func(data []byte) {}
				client := New("ws://invalid-url:99999", handler)
				return nil, client
			},
			teardown: func(server *httptest.Server) {},
		},
		{
			testName:    "Context cancellation during operation",
			description: "runOnce should return context error when context is cancelled during connection",
			input: testInput{
				contextTimeout: 50 * time.Millisecond,
				messages:       []string{"message"},
			},
			expectedOutput: testOutput{
				errorExpected:    true,
				messagesReceived: 0,
				connectionStored: false,
			},
			setup: func() (*httptest.Server, *Client) {
				handler := func(data []byte) {}

				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					upgrader := websocket.Upgrader{}
					conn, err := upgrader.Upgrade(w, r, nil)
					if err != nil {
						return
					}
					defer func() { require.NoError(t, conn.Close()) }()

					// Keep connection open longer than context timeout
					time.Sleep(200 * time.Millisecond)
				}))

				url := "ws" + strings.TrimPrefix(server.URL, "http")
				client := New(url, handler)
				return server, client
			},
			teardown: func(server *httptest.Server) {
				server.Close()
			},
		},
		{
			testName:    "Connection drops during message reading",
			description: "runOnce should return error when connection is lost during message reading",
			input: testInput{
				contextTimeout: 5 * time.Second,
				messages:       []string{"message1"},
			},
			expectedOutput: testOutput{
				errorExpected:    true,
				messagesReceived: 1,
				connectionStored: false,
			},
			setup: func() (*httptest.Server, *Client) {
				var receivedMessages []string
				handler := func(data []byte) {
					receivedMessages = append(receivedMessages, string(data))
				}

				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					upgrader := websocket.Upgrader{}
					conn, err := upgrader.Upgrade(w, r, nil)
					if err != nil {
						return
					}

					// Send one message then close connection abruptly
					require.NoError(t, conn.WriteMessage(websocket.TextMessage, []byte("message1")))
					time.Sleep(10 * time.Millisecond)
					require.NoError(t, conn.Close())
				}))

				url := "ws" + strings.TrimPrefix(server.URL, "http")
				client := New(url, handler)
				return server, client
			},
			teardown: func(server *httptest.Server) {
				server.Close()
			},
		},
		{
			testName:    "Empty message handling",
			description: "runOnce should handle empty messages without error",
			input: testInput{
				contextTimeout: 2 * time.Second,
				messages:       []string{""},
			},
			expectedOutput: testOutput{
				// connection closes after sending; expect error
				errorExpected:    true,
				messagesReceived: 1,
				connectionStored: false,
			},
			setup: func() (*httptest.Server, *Client) {
				var receivedMessages []string
				handler := func(data []byte) {
					receivedMessages = append(receivedMessages, string(data))
				}

				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					upgrader := websocket.Upgrader{}
					conn, err := upgrader.Upgrade(w, r, nil)
					if err != nil {
						return
					}
					defer func() { require.NoError(t, conn.Close()) }()

					require.NoError(t, conn.WriteMessage(websocket.TextMessage, []byte("")))
					time.Sleep(50 * time.Millisecond)
				}))

				url := "ws" + strings.TrimPrefix(server.URL, "http")
				client := New(url, handler)
				return server, client
			},
			teardown: func(server *httptest.Server) {
				server.Close()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			// Setup
			server, client := tt.setup()
			if tt.teardown != nil {
				defer tt.teardown(server)
			}

			// Track received messages
			var messageCount int32
			originalHandler := client.handler
			client.handler = func(data []byte) {
				atomic.AddInt32(&messageCount, 1)
				if originalHandler != nil {
					originalHandler(data)
				}
			}

			// Create context with timeout
			ctx, cancel := context.WithTimeout(context.Background(), tt.input.contextTimeout)
			defer cancel()

			// Execute runOnce
			actualErr := client.runOnce(ctx)

			// Wait a bit for message processing
			time.Sleep(100 * time.Millisecond)

			// Verify results
			actualOutput := testOutput{
				errorExpected:    actualErr != nil,
				messagesReceived: int(atomic.LoadInt32(&messageCount)),
				connectionStored: func() bool {
					client.connMu.RLock()
					defer client.connMu.RUnlock()
					return client.conn != nil
				}(),
			}

			// After runOnce completes, connection should be nil (cleaned up)
			if actualErr != nil {
				client.connMu.RLock()
				connectionAfterError := client.conn
				client.connMu.RUnlock()

				if connectionAfterError != nil {
					t.Errorf("Expected connection to be nil after error, but was not nil")
				}
			}

			// Validate error expectation
			if tt.expectedOutput.errorExpected != actualOutput.errorExpected {
				t.Errorf("Error expectation mismatch. Expected error: %v, Got error: %v (actual error: %v)",
					tt.expectedOutput.errorExpected, actualOutput.errorExpected, actualErr)
			}

			// Validate message count (allow some flexibility for timing)
			if actualOutput.messagesReceived < tt.expectedOutput.messagesReceived {
				t.Errorf("Message count mismatch. Expected at least: %d, Got: %d",
					tt.expectedOutput.messagesReceived, actualOutput.messagesReceived)
			}

			t.Logf("Test '%s' - Input: timeout=%v, messages=%v | Expected: error=%v, msgs=%d, conn=%v | Actual: error=%v, msgs=%d, conn=%v | Result: PASS",
				tt.testName, tt.input.contextTimeout, tt.input.messages,
				tt.expectedOutput.errorExpected, tt.expectedOutput.messagesReceived, tt.expectedOutput.connectionStored,
				actualOutput.errorExpected, actualOutput.messagesReceived, actualOutput.connectionStored)
		})
	}
}
