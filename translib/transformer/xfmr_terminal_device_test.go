package transformer_test

import (
	"testing"
	"time"

	"github.com/Azure/sonic-mgmt-common/translib/db"
)

/* terminal device counters, state DB fields and values */
func Test_terminal_device_config_state_counter_DB_key_and_field_xfmr(t *testing.T) {
	var pre_req_map, cleanuptbl, cleanupstatetbl, cleanupcountertbl map[string]interface{}
	//var pre_req_map map[string]interface{}
	var expected_get_json, url []string

	/* terminal device config DB test */
	t.Log("\n\n+++++++++++++ Performing Get on configDB terminal device Key-Xfmr and Field-Xfmr ++++++++++++")
	pre_req_map = map[string]interface{}{"OC_COMPONENT_TABLE": map[string]interface{}{"OCH-1-1-L2": map[string]interface{}{"name": "OCH-1-1-L2"}}, "OCH": map[string]interface{}{"OCH-1-1-L2": map[string]interface{}{"operational-mode": "test_data", "frequency": "2", "target-output-power": "0.0"}}, "LOGICAL_CHANNEL": map[string]interface{}{"CH102": map[string]interface{}{"test-signal": "false", "loopback-mode": "NONE", "admin-state": "MAINT"}}, "ETHERNET": map[string]interface{}{"CH102": map[string]interface{}{"index": "102", "client-als": "LASER_SHUTDOWN", "als-delay": "1000"}}, "LLDP": map[string]interface{}{"CH102": map[string]interface{}{"index": "102", "enabled": "true", "snooping": "true"}}, "MODE": map[string]interface{}{"102": map[string]interface{}{"mode-id": "102", "description": "600G::DP-64QAM::FEC-15::63.054::BICHM[0]", "vendor-id": "Accelink"}}}
	loadDB(db.ConfigDB, pre_req_map)
	time.Sleep(1 * time.Second)
	expected_get_json = []string{"{\"openconfig-terminal-device:config\":{\"admin-state\":\"MAINT\",\"index\":102,\"loopback-mode\":\"NONE\",\"test-signal\":false}}", "{\"openconfig-terminal-device:config\":{\"enabled\":true,\"snooping\":true}}", "{\"openconfig-terminal-device:config\":{\"als-delay\":1000,\"client-als\":\"LASER_SHUTDOWN\"}}", "{\"openconfig-terminal-device:state\":{\"description\":\"600G::DP-64QAM::FEC-15::63.054::BICHM[0]\",\"vendor-id\":\"Accelink\"}}"}
	url = []string{"/openconfig-terminal-device:terminal-device/logical-channels/channel[index=102]/config", "/openconfig-terminal-device:terminal-device/logical-channels/channel[index=102]/ethernet/lldp/config", "/openconfig-terminal-device:terminal-device/logical-channels/channel[index=102]/ethernet/config", "/openconfig-terminal-device:terminal-device/operational-modes/mode[mode-id=102]/state"}
	for idx, url := range url {
		t.Run("Test get on config DB Key-Xfmr and Field-Xfmr for "+url, processGetRequest(url, nil, expected_get_json[idx], false))
		time.Sleep(1 * time.Second)
	}
	t.Log("\n\n+++++++++++++ Done Performing Get on configDB terminal device Key-Xfmr and Field-Xfmr ++++++++++++")

	/* terminal device state DB test */
	t.Log("\n\n+++++++++++++ Performing Get on stateDB terminal device Key-Xfmr and Field-Xfmr ++++++++++++")
	pre_req_map = map[string]interface{}{"LOGICAL_CHANNEL": map[string]interface{}{"CH102": map[string]interface{}{"index": "102", "rate-class": "TRIB_RATE_200G", "description": "GE-1-1-C3", "loopback-mode": "NONE", "test-signal": "false", "admin-state": "MAINT", "trib-protocol": "PROT_1GE", "link-state": "DOWN"}}, "ETHERNET": map[string]interface{}{"CH102": map[string]interface{}{"index": "102", "client-als": "LASER_SHUTDOWN", "als-delay": "1000"}}, "LLDP": map[string]interface{}{"CH102": map[string]interface{}{"index": "102", "enabled": "true", "snooping": "true"}}, "NEIGHBOR": map[string]interface{}{"CH102|1": map[string]interface{}{"id": "1", "system-name": "test_sysname", "system-description": "test_desc", "chassis-id": "test_chassis_id", "chassis-id-type": "PORT_COMPONENT", "age": "12", "last-update": "12", "ttl": "100", "port-id": "test_portid", "port-id-type": "PORT_COMPONENT", "port-description": "test_pdes", "management-address": "test_mg_add", "management-address-type": "test_mg_type"}}, "ASSIGNMENT": map[string]interface{}{"CH102|ASS1": map[string]interface{}{"index": "102", "assignment-type": "LOGICAL_CHANNEL", "description": "ASSIGNMENT-1-1-18", "tributary-slot-index": "1", "mapping": "GMP", "allocation": "400.0", "logical-channel": "102", "optical-channel": "102"}}}
	loadDB(db.StateDB, pre_req_map)
	time.Sleep(1 * time.Second)
	expected_get_json = []string{"{\"openconfig-terminal-device:state\":{\"admin-state\":\"MAINT\",\"description\":\"GE-1-1-C3\",\"index\":102,\"link-state\":\"DOWN\",\"loopback-mode\":\"NONE\",\"rate-class\":\"openconfig-transport-types:TRIB_RATE_200G\",\"test-signal\":false,\"trib-protocol\":\"openconfig-transport-types:PROT_1GE\"}}", "{\"openconfig-terminal-device:client-als\":\"LASER_SHUTDOWN\"}", "{\"openconfig-terminal-device:state\":{\"description\":\"600G::DP-64QAM::FEC-15::63.054::BICHM[0]\",\"vendor-id\":\"Accelink\"}}", "{\"openconfig-terminal-device:state\":{\"allocation\":\"400\",\"assignment-type\":\"LOGICAL_CHANNEL\",\"description\":\"ASSIGNMENT-1-1-18\",\"index\":1,\"logical-channel\":0,\"tributary-slot-index\":1}}", "{\"openconfig-terminal-device:state\":{\"age\":\"12\",\"chassis-id\":\"test_chassis_id\",\"chassis-id-type\":\"PORT_COMPONENT\",\"id\":\"1\",\"last-update\":\"12\",\"management-address\":\"test_mg_add\",\"management-address-type\":\"test_mg_type\",\"port-description\":\"test_pdes\",\"port-id\":\"test_portid\",\"port-id-type\":\"PORT_COMPONENT\",\"system-description\":\"test_desc\",\"system-name\":\"test_sysname\",\"ttl\":100}}"}
	url = []string{"/openconfig-terminal-device:terminal-device/logical-channels/channel[index=102]/state", "/openconfig-terminal-device:terminal-device/logical-channels/channel[index=102]/ethernet/state/client-als", "/openconfig-terminal-device:terminal-device/operational-modes/mode[mode-id=102]/state", "/openconfig-terminal-device:terminal-device/logical-channels/channel[index=102]/logical-channel-assignments/assignment[index=1]/state", "/openconfig-terminal-device:terminal-device/logical-channels/channel[index=102]/ethernet/lldp/neighbors/neighbor[id=1]/state"}
	for idx, url := range url {
		t.Run("Test get on state DB Key-Xfmr and Field-Xfmr for "+url, processGetRequest(url, nil, expected_get_json[idx], false))
		time.Sleep(1 * time.Second)
	}
	t.Log("\n\n+++++++++++++ Done Performing Get on configDB terminal device Key-Xfmr and Field-Xfmr ++++++++++++")

	/* terminal monitor Counters DB test */
	t.Log("\n\n+++++++++++++ Performing Get on counterDB terminal device Key-Xfmr and Field-Xfmr ++++++++++++")
	pre_req_map = map[string]interface{}{"ETHERNET": map[string]interface{}{"CH102:15_pm_current": map[string]interface{}{"in-pcs-errored-seconds": "0", "out-pcs-bip-errors": "0", "in-block-errors": "0", "out-block-errors": "0", "in-pcs-severely-errored-seconds": "0", "in-pcs-bip-errors": "0", "in-jabber-frames": "0", "in-fragment-frames": "0", "in-oversize-frames": "0", "out-crc-errors": "0", "in-pcs-unavailable-seconds": "0", "in-undersize-frames": "0", "in-crc-errors": "0"}}}
	loadDB(db.CountersDB, pre_req_map)
	expected_get_json = []string{"{\"openconfig-terminal-device:state\":{\"in-block-errors\":\"0\",\"in-fragment-frames\":\"0\",\"in-jabber-frames\":\"0\",\"in-pcs-bip-errors\":\"0\",\"in-pcs-errored-seconds\":\"0\",\"in-pcs-severely-errored-seconds\":\"0\",\"out-block-errors\":\"0\",\"out-pcs-bip-errors\":\"0\"}}"}
	url = []string{"/openconfig-terminal-device:terminal-device/logical-channels/channel[index=102]/ethernet/state"}
	t.Run("Test get on counter DB Ethernet Key-Xfmr and Field-Xfmr.", processGetRequest(url[0], nil, expected_get_json[0], false))
	time.Sleep(1 * time.Second)
	t.Log("\n\n+++++++++++++ Done Performing Get on counterDB <tableName> Key-Xfmr and Field-Xfmr ++++++++++++")

	cleanuptbl = map[string]interface{}{"OCH": map[string]interface{}{"OCH-1-1-L2": ""}, "OC_COMPONENT_TABLE": map[string]interface{}{"OCH-1-1-L2": ""}, "LOGICAL_CHANNEL": map[string]interface{}{"CH102": ""}, "ETHERNET": map[string]interface{}{"CH102": ""}, "LLDP": map[string]interface{}{"CH102": ""}, "MODE": map[string]interface{}{"102": ""}}
	unloadDB(db.ConfigDB, cleanuptbl)
	cleanupstatetbl = map[string]interface{}{"LOGICAL_CHANNEL": map[string]interface{}{"CH102": ""}, "ETHERNET": map[string]interface{}{"CH102": ""}, "LLDP": map[string]interface{}{"CH102": ""}, "NEIGHBOR": map[string]interface{}{"CH102|1": ""}, "ASSIGNMENT": map[string]interface{}{"CH102|ASS1": ""}}
	unloadDB(db.StateDB, cleanupstatetbl)
	cleanupcountertbl = map[string]interface{}{"ETHERNET": map[string]interface{}{"CH101:15_pm_current": ""}}
	unloadDB(db.CountersDB, cleanupcountertbl)
}
