//go:build xfmrtest
// +build xfmrtest

package transformer_test

import (
	"fmt"
	"github.com/Azure/sonic-mgmt-common/translib/db"
	"testing"
	"time"
)

/*  MEDIA_CHANNEL config, state DB fields and values */
func Test_media_channel_config_state_DB_key(t *testing.T) {
	var pre_req_map, cleanuptbl, cleanupstatetbl map[string]interface{}
	var expected_get_json, url []string

	for idx, dbName := range getMdbNames() {
		/* MEDIA_CHANNEL Config DB test */
		t.Log("\n\n+++++++++++++ Performing Get on configDB Media_channel Key ++++++++++++")
		pre_req_map = map[string]interface{}{
			"MEDIA_CHANNEL": map[string]interface{}{
				fmt.Sprintf("%d10", idx+1): map[string]interface{}{
					"index":                    fmt.Sprintf("%d10", idx+1),
					"name":                     "Channel",
					"lower-frequency":          193775000,
					"upper-frequency":          193775000,
					"admin-status":             "ENABLED",
					"super-channel":            "true",
					"super-channel-parent":     1,
					"attenuation-control-mode": "ATTENUATION_FIXED_LOSS",
					"source-port-name":         "PORT-1-1-WSS1",
					"dest-port-name":           "PORT-1-1-WSS2",
					"ase-control-mode":         "ASE_ENABLED",
					"ase-injection-threshold":  10.0,
				},
			},
			"MEDIA_CHANNEL_DISTRIBUTION": map[string]interface{}{
				fmt.Sprintf("%d10|193775000|193775000", idx+1): map[string]interface{}{
					"lower-frequency": 193775000,
					"upper-frequency": 193775000,
					"target-power":    -5.0,
				},
			},
		}

		loadDB(dbName, db.ConfigDB, pre_req_map)
		expected_get_json = []string{
			fmt.Sprintf("{\"openconfig-wavelength-router:config\":{\"admin-status\":\"ENABLED\",\"ase-control-mode\":\"openconfig-wavelength-router:ASE_ENABLED\",\"ase-injection-threshold\":\"10\",\"attenuation-control-mode\":\"openconfig-wavelength-router:ATTENUATION_FIXED_LOSS\",\"index\":%d10,\"lower-frequency\":\"193775000\",\"name\":\"Channel\",\"super-channel\":true,\"super-channel-parent\":1,\"upper-frequency\":\"193775000\"}}", idx+1),
			"{\"openconfig-wavelength-router:config\":{\"lower-frequency\":\"193775000\",\"target-power\":\"-5\",\"upper-frequency\":\"193775000\"}}",
		}

		url = []string{
			fmt.Sprintf("/openconfig-wavelength-router:wavelength-router/media-channels/channel[index=%d10]/config", idx+1),
			fmt.Sprintf("/openconfig-wavelength-router:wavelength-router/media-channels/channel[index=%d10]/spectrum-power-profile/distribution[lower-frequency=193775000][upper-frequency=193775000]/config", idx+1),
		}
		for idx, url := range url {
			t.Run("Test get on config DB Key for "+url, processGetRequest(url, nil, expected_get_json[idx], false))
			time.Sleep(1 * time.Second)
		}
		t.Log("\n\n+++++++++++++ Done Performing Get on configDB Media_channel Key ++++++++++++")

		/* MEDIA_CHANNEL State DB test */
		t.Log("\n\n+++++++++++++ Performing Get on stateDB Media_channel Key ++++++++++++")
		pre_req_map = map[string]interface{}{
			"MEDIA_CHANNEL_TABLE": map[string]interface{}{
				fmt.Sprintf("%d10", idx+1): map[string]interface{}{
					"index":                    fmt.Sprintf("%d10", idx+1),
					"name":                     "WavelengthRouter1",
					"lower-frequency":          193775000,
					"upper-frequency":          193775000,
					"admin-status":             "ENABLED",
					"super-channel":            "false",
					"super-channel-parent":     0,
					"ase-control-mode":         "ASE_ENABLED",
					"ase-injection-mode":       "MODE_THRESHOLD",
					"ase-injection-threshold":  10.0,
					"ase-injection-delta":      2.0,
					"attenuation-control-mode": "ATTENUATION_FIXED_LOSS",
					"source-port-name":         "PORT-1-1-WSS1",
					"dest-port-name":           "PORT-1-1-WSS2",
					"oper-status":              "UP",
					"ase-status":               "PRESENT",
				},
			},
			"MEDIA_CHANNEL_DISTRIBUTION_TABLE": map[string]interface{}{
				fmt.Sprintf("%d10|193775000|193775000", idx+1): map[string]interface{}{
					"lower-frequency": 193775000,
					"upper-frequency": 193775000,
					"target-power":    -5.0,
				},
			},
		}

		loadDB(dbName, db.StateDB, pre_req_map)
		expected_get_json = []string{
			fmt.Sprintf("{\"openconfig-wavelength-router:state\":{\"admin-status\":\"ENABLED\",\"ase-control-mode\":\"openconfig-wavelength-router:ASE_ENABLED\",\"ase-injection-delta\":\"2\",\"ase-injection-mode\":\"MODE_THRESHOLD\",\"ase-injection-threshold\":\"10\",\"ase-status\":\"PRESENT\",\"attenuation-control-mode\":\"openconfig-wavelength-router:ATTENUATION_FIXED_LOSS\",\"index\":%d10,\"lower-frequency\":\"193775000\",\"name\":\"WavelengthRouter1\",\"oper-status\":\"UP\",\"super-channel\":false,\"super-channel-parent\":0,\"upper-frequency\":\"193775000\"}}", idx+1),
			"{\"openconfig-wavelength-router:state\":{\"lower-frequency\":\"193775000\",\"target-power\":\"-5\",\"upper-frequency\":\"193775000\"}}",
		}
		url = []string{
			fmt.Sprintf("/openconfig-wavelength-router:wavelength-router/media-channels/channel[index=%d10]/state", idx+1),
			fmt.Sprintf("/openconfig-wavelength-router:wavelength-router/media-channels/channel[index=%d10]/spectrum-power-profile/distribution[lower-frequency=193775000][upper-frequency=193775000]/state", idx+1),
		}
		for idx, url := range url {
			t.Run("Test get on state DB Key for "+url, processGetRequest(url, nil, expected_get_json[idx], false))
			time.Sleep(1 * time.Second)
		}
		t.Log("\n\n+++++++++++++ Done Get on stateDB Media_channel Key ++++++++++++")

		//Unload the Data
		cleanuptbl = map[string]interface{}{
			"MEDIA_CHANNEL": map[string]interface{}{
				fmt.Sprintf("%d10", idx+1): "",
			},
			"MEDIA_CHANNEL_DISTRIBUTION": map[string]interface{}{
				fmt.Sprintf("%d10|193775000|193775000", idx+1): "",
			},
		}
		unloadDB(dbName, db.ConfigDB, cleanuptbl)

		cleanupstatetbl = map[string]interface{}{
			"MEDIA_CHANNEL_TABLE": map[string]interface{}{
				fmt.Sprintf("%d10", idx+1): "",
			},
			"MEDIA_CHANNEL_DISTRIBUTION_TABLE": map[string]interface{}{
				fmt.Sprintf("%d10|193775000|193775000", idx+1): "",
			},
		}
		unloadDB(dbName, db.StateDB, cleanupstatetbl)
	}
}

