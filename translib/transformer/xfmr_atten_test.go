//go:build xfmrtest
// +build xfmrtest

package transformer_test

import (
	"fmt"
	"github.com/Azure/sonic-mgmt-common/translib/db"
	"testing"
	"time"
)

/* Attenuator config, state, counters DB fields and values */
func Test_get_attenuator_config_state_counter_DB_key(t *testing.T) {
	var pre_req_map, cleanuptbl, cleanupstatetbl, cleanupcountertbl map[string]interface{}
	var url string

	for idx, dbName := range getMdbNames() {

		/* ATTENUATOR Config DB test */
		t.Log("\n\n+++++++++++++ Performing Get on configDB Attenuator Key ++++++++++++")
		pre_req_map = map[string]interface{}{
			"ATTENUATOR": map[string]interface{}{
				fmt.Sprintf("ATTENUATOR-1-%d-3", idx+1): map[string]interface{}{
					"attenuation":         "1.5",
					"attenuation-mode":    "CONSTANT_ATTENUATION",
					"target-output-power": "2.5",
					"enabled":             "true",
				},
			},
		}
		loadDB(dbName, db.ConfigDB, pre_req_map)
		expected_get_json := fmt.Sprintf("{\"openconfig-optical-attenuator:config\":{\"attenuation\":\"1.5\",\"attenuation-mode\":\"openconfig-optical-attenuator:CONSTANT_ATTENUATION\",\"enabled\":true,\"name\":\"ATTENUATOR-1-%d-3\",\"target-output-power\":\"2.5\"}}", idx+1)
		url = fmt.Sprintf("/openconfig-optical-attenuator:optical-attenuator/attenuators/attenuator[name=ATTENUATOR-1-%d-3]/config", idx+1)

		t.Run("Test get on ConfigDB Attnuator Key", processGetRequest(url, nil, expected_get_json, false))
		time.Sleep(1 * time.Second)
		t.Log("\n\n+++++++++++++ Done Performing Get on configDB Attenuator Key ++++++++++++")

		/* ATTENUATOR State DB test */
		t.Log("\n\n+++++++++++++ Performing Get on stateDB Attenuator Key ++++++++++++")
		pre_req_map = map[string]interface{}{
			"ATTENUATOR_TABLE": map[string]interface{}{
				fmt.Sprintf("ATTENUATOR-1-%d-3", idx+1): map[string]interface{}{
					"attenuation":         "2.5",
					"name":                fmt.Sprintf("ATTENUATOR-1-%d-3", idx+1),
					"target-output-power": "1.5",
					"ingress-port":        "PORT-1-2-2-VOAIN",
					"component":           fmt.Sprintf("ATTENUATOR-1-%d-3", idx+1),
					"enabled":             "true",
					"egress-port":         "PORT-1-2-2-VOAOUT",
					"fix-attenuation":     "0.55",
					"attenuation-mode":    "CONSTANT_ATTENUATION",
				},
			},
		}
		loadDB(dbName, db.StateDB, pre_req_map)
		expected_get_json = fmt.Sprintf("{\"openconfig-optical-attenuator:state\":{\"attenuation\":\"2.5\",\"attenuation-mode\":\"openconfig-optical-attenuator:CONSTANT_ATTENUATION\",\"component\":\"ATTENUATOR-1-%d-3\",\"egress-port\":\"PORT-1-2-2-VOAOUT\",\"enabled\":true,\"ingress-port\":\"PORT-1-2-2-VOAIN\",\"name\":\"ATTENUATOR-1-%d-3\",\"target-output-power\":\"1.5\"}}", idx+1, idx+1)
		url = fmt.Sprintf("/openconfig-optical-attenuator:optical-attenuator/attenuators/attenuator[name=ATTENUATOR-1-%d-3]/state", idx+1)
		t.Run("Test get on StateDB Attnuator Key ", processGetRequest(url, nil, expected_get_json, false))
		time.Sleep(1 * time.Second)
		t.Log("\n\n+++++++++++++ Done Performing Get on stateDB Attnuator Key ++++++++++++")

		/* ATTENUATOR Counter DB test */
		t.Log("\n\n+++++++++++++ Performing Get on CounterDB Attenuator Key ++++++++++++")
		pre_req_map = map[string]interface{}{
			"ATTENUATOR": map[string]interface{}{
				fmt.Sprintf("ATTENUATOR-1-%d-3_OutputPowerTotal:15_pm_current", idx+1): map[string]interface{}{
					"instant":          "2.5",
					"min-time":         "1680168600865596372",
					"min":              "5.1",
					"max":              "10.2",
					"interval":         "900000000000",
					"validity":         "incomplete",
					"starttime":        "1680168600000000000",
					"max-time":         "1680168600865596372",
					"current_validity": "complete",
					"avg":              "5.5",
				},
			},
		}
		loadDB(dbName, db.CountersDB, pre_req_map)
		expected_get_json = "{\"openconfig-optical-attenuator:output-power-total\":{\"avg\":\"5.5\",\"instant\":\"2.5\",\"interval\":\"900000000000\",\"max\":\"10.2\",\"max-time\":\"1680168600865596372\",\"min\":\"5.1\",\"min-time\":\"1680168600865596372\"}}"
		url = fmt.Sprintf("/openconfig-optical-attenuator:optical-attenuator/attenuators/attenuator[name=ATTENUATOR-1-%d-3]/state/output-power-total", idx+1)
		t.Run("Test get on StateDB Attnuator Key-Xfmr and Field-Xfmr.", processGetRequest(url, nil, expected_get_json, false))
		time.Sleep(1 * time.Second)
		t.Log("\n\n+++++++++++++ Done Performing Get on CounterDB Attenuator Key ++++++++++++")

		/* Unload the Data */
		cleanuptbl = map[string]interface{}{
			"ATTENUATOR": map[string]interface{}{
				fmt.Sprintf("ATTENUATOR-1-%d-3", idx+1): "",
			},
		}
		unloadDB(dbName, db.ConfigDB, cleanuptbl)

		cleanupstatetbl = map[string]interface{}{
			"ATTENUATOR_TABLE": map[string]interface{}{
				fmt.Sprintf("ATTENUATOR-1-%d-3", idx+1): "",
			},
		}
		unloadDB(dbName, db.StateDB, cleanupstatetbl)

		cleanupcountertbl = map[string]interface{}{
			"ATTENUATOR": map[string]interface{}{
				fmt.Sprintf("ATTENUATOR-1-%d-3_OutputPowerTotal:15_pm_current", idx+1): "",
			},
		}
		unloadDB(dbName, db.CountersDB, cleanupcountertbl)
	}
}

