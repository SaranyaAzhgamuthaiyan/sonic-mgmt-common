//go:build xfmrtest
// +build xfmrtest

package transformer_test

import (
	"fmt"
	"github.com/Azure/sonic-mgmt-common/translib/db"
	"testing"
	"time"
)

/* APS config,counters, state DB fields and values */
func Test_aps_config_state_counter_DB_key(t *testing.T) {
	var pre_req_map, cleanuptbl, cleanupstatetbl, cleanupcountertbl map[string]interface{}
	var expected_get_json, url []string

	for idx, dbName := range getMdbNames() {
		/* APS Config DB test */
		t.Log("\n\n+++++++++++++ Performing Get on configDB APS Key ++++++++++++")
		pre_req_map = map[string]interface{}{
			"APS": map[string]interface{}{
				fmt.Sprintf("APS-1-%d-2", idx+1): map[string]interface{}{
					"force-to-port": "PRIMARY",
				},
			},
		}
		loadDB(dbName, db.ConfigDB, pre_req_map)
		expected_get_json = []string{fmt.Sprintf("{\"openconfig-transport-line-protection:config\":{\"force-to-port\":\"PRIMARY\",\"name\":\"APS-1-%d-2\"}}", idx+1)}
		url = []string{fmt.Sprintf("/openconfig-transport-line-protection:aps/aps-modules/aps-module[name=APS-1-%d-2]/config", idx+1)}

		t.Run("Test get on configDB APS Key", processGetRequest(url[0], nil, expected_get_json[0], false))
		time.Sleep(1 * time.Second)
		t.Log("\n\n+++++++++++++ Done Performing Get on configDB APS Key ++++++++++++")

		/* APS and APS_PORT State DB test */
		t.Log("\n\n+++++++++++++ Performing Get on stateDB APS Key ++++++++++++")
		pre_req_map = map[string]interface{}{
			"APS_TABLE": map[string]interface{}{
				fmt.Sprintf("APS-1-%d-2", idx+1): map[string]interface{}{
					"revertive":                        true,
					"wait-to-restore-time":             10,
					"hold-off-time":                    5,
					"primary-switch-threshold":         1.5,
					"primary-switch-hysteresis":        3.5,
					"secondary-switch-threshold":       2.4,
					"relative-switch-threshold":        4.5,
					"relative-switch-threshold-offset": 6.7,
					"force-to-port":                    "NONE",
					"active-path":                      "PRIMARY",
				},
			},
			"APS_PORT_TABLE": map[string]interface{}{
				fmt.Sprintf("APS-1-%d-2_CommonIn", idx+1): map[string]interface{}{
					"enabled":            true,
					"target-attenuation": 6.1,
					"attenuation":        10.5,
				},
			},
		}
		loadDB(dbName, db.StateDB, pre_req_map)
		expected_get_json = []string{
			fmt.Sprintf("{\"openconfig-transport-line-protection:state\":{\"active-path\":\"openconfig-transport-line-protection:PRIMARY\",\"force-to-port\":\"NONE\",\"hold-off-time\":5,\"name\":\"APS-1-%d-2\",\"primary-switch-hysteresis\":\"3.5\",\"primary-switch-threshold\":\"1.5\",\"relative-switch-threshold\":\"4.5\",\"relative-switch-threshold-offset\":\"6.7\",\"revertive\":true,\"secondary-switch-threshold\":\"2.4\",\"wait-to-restore-time\":10}}", idx+1),
			"{\"openconfig-transport-line-protection:state\":{\"attenuation\":\"10.5\",\"enabled\":true,\"target-attenuation\":\"6.1\"}}",
		}
		url = []string{
			fmt.Sprintf("/openconfig-transport-line-protection:aps/aps-modules/aps-module[name=APS-1-%d-2]/state", idx+1),
			fmt.Sprintf("/openconfig-transport-line-protection:aps/aps-modules/aps-module[name=APS-1-%d-2]/ports/common-in/state", idx+1),
		}

		for idx, url := range url {
			t.Run("Test get on state DB Key for "+url, processGetRequest(url, nil, expected_get_json[idx], false))
			time.Sleep(1 * time.Second)
		}
		t.Log("\n\n+++++++++++++ Done Performing Get on stateDB APS Key ++++++++++++")

		/* APS_PORT Counters DB test */
		t.Log("\n\n+++++++++++++ Performing Get on countersDB APS Key ++++++++++++")
		pre_req_map = map[string]interface{}{
			"APS_PORT": map[string]interface{}{
				fmt.Sprintf("APS-1-%d-2_CommonIn_OpticalPower:15_pm_current", idx+1): map[string]interface{}{
					"instant":          "0.55",
					"min-time":         "1680134400477338361",
					"min":              "0.55",
					"max":              "0.6",
					"interval":         "86400000000000",
					"validity":         "incomplete",
					"starttime":        "1680134400000000000",
					"max-time":         "1680136923116487196",
					"current_validity": "complete",
					"avg":              "0.55",
				},
			},
		}
		loadDB(dbName, db.CountersDB, pre_req_map)
		expected_get_json = []string{"{\"openconfig-transport-line-protection:optical-power\":{\"avg\":\"0.55\",\"instant\":\"0.55\",\"interval\":\"86400000000000\",\"max\":\"0.6\",\"max-time\":\"1680136923116487196\",\"min\":\"0.55\",\"min-time\":\"1680134400477338361\"}}"}
		url = []string{fmt.Sprintf("/openconfig-transport-line-protection:aps/aps-modules/aps-module[name=APS-1-%d-2]/ports/common-in/state/optical-power", idx+1)}

		t.Run("Test get on countersDB APS Key", processGetRequest(url[0], nil, expected_get_json[0], false))
		time.Sleep(1 * time.Second)
		t.Log("\n\n+++++++++++++ Done Performing Get on countersDB APS Key ++++++++++++")

		/* Unload the data */
		cleanuptbl = map[string]interface{}{
			"APS": map[string]interface{}{
				fmt.Sprintf("APS-1-%d-2", idx+1): "",
			},
		}
		unloadDB(dbName, db.ConfigDB, cleanuptbl)

		cleanupstatetbl = map[string]interface{}{
			"APS_TABLE":      map[string]interface{}{fmt.Sprintf("APS-1-%d-2", idx+1): ""},
			"APS_PORT_TABLE": map[string]interface{}{fmt.Sprintf("APS-1-%d-2_CommonIn", idx+1): ""},
		}
		unloadDB(dbName, db.StateDB, cleanupstatetbl)

		cleanupcountertbl = map[string]interface{}{
			"APS_PORT": map[string]interface{}{
				fmt.Sprintf("APS-1-%d-2_CommonIn_OpticalPower:15_pm_current", idx+1): "",
			},
		}
		unloadDB(dbName, db.CountersDB, cleanupcountertbl)
	}
}

