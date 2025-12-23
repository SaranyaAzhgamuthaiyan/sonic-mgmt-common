//go:build xfmrtest
// +build xfmrtest

package transformer_test

import (
	"fmt"
	"github.com/Azure/sonic-mgmt-common/translib/db"
	"testing"
	"time"
)

/* TRANSCEIVER config, state DB fields and values */
func Test_pfm_transceiver_config_state_DB_key(t *testing.T) {
	var pre_req_map, cleanuptbl, cleanupstatetbl map[string]interface{}
	var expected_get_json, url []string

	for idx, dbName := range getMdbNames() {
		/* TRANSCEIVER config DB test */
		t.Log("\n\n+++++++++++++ Performing Get on configDB TRANSCEIVER Key ++++++++++++")
		pre_req_map = map[string]interface{}{
			"OC_COMPONENT": map[string]interface{}{
				fmt.Sprintf("TRANSCEIVER-1-%d-C2", idx+1): map[string]interface{}{
					"name": fmt.Sprintf("TRANSCEIVER-1-%d-C2", idx+1),
				},
			},
			"TRANSCEIVER": map[string]interface{}{
				fmt.Sprintf("TRANSCEIVER-1-%d-C2", idx+1): map[string]interface{}{
					"enabled":              "true",
					"fec-mode":             "FEC_AUTO",
					"ethernet-pmd-preconf": "ETH_100GBASE_SR10",
				},
			},
		}
		loadDB(dbName, db.ConfigDB, pre_req_map)
		expected_get_json = []string{"{\"openconfig-platform-transceiver:config\":{\"enabled\":true,\"ethernet-pmd-preconf\":\"openconfig-transport-types:ETH_100GBASE_SR10\",\"fec-mode\":\"openconfig-platform-types:FEC_AUTO\"}}"}

		url = []string{fmt.Sprintf("/openconfig-platform:components/component[name=TRANSCEIVER-1-%d-C2]/openconfig-platform-transceiver:transceiver/config", idx+1)}
		t.Run("Test get on configDB TRANSCEIVER Key", processGetRequest(url[0], nil, expected_get_json[0], false))
		time.Sleep(1 * time.Second)
		t.Log("\n\n+++++++++++++ Done Performing Get on configDB TRANSCEIVER Key ++++++++++++")

		/* TRANSCEIVER State DB test */
		t.Log("\n\n+++++++++++++ Performing Get on stateDB TRANSCEIVER Key ++++++++++++")
		pre_req_map = map[string]interface{}{
			"TRANSCEIVER_TABLE": map[string]interface{}{
				fmt.Sprintf("TRANSCEIVER-1-%d-C2", idx+1): map[string]interface{}{
					"vendor":          "VendorX",
					"vendor-rev":      "1.0",
					"date-code":       "2024-10-03",
					"vendor-part":     "VX-1234",
					"fec-mode":        "FEC_ENABLED",
					"form-factor":     "QSFP28",
					"present":         "PRESENT",
					"enabled":         "true",
					"fault-condition": "false",
					"connector-type":  "LC_CONNECTOR",
					"ethernet-pmd":    "ETH_100GBASE_CWDM4",
				},
			},
			"PHYSICAL_CHANNEL_TABLE": map[string]interface{}{
				fmt.Sprintf("TRANSCEIVER-1-%d-C2|CH-1", idx+1): map[string]interface{}{
					"index":            1,
					"output-frequency": "193000000",
					"tx-laser":         "true",
				},
			},
		}
		loadDB(dbName, db.StateDB, pre_req_map)
		expected_get_json = []string{
			"{\"openconfig-platform-transceiver:state\":{\"connector-type\":\"openconfig-transport-types:LC_CONNECTOR\",\"date-code\":\"2024-10-03\",\"enabled\":true,\"ethernet-pmd\":\"openconfig-transport-types:ETH_100GBASE_CWDM4\",\"fault-condition\":false,\"fec-mode\":\"openconfig-platform-types:FEC_ENABLED\",\"form-factor\":\"openconfig-transport-types:QSFP28\",\"present\":\"PRESENT\",\"vendor\":\"VendorX\",\"vendor-part\":\"VX-1234\",\"vendor-rev\":\"1.0\"}}",
			"{\"openconfig-platform-transceiver:state\":{\"index\":1,\"output-frequency\":\"193000000\",\"tx-laser\":true}}",
		}
		url = []string{
			fmt.Sprintf("/openconfig-platform:components/component[name=TRANSCEIVER-1-%d-C2]/openconfig-platform-transceiver:transceiver/state", idx+1),
			fmt.Sprintf("/openconfig-platform:components/component[name=TRANSCEIVER-1-%d-C2]/openconfig-platform-transceiver:transceiver/physical-channels/channel[index=1]/state", idx+1),
		}
		for idx, url := range url {
			t.Run("Test get on state DB TRANSCEIVER Key for "+url, processGetRequest(url, nil, expected_get_json[idx], false))
			time.Sleep(1 * time.Second)
		}

		t.Log("\n\n+++++++++++++ Done Performing Get on stateDB TRANSCEIVER Key ++++++++++++")

		//Unload the Data
		cleanuptbl = map[string]interface{}{
			"OC_COMPONENT": map[string]interface{}{
				fmt.Sprintf("TRANSCEIVER-1-%d-C2", idx+1): "",
			},
			"TRANSCEIVER": map[string]interface{}{
				fmt.Sprintf("TRANSCEIVER-1-%d-C2", idx+1): "",
			},
		}
		unloadDB(dbName, db.ConfigDB, cleanuptbl)

		cleanupstatetbl = map[string]interface{}{
			"TRANSCEIVER_TABLE": map[string]interface{}{
				fmt.Sprintf("TRANSCEIVER-1-%d-C2", idx+1): "",
			},
			"PHYSICAL_CHANNEL_TABLE": map[string]interface{}{
				fmt.Sprintf("TRANSCEIVER-1-%d-C2|CH-1", idx+1): "",
			},
		}
		unloadDB(dbName, db.StateDB, cleanupstatetbl)
	}
}