func Test_set_media_channel_config_DB_key_and_field_xfmr(t *testing.T) {
	t.Log("\n\n+++++++++++++ Performing Create/Replace/Delete on Media_channel ++++++++++++")

	for idx, dbName := range getMdbNames() {
		url := fmt.Sprintf("/openconfig-wavelength-router:wavelength-router/media-channels/channel[index=%d10]/config", idx+1)
		url_body_json := "{\"config\":{\"admin-status\":\"ENABLED\",\"ase-control-mode\":\"openconfig-wavelength-router:ASE_ENABLED\",\"ase-injection-threshold\":\"20\",\"attenuation-control-mode\":\"openconfig-wavelength-router:ATTENUATION_FIXED_LOSS\",\"lower-frequency\":\"193775000\",\"name\":\"Media_Channel\",\"super-channel\":false,\"super-channel-parent\":150,\"upper-frequency\":\"193775000\"}}"

		pre_req_map := map[string]interface{}{
			"MEDIA_CHANNEL": map[string]interface{}{
				fmt.Sprintf("%d10", idx+1): map[string]interface{}{
					"name":                     "Channel",
					"lower-frequency":          193775000,
					"upper-frequency":          193775000,
					"admin-status":             "ENABLED",
					"super-channel":            "true",
					"super-channel-parent":     1,
					"attenuation-control-mode": "ATTENUATION_FIXED_LOSS",
					"ase-control-mode":         "ASE_ENABLED",
					"ase-injection-threshold":  10.0,
				},
			},
		}
		loadDB(dbName, db.ConfigDB, pre_req_map)
		expected_map := map[string]interface{}{
			"MEDIA_CHANNEL": map[string]interface{}{
				fmt.Sprintf("%d10", idx+1): map[string]interface{}{
					"name":                     "Media_Channel",
					"lower-frequency":          193775000,
					"upper-frequency":          193775000,
					"admin-status":             "ENABLED",
					"super-channel":            "false",
					"super-channel-parent":     150,
					"attenuation-control-mode": "ATTENUATION_FIXED_LOSS",
					"ase-control-mode":         "ASE_ENABLED",
					"ase-injection-threshold":  20.0,
				},
			},
		}

		time.Sleep(1 * time.Second)

		//PUT Test
		t.Run("Replace Test on MEDIA_CHANNEL yang", processSetRequest(url, url_body_json, "PUT", false))
		time.Sleep(1 * time.Second)

		//GET Test
		t.Run("Verify replace on MEDIA_CHANNEL yang", verifyDbResult(rclient[dbName], fmt.Sprintf("MEDIA_CHANNEL|%d10", idx+1), expected_map, false))
		time.Sleep(1 * time.Second)

		//Delete Test
		url = fmt.Sprintf("/openconfig-wavelength-router:wavelength-router/media-channels/channel[index=%d10]/config", idx+1)
		t.Run("DELETE Test on MEDIA_CHANNEL yang", processDeleteRequest(url, false))
		time.Sleep(1 * time.Second)

		//POST Test
		url = "/openconfig-wavelength-router:wavelength-router/media-channels"
		post_payload := fmt.Sprintf("{\"channel\":[{\"index\":%d41,\"config\": {\"name\": \"NewChannel\",\"lower-frequency\": \"193775000\", \"upper-frequency\": \"193775000\", \"admin-status\": \"ENABLED\", \"super-channel\": false, \"super-channel-parent\": 130, \"attenuation-control-mode\": \"ATTENUATION_FIXED_LOSS\",\"ase-control-mode\":\"ASE_ENABLED\",\"ase-injection-threshold\":\"10.40\"}}]}", idx+1)
		t.Run("Create test on MEDIA_CHANNEL yang", processSetRequest(url, post_payload, "POST", false))
		time.Sleep(1 * time.Second)

		expected_map = map[string]interface{}{
			"MEDIA_CHANNEL": map[string]interface{}{
				fmt.Sprintf("%d41", idx+1): map[string]interface{}{
					"name":                     "NewChannel",
					"lower-frequency":          193775000,
					"upper-frequency":          193775000,
					"admin-status":             "ENABLED",
					"super-channel":            "false",
					"super-channel-parent":     130,
					"attenuation-control-mode": "ATTENUATION_FIXED_LOSS",
					"ase-control-mode":         "ASE_ENABLED",
					"ase-injection-threshold":  10.40,
				},
			},
		}

		//Verify After the POST test
		t.Run("Verify After the POST test on MEDIA_CHANNEL yang", verifyDbResult(rclient[dbName], fmt.Sprintf("MEDIA_CHANNEL|%d41", idx+1), expected_map, false))

		//Delete After the POST test
		url = fmt.Sprintf("/openconfig-wavelength-router:wavelength-router/media-channels/channel[index=%d41]/config", idx+1)
		t.Run("DELETE After the POST test on MEDIA_CHANNEL yang", processDeleteRequest(url, false))
		time.Sleep(1 * time.Second)

		url = "/openconfig-wavelength-router:wavelength-router/media-channels"

		//Bulk POST Test
		Bulk_post_payload := fmt.Sprintf("{\"channel\":[{\"index\":%d51,\"config\": {\"name\": \"NewChannel\",\"lower-frequency\": \"193775000\", \"upper-frequency\": \"193775000\", \"admin-status\": \"ENABLED\", \"super-channel\": true, \"super-channel-parent\": 150, \"attenuation-control-mode\": \"ATTENUATION_FIXED_LOSS\",\"ase-control-mode\":\"ASE_ENABLED\",\"ase-injection-threshold\":\"15.40\"}},{\"index\":%d61,\"config\": {\"name\": \"NewChannel\",\"lower-frequency\": \"193775000\", \"upper-frequency\": \"193775000\", \"admin-status\": \"ENABLED\", \"super-channel\": false, \"super-channel-parent\": 160, \"attenuation-control-mode\": \"ATTENUATION_FIXED_LOSS\",\"ase-control-mode\":\"ASE_ENABLED\",\"ase-injection-threshold\":\"16.40\"}}]}", idx+1, idx+1)
		t.Run("Create test on MEDIA_CHANNEL yang", processSetRequest(url, Bulk_post_payload, "POST", false))

		expected_map = map[string]interface{}{
			"MEDIA_CHANNEL": map[string]interface{}{
				fmt.Sprintf("%d51", idx+1): map[string]interface{}{
					"name":                     "NewChannel",
					"lower-frequency":          193775000,
					"upper-frequency":          193775000,
					"admin-status":             "ENABLED",
					"super-channel":            "true",
					"super-channel-parent":     150,
					"attenuation-control-mode": "ATTENUATION_FIXED_LOSS",
					"ase-control-mode":         "ASE_ENABLED",
					"ase-injection-threshold":  15.40,
				},
			},
		}

		time.Sleep(1 * time.Second)
		//Verify After the Bulk test
		t.Run("Verify After the Bulk test on MEDIA_CHANNEL yang", verifyDbResult(rclient[dbName], fmt.Sprintf("MEDIA_CHANNEL|%d51", idx+1), expected_map, false))

		expected_map = map[string]interface{}{
			"MEDIA_CHANNEL": map[string]interface{}{
				fmt.Sprintf("%d61", idx+1): map[string]interface{}{
					"name":                     "NewChannel",
					"lower-frequency":          193775000,
					"upper-frequency":          193775000,
					"admin-status":             "ENABLED",
					"super-channel":            "false",
					"super-channel-parent":     160,
					"attenuation-control-mode": "ATTENUATION_FIXED_LOSS",
					"ase-control-mode":         "ASE_ENABLED",
					"ase-injection-threshold":  16.40,
				},
			},
		}

		time.Sleep(1 * time.Second)
		//Verify After the POST test
		t.Run("Verify After the Bulk test on MEDIA_CHANNEL yang", verifyDbResult(rclient[dbName], fmt.Sprintf("MEDIA_CHANNEL|%d61", idx+1), expected_map, false))

		//Delete keys After the Bulk Test
		url = fmt.Sprintf("/openconfig-wavelength-router:wavelength-router/media-channels/channel[index=%d51]/config", idx+1)
		t.Run("DELETE After the Bulk Test on MEDIA_CHANNEL yang", processDeleteRequest(url, false))
		time.Sleep(1 * time.Second)

		url = fmt.Sprintf("/openconfig-wavelength-router:wavelength-router/media-channels/channel[index=%d61]/config", idx+1)
		t.Run("DELETE After the Bulk Test on MEDIA_CHANNEL yang", processDeleteRequest(url, false))
		time.Sleep(1 * time.Second)
	}
	t.Log("\n\n+++++++++++++ Done Performing Create/Replace/Delete on Media_channel ++++++++++++")
}

