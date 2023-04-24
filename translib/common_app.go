////////////////////////////////////////////////////////////////////////////////
//                                                                            //
//  Copyright 2019 Dell, Inc.                                                 //
//                                                                            //
//  Licensed under the Apache License, Version 2.0 (the "License");           //
//  you may not use this file except in compliance with the License.          //
//  You may obtain a copy of the License at                                   //
//                                                                            //
//  http://www.apache.org/licenses/LICENSE-2.0                                //
//                                                                            //
//  Unless required by applicable law or agreed to in writing, software       //
//  distributed under the License is distributed on an "AS IS" BASIS,         //
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.  //
//  See the License for the specific language governing permissions and       //
//  limitations under the License.                                            //
//                                                                            //
////////////////////////////////////////////////////////////////////////////////

package translib

import (
	"errors"
	"strings"
	"github.com/golang/glog"
	"github.com/openconfig/ygot/ygot"
	"github.com/openconfig/ygot/ytypes"
	"github.com/openconfig/ygot/util"
	"reflect"
	"github.com/Azure/sonic-mgmt-common/translib/db"
	"github.com/Azure/sonic-mgmt-common/translib/ocbinds"
	"github.com/Azure/sonic-mgmt-common/translib/tlerr"
	"github.com/Azure/sonic-mgmt-common/translib/transformer"
	"github.com/Azure/sonic-mgmt-common/translib/utils"
	"sync"
)

var ()

type CommonApp struct {
	pathInfo       *PathInfo
	body           []byte
	ygotRoot       *ygot.GoStruct
	ygotTarget     *interface{}
	skipOrdTableChk bool
	cmnAppTableMap map[int]map[db.DBNum]map[string]map[string]db.Value
	cmnAppYangDefValMap map[string]map[string]db.Value
	cmnAppYangAuxMap map[string]map[string]db.Value
	appOptions
}

func (app *CommonApp) translateMDBReplace(numDB db.NumberDB) ([]db.WatchKeys, error) {
	return nil, tlerr.NotSupported("Unsupported")
}

func (app *CommonApp) translateMDBGet(mdb db.MDB) error {
	return tlerr.NotSupported("Unsupported")
}

func (app *CommonApp) translateGetRegex(mdb db.MDB) error {
	return tlerr.NotSupported("Unsupported")
}

func (app *CommonApp) processMDBReplace(numDB db.NumberDB) (SetResponse, error) {
	return SetResponse{}, tlerr.NotSupported("Unsupported")
}

func (app *CommonApp) processMDBGet(mdb db.MDB) (GetResponse, error) {
	return GetResponse{}, tlerr.NotSupported("Unsupported")
}

func (app *CommonApp) processGetRegex(mdb db.MDB) ([]GetResponseRegex, error) {
	return nil, tlerr.NotSupported("Unsupported")
}

var cmnAppInfo = appInfo{appType: reflect.TypeOf(CommonApp{}),
	ygotRootType:  nil,
	isNative:      false,
	tablesToWatch: nil}

func init() {

	register_model_path := []string{"/sonic-", "*"} // register yang model path(s) to be supported via common app
	for _, mdl_pth := range register_model_path {
		err := register(mdl_pth, &cmnAppInfo)

		if err != nil {
			glog.Fatal("Register Common app module with App Interface failed with error=", err, "for path=", mdl_pth)
		}
	}
	mdlCpblt := transformer.AddModelCpbltInfo()
	if mdlCpblt == nil {
		glog.Warning("Failure in fetching model capabilities data.")
	} else {
		for yngMdlNm, mdlDt := range(mdlCpblt) {
			err := addModel(&ModelData{Name: yngMdlNm, Org: mdlDt.Org, Ver: mdlDt.Ver})
			if err != nil {
				glog.Warningf("Adding model data for module %v to appinterface failed with error=%v", yngMdlNm, err)
			}
		}
	}
}

func (app *CommonApp) initialize(data appData) {
	pathInfo := NewPathInfo(data.path)
	*app = CommonApp{pathInfo: pathInfo, body: data.payload, ygotRoot: data.ygotRoot, ygotTarget: data.ygotTarget, skipOrdTableChk: false}
	app.appOptions = data.appOptions

}

func (app *CommonApp) translateCreate(d *db.DB) ([]db.WatchKeys, error) {
	var err error
	var keys []db.WatchKeys
	glog.Info("translateCreate:path =", app.pathInfo.Path)

	keys, err = app.translateCRUDCommon(d, CREATE)

	return keys, err
}