func pfmMapDetails() map[string]interface{} {
	return map[string]interface{}{
		"serial-no":        "test_serial_no",
		"part-no":          "10",
		"mfg-name":         "test_mfg_name",
		"mfg-date":         "09-10-24",
		"hardware-version": "test_hw_version",
		"oper-status":      "ACTIVE",
		"empty":            "false",
		"removable":        "true",
		"parent":           "ACTIVE",
		"software-version": "test_software_version",
		"memory-available": "100",
		"memory-utilized":  "200",
	}
}

/* FAN state DB fields and values */
func Test_fan_state_DB_key(t *testing.T) {
	var pre_req_map, cleanuptbl, cleanupstatetbl map[string]interface{}
	var expected_get_json, url string

	pfmValues := pfmMapDetails()

	/* FAN State DB test */
	t.Log("\n\n+++++++++++++ Performing Get on stateDB FAN Key ++++++++++++")
	pre_req_map = map[string]interface{}{
		"FAN": map[string]interface{}{
			"FAN-1-7": pfmValues,
		},
	}
	loadDB(hostDBName, db.ConfigDB, pre_req_map)

	pre_req_map = map[string]interface{}{
		"FAN_TABLE": map[string]interface{}{
			"FAN-1-7": pfmValues,
		},
	}
	time.Sleep(1 * time.Second)
	loadDB(hostDBName, db.StateDB, pre_req_map)
	expected_get_json = "{\"openconfig-platform:state\":{\"empty\":false,\"hardware-version\":\"test_hw_version\",\"memory\":{\"available\":\"100\",\"utilized\":\"200\"},\"mfg-date\":\"09-10-24\",\"mfg-name\":\"test_mfg_name\",\"name\":\"FAN-1-7\",\"oper-status\":\"openconfig-platform-types:ACTIVE\",\"parent\":\"ACTIVE\",\"part-no\":\"10\",\"removable\":true,\"serial-no\":\"test_serial_no\",\"software-version\":\"test_software_version\"}}"
	url = "/openconfig-platform:components/component[name=FAN-1-7]/state"
	t.Run("Test get on state DB FAN Key for "+url, processGetRequest(url, nil, expected_get_json, false))
	time.Sleep(1 * time.Second)
	t.Log("\n\n+++++++++++++ Done Performing Get on stateDB FAN Key ++++++++++++")

	//Unload the data
	cleanuptbl = map[string]interface{}{
		"FAN": map[string]interface{}{
			"FAN-1-7": "",
		},
	}
	unloadDB(hostDBName, db.ConfigDB, cleanuptbl)
	cleanupstatetbl = map[string]interface{}{
		"FAN_TABLE": map[string]interface{}{
			"FAN-1-7": "",
		},
	}
	unloadDB(hostDBName, db.StateDB, cleanupstatetbl)
}

