package transformer

import (
	"errors"
	"fmt"
	"github.com/Azure/sonic-mgmt-common/translib/db"
	log "github.com/golang/glog"
	"strconv"
	"strings"
)

const (
	OC_COMPONENT_TBL = "OC_COMPONENT"
	OCH_TBL          = "OCH"
	TRANSCEIVER_TBL  = "TRANSCEIVER"
)

func init() {
	/* Table transformer for tables in component container*/
	XlateFuncBind("pfm_table_xfmr", pfm_table_xfmr)

	/* Key transformer for OC_COMPONENT_TABLE table*/
	XlateFuncBind("YangToDb_component_key_xfmr", YangToDb_component_key_xfmr)
	XlateFuncBind("DbToYang_component_key_xfmr", DbToYang_component_key_xfmr)

	/* Field transformer for OC_COMPONENT_TABLE table*/
	XlateFuncBind("YangToDb_component_field_xfmr", YangToDb_component_field_xfmr)
	XlateFuncBind("DbToYang_component_field_xfmr", DbToYang_component_field_xfmr)

	/* Key transformer for transceiver table*/
	XlateFuncBind("YangToDb_pfm_transceiver_key_xfmr", YangToDb_pfm_transceiver_key_xfmr)
	XlateFuncBind("DbToYang_pfm_transceiver_key_xfmr", DbToYang_pfm_transceiver_key_xfmr)

	XlateFuncBind("pfm_post_xfmr", pfm_post_xfmr)

}

var YangToDb_component_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.Infof("YangToDb_component_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)
	ochkey := pathInfo.Var("name")
	log.Infof("YangToDb_component_key_xfmr: :", ochkey)

	return ochkey, nil
}

var DbToYang_component_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{}, 1)
	log.Infof("DbToYang_component_key_xfmr: ", inParams.key)
	rmap["name"] = inParams.key

	return rmap, nil
}

var YangToDb_component_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {

	rmap := make(map[string]string)
	rmap["NULL"] = "NULL"

	return rmap, nil
}

var DbToYang_component_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{})
	rmap["name"] = inParams.key

	return rmap, nil
}

var pfm_table_xfmr TableXfmrFunc = func(inParams XfmrParams) ([]string, error) {
	var tblList []string
	pathInfo := NewPathInfo(inParams.uri)
	name := pathInfo.Var("name")

	log.Infof("Pfm_table_xfmr inParams.d:%v ", inParams.d)

	// Determine the suffix based on DBNo
	suffix := ""
	if inParams.d.Opts.DBNo == db.StateDB {
		suffix = "_TABLE"
	}

	// Base entries for the case when name is empty
	baseEntries := []string{"FAN", "PSU", "CU", "CHASSIS", "OC_COMPONENT"}

	if len(name) == 0 {
		// Append the suffix
		tblName := "OC_COMPONENT" + suffix
		tblList = append(tblList, tblName)
		log.Info("Pfm_table_xfmr NO KEY tblList= ", tblList)
		return tblList, nil
	}

	// Check the prefix and append the corresponding entry
	for _, entry := range baseEntries {
		if strings.HasPrefix(name, entry) {
			tblList = append(tblList, entry)
			break
		}
	}

	// If no prefix matched, default to OC_COMPONENT
	if len(tblList) == 0 {
		tblList = append(tblList, "OC_COMPONENT")
	}

	// Append "_TABLE" if its STATE_DB
	if len(tblList) > 0 && suffix != "" {
		tblList[len(tblList)-1] += suffix
	}

	log.Info("Pfm_table_xfmr tblList= ", tblList)
	return tblList, nil
}

var YangToDb_pfm_transceiver_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.Infof("YangToDb_pfm_transceiver_key_xfmr : root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	key := pathInfo.Var("name")
	index := pathInfo.Var("index")

	if strings.Contains(inParams.uri, "/channel") {
		if index != "" {
			if strings.Contains(index, "CH") {
				key += "|" + index
			} else {
				key += "|CH-" + index
			}
		}
	}

	log.Infof("YangToDb_pfm_transceiver_key_xfmr : key:", key)

	return key, nil
}

