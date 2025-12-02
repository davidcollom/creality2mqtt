package mqttclient

import (
	"fmt"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Client struct {
	client mqtt.Client
	mu     sync.RWMutex

	// rate limiting
	minInterval   time.Duration
	lastPublished map[string]time.Time
	// pending payload to coalesce within interval
	lastPayload map[string]string
	// test hook to bypass IsConnected checks
	testBypassConnection bool
}

func New(brokerURL, clientID, username, password, willTopic, willPayload string) (*Client, error) {
	opts := mqtt.NewClientOptions().
		AddBroker(brokerURL).
		SetClientID(clientID).
		SetUsername(username).
		SetPassword(password).
		SetCleanSession(true).
		SetConnectTimeout(10 * time.Second).
		SetAutoReconnect(true)

	// Set Last Will Testament
	if willTopic != "" {
		opts.SetWill(willTopic, willPayload, 0, true)
	}

	c := mqtt.NewClient(opts)
	token := c.Connect()

	// Fail fast on initial connection
	if ok := token.WaitTimeout(10 * time.Second); !ok {
		return nil, fmt.Errorf("mqtt connection timeout after 10 seconds")
	} else if token.Error() != nil {
		return nil, fmt.Errorf("mqtt connection failed: %w", token.Error())
	}

	log.Info("MQTT connected", "broker", brokerURL)
	return &Client{client: c, minInterval: 0, lastPublished: make(map[string]time.Time), lastPayload: make(map[string]string)}, nil
}

func (c *Client) Disconnect() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.client.IsConnected() {
		log.Debug("Disconnecting MQTT client")
		c.client.Disconnect(250)
	}
}

func (c *Client) Publish(topic, payload string, retain bool) {
	// Shortcut for disconnected client
	c.mu.RLock()
	connected := c.testBypassConnection || c.client.IsConnected()
	c.mu.RUnlock()
	if !connected {
		log.Warn("MQTT not connected, dropping message", "topic", topic, "payload", payload)
		return
	}

	// Do not throttle retained messages (discovery, availability)
	if retain || c.minInterval <= 0 {
		token := c.client.Publish(topic, 0, retain, payload)
		ok := token.WaitTimeout(5 * time.Second)
		if !ok || token.Error() != nil {
			log.Error("MQTT publish failed", "topic", topic, "error", token.Error())
		}
		return
	}

	// Rate limiting with coalescing: at most one publish per interval per topic
	now := time.Now()

	// Critical section to decide whether to publish now or defer
	c.mu.Lock()
	last := c.lastPublished[topic]
	elapsed := now.Sub(last)
	if elapsed < c.minInterval {
		// Within interval: store latest payload and skip publish
		c.lastPayload[topic] = payload
		c.mu.Unlock()
		return
	}

	// Interval elapsed: choose the latest known payload for this topic
	toSend := payload
	if pending, ok := c.lastPayload[topic]; ok && pending != "" {
		toSend = pending
		delete(c.lastPayload, topic)
	}
	// Update lastPublished before unlock to avoid duplicate sends by racers
	c.lastPublished[topic] = now
	c.mu.Unlock()

	token := c.client.Publish(topic, 0, retain, toSend)
	ok := token.WaitTimeout(5 * time.Second)
	if !ok || token.Error() != nil {
		log.Error("MQTT publish failed", "topic", topic, "error", token.Error())
	}
}

// Subscribe subscribes to an MQTT topic with a message handler
func (c *Client) Subscribe(topic string, handler mqtt.MessageHandler) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.client.IsConnected() {
		return fmt.Errorf("mqtt not connected")
	}

	token := c.client.Subscribe(topic, 0, handler)
	ok := token.WaitTimeout(5 * time.Second)
	if !ok {
		return fmt.Errorf("subscribe timeout for topic: %s", topic)
	}
	if token.Error() != nil {
		return fmt.Errorf("subscribe failed for topic %s: %w", topic, token.Error())
	}

	log.Info("Subscribed to MQTT topic", "topic", topic)
	return nil
}

// SetMinInterval sets a minimum interval between publishes per topic.
// If set to 0, rate limiting is disabled.
func (c *Client) SetMinInterval(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.minInterval = d
}

// SetTestBypassConnection enables bypassing IsConnected checks (for unit tests only).
func (c *Client) SetTestBypassConnection(b bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.testBypassConnection = b
}