/* PSU state DB fields and values */
func Test_psu_state_DB_key(t *testing.T) {
	var pre_req_map, cleanuptbl, cleanupstatetbl map[string]interface{}
	var expected_get_json, url string

	pfmValues := pfmMapDetails()

	/* PSU State DB test */
	t.Log("\n\n+++++++++++++ Performing Get on stateDB PSU Key ++++++++++++")
	pre_req_map = map[string]interface{}{
		"PSU": map[string]interface{}{
			"PSU-1-5": pfmValues,
		},
	}
	loadDB(hostDBName, db.ConfigDB, pre_req_map)

	pre_req_map = map[string]interface{}{
		"PSU_TABLE": map[string]interface{}{
			"PSU-1-5": pfmValues,
		},
	}
	time.Sleep(1 * time.Second)
	loadDB(hostDBName, db.StateDB, pre_req_map)
	expected_get_json = "{\"openconfig-platform:state\":{\"empty\":false,\"hardware-version\":\"test_hw_version\",\"memory\":{\"available\":\"100\",\"utilized\":\"200\"},\"mfg-date\":\"09-10-24\",\"mfg-name\":\"test_mfg_name\",\"name\":\"PSU-1-5\",\"oper-status\":\"openconfig-platform-types:ACTIVE\",\"parent\":\"ACTIVE\",\"part-no\":\"10\",\"removable\":true,\"serial-no\":\"test_serial_no\",\"software-version\":\"test_software_version\"}}"
	url = "/openconfig-platform:components/component[name=PSU-1-5]/state"
	t.Run("Test get on state DB PSU Key for "+url, processGetRequest(url, nil, expected_get_json, false))
	time.Sleep(1 * time.Second)
	t.Log("\n\n+++++++++++++ Done Performing Get on stateDB PSU Key ++++++++++++")

	//Unload the data
	cleanuptbl = map[string]interface{}{
		"PSU": map[string]interface{}{
			"PSU-1-5": "",
		},
	}
	unloadDB(hostDBName, db.ConfigDB, cleanuptbl)
	cleanupstatetbl = map[string]interface{}{
		"PSU_TABLE": map[string]interface{}{
			"PSU-1-5": "",
		},
	}
	unloadDB(hostDBName, db.StateDB, cleanupstatetbl)
}

/* CU state DB fields and values */
func Test_CU_state_DB_key(t *testing.T) {
	var pre_req_map, cleanuptbl, cleanupstatetbl map[string]interface{}
	var expected_get_json, url string

	pfmValues := pfmMapDetails()

	/* CU State DB test */
	t.Log("\n\n+++++++++++++ Performing Get on stateDB CU Key ++++++++++++")
	pre_req_map = map[string]interface{}{
		"CU": map[string]interface{}{
			"CU-1": pfmValues,
		},
	}
	loadDB(hostDBName, db.ConfigDB, pre_req_map)

	pre_req_map = map[string]interface{}{
		"CU_TABLE": map[string]interface{}{
			"CU-1": pfmValues,
		},
	}
	time.Sleep(1 * time.Second)
	loadDB(hostDBName, db.StateDB, pre_req_map)
	expected_get_json = "{\"openconfig-platform:state\":{\"empty\":false,\"hardware-version\":\"test_hw_version\",\"memory\":{\"available\":\"100\",\"utilized\":\"200\"},\"mfg-date\":\"09-10-24\",\"mfg-name\":\"test_mfg_name\",\"name\":\"CU-1\",\"oper-status\":\"openconfig-platform-types:ACTIVE\",\"parent\":\"ACTIVE\",\"part-no\":\"10\",\"removable\":true,\"serial-no\":\"test_serial_no\",\"software-version\":\"test_software_version\"}}"
	url = "/openconfig-platform:components/component[name=CU-1]/state"
	t.Run("Test get on state DB CU Key for "+url, processGetRequest(url, nil, expected_get_json, false))
	time.Sleep(1 * time.Second)
	t.Log("\n\n+++++++++++++ Done Performing Get on stateDB CU Key ++++++++++++")

	//Unload the data
	cleanuptbl = map[string]interface{}{
		"CU": map[string]interface{}{
			"CU-1": "",
		},
	}
	unloadDB(hostDBName, db.ConfigDB, cleanuptbl)
	cleanupstatetbl = map[string]interface{}{
		"CU_TABLE": map[string]interface{}{
			"CU-1": "",
		},
	}
	unloadDB(hostDBName, db.StateDB, cleanupstatetbl)
}

