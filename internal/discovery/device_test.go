package discovery

import (
	"testing"

	assert "github.com/stretchr/testify/assert"
)

func TestExtractDeviceInfo(t *testing.T) {
	in := map[string]any{
		"deviceName":  "K1 SE",
		"deviceModel": "K1 SE",
		"deviceId":    "K1-SE 01",
	}
	id, name, model := ExtractDeviceInfo(in)
	assert.Equal(t, "k1_se_01", id)
	assert.Equal(t, true, name != "")
	assert.Equal(t, true, model != "")
}
