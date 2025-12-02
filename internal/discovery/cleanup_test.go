package discovery

import "testing"

func TestCleanupOldEntities(t *testing.T) {
	cfg := Config{DiscoveryPrefix: "ha", DeviceID: "dev"}
	msgs := CleanupOldEntities(cfg)
	if len(msgs) == 0 {
		t.Fatalf("expected cleanup messages")
	}
	for _, m := range msgs {
		if m.Payload != "" || !m.Retain {
			t.Fatalf("cleanup message must be retained with empty payload")
		}
	}
}