/* chassis state DB fields and values */
func Test_chassis_state_DB_key(t *testing.T) {
	var pre_req_map, cleanuptbl, cleanupstatetbl map[string]interface{}
	var expected_get_json, url string

	/* chassis State DB test */
	t.Log("\n\n+++++++++++++ Performing Get on stateDB chassis Key ++++++++++++")
	pre_req_map = map[string]interface{}{
		"CHASSIS": map[string]interface{}{
			"CHASSIS-1": map[string]interface{}{
				"oper-status": "ACTIVE",
			},
		},
	}
	loadDB(hostDBName, db.ConfigDB, pre_req_map)

	pre_req_map = map[string]interface{}{
		"CHASSIS_TABLE": map[string]interface{}{
			"CHASSIS-1": map[string]interface{}{
				"oper-status": "ACTIVE",
			},
		},
	}
	time.Sleep(1 * time.Second)
	loadDB(hostDBName, db.StateDB, pre_req_map)
	expected_get_json = "{\"openconfig-platform:state\":{\"name\":\"CHASSIS-1\",\"oper-status\":\"openconfig-platform-types:ACTIVE\"}}"
	url = "/openconfig-platform:components/component[name=CHASSIS-1]/state"
	t.Run("Test get on state DB chassis Key for "+url, processGetRequest(url, nil, expected_get_json, false))
	time.Sleep(1 * time.Second)
	t.Log("\n\n+++++++++++++ Done Performing Get on stateDB chassis Key ++++++++++++")

	//Unload the data
	cleanuptbl = map[string]interface{}{
		"CHASSIS": map[string]interface{}{
			"CHASSIS-1": "",
		},
	}
	unloadDB(hostDBName, db.ConfigDB, cleanuptbl)
	cleanupstatetbl = map[string]interface{}{
		"CHASSIS_TABLE": map[string]interface{}{
			"CHASSIS-1": "",
		},
	}
	unloadDB(hostDBName, db.StateDB, cleanupstatetbl)
}

