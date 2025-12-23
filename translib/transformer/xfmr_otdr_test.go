//go:build xfmrtest
// +build xfmrtest

package transformer_test

import (
	"fmt"
	"github.com/Azure/sonic-mgmt-common/translib/db"
	"testing"
	"time"
)

/* otdr config,counters, state DB fields and values */
func Test_otdr_config_state_DB_key(t *testing.T) {
	var pre_req_map, cleanuptbl, cleanupstatetbl map[string]interface{}
	var expected_get_json, url []string

	for idx, dbName := range getMdbNames() {
		/* OTDR config DB test */
		t.Log("\n\n+++++++++++++ Performing Get on configDB OTDR Key ++++++++++++")
		pre_req_map = map[string]interface{}{
			"OTDR": map[string]interface{}{
				fmt.Sprintf("OTDR-1-%d-2", idx+1): map[string]interface{}{
					"average-time":           "2",
					"backscatter-index":      "2.0",
					"distance-range":         "2",
					"enable":                 "true",
					"end-of-fiber-threshold": "0.0",
					"name":                   fmt.Sprintf("OTDR-1-%d-2", idx+1),
					"output-frequency":       "88",
					"parent-port":            "test_parent_port",
					"period":                 "1",
					"pulse-width":            "2",
					"refractive-index":       "0.0",
					"reflection-threshold":   "0.0",
					"splice-loss-threshold":  "0.0",
					"start-time":             "0:10",
				},
			},
		}
		loadDB(dbName, db.ConfigDB, pre_req_map)
		expected_get_json = []string{fmt.Sprintf("{\"openconfig-optical-time-domain-reflectometer:config\":{\"fiber-profile\":{\"backscatter-index\":\"2\",\"end-of-fiber-threshold\":\"0\",\"reflection-threshold\":\"0\",\"refractive-index\":\"0\",\"splice-loss-threshold\":\"0\"},\"name\":\"OTDR-1-%d-2\",\"parent-port\":\"test_parent_port\",\"repetition\":{\"enable\":true,\"period\":1,\"start-time\":\"0:10\"},\"scanning-profile\":{\"average-time\":2,\"distance-range\":2,\"output-frequency\":\"88\",\"pulse-width\":2}}}", idx+1)}

		url = []string{fmt.Sprintf("/openconfig-optical-time-domain-reflectometer:otdrs/otdr[name=OTDR-1-%d-2]/config", idx+1)}
		t.Run("Test get on config DB otdr Key.", processGetRequest(url[0], nil, expected_get_json[0], false))
		time.Sleep(1 * time.Second)
		t.Log("\n\n+++++++++++++ Done Performing Get on configDB otdr Key ++++++++++++")

		/* OTDR State DB test */
		t.Log("\n\n+++++++++++++ Performing Get on stateDB otdr Key ++++++++++++")
		pre_req_map = map[string]interface{}{
			"OTDR_TABLE": map[string]interface{}{
				fmt.Sprintf("OTDR-1-%d-2", idx+1): map[string]interface{}{
					"average-time":           "2",
					"backscatter-index":      "2.0",
					"distance-accuracy":      "0.0",
					"distance-range":         "2",
					"dynamic-range":          "1",
					"enable":                 "true",
					"end-of-fiber-threshold": "0.0",
					"firmware-version":       "test_f_version",
					"hardware-version":       "test_version",
					"loss-dead-zone":         "2.0",
					"mfg-date":               "test-mfg-date",
					"mfg-name":               "test-mfg-name",
					"name":                   "OTDR-1-1-2",
					"oper-status":            "test_oper",
					"output-frequency":       "88",
					"parent":                 "test_parent",
					"parent-port":            "test_parent_port",
					"part-no":                "test_part-no",
					"period":                 "1",
					"pulse-width":            "2",
					"reflection-dead-zone":   "2.0",
					"reflection-threshold":   "0.0",
					"refractive-index":       "0.0",
					"removable":              "false",
					"sampling-resolution":    "2.0",
					"scanning-status":        "ACTIVE",
					"serial-no":              "test_serial_no",
					"software-version":       "test_soft",
					"splice-loss-threshold":  "0.0",
					"start-time":             "0:10",
				},
				fmt.Sprintf("OTDR-1-%d-2|BASELINE", idx+1): map[string]interface{}{
					"scan-time":        "2024-06-14T15:30:45.123Z+05:30",
					"span-distance":    "0.0",
					"span-loss":        "0.0",
					"distance-range":   "12",
					"pulse-width":      "12",
					"average-time":     "0",
					"output-frequency": "2",
					"update-time":      "0:0:1",
					"data":             "test_data",
				},
				fmt.Sprintf("OTDR-1-%d-2|CURRENT", idx+1): map[string]interface{}{
					"scan-time":        "2024-06-14T15:30:45.123Z+05:30",
					"span-distance":    "0.0",
					"span-loss":        "0.0",
					"distance-range":   "12",
					"pulse-width":      "12",
					"average-time":     "0",
					"output-frequency": "2",
					"update-time":      "0:0:1",
					"data":             "test_data",
				},
				fmt.Sprintf("OTDR-1-%d-2|2024-06-14T15:30:45.123Z+05:30", idx+1): map[string]interface{}{
					"scan-time":        "2024-06-14T15:30:45.123Z+05:30",
					"span-distance":    "0.0",
					"span-loss":        "0.0",
					"distance-range":   "12",
					"pulse-width":      "12",
					"average-time":     "0",
					"output-frequency": "2",
					"update-time":      "0:0:1",
					"data":             "test_data",
				},
			},
			"OTDR_EVENT_TABLE": map[string]interface{}{
				fmt.Sprintf("OTDR-1-%d-2|BASELINE|101", idx+1): map[string]interface{}{
					"index":           "101",
					"length":          "0.0",
					"loss":            "0.0",
					"accumulate-loss": "0.0",
					"type":            "REFLECTION",
					"reflection":      "0.0",
				},
				fmt.Sprintf("OTDR-1-%d-2|CURRENT|101", idx+1): map[string]interface{}{
					"index":           "101",
					"length":          "0.0",
					"loss":            "0.0",
					"accumulate-loss": "0.0",
					"type":            "REFLECTION",
					"reflection":      "0.0",
				},
				fmt.Sprintf("OTDR-1-%d-2|2024-06-14T15:30:45.123Z+05:30|101", idx+1): map[string]interface{}{
					"index":           "101",
					"length":          "0.0",
					"loss":            "0.0",
					"accumulate-loss": "0.0",
					"type":            "REFLECTION",
					"reflection":      "0.0",
				},
			},
		}
		loadDB(dbName, db.StateDB, pre_req_map)
		expected_get_json = []string{
			fmt.Sprintf("{\"openconfig-optical-time-domain-reflectometer:state\":{\"fiber-profile\":{\"backscatter-index\":\"2\",\"end-of-fiber-threshold\":\"0\",\"reflection-threshold\":\"0\",\"refractive-index\":\"0\",\"splice-loss-threshold\":\"0\"},\"name\":\"OTDR-1-%d-2\",\"parent-port\":\"test_parent_port\",\"repetition\":{\"enable\":true,\"period\":1,\"start-time\":\"0:10\"},\"scanning-profile\":{\"average-time\":2,\"distance-range\":2,\"output-frequency\":\"88\",\"pulse-width\":2},\"scanning-status\":\"openconfig-platform-types:ACTIVE\",\"specification\":{\"distance-accuracy\":\"0\",\"dynamic-range\":1,\"loss-dead-zone\":\"2\",\"reflection-dead-zone\":\"2\",\"sampling-resolution\":\"2\"}}}", idx+1),

			"{\"openconfig-optical-time-domain-reflectometer:baseline-result\":{\"events\":{\"event\":[{\"accumulate-loss\":\"0\",\"index\":101,\"length\":\"0\",\"loss\":\"0\",\"reflection\":\"0\",\"type\":\"REFLECTION\"}],\"scan-time\":\"2024-06-14T15:30:45.123Z+05:30\",\"span-distance\":\"0\",\"span-loss\":\"0\"},\"scanning-profile\":{\"average-time\":0,\"distance-range\":12,\"output-frequency\":\"2\",\"pulse-width\":12},\"trace\":{\"data\":\"test_data\",\"update-time\":\"0:0:1\"}}}",

			"{\"openconfig-optical-time-domain-reflectometer:event\":[{\"accumulate-loss\":\"0\",\"index\":101,\"length\":\"0\",\"loss\":\"0\",\"reflection\":\"0\",\"type\":\"REFLECTION\"}]}",

			"{\"openconfig-optical-time-domain-reflectometer:events\":{\"event\":[{\"accumulate-loss\":\"0\",\"index\":101,\"length\":\"0\",\"loss\":\"0\",\"reflection\":\"0\",\"type\":\"REFLECTION\"}],\"scan-time\":\"2024-06-14T15:30:45.123Z+05:30\",\"span-distance\":\"0\",\"span-loss\":\"0\"}}",

			"{\"openconfig-optical-time-domain-reflectometer:current-result\":{\"events\":{\"event\":[{\"accumulate-loss\":\"0\",\"index\":101,\"length\":\"0\",\"loss\":\"0\",\"reflection\":\"0\",\"type\":\"REFLECTION\"}],\"scan-time\":\"2024-06-14T15:30:45.123Z+05:30\",\"span-distance\":\"0\",\"span-loss\":\"0\"},\"scanning-profile\":{\"average-time\":0,\"distance-range\":12,\"output-frequency\":\"2\",\"pulse-width\":12},\"trace\":{\"data\":\"test_data\",\"update-time\":\"0:0:1\"}}}",

			"{\"openconfig-optical-time-domain-reflectometer:event\":[{\"accumulate-loss\":\"0\",\"index\":101,\"length\":\"0\",\"loss\":\"0\",\"reflection\":\"0\",\"type\":\"REFLECTION\"}]}",

			"{\"openconfig-optical-time-domain-reflectometer:events\":{\"event\":[{\"accumulate-loss\":\"0\",\"index\":101,\"length\":\"0\",\"loss\":\"0\",\"reflection\":\"0\",\"type\":\"REFLECTION\"}],\"scan-time\":\"2024-06-14T15:30:45.123Z+05:30\",\"span-distance\":\"0\",\"span-loss\":\"0\"}}",
		}

		url = []string{
			fmt.Sprintf("/openconfig-optical-time-domain-reflectometer:otdrs/otdr[name=OTDR-1-%d-2]/state", idx+1),
			fmt.Sprintf("/openconfig-optical-time-domain-reflectometer:otdrs/otdr[name=OTDR-1-%d-2]/baseline-result", idx+1),
			fmt.Sprintf("/openconfig-optical-time-domain-reflectometer:otdrs/otdr[name=OTDR-1-%d-2]/baseline-result/events/event[index=101]", idx+1),
			fmt.Sprintf("/openconfig-optical-time-domain-reflectometer:otdrs/otdr[name=OTDR-1-%d-2]/baseline-result/events", idx+1),
			fmt.Sprintf("/openconfig-optical-time-domain-reflectometer:otdrs/otdr[name=OTDR-1-%d-2]/current-result", idx+1),
			fmt.Sprintf("/openconfig-optical-time-domain-reflectometer:otdrs/otdr[name=OTDR-1-%d-2]/current-result/events/event[index=101]", idx+1),
			fmt.Sprintf("/openconfig-optical-time-domain-reflectometer:otdrs/otdr[name=OTDR-1-%d-2]/current-result/events", idx+1),
		}
		for idx, url := range url {
			t.Run("Test get on state DB otdr Key for "+url, processGetRequest(url, nil, expected_get_json[idx], false))
			time.Sleep(1 * time.Second)
		}

		t.Log("\n\n+++++++++++++ Done Performing Get on stateDB otdr Key ++++++++++++")
		//Unload the Data
		cleanuptbl = map[string]interface{}{
			"OTDR": map[string]interface{}{
				fmt.Sprintf("OTDR-1-%d-2", idx+1): "",
			},
		}
		unloadDB(dbName, db.ConfigDB, cleanuptbl)

		cleanupstatetbl = map[string]interface{}{
			"OTDR_TABLE": map[string]interface{}{
				fmt.Sprintf("OTDR-1-%d-2", idx+1):                                "",
				fmt.Sprintf("OTDR-1-%d-2|BASELINE", idx+1):                       "",
				fmt.Sprintf("OTDR-1-%d-2|2024-06-14T15:30:45.123Z+05:30", idx+1): "",
				fmt.Sprintf("OTDR-1-%d-2|CURRENT", idx+1):                        "",
			},
			"OTDR_EVENT_TABLE": map[string]interface{}{
				fmt.Sprintf("OTDR-1-%d-2|CURRENT|101", idx+1):                        "",
				fmt.Sprintf("OTDR-1-%d-2|BASELINE|101", idx+1):                       "",
				fmt.Sprintf("OTDR-1-%d-2|2024-06-14T15:30:45.123Z+05:30|101", idx+1): "",
			},
		}
		unloadDB(dbName, db.StateDB, cleanupstatetbl)
	}
}

