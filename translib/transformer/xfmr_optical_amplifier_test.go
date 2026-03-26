//go:build xfmrtest
// +build xfmrtest

package transformer_test

import (
	"fmt"
	"github.com/Azure/sonic-mgmt-common/translib/db"
	"testing"
	"time"
)

/* AMPLIFIER config,counters, state DB fields and values */
func Test_amplifier_config_state_counter_DB_key(t *testing.T) {
	var pre_req_map, cleanuptbl, cleanupstatetbl, cleanupcountertbl map[string]interface{}
	var expected_get_json, url []string

	for idx, dbName := range getMdbNames() {
		/* AMPLIFIER Config DB test */
		t.Log("\n\n+++++++++++++ Performing Get on configDB AMPLIFIER Key ++++++++++++")
		pre_req_map = map[string]interface{}{
			"AMPLIFIER": map[string]interface{}{
				fmt.Sprintf("AMPLIFIER-1-%d-2", idx+1): map[string]interface{}{
					"type":                "EDFA",
					"target-gain":         15.5,
					"max-gain":            20.0,
					"min-gain":            10.0,
					"target-gain-tilt":    2.5,
					"gain-range":          "MID_GAIN_RANGE",
					"amp-mode":            "CONSTANT_GAIN",
					"target-output-power": 0.0,
					"max-output-power":    5.0,
					"enabled":             true,
					"fiber-type-profile":  "SSMF",
				},
			},
		}
		loadDB(dbName, db.ConfigDB, pre_req_map)
		expected_get_json = []string{fmt.Sprintf("{\"openconfig-optical-amplifier:config\":{\"amp-mode\":\"openconfig-optical-amplifier:CONSTANT_GAIN\",\"enabled\":true,\"fiber-type-profile\":\"openconfig-optical-amplifier:SSMF\",\"gain-range\":\"openconfig-optical-amplifier:MID_GAIN_RANGE\",\"max-gain\":\"20\",\"max-output-power\":\"5\",\"min-gain\":\"10\",\"name\":\"AMPLIFIER-1-%d-2\",\"target-gain\":\"15.5\",\"target-gain-tilt\":\"2.5\",\"target-output-power\":\"0\",\"type\":\"openconfig-optical-amplifier:EDFA\"}}", idx+1)}

		url = []string{fmt.Sprintf("/openconfig-optical-amplifier:optical-amplifier/amplifiers/amplifier[name=AMPLIFIER-1-%d-2]/config", idx+1)}
		t.Run("Test get on configDB AMPLIFIER Key ", processGetRequest(url[0], nil, expected_get_json[0], false))
		time.Sleep(1 * time.Second)
		t.Log("\n\n+++++++++++++ Done Performing Get on configDB AMPLIFIER Key ++++++++++++")

		/* AMPLIFIER and OSC State DB test */
		t.Log("\n\n+++++++++++++ Performing Get on stateDB AMPLIFIER , OSC Key ++++++++++++")
		pre_req_map = map[string]interface{}{
			"AMPLIFIER_TABLE": map[string]interface{}{
				fmt.Sprintf("AMPLIFIER-1-%d-2", idx+1): map[string]interface{}{
					"component":           fmt.Sprintf("AMPLIFIER-1-%d-2", idx+1),
					"location":            "1-4",
					"name":                fmt.Sprintf("AMPLIFIER-1-%d-2", idx+1),
					"amp-mode":            "CONSTANT_POWER",
					"egress-port":         "test data",
					"enabled":             "true",
					"fiber-type-profile":  "DSF",
					"gain-range":          "LOW_GAIN_RANGE",
					"ingress-port":        "test data",
					"max-gain":            "10.5",
					"max-output-power":    "11.2",
					"min-gain":            "5.6",
					"target-gain":         "11.15",
					"target-gain-tilt":    "16.10",
					"target-output-power": "10.30",
					"type":                "EDFA",
				},
			},
			"OSC_TABLE": map[string]interface{}{
				fmt.Sprintf("OSC-1-%d-2", idx+1): map[string]interface{}{
					"output-frequency": "198538051",
				},
			},
		}
		loadDB(dbName, db.StateDB, pre_req_map)
		expected_get_json = []string{fmt.Sprintf("{\"openconfig-optical-amplifier:state\":{\"amp-mode\":\"openconfig-optical-amplifier:CONSTANT_POWER\",\"component\":\"AMPLIFIER-1-%d-2\",\"egress-port\":\"test data\",\"enabled\":true,\"fiber-type-profile\":\"openconfig-optical-amplifier:DSF\",\"gain-range\":\"openconfig-optical-amplifier:LOW_GAIN_RANGE\",\"ingress-port\":\"test data\",\"max-gain\":\"10.5\",\"max-output-power\":\"11.2\",\"min-gain\":\"5.6\",\"name\":\"AMPLIFIER-1-%d-2\",\"target-gain\":\"11.15\",\"target-gain-tilt\":\"16.1\",\"target-output-power\":\"10.3\",\"type\":\"openconfig-optical-amplifier:EDFA\"}}", idx+1, idx+1), fmt.Sprintf("{\"openconfig-optical-amplifier:state\":{\"interface\":\"OSC-1-%d-2\",\"output-frequency\":\"198538051\"}}", idx+1)}
		url = []string{
			fmt.Sprintf("/openconfig-optical-amplifier:optical-amplifier/amplifiers/amplifier[name=AMPLIFIER-1-%d-2]/state", idx+1),
			fmt.Sprintf("/openconfig-optical-amplifier:optical-amplifier/supervisory-channels/supervisory-channel[interface=OSC-1-%d-2]/state", idx+1),
		}

		for idx, url := range url {
			t.Run("Test get on state DB Key for "+url, processGetRequest(url, nil, expected_get_json[idx], false))
			time.Sleep(1 * time.Second)
		}
		time.Sleep(1 * time.Second)

		t.Log("\n\n+++++++++++++ Done Performing Get on stateDB AMPLIFIER, OSC Key ++++++++++++")

		/* AMPLIFIER and OSC counters DB test */
		t.Log("\n\n+++++++++++++ Performing Get on counters AMPLIFIER , OSC Key ++++++++++++")
		pre_req_map = map[string]interface{}{
			"AMPLIFIER": map[string]interface{}{
				fmt.Sprintf("AMPLIFIER-1-%d-2_ActualGain:15_pm_current", idx+1): map[string]interface{}{
					"interval":         "86400000000000",
					"starttime":        "1709251200000000000",
					"max":              "1.0",
					"max-time":         "1909276854829015912",
					"min":              "1.0",
					"min-time":         "1909276854829015912",
					"instant":          "1.0",
					"avg":              "1.0",
					"current_validity": "complete",
					"validity":         "incomplete",
				},
			},
			"OSC": map[string]interface{}{
				fmt.Sprintf("OSC-1-%d-2_OutputPower:15_pm_current", idx+1): map[string]interface{}{
					"interval":         "900000000000",
					"starttime":        "1709251200000000000",
					"max":              "2.5",
					"max-time":         "1680168786273045375",
					"min":              "2.12",
					"min-time":         "1680168676158235913",
					"instant":          "2.29",
					"avg":              "2.32",
					"current_validity": "complete",
					"validity":         "incomplete",
				},
			},
		}
		loadDB(dbName, db.CountersDB, pre_req_map)
		expected_get_json = []string{"{\"openconfig-optical-amplifier:actual-gain\":{\"avg\":\"1\",\"instant\":\"1\",\"interval\":\"86400000000000\",\"max\":\"1\",\"max-time\":\"1909276854829015912\",\"min\":\"1\",\"min-time\":\"1909276854829015912\"}}", "{\"openconfig-optical-amplifier:output-power\":{\"avg\":\"2.32\",\"instant\":\"2.29\",\"interval\":\"900000000000\",\"max\":\"2.5\",\"max-time\":\"1680168786273045375\",\"min\":\"2.12\",\"min-time\":\"1680168676158235913\"}}"}

		url = []string{
			fmt.Sprintf("/openconfig-optical-amplifier:optical-amplifier/amplifiers/amplifier[name=AMPLIFIER-1-%d-2]/state/actual-gain", idx+1),
			fmt.Sprintf("/openconfig-optical-amplifier:optical-amplifier/supervisory-channels/supervisory-channel[interface=OSC-1-%d-2]/state/output-power", idx+1),
		}
		for idx, url := range url {
			t.Run("Test get on counters DB Key "+url, processGetRequest(url, nil, expected_get_json[idx], false))
			time.Sleep(1 * time.Second)
		}
		t.Log("\n\n+++++++++++++ Performing Get on counters AMPLIFIER, OSC Key ++++++++++++")

		//Unload the Data
		cleanuptbl = map[string]interface{}{
			"AMPLIFIER": map[string]interface{}{fmt.Sprintf("AMPLIFIER-1-%d-2", idx+1): ""},
		}
		unloadDB(dbName, db.ConfigDB, cleanuptbl)

		cleanupstatetbl = map[string]interface{}{
			"AMPLIFIER_TABLE": map[string]interface{}{fmt.Sprintf("AMPLIFIER-1-%d-2", idx+1): ""},
			"OSC_TABLE":       map[string]interface{}{fmt.Sprintf("OSC-1-%d-2", idx+1): ""},
		}
		unloadDB(dbName, db.StateDB, cleanupstatetbl)

		cleanupcountertbl = map[string]interface{}{
			"AMPLIFIER": map[string]interface{}{fmt.Sprintf("AMPLIFIER-1-%d-2_ActualGain:15_pm_current", idx+1): ""},
			"OSC":       map[string]interface{}{fmt.Sprintf("OSC-1-%d-2_OutputPower:15_pm_current", idx+1): ""},
		}
		unloadDB(dbName, db.CountersDB, cleanupcountertbl)
	}
}