func (app *CommonApp) translateUpdate(d *db.DB) ([]db.WatchKeys, error) {
	var err error
	var keys []db.WatchKeys
	glog.Info("translateUpdate:path =", app.pathInfo.Path)

	keys, err = app.translateCRUDCommon(d, UPDATE)

	return keys, err
}

func (app *CommonApp) translateReplace(d *db.DB) ([]db.WatchKeys, error) {
	var err error
	var keys []db.WatchKeys
	glog.Info("translateReplace:path =", app.pathInfo.Path)

	keys, err = app.translateCRUDCommon(d, REPLACE)

	return keys, err
}

func (app *CommonApp) translateDelete(d *db.DB) ([]db.WatchKeys, error) {
	var err error
	var keys []db.WatchKeys
	glog.Info("translateDelete:path =", app.pathInfo.Path)
	keys, err = app.translateCRUDCommon(d, DELETE)

	return keys, err
}

func (app *CommonApp) translateGet(dbs [db.MaxDB]*db.DB) error {
	var err error
	glog.Info("translateGet:path =", app.pathInfo.Path)
	return err
}

func (app *CommonApp) translateSubscribe(dbs [db.MaxDB]*db.DB, path string) (*notificationOpts, *notificationInfo, error) {
    var err error
    var subscDt transformer.XfmrTranslateSubscribeInfo
    var notifInfo notificationInfo
    var notifOpts notificationOpts
    txCache := new(sync.Map)
    err = tlerr.NotSupportedError{Format: "Subscribe not supported", Path: path}

    glog.Info("tranlateSubscribe:path", path)
    subscDt, err = transformer.XlateTranslateSubscribe(path, dbs, txCache)
    if subscDt.PType == transformer.OnChange {
        notifOpts.pType = OnChange
    } else {
        notifOpts.pType = Sample
    }
    notifOpts.mInterval = subscDt.MinInterval
    notifOpts.isOnChangeSupported = subscDt.OnChange
    if err != nil {
        glog.Infof("returning: notificationOpts - %v, nil, error - %v", notifOpts, err)
        return &notifOpts, nil, err
    }
    if subscDt.DbDataMap == nil {
        glog.Infof("DB data is nil so returning: notificationOpts - %v, nil, error - %v", notifOpts, err)
        return &notifOpts, nil, err
    } else {
        for dbNo, dbDt := range(subscDt.DbDataMap) {
            if (len(dbDt) == 0) { //ideally all tables for a given uri should be from same DB
                continue
            }
            glog.Infof("Adding to notifInfo, Db Data - %v for DB No - %v", dbDt, dbNo)
            notifInfo.dbno = dbNo
            // in future there will be, multi-table in a DB, support from translib, for now its just single table
            for tblNm, tblDt := range(dbDt) {
                notifInfo.table = db.TableSpec{Name:tblNm}
                if (len(tblDt) == 1) {
                    for tblKy := range(tblDt) {
                        notifInfo.key = asKey(tblKy)
                        notifInfo.needCache = subscDt.NeedCache
                    }
                } else {
                    if (len(tblDt) >  1) {
                        glog.Warningf("More than one DB key found for subscription path - %v", path)
                    } else {
                        glog.Warningf("No DB key found for subscription path - %v", path)
                    }
                    return &notifOpts, nil, err
                }

            }
        }
    }
    glog.Infof("For path - %v, returning: notifOpts - %v, notifInfo - %v, error - nil", path, notifOpts, notifInfo)
    return &notifOpts, &notifInfo, nil
}

func (app *CommonApp) translateAction(mdb db.MDB) error {
    var err error
    glog.Info("translateAction:path =", app.pathInfo.Path, app.body)
    return err
}

func (app *CommonApp) processCreate(d *db.DB) (SetResponse, error) {
	var err error
	var resp SetResponse

	glog.Info("processCreate:path =", app.pathInfo.Path)
	targetType := reflect.TypeOf(*app.ygotTarget)
	glog.Infof("processCreate: Target object is a <%s> of Type: %s", targetType.Kind().String(), targetType.Elem().Name())
	if err = app.processCommon(d, CREATE); err != nil {
		glog.Warning(err)
		resp = SetResponse{ErrSrc: AppErr}
	}

	return resp, err
}

func (app *CommonApp) processUpdate(d *db.DB) (SetResponse, error) {
	var err error
	var resp SetResponse
	glog.Info("processUpdate:path =", app.pathInfo.Path)
	if err = app.processCommon(d, UPDATE); err != nil {
		glog.Warning(err)
		resp = SetResponse{ErrSrc: AppErr}
	}

	return resp, err
}