func Test_media_allow_write_multiple_namespace(t *testing.T) {
	t.Log("\n\n+++++++++++++ Performing Media_channel write for multiple namespace (should be allowed) ++++++++++++")

	url := "/openconfig-wavelength-router:wavelength-router/media-channels"
	url_body_json := `{
        "channel": [
            {
                "index": 121,
                "config": {
                    "name": "NewChannel",
                    "lower-frequency": "193775000",
                    "upper-frequency": "193775000",
                    "admin-status": "ENABLED",
                    "super-channel": true,
                    "super-channel-parent": 150,
                    "attenuation-control-mode": "ATTENUATION_FIXED_LOSS",
                    "ase-control-mode": "ASE_ENABLED",
                    "ase-injection-threshold": "15.40"
                }
            },
            {
                "index": 210,
                "config": {
                    "name": "NewChannel",
                    "lower-frequency": "193775000",
                    "upper-frequency": "193775000",
                    "admin-status": "ENABLED",
                    "super-channel": false,
                    "super-channel-parent": 160,
                    "attenuation-control-mode": "ATTENUATION_FIXED_LOSS",
                    "ase-control-mode": "ASE_ENABLED",
                    "ase-injection-threshold": "16.40"
                }
            }
        ]
    }`

	// false = no error expected
	t.Run("Test allow writing multiple namespace for MEDIA_CHANNEL", processSetRequest(url, url_body_json, "POST", false, nil))

	t.Log("\n\n+++++++++++++ Done Performing Media_channel write for multiple namespace ++++++++++++")
}