func Test_set_amplifier_config_DB_key_and_field_xfmr(t *testing.T) {
	t.Log("\n\n+++++++++++++ Performing Create/Replace/Delete on amplifier ++++++++++++")

	for idx, dbName := range getMdbNames() {
		url := fmt.Sprintf("/openconfig-optical-amplifier:optical-amplifier/amplifiers/amplifier[name=AMPLIFIER-1-%d-3]/config", idx+1)
		url_body_json := fmt.Sprintf("{\"config\":{\"amp-mode\":\"openconfig-optical-amplifier:CONSTANT_GAIN\",\"enabled\":false,\"fiber-type-profile\":\"openconfig-optical-amplifier:SSMF\",\"gain-range\":\"openconfig-optical-amplifier:MID_GAIN_RANGE\",\"max-gain\":\"30\",\"max-output-power\":\"5\",\"min-gain\":\"15\",\"name\":\"AMPLIFIER-1-%d-3\",\"target-gain\":\"10.5\",\"target-gain-tilt\":\"12.5\",\"target-output-power\":\"0\",\"type\":\"openconfig-optical-amplifier:EDFA\"}}", idx+1)

		pre_req_map := map[string]interface{}{
			"AMPLIFIER": map[string]interface{}{
				fmt.Sprintf("AMPLIFIER-1-%d-3", idx+1): map[string]interface{}{
					"type":                "EDFA",
					"target-gain":         15.5,
					"max-gain":            20.0,
					"min-gain":            10.0,
					"target-gain-tilt":    2.5,
					"gain-range":          "MID_GAIN_RANGE",
					"amp-mode":            "CONSTANT_GAIN",
					"target-output-power": 0.0,
					"max-output-power":    5.0,
					"enabled":             true,
					"fiber-type-profile":  "SSMF",
				},
			},
		}
		loadDB(dbName, db.ConfigDB, pre_req_map)
		expected_map := map[string]interface{}{
			"AMPLIFIER": map[string]interface{}{
				fmt.Sprintf("AMPLIFIER-1-%d-3", idx+1): map[string]interface{}{
					"type":                "EDFA",
					"target-gain":         10.5,
					"max-gain":            30.0,
					"min-gain":            15.0,
					"target-gain-tilt":    12.5,
					"gain-range":          "MID_GAIN_RANGE",
					"amp-mode":            "CONSTANT_GAIN",
					"target-output-power": 0.0,
					"max-output-power":    5.0,
					"enabled":             false,
					"fiber-type-profile":  "SSMF",
				},
			},
		}

		time.Sleep(1 * time.Second)

		//PUT Test
		t.Run("Replace Test on Amplifier yang", processSetRequest(url, url_body_json, "PUT", false))
		time.Sleep(1 * time.Second)

		//GET Test
		t.Run("Verify replace on Amplifier yang", verifyDbResult(rclient[dbName], fmt.Sprintf("AMPLIFIER|AMPLIFIER-1-%d-3", idx+1), expected_map, false))
		time.Sleep(1 * time.Second)

		//Delete Test
		url = fmt.Sprintf("/openconfig-optical-amplifier:optical-amplifier/amplifiers/amplifier[name=AMPLIFIER-1-%d-3]/config", idx+1)
		t.Run("DELETE Test on Amplifier yang", processDeleteRequest(url, false))
		time.Sleep(1 * time.Second)

		//POST Test
		url = "/openconfig-optical-amplifier:optical-amplifier/amplifiers"
		post_payload := fmt.Sprintf("{\"amplifier\":[{\"name\":\"AMPLIFIER-1-%d-4\",\"config\":{\"target-gain\":\"30.22\",\"max-gain\":\"41.42\",\"min-gain\":\"20.92\",\"target-gain-tilt\":\"59.22\",\"gain-range\":\"LOW_GAIN_RANGE\",\"amp-mode\":\"CONSTANT_POWER\",\"target-output-power\":\"82.33\",\"max-output-power\":\"199.99\",\"enabled\":false,\"fiber-type-profile\":\"DSF\",\"type\":\"openconfig-optical-amplifier:EDFA\"}}]}", idx+1)
		t.Run("Create test on Amplifier yang", processSetRequest(url, post_payload, "POST", false))
		time.Sleep(1 * time.Second)

		expected_map = map[string]interface{}{
			"AMPLIFIER": map[string]interface{}{
				fmt.Sprintf("AMPLIFIER-1-%d-4", idx+1): map[string]interface{}{
					"type":                "EDFA",
					"target-gain":         30.22,
					"max-gain":            41.42,
					"min-gain":            20.92,
					"target-gain-tilt":    59.22,
					"gain-range":          "LOW_GAIN_RANGE",
					"amp-mode":            "CONSTANT_POWER",
					"target-output-power": 82.33,
					"max-output-power":    199.99,
					"enabled":             false,
					"fiber-type-profile":  "DSF",
				},
			},
		}

		//Verify After the POST test
		t.Run("Verify After the POST test on Amplifier yang", verifyDbResult(rclient[dbName], fmt.Sprintf("AMPLIFIER|AMPLIFIER-1-%d-4", idx+1), expected_map, false))

		//Delete After the POST test
		url = fmt.Sprintf("/openconfig-optical-amplifier:optical-amplifier/amplifiers/amplifier[name=AMPLIFIER-1-%d-4]/config", idx+1)
		t.Run("DELETE After the POST test on Amplifier yang", processDeleteRequest(url, false))
		time.Sleep(1 * time.Second)

		url = "/openconfig-optical-amplifier:optical-amplifier/amplifiers"

		//Bulk POST Test
		Bulk_post_payload := fmt.Sprintf("{\"amplifier\":[{\"name\":\"AMPLIFIER-1-%d-5\",\"config\":{\"target-gain\":\"30.22\",\"max-gain\":\"41.42\",\"min-gain\":\"20.92\",\"target-gain-tilt\":\"59.22\",\"gain-range\":\"LOW_GAIN_RANGE\",\"amp-mode\":\"CONSTANT_POWER\",\"target-output-power\":\"82.33\",\"max-output-power\":\"199.99\",\"enabled\":false,\"fiber-type-profile\":\"DSF\",\"type\":\"openconfig-optical-amplifier:EDFA\"}},{\"name\":\"AMPLIFIER-1-%d-6\",\"config\":{\"target-gain\":\"20.22\",\"max-gain\":\"31.42\",\"min-gain\":\"15.92\",\"target-gain-tilt\":\"50.22\",\"gain-range\":\"MID_GAIN_RANGE\",\"amp-mode\":\"CONSTANT_POWER\",\"target-output-power\":\"12.33\",\"max-output-power\":\"100.99\",\"enabled\":true,\"fiber-type-profile\":\"DSF\",\"type\":\"openconfig-optical-amplifier:EDFA\"}}]}", idx+1, idx+1)
		t.Run("Create test on Amplifier yang", processSetRequest(url, Bulk_post_payload, "POST", false))

		expected_map = map[string]interface{}{
			"AMPLIFIER": map[string]interface{}{
				fmt.Sprintf("AMPLIFIER-1-%d-5", idx+1): map[string]interface{}{
					"type":                "EDFA",
					"target-gain":         30.22,
					"max-gain":            41.42,
					"min-gain":            20.92,
					"target-gain-tilt":    59.22,
					"gain-range":          "LOW_GAIN_RANGE",
					"amp-mode":            "CONSTANT_POWER",
					"target-output-power": 82.33,
					"max-output-power":    199.99,
					"enabled":             false,
					"fiber-type-profile":  "DSF",
				},
			},
		}

		time.Sleep(1 * time.Second)
		//Verify After the Bulk test
		t.Run("Verify After the Bulk test on Amplifier yang", verifyDbResult(rclient[dbName], fmt.Sprintf("AMPLIFIER|AMPLIFIER-1-%d-5", idx+1), expected_map, false))

		expected_map = map[string]interface{}{
			"AMPLIFIER": map[string]interface{}{
				fmt.Sprintf("AMPLIFIER-1-%d-6", idx+1): map[string]interface{}{
					"type":                "EDFA",
					"target-gain":         20.22,
					"max-gain":            31.42,
					"min-gain":            15.92,
					"target-gain-tilt":    50.22,
					"gain-range":          "MID_GAIN_RANGE",
					"amp-mode":            "CONSTANT_POWER",
					"target-output-power": 12.33,
					"max-output-power":    100.99,
					"enabled":             true,
					"fiber-type-profile":  "DSF",
				},
			},
		}

		time.Sleep(1 * time.Second)
		//Verify After the POST test
		t.Run("Verify After the Bulk test on Amplifier yang", verifyDbResult(rclient[dbName], fmt.Sprintf("AMPLIFIER|AMPLIFIER-1-%d-6", idx+1), expected_map, false))

		//Delete keys After the Bulk Test
		url = fmt.Sprintf("/openconfig-optical-amplifier:optical-amplifier/amplifiers/amplifier[name=AMPLIFIER-1-%d-5]/config", idx+1)
		t.Run("DELETE After the Bulk Test on Amplifier yang", processDeleteRequest(url, false))
		time.Sleep(1 * time.Second)

		url = fmt.Sprintf("/openconfig-optical-amplifier:optical-amplifier/amplifiers/amplifier[name=AMPLIFIER-1-%d-6]/config", idx+1)
		t.Run("DELETE After the Bulk Test on Amplifier yang", processDeleteRequest(url, false))
		time.Sleep(1 * time.Second)
	}
	t.Log("\n\n+++++++++++++ Done Performing Create/Replace/Delete on Amplifier ++++++++++++")
}