func (app *CommonApp) processReplace(d *db.DB) (SetResponse, error) {
	var err error
	var resp SetResponse
	glog.Info("processReplace:path =", app.pathInfo.Path)
	if err = app.processCommon(d, REPLACE); err != nil {
		glog.Warning(err)
		resp = SetResponse{ErrSrc: AppErr}
	}
	return resp, err
}

func (app *CommonApp) processDelete(d *db.DB) (SetResponse, error) {
	var err error
	var resp SetResponse

	glog.Infof("processDelete:path = %s, deleteEmptyEntry = %v", app.pathInfo.Path, app.deleteEmptyEntry)

	if err = app.processCommon(d, DELETE); err != nil {
		glog.Warning(err)
		resp = SetResponse{ErrSrc: AppErr}
	}

	return resp, err
}

func (app *CommonApp) processGet(dbs [db.MaxDB]*db.DB) (GetResponse, error) {
    var err error
    var payload []byte
    var resPayload []byte
    glog.Info("processGet:path =", app.pathInfo.Path)
    txCache := new(sync.Map)

    for {
	    origXfmrYgotRoot, _ := ygot.DeepCopy((*app.ygotRoot).(ygot.GoStruct))

        isEmptyPayload  := false
		appYgotStruct := (*app.ygotRoot).(ygot.GoStruct)        
	    payload, isEmptyPayload, err = transformer.GetAndXlateFromDB(app.pathInfo.Path, &appYgotStruct, dbs, txCache)
	    if err != nil {
		    glog.Warning("transformer.GetAndXlateFromDB() returned : ", err)
		    resPayload = payload
		    break
            }
	    if strings.HasPrefix(app.pathInfo.Path, "/sonic") && isEmptyPayload {
		    glog.Info("transformer.GetAndXlateFromDB() returned EmptyPayload")
		    resPayload = payload
		    break
	    }

	    targetObj, tgtObjCastOk := (*app.ygotTarget).(ygot.GoStruct)
	    if !tgtObjCastOk {
		    /*For ygotTarget populated by tranlib, for query on leaf level and list(without instance) level, 
		      casting to GoStruct fails so use the parent node of ygotTarget to Unmarshall the payload into*/
		    glog.Infof("Use GetParentNode() instead of casting ygotTarget to GoStruct, uri - %v", app.pathInfo.Path)
		    targetUri := app.pathInfo.Path
		    parentTargetObj, _, getParentNodeErr := getParentNode(&targetUri, (*app.ygotRoot).(*ocbinds.Device))
		    if getParentNodeErr != nil {
			    glog.Warningf("getParentNode() failure for uri %v", app.pathInfo.Path)
			    resPayload = payload
			    break
		    }
		    if parentTargetObj != nil {
			    targetObj, tgtObjCastOk = (*parentTargetObj).(ygot.GoStruct)
			    if !tgtObjCastOk {
				    glog.Warningf("Casting of parent object returned from getParentNode() to GoStruct failed(uri - %v)", app.pathInfo.Path)
				    resPayload = payload
				    break
			    }
		    } else {
			    glog.Warningf("getParentNode() returned a nil Object for uri %v", app.pathInfo.Path)
                            resPayload = payload
                            break
		    }
	    }
	    if targetObj != nil {
		    updateListEntriesOpt := ytypes.AllowUpdateInListMap{}
		    err = ocbinds.Unmarshal(payload, targetObj, &updateListEntriesOpt)
		    if err != nil {
			    glog.Warning("ocbinds.Unmarshal()  returned : ", err)
			    resPayload = payload
			    break
		    }

		    resYgot := (*app.ygotRoot)
		    if !strings.HasPrefix(app.pathInfo.Path, "/sonic") {
			    if isEmptyPayload {
				    if areEqual(appYgotStruct, origXfmrYgotRoot) {
					    glog.Info("origXfmrYgotRoot and appYgotStruct are equal.")
					    // No data available in appYgotStruct.
					    if transformer.IsLeafNode(app.pathInfo.Path) {
						    //if leaf not exist in DB subtree won't fill ygotRoot, as per RFC return err
						    resPayload = payload
						    glog.Info("No data found for leaf.")
						    err = tlerr.NotFound("Resource not found")
						    break
					    }
					    resPayload = payload
					    glog.Info("No data available")
					    //TODO: Return not found error
					    //err = tlerr.NotFound("Resource not found")
					    break

				    }
				    resYgot = appYgotStruct
			    }
		    }
		    if resYgot != nil {
			    resPayload, err = generateGetResponsePayload(app.pathInfo.Path, resYgot.(*ocbinds.Device), app.ygotTarget)
			    if err != nil {
				    glog.Warning("generateGetResponsePayload() couldn't generate payload.")
				    resPayload = payload
			    }
		    } else {
			resPayload = payload
		    }

		    break
	    } else {
		glog.Warning("processGet. targetObj is null. Unable to Unmarshal payload")
		resPayload = payload
		break
	    }
    }

    return GetResponse{Payload: resPayload}, err
}

