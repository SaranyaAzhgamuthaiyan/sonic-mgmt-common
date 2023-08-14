package translib

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Azure/sonic-mgmt-common/translib/db"
	"github.com/Azure/sonic-mgmt-common/translib/ocbinds"
	"github.com/Azure/sonic-mgmt-common/translib/tlerr"
	"github.com/golang/glog"
	"github.com/openconfig/ygot/util"
	"github.com/openconfig/ygot/ygot"
	"github.com/openconfig/ygot/ytypes"
	"math"
	"reflect"
	"strconv"
	"strings"
)

type regexPathKeyParams struct {
	tableName    string
	listNodeName []string
	keyName      []string
	redisPrefix  []string
}

/*
 * translib northbound api
 */

func transToUint64(value interface{}) (uint64, error) {
	switch val := value.(type) {
	case string:
		tmp, err := strconv.ParseInt(value.(string), 10, 64)
		if err != nil {
			glog.Errorf("body field is not uint type, err %v", err)
			return 0,err
		}
		return uint64(tmp), nil
	case float64:
		return uint64(val), nil
	default:
		return 0, errors.New("invalid type")
	}
}

func parsePostBody(reqData []byte, acceptStruct interface{}) error {
	body := make(map[string]interface{})
	err := json.Unmarshal(reqData, &body)
	if err != nil {
		glog.Errorf("decode post body failed as %v", err)
		return err
	}

	input := body["input"].(map[string]interface{})
	for k, v := range input {
		cType := reflect.TypeOf(acceptStruct).Elem()
		cValue := reflect.ValueOf(acceptStruct).Elem()
		for i := 0; i < cType.NumField(); i++ {
			fType := cType.Field(i)
			fValue := cValue.Field(i)
			if fType.Tag.Get("json") != k {
				continue
			}

			if fValue.IsNil() {
				fValue.Set(reflect.New(fType.Type.Elem()))
			}

			switch fType.Type.String() {
			case "*string":
				fValue.Elem().SetString(v.(string))
			case "*uint16":
				tmp, err := transToUint64(v)
				if err != nil {
					glog.Errorf("body field is not uint type, err %v", err)
					return err
				}
				fValue.Elem().SetUint(uint64(tmp))
			default:
				glog.V(2).Infof("input field %s type %d is not supported", k, reflect.TypeOf(v).Kind())
			}
		}
	}

	return nil
}

func float32StringToByte(floatStr string) ([]byte, error) {
	f, err := strconv.ParseFloat(floatStr, 32)
	if err != nil {
		return nil, err
	}
	bits := math.Float32bits(float32(f))
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, bits)

	return bytes, err
}

func buildGoStruct(gs interface{}, data db.Value) {
	if reflect.ValueOf(gs).IsNil() || len(data.Field) == 0 {
		glog.V(1).Infof("buildGoStruct with invalid params for %s", reflect.TypeOf(gs).String())
		return
	}

	rt := reflect.TypeOf(gs).Elem()
	rv := reflect.ValueOf(gs).Elem()
	for i := 0; i < rt.NumField(); i++ {
		fType := rt.Field(i)
		fVal := rv.Field(i)
		val := data.Get(fType.Tag.Get("path"))
		if len(val) == 0 {
			continue
		}

		if util.IsTypeStructPtr(fType.Type) {
			//err = buildGoStruct(fVal.Interface(), data)
		} else if util.IsTypeMap(fType.Type) {
			glog.Error("map type need to support")
		} else if util.IsTypeStruct(fType.Type) {
			glog.Error("struct type need to support")
		} else if util.IsTypeInterface(fType.Type) {
			glog.V(1).Infof("GoStruct %s field %s type %s is supported separately in app.", reflect.TypeOf(gs).String(), fType.Name, fType.Type.Name())
		} else {
			err := buildGoStructField(fType, fVal, val)
			if err != nil {
				glog.Errorf("build GoStruct %s failed as %v", reflect.TypeOf(gs).String(), err)
			}
		}
	}
}

