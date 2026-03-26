//go:build xfmrtest
// +build xfmrtest

package transformer_test

import (
	"fmt"
	"github.com/Azure/sonic-mgmt-common/translib/db"
	"testing"
	"time"
)

/* PORT counters, state DB fields and values */
func Test_port_config_state_counter_DB_key(t *testing.T) {
	var pre_req_map, cleanuptbl, cleanupstatetbl, cleanupcountertbl map[string]interface{}
	var url string

	for idx, dbName := range getMdbNames() {

		/* PORT config DB data load */
		pre_req_map = map[string]interface{}{
			"OC_COMPONENT": map[string]interface{}{
				fmt.Sprintf("PORT-1-%d-L1IN", idx+1): map[string]interface{}{
					"name": fmt.Sprintf("PORT-1-%d-L1IN", idx+1),
				},
			},
		}
		loadDB(dbName, db.ConfigDB, pre_req_map)

		/* PORT State DB test */
		t.Log("\n\n+++++++++++++ Performing Get on stateDB PORT Key ++++++++++++")
		pre_req_map = map[string]interface{}{
			"PORT_TABLE": map[string]interface{}{
				fmt.Sprintf("PORT-1-%d-L1IN", idx+1): map[string]interface{}{
					"admin-state":       "ENABLED",
					"optical-port-type": "ADD",
				},
			},
		}
		loadDB(dbName, db.StateDB, pre_req_map)
		expected_get_json := "{\"openconfig-transport-line-common:state\":{\"admin-state\":\"ENABLED\",\"optical-port-type\":\"openconfig-transport-types:ADD\"}}"
		url = fmt.Sprintf("/openconfig-platform:components/component[name=PORT-1-%d-L1IN]/port/openconfig-transport-line-common:optical-port/state", idx+1)

		t.Run("Test get on state DB PORT Key", processGetRequest(url, nil, expected_get_json, false))
		time.Sleep(1 * time.Second)
		t.Log("\n\n+++++++++++++ Done Performing Get on stateDB PORT  Key ++++++++++++")

		/* PORT Counters DB test */
		t.Log("\n\n+++++++++++++ Performing Get on counterDB PORT Key ++++++++++++")
		pre_req_map = map[string]interface{}{
			"PORT": map[string]interface{}{
				fmt.Sprintf("PORT-1-%d-L1IN_InputPower:15_pm_current", idx+1): map[string]interface{}{
					"instant":          "1.5",
					"min-time":         "1680168600865596372",
					"min":              "5.4",
					"max":              "10.5",
					"interval":         "900000000000",
					"validity":         "incomplete",
					"starttime":        "1680168600000000000",
					"max-time":         "1680168600865596372",
					"current_validity": "complete",
					"avg":              "11.12",
				},
			},
		}
		loadDB(dbName, db.CountersDB, pre_req_map)
		expected_get_json = "{\"openconfig-transport-line-common:input-power\":{\"avg\":\"11.12\",\"instant\":\"1.5\",\"interval\":\"900000000000\",\"max\":\"10.5\",\"max-time\":\"1680168600865596372\",\"min\":\"5.4\",\"min-time\":\"1680168600865596372\"}}"
		url = fmt.Sprintf("/openconfig-platform:components/component[name=PORT-1-%d-L1IN]/port/openconfig-transport-line-common:optical-port/state/input-power", idx+1)

		t.Run("Test get on counter DB PORT Key", processGetRequest(url, nil, expected_get_json, false))
		time.Sleep(1 * time.Second)
		t.Log("\n\n+++++++++++++ Done Performing Get on counterDB PORT Key ++++++++++++")

		// Unload the Data
		cleanuptbl = map[string]interface{}{
			"OC_COMPONENT": map[string]interface{}{
				fmt.Sprintf("PORT-1-%d-L1IN", idx+1): "",
			},
		}
		unloadDB(dbName, db.ConfigDB, cleanuptbl)
		cleanupstatetbl = map[string]interface{}{
			"PORT_TABLE": map[string]interface{}{
				fmt.Sprintf("PORT-1-%d-L1IN", idx+1): "",
			},
		}
		unloadDB(dbName, db.StateDB, cleanupstatetbl)
		cleanupcountertbl = map[string]interface{}{
			"PORT": map[string]interface{}{
				fmt.Sprintf("PORT-1-%d-L1IN_InputPower:15_pm_current", idx+1): "",
			},
		}
		unloadDB(dbName, db.CountersDB, cleanupcountertbl)
	}
}