func (app *CommonApp) processAction(mdb db.MDB) (ActionResponse, error) {
        var resp ActionResponse
        var err error

	//resp.Payload, err = transformer.CallRpcMethod(app.pathInfo.Path, app.body, dbs)
	glog.Info("transformer.CallRpcMethod() returned")

	return resp, err
}

func (app *CommonApp) translateCRUDCommon(d *db.DB, opcode int) ([]db.WatchKeys, error) {
	var err error
	var keys []db.WatchKeys
	var tblsToWatch []*db.TableSpec
	txCache := new(sync.Map)
	glog.Info("translateCRUDCommon:path =", app.pathInfo.Path)

	// translate yang to db
	result, defValMap, auxMap, err := transformer.XlateToDb(app.pathInfo.Path, opcode, d, (*app).ygotRoot, (*app).ygotTarget, (*app).body, txCache, &app.skipOrdTableChk)
	glog.Info("transformer.XlateToDb() returned result DB map - ", result, "\nDefault value Db Map - ", defValMap, "\nAux Db Map - ", auxMap)


	if err != nil {
		glog.Warning(err)
		return keys, err
	}
	app.cmnAppTableMap = result
	app.cmnAppYangDefValMap = defValMap
	app.cmnAppYangAuxMap = auxMap //used for Replace case
	if len(result) == 0 {
		glog.Info("XlatetoDB() returned empty map")
		//Note: Get around for no redis ABNF Schema for set(temporary)
		//`err = errors.New("transformer.XlatetoDB() returned empty map")
		return keys, err
	}

	moduleNm, err := transformer.GetModuleNmFromPath(app.pathInfo.Path)
        if (err != nil) || (len(moduleNm) == 0) {
                glog.Warning("GetModuleNmFromPath() couldn't fetch module name.")
                return keys, err
        }

	var resultTblList []string
        for _, dbMap := range result { //Get dependency list for all tables in result
		for _, resMap := range dbMap { //Get dependency list for all tables in result
		        for tblnm := range resMap { //Get dependency list for all tables in result
				resultTblList = append(resultTblList, tblnm)
			}
		}
	}
        glog.Info("Result Tables List", resultTblList)

	// Get list of tables to watch
	if len(resultTblList) > 0 {
		depTbls := transformer.GetTablesToWatch(resultTblList, moduleNm)
		if len(depTbls) == 0 {
			glog.Warningf("Couldn't get Tables to watch for module %v", moduleNm)
			err = errors.New("GetTablesToWatch returned empty slice")
			return keys, err
		}
		for _, tbl := range depTbls {
			tblsToWatch = append(tblsToWatch, &db.TableSpec{Name: tbl})
		}
	}
        glog.Info("Tables to watch", tblsToWatch)
        cmnAppInfo.tablesToWatch = tblsToWatch

	keys, err = app.generateDbWatchKeys(d, false)
	return keys, err
}

func (app *CommonApp) processCommon(d *db.DB, opcode int) error {

	var err error
	if len(app.cmnAppTableMap) == 0 {
		return err
	}

	glog.Info("Processing DB operation for ", app.cmnAppTableMap)
	switch opcode {
		case CREATE:
			glog.Info("CREATE case")
		case UPDATE:
			glog.Info("UPDATE case")
		case REPLACE:
			glog.Info("REPLACE case")
		case DELETE:
			glog.Info("DELETE case")
	}

	// Handle delete first if any available
	if _, ok := app.cmnAppTableMap[DELETE][db.ConfigDB]; ok {
		err = app.cmnAppDelDbOpn(d, DELETE, app.cmnAppTableMap[DELETE][db.ConfigDB])
		if err != nil {
			glog.Info("Process delete fail. cmnAppDelDbOpn error:", err)
			return err
		}
	}
	// Handle create operation next
	if _, ok := app.cmnAppTableMap[CREATE][db.ConfigDB]; ok {
		err = app.cmnAppCRUCommonDbOpn(d, CREATE, app.cmnAppTableMap[CREATE][db.ConfigDB])
		if err != nil {
			glog.Info("Process create fail. cmnAppCRUCommonDbOpn error:", err)
			return err
		}
	}
	// Handle update and replace operation next
	if _, ok := app.cmnAppTableMap[UPDATE][db.ConfigDB]; ok {
		err = app.cmnAppCRUCommonDbOpn(d, UPDATE, app.cmnAppTableMap[UPDATE][db.ConfigDB])
		if err != nil {
			glog.Info("Process update fail. cmnAppCRUCommonDbOpn error:", err)
			return err
		}
	}
	if _, ok := app.cmnAppTableMap[REPLACE][db.ConfigDB]; ok {
		err = app.cmnAppCRUCommonDbOpn(d, REPLACE, app.cmnAppTableMap[REPLACE][db.ConfigDB])
		if err != nil {
			glog.Info("Process replace fail. cmnAppCRUCommonDbOpn error:", err)
			return err
		}
	}
	glog.Info("Returning from processCommon() - success")
	return err
}