var DbToYang_pfm_transceiver_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.Info("DbToYang_pfm_transceiver_key_xfmr:  ", inParams.key)
	rmap := make(map[string]interface{}, 1)
	if strings.Contains(inParams.uri, "/channel") {
		TableKeys := strings.Split(inParams.key, "|")
		if len(TableKeys) < 2 {
			log.Errorf("length of Table Key is less than 2: %v", TableKeys)
			return nil, errors.New("length of Table Key is less than 2")
		}

		indexKey := strings.Split(TableKeys[1], "-")
		if len(indexKey) < 2 {
			log.Errorf("length of Index Key is less than 2: %v", indexKey)
			return nil, errors.New("length of Index Key is less than 2")
		}

		index, err := strconv.ParseUint(indexKey[1], 10, 16)
		if err != nil {
			log.Errorf("Error parsing index: %v", err)
			return nil, err
		}
		rmap["index"] = index
	} else {
		rmap["name"] = inParams.key
	}

	log.Info("DbToYang_pfm_transceiver_key_xfmr:  rmap", rmap)
	return rmap, nil
}

var pfm_post_xfmr PostXfmrFunc = func(inParams XfmrParams) (map[string]map[string]db.Value, error) {

	if inParams.dbDataMap == nil || *inParams.dbDataMap == nil {
		return nil, fmt.Errorf("dbDataMap is nil")
	}

	retDbDataMap := (*inParams.dbDataMap)[inParams.curDb]
	if retDbDataMap == nil {
		retDbDataMap = make(map[string]map[string]db.Value)
		(*inParams.dbDataMap)[inParams.curDb] = retDbDataMap
	}

	log.Infof("pfm_post_xfmr: Entering Request URI path = %v", inParams.requestUri)

	// ---------- CREATE ----------
	if inParams.oper == CREATE {
		if strings.Contains(inParams.requestUri, "components") {
			for _, table := range []string{OCH_TBL, TRANSCEIVER_TBL} {
				if retDbDataMap[table] == nil {
					continue
				}

				if retDbDataMap[OC_COMPONENT_TBL] == nil {
					retDbDataMap[OC_COMPONENT_TBL] = make(map[string]db.Value)
				}

				for keyName := range retDbDataMap[table] {
					// Only add if not already present in dbDataMap
					if _, exists := retDbDataMap[OC_COMPONENT_TBL][keyName]; !exists {
						log.Infof("pfm_post_xfmr: Adding key %v to OC_COMPONENT_TBL", keyName)
						retDbDataMap[OC_COMPONENT_TBL][keyName] = db.Value{} // empty value marks new entry
					}
				}
			}
		}
	}

	// ---------- DELETE ----------
	if inParams.oper == DELETE {
		pathInfo := NewPathInfo(inParams.requestUri)
		keyName := pathInfo.Var("name")

		// Determine table based on URI (similar to original logic)
		var tableTsName string
		switch {
		case strings.Contains(inParams.requestUri, "OCH") && !strings.Contains(inParams.requestUri, "optical-channel"):
			tableTsName = OCH_TBL
		case strings.Contains(inParams.requestUri, "OCH") && strings.Contains(inParams.requestUri, "optical-channel"):
			tableTsName = OC_COMPONENT_TBL
		case strings.Contains(inParams.requestUri, "TRANSCEIVER") && strings.Contains(inParams.requestUri, "transceiver"):
			tableTsName = OC_COMPONENT_TBL
		case strings.Contains(inParams.requestUri, "TRANSCEIVER") && !strings.Contains(inParams.requestUri, "transceiver"):
			tableTsName = TRANSCEIVER_TBL
		}

		if tableTsName != "" {
			if retDbDataMap[tableTsName] == nil {
				retDbDataMap[tableTsName] = make(map[string]db.Value)
			}

			if keyName != "" {
				log.Infof("pfm_post_xfmr: Marking key %v in %v for deletion", keyName, tableTsName)
				retDbDataMap[tableTsName][keyName] = db.Value{} // marks deletion
			} else {
				log.Infof("pfm_post_xfmr: Clearing all entries in table %v", tableTsName)
				retDbDataMap[tableTsName] = make(map[string]db.Value) // clears table
			}
		}
	}

	log.Infof("pfm_post_xfmr: Updated dbDataMap: %+v", retDbDataMap)
	return retDbDataMap, nil
}