func Test_set_attenuator_config_DB_key_and_field_xfmr(t *testing.T) {
	t.Log("\n\n+++++++++++++ Performing Create/Replace/Delete on platform ATTENUATOR ++++++++++++")

	for idx, dbName := range getMdbNames() {
		url := fmt.Sprintf("/openconfig-optical-attenuator:optical-attenuator/attenuators/attenuator[name=ATTENUATOR-1-%d-3]/config", idx+1)
		url_body_json := "{\"config\":{\"attenuation\":\"1.5\",\"attenuation-mode\":\"CONSTANT_ATTENUATION\",\"enabled\":true,\"target-output-power\":\"2.5\"}}"
		pre_req_map := map[string]interface{}{
			"ATTENUATOR": map[string]interface{}{
				fmt.Sprintf("ATTENUATOR-1-%d-3", idx+1): map[string]interface{}{
					"attenuation":         "12.5",
					"attenuation-mode":    "CONSTANT_POWER",
					"target-output-power": "12.15",
					"enabled":             "false",
				},
			},
		}
		loadDB(dbName, db.ConfigDB, pre_req_map)
		expected_map := map[string]interface{}{
			"ATTENUATOR": map[string]interface{}{
				fmt.Sprintf("ATTENUATOR-1-%d-3", idx+1): map[string]interface{}{
					"attenuation":         "1.5",
					"attenuation-mode":    "CONSTANT_ATTENUATION",
					"target-output-power": "2.5",
					"enabled":             "true",
				},
			},
		}
		time.Sleep(1 * time.Second)

		//PUT Test
		t.Run("Replace Test on ATTENUATOR yang", processSetRequest(url, url_body_json, "PUT", false))
		time.Sleep(1 * time.Second)

		//GET Test
		t.Run("Verify replace on ATTENUATOR yang", verifyDbResult(rclient[dbName], fmt.Sprintf("ATTENUATOR|ATTENUATOR-1-%d-3", idx+1), expected_map, false))
		time.Sleep(1 * time.Second)

		//Delete Test
		url = fmt.Sprintf("/openconfig-optical-attenuator:optical-attenuator/attenuators/attenuator[name=ATTENUATOR-1-%d-3]/config", idx+1)
		t.Run("DELETE Test on ATTENUATOR yang", processDeleteRequest(url, false))
		time.Sleep(1 * time.Second)

		//POST Test
		url = "/openconfig-optical-attenuator:optical-attenuator/attenuators"
		post_payload := fmt.Sprintf("{\"attenuator\":[{\"name\":\"ATTENUATOR-1-%d-21\",\"config\": {\"attenuation-mode\": \"CONSTANT_POWER\",\"target-output-power\": \"55.55\",\"attenuation\": \"44.44\",\"enabled\":false}}]}", idx+1)
		t.Run("Create test on ATTENUATOR yang", processSetRequest(url, post_payload, "POST", false))
		time.Sleep(1 * time.Second)

		expected_map = map[string]interface{}{
			"ATTENUATOR": map[string]interface{}{
				fmt.Sprintf("ATTENUATOR-1-%d-21", idx+1): map[string]interface{}{
					"attenuation":         "44.44",
					"attenuation-mode":    "CONSTANT_POWER",
					"target-output-power": "55.55",
					"enabled":             "false",
				},
			},
		}

		//Verify After the POST test
		t.Run("Verify After the POST test on ATTENUATOR yang", verifyDbResult(rclient[dbName], fmt.Sprintf("ATTENUATOR|ATTENUATOR-1-%d-21", idx+1), expected_map, false))

		//Delete After the POST test
		url = fmt.Sprintf("/openconfig-optical-attenuator:optical-attenuator/attenuators/attenuator[name=ATTENUATOR-1-%d-21]/config", idx+1)
		t.Run("DELETE After the POST test on ATTENUATOR yang", processDeleteRequest(url, false))
		time.Sleep(1 * time.Second)

		url = "/openconfig-optical-attenuator:optical-attenuator/attenuators"
		//Bulk POST Test
		Bulk_post_payload := fmt.Sprintf("{\"attenuator\":[{\"name\":\"ATTENUATOR-1-%d-22\",\"config\": {\"attenuation-mode\": \"CONSTANT_POWER\",\"target-output-power\": \"55.55\",\"attenuation\": \"44.44\",\"enabled\":false}},{\"name\":\"ATTENUATOR-1-%d-23\",\"config\": {\"attenuation-mode\": \"CONSTANT_POWER\",\"target-output-power\": \"5.5\",\"attenuation\": \"4.4\",\"enabled\":true}}]}", idx+1, idx+1)
		t.Run("Create test on ATTENUATOR yang", processSetRequest(url, Bulk_post_payload, "POST", false))

		expected_map = map[string]interface{}{
			"ATTENUATOR": map[string]interface{}{
				fmt.Sprintf("ATTENUATOR-1-%d-23", idx+1): map[string]interface{}{
					"attenuation":         "4.4",
					"attenuation-mode":    "CONSTANT_POWER",
					"target-output-power": "5.5",
					"enabled":             "true",
				},
			},
		}

		time.Sleep(1 * time.Second)
		//Verify After the Bulk test
		t.Run("Verify After the Bulk test on ATTENUATOR yang", verifyDbResult(rclient[dbName], fmt.Sprintf("ATTENUATOR|ATTENUATOR-1-%d-23", idx+1), expected_map, false))

		expected_map = map[string]interface{}{
			"ATTENUATOR": map[string]interface{}{
				fmt.Sprintf("ATTENUATOR-1-%d-22", idx+1): map[string]interface{}{
					"attenuation":         "44.44",
					"attenuation-mode":    "CONSTANT_POWER",
					"target-output-power": "55.55",
					"enabled":             "false",
				},
			},
		}

		time.Sleep(1 * time.Second)
		//Verify After the POST test
		t.Run("Verify After the Bulk test on ATTENUATOR yang", verifyDbResult(rclient[dbName], fmt.Sprintf("ATTENUATOR|ATTENUATOR-1-%d-22", idx+1), expected_map, false))

		//Delete keys After the Bulk Test
		url = fmt.Sprintf("/openconfig-optical-attenuator:optical-attenuator/attenuators/attenuator[name=ATTENUATOR-1-%d-22]/config", idx+1)
		t.Run("DELETE After the Bulk Test on ATTENUATOR yang", processDeleteRequest(url, false))
		time.Sleep(1 * time.Second)

		url = fmt.Sprintf("/openconfig-optical-attenuator:optical-attenuator/attenuators/attenuator[name=ATTENUATOR-1-%d-23]/config", idx+1)
		t.Run("DELETE After the Bulk Test on ATTENUATOR yang", processDeleteRequest(url, false))
		time.Sleep(1 * time.Second)
	}

	t.Log("\n\n+++++++++++++ Done Performing Create/Replace/Delete on ATTENUATOR ++++++++++++")
}

func Test_attenuator_allow_write_multiple_namespace(t *testing.T) {
	t.Log("\n\n+++++++++++++ Performing attenuator write for multiple namespace (should be allowed) ++++++++++++")

	url := "/openconfig-optical-attenuator:optical-attenuator/attenuators"

	url_body_json := `{
        "attenuator": [
            {
                "name": "ATTENUATOR-1-2-22",
                "config": {
                    "attenuation-mode": "CONSTANT_POWER",
                    "target-output-power": "55.55",
                    "attenuation": "44.44",
                    "enabled": false
                }
            },
            {
                "name": "ATTENUATOR-1-1-23",
                "config": {
                    "attenuation-mode": "CONSTANT_POWER",
                    "target-output-power": "5.5",
                    "attenuation": "4.4",
                    "enabled": true
                }
            }
        ]
    }`

	// false = no error expected
	t.Run("Test allow writing multiple namespace for attenuator", processSetRequest(url, url_body_json, "POST", false, nil))

	t.Log("\n\n+++++++++++++ Done Performing attenuator write for multiple namespace ++++++++++++")
}