func (app *CommonApp) cmnAppCRUCommonDbOpn(d *db.DB, opcode int, dbMap map[string]map[string]db.Value) error {
	var err error
	var cmnAppTs *db.TableSpec
	var xfmrTblLst []string
	var resultTblLst []string

	for tblNm := range(dbMap) {
		xfmrTblLst = append(xfmrTblLst, tblNm)
	}
	resultTblLst, err = utils.SortAsPerTblDeps(xfmrTblLst)
	if err != nil {
		return err
	}

	/* CVL sorted order is in child first, parent later order. CRU ops from parent first order */
	for idx := len(resultTblLst)-1; idx >= 0; idx-- {
		tblNm := resultTblLst[idx]
		glog.Info("In Yang to DB map returned from transformer looking for table = ", tblNm)
		if tblVal, ok := dbMap[tblNm]; ok {
			cmnAppTs = &db.TableSpec{Name: tblNm}
			glog.Info("Found table entry in yang to DB map")
			if ((tblVal == nil) || (len(tblVal) == 0)) {
				glog.Info("No table instances/rows found.")
				continue
			}
			for tblKey, tblRw := range tblVal {
				glog.Info("Processing Table key ", tblKey)
				// REDIS doesn't allow to create a table instance without any fields
				if tblRw.Field == nil {
					tblRw.Field = map[string]string{"NULL": "NULL"}
				}
				if len(tblRw.Field) == 0 {
					tblRw.Field["NULL"] = "NULL"
				}
				if len(tblRw.Field) > 1 {
					delete(tblRw.Field, "NULL")
				}
				existingEntry, _ := d.GetEntry(cmnAppTs, db.Key{Comp: []string{tblKey}})
				switch opcode {
				case CREATE:
					if existingEntry.IsPopulated() {
						glog.Info("Create case - Entry ", tblKey, " already exists hence modifying it.")
						/* Handle leaf-list merge if any leaf-list exists 
						A leaf-list field in redis has "@" suffix as per swsssdk convention.
						*/
						resTblRw := db.Value{Field: map[string]string{}}
						resTblRw = checkAndProcessLeafList(existingEntry, tblRw, UPDATE, d, tblNm, tblKey)
						glog.Info("Processing Table row ", resTblRw)
						err = d.ModEntry(cmnAppTs, db.Key{Comp: []string{tblKey}}, resTblRw)
						if err != nil {
							glog.Warning("CREATE case - d.ModEntry() failure")
							return err
						}
					} else {
						if tblRwDefaults, defaultOk := app.cmnAppYangDefValMap[tblNm][tblKey]; defaultOk {
							glog.Info("Entry ", tblKey, " doesn't exist so fill defaults - ", tblRwDefaults)
							for fld, val := range tblRwDefaults.Field {
								tblRw.Field[fld] = val
							}
						}
						glog.Info("Processing Table row ", tblRw)
						err = d.CreateEntry(cmnAppTs, db.Key{Comp: []string{tblKey}}, tblRw)
						if err != nil {
							glog.Warning("CREATE case - d.CreateEntry() failure")
							return err
						}
					}
				case UPDATE:
					if existingEntry.IsPopulated() {
						glog.Info("Entry already exists hence modifying it.")
						/* Handle leaf-list merge if any leaf-list exists 
						A leaf-list field in redis has "@" suffix as per swsssdk convention.
						*/
						resTblRw := db.Value{Field: map[string]string{}}
						resTblRw = checkAndProcessLeafList(existingEntry, tblRw, UPDATE, d, tblNm, tblKey)
						err = d.ModEntry(cmnAppTs, db.Key{Comp: []string{tblKey}}, resTblRw)
						if err != nil {
							glog.Warning("UPDATE case - d.ModEntry() failure")
							return err
						}
					} else {
						// workaround to patch operation from CLI
						glog.Info("Create(pathc) an entry.")
						if tblRwDefaults, defaultOk := app.cmnAppYangDefValMap[tblNm][tblKey]; defaultOk {
							glog.Info("Entry ", tblKey, " doesn't exist so fill defaults - ", tblRwDefaults)
                                                        for fld, val := range tblRwDefaults.Field {
                                                                tblRw.Field[fld] = val
                                                        }
                                                }
						glog.Info("Processing Table row ", tblRw)
						err = d.CreateEntry(cmnAppTs, db.Key{Comp: []string{tblKey}}, tblRw)
						if err != nil {
							glog.Warning("UPDATE case - d.CreateEntry() failure")
							return err
						}
					}
				case REPLACE:
					if tblRwDefaults, defaultOk := app.cmnAppYangDefValMap[tblNm][tblKey]; defaultOk {
						glog.Info("For entry ", tblKey, ", being replaced, fill defaults - ", tblRwDefaults)
						for fld, val := range tblRwDefaults.Field {
							tblRw.Field[fld] = val
						}
					}
					glog.Info("Processing Table row ", tblRw)
					if existingEntry.IsPopulated() {
						glog.Info("Entry already exists.")
						auxRwOk := false
						auxRw := db.Value{Field: map[string]string{}}
						auxRw, auxRwOk = app.cmnAppYangAuxMap[tblNm][tblKey]
						glog.Info("Process Aux row ", auxRw)
						isTlNd := false
						if !strings.HasPrefix(app.pathInfo.Path, "/sonic") {
							isTlNd, err = transformer.IsTerminalNode(app.pathInfo.Path)
							glog.Info("transformer.IsTerminalNode() returned - ", isTlNd, " error ", err)
							if err != nil {
								return err
							}
						}
						if isTlNd && isPartialReplace(existingEntry, tblRw, auxRw) {
							glog.Info("Since its partial replace modifying fields - ", tblRw)
							err = d.ModEntry(cmnAppTs, db.Key{Comp: []string{tblKey}}, tblRw)
							if err != nil {
								glog.Warning("REPLACE case - d.ModEntry() failure")
								return err
							}
							if auxRwOk {
								if len(auxRw.Field) > 0 {
									glog.Info("Since its partial replace delete aux fields - ", auxRw)
									err := d.DeleteEntryFields(cmnAppTs, db.Key{Comp: []string{tblKey}}, auxRw)
									if err != nil {
										glog.Warning("REPLACE case - d.DeleteEntryFields() failure")
										return err
									}
								}
							}
						} else {
							err := d.SetEntry(cmnAppTs, db.Key{Comp: []string{tblKey}}, tblRw)
							if err != nil {
								glog.Warning("REPLACE case - d.SetEntry() failure")
								return err
							}
						}
					} else {
						glog.Info("Entry doesn't exist hence create it.")
						err = d.CreateEntry(cmnAppTs, db.Key{Comp: []string{tblKey}}, tblRw)
						if err != nil {
							glog.Warning("REPLACE case - d.CreateEntry() failure")
							return err
						}
					}
				}
			}
		}
	}
	return err
}

