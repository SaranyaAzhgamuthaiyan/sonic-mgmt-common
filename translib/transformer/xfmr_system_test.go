//go:build xfmrtest
// +build xfmrtest

package transformer_test

import (
	"testing"
	"time"

	"github.com/Azure/sonic-mgmt-common/translib/db"
)

/* ALARM counters, state DB fields and values */
func Test_alarm_state_DB_key(t *testing.T) {
	var pre_req_map, cleanupstatetbl map[string]interface{}
	var url string

	/* PORT State DB test */
	t.Log("\n\n+++++++++++++ Performing Get on stateDB ALARM Key ++++++++++++")
	pre_req_map = map[string]interface{}{
		"CURALARM_TABLE": map[string]interface{}{
			"ALARM1": map[string]interface{}{
				"time-created": "12",
				"resource":     "test_res",
				"text":         "test_text",
				"severity":     "CRITICAL",
				"type-id":      "EQPT",
			},
		},
	}
	loadDB(hostDBName, db.StateDB, pre_req_map)
	expected_get_json := "{\"openconfig-system:alarm\":[{\"id\":\"ALARM1\",\"state\":{\"id\":\"ALARM1\",\"resource\":\"test_res\",\"severity\":\"openconfig-alarm-types:CRITICAL\",\"text\":\"test_text\",\"time-created\":\"12\",\"type-id\":\"openconfig-alarm-types:EQPT\"}}]}"
	url = "/openconfig-system:system/alarms/alarm[id=ALARM1]"
	t.Run("Test get on state DB ALARM Key", processGetRequest(url, nil, expected_get_json, false))
	time.Sleep(1 * time.Second)
	t.Log("\n\n+++++++++++++ Done Performing Get on stateDB ALARM Key ++++++++++++")

	// Unload the Data
	cleanupstatetbl = map[string]interface{}{
		"CURALARM_TABLE": map[string]interface{}{
			"ALARM1": "",
		},
	}
	unloadDB(hostDBName, db.StateDB, cleanupstatetbl)
}