func buildGoStructField(sf reflect.StructField, rv reflect.Value, val string) error {
	if len(val) == 0 {
		return nil
	}

	k := sf.Type.Kind()
	if k == reflect.Ptr {
		et := sf.Type.Elem()
		bitSize := int(et.Size()) * 8
		switch et.Kind() {
		case reflect.Float32, reflect.Float64:
			tmp, err := strconv.ParseFloat(val, bitSize)
			if err != nil {
				return fmt.Errorf("%s has invalid value %s in database", sf.Name, val)
			}
			// initialize the leaf node which is Ptr
			rv.Set(reflect.New(sf.Type.Elem()))
			rv.Elem().SetFloat(tmp)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			tmp, err := strconv.ParseInt(val, 10, bitSize)
			if err != nil {
				return fmt.Errorf("%s has invalid value %s in database", sf.Name, val)
			}
			// initialize the leaf node which is Ptr
			rv.Set(reflect.New(sf.Type.Elem()))
			rv.Elem().SetInt(tmp)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			tmp, err := strconv.ParseUint(val, 10, bitSize)
			if err != nil {
				return fmt.Errorf("%s has invalid value %s in database", sf.Name, val)
			}
			// initialize the leaf node which is Ptr
			rv.Set(reflect.New(sf.Type.Elem()))
			rv.Elem().SetUint(tmp)
		case reflect.Bool:
			tmp, err := strconv.ParseBool(val)
			if err != nil {
				return fmt.Errorf("%s has invalid value %s in database", sf.Name, val)
			}
			// initialize the leaf node which is Ptr
			rv.Set(reflect.New(sf.Type.Elem()))
			rv.Elem().SetBool(tmp)
		case reflect.String:
			// initialize the leaf node which is Ptr
			rv.Set(reflect.New(sf.Type.Elem()))
			rv.Elem().SetString(val)
		default:
			return fmt.Errorf("field %s type %s is not supported", sf.Name, sf.Type.Name())
		}
	} else if k == reflect.Int64 {
		if enum, isEnum := rv.Interface().(ygot.GoEnum); isEnum {

			tmp, err := getEnumVal(enum, val, sf.Name)
			if err != nil {
				return fmt.Errorf("%s has invalid value %s in database", sf.Name, val)
			}

			rv.SetInt(tmp)
		} else {
			glog.Errorf("build field %s failed", sf.Name)
		}
	} else if k == reflect.Slice {
		//only support binary now
		switch rv.Interface().(type) {
		case ocbinds.Binary:
			byteVal, err := float32StringToByte(val)
			if err != nil {
				glog.Errorf("ieeefloat32 encode failed as %v", err)
			}
			rv.SetBytes(byteVal)
		default:
			glog.Error("slice type need to support")
		}
	} else {
		glog.Errorf("field %s type %s is not supported", sf.Name, sf.Type.Name())
	}

	return nil
}

func constructCountersTableKey(parentKey db.Key, suffix string, lastKey string) db.Key {
	var keys []string
	for i := 0; i < parentKey.Len(); i++ {
		elmt := parentKey.Get(i)
		if i + 1 == parentKey.Len() && len(suffix) != 0 {
			elmt += "_" + suffix
		}
		keys = append(keys, elmt)
	}

	keys = append(keys, lastKey)
	return db.Key{Comp: keys}
}

func buildEnclosedCountersNodes(gs interface{}, dbCl *db.DB, ts *db.TableSpec, key db.Key) error {
	if gs == nil || key.Len() == 0 {
		return tlerr.InvalidArgs("input params invalid")
	}

	if v, ok := gs.(ygot.GoStruct); ok {
		ygot.BuildEmptyTree(v)
	}

    // 统计值
	dbKey := constructCountersTableKey(key, "", PMCurrent)
	if data, err := getRedisData(dbCl, ts, dbKey); err == nil && data.IsPopulated() {
		buildGoStruct(gs, data)
	}

	rt := reflect.TypeOf(gs).Elem()
	rv := reflect.ValueOf(gs).Elem()
	for i := 0; i < rt.NumField(); i++ {
		fType := rt.Field(i)
		fVal := rv.Field(i)

		if _, ok := fVal.Interface().(ygot.GoEnum);  ok {
			continue
		}

		// OCH:OCH-1-1-L2_OutputPower:15_pm_current
		dbKey := constructCountersTableKey(key, fType.Name, PMCurrent15min)
		if data, err := getRedisData(dbCl, ts, dbKey); err == nil {
			if util.IsTypeStructPtr(fType.Type) {
				// 模拟值 container
				buildGoStruct(fVal.Interface(), data)
			} else {
				fieldVal := data.Get("instant")
				if len(fieldVal) == 0 {
					glog.Errorf("get instant field from table %s failed as %s", ts.Name, err)
					continue
				}
				// 模拟值 leaf
				err = buildGoStructField(fType, fVal, fieldVal)
				if err != nil {
					glog.Errorf("build field %s failed as %v", fType.Name, err)
					continue
				}
			}
		}
	}

	return nil
}

