package transformer_test

import (
	"testing"
	"time"

	"github.com/Azure/sonic-mgmt-common/translib/db"
)

/* AMPLIFIER config,counters, state DB fields and values */
func Test_amplifier_config_state_counter_DB_key_and_field_xfmr(t *testing.T) {
	var pre_req_map, cleanuptbl, cleanupstatetbl, cleanupcountertbl map[string]interface{}
	var expected_get_json, url []string

	/* AMPLIFIER Config DB test */
	t.Log("\n\n+++++++++++++ Performing Get on configDB AMPLIFIER Key-Xfmr and Field-Xfmr ++++++++++++")
	pre_req_map = map[string]interface{}{"AMPLIFIER": map[string]interface{}{"AMPLIFIER-1-1-2": map[string]interface{}{"type": "EDFA", "target-gain": 15.5, "max-gain": 20.0, "min-gain": 10.0, "target-gain-tilt": 2.5, "gain-range": "MID_GAIN_RANGE", "amp-mode": "CONSTANT_GAIN", "target-output-power": 0.0, "max-output-power": 5.0, "enabled": true, "fiber-type-profile": "SSMF"}}}
	loadDB(db.ConfigDB, pre_req_map)
	expected_get_json = []string{"{\"openconfig-optical-amplifier:config\":{\"amp-mode\":\"openconfig-optical-amplifier:CONSTANT_GAIN\",\"enabled\":true,\"fiber-type-profile\":\"openconfig-optical-amplifier:SSMF\",\"gain-range\":\"openconfig-optical-amplifier:MID_GAIN_RANGE\",\"max-gain\":\"20\",\"max-output-power\":\"5\",\"min-gain\":\"10\",\"name\":\"AMPLIFIER-1-1-2\",\"target-gain\":\"15.5\",\"target-gain-tilt\":\"2.5\",\"target-output-power\":\"0\",\"type\":\"openconfig-optical-amplifier:EDFA\"}}"}

	url = []string{"/openconfig-optical-amplifier:optical-amplifier/amplifiers/amplifier[name=AMPLIFIER-1-1-2]/config"}
	t.Run("Test get on configDB AMPLIFIER Key-Xfmr and Field-Xfmr.", processGetRequest(url[0], nil, expected_get_json[0], false))
	time.Sleep(1 * time.Second)
	t.Log("\n\n+++++++++++++ Done Performing Get on configDB AMPLIFIER Key-Xfmr and Field-Xfmr ++++++++++++")

	/* AMPLIFIER and OSC State DB test */
	t.Log("\n\n+++++++++++++ Performing Get on stateDB AMPLIFIER , OSC Key-Xfmr and Field-Xfmr ++++++++++++")
	pre_req_map = map[string]interface{}{"AMPLIFIER": map[string]interface{}{"AMPLIFIER-1-1-2": map[string]interface{}{"component": "AMPLIFIER-1-1-2", "location": "1-4", "name": "AMPLIFIER-1-1-2", "amp-mode": "CONSTANT_POWER", "egress-port": "test data", "enabled": "true", "fiber-type-profile": "DSF", "gain-range": "LOW_GAIN_RANGE", "ingress-port": "test data", "max-gain": "0.0", "max-output-power": "0.0", "min-gain": "0.0", "target-gain": "0.0", "target-gain-tilt": "0.0", "target-output-power": "0.0", "type": "EDFA"}}, "OSC": map[string]interface{}{"OSC-1-1-2": map[string]interface{}{"output-frequency": "198538051"}}}
	loadDB(db.StateDB, pre_req_map)
	expected_get_json = []string{"{\"openconfig-optical-amplifier:state\":{\"amp-mode\":\"openconfig-optical-amplifier:CONSTANT_POWER\",\"component\":\"AMPLIFIER-1-1-2\",\"egress-port\":\"test data\",\"enabled\":true,\"fiber-type-profile\":\"openconfig-optical-amplifier:DSF\",\"gain-range\":\"openconfig-optical-amplifier:LOW_GAIN_RANGE\",\"ingress-port\":\"test data\",\"max-gain\":\"0\",\"min-gain\":\"0\",\"name\":\"AMPLIFIER-1-1-2\",\"target-gain\":\"0\",\"target-gain-tilt\":\"0\",\"target-output-power\":\"0\",\"type\":\"openconfig-optical-amplifier:EDFA\"}}", "{\"openconfig-optical-amplifier:state\":{\"interface\":\"OSC-1-1-2\",\"output-frequency\":\"198538051\"}}"}
	url = []string{"/openconfig-optical-amplifier:optical-amplifier/amplifiers/amplifier[name=AMPLIFIER-1-1-2]/state", "/openconfig-optical-amplifier:optical-amplifier/supervisory-channels/supervisory-channel[interface=OSC-1-1-2]/state"}
	for idx, url := range url {
		t.Run("Test get on state DB Key-Xfmr and Field-Xfmr for "+url, processGetRequest(url, nil, expected_get_json[idx], false))
		time.Sleep(1 * time.Second)
	}
	t.Log("\n\n+++++++++++++ Done Performing Get on stateDB AMPLIFIER, OSC Key-Xfmr and Field-Xfmr ++++++++++++")

	/* AMPLIFIER and OSC counters DB test */
	t.Log("\n\n+++++++++++++ Performing Get on counters AMPLIFIER , OSC Key-Xfmr and Field-Xfmr ++++++++++++")
	pre_req_map = map[string]interface{}{"AMPLIFIER": map[string]interface{}{"AMPLIFIER-1-1-2_ActualGain:15_pm_current": map[string]interface{}{"interval": "86400000000000", "starttime": "1709251200000000000", "max": "1.0", "max-time": "1909276854829015912", "min": "1.0", "min-time": "1909276854829015912", "instant": "1.0", "avg": "1.0", "current_validity": "complete", "validity": "incomplete"}}, "OSC": map[string]interface{}{"OSC-1-1-2_OutputPower:15_pm_current": map[string]interface{}{"interval": "900000000000", "starttime": "1709251200000000000", "max": "2.5", "max-time": "1680168786273045375", "min": "2.12", "min-time": "1680168676158235913", "instant": "2.29", "avg": "2.32", "current_validity": "complete", "validity": "incomplete"}}}
	loadDB(db.CountersDB, pre_req_map)
	expected_get_json = []string{"{\"openconfig-optical-amplifier:actual-gain\":{\"avg\":\"1\",\"instant\":\"1\",\"interval\":\"86400000000000\",\"max\":\"1\",\"max-time\":\"1909276854829015912\",\"min\":\"1\",\"min-time\":\"1909276854829015912\"}}", "{\"openconfig-optical-amplifier:output-power\":{\"avg\":\"2.32\",\"instant\":\"2.29\",\"interval\":\"900000000000\",\"max\":\"2.5\",\"max-time\":\"1680168786273045375\",\"min\":\"2.12\",\"min-time\":\"1680168676158235913\"}}"}
	url = []string{"/openconfig-optical-amplifier:optical-amplifier/amplifiers/amplifier[name=AMPLIFIER-1-1-2]/state/actual-gain", "/openconfig-optical-amplifier:optical-amplifier/supervisory-channels/supervisory-channel[interface=OSC-1-1-2]/state/output-power"}
	for idx, url := range url {
		t.Run("Test get on counters DB Key-Xfmr and Field-Xfmr for "+url, processGetRequest(url, nil, expected_get_json[idx], false))
		time.Sleep(1 * time.Second)
	}
	t.Log("\n\n+++++++++++++ Performing Get on counters AMPLIFIER , OSC Key-Xfmr and Field-Xfmr ++++++++++++")

	cleanuptbl = map[string]interface{}{"AMPLIFIER": map[string]interface{}{"AMPLIFIER-1-1-2": ""}}
	unloadDB(db.ConfigDB, cleanuptbl)
	cleanupstatetbl = map[string]interface{}{"AMPLIFIER": map[string]interface{}{"AMPLIFIER-1-1-2": ""}, "OSC": map[string]interface{}{"OSC-1-1-2": ""}}
	unloadDB(db.StateDB, cleanupstatetbl)
	cleanupcountertbl = map[string]interface{}{"AMPLIFIER": map[string]interface{}{"AMPLIFIER-1-1-2_ActualGain:15_pm_current": ""}, "OSC": map[string]interface{}{"OSC-1-1-2_OutputPower:15_pm_current": ""}}
	unloadDB(db.CountersDB, cleanupcountertbl)
}