func Test_set_aps_config_DB_key_and_field_xfmr(t *testing.T) {
	t.Log("\n\n+++++++++++++ Performing Create/Replace/Delete on aps ++++++++++++")

	for idx, dbName := range getMdbNames() {
		url := fmt.Sprintf("/openconfig-transport-line-protection:aps/aps-modules/aps-module[name=APS-1-%d-2]/config", idx+1)
		url_body_json := "{\"config\": {\"force-to-port\": \"SECONDARY\"}}"

		pre_req_map := map[string]interface{}{
			"APS": map[string]interface{}{
				fmt.Sprintf("APS-1-%d-2", idx+1): map[string]interface{}{
					"force-to-port": "PRIMARY",
				},
			},
		}

		loadDB(dbName, db.ConfigDB, pre_req_map)
		expected_map := map[string]interface{}{
			"APS": map[string]interface{}{
				fmt.Sprintf("APS-1-%d-2", idx+1): map[string]interface{}{
					"force-to-port": "SECONDARY",
				},
			},
		}

		time.Sleep(1 * time.Second)

		//PUT Test
		t.Run("Replace Test on APS yang", processSetRequest(url, url_body_json, "PUT", false))
		time.Sleep(1 * time.Second)

		//GET Test
		t.Run("Verify replace on APS yang", verifyDbResult(rclient[dbName], fmt.Sprintf("APS|APS-1-%d-2", idx+1), expected_map, false))
		time.Sleep(1 * time.Second)

		//Delete Test
		url = fmt.Sprintf("/openconfig-transport-line-protection:aps/aps-modules/aps-module[name=APS-1-%d-2]", idx+1)
		t.Run("DELETE Test on APS yang", processDeleteRequest(url, false))
		time.Sleep(1 * time.Second)

		//POST Test
		url = "/openconfig-transport-line-protection:aps/aps-modules"
		post_payload := fmt.Sprintf("{\"aps-module\":[{\"name\":\"APS-1-%d-13\",\"config\": {\"force-to-port\": \"SECONDARY\"}}]}", idx+1)
		t.Run("Create test on APS yang", processSetRequest(url, post_payload, "POST", false))
		time.Sleep(1 * time.Second)

		expected_map = map[string]interface{}{
			"APS": map[string]interface{}{
				fmt.Sprintf("APS-1-%d-13", idx+1): map[string]interface{}{
					"force-to-port": "SECONDARY",
				},
			},
		}

		//Verify After the POST test
		t.Run("Verify After the POST test on APS yang", verifyDbResult(rclient[dbName], fmt.Sprintf("APS|APS-1-%d-13", idx+1), expected_map, false))

		//Delete After the POST test
		url = fmt.Sprintf("/openconfig-transport-line-protection:aps/aps-modules/aps-module[name=APS-1-%d-13]", idx+1)
		t.Run("DELETE After the POST test on APS yang", processDeleteRequest(url, false))
		time.Sleep(1 * time.Second)

		url = "/openconfig-transport-line-protection:aps/aps-modules"

		//Bulk POST Test
		Bulk_post_payload := fmt.Sprintf("{\"aps-module\":[{\"name\":\"APS-1-%d-14\",\"config\": {\"force-to-port\": \"PRIMARY\"}},{\"name\":\"APS-1-%d-15\",\"config\": {\"force-to-port\": \"SECONDARY\"}}]}", idx+1, idx+1)
		t.Run("Create test on APS yang", processSetRequest(url, Bulk_post_payload, "POST", false))

		expected_map = map[string]interface{}{
			"APS": map[string]interface{}{
				fmt.Sprintf("APS-1-%d-14", idx+1): map[string]interface{}{
					"force-to-port": "PRIMARY",
				},
			},
		}

		time.Sleep(1 * time.Second)
		//Verify After the Bulk test
		t.Run("Verify After the Bulk test on APS yang", verifyDbResult(rclient[dbName], fmt.Sprintf("APS|APS-1-%d-14", idx+1), expected_map, false))

		expected_map = map[string]interface{}{
			"APS": map[string]interface{}{
				fmt.Sprintf("APS-1-%d-15", idx+1): map[string]interface{}{
					"force-to-port": "SECONDARY",
				},
			},
		}

		time.Sleep(1 * time.Second)
		//Verify After the POST test
		t.Run("Verify After the Bulk test on APS yang", verifyDbResult(rclient[dbName], fmt.Sprintf("APS|APS-1-%d-15", idx+1), expected_map, false))

		//Delete keys After the Bulk Test
		url = fmt.Sprintf("/openconfig-transport-line-protection:aps/aps-modules/aps-module[name=APS-1-%d-14]", idx+1)
		t.Run("DELETE After the Bulk Test on APS yang", processDeleteRequest(url, false))
		time.Sleep(1 * time.Second)

		url = fmt.Sprintf("/openconfig-transport-line-protection:aps/aps-modules/aps-module[name=APS-1-%d-15]", idx+1)
		t.Run("DELETE After the Bulk Test on APS yang", processDeleteRequest(url, false))
		time.Sleep(1 * time.Second)
	}
	t.Log("\n\n+++++++++++++ Done Performing Create/Replace/Delete on APS ++++++++++++")
}

func Test_aps_allow_write_multiple_namespace(t *testing.T) {
	t.Log("\n\n+++++++++++++ Performing aps write for multiple namespace (should be allowed) ++++++++++++")

	url := "/openconfig-transport-line-protection:aps/aps-modules"
	url_body_json := `{
        "aps-module": [
            {"name": "APS-1-1-4", "config": {"force-to-port": "PRIMARY"}},
            {"name": "APS-1-2-5", "config": {"force-to-port": "SECONDARY"}}
        ]
    }`

	// false = no error expected
	t.Run("Test allow writing multiple namespace for aps", processSetRequest(url, url_body_json, "POST", false, nil))

	t.Log("\n\n+++++++++++++ Done Performing aps write for multiple namespace ++++++++++++")
}
