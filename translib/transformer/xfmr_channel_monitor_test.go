package transformer_test

import (
	"testing"
	"time"

	"github.com/Azure/sonic-mgmt-common/translib/db"
)

/* channel_monitor config,counters, state DB fields and values */
func Test_channel_monitor_config_state_DB_key_and_field_xfmr(t *testing.T) {
	var pre_req_map, cleanuptbl, cleanupstatetbl map[string]interface{}
	var expected_get_json, url []string

	/* channel_monitor config DB test */
	t.Log("\n\n+++++++++++++ Performing Get on configDB channel monitor  Key-Xfmr and Field-Xfmr ++++++++++++")
	pre_req_map = map[string]interface{}{"OCM": map[string]interface{}{"OCM-1-1-2": map[string]interface{}{"monitor-port": "test_monitor_port"}}}
	loadDB(db.ConfigDB, pre_req_map)
	//expected_get_json = []string{}
	//url = []string{"/openconfig-channel-monitor:channel-monitors/channel-monitor[name=OCM-1-1-2]/config"}
	//t.Run("Test get on config DB channel monitor Key-Xfmr and Field-Xfmr.", processGetRequest(url[0], nil, expected_get_json[0], false))
	//time.Sleep(1 * time.Second)
	t.Log("\n\n+++++++++++++ Done Performing Get on configDB channel monitor Key-Xfmr and Field-Xfmr ++++++++++++")

	/* channel monitor State DB test */
	t.Log("\n\n+++++++++++++ Performing Get on stateDB channel monitor Key-Xfmr and Field-Xfmr ++++++++++++")
	pre_req_map = map[string]interface{}{"OCM_TABLE": map[string]interface{}{"OCM-1-1-2": map[string]interface{}{"monitor-port": "test_monitor_port"}}, "OCM_CHANNEL": map[string]interface{}{"OCM-1-1-2|25|100": map[string]interface{}{"power": "5.0"}}}
	loadDB(db.StateDB, pre_req_map)
	expected_get_json = []string{"{\"openconfig-channel-monitor:state\":{\"monitor-port\":\"test_monitor_port\"}}", "{\"openconfig-channel-monitor:state\":{\"lower-frequency\":\"0\",\"power\":\"5\",\"upper-frequency\":\"25\"}}"}
	url = []string{"/openconfig-channel-monitor:channel-monitors/channel-monitor[name=OCM-1-1-2]/state", "/openconfig-channel-monitor:channel-monitors/channel-monitor[name=OCM-1-1-2]/channels/channel[lower-frequency=25][upper-frequency=100]/state"}
	for idx, url := range url {
		t.Run("Test get on state DB channel monitor Key-Xfmr and Field-Xfmr for "+url, processGetRequest(url, nil, expected_get_json[idx], false))
		time.Sleep(1 * time.Second)
	}
	t.Log("\n\n+++++++++++++ Done Performing Get on stateDB channel monitor Key-Xfmr and Field-Xfmr ++++++++++++")
	cleanuptbl = map[string]interface{}{"OCM": map[string]interface{}{"OCM-1-1-2": ""}}
	unloadDB(db.ConfigDB, cleanuptbl)
	cleanupstatetbl = map[string]interface{}{"OCM_TABLE": map[string]interface{}{"OCM-1-1-2": ""}, "OCM_CHANNEL": map[string]interface{}{"OCM-1-1-2|25|100": ""}}
	unloadDB(db.StateDB, cleanupstatetbl)
}