func Test_set_pfm_transcevier_config_DB_key_and_field_xfmr(t *testing.T) {
	t.Log("\n\n+++++++++++++ Performing Create/Replace/Delete on pfm_transcevier ++++++++++++")

	for idx, dbName := range getMdbNames() {
		url := fmt.Sprintf("/openconfig-platform:components/component[name=TRANSCEIVER-1-%d-C2]/openconfig-platform-transceiver:transceiver/config", idx+1)
		url_body_json := "{\"config\":{\"enabled\":false,\"ethernet-pmd-preconf\":\"openconfig-transport-types:ETH_100GBASE_SR4\",\"fec-mode\":\"openconfig-platform-types:FEC_AUTO\"}}"

		pre_req_map := map[string]interface{}{
			"OC_COMPONENT": map[string]interface{}{
				fmt.Sprintf("TRANSCEIVER-1-%d-C2", idx+1): map[string]interface{}{
					"enabled":              "true",
					"fec-mode":             "FEC_AUTO",
					"ethernet-pmd-preconf": "ETH_100GBASE_SR10",
				},
			},
			"TRANSCEIVER": map[string]interface{}{
				fmt.Sprintf("TRANSCEIVER-1-%d-C2", idx+1): map[string]interface{}{
					"enabled":              "true",
					"fec-mode":             "FEC_AUTO",
					"ethernet-pmd-preconf": "ETH_100GBASE_SR10",
				},
			},
		}

		loadDB(dbName, db.ConfigDB, pre_req_map)
		expected_map := map[string]interface{}{
			"TRANSCEIVER": map[string]interface{}{
				fmt.Sprintf("TRANSCEIVER-1-%d-C2", idx+1): map[string]interface{}{
					"enabled":              "false",
					"fec-mode":             "FEC_AUTO",
					"ethernet-pmd-preconf": "ETH_100GBASE_SR4",
				},
			},
		}
		time.Sleep(1 * time.Second)

		//PUT Test
		t.Run("Replace Test on pfm_transcevier yang", processSetRequest(url, url_body_json, "PUT", false))
		time.Sleep(1 * time.Second)

		//GET Test
		t.Run("Verify replace on pfm_transcevier yang", verifyDbResult(rclient[dbName], fmt.Sprintf("TRANSCEIVER|TRANSCEIVER-1-%d-C2", idx+1), expected_map, false))
		time.Sleep(1 * time.Second)

		//Delete Test
		url = fmt.Sprintf("/openconfig-platform:components/component[name=TRANSCEIVER-1-%d-C2]/openconfig-platform-transceiver:transceiver", idx+1)
		t.Run("DELETE Test on pfm_transcevier yang", processDeleteRequest(url, false))
		time.Sleep(1 * time.Second)

		//POST Test
		url = "/openconfig-platform:components"
		post_payload := fmt.Sprintf("{\"component\":[{\"name\": \"TRANSCEIVER-1-%d-C3\", \"openconfig-platform-transceiver:transceiver\": {\"config\":{\"enabled\":false,\"ethernet-pmd-preconf\":\"openconfig-transport-types:ETH_100GBASE_SR4\",\"fec-mode\":\"openconfig-platform-types:FEC_AUTO\"}}}]}", idx+1)

		t.Run("Create test on Transceiver yang", processSetRequest(url, post_payload, "POST", false))
		time.Sleep(1 * time.Second)

		expected_map = map[string]interface{}{
			"TRANSCEIVER": map[string]interface{}{
				fmt.Sprintf("TRANSCEIVER-1-%d-C3", idx+1): map[string]interface{}{
					"enabled":              false,
					"ethernet-pmd-preconf": "ETH_100GBASE_SR4",
					"fec-mode":             "FEC_AUTO",
				},
			},
		}

		//Verify After the POST test
		t.Run("Verify After the POST test on Transceiver yang", verifyDbResult(rclient[dbName], fmt.Sprintf("TRANSCEIVER|TRANSCEIVER-1-%d-C3", idx+1), expected_map, false))

		//Delete After the POST test
		url = fmt.Sprintf("/openconfig-platform:components/component[name=TRANSCEIVER-1-%d-C3]/openconfig-platform-transceiver:transceiver", idx+1)
		t.Run("DELETE After the POST test on Transceiver yang", processDeleteRequest(url, false))
		time.Sleep(1 * time.Second)
	}
}

func Test_transceiver_allow_write_multiple_namespace(t *testing.T) {
	t.Log("\n\n+++++++++++++ Performing transceiver write for multiple namespace (should be allowed) ++++++++++++")

	url := "/openconfig-platform:components"
	url_body_json := `{
        "component": [
            {
                "name": "TRANSCEIVER-1-1-C1",
                "openconfig-platform-transceiver:transceiver": {
                    "config": {
                        "enabled": false,
                        "ethernet-pmd-preconf": "openconfig-transport-types:ETH_100GBASE_SR4",
                        "fec-mode": "openconfig-platform-types:FEC_AUTO"
                    }
                }
            },
            {
                "name": "TRANSCEIVER-1-2-C1",
                "openconfig-platform-transceiver:transceiver": {
                    "config": {
                        "enabled": false,
                        "ethernet-pmd-preconf": "openconfig-transport-types:ETH_100GBASE_SR4",
                        "fec-mode": "openconfig-platform-types:FEC_AUTO"
                    }
                }
            }
        ]
    }`

	// false = no error expected
	t.Run("Test allow writing multiple namespace for transceiver", processSetRequest(url, url_body_json, "POST", false, nil))

	t.Log("\n\n+++++++++++++ Done Performing transceiver write for multiple namespace ++++++++++++")
}
