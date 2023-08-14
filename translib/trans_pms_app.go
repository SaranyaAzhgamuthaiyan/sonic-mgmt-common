package translib

import (
	"encoding/json"
	"github.com/Azure/sonic-mgmt-common/translib/db"
	"github.com/Azure/sonic-mgmt-common/translib/tlerr"
	"github.com/golang/glog"
	"github.com/openconfig/ygot/ygot"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type PmpInfo struct {
	PmParameter string
	NodeName string
	IsGuage bool
	OnlyInstant bool
}

type PmDbInfo struct {
	dbName string
	dbNum db.DBNum
	ts *db.TableSpec
	key db.Key
}

type TransPmsApp struct {
	path       *PathInfo
	reqData    []byte
	ygotRoot   *ygot.GoStruct
	ygotTarget *interface{}
}

type GetPmInput struct {
	PmBinType *string `json:"pm-bin-type"`
	EntityName *string `json:"entity-name"`
	PmTimePeriod *string `json:"pm-time-period"`
	PmpType *string `json:"pmp-type"`
	PmParameter *string `json:"pm-parameter"`
	Validity *string `json:"validity"`
	StartNumberOfBin *uint16 `json:"start-number-of-bin"`
	EndNumberOfBin *uint16 `json:"end-number-of-bin"`
}

type PmResults struct {
	PmResult []*PmResult `json:"pm-result"`
}

type PmResult struct {
	PmBinNo uint16 `json:"pm-bin-no"`
	PmBinState *PmBinState `json:"pm-bin-state"`
	PmInstances *PmInstances `json:"pm-instances,omitempty"`
}

type PmBinState struct {
	CollectStartTime string `json:"collect-start-time,omitempty"`
	CollectEndTime string `json:"collect-end-time,omitempty"`
	BinTimePeriod string `json:"bin-time-period"`
}

type PmInstances struct {
	PmInstance []*PmInstance `json:"pm-instance"`
}

type PmInstance struct {
	EntityName string `json:"entity-name"`
	PmItems *PmItems `json:"pm-items"`
}

type PmItems struct {
	PmItem []*PmItem `json:"pm-item"`
}

type PmItem struct {
	PmTimePeriod string `json:"pm-time-period"`
	PmpType string `json:"pmp-type"`
	PmParameter string `json:"pm-parameter"`
	MonitoringDateTime string `json:"monitoring-date-time,omitempty"`
	PmValue string `json:"pm-value"`
	Validity string `json:"validity"`
}

type PmQueryCondition struct {
	binNo uint16
	binType string
	entityName string
	timePeriod string
	tableLastKey string
}

var (
	entityPmpTypeMap     map[PMEntityType][]string
	pmpTypePathMap       map[string][]string
	pathPmpMap           map[string][]PmpInfo
	logicalChannelMap    map[string]*logicalChannelInfo
	componentList        []string
)

const (
	XPATH_TEMP_CHASSIS        = "/openconfig-platform/components/component/state/1"
	XPATH_TEMP_EQUIPMENT      = "/openconfig-platform/components/component/state/2"
	XPATH_FAN                 = "/openconfig-platform/components/component/fan/state"
	XPATH_PSU                 = "/openconfig-platform/components/component/power-supply/state"
	XPATH_OCH_FIRST           = "/openconfig-platform/components/component/optical-channel/state/1"
	XPATH_OCH_SECOND          = "/openconfig-platform/components/component/optical-channel/state/2"
	XPATH_LOGICAL_CHANNEL_OTN = "/openconfig-terminal-device/terminal-device/logical-channels/channel/otn/state"
	XPATH_LOGICAL_CHANNEL_OTU = "/openconfig-terminal-device/terminal-device/logical-channels/channel/otn/state/1"
	XPATH_LOGICAL_CHANNEL_ETH = "/openconfig-terminal-device/terminal-device/logical-channels/channel/ethernet/state"
	XPATH_INTERFACE           = "/openconfig-interfaces/interfaces/interface/state/counters"
	XPATH_TRANSCEIVER_CLIENT  = "/openconfig-platform/components/component/transceiver/state"
    XPATH_PHYSICAL_CHANNEL_1  = "/openconfig-platform/components/component/transceiver/physical-channels/channel=1/state"
    XPATH_PHYSICAL_CHANNEL_2  = "/openconfig-platform/components/component/transceiver/physical-channels/channel=2/state"
    XPATH_PHYSICAL_CHANNEL_3  = "/openconfig-platform/components/component/transceiver/physical-channels/channel=3/state"
    XPATH_PHYSICAL_CHANNEL_4  = "/openconfig-platform/components/component/transceiver/physical-channels/channel=4/state"
	XPATH_OPTICAL_PORT        = "/openconfig-platform:components/component/port/openconfig-transport-line-common:optical-port/state"
	XPATH_OSC                 = "/openconfig-optical-amplifier:optical-amplifier/supervisory-channels/supervisory-channel/state"
	XPATH_EDFA                = "/openconfig-optical-amplifier:optical-amplifier/amplifiers/amplifier/state"

)

type logicalChannelInfo struct {
	index  string
}

func init() {
	err := register("/openconfig-transport-pms:get-pm",
		&appInfo{appType: reflect.TypeOf(TransPmsApp{}),
			ygotRootType: nil,
			isNative:     false})
	if err != nil {
		glog.Fatal("TransPmsApp:  Register openconfig-transport-pms app module with App Interface failed with error=", err)
	}

	err = register("/openconfig-transport-pms:clear-pm-data",
		&appInfo{appType: reflect.TypeOf(TransPmsApp{}),
			ygotRootType: nil,
			isNative:     false})
	if err != nil {
		glog.Fatal("TransPmsApp:  Register openconfig-transport-pms app module with App Interface failed with error=", err)
	}

	err = addModel(&ModelData{Name: "openconfig-transport-pms",
		Org: "Alibaba transport working group",
		Ver: "2021-11-23"})
	if err != nil {
		glog.Fatal("TransPmsApp:  Adding model data to appinterface failed with error=", err)
	}
}

func (app *TransPmsApp) initialize(data appData) {
	app.path = NewPathInfo(data.path)
	app.reqData = data.payload
	app.ygotRoot = data.ygotRoot
	app.ygotTarget = data.ygotTarget
}

func (app *TransPmsApp) translateCreate(d *db.DB) ([]db.WatchKeys, error) {
	return nil, tlerr.NotSupported("Unsupported")
}

func (app *TransPmsApp) translateUpdate(d *db.DB) ([]db.WatchKeys, error) {
	return nil, tlerr.NotSupported("Unsupported")
}

func (app *TransPmsApp) translateReplace(d *db.DB) ([]db.WatchKeys, error) {
	return nil, tlerr.NotSupported("Unsupported")
}

func (app *TransPmsApp) translateMDBReplace(numDB db.NumberDB) ([]db.WatchKeys, error) {
	return nil, tlerr.NotSupported("Unsupported")
}

func (app *TransPmsApp) translateDelete(d *db.DB) ([]db.WatchKeys, error) {
	return nil, tlerr.NotSupported("Unsupported")
}

func (app *TransPmsApp) translateGet(dbs [db.MaxDB]*db.DB) error {
	return tlerr.NotSupported("Unsupported")
}

func (app *TransPmsApp) translateMDBGet(mdb db.MDB) error {
	return tlerr.NotSupported("Unsupported")
}

func (app *TransPmsApp) translateGetRegex(mdb db.MDB) error {
	return tlerr.NotSupported("Unsupported")
}

func (app *TransPmsApp) translateAction(mdb db.MDB) error {
	glog.V(3).Info("TransPmsApp translateAction : init maps for pm query")
	initComponents(mdb)
	initLogicalChannel(mdb)
	initEntityPmpTypeMap()
	initPmpTypePathMap()
	initPathPmpMap()
	return nil
}

func (app *TransPmsApp) translateSubscribe(dbs [db.MaxDB]*db.DB, path string) (*notificationOpts, *notificationInfo, error) {
	return nil, nil, tlerr.NotSupported("Unsupported")
}

func (app *TransPmsApp) processCreate(d *db.DB) (SetResponse, error) {
	return SetResponse{}, tlerr.NotSupported("Unsupported")
}

func (app *TransPmsApp) processUpdate(d *db.DB) (SetResponse, error) {
	return SetResponse{}, tlerr.NotSupported("Unsupported")
}

func (app *TransPmsApp) processReplace(d *db.DB) (SetResponse, error) {
	return SetResponse{}, tlerr.NotSupported("Unsupported")
}

func (app *TransPmsApp) processMDBReplace(numDB db.NumberDB) (SetResponse, error) {
	return SetResponse{}, tlerr.NotSupported("Unsupported")
}

func (app *TransPmsApp) processDelete(d *db.DB) (SetResponse, error) {
	return SetResponse{}, tlerr.NotSupported("Unsupported")
}

func (app *TransPmsApp) processGet(dbs [db.MaxDB]*db.DB) (GetResponse, error) {
	return GetResponse{}, tlerr.NotSupported("Unsupported")
}

func (app *TransPmsApp) processMDBGet(mdb db.MDB) (GetResponse, error) {
	return GetResponse{}, tlerr.NotSupported("Unsupported")
}

func (app *TransPmsApp) processGetRegex(mdb db.MDB) ([]GetResponseRegex, error) {
	return nil, tlerr.NotSupported("Unsupported")
}

func (app *TransPmsApp) processAction(mdb db.MDB) (ActionResponse, error) {
	input,err := parseGetPmInput(app.reqData)
	if err != nil {
		return ActionResponse{ErrSrc: AppErr}, err
	}

	err = validateGetPmRequest(input)
	if err != nil {
		return ActionResponse{ErrSrc: AppErr}, err
	}

	type Output struct {
		PmResults *PmResults `json:"pm-results"`
	}

	type rpcResponse struct {
		Output *Output `json:"output"`
	}

	var result rpcResponse
	result.Output = new(Output)
	result.Output.PmResults = buildPmResults(input, mdb)

	payload, err := json.Marshal(result)
	if err != nil {
		glog.Errorf("encode rpc response failed as %v", err)
		return ActionResponse{ErrSrc: AppErr}, err
	}

	return ActionResponse{Payload: payload}, nil
}

func validateGetPmRequest(input *GetPmInput) error {
	if input.PmBinType == nil {
		return tlerr.New("pm-bin-type is mandatory")
	}

	if *input.PmBinType != "current" && *input.PmBinType != "history" {
		return tlerr.New("the value of pm-bin-type is invalid")
	}

	if input.PmTimePeriod == nil {
		return tlerr.New("pm-time-period is mandatory")
	}

	if *input.PmTimePeriod != "15min" && *input.PmTimePeriod != "24h" {
		return tlerr.New("the value of pm-bin-period is invalid")
	}

	if *input.PmBinType == "current" {
		if input.StartNumberOfBin != nil || input.EndNumberOfBin != nil {
			return tlerr.New("start-number-of-bin or end-number-of-bin is invalid as pm-bin-type is current")
		}
	}

	return nil
}

func parseGetPmInput(reqData []byte) (*GetPmInput, error) {
	condition := new(GetPmInput)

	err := parsePostBody(reqData, condition)
	if err != nil {
		return nil, err
	}

	return condition, nil
}

func updateGetPmInputWithDefaultValue(input *GetPmInput) {
	defaultVal := "all"
	if input.EntityName == nil {
		input.EntityName = &defaultVal
	}

	if input.PmpType == nil {
		input.PmpType = &defaultVal
	}

	if input.PmParameter == nil {
		input.PmParameter = &defaultVal
	}

	defaultBinNo := uint16(1)
	if input.StartNumberOfBin == nil {
		input.StartNumberOfBin = &defaultBinNo
	}

	if input.EndNumberOfBin == nil {
		input.EndNumberOfBin = &defaultBinNo
	}
}

func buildPmResults(input *GetPmInput, mdb db.MDB) *PmResults {
	var prs *PmResults

	updateGetPmInputWithDefaultValue(input)
	condition := PmQueryCondition {
		binNo:      0,
		binType:    *input.PmBinType,
		entityName: *input.EntityName,
		timePeriod: *input.PmTimePeriod,
	}

	if *input.PmBinType == "current" {
		condition.binNo = 0
		pr := buildPmResult(condition, mdb)
		if pr == nil {
			return nil
		}
		prs = new(PmResults)
		prs.PmResult = make([]*PmResult, 0)
		prs.PmResult = append(prs.PmResult, pr)
		return prs
	}

	start := *input.StartNumberOfBin
	end := *input.EndNumberOfBin

	var prArray []*PmResult
	for i := start; i <= end; i++ {
		condition.binNo = i
		pr := buildPmResult(condition, mdb)
		if pr == nil {
			return nil
		}

		prArray = append(prArray, pr)
	}

	prs = new(PmResults)
	prs.PmResult = prArray
	return prs
}

func buildPmResult(condition PmQueryCondition, mdb db.MDB) *PmResult {
	var start, end time.Time

	if condition.binNo == 0 {
		start, end = getCurrentCollectStartEndTime(condition)
	} else {
		start, end = getHistoryCollectStartEndTime(condition)
	}
	condition.tableLastKey = getCountersTableLastKey(condition)

	pr := new(PmResult)
	pr.PmBinNo = condition.binNo
	pr.PmBinState = new(PmBinState)
	pr.PmBinState.CollectStartTime = start.Format(CUSTOM_TIME_FORMAT)
	pr.PmBinState.CollectEndTime = end.Format(CUSTOM_TIME_FORMAT)
	pr.PmBinState.BinTimePeriod = condition.timePeriod
	pr.PmInstances = buildPmInstances(condition, mdb)
	return pr
}

func buildPmInstances(condition PmQueryCondition, mdb db.MDB) *PmInstances {
	var pmInstArray []*PmInstance

	buildOneInstance := func(condition PmQueryCondition, mdb db.MDB) {
		inst := buildPmInstance(condition, mdb)
		if inst == nil {
			return
		}
		pmInstArray = append(pmInstArray, inst)
	}

	if len(condition.entityName) != 0 && condition.entityName != "all" {
		buildOneInstance(condition, mdb)
		goto buildInsts
	}

	// build components
	for _, cpt := range componentList {
		condition.entityName = cpt
		dbName := db.GetMDBNameFromEntity(cpt)
		if condition.binType == PMCurrent && !isComponentActive(mdb[dbName][db.StateDB], cpt) {
			continue
		}
		buildOneInstance(condition, mdb)
	}

	// build logical-channels
	for lcName, _ := range logicalChannelMap {
		condition.entityName = lcName
		buildOneInstance(condition, mdb)
	}

buildInsts:
	if len(pmInstArray) == 0 {
		return nil
	}
	insts := new(PmInstances)
	insts.PmInstance = pmInstArray
	return insts
}

func buildPmInstance(condition PmQueryCondition, mdb db.MDB) *PmInstance {
	iterms := buildItems(condition, mdb)
	if iterms == nil {
		return nil
	}
	inst := new(PmInstance)
	inst.EntityName = condition.entityName
	inst.PmItems = iterms
	return inst
}

func buildNoneGuageItems(condition PmQueryCondition, mdb db.MDB, pmDbInfo *PmDbInfo, pmpInfos []PmpInfo, pmpType string) []*PmItem {
	var items []*PmItem

	data, err := getRedisData(mdb[pmDbInfo.dbName][pmDbInfo.dbNum], pmDbInfo.ts, pmDbInfo.key)
	if err != nil {
		return items
	}

	for _, pmpInfo := range pmpInfos {
		if pmpInfo.IsGuage {
			continue
		}
		item := buildItem(condition, data, &pmpInfo, pmpType)
		if item == nil {
			continue
		}
		items = append(items, item)
	}

	return items
}

func genGuageNodeKey(oldKey db.Key, suffix string) db.Key {
	var newKey []string
	num := oldKey.Len()
	for i := 0; i < num; i++ {
		elmt := oldKey.Get(i)
		if i + 1 == num - 1 && len(suffix) != 0 { // 倒数第二个key需要添加后缀
			elmt += "_" + suffix
		}
		newKey = append(newKey, elmt)
	}

	return db.Key{Comp: newKey}
}

func buildGuageItems(condition PmQueryCondition, mdb db.MDB, pmDbInfo *PmDbInfo, pmpInfos []PmpInfo, pmpType string) []*PmItem {
	var items []*PmItem

	suffixNodeMap := map[string]string {
		"max"  : "max",
		"min"  : "min",
		"avg"  : "avg",
		"inst" : "instant",
	}

	for _, pmpInfo := range pmpInfos {
		if !pmpInfo.IsGuage {
			continue
		}

		suffix := getFieldNameByNodeName(pmpInfo.NodeName)
		newKey := genGuageNodeKey(pmDbInfo.key, suffix)
		data, err := getRedisData(mdb[pmDbInfo.dbName][pmDbInfo.dbNum], pmDbInfo.ts, newKey)
		if err != nil {
			continue
		}

		baseParam := pmpInfo.PmParameter
		if pmpInfo.OnlyInstant {
			pmpInfo.NodeName = "instant"
			if item := buildItem(condition, data, &pmpInfo, pmpType); item != nil {
				items = append(items, item)
			}
			continue
		}

		for k, v := range suffixNodeMap {
			pmpInfo.PmParameter = baseParam + "-" + k
			pmpInfo.NodeName = v
			item := buildItem(condition, data, &pmpInfo, pmpType)
			if item == nil {
				continue
			}
			items = append(items, item)
		}
	}

	return items
}

func buildItems(condition PmQueryCondition, mdb db.MDB) *PmItems {
	eType := getEntityType(condition.entityName)
	if eType == UNKNOWN {
		return nil
	}

	items := new(PmItems)
	items.PmItem = make([]*PmItem, 0)
	for _, pmpType := range entityPmpTypeMap[eType] {
		for _, path := range pmpTypePathMap[pmpType] {
			pmDbInfo := getPmDbInfo(condition, path)
			if pmDbInfo == nil {
				continue
			}

			pmpInfos := pathPmpMap[path]
			if strings.HasPrefix(pmpType, "fec-") {
				pmpInfos = pathPmpMap["fec-*"]
			}

		    if noGuageItems := buildNoneGuageItems(condition, mdb, pmDbInfo, pmpInfos, pmpType); noGuageItems != nil {
				items.PmItem = append(items.PmItem, noGuageItems...)
			}

			if guageItems := buildGuageItems(condition, mdb, pmDbInfo, pmpInfos, pmpType); guageItems != nil {
				items.PmItem = append(items.PmItem, guageItems...)
			}
		}
	}

	if len(items.PmItem) == 0 {
		return nil
	}

	return items
}

func buildItem(condition PmQueryCondition, data db.Value, pmInfo *PmpInfo, pmpType string) *PmItem {
	pmValue := data.Get(pmInfo.NodeName)
	if len(pmValue) == 0 {
		return nil
	}

	validity := data.Get("validity")
	if len(validity) == 0 {
		validity = "invalid"
	}

	if strings.HasPrefix(pmpType, "fec-") {
		pmpType = "fec"
	}

	mdt := getMonitoringDateTime(data, pmInfo)
	item := PmItem{
		PmTimePeriod:       condition.timePeriod,
		PmpType:            pmpType,
		PmParameter:        pmInfo.PmParameter,
		MonitoringDateTime: mdt,
		PmValue:            pmValue,
		Validity:           validity,
	}

	return &item
}

func getMonitoringDateTime(data db.Value, pmInfo *PmpInfo) string {
	var mdt string
	var field string
	if strings.Contains(pmInfo.PmParameter, "-min") {
		field = "min-time"
	} else if strings.Contains(pmInfo.PmParameter, "-max") {
		field = "max-time"
	}

	if len(field) != 0 {
		timeStr := data.Get(field)
		if len(timeStr) != 0 {
			timeInt64, err := strconv.ParseInt(timeStr, 10, 64)
			if err != nil {
				glog.Errorf("parse the %s failed as %v", field, err)
				return ""
			}
			mdt = time.Unix(0, timeInt64).Format(CUSTOM_TIME_FORMAT)
		}
	}

	return mdt
}

func getPmTimePeriod(timePeriodStr string) PMTimePeriod {
	if timePeriodStr == "15min" {
		return MIN_15
	}

	if timePeriodStr == "24h" {
		return HOUR_24
	}

	return NOT_APPLICABLE
}

func getEntityType(entityName string) PMEntityType {
	entityType := UNKNOWN
	prefix := strings.Split(entityName, "-")[0]
	switch prefix {
	case "CHASSIS":
		entityType = CHASSIS
	case "LINECARD":
		entityType = LINECARD
	case "FAN":
		entityType = FAN
	case "PSU":
		entityType = PSU
	case "OCH":
		entityType = OCH
	case "TRANSCEIVER":
		if ok, _ := regexp.MatchString("TRANSCEIVER-1-[1-4]-L[1-2]", entityName); ok {
			entityType = TRANSCEIVER_L
		} else {
			entityType = TRANSCEIVER_C
		}
	case "ODUCN":
		entityType = LOGICAL_CHANNEL_ODU
	case "OTUCN":
		entityType = LOGICAL_CHANNEL_OTU
	case "GE":
		entityType = LOGICAL_CHANNEL_ETH
	case "PORT", "OSC", "AMPLIFIER":
		entityType = OPTICAL_DEVICE
	default:
		glog.Errorf("get entity %s type failed", entityName)
	}

	return entityType
}

func initEntityPmpTypeMap() {
	entityPmpTypeMap = make(map[PMEntityType][]string)
	entityPmpTypeMap[CHASSIS] = []string{"chassis-temperature"}
	entityPmpTypeMap[FAN] = []string{"equipment-temperature", "fan-speed"}
	entityPmpTypeMap[LINECARD] = []string{"equipment-temperature"}
	entityPmpTypeMap[PSU] = []string{"equipment-temperature", "power-supply"}
	entityPmpTypeMap[OCH] = []string{"coherent-optical-interface", "optical-power"}
	entityPmpTypeMap[TRANSCEIVER_L] = []string{"equipment-temperature"}
	entityPmpTypeMap[TRANSCEIVER_C] = []string{"equipment-temperature", "fec-transceiver", "optical-power-lane"}
	entityPmpTypeMap[LOGICAL_CHANNEL_OTU] = []string{"otu", "fec-otu", "coherent-optical-interface"}
	entityPmpTypeMap[LOGICAL_CHANNEL_ODU] = []string{"odu"}
	entityPmpTypeMap[LOGICAL_CHANNEL_ETH] = []string{"ethernet"}
	entityPmpTypeMap[OPTICAL_DEVICE] = []string{"optical-power"}
}

func initPmpTypePathMap() {
	pmpTypePathMap = make(map[string][]string)

	pmpTypePathMap["chassis-temperature"] = []string{
		XPATH_TEMP_CHASSIS,
	}

	pmpTypePathMap["equipment-temperature"] = []string{
		XPATH_TEMP_EQUIPMENT,
	}

	pmpTypePathMap["fan-speed"] = []string{
		XPATH_FAN,
	}

	pmpTypePathMap["power-supply"] = []string{
		XPATH_PSU,
	}

	pmpTypePathMap["coherent-optical-interface"] = []string{
		XPATH_OCH_FIRST,
		XPATH_LOGICAL_CHANNEL_OTU,
	}

	pmpTypePathMap["optical-power-lane"] = []string{
		XPATH_TRANSCEIVER_CLIENT,
		XPATH_PHYSICAL_CHANNEL_1,
		XPATH_PHYSICAL_CHANNEL_2,
		XPATH_PHYSICAL_CHANNEL_3,
		XPATH_PHYSICAL_CHANNEL_4,
	}

	pmpTypePathMap["otu"] = []string{
		XPATH_LOGICAL_CHANNEL_OTN,
	}

	pmpTypePathMap["odu"] = []string{
		XPATH_LOGICAL_CHANNEL_OTN,
	}

	pmpTypePathMap["ethernet"] = []string{
		XPATH_LOGICAL_CHANNEL_ETH,
	}

	pmpTypePathMap["fec-otu"] = []string{
		XPATH_LOGICAL_CHANNEL_OTN,
	}

	pmpTypePathMap["fec-transceiver"] = []string{
		XPATH_TRANSCEIVER_CLIENT,
	}

	pmpTypePathMap["optical-power"] = []string{
		XPATH_OCH_SECOND,
		XPATH_OPTICAL_PORT,
		XPATH_OSC,
		XPATH_EDFA,
	}
}

func initPathPmpMap() {
	pathPmpMap = make(map[string][]PmpInfo)
	pathPmpMap[XPATH_LOGICAL_CHANNEL_OTN] = []PmpInfo{
		{PmParameter: "EB", NodeName: "errored-blocks"},
		{PmParameter: "ES", NodeName: "errored-seconds"},
		{PmParameter: "SES", NodeName: "severely-errored-seconds"},
		{PmParameter: "UAS", NodeName: "unavailable-seconds"},
		{PmParameter: "CV", NodeName: "code-violations"},
		{PmParameter: "BBE", NodeName: "background-block-errors"},
		{PmParameter: "DELAY", NodeName: "delay", IsGuage: true},
	}

	pathPmpMap[XPATH_LOGICAL_CHANNEL_OTU] = []PmpInfo{
		{PmParameter: "ESNR", NodeName: "esnr", IsGuage: true},
		{PmParameter: "Q-factor", NodeName: "q-value", IsGuage: true},
	}

	pathPmpMap["fec-*"] = []PmpInfo{
		{PmParameter: "BE-FEC", NodeName: "fec-corrected-bits"},
		{PmParameter: "UBE-FEC", NodeName: "fec-uncorrectable-blocks"},
		{PmParameter: "PRE-FEC-BER", NodeName: "pre-fec-ber", IsGuage: true},
		{PmParameter: "POST-FEC-BER", NodeName: "post-fec-ber", IsGuage: true},
	}

	pathPmpMap[XPATH_OCH_FIRST] = []PmpInfo{
		{PmParameter: "tx-laser-age", NodeName: "tx-laser-age", IsGuage: true, OnlyInstant: true},
		{PmParameter: "tx-mod-bias-x-i", NodeName: "tx-mod-bias-x-i", IsGuage: true},
		{PmParameter: "tx-mod-bias-x-q", NodeName: "tx-mod-bias-x-q", IsGuage: true},
		{PmParameter: "tx-mod-bias-x-ph", NodeName: "tx-mod-bias-x-ph", IsGuage: true},
		{PmParameter: "tx-mod-bias-y-i", NodeName: "tx-mod-bias-y-i", IsGuage: true},
		{PmParameter: "tx-mod-bias-y-q", NodeName: "tx-mod-bias-y-q", IsGuage: true},
		{PmParameter: "tx-mod-bias-y-ph", NodeName: "tx-mod-bias-y-ph", IsGuage: true},
		{PmParameter: "CD", NodeName: "chromatic-dispersion", IsGuage: true},
		{PmParameter: "OSNR", NodeName: "osnr", IsGuage: true},
		{PmParameter: "DGD", NodeName: "polarization-mode-dispersion", IsGuage: true},
		{PmParameter: "PDL", NodeName: "polarization-dependent-loss", IsGuage: true},
		{PmParameter: "sop-change-rate", NodeName: "sop-change-rate", IsGuage: true},
		//{PmParameter: "ESNR", NodeName: "esnr", IsGuage: true},
		//{PmParameter: "Q-factor", NodeName: "q-value", IsGuage: true},
	}

	pathPmpMap[XPATH_OCH_SECOND] = []PmpInfo{
		{PmParameter: "OPT", NodeName: "output-power", IsGuage: true},
		{PmParameter: "OPR", NodeName: "input-power", IsGuage: true},
		{PmParameter: "BIAS", NodeName: "laser-bias-current", IsGuage: true},
	}

	pathPmpMap[XPATH_TEMP_CHASSIS] = []PmpInfo{
		{PmParameter: "Tinlet", NodeName: ""},
		{PmParameter: "Toutlet", NodeName: "temperature", IsGuage: true},
	}

	pathPmpMap[XPATH_TEMP_EQUIPMENT] = []PmpInfo{
		{PmParameter: "Tmodule", NodeName: "temperature", IsGuage: true},
	}

	pathPmpMap[XPATH_TRANSCEIVER_CLIENT] = []PmpInfo{
		//{PmParameter: "BIAS", NodeName: "laser-bias-current", IsGuage: true},
		{PmParameter: "OPR-total", NodeName: "input-power", IsGuage: true},
		{PmParameter: "OPT-total", NodeName: "output-power", IsGuage: true},
	}

	pathPmpMap[XPATH_PHYSICAL_CHANNEL_1] = []PmpInfo{
		{PmParameter: "BIAS-lane-1", NodeName: "laser-bias-current", IsGuage: true},
		{PmParameter: "OPR-lane-1", NodeName: "input-power", IsGuage: true},
		{PmParameter: "OPT-lane-1", NodeName: "output-power", IsGuage: true},
	}
	pathPmpMap[XPATH_PHYSICAL_CHANNEL_2] = []PmpInfo{
		{PmParameter: "BIAS-lane-2", NodeName: "laser-bias-current", IsGuage: true},
		{PmParameter: "OPR-lane-2", NodeName: "input-power", IsGuage: true},
		{PmParameter: "OPT-lane-2", NodeName: "output-power", IsGuage: true},
	}
	pathPmpMap[XPATH_PHYSICAL_CHANNEL_3] = []PmpInfo{
		{PmParameter: "BIAS-lane-3", NodeName: "laser-bias-current", IsGuage: true},
		{PmParameter: "OPR-lane-3", NodeName: "input-power", IsGuage: true},
		{PmParameter: "OPT-lane-3", NodeName: "output-power", IsGuage: true},
	}
	pathPmpMap[XPATH_PHYSICAL_CHANNEL_4] = []PmpInfo{
		{PmParameter: "BIAS-lane-4", NodeName: "laser-bias-current", IsGuage: true},
		{PmParameter: "OPR-lane-4", NodeName: "input-power", IsGuage: true},
		{PmParameter: "OPT-lane-4", NodeName: "output-power", IsGuage: true},
	}

	pathPmpMap[XPATH_PSU] = []PmpInfo{
		{PmParameter: "input-current", NodeName: "input-current", IsGuage: true},
		{PmParameter: "output-current", NodeName: "output-current", IsGuage: true},
		{PmParameter: "input-voltage", NodeName: "input-voltage", IsGuage: true},
		{PmParameter: "output-voltage", NodeName: "output-voltage", IsGuage: true},
		{PmParameter: "output-power", NodeName: "output-power", IsGuage: true},
	}

	pathPmpMap[XPATH_FAN] = []PmpInfo{
		{PmParameter: "fan-speed", NodeName: "speed", IsGuage: true},
	}

	pathPmpMap[XPATH_OPTICAL_PORT] = []PmpInfo{
		{PmParameter: "OPT", NodeName: "output-power", IsGuage: true},
		{PmParameter: "OPR", NodeName: "input-power", IsGuage: true},
	}

	pathPmpMap[XPATH_OSC] = []PmpInfo{
		{PmParameter: "OPT", NodeName: "output-power", IsGuage: true},
		{PmParameter: "OPR", NodeName: "input-power", IsGuage: true},
		{PmParameter: "BIAS", NodeName: "laser-bias-current", IsGuage: true},
	}

	pathPmpMap[XPATH_EDFA] = []PmpInfo{
		{PmParameter: "OPT", NodeName: "output-power-total", IsGuage: true},
		{PmParameter: "OPR", NodeName: "input-power-total", IsGuage: true},
		{PmParameter: "BIAS", NodeName: "laser-bias-current", IsGuage: true},
	}

	pathPmpMap[XPATH_INTERFACE] = []PmpInfo {
		{PmParameter: "InOctets", NodeName: "in-octets"},
		{PmParameter: "InPkts", NodeName: "in-pkts"},
		{PmParameter: "InUnicastPkts", NodeName: "in-unicast-pkts"},
		{PmParameter: "InBroadcastPkts", NodeName: "in-broadcast-pkts"},
		{PmParameter: "InMulticastPkts", NodeName: "in-multicast-pkts"},
		{PmParameter: "InDiscards", NodeName: "in-discards"},
		{PmParameter: "InErrors", NodeName: "in-errors"},
		{PmParameter: "InUnknowns", NodeName: "in-unknown-protos"},
		{PmParameter: "InFcsErrors", NodeName: "in-fcs-errors"},
		{PmParameter: "OutOctets", NodeName: "out-octets"},
		{PmParameter: "OutPkts", NodeName: "out-pkts"},
		{PmParameter: "OutUnicastPkts", NodeName: "out-unicast-pkts"},
		{PmParameter: "OutBroadcastPkts", NodeName: "out-broadcast-pkts"},
		{PmParameter: "OutMulticastPkts", NodeName: "out-multicast-pkts"},
		{PmParameter: "OutDiscards", NodeName: "out-discards"},
		{PmParameter: "OutErrors", NodeName: "out-errors"},
		{PmParameter: "InControls", NodeName: "in-mac-control-frames"},
		{PmParameter: "InPauses", NodeName: "in-mac-pause-frames"},
		{PmParameter: "In802q", NodeName: "in-8021q-frames"},
		{PmParameter: "InOversizePkts", NodeName: "in-oversize-frames"},
		{PmParameter: "InUndersizePkts", NodeName: "in-undersize-frames"},
		{PmParameter: "InJabbers", NodeName: "in-jabber-frames"},
		{PmParameter: "InFragments", NodeName: "in-fragment-frames"},
		{PmParameter: "InCRCAlignErrors", NodeName: "in-crc-errors"},
		{PmParameter: "OutControls", NodeName: "out-mac-control-frames"},
		{PmParameter: "OutPauses", NodeName: "out-mac-pause-frames"},
		{PmParameter: "Out802q", NodeName: "out-8021q-frames"},
		{PmParameter: "InPkts64Octets", NodeName: "in-frames-64-octets"},
		{PmParameter: "InPkts65to127Octets", NodeName: "in-frames-65-127-octets"},
		{PmParameter: "InPkts128to255Octets", NodeName: "in-frames-128-255-octets"},
		{PmParameter: "InPkts256to511Octets", NodeName: "in-frames-256-511-octets"},
		{PmParameter: "InPkts512to1023Octets", NodeName: "in-frames-512-1023-octets"},
		{PmParameter: "InPkts1024to1518Octets", NodeName: "in-frames-1024-1518-octets"},
	}

	pathPmpMap[XPATH_LOGICAL_CHANNEL_ETH] = []PmpInfo{
		{PmParameter: "InOctets", NodeName: "rx-octets"},
		{PmParameter: "InPkts", NodeName: "rx-frame"},
		{PmParameter: "InBroadcastPkts", NodeName: "rx-broadcast"},
		{PmParameter: "InMulticastPkts", NodeName: "rx-multicast"},
		{PmParameter: "InOversizePkts", NodeName: "in-oversize-frames"},
		{PmParameter: "InUndersizePkts", NodeName: "in-undersize-frames"},
		{PmParameter: "InJabbers", NodeName: "in-jabber-frames"},
		{PmParameter: "InFragments", NodeName: "in-fragment-frames"},
		{PmParameter: "InCRCAlignErrors", NodeName: "rx-crc-align"},
		{PmParameter: "OutErrors", NodeName: "out-crc-errors"},
		{PmParameter: "OutOctets", NodeName: "tx-octets"},
		{PmParameter: "OutPkts", NodeName: "tx-frame"},
		{PmParameter: "OutBroadcastPkts", NodeName: "tx-broadcast"},
		{PmParameter: "OutMulticastPkts", NodeName: "tx-multicast"},
		{PmParameter: "CV", NodeName: "in-pcs-bip-errors"},
		{PmParameter: "ES", NodeName: "in-pcs-errored-seconds"},
		{PmParameter: "SES", NodeName: "in-pcs-severely-errored-seconds"},
		{PmParameter: "UAS", NodeName: "in-pcs-unavailable-seconds"},
		{PmParameter: "InPkts64Octets", NodeName: "rx-64b"},
		{PmParameter: "InPkts65to127Octets", NodeName: "rx-65b-127b"},
		{PmParameter: "InPkts128to255Octets", NodeName: "rx-128b-255b"},
		{PmParameter: "InPkts256to511Octets", NodeName: "rx-256b-511b"},
		{PmParameter: "InPkts512to1023Octets", NodeName: "rx-512b-1023b"},
		{PmParameter: "InPkts1024to1518Octets", NodeName: "rx-1024b-1518b"},
	}
}

func isOpticalPort(name string) bool {
	if !strings.HasPrefix(name, "PORT") {
		return false
	}

	if !strings.HasSuffix(name, "IN") && !strings.HasSuffix(name, "OUT") {
		return false
	}

	// edfa ports optical-power was displayed by AMPLIFIER
	if strings.Contains(name, "EDFA") {
		return false
	}

	return true
}

func initLogicalChannel(mdb db.MDB) {
	logicalChannelMap = make(map[string]*logicalChannelInfo)
	for _, dbs := range mdb {
		d := dbs[db.StateDB]
		ts := asTableSpec("LOGICAL_CHANNEL")
		keys, _ := d.GetKeys(ts)
		if len(keys) == 0 {
			continue
		}

		for _, key := range keys {
			data, _ := getRedisData(d, ts, key)
			logicalChannelMap[data.Get("description")] = &logicalChannelInfo {
				index: data.Get("index"),
			}
		}
	}
	glog.V(2).Infof("entity for logical-channel %v", logicalChannelMap)
}

func initComponents(mdb db.MDB) {
	cmpTables := []string{"FAN", "PSU", "LINECARD", "CHASSIS", "TRANSCEIVER", "OCH", "PORT", "OSC", "AMPLIFIER"}
	componentList = make([]string, 0)
	for _, tbl := range cmpTables {
		keys, err := db.GetTableKeysByDbNum(mdb, asTableSpec(tbl), db.StateDB)
		if err != nil {
			continue
		}

		for _, key := range keys {
			if key.Len() != 1 {
				continue
			}
			name := key.Get(0)
			// only optical port need to add into list
			if tbl == "PORT" && !isOpticalPort(name) {
				continue
			}
			componentList = append(componentList, name)
		}
	}
	glog.V(2).Infof("entity for component %v", componentList)
}

func getPmDbInfo(condition PmQueryCondition, path string) *PmDbInfo {
	dbInfo := new(PmDbInfo)
	entityName := condition.entityName
	dbName := db.GetMDBNameFromEntity(entityName)
	dbInfo.dbName = dbName
	dbInfo.dbNum = db.CountersDB
	switch path {
	case XPATH_LOGICAL_CHANNEL_OTN, XPATH_LOGICAL_CHANNEL_OTU:
		if logicalChannelMap[entityName] == nil {
			return nil
		}
		dbInfo.ts = asTableSpec("OTN")
		dbInfo.key = asKey("CH" + logicalChannelMap[entityName].index)
	case XPATH_LOGICAL_CHANNEL_ETH:
		if logicalChannelMap[entityName] == nil {
			return nil
		}
		dbInfo.ts = asTableSpec("ETHERNET")
		dbInfo.key = asKey("CH" + logicalChannelMap[entityName].index)
	case XPATH_TEMP_EQUIPMENT, XPATH_TEMP_CHASSIS:
		elmts := strings.Split(entityName, "-")
		dbInfo.ts = asTableSpec(elmts[0])
		dbInfo.key = asKey(entityName)
	case XPATH_OCH_FIRST, XPATH_OCH_SECOND:
		dbInfo.ts = asTableSpec("OCH")
		dbInfo.key = asKey(entityName)
	case XPATH_TRANSCEIVER_CLIENT:
		dbInfo.ts = asTableSpec("TRANSCEIVER")
		dbInfo.key = asKey(entityName)
	case XPATH_PHYSICAL_CHANNEL_1:
		dbInfo.ts = asTableSpec("TRANSCEIVER")
		dbInfo.key = asKey(entityName, "CH-1")
	case XPATH_PHYSICAL_CHANNEL_2:
		dbInfo.ts = asTableSpec("TRANSCEIVER")
		dbInfo.key = asKey(entityName, "CH-2")
	case XPATH_PHYSICAL_CHANNEL_3:
		dbInfo.ts = asTableSpec("TRANSCEIVER")
		dbInfo.key = asKey(entityName, "CH-3")
	case XPATH_PHYSICAL_CHANNEL_4:
		dbInfo.ts = asTableSpec("TRANSCEIVER")
		dbInfo.key = asKey(entityName, "CH-4")
	case XPATH_PSU:
		dbInfo.ts = asTableSpec("PSU")
		dbInfo.key = asKey(entityName)
	case XPATH_FAN:
		dbInfo.ts = asTableSpec("FAN")
		dbInfo.key = asKey(entityName)
	case XPATH_OPTICAL_PORT:
		dbInfo.ts = asTableSpec("PORT")
		dbInfo.key = asKey(entityName)
	case XPATH_OSC:
		dbInfo.ts = asTableSpec("OSC")
		dbInfo.key = asKey(entityName)
	case XPATH_EDFA:
		dbInfo.ts = asTableSpec("AMPLIFIER")
		dbInfo.key = asKey(entityName)
	case "unknown":
		return nil
	default:
		glog.Errorf("entity %s path %s is not matched", entityName, path)
		return nil
	}

	if condition.binType == "history" {
		dbInfo.dbNum = db.HistoryDB
	}

	// add counters table last key. like 15_pm_current or 24_pm_history_1648483200000000000
	dbInfo.key.Comp = append(dbInfo.key.Comp, condition.tableLastKey)

	return dbInfo
}

func getPmType(condition PmQueryCondition) string {
	if condition.timePeriod == "15min" && condition.binType == "current" {
		return PMCurrent15min
	}
	if condition.timePeriod == "24h" && condition.binType == "current" {
		return PMCurrent24h
	}
	if condition.timePeriod == "15min" && condition.binType == "history" {
		return PMHistory15min
	}
	if condition.timePeriod == "24h" && condition.binType == "history" {
		return PMHistory24h
	}
	return ""
}

func getCountersTableLastKey(condition PmQueryCondition) string {
	lastKey := getPmType(condition)
	if condition.binNo == 0 && condition.binType == "current" {
		return lastKey
	}

	if condition.binType != "history" {
		return lastKey
	}

	samplingStartTime, _ := getHistoryCollectStartEndTime(condition)
	lastKey += "_" + strconv.FormatInt(samplingStartTime.UnixNano(), 10)
	glog.Infof("history PM sampling timestamp %s, lastKey %s", samplingStartTime.Format(CUSTOM_TIME_FORMAT), lastKey)

	return lastKey
}

func getLastSamplingTimestamp(period string, currentTime time.Time) time.Time {
	var timestamp time.Time
	if period == "15min" {
		minuteDelta := time.Duration((currentTime.Minute() / 15) * 15) * time.Minute
		timestamp = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(),
			currentTime.Hour(), 0, 0, 0, currentTime.Location()).Add(minuteDelta)
	} else if period == "24h" {
		_, offset := currentTime.Zone()
		timestamp = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(),
			       0, 0, offset, 0, currentTime.Location())
	}

	return timestamp
}

func getCurrentCollectStartEndTime(condition PmQueryCondition) (time.Time, time.Time) {
	currentTime := time.Now()
	timestamp := getLastSamplingTimestamp(condition.timePeriod, currentTime)
	start := timestamp
	end := currentTime
	return start, end
}

func getHistoryCollectStartEndTime(condition PmQueryCondition) (time.Time, time.Time) {
	var startDuration time.Duration
	var endDuration time.Duration
	currentTime := time.Now()
	timestamp := getLastSamplingTimestamp(condition.timePeriod, currentTime)

	startDelta := int(1 - int(condition.binNo) - 1)
	endDelta := int(1 - int(condition.binNo))
	if condition.timePeriod == "15min" {
		startDuration = time.Duration(startDelta * 15) * time.Minute
		endDuration = time.Duration(endDelta * 15) * time.Minute
	} else if condition.timePeriod == "24h" {
		startDuration = time.Duration(startDelta * 24) * time.Hour
		endDuration = time.Duration(endDelta * 24) * time.Hour
	}

	start := timestamp.Add(startDuration)
	end := timestamp.Add(endDuration)

	return start, end
}
