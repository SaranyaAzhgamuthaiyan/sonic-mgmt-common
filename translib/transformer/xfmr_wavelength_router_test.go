package transformer_test

import (
	"testing"
	"time"

	"github.com/Azure/sonic-mgmt-common/translib/db"
)

/*  MEDIA_CHANNEL config, state DB fields and values */
func Test_media_channel_config_state_DB_key_and_field_xfmr(t *testing.T) {
	var pre_req_map, cleanuptbl, cleanupstatetbl map[string]interface{}
	var expected_get_json, url []string

	/* MEDIA_CHANNEL Config DB test */
	t.Log("\n\n+++++++++++++ Performing Get on configDB Media_channel Key-Xfmr and Field-Xfmr ++++++++++++")
	pre_req_map = map[string]interface{}{"MEDIA_CHANNEL": map[string]interface{}{"110": map[string]interface{}{"index": 110, "name": "Channel", "lower-frequency": 193775000, "upper-frequency": 193775000, "admin-status": "ENABLED", "super-channel": "true", "super-channel-parent": 1, "attenuation-control-mode": "ATTENUATION_FIXED_LOSS", "source-port-name": "PORT-1-1-WSS1", "dest-port-name": "PORT-1-1-WSS2"}, "110|193775000|193775000": map[string]interface{}{"lower-frequency": 193775000, "upper-frequency": 193775000, "attenuation-value": 10.5, "wait-to-restore-time": 300, "target-power": -5.0}}}
	loadDB(db.ConfigDB, pre_req_map)
	expected_get_json = []string{"{\"openconfig-wavelength-router:config\":{\"admin-status\":\"ENABLED\",\"attenuation-control-mode\":\"openconfig-wavelength-router:ATTENUATION_FIXED_LOSS\",\"index\":110,\"lower-frequency\":\"193775000\",\"name\":\"Channel\",\"super-channel\":true,\"super-channel-parent\":1,\"upper-frequency\":\"193775000\"}}", "{\"openconfig-wavelength-router:config\":{\"lower-frequency\":\"193775000\",\"target-power\":\"-5\",\"upper-frequency\":\"193775000\"}}"}
	url = []string{"/openconfig-wavelength-router:wavelength-router/media-channels/channel[index=110]/config", "/openconfig-wavelength-router:wavelength-router/media-channels/channel[index=110]/spectrum-power-profile/distribution[lower-frequency=193775000][upper-frequency=193775000]/config"}
	for idx, url := range url {
		t.Run("Test get on config DB Key-Xfmr and Field-Xfmr for "+url, processGetRequest(url, nil, expected_get_json[idx], false))
		time.Sleep(1 * time.Second)
	}
	t.Log("\n\n+++++++++++++ Done Performing Get on configDB Media_channel Key-Xfmr and Field-Xfmr++++++++++++")

	/* MEDIA_CHANNEL State DB test */
	t.Log("\n\n+++++++++++++ Performing Get on stateDB Media_channel Key-Xfmr and Field-Xfmr ++++++++++++")
	pre_req_map = map[string]interface{}{"MEDIA_CHANNEL": map[string]interface{}{"110": map[string]interface{}{"index": 110, "name": "WavelengthRouter1", "lower-frequency": 193775000, "upper-frequency": 193775000, "admin-status": "ENABLED", "super-channel": "false", "super-channel-parent": 0, "ase-control-mode": "ASE_ENABLED", "ase-injection-mode": "MODE_THRESHOLD", "ase-injection-threshold": 10.0, "ase-injection-delta": 2.0, "attenuation-control-mode": "ATTENUATION_FIXED_LOSS", "source-port-name": "PORT-1-1-WSS1", "dest-port-name": "PORT-1-1-WSS2", "oper-status": "UP", "ase-status": "PRESENT"}, "110|193775000|193775000": map[string]interface{}{"lower-frequency": 193775000, "upper-frequency": 193775000, "attenuation-value": 10.5, "wait-to-restore-time": 300, "target-power": -5.0}}}
	loadDB(db.StateDB, pre_req_map)
	expected_get_json = []string{"{\"openconfig-wavelength-router:state\":{\"admin-status\":\"ENABLED\",\"ase-control-mode\":\"openconfig-wavelength-router:ASE_ENABLED\",\"ase-injection-delta\":\"2\",\"ase-injection-mode\":\"MODE_THRESHOLD\",\"ase-injection-threshold\":\"10\",\"ase-status\":\"PRESENT\",\"attenuation-control-mode\":\"openconfig-wavelength-router:ATTENUATION_FIXED_LOSS\",\"index\":110,\"lower-frequency\":\"193775000\",\"name\":\"WavelengthRouter1\",\"oper-status\":\"UP\",\"super-channel\":false,\"super-channel-parent\":0,\"upper-frequency\":\"193775000\"}}", "{\"openconfig-wavelength-router:state\":{\"lower-frequency\":\"193775000\",\"target-power\":\"-5\",\"upper-frequency\":\"193775000\"}}"}
	url = []string{"/openconfig-wavelength-router:wavelength-router/media-channels/channel[index=110]/state", "/openconfig-wavelength-router:wavelength-router/media-channels/channel[index=110]/spectrum-power-profile/distribution[lower-frequency=193775000][upper-frequency=193775000]/state"}
	for idx, url := range url {
		t.Run("Test get on state DB Key-Xfmr and Field-Xfmr for "+url, processGetRequest(url, nil, expected_get_json[idx], false))
		time.Sleep(1 * time.Second)
	}
	t.Log("\n\n+++++++++++++ Done Get on stateDB Media_channel Key-Xfmr and Field-Xfmr ++++++++++++")

	cleanuptbl = map[string]interface{}{"MEDIA_CHANNEL": map[string]interface{}{"110": "", "110|193775000|193775000": ""}}
	unloadDB(db.ConfigDB, cleanuptbl)
	cleanupstatetbl = map[string]interface{}{"MEDIA_CHANNEL": map[string]interface{}{"110": "", "110|193775000|193775000": ""}}
	unloadDB(db.StateDB, cleanupstatetbl)
}
