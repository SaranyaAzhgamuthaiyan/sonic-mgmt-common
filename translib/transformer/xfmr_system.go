package transformer

import (
	"encoding/json"
	"fmt"
	"github.com/Azure/sonic-mgmt-common/translib/db"
	"github.com/Azure/sonic-mgmt-common/translib/tlerr"
	log "github.com/golang/glog"
)

var supportRebootMethod map[string]int

func init () {
	XlateFuncBind("rpc_reboot", rpc_reboot)

	supportRebootMethod = make(map[string]int)
	supportRebootMethod["UNKNOWN"] = 1
	supportRebootMethod["COLD"] = 1
	supportRebootMethod["POWERDOWN"] = 1
	supportRebootMethod["HALT"] = 1
	supportRebootMethod["WARM"] = 1
	supportRebootMethod["NSF"] = 1
	supportRebootMethod["RESET"] = 1
	supportRebootMethod["POWERUP"] = 1
}

/* RPC for reboot */
var rpc_reboot RpcCallpoint = func(body []byte, dbs [db.MaxDB]*db.DB) ([]byte, error) {
	var err error
	var result struct {
		Output struct {
			Message string `json:"message"`
		} `json:"openconfig-system:output"`
	}

	/* Get input data */
	var mapData map[string]interface{}
	err = json.Unmarshal(body, &mapData)
	if err != nil {
		log.Info("Failed to unmarshall given input data")
		result.Output.Message = fmt.Sprintf("Error: Failed to unmarshall given input data")
		return json.Marshal(&result)
	}

	input, ok := mapData["openconfig-system:input"]
	if ok == false {
		return nil, tlerr.InvalidArgs("payload don't has openconfig-system:input node")
	}

	mapData = input.(map[string]interface{})

	if input, ok = mapData["method"]; ok == false {
		return nil, tlerr.InvalidArgs("payload don't has method node")
	}
	method := input.(string)
	if _, ok := supportRebootMethod[method]; ok == false {
		return nil, tlerr.NotSupported("method '%s' not supported", method)
	}

	if input, ok = mapData["message"]; ok == false {
		return nil, tlerr.InvalidArgs("payload don't has message node")
	}
	message := input.(string)

	if input, ok = mapData["entity-name"]; ok == false {
		return nil, tlerr.InvalidArgs("payload don't has entity-name node")
	}
	entityName := input.(string)

	log.Infof("oc-system:reboot method = %s, message = %s, entityName = %s", method, message, entityName)

	result.Output.Message = "system reboot succeed."

	return json.Marshal(&result)
}