func (app *CommonApp) cmnAppDelDbOpn(d *db.DB, opcode int, dbMap map[string]map[string]db.Value) error {
	var err error
	var cmnAppTs, dbTblSpec *db.TableSpec
	var moduleNm string
	var xfmrTblLst []string
	var resultTblLst []string
	var ordTblList []string

	for tblNm := range(dbMap) {
		xfmrTblLst = append(xfmrTblLst, tblNm)
	}
	resultTblLst, err = utils.SortAsPerTblDeps(xfmrTblLst)
	if err != nil {
		return err
	}


	/* Retrieve module Name */
	moduleNm, err = transformer.GetModuleNmFromPath(app.pathInfo.Path)
	if (err != nil) || (len(moduleNm) == 0) {
		glog.Warning("GetModuleNmFromPath() failed")
		return err
	}
	glog.Info("getModuleNmFromPath() returned module name = ", moduleNm)

	/* resultTblLst has child first, parent later order */
	for _, tblNm := range resultTblLst {
		glog.Info("In Yang to DB map returned from transformer looking for table = ", tblNm)
		if tblVal, ok := dbMap[tblNm]; ok {
			cmnAppTs = &db.TableSpec{Name: tblNm}
			glog.Info("Found table entry in yang to DB map")
			if !app.skipOrdTableChk {
				ordTblList = transformer.GetXfmrOrdTblList(tblNm)
				if len(ordTblList) == 0 {
					ordTblList = transformer.GetOrdTblList(tblNm, moduleNm)
				}
				if len(ordTblList) == 0 {
					glog.Warning("GetOrdTblList returned empty slice")
					err = errors.New("GetOrdTblList returned empty slice. Insufficient information to process request")
					return err
				}
				glog.Infof("GetOrdTblList for table - %v, module %v returns %v", tblNm, moduleNm, ordTblList)
			}
			if len(tblVal) == 0 {
				glog.Info("DELETE case - No table instances/rows found hence delete entire table = ", tblNm)
				if !app.skipOrdTableChk {
					for _, ordtbl := range ordTblList {
						if ordtbl == tblNm {
							// Handle the child tables only till you reach the parent table entry
							break
						}
						glog.Info("Since parent table is to be deleted, first deleting child table = ", ordtbl)
						dbTblSpec = &db.TableSpec{Name: ordtbl}
						err = d.DeleteTable(dbTblSpec)
						if err != nil {
							glog.Warning("DELETE case - d.DeleteTable() failure for Table = ", ordtbl)
							return err
						}
					}
				}
				err = d.DeleteTable(cmnAppTs)
				if err != nil {
					glog.Warning("DELETE case - d.DeleteTable() failure for Table = ", tblNm)
					return err
				}
				glog.Info("DELETE case - Deleted entire table = ", tblNm)
				// Continue to repeat ordered deletion for all tables
				continue

			}

			for tblKey, tblRw := range tblVal {
				if len(tblRw.Field) == 0 {
					glog.Info("DELETE case - no fields/cols to delete hence delete the entire row.")
					glog.Info("First, delete child table instances that correspond to parent table instance to be deleted = ", tblKey)
					if !app.skipOrdTableChk {
						for _, ordtbl := range ordTblList {
							if ordtbl == tblNm {
								// Handle the child tables only till you reach the parent table entry
								break;
							}
							dbTblSpec = &db.TableSpec{Name: ordtbl}
							keyPattern := tblKey + "|*"
							glog.Info("Key pattern to be matched for deletion = ", keyPattern)
							err = d.DeleteKeys(dbTblSpec, db.Key{Comp: []string{keyPattern}})
							if err != nil {
								glog.Warning("DELETE case - d.DeleteTable() failure for Table = ", ordtbl)
								return err
							}
							glog.Info("Deleted keys matching parent table key pattern for child table = ", ordtbl)
						}
					}
					err = d.DeleteEntry(cmnAppTs, db.Key{Comp: []string{tblKey}})
					if err != nil {
						glog.Warning("DELETE case - d.DeleteEntry() failure")
						return err
					}
					glog.Info("Finally deleted the parent table row with key = ", tblKey)
				} else {
					glog.Info("DELETE case - fields/cols to delete hence delete only those fields.")
					existingEntry, exstErr := d.GetEntry(cmnAppTs, db.Key{Comp: []string{tblKey}})
					if exstErr != nil {
						glog.Info("Table Entry from which the fields are to be deleted does not exist")
						err = exstErr
						return err
					}
					/* handle leaf-list merge if any leaf-list exists */
					resTblRw := checkAndProcessLeafList(existingEntry, tblRw, DELETE, d, tblNm, tblKey)
					glog.Info("DELETE case - checkAndProcessLeafList() returned table row ", resTblRw)
					if len(resTblRw.Field) > 0 {
						if !app.deleteEmptyEntry {
							/* add the NULL field if the last field gets deleted && deleteEmpyEntry is false */
							deleteCount := 0
							for field := range existingEntry.Field {
								if resTblRw.Has(field) {
									deleteCount++
								}
							}
							if deleteCount == len(existingEntry.Field) {
								nullTblRw := db.Value{Field: map[string]string{"NULL": "NULL"}}
								glog.Info("Last field gets deleted, add NULL field to keep an db entry")
								err = d.ModEntry(cmnAppTs, db.Key{Comp: []string{tblKey}}, nullTblRw)
								if err != nil {
									glog.Warning("UPDATE case - d.ModEntry() failure")
									return err
								}
							}
						}
						/* deleted fields */
						err := d.DeleteEntryFields(cmnAppTs, db.Key{Comp: []string{tblKey}}, resTblRw)
						if err != nil {
							glog.Warning("DELETE case - d.DeleteEntryFields() failure")
							return err
						}
					}
				}
			}
		}
	} /* end of ordered table list for loop */
	return err
}

