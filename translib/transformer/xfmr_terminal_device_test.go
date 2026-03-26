//go:build xfmrtest
// +build xfmrtest

package transformer_test

import (
	"fmt"
	"github.com/Azure/sonic-mgmt-common/translib/db"
	"testing"
	"time"
)

/* terminal device counters, state DB fields and values */
func Test_terminal_device_config_state_counter_DB_key(t *testing.T) {
	var pre_req_map, cleanuptbl, cleanupstatetbl, cleanupcountertbl map[string]interface{}
	var expected_get_json, url []string

	for idx, dbName := range getMdbNames() {
		/* terminal device config DB test */
		t.Log("\n\n+++++++++++++ Performing Get on configDB terminal device Key ++++++++++++")
		pre_req_map = map[string]interface{}{
			"OC_COMPONENT": map[string]interface{}{
				fmt.Sprintf("OCH-1-%d-L2", idx+1): map[string]interface{}{
					"name": fmt.Sprintf("OCH-1-%d-L2", idx+1),
				},
			},
			"OCH": map[string]interface{}{
				fmt.Sprintf("OCH-1-%d-L2", idx+1): map[string]interface{}{
					"operational-mode":    "test_data",
					"frequency":           "2",
					"target-output-power": "0.0",
				},
			},
			"LOGICAL_CHANNEL": map[string]interface{}{
				fmt.Sprintf("CH%d02", idx+1): map[string]interface{}{
					"test-signal":   "false",
					"loopback-mode": "NONE",
					"admin-state":   "MAINT",
				},
			},
			"ETHERNET": map[string]interface{}{
				fmt.Sprintf("CH%d02", idx+1): map[string]interface{}{
					"index":      fmt.Sprintf("%d02", idx+1),
					"client-als": "LASER_SHUTDOWN",
					"als-delay":  1000, // float to string

				},
			},
			"LLDP": map[string]interface{}{
				fmt.Sprintf("CH%d02", idx+1): map[string]interface{}{
					"index":    fmt.Sprintf("%d02", idx+1),
					"enabled":  "true",
					"snooping": "true",
				},
			},
		}
		loadDB(dbName, db.ConfigDB, pre_req_map)
		time.Sleep(1 * time.Second)
		expected_get_json = []string{fmt.Sprintf("{\"openconfig-terminal-device:config\":{\"admin-state\":\"MAINT\",\"index\":%d02,\"loopback-mode\":\"NONE\",\"test-signal\":false}}", idx+1), "{\"openconfig-terminal-device:config\":{\"als-delay\":1000,\"client-als\":\"LASER_SHUTDOWN\"}}"}
		url = []string{
			fmt.Sprintf("/openconfig-terminal-device:terminal-device/logical-channels/channel[index=%d02]/config", idx+1),
			fmt.Sprintf("/openconfig-terminal-device:terminal-device/logical-channels/channel[index=%d02]/ethernet/config", idx+1),
		}

		for idx, url := range url {
			t.Run("Test get on config DB Key for "+url, processGetRequest(url, nil, expected_get_json[idx], false))
			time.Sleep(1 * time.Second)
		}
		t.Log("\n\n+++++++++++++ Done Performing Get on configDB terminal device Key ++++++++++++")

		/* terminal device state DB test */
		t.Log("\n\n+++++++++++++ Performing Get on stateDB terminal device Key ++++++++++++")
		pre_req_map = map[string]interface{}{
			"LOGICAL_CHANNEL_TABLE": map[string]interface{}{
				fmt.Sprintf("CH%d02", idx+1): map[string]interface{}{
					"index":         fmt.Sprintf("%d02", idx+1),
					"rate-class":    "TRIB_RATE_200G",
					"description":   "GE-1-1-C3",
					"loopback-mode": "NONE",
					"test-signal":   "false",
					"admin-state":   "MAINT",
					"trib-protocol": "PROT_1GE",
					"link-state":    "DOWN",
				},
			},
			"ETHERNET_TABLE": map[string]interface{}{
				fmt.Sprintf("CH%d02", idx+1): map[string]interface{}{
					"index":      fmt.Sprintf("%d02", idx+1),
					"client-als": "LASER_SHUTDOWN",
					"als-delay":  1000,
				},
			},
			"LLDP_TABLE": map[string]interface{}{
				fmt.Sprintf("CH%d02", idx+1): map[string]interface{}{
					"index":    fmt.Sprintf("%d02", idx+1),
					"enabled":  "true",
					"snooping": "true",
				},
			},
			"NEIGHBOR_TABLE": map[string]interface{}{
				fmt.Sprintf("CH%d02|1", idx+1): map[string]interface{}{
					"id":                      "1",
					"system-name":             "test_sysname",
					"system-description":      "test_desc",
					"chassis-id":              "test_chassis_id",
					"chassis-id-type":         "PORT_COMPONENT",
					"age":                     "12",
					"last-update":             "12",
					"ttl":                     "100",
					"port-id":                 "test_portid",
					"port-id-type":            "PORT_COMPONENT",
					"port-description":        "test_pdes",
					"management-address":      "test_mg_add",
					"management-address-type": "test_mg_type",
				},
			},
			"ASSIGNMENT_TABLE": map[string]interface{}{
				fmt.Sprintf("CH%d02|ASS1", idx+1): map[string]interface{}{
					"index":                fmt.Sprintf("%d02", idx+1),
					"assignment-type":      "LOGICAL_CHANNEL",
					"description":          "ASSIGNMENT-1-1-18",
					"tributary-slot-index": "1",
					"mapping":              "GMP",
					"allocation":           "400.0",
					"logical-channel":      "102",
					"optical-channel":      "102",
				},
			},
			"MODE_TABLE": map[string]interface{}{
				fmt.Sprintf("%d02", idx+1): map[string]interface{}{
					"mode-id":     fmt.Sprintf("%d02", idx+1),
					"description": "600G::DP-64QAM::FEC-15::63.054::BICHM[0]",
					"vendor-id":   "Accelink",
				},
			},
		}
		loadDB(dbName, db.StateDB, pre_req_map)
		time.Sleep(1 * time.Second)
		expected_get_json = []string{
			fmt.Sprintf("{\"openconfig-terminal-device:state\":{\"admin-state\":\"MAINT\",\"description\":\"GE-1-1-C3\",\"index\":%d02,\"link-state\":\"DOWN\",\"loopback-mode\":\"NONE\",\"rate-class\":\"openconfig-transport-types:TRIB_RATE_200G\",\"test-signal\":false,\"trib-protocol\":\"openconfig-transport-types:PROT_1GE\"}}", idx+1),

			fmt.Sprintf("{\"openconfig-terminal-device:state\":{\"description\":\"600G::DP-64QAM::FEC-15::63.054::BICHM[0]\",\"mode-id\":%d02,\"vendor-id\":\"Accelink\"}}", idx+1),

			"{\"openconfig-terminal-device:state\":{\"age\":\"12\",\"chassis-id\":\"test_chassis_id\",\"chassis-id-type\":\"PORT_COMPONENT\",\"id\":\"1\",\"last-update\":\"12\",\"management-address\":\"test_mg_add\",\"management-address-type\":\"test_mg_type\",\"port-description\":\"test_pdes\",\"port-id\":\"test_portid\",\"port-id-type\":\"PORT_COMPONENT\",\"system-description\":\"test_desc\",\"system-name\":\"test_sysname\",\"ttl\":100}}",
		}
		url = []string{
			fmt.Sprintf("/openconfig-terminal-device:terminal-device/logical-channels/channel[index=%d02]/state", idx+1),
			fmt.Sprintf("/openconfig-terminal-device:terminal-device/operational-modes/mode[mode-id=%d02]/state", idx+1),
			fmt.Sprintf("/openconfig-terminal-device:terminal-device/logical-channels/channel[index=%d02]/ethernet/lldp/neighbors/neighbor[id=1]/state", idx+1),
		}
		for idx, url := range url {
			t.Run("Test get on state DB Key for "+url, processGetRequest(url, nil, expected_get_json[idx], false))
			time.Sleep(1 * time.Second)
		}
		t.Log("\n\n+++++++++++++ Done Performing Get on configDB terminal device Key ++++++++++++")

		/* terminal monitor Counters DB test */
		t.Log("\n\n+++++++++++++ Performing Get on counterDB terminal device Key ++++++++++++")
		pre_req_map = map[string]interface{}{
			"ETHERNET": map[string]interface{}{
				fmt.Sprintf("CH%d02:15_pm_current", idx+1): map[string]interface{}{
					"in-pcs-errored-seconds":          "0",
					"out-pcs-bip-errors":              "0",
					"in-block-errors":                 "0",
					"out-block-errors":                "0",
					"in-pcs-severely-errored-seconds": "0",
					"in-pcs-bip-errors":               "0",
					"in-jabber-frames":                "0",
					"in-fragment-frames":              "0",
					"in-oversize-frames":              "0",
					"out-crc-errors":                  "0",
					"in-pcs-unavailable-seconds":      "0",
					"in-undersize-frames":             "0",
					"in-crc-errors":                   "0",
				},
			},
		}
		loadDB(dbName, db.CountersDB, pre_req_map)
		expected_get_json = []string{"{\"openconfig-terminal-device:state\":{\"als-delay\":1000,\"client-als\":\"LASER_SHUTDOWN\",\"in-block-errors\":\"0\",\"in-fragment-frames\":\"0\",\"in-jabber-frames\":\"0\",\"in-pcs-bip-errors\":\"0\",\"in-pcs-errored-seconds\":\"0\",\"in-pcs-severely-errored-seconds\":\"0\",\"out-block-errors\":\"0\",\"out-pcs-bip-errors\":\"0\"}}"}
		url = []string{fmt.Sprintf("/openconfig-terminal-device:terminal-device/logical-channels/channel[index=%d02]/ethernet/state", idx+1)}
		t.Run("Test get on counter DB Ethernet Key-Xfmr and Field-Xfmr.", processGetRequest(url[0], nil, expected_get_json[0], false))
		time.Sleep(1 * time.Second)
		t.Log("\n\n+++++++++++++ Done Performing Get on counterDB terminal device Key ++++++++++++")

		//Unload the Data
		cleanuptbl = map[string]interface{}{
			"OCH": map[string]interface{}{
				fmt.Sprintf("OCH-1-%d-L2", idx+1): "",
			},
			"OC_COMPONENT": map[string]interface{}{
				fmt.Sprintf("OCH-1-%d-L2", idx+1): "",
			},
			"LOGICAL_CHANNEL": map[string]interface{}{
				fmt.Sprintf("CH%d02", idx+1): "",
			},
			"ETHERNET": map[string]interface{}{
				fmt.Sprintf("CH%d02", idx+1): "",
			},
			"LLDP": map[string]interface{}{
				fmt.Sprintf("CH%d02", idx+1): "",
			},
		}
		unloadDB(dbName, db.ConfigDB, cleanuptbl)

		cleanupstatetbl = map[string]interface{}{
			"LOGICAL_CHANNEL_TABLE": map[string]interface{}{
				fmt.Sprintf("CH%d02", idx+1): "",
			},
			"ETHERNET_TABLE": map[string]interface{}{
				fmt.Sprintf("CH%d02", idx+1): "",
			},
			"LLDP_TABLE": map[string]interface{}{
				fmt.Sprintf("CH%d02", idx+1): "",
			},
			"NEIGHBOR_TABLE": map[string]interface{}{
				fmt.Sprintf("CH%d02|1", idx+1): "",
			},
			"ASSIGNMENT_TABLE": map[string]interface{}{
				fmt.Sprintf("CH%d02|ASS1", idx+1): "",
			},
			"MODE_TABLE": map[string]interface{}{
				fmt.Sprintf("%d02", idx+1): "",
			},
		}
		unloadDB(dbName, db.StateDB, cleanupstatetbl)

		cleanupcountertbl = map[string]interface{}{
			"ETHERNET": map[string]interface{}{
				fmt.Sprintf("CH%d02:15_pm_current", idx+1): "",
			},
		}
		unloadDB(dbName, db.CountersDB, cleanupcountertbl)
	}
}

