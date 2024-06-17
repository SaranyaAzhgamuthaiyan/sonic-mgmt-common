package transformer_test

import (
	"testing"
	"time"

	"github.com/Azure/sonic-mgmt-common/translib/db"
)

/* PORT counters, state DB fields and values */
func Test_port_config_state_counter_DB_key_and_field_xfmr(t *testing.T) {
	var pre_req_map, cleanuptbl, cleanupstatetbl, cleanupcountertbl map[string]interface{}
	var url string

	/* PORT config DB data load */
	pre_req_map = map[string]interface{}{"OC_COMPONENT_TABLE": map[string]interface{}{"PORT-1-1-L1IN": map[string]interface{}{"name": "PORT-1-1-L1IN"}}}
	loadDB(db.ConfigDB, pre_req_map)

	/* PORT State DB test */
	t.Log("\n\n+++++++++++++ Performing Get on stateDB PORT Key-Xfmr and Field-Xfmr ++++++++++++")
	pre_req_map = map[string]interface{}{"PORT": map[string]interface{}{"PORT-1-1-L1IN": map[string]interface{}{"admin-state": "ENABLED", "optical-port-type": "ADD"}}}
	loadDB(db.StateDB, pre_req_map)
	expected_get_json := "{\"openconfig-transport-line-common:state\":{\"admin-state\":\"ENABLED\",\"optical-port-type\":\"openconfig-transport-types:ADD\"}}"
	url = "/openconfig-platform:components/component[name=PORT-1-1-L1IN]/port/openconfig-transport-line-common:optical-port/state"
	t.Run("Test get on state DB PORT Key-Xfmr and Field-Xfmr.", processGetRequest(url, nil, expected_get_json, false))
	time.Sleep(1 * time.Second)
	t.Log("\n\n+++++++++++++ Done Performing Get on stateDB PORT  Key-Xfmr and Field-Xfmr ++++++++++++")

	/* PORT Counters DB test */
	t.Log("\n\n+++++++++++++ Performing Get on counterDB PORT Key-Xfmr and Field-Xfmr ++++++++++++")
	pre_req_map = map[string]interface{}{"PORT": map[string]interface{}{"PORT-1-1-L1IN_InputPower:15_pm_current": map[string]interface{}{"instant": "0.0", "min-time": "1680168600865596372", "min": "0.0", "max": "0.0", "interval": "900000000000", "validity": "incomplete", "starttime": "1680168600000000000", "max-time": "1680168600865596372", "current_validity": "complete", "avg": "0.0"}}}
	loadDB(db.CountersDB, pre_req_map)
	expected_get_json = "{\"openconfig-transport-line-common:input-power\":{\"avg\":\"0\",\"instant\":\"0\",\"interval\":\"900000000000\",\"max\":\"0\",\"max-time\":\"1680168600865596372\",\"min\":\"0\",\"min-time\":\"1680168600865596372\"}}"
	url = "/openconfig-platform:components/component[name=PORT-1-1-L1IN]/port/openconfig-transport-line-common:optical-port/state/input-power"
	t.Run("Test get on counter DB PORT Key-Xfmr and Field-Xfmr.", processGetRequest(url, nil, expected_get_json, false))
	time.Sleep(1 * time.Second)
	t.Log("\n\n+++++++++++++ Done Performing Get on counterDB PORT  Key-Xfmr and Field-Xfmr ++++++++++++")

	cleanuptbl = map[string]interface{}{"OC_COMPONENT_TABLE": map[string]interface{}{"PORT-1-1-L1IN": ""}}
	unloadDB(db.ConfigDB, cleanuptbl)
	cleanupstatetbl = map[string]interface{}{"PORT": map[string]interface{}{"PORT-1-1-L1IN": ""}}
	unloadDB(db.StateDB, cleanupstatetbl)
	cleanupcountertbl = map[string]interface{}{"PORT": map[string]interface{}{"PORT-1-1-L1IN_InputPower:15_pm_current": ""}}
	unloadDB(db.CountersDB, cleanupcountertbl)
}
