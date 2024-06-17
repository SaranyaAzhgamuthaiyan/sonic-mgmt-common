package transformer_test

import (
	"testing"
	"time"

	"github.com/Azure/sonic-mgmt-common/translib/db"
)

/* Attenuator config, state, counters DB fields and values */
func Test_attenuator_config_state_counter_DB_key_and_field_xfmr(t *testing.T) {
	var pre_req_map, cleanuptbl, cleanupstatetbl, cleanupcountertbl map[string]interface{}
	var url string

	/* ATTENUATOR Config DB test */
	t.Log("\n\n+++++++++++++ Performing Get on configDB Attenuator Key-Xfmr and Field-Xfmr ++++++++++++")
	pre_req_map = map[string]interface{}{"ATTENUATOR": map[string]interface{}{"ATTENUATOR-1-1-3": map[string]interface{}{"attenuation": "0.0", "attenuation-mode": "CONSTANT_ATTENUATION", "target-output-power": "0.0", "enabled": "true"}}}
	loadDB(db.ConfigDB, pre_req_map)
	expected_get_json := "{\"openconfig-optical-attenuator:config\":{\"attenuation\":\"0\",\"attenuation-mode\":\"openconfig-optical-attenuator:CONSTANT_ATTENUATION\",\"enabled\":true,\"name\":\"ATTENUATOR-1-1-3\",\"target-output-power\":\"0\"}}"
	url = "/openconfig-optical-attenuator:optical-attenuator/attenuators/attenuator[name=ATTENUATOR-1-1-3]/config"
	t.Run("Test get on ConfigDB Attnuator Key-Xfmr and Field-Xfmr.", processGetRequest(url, nil, expected_get_json, false))
	time.Sleep(1 * time.Second)
	t.Log("\n\n+++++++++++++ Done Performing Get on configDB Attenuator Key-Xfmr and Field-Xfmr++++++++++++")

	/* ATTENUATOR State DB test */
	t.Log("\n\n+++++++++++++ Performing Get on stateDB Attenuator Key-Xfmr and Field-Xfmr ++++++++++++")
	pre_req_map = map[string]interface{}{"ATTENUATOR": map[string]interface{}{"ATTENUATOR-1-1-3": map[string]interface{}{"attenuation": "0.0", "name": "ATTENUATOR-1-1-3", "target-output-power": "0.0", "ingress-port": "PORT-1-2-2-VOAIN", "component": "ATTENUATOR-1-1-3", "enabled": "true", "egress-port": "PORT-1-2-2-VOAOUT", "fix-attenuation": "0.55", "attenuation-mode": "CONSTANT_ATTENUATION"}}}
	loadDB(db.StateDB, pre_req_map)
	expected_get_json = "{\"openconfig-optical-attenuator:state\":{\"attenuation\":\"0\",\"attenuation-mode\":\"openconfig-optical-attenuator:CONSTANT_ATTENUATION\",\"component\":\"ATTENUATOR-1-1-3\",\"egress-port\":\"PORT-1-2-2-VOAOUT\",\"enabled\":true,\"ingress-port\":\"PORT-1-2-2-VOAIN\",\"name\":\"ATTENUATOR-1-1-3\",\"target-output-power\":\"0\"}}"
	url = "/openconfig-optical-attenuator:optical-attenuator/attenuators/attenuator[name=ATTENUATOR-1-1-3]/state"
	t.Run("Test get on StateDB Attnuator Key-Xfmr and Field-Xfmr.", processGetRequest(url, nil, expected_get_json, false))
	time.Sleep(1 * time.Second)
	t.Log("\n\n+++++++++++++ Done Performing Get on stateDB Attnuator Key-Xfmr and Field-Xfmr++++++++++++")

	/* ATTENUATOR Counter DB test */
	t.Log("\n\n+++++++++++++ Performing Get on CounterDB Attenuator Key-Xfmr and Field-Xfmr ++++++++++++")
	pre_req_map = map[string]interface{}{"ATTENUATOR": map[string]interface{}{"ATTENUATOR-1-1-3_OutputPowerTotal:15_pm_current": map[string]interface{}{"instant": "0.0", "min-time": "1680168600865596372", "min": "0.0", "max": "0.0", "interval": "900000000000", "validity": "incomplete", "starttime": "1680168600000000000", "max-time": "1680168600865596372", "current_validity": "complete", "avg": "0.0"}}}
	loadDB(db.CountersDB, pre_req_map)
	expected_get_json = "{\"openconfig-optical-attenuator:output-power-total\":{\"avg\":\"0\",\"instant\":\"0\",\"interval\":\"900000000000\",\"max\":\"0\",\"max-time\":\"1680168600865596372\",\"min\":\"0\",\"min-time\":\"1680168600865596372\"}}"
	url = "/openconfig-optical-attenuator:optical-attenuator/attenuators/attenuator[name=ATTENUATOR-1-1-3]/state/output-power-total"
	t.Run("Test get on StateDB Attnuator Key-Xfmr and Field-Xfmr.", processGetRequest(url, nil, expected_get_json, false))
	time.Sleep(1 * time.Second)
	t.Log("\n\n+++++++++++++ Done Performing Get on CounterDB Attnuator Key-Xfmr and Field-Xfmr++++++++++++")

	cleanuptbl = map[string]interface{}{"ATTENUATOR": map[string]interface{}{"ATTENUATOR-1-1-3": ""}}
	unloadDB(db.ConfigDB, cleanuptbl)
	cleanupstatetbl = map[string]interface{}{"ATTENUATOR": map[string]interface{}{"ATTENUATOR-1-1-3": ""}}
	unloadDB(db.StateDB, cleanupstatetbl)
	cleanupcountertbl = map[string]interface{}{"ATTENUATOR": map[string]interface{}{"ATTENUATOR-1-1-3_OutputPowerTotal:15_pm_current": ""}}
	unloadDB(db.CountersDB, cleanupcountertbl)
}