func Test_set_OTDR_config_DB_key_and_field_xfmr(t *testing.T) {
	t.Log("\n\n+++++++++++++ Performing Create/Replace/Delete on OTDR ++++++++++++")

	for idx, dbName := range getMdbNames() {
		url := fmt.Sprintf("/openconfig-optical-time-domain-reflectometer:otdrs/otdr[name=OTDR-1-%d-2]/config", idx+1)
		url_body_json := "{\"config\":{\"fiber-profile\":{\"backscatter-index\":\"28.78\",\"refractive-index\":\"11.11\",\"reflection-threshold\":\"22.22\",\"splice-loss-threshold\":\"33.33\",\"end-of-fiber-threshold\":\"44.44\"},\"scanning-profile\":{\"distance-range\":8,\"pulse-width\":7,\"average-time\":9,\"output-frequency\":\"20\"},\"repetition\":{\"period\":10,\"start-time\":\"2024-07-03T14:30:00.123Z+03:00\",\"enable\":false}}}"

		pre_req_map := map[string]interface{}{
			"OTDR": map[string]interface{}{
				fmt.Sprintf("OTDR-1-%d-2", idx+1): map[string]interface{}{
					"average-time":           "2",
					"backscatter-index":      "2.0",
					"distance-range":         "2",
					"enable":                 "true",
					"end-of-fiber-threshold": "0.0",
					"output-frequency":       "88",
					"period":                 "1",
					"pulse-width":            "2",
					"refractive-index":       "0.0",
					"reflection-threshold":   "0.0",
					"splice-loss-threshold":  "0.0",
					"start-time":             "0:10",
				},
			},
		}

		loadDB(dbName, db.ConfigDB, pre_req_map)
		expected_map := map[string]interface{}{
			"OTDR": map[string]interface{}{
				fmt.Sprintf("OTDR-1-%d-2", idx+1): map[string]interface{}{
					"average-time":           "9",
					"backscatter-index":      "28.78",
					"distance-range":         "8",
					"enable":                 "false",
					"end-of-fiber-threshold": "44.44",
					"output-frequency":       "20",
					"period":                 "10",
					"pulse-width":            "7",
					"refractive-index":       "11.11",
					"reflection-threshold":   "22.22",
					"splice-loss-threshold":  "33.33",
					"start-time":             "2024-07-03T14:30:00.123Z+03:00",
				},
			},
		}
		time.Sleep(1 * time.Second)

		//PUT Test
		t.Run("Replace Test on OTDR yang", processSetRequest(url, url_body_json, "PUT", false))
		time.Sleep(1 * time.Second)

		//GET Test
		t.Run("Verify replace on OTDR yang", verifyDbResult(rclient[dbName], fmt.Sprintf("OTDR|OTDR-1-%d-2", idx+1), expected_map, false))
		time.Sleep(1 * time.Second)

		//Delete Test
		url = fmt.Sprintf("/openconfig-optical-time-domain-reflectometer:otdrs/otdr[name=OTDR-1-%d-2]/config", idx+1)
		t.Run("DELETE Test on OTDR yang", processDeleteRequest(url, false))
		time.Sleep(1 * time.Second)

		//POST Test
		url = "/openconfig-optical-time-domain-reflectometer:otdrs"
		post_payload := fmt.Sprintf("{\"otdr\":[{\"name\":\"OTDR-1-%d-3\",\"config\":{\"fiber-profile\":{\"backscatter-index\":\"18.18\",\"refractive-index\":\"11.11\",\"reflection-threshold\":\"22.22\",\"splice-loss-threshold\":\"33.33\",\"end-of-fiber-threshold\":\"44.44\"},\"scanning-profile\":{\"distance-range\":15,\"pulse-width\":7,\"average-time\":10,\"output-frequency\":\"31\"},\"repetition\":{\"period\":10,\"start-time\":\"2024-07-03T14:30:00.123Z+03:00\",\"enable\":true}}}]}", idx+1)
		t.Run("Create test on OTDR yang", processSetRequest(url, post_payload, "POST", false))
		time.Sleep(1 * time.Second)

		expected_map = map[string]interface{}{
			"OTDR": map[string]interface{}{
				fmt.Sprintf("OTDR-1-%d-3", idx+1): map[string]interface{}{
					"average-time":           "10",
					"backscatter-index":      "18.18",
					"distance-range":         "15",
					"enable":                 "true",
					"end-of-fiber-threshold": "44.44",
					"output-frequency":       "31",
					"period":                 "10",
					"pulse-width":            "7",
					"refractive-index":       "11.11",
					"reflection-threshold":   "22.22",
					"splice-loss-threshold":  "33.33",
					"start-time":             "2024-07-03T14:30:00.123Z+03:00",
				},
			},
		}

		//Verify After the POST test
		t.Run("Verify After the POST test on OTDR yang", verifyDbResult(rclient[dbName], fmt.Sprintf("OTDR|OTDR-1-%d-3", idx+1), expected_map, false))

		//Delete After the POST test
		url = fmt.Sprintf("/openconfig-optical-time-domain-reflectometer:otdrs/otdr[name=OTDR-1-%d-3]/config", idx+1)
		t.Run("DELETE After the POST test on OTDR yang", processDeleteRequest(url, false))
		time.Sleep(1 * time.Second)

		url = "/openconfig-optical-time-domain-reflectometer:otdrs"

		//Bulk POST Test
		Bulk_post_payload := fmt.Sprintf("{\"otdr\":[{\"name\":\"OTDR-1-%d-4\",\"config\":{\"fiber-profile\":{\"backscatter-index\":\"18.18\",\"refractive-index\":\"11.11\",\"reflection-threshold\":\"22.22\",\"splice-loss-threshold\":\"33.33\",\"end-of-fiber-threshold\":\"44.44\"},\"scanning-profile\":{\"distance-range\":15,\"pulse-width\":7,\"average-time\":10,\"output-frequency\":\"31\"},\"repetition\":{\"period\":10,\"start-time\":\"2024-07-03T14:30:00.123Z+03:00\",\"enable\":true}}},{\"name\":\"OTDR-1-%d-5\",\"config\":{\"fiber-profile\":{\"backscatter-index\":\"28.78\",\"refractive-index\":\"11.11\",\"reflection-threshold\":\"22.22\",\"splice-loss-threshold\":\"33.33\",\"end-of-fiber-threshold\":\"44.44\"},\"scanning-profile\":{\"distance-range\":8,\"pulse-width\":7,\"average-time\":9,\"output-frequency\":\"20\"},\"repetition\":{\"period\":10,\"start-time\":\"2024-07-03T14:30:00.123Z+03:00\",\"enable\":false}}}]}", idx+1, idx+1)
		t.Run("Create test on terminal_device yang", processSetRequest(url, Bulk_post_payload, "POST", false))

		expected_map = map[string]interface{}{
			"OTDR": map[string]interface{}{
				fmt.Sprintf("OTDR-1-%d-4", idx+1): map[string]interface{}{
					"average-time":           "10",
					"backscatter-index":      "18.18",
					"distance-range":         "15",
					"enable":                 "true",
					"end-of-fiber-threshold": "44.44",
					"output-frequency":       "31",
					"period":                 "10",
					"pulse-width":            "7",
					"refractive-index":       "11.11",
					"reflection-threshold":   "22.22",
					"splice-loss-threshold":  "33.33",
					"start-time":             "2024-07-03T14:30:00.123Z+03:00",
				},
			},
		}

		time.Sleep(1 * time.Second)
		//Verify After the Bulk test
		t.Run("Verify After the Bulk test on OTDR yang", verifyDbResult(rclient[dbName], fmt.Sprintf("OTDR|OTDR-1-%d-4", idx+1), expected_map, false))

		expected_map = map[string]interface{}{
			"OTDR": map[string]interface{}{
				fmt.Sprintf("OTDR-1-%d-5", idx+1): map[string]interface{}{
					"average-time":           "9",
					"backscatter-index":      "28.78",
					"distance-range":         "8",
					"enable":                 "false",
					"end-of-fiber-threshold": "44.44",
					"output-frequency":       "20",
					"period":                 "10",
					"pulse-width":            "7",
					"refractive-index":       "11.11",
					"reflection-threshold":   "22.22",
					"splice-loss-threshold":  "33.33",
					"start-time":             "2024-07-03T14:30:00.123Z+03:00",
				},
			},
		}

		time.Sleep(1 * time.Second)
		//Verify After the POST test
		t.Run("Verify After the Bulk test on OTDR yang", verifyDbResult(rclient[dbName], fmt.Sprintf("OTDR|OTDR-1-%d-5", idx+1), expected_map, false))

		//Delete keys After the Bulk Test
		url = fmt.Sprintf("/openconfig-optical-time-domain-reflectometer:otdrs/otdr[name=OTDR-1-%d-4]/config", idx+1)
		t.Run("DELETE After the Bulk Test on OTDR yang", processDeleteRequest(url, false))
		time.Sleep(1 * time.Second)

		url = fmt.Sprintf("/openconfig-optical-time-domain-reflectometer:otdrs/otdr[name=OTDR-1-%d-5]/config", idx+1)
		t.Run("DELETE After the Bulk Test on OTDR yang", processDeleteRequest(url, false))
		time.Sleep(1 * time.Second)
	}

	t.Log("\n\n+++++++++++++ Done Performing Create/Replace/Delete on OTDR ++++++++++++")
}