func getEnumVal(enum ygot.GoEnum, val string, name string) (int64, error) {
	e := reflect.ValueOf(enum)
	lookup, ok := enum.ΛMap()[e.Type().Name()]
	if !ok {
		return 0, fmt.Errorf("build %s failed for %s", val, name)
	}

	for idx, ed := range lookup {
		if ed.Name == val {
			return idx, nil
		}
	}

	return 0, fmt.Errorf("%s is invalid", val)
}

func requestBodyHasField(t reflect.Type, v reflect.Value) bool {
	if t.Kind() == reflect.Ptr && !v.IsNil() {
		return true
	}

	_, isEnum := v.Interface().(ygot.GoEnum)
	if isEnum && !v.IsZero() {
		return true
	}

	return false
}

func convertRequestBodyToInternal(gs interface{}) db.Value {
	rt := reflect.TypeOf(gs)
	rv := reflect.ValueOf(gs)
	data := db.Value{Field: map[string]string{}}
	convert(rt, rv, &data)
	return data
}
func convert(rt reflect.Type, rv reflect.Value, data *db.Value) error {
	var err error
	if rt.Elem().Kind() == reflect.Struct {
		for i := 0; i < rt.Elem().NumField(); i++ {
			if !requestBodyHasField(rt.Elem().Field(i).Type, rv.Elem().Field(i)) {
				continue
			}
			if rt.Elem().Field(i).Type.Elem().Kind() == reflect.Struct {
				fType := rt.Elem().Field(i).Type
				fVal := rv.Elem().Field(i)
				err = convert(fType, fVal, data)
			} else {
				fType := rt.Elem().Field(i).Type.Elem()
				fVal := rv.Elem().Field(i).Elem()
				val, err := getFieldStringValue(fType, fVal)
				if err != nil {
					glog.Error(err)
					continue
				}
				path := rt.Elem().Field(i).Tag.Get("path")
				data.Field[path] = val
			}
		}
	} else {
		fType := rt.Elem()
		fVal := rv.Elem()
		val, err := getFieldStringValue(fType, fVal)
		if err != nil {
			glog.Error(err)
		}
		path := rt.Elem().Field(0).Tag.Get("path")
		data.Field[path] = val
	}
	return err
}

func createOrMergeToDb(gs interface{}, d *db.DB, ts *db.TableSpec, key db.Key) error {
	if util.IsValueNil(gs) {
		return nil
	}

	data := convertRequestBodyToInternal(gs)
	if !data.IsPopulated() {
		glog.V(2).Infof("createOrMergeToDb did nothing as the request body has no data")
		return nil
	}

	err := d.ModEntry(ts, key, data)
	if err != nil {
		glog.V(2).Infof("createOrMergeRecord failed as %v", err)
		return err
	}

	return err
}