func (app *CommonApp) generateDbWatchKeys(d *db.DB, isDeleteOp bool) ([]db.WatchKeys, error) {
	var err error
	var keys []db.WatchKeys

	return keys, err
}

/*check if any field is leaf-list , if yes perform merge*/
func checkAndProcessLeafList(existingEntry db.Value, tblRw db.Value, opcode int, d *db.DB, tblNm string, tblKey string) db.Value {
	dbTblSpec := &db.TableSpec{Name: tblNm}
	mergeTblRw := db.Value{Field: map[string]string{}}
	for field, value := range tblRw.Field {
		if strings.HasSuffix(field, "@") {
			exstLst := existingEntry.GetList(field)
			glog.Infof("Existing DB value for field %v - %v", field, exstLst)
			var valueLst []string
			if value != "" { //zero len string as leaf-list value is treated as delete entire leaf-list
				valueLst = strings.Split(value, ",")
			}
			glog.Infof("Incoming value for field %v - %v", field, valueLst)
			if len(exstLst) != 0 {
				glog.Infof("Existing list is not empty for field %v", field)
				for _, item := range valueLst {
					if !contains(exstLst, item) {
						if opcode == UPDATE {
							exstLst = append(exstLst, item)
						}
					} else {
						if opcode == DELETE {
                                                        exstLst = utils.RemoveElement(exstLst, item)
                                                }

					}
				}
				glog.Infof("For field %v value after merging incoming with existing %v", field, exstLst)
				if opcode == DELETE {
					if len(valueLst) > 0 {
						mergeTblRw.SetList(field, exstLst)
						if len(exstLst) == 0 {
							tblRw.Field[field] = ""
						} else {
							delete(tblRw.Field, field)
						}
					}
				} else if opcode == UPDATE {
					tblRw.SetList(field, exstLst)
				}
			} else { //when existing list is empty(either empty string val in field or no field at all n entry)
				glog.Infof("Existing list is empty for field %v", field)
				if opcode == UPDATE {
					if len(valueLst) > 0 {
						exstLst = valueLst
						tblRw.SetList(field, exstLst)
					} else {
						tblRw.Field[field] = ""
					}
				} else if opcode == DELETE {
					_, fldExistsOk := existingEntry.Field[field]
					if (fldExistsOk && (len(valueLst) == 0)) {
						tblRw.Field[field] = ""
					} else {
						delete(tblRw.Field, field)
					}
				}
                        }
		}
	}
	/* delete specific item from leaf-list */
	if opcode == DELETE {
		if len(mergeTblRw.Field) == 0 {
			glog.Infof("mergeTblRow is empty - Returning Table Row %v", tblRw)
			return tblRw
		}
		err := d.ModEntry(dbTblSpec, db.Key{Comp: []string{tblKey}}, mergeTblRw)
		if err != nil {
			glog.Warning("DELETE case(merge leaf-list) - d.ModEntry() failure")
		}
	}
	glog.Infof("Returning Table Row %v", tblRw)
	return tblRw
}