func Test_OTDR_allow_write_multiple_namespace(t *testing.T) {
	t.Log("\n\n+++++++++++++ Performing OTDR write for multiple namespace (should be allowed) ++++++++++++")

	url := "/openconfig-optical-time-domain-reflectometer:otdrs"
	url_body_json := `{
        "otdr": [
            {
                "name": "OTDR-1-1-4",
                "config": {
                    "fiber-profile": {
                        "backscatter-index": "18.18",
                        "refractive-index": "11.11",
                        "reflection-threshold": "22.22",
                        "splice-loss-threshold": "33.33",
                        "end-of-fiber-threshold": "44.44"
                    },
                    "scanning-profile": {
                        "distance-range": 15,
                        "pulse-width": 7,
                        "average-time": 10,
                        "output-frequency": "31"
                    },
                    "repetition": {
                        "period": 10,
                        "start-time": "2024-07-03T14:30:00.123Z+03:00",
                        "enable": true
                    }
                }
            },
            {
                "name": "OTDR-1-3-5",
                "config": {
                    "fiber-profile": {
                        "backscatter-index": "28.78",
                        "refractive-index": "11.11",
                        "reflection-threshold": "22.22",
                        "splice-loss-threshold": "33.33",
                        "end-of-fiber-threshold": "44.44"
                    },
                    "scanning-profile": {
                        "distance-range": 8,
                        "pulse-width": 7,
                        "average-time": 9,
                        "output-frequency": "20"
                    },
                    "repetition": {
                        "period": 10,
                        "start-time": "2024-07-03T14:30:00.123Z+03:00",
                        "enable": false
                    }
                }
            }
        ]
    }`

	// false = no error expected
	t.Run("Test allow writing multiple namespace for OTDR", processSetRequest(url, url_body_json, "POST", false, nil))

	t.Log("\n\n+++++++++++++ Done Performing OTDR write for multiple namespace ++++++++++++")
}
