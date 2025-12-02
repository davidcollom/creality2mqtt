package mapper

import (
	"encoding/json"
	"testing"

	"github.com/davidcollom/creality2mqtt/internal/types"
)

func toTopicMap(msgs []types.MqttMessage) map[string]string {
	m := make(map[string]string, len(msgs))
	for _, mm := range msgs {
		m[mm.Topic] = mm.Payload
	}
	return m
}

func TestMapMessageToMqtt_InitialSnapshot(t *testing.T) {
	t.Parallel()

	baseTopic := "3dprinter/k1se"

	// This is a reduced version of your very first message: same keys, fewer fields.
	jsonStr := `{
      "TotalLayer": 0,
      "accelToDecelLimits": 2500,
      "accelerationLimits": 5000,
      "autoLevelResult": "25:-0.06",
      "autohome": "X:1 Y:1 Z:1",
      "auxiliaryFanPct": 0,
      "bedTemp0": "59.340000",
      "bedTemp1": "0.000000",
      "bedTemp2": "0.000000",
      "boxTemp": 0,
      "caseFanPct": 0,
      "cfsConnect": 1,
      "connect": 1,
      "curFeedratePct": 100,
      "curFlowratePct": 100,
      "curPosition": "X:142.82 Y:64.43 Z:5.33",
      "deviceState": 1,
      "err": { "errcode": 0, "key": 0 },
      "hostname": "K1 SE-7E0E",
      "layer": 62,
      "lightSw": 1,
      "maxBedTemp": 115,
      "maxNozzleTemp": 320,
      "model": "K1 SE",
      "modelFanPct": 100,
      "modelVersion": "printer hw ver:;printer sw ver:;DWIN hw ver:CR4CU220812S11;DWIN sw ver:2.3.5.33;",
      "nozzleTemp": "219.900000",
      "pressureAdvance": "0.044000",
      "printFileName": "/usr/data/printer_data/gcodes/Meta.gcode",
      "printJobTime": 1272,
      "printLeftTime": 767,
      "printProgress": 58,
      "printStartTime": 1764623088,
      "realTimeFlow": "3.170000",
      "realTimeSpeed": "210.670000",
      "state": 1,
      "targetBedTemp0": 60,
      "targetNozzleTemp": 220,
      "tfCard": 1,
      "velocityLimits": 600,
      "video": 1,
      "video1": 0,
      "videoElapse": 1,
      "videoElapseFrame": 15,
      "videoElapseInterval": 1
    }`

	var msg map[string]any
	if err := json.Unmarshal([]byte(jsonStr), &msg); err != nil {
		t.Fatalf("failed to unmarshal test json: %v", err)
	}

	got := MapMessageToMqtt(msg, baseTopic)
	tp := toTopicMap(got)

	// We don't assert every topic – just the ones we care most about.
	expectations := map[string]string{
		baseTopic + "/nozzle_temp":        "219.900",
		baseTopic + "/bed_temp0":          "59.340",
		baseTopic + "/box_temp":           "0",
		baseTopic + "/print_progress":     "58",
		baseTopic + "/print_job_time":     "1272",
		baseTopic + "/print_left_time":    "767",
		baseTopic + "/hostname":           "K1 SE-7E0E",
		baseTopic + "/model":              "K1 SE",
		baseTopic + "/device_state":       "1",
		baseTopic + "/state":              "1",
		baseTopic + "/cur_position":       "X:142.82 Y:64.43 Z:5.33",
		baseTopic + "/target_nozzle_temp": "220",
		baseTopic + "/target_bed_temp0":   "60",
	}

	for topic, want := range expectations {
		gotVal, ok := tp[topic]
		if !ok {
			t.Errorf("expected topic %s to be present", topic)
			continue
		}
		if gotVal != want {
			t.Errorf("topic %s payload = %q, want %q", topic, gotVal, want)
		}
	}

	// Sanity check that noisy video keys were skipped
	for _, noisy := range []string{"video", "video1", "video_elapse_frame", "video_elapse_interval"} {
		full := baseTopic + "/" + noisy
		if _, ok := tp[full]; ok {
			t.Errorf("expected noisy key %s to be skipped; found topic %s", noisy, full)
		}
	}
}

func TestMapMessageToMqtt_DeltaMessage(t *testing.T) {
	t.Parallel()

	baseTopic := "3dprinter/k1se"
	jsonStr := `{
        "nozzleTemp": "219.790000",
        "bedTemp0": "59.380000",
        "printJobTime": 1273,
        "usedMaterialLength": 2654,
        "realTimeFlow": "17.710000",
        "curPosition": "X:149.82 Y:72.50 Z:5.11"
    }`

	var msg map[string]any
	if err := json.Unmarshal([]byte(jsonStr), &msg); err != nil {
		t.Fatalf("failed to unmarshal test json: %v", err)
	}

	got := MapMessageToMqtt(msg, baseTopic)
	tp := toTopicMap(got)

	expectations := map[string]string{
		baseTopic + "/nozzle_temp":          "219.790",
		baseTopic + "/bed_temp0":            "59.380",
		baseTopic + "/print_job_time":       "1273",
		baseTopic + "/used_material_length": "2654",
		baseTopic + "/real_time_flow":       "17.710",
		baseTopic + "/cur_position":         "X:149.82 Y:72.50 Z:5.11",
	}

	for topic, want := range expectations {
		gotVal, ok := tp[topic]
		if !ok {
			t.Errorf("expected topic %s to be present", topic)
			continue
		}
		if gotVal != want {
			t.Errorf("topic %s payload = %q, want %q", topic, gotVal, want)
		}
	}
}