// This function is a copy of the function areEqual in ygot.util package.
// areEqual compares a and b. If a and b are both pointers, it compares the
// values they are pointing to.
func areEqual(a, b interface{}) bool {
        if util.IsValueNil(a) && util.IsValueNil(b) {
                return true
        }
        va, vb := reflect.ValueOf(a), reflect.ValueOf(b)
        if va.Kind() == reflect.Ptr && vb.Kind() == reflect.Ptr {
                return reflect.DeepEqual(va.Elem().Interface(), vb.Elem().Interface())
        }

        return reflect.DeepEqual(a, b)
}

func isPartialReplace(exstRw db.Value, replTblRw db.Value, auxRw db.Value) bool {
	/* if existing entry contains field thats not present in result,
           default and auxillary map then its a partial replace
         */
	partialReplace := false
	for exstFld := range exstRw.Field {
		if exstFld == "NULL" {
			continue
		}
		isIncomingFld := false
		if replTblRw.Has(exstFld) {
			continue
		}
		if auxRw.Has(exstFld) {
			continue
		}
		if !isIncomingFld {
			glog.Info("Entry contains field ", exstFld, " not found in result, default and aux fields hence its partial replace.")
			partialReplace = true
			break
		}
	}
	glog.Info("returning partialReplace - ", partialReplace)
	return partialReplace
}



