//go:build xfmrtest
// +build xfmrtest

package transformer_test

import (
	"fmt"
	"github.com/Azure/sonic-mgmt-common/translib/db"
	"testing"
	"time"
)

/* channel_monitor config,counters, state DB fields and values */
func Test_channel_monitor_config_state_DB_key(t *testing.T) {
	var pre_req_map, cleanuptbl, cleanupstatetbl map[string]interface{}
	var expected_get_json, url []string

	for idx, dbName := range getMdbNames() {

		/* channel_monitor config DB test */
		t.Log("\n\n+++++++++++++ Performing Get on configDB channel monitor Key ++++++++++++")
		pre_req_map = map[string]interface{}{
			"OCM": map[string]interface{}{
				fmt.Sprintf("OCM-1-%d-2", idx+1): map[string]interface{}{
					"monitor-port": "test_monitor_port",
				},
			},
		}
		loadDB(dbName, db.ConfigDB, pre_req_map)
		expected_get_json = []string{
			fmt.Sprintf("{\"openconfig-channel-monitor:config\":{\"monitor-port\":\"test_monitor_port\",\"name\":\"OCM-1-%d-2\"}}", idx+1),
		}
		url = []string{
			fmt.Sprintf("/openconfig-channel-monitor:channel-monitors/channel-monitor[name=OCM-1-%d-2]/config", idx+1),
		}

		t.Run("Test get on config DB channel monitor Key", processGetRequest(url[0], nil, expected_get_json[0], false))
		time.Sleep(1 * time.Second)
		t.Log("\n\n+++++++++++++ Done Performing Get on configDB channel monitor Key ++++++++++++")

		/* channel monitor State DB test */
		t.Log("\n\n+++++++++++++ Performing Get on stateDB channel monitor Key ++++++++++++")
		pre_req_map = map[string]interface{}{
			"OCM_TABLE": map[string]interface{}{
				fmt.Sprintf("OCM-1-%d-2", idx+1): map[string]interface{}{
					"monitor-port": "test_monitor_port",
				},
			},
			"OCM_CHANNEL_TABLE": map[string]interface{}{
				fmt.Sprintf("OCM-1-%d-2|25|100", idx+1): map[string]interface{}{
					"power": "5.0",
				},
			},
		}
		loadDB(dbName, db.StateDB, pre_req_map)
		expected_get_json = []string{
			fmt.Sprintf("{\"openconfig-channel-monitor:state\":{\"monitor-port\":\"test_monitor_port\",\"name\":\"OCM-1-%d-2\"}}", idx+1),
			fmt.Sprintf("{\"openconfig-channel-monitor:state\":{\"lower-frequency\":\"25\",\"power\":\"5\",\"upper-frequency\":\"100\"}}"),
		}
		url = []string{
			fmt.Sprintf("/openconfig-channel-monitor:channel-monitors/channel-monitor[name=OCM-1-%d-2]/state", idx+1),
			fmt.Sprintf("/openconfig-channel-monitor:channel-monitors/channel-monitor[name=OCM-1-%d-2]/channels/channel[lower-frequency=25][upper-frequency=100]/state", idx+1),
		}

		for idx, url := range url {
			t.Run("Test get on state DB channel monitor Key for "+url, processGetRequest(url, nil, expected_get_json[idx], false))
			time.Sleep(1 * time.Second)
		}
		t.Log("\n\n+++++++++++++ Done Performing Get on stateDB channel monitor Key ++++++++++++")

		/* Unload the Data */
		cleanuptbl = map[string]interface{}{
			"OCM": map[string]interface{}{
				fmt.Sprintf("OCM-1-%d-2", idx+1): "",
			},
		}
		unloadDB(dbName, db.ConfigDB, cleanuptbl)
		cleanupstatetbl = map[string]interface{}{
			"OCM_TABLE": map[string]interface{}{
				fmt.Sprintf("OCM-1-%d-2", idx+1): "",
			},
			"OCM_CHANNEL_TABLE": map[string]interface{}{
				fmt.Sprintf("OCM-1-%d-2|25|100", idx+1): "",
			},
		}
		unloadDB(dbName, db.StateDB, cleanupstatetbl)
	}
}