func Test_Amplifier_allow_write_multiple_namespace(t *testing.T) {
	t.Log("\n\n+++++++++++++ Performing Amplifier write for multiple namespace (should be allowed) ++++++++++++")

	url := "/openconfig-optical-amplifier:optical-amplifier/amplifiers"
	url_body_json := `{
        "amplifier": [
            {
                "name": "AMPLIFIER-1-1-5",
                "config": {
                    "target-gain": "30.22",
                    "max-gain": "41.42",
                    "min-gain": "20.92",
                    "target-gain-tilt": "59.22",
                    "gain-range": "LOW_GAIN_RANGE",
                    "amp-mode": "CONSTANT_POWER",
                    "target-output-power": "82.33",
                    "max-output-power": "199.99",
                    "enabled": false,
                    "fiber-type-profile": "DSF"
                }
            },
            {
                "name": "AMPLIFIER-1-2-6",
                "config": {
                    "target-gain": "20.22",
                    "max-gain": "31.42",
                    "min-gain": "15.92",
                    "target-gain-tilt": "50.22",
                    "gain-range": "MID_GAIN_RANGE",
                    "amp-mode": "CONSTANT_POWER",
                    "target-output-power": "12.33",
                    "max-output-power": "100.99",
                    "enabled": true,
                    "fiber-type-profile": "DSF"
                }
            }
        ]
    }`

	// false = no error expected
	t.Run("Test allow writing multiple namespace for Amplifier", processSetRequest(url, url_body_json, "POST", false, nil))
	t.Log("\n\n+++++++++++++ Done Performing Amplifier write for multiple namespace ++++++++++++")
}