func getFieldStringValue(t reflect.Type, v reflect.Value) (string, error) {
	var valStr string
	var err error

	bitSize := int(t.Size()) * 8
	switch t.Kind() {
	case reflect.Float32, reflect.Float64:
		valStr = strconv.FormatFloat(v.Float(), 'f', -1, bitSize)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		valStr = strconv.FormatInt(v.Int(), 10)
		if enum, isEnum := v.Interface().(ygot.GoEnum); isEnum {
			e := reflect.ValueOf(enum)
			lookup, ok := enum.ΛMap()[e.Type().Name()]
			if ok {
				valStr = lookup[v.Int()].Name
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		valStr = strconv.FormatUint(v.Uint(), 10)
	case reflect.String:
		valStr = v.String()
	case reflect.Bool:
		valStr = strconv.FormatBool(v.Bool())
	default:
		err = tlerr.New("failed to get the field value string as the type is %s", t.Name())
	}

	return valStr, err
}

func getGoEnumFieldString(field interface{}, val string) string {
	rt := reflect.TypeOf(field)
	rv := reflect.ValueOf(field)

	enumTypeName := rt.Name()

	// is GoEnum
	if enumVal, isEnum := rv.Interface().(ygot.GoEnum); isEnum {
		tmpInt, err := strconv.ParseInt(val, 10, 0)
		if err == nil {
			lookup, ok := enumVal.ΛMap()[enumTypeName]
			if !ok {
				glog.Errorf("cannot map enumerated value as type %s was unknown", enumTypeName)
				return ""
			}
			def, ok := lookup[tmpInt]
			if !ok {
				glog.Errorf("cannot map enumerated value as type %s has unknown value %d", enumTypeName, enumVal)
				return ""
			}
			return def.Name
		} else {
			glog.Error(rt.Field(0).Name, "has invalid value in database")
		}
	}

	return ""
}

func getRedisKey(prefix string, mdlKey interface{}) string {
	var tblKey string

	rt := reflect.TypeOf(mdlKey)
	rv := reflect.ValueOf(mdlKey)

	switch rt.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		tblKey = prefix + strconv.FormatInt(rv.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		tblKey = prefix + strconv.FormatUint(rv.Uint(), 10)
	case reflect.String:
		tblKey = prefix + rv.String()
	default:
		return ""
	}

	return tblKey
}

func getYangMdlKey(prefix string, tblKey string, keyType reflect.Type) (interface{}, error) {
	index := len(prefix)
	keyStr := tblKey[index:]
	keyVal, err := ytypes.StringToType(keyType, keyStr)
	if err != nil {
		return nil, err
	}

	return keyVal.Interface(), err
}

func internalRecordExist(tblMap interface{}, keys db.Key) bool {
	data := reflect.ValueOf(tblMap)
	if data.Kind() != reflect.Map || keys.Len() == 0 {
		return false
	}

	for _, key := range keys.Comp {
		if data.Kind() != reflect.Map {
			return false
		}
		data = data.MapIndex(reflect.ValueOf(key))
		if !data.IsValid() {
			return false
		}
	}

	return true
}

/* Check if targetUriPath is child (subtree) of nodePath
The return value can be used to decide if subtrees needs
to visited to fill the data or not.
*/
func isSubtreeRequest(targetUriPath string, nodePath string) bool {
	return strings.HasPrefix(targetUriPath, nodePath)
}

func isRootRequest(targetUriPath string, nodePath string) bool {
	return strings.Compare(targetUriPath, nodePath) == 0
}

/*
 * translib db api
 */
func getRedisData(dbCl *db.DB, ts *db.TableSpec, key db.Key) (db.Value, error) {
	var value db.Value

	if dbCl == nil || ts == nil || key.Len() == 0 {
		return value, errors.New("input invalid")
	}

	if !dbCl.KeyExists(ts, key) {
		return value, tlerr.NotFound("table %s with key %s does not exist in db %d", ts.Name, key.Comp, dbCl.Opts.DBNo)
	}

	entry, _ := dbCl.GetEntry(ts, key)
	glog.V(1).Infof("query db %d table %s by key %v", dbCl.Opts.DBNo, ts.Name, key.Comp)

	return entry, nil
}

func deleteRedisData(d *db.DB, ts *db.TableSpec, key db.Key) error {
	err := d.StartTx(nil, nil)
	if err != nil {
		return err
	}

	err = d.DeleteEntry(ts, key)
	if err != nil {
		return err
	}

	err = d.CommitTx()
	if err != nil {
		return err
	}

	return nil
}

func setRedisData(d *db.DB, ts *db.TableSpec, key db.Key, data db.Value) error {
	err := d.StartTx(nil, nil)
	if err != nil {
		return err
	}

	err = d.SetEntry(ts, key, data)
	if err != nil {
		return err
	}

	err = d.CommitTx()
	if err != nil {
		return err
	}

	return nil
}

func splitPrefix(value string, prefix string) string {
	return value[len(prefix):]
}

func getTableFieldStringValue(dbCl *db.DB, ts *db.TableSpec, key db.Key, fieldName string) (string, error) {
	data, err := getRedisData(dbCl, ts, key)
	if err == nil && len(data.Field) > 0 {
		if data.Has(fieldName) {
			return data.Get(fieldName), nil
		}
		return "", nil
	}

	return "", tlerr.New("get table %s field %s`s value from db %d failed", ts.Name, fieldName, dbCl.Opts.DBNo)
}

func isQuerySupervisoryChannel(path string, key string) bool {
	if strings.HasPrefix(key, "OSC") {
		if strings.Contains(path, "supervisory-channels") {
			return true
		}
	}
	return false
}

func constructRegexPathWithKey(mdb db.MDB, num db.DBNum, path string, params *regexPathKeyParams) []string {
	var pathWithKey []string

	for _, dbs := range mdb {
		dbCl := dbs[num]

		keys, _ := dbCl.GetKeys(asTableSpec(params.tableName))
		if len(keys) == 0 {
			continue
		}

		for _, key := range keys {
			if key.Len() != len(params.keyName) {
				continue
			}

			p := path
			for i, eachKey := range key.Comp {
				mdlKey, err := getYangMdlKey(params.redisPrefix[i], eachKey, reflect.TypeOf(""))
				if err != nil {
					glog.Error("construct path with key failed")
					return nil
				}
				mdlKeyStr, _ := mdlKey.(string)
				if isQuerySupervisoryChannel(path, mdlKeyStr) {
					mdlKeyStr = oscName2InterfaceName(mdlKeyStr)
				}
				oldStr := fmt.Sprintf("/%s/", params.listNodeName[i])
				newStr := fmt.Sprintf("/%s[%s=%s]/", params.listNodeName[i], params.keyName[i], mdlKeyStr)
				p = strings.ReplaceAll(p, oldStr, newStr)
			}
			pathWithKey = append(pathWithKey, p)
		}
	}

	return pathWithKey
}

func getFieldNameByNodeName(nodeName string) string {
	var suffix string
	elmts := strings.Split(nodeName, "-")
	for _, elmt := range elmts {
		//if len(suffix) == 0 {
		//	suffix = "_"
		//}
		tmp := []rune(elmt)
		if tmp[0] >= 97 && tmp[0] <= 122 { // ASCII a~z
			tmp[0] -= 32
		}
		elmt = string(tmp)
		suffix += elmt
	}
	return suffix
}

func constructRouteByPath(path string) string {
	var route string
	pathElmts := strings.Split(path[1:], "/")
	for i, elmt := range pathElmts {
		if index := strings.Index(elmt, ":"); index > 0 {
			ns := elmt[:index]
			node := elmt[index + 1 :]
			if i == 0 {
				route += getFieldNameByNodeName(ns)
			}
			route += getFieldNameByNodeName(node)
		} else if index = strings.Index(elmt, "["); index > 0  {
			route += getFieldNameByNodeName(elmt[:index])
		} else {
			route += getFieldNameByNodeName(elmt)
		}
	}

	return route
}

func parentNeedQuery(parentRoute string, targetRoute string) bool {
	// should query the list parent as the key node is necessary for the building of the target node.
	// should not query the container parent，which is to be supported
    // currently, ignore the type of parent and query parent always
	return strings.HasPrefix(targetRoute, parentRoute)
}

func childNeedQuery(childRoute string, targetRoute string) bool {
	return strings.HasPrefix(childRoute, targetRoute)
}

func siblingNeedQuery(siblingRoute string, targetRoute string) bool {
    return siblingRoute == targetRoute
}

// /openconfig-platform:components/component/openconfig-platform-transceiver:transceiver/physical-channels/channel/state
func needQuery(targetPath string, curNode interface{}) bool {
	targetRoute := constructRouteByPath(targetPath)
	if len(targetRoute) == 0 {
		return false
	}
	targetDepth := len(strings.Split(targetPath[1:], "/")) + 1 // plus the namespace element

	typeName := strings.TrimLeft(reflect.TypeOf(curNode).String(), "*ocbinds.")
	curNodeDepth := len(strings.Split(typeName, "_"))
	curNodeRoute := strings.ReplaceAll(typeName, "_", "")
	glog.V(2).Infof("target route depth %d value %s", targetDepth, targetRoute)
	glog.V(2).Infof("current route depth %d value %s", curNodeDepth, curNodeRoute)

	if targetDepth > curNodeDepth {
		return parentNeedQuery(curNodeRoute, targetRoute)
	}

	if targetDepth < curNodeDepth {
		return childNeedQuery(curNodeRoute, targetRoute)
	}

	return siblingNeedQuery(curNodeRoute, targetRoute)
}