func Test_set_channel_monitor_config_DB_key_and_field_xfmr(t *testing.T) {
	t.Log("\n\n+++++++++++++ Performing Create/Replace/Delete on channel monitor ++++++++++++")
	for idx, dbName := range getMdbNames() {
		url := fmt.Sprintf("/openconfig-channel-monitor:channel-monitors/channel-monitor[name=OCM-1-%d-2]/config", idx+1)
		url_body_json := "{\"config\": {\"monitor-port\": \"monitor_data\"}}"

		pre_req_map := map[string]interface{}{
			"OCM": map[string]interface{}{
				fmt.Sprintf("OCM-1-%d-2", idx+1): map[string]interface{}{
					"monitor-port": "test_monitor_port",
				},
			},
		}

		loadDB(dbName, db.ConfigDB, pre_req_map)
		expected_map := map[string]interface{}{
			"OCM": map[string]interface{}{
				fmt.Sprintf("OCM-1-%d-2", idx+1): map[string]interface{}{
					"monitor-port": "monitor_data",
				},
			},
		}

		time.Sleep(1 * time.Second)

		//PUT Test
		t.Run("Replace Test on Channel Monitor yang", processSetRequest(url, url_body_json, "PUT", false))
		time.Sleep(1 * time.Second)

		//GET Test
		t.Run("Verify replace on Channel Monitor yang", verifyDbResult(rclient[dbName], fmt.Sprintf("OCM|OCM-1-%d-2", idx+1), expected_map, false))
		time.Sleep(1 * time.Second)

		//Delete Test
		url = fmt.Sprintf("/openconfig-channel-monitor:channel-monitors/channel-monitor[name=OCM-1-%d-2]/config", idx+1)
		t.Run("DELETE Test on Channel Monitor yang", processDeleteRequest(url, false))
		time.Sleep(1 * time.Second)

		//POST Test
		url = "/openconfig-channel-monitor:channel-monitors"
		post_payload := fmt.Sprintf("{\"channel-monitor\":[{\"name\":\"OCM-1-%d-13\",\"config\":{\"monitor-port\":\"test_port_data\"}}]}", idx+1)
		t.Run("Create test on Channel Monitor yang", processSetRequest(url, post_payload, "POST", false))
		time.Sleep(1 * time.Second)

		expected_map = map[string]interface{}{
			"OCM": map[string]interface{}{
				fmt.Sprintf("OCM-1-%d-13", idx+1): map[string]interface{}{
					"monitor-port": "test_port_data",
				},
			},
		}

		//Verify After the POST test
		t.Run("Verify After the POST test on Channel Monitor yang", verifyDbResult(rclient[dbName], fmt.Sprintf("OCM|OCM-1-%d-13", idx+1), expected_map, false))

		//Delete After the POST test
		url = fmt.Sprintf("/openconfig-channel-monitor:channel-monitors/channel-monitor[name=OCM-1-%d-13]/config", idx+1)
		t.Run("DELETE After the POST test on Channel Monitor yang", processDeleteRequest(url, false))
		time.Sleep(1 * time.Second)

		url = "/openconfig-channel-monitor:channel-monitors"

		//Bulk POST Test
		Bulk_post_payload := fmt.Sprintf("{\"channel-monitor\":[{\"name\":\"OCM-1-%d-14\",\"config\":{\"monitor-port\":\"Port_data_1\"}},{\"name\":\"OCM-1-%d-15\",\"config\":{\"monitor-port\":\"Port_data_2\"}}]}", idx+1, idx+1)

		t.Run("Create test on Channel monitor yang", processSetRequest(url, Bulk_post_payload, "POST", false))

		expected_map = map[string]interface{}{
			"OCM": map[string]interface{}{
				fmt.Sprintf("OCM-1-%d-14", idx+1): map[string]interface{}{
					"monitor-port": "Port_data_1",
				},
			},
		}

		time.Sleep(1 * time.Second)
		//Verify After the Bulk test
		t.Run("Verify After the Bulk test on Channel Monitor yang", verifyDbResult(rclient[dbName], fmt.Sprintf("OCM|OCM-1-%d-14", idx+1), expected_map, false))

		expected_map = map[string]interface{}{
			"OCM": map[string]interface{}{
				fmt.Sprintf("OCM-1-%d-15", idx+1): map[string]interface{}{
					"monitor-port": "Port_data_2",
				},
			},
		}

		time.Sleep(1 * time.Second)
		//Verify After the POST test
		t.Run("Verify After the Bulk test on Channel Monitor yang", verifyDbResult(rclient[dbName], fmt.Sprintf("OCM|OCM-1-%d-15", idx+1), expected_map, false))

		//Delete keys After the Bulk Test
		url = fmt.Sprintf("/openconfig-channel-monitor:channel-monitors/channel-monitor[name=OCM-1-%d-14]/config", idx+1)
		t.Run("DELETE After the Bulk Test on Channel Monitor yang", processDeleteRequest(url, false))
		time.Sleep(1 * time.Second)

		url = fmt.Sprintf("/openconfig-channel-monitor:channel-monitors/channel-monitor[name=OCM-1-%d-15]/config", idx+1)
		t.Run("DELETE After the Bulk Test on Channel Monitor yang", processDeleteRequest(url, false))
		time.Sleep(1 * time.Second)
	}

	t.Log("\n\n+++++++++++++ Done Performing Create/Replace/Delete on Channel Monitor ++++++++++++")
}

func Test_channel_monitor_allow_write_multiple_namespace(t *testing.T) {
	t.Log("\n\n+++++++++++++ Performing channel monitor write for multiple namespace (should be allowed) ++++++++++++")

	url := "/openconfig-channel-monitor:channel-monitors"
	url_body_json := `{
        "channel-monitor": [
            {
                "name": "OCM-1-2-14",
                "config": {
                    "monitor-port": "Port_data_1"
                }
            },
            {
                "name": "OCM-1-3-15",
                "config": {
                    "monitor-port": "Port_data_2"
                }
            }
        ]
    }`

	// false = no error expected
	t.Run("Test allow writing multiple namespace for channel monitor", processSetRequest(url, url_body_json, "POST", false, nil))

	t.Log("\n\n+++++++++++++ Done Performing channel monitor write for multiple namespace ++++++++++++")
}