func Test_set_terminal_device_config_DB_key_and_field_xfmr(t *testing.T) {
	t.Log("\n\n+++++++++++++ Performing Create/Replace/Delete on terminal_device ++++++++++++")

	for idx, dbName := range getMdbNames() {
		url := fmt.Sprintf("/openconfig-terminal-device:terminal-device/logical-channels/channel[index=%d02]/config", idx+1)
		url_body_json := "{\"config\":{\"loopback-mode\":\"NONE\",\"admin-state\":\"ENABLED\",\"test-signal\":true}}"

		pre_req_map := map[string]interface{}{
			"LOGICAL_CHANNEL": map[string]interface{}{
				fmt.Sprintf("CH%d02", idx+1): map[string]interface{}{
					"test-signal":   "false",
					"loopback-mode": "NONE",
					"admin-state":   "MAINT",
				},
			},
		}
		loadDB(dbName, db.ConfigDB, pre_req_map)
		expected_map := map[string]interface{}{
			"LOGICAL_CHANNEL": map[string]interface{}{
				fmt.Sprintf("CH%d02", idx+1): map[string]interface{}{
					"test-signal":   "true",
					"loopback-mode": "NONE",
					"admin-state":   "ENABLED",
				},
			},
		}
		time.Sleep(1 * time.Second)

		//PUT Test
		t.Run("Replace Test on terminal_device yang", processSetRequest(url, url_body_json, "PUT", false))
		time.Sleep(1 * time.Second)

		//GET Test
		t.Run("Verify replace on terminal_device yang", verifyDbResult(rclient[dbName], fmt.Sprintf("LOGICAL_CHANNEL|CH%d02", idx+1), expected_map, false))
		time.Sleep(1 * time.Second)

		//Delete Test
		url = fmt.Sprintf("/openconfig-terminal-device:terminal-device/logical-channels/channel[index=%d02]/config", idx+1)
		t.Run("DELETE Test on terminal_device yang", processDeleteRequest(url, false))
		time.Sleep(1 * time.Second)

		//POST Test
		url = "/openconfig-terminal-device:terminal-device/logical-channels"
		post_payload := fmt.Sprintf("{\"channel\":[{\"index\":%d03,\"config\":{\"loopback-mode\":\"NONE\",\"admin-state\":\"ENABLED\",\"test-signal\":true}}]}", idx+1)
		t.Run("Create test on terminal_device yang", processSetRequest(url, post_payload, "POST", false))
		time.Sleep(1 * time.Second)

		expected_map = map[string]interface{}{
			"LOGICAL_CHANNEL": map[string]interface{}{
				fmt.Sprintf("CH%d03", idx+1): map[string]interface{}{
					"test-signal":   "true",
					"loopback-mode": "NONE",
					"admin-state":   "ENABLED",
				},
			},
		}

		//Verify After the POST test
		t.Run("Verify After the POST test on terminal_device yang", verifyDbResult(rclient[dbName], fmt.Sprintf("LOGICAL_CHANNEL|CH%d03", idx+1), expected_map, false))

		//Delete After the POST test
		url = fmt.Sprintf("/openconfig-terminal-device:terminal-device/logical-channels/channel[index=%d03]/config", idx+1)
		t.Run("DELETE After the POST test on terminal_device yang", processDeleteRequest(url, false))
		time.Sleep(1 * time.Second)

		url = "/openconfig-terminal-device:terminal-device/logical-channels"
		//Bulk POST Test
		Bulk_post_payload := fmt.Sprintf("{\"channel\":[{\"index\":%d04,\"config\":{\"loopback-mode\":\"NONE\",\"admin-state\":\"ENABLED\",\"test-signal\":true}},{\"index\":%d05,\"config\":{\"loopback-mode\":\"NONE\",\"admin-state\":\"ENABLED\",\"test-signal\":false}}]}", idx+1, idx+1)
		t.Run("Create test on terminal_device yang", processSetRequest(url, Bulk_post_payload, "POST", false))

		expected_map = map[string]interface{}{
			"LOGICAL_CHANNEL": map[string]interface{}{
				fmt.Sprintf("CH%d04", idx+1): map[string]interface{}{
					"test-signal":   "true",
					"loopback-mode": "NONE",
					"admin-state":   "ENABLED",
				},
			},
		}

		time.Sleep(1 * time.Second)
		//Verify After the Bulk test
		t.Run("Verify After the Bulk test on terminal_device yang", verifyDbResult(rclient[dbName], fmt.Sprintf("LOGICAL_CHANNEL|CH%d04", idx+1), expected_map, false))

		expected_map = map[string]interface{}{
			"LOGICAL_CHANNEL": map[string]interface{}{
				fmt.Sprintf("CH%d05", idx+1): map[string]interface{}{
					"test-signal":   "false",
					"loopback-mode": "NONE",
					"admin-state":   "ENABLED",
				},
			},
		}

		time.Sleep(1 * time.Second)
		//Verify After the POST test
		t.Run("Verify After the Bulk test on terminal_device yang", verifyDbResult(rclient[dbName], fmt.Sprintf("LOGICAL_CHANNEL|CH%d05", idx+1), expected_map, false))

		//Delete keys After the Bulk Test
		url = fmt.Sprintf("/openconfig-terminal-device:terminal-device/logical-channels/channel[index=%d04]/config", idx+1)
		t.Run("DELETE After the Bulk Test on terminal_device yang", processDeleteRequest(url, false))
		time.Sleep(1 * time.Second)

		url = fmt.Sprintf("/openconfig-terminal-device:terminal-device/logical-channels/channel[index=%d05]/config", idx+1)
		t.Run("DELETE After the Bulk Test on terminal_device yang", processDeleteRequest(url, false))
		time.Sleep(1 * time.Second)
	}

	t.Log("\n\n+++++++++++++ Done Performing Create/Replace/Delete on terminal_device ++++++++++++")
}

func Test_terminal_device_allow_write_multiple_namespace(t *testing.T) {
	t.Log("\n\n+++++++++++++ Performing terminal_device write for multiple namespace (should be allowed) ++++++++++++")

	url := "/openconfig-terminal-device:terminal-device/logical-channels"
	url_body_json := `{
        "channel": [
            {
                "index": 104,
                "config": {
                    "loopback-mode": "NONE",
                    "admin-state": "ENABLED",
                    "test-signal": true
                }
            },
            {
                "index": 205,
                "config": {
                    "loopback-mode": "NONE",
                    "admin-state": "ENABLED",
                    "test-signal": false
                }
            }
        ]
    }`

	// false = no error expected
	t.Run("Test allow writing multiple namespace for terminal_device", processSetRequest(url, url_body_json, "POST", false, nil))

	t.Log("\n\n+++++++++++++ Done Performing terminal_device write for multiple namespace ++++++++++++")
}
