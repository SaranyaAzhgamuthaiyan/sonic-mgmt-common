package transformer_test

import (
	"testing"
	"time"

	"github.com/Azure/sonic-mgmt-common/translib/db"
)

/* APS config,counters, state DB fields and values */
func Test_aps_config_state_counter_DB_key_and_field_xfmr(t *testing.T) {
	var pre_req_map, cleanuptbl, cleanupstatetbl, cleanupcountertbl map[string]interface{}
	var expected_get_json, url []string

	/* APS Config DB test */
	t.Log("\n\n+++++++++++++ Performing Get on configDB APS Key-Xfmr and Field-Xfmr ++++++++++++")
	pre_req_map = map[string]interface{}{"APS": map[string]interface{}{"APS-1-1-2": map[string]interface{}{"revertive": true, "wait-to-restore-time": 0, "hold-off-time": 0, "primary-switch-threshold": 0.0, "primary-switch-hysteresis": 0.0, "secondary-switch-threshold": 0.0, "relative-switch-threshold": 0.0, "relative-switch-threshold-offset": 0.0, "force-to-port": "PRIMARY"}}}
	loadDB(db.ConfigDB, pre_req_map)
	expected_get_json = []string{"{\"openconfig-transport-line-protection:config\":{\"force-to-port\":\"PRIMARY\",\"hold-off-time\":0,\"name\":\"APS-1-1-2\",\"primary-switch-hysteresis\":\"0\",\"primary-switch-threshold\":\"0\",\"relative-switch-threshold\":\"0\",\"relative-switch-threshold-offset\":\"0\",\"revertive\":true,\"secondary-switch-threshold\":\"0\",\"wait-to-restore-time\":0}}"}
	url = []string{"/openconfig-transport-line-protection:aps/aps-modules/aps-module[name=APS-1-1-2]/config"}
	t.Run("Test get on configDB APS Key-Xfmr and Field-Xfmr.", processGetRequest(url[0], nil, expected_get_json[0], false))
	time.Sleep(1 * time.Second)
	t.Log("\n\n+++++++++++++ Done Performing Get on configDB APS Key-Xfmr and Field-Xfmr ++++++++++++")

	/* APS and APS_PORT State DB test */
	t.Log("\n\n+++++++++++++ Performing Get on  stateDB APS Key-Xfmr and Field-Xfmr ++++++++++++")
	pre_req_map = map[string]interface{}{"APS": map[string]interface{}{"APS-1-1-2": map[string]interface{}{"revertive": true, "wait-to-restore-time": 0, "hold-off-time": 0, "primary-switch-threshold": 0.0, "primary-switch-hysteresis": 0.0, "secondary-switch-threshold": 0.0, "relative-switch-threshold": 0.0, "relative-switch-threshold-offset": 0.0, "force-to-port": "NONE", "active-path": "PRIMARY"}}, "APS_PORT": map[string]interface{}{"APS-1-1-2_CommonIn": map[string]interface{}{"enabled": true, "target-attenuation": 0.0, "attenuation": 0.0}}}
	loadDB(db.StateDB, pre_req_map)
	expected_get_json = []string{"{\"openconfig-transport-line-protection:state\":{\"active-path\":\"openconfig-transport-line-protection:PRIMARY\",\"force-to-port\":\"NONE\",\"hold-off-time\":0,\"name\":\"APS-1-1-2\",\"primary-switch-hysteresis\":\"0\",\"primary-switch-threshold\":\"0\",\"relative-switch-threshold\":\"0\",\"relative-switch-threshold-offset\":\"0\",\"revertive\":true,\"secondary-switch-threshold\":\"0\",\"wait-to-restore-time\":0}}", "{\"openconfig-transport-line-protection:state\":{\"attenuation\":\"0\",\"enabled\":true,\"target-attenuation\":\"0\"}}"}
	url = []string{"/openconfig-transport-line-protection:aps/aps-modules/aps-module[name=APS-1-1-2]/state", "/openconfig-transport-line-protection:aps/aps-modules/aps-module[name=APS-1-1-2]/ports/common-in/state"}
	for idx, url := range url {
		t.Run("Test get on state DB Key-Xfmr and Field-Xfmr for "+url, processGetRequest(url, nil, expected_get_json[idx], false))
		time.Sleep(1 * time.Second)
	}
	t.Log("\n\n+++++++++++++ Done Performing Get on stateDB APS Key-Xfmr and Field-Xfmr ++++++++++++")

	/* APS_PORT Counters DB test */
	t.Log("\n\n+++++++++++++ Performing Get on  countersDB APS Key-Xfmr and Field-Xfmr ++++++++++++")
	pre_req_map = map[string]interface{}{"APS_PORT": map[string]interface{}{"APS-1-1-2_CommonIn_OpticalPower:15_pm_current": map[string]interface{}{"instant": "0.55", "min-time": "1680134400477338361", "min": "0.55", "max": "0.6", "interval": "86400000000000", "validity": "incomplete", "starttime": "1680134400000000000", "max-time": "1680136923116487196", "current_validity": "complete", "avg": "0.55"}}}
	loadDB(db.CountersDB, pre_req_map)
	expected_get_json = []string{"{\"openconfig-transport-line-protection:state\":{\"attenuation\":\"0\",\"enabled\":true,\"optical-power\":{\"avg\":\"0.55\",\"instant\":\"0.55\",\"interval\":\"86400000000000\",\"max\":\"0.6\",\"max-time\":\"1680136923116487196\",\"min\":\"0.55\",\"min-time\":\"1680134400477338361\"},\"target-attenuation\":\"0\"}}"}
	url = []string{"/openconfig-transport-line-protection:aps/aps-modules/aps-module[name=APS-1-1-2]/ports/common-in/state"}
	t.Run("Test get on countersDB APS Key-Xfmr and Field-Xfmr.", processGetRequest(url[0], nil, expected_get_json[0], false))
	time.Sleep(1 * time.Second)
	t.Log("\n\n+++++++++++++ Done Performing Get on countersDB APS Key-Xfmr and Field-Xfmr ++++++++++++")

	cleanuptbl = map[string]interface{}{"APS": map[string]interface{}{"APS-1-1-2": ""}}
	unloadDB(db.ConfigDB, cleanuptbl)
	cleanupstatetbl = map[string]interface{}{"APS": map[string]interface{}{"APS-1-1-2": ""}, "APS_PORT": map[string]interface{}{"APS-1-1-2_CommonIn": ""}}
	unloadDB(db.StateDB, cleanupstatetbl)
	cleanupcountertbl = map[string]interface{}{"APS_PORT": map[string]interface{}{"APS-1-1-2_CommonIn_OpticalPower:15_pm_current": ""}}
	unloadDB(db.CountersDB, cleanupcountertbl)
}
