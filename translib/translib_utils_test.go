package translib

import (
	"github.com/Azure/sonic-mgmt-common/translib/db"
	"github.com/Azure/sonic-mgmt-common/translib/ocbinds"
	"github.com/agiledragon/gomonkey"
	"reflect"
	"testing"
)

func TestApp_buildGoStruct(t *testing.T) {
	type args struct {
		gs   interface{}
		data db.Value
	}

	description := "ODU4-1-1-C12"
	index := uint32(10)
	testSignal := false

	config := &ocbinds.OpenconfigTerminalDevice_TerminalDevice_LogicalChannels_Channel_Config{
		AdminState:         ocbinds.OpenconfigTerminalDevice_AdminStateType_DISABLED,
		ClientMappingMode:  ocbinds.OpenconfigTransportTypes_CLIENT_MAPPING_MODE_MODE_1X100G,
		Description:        &description,
		Index:              &index,
		LogicalChannelType: ocbinds.OpenconfigTransportTypes_LOGICAL_ELEMENT_PROTOCOL_TYPE_PROT_ETHERNET,
		LoopbackMode:       ocbinds.OpenconfigTerminalDevice_LoopbackModeType_NONE,
		RateClass:          ocbinds.OpenconfigTransportTypes_TRIBUTARY_RATE_CLASS_TYPE_TRIB_RATE_1000G,
		TestSignal:         &testSignal,
		TribProtocol:       ocbinds.OpenconfigTransportTypes_TRIBUTARY_PROTOCOL_TYPE_PROT_100GE,
	}

	data := db.Value{Field: map[string]string{
		"admin-state" : "DISABLED",
		"client-mapping-mode" : "MODE_1X100G",
		"description" : "ODU4-1-1-C12",
		"index" : "10",
		"logical-channel-type" : "PROT_ETHERNET",
		"loopback-mode" : "NONE",
		"rate-class" : "TRIB_RATE_1000G",
		"test-signal" : "false",
		"trib-protocol" : "PROT_100GE",
	}}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "channel config",
			args:    args{
				gs:   &ocbinds.OpenconfigTerminalDevice_TerminalDevice_LogicalChannels_Channel_Config{},
				data: data,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buildGoStruct(tt.args.gs, tt.args.data)
			if !reflect.DeepEqual(tt.args.gs, config) {
				t.Errorf("buildGoStruct() got = %v, want %v", tt.args.gs, *config)
			}
		})
	}
}

func TestApp_buildGoStructField(t *testing.T) {
	type args struct {
		sf  reflect.StructField
		rv  reflect.Value
		val string
	}

	type yangNode struct {
		TestFloat  *float32
		TestInt    *int32
		TestUint   *uint
		TestString *string
		TestBool   *bool
		TestEnum   ocbinds.E_OpenconfigTerminalDevice_AdminStateType
	}

	testFloat := float32(0.123)
	testInt := int32(-123)
	testUint := uint(123)
	testString := "hello world"
	testBool := true
	testEnum := ocbinds.OpenconfigTerminalDevice_AdminStateType_ENABLED
	tgt := &yangNode {
		TestFloat:  &testFloat,
		TestInt:    &testInt,
		TestUint:   &testUint,
		TestString: &testString,
		TestBool:   &testBool,
		TestEnum:   testEnum,
	}

	testFloat1 := float32(0)
	testInt1 := int32(0)
	testUint1 := uint(0)
	testString1 := ""
	testBool1 := false
	testEnum1 := ocbinds.OpenconfigTerminalDevice_AdminStateType_UNSET
	var dest interface{}
	dest = &yangNode {
		TestFloat:  &testFloat1,
		TestInt:    &testInt1,
		TestUint:   &testUint1,
		TestString: &testString1,
		TestBool:   &testBool1,
		TestEnum:   testEnum1,
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "float32",
			args:    args{
				sf:  reflect.TypeOf(dest).Elem().Field(0),
				rv:  reflect.ValueOf(dest).Elem().Field(0),
				val: "0.123",
			},
			wantErr: false,
		},
		{
			name:    "int",
			args:    args{
				sf:  reflect.TypeOf(dest).Elem().Field(1),
				rv:  reflect.ValueOf(dest).Elem().Field(1),
				val: "-123",
			},
			wantErr: false,
		},
		{
			name:    "uint",
			args:    args{
				sf:  reflect.TypeOf(dest).Elem().Field(2),
				rv:  reflect.ValueOf(dest).Elem().Field(2),
				val: "123",
			},
			wantErr: false,
		},
		{
			name:    "string",
			args:    args{
				sf:  reflect.TypeOf(dest).Elem().Field(3),
				rv:  reflect.ValueOf(dest).Elem().Field(3),
				val: "hello world",
			},
			wantErr: false,
		},
		{
			name:    "bool",
			args:    args{
				sf:  reflect.TypeOf(dest).Elem().Field(4),
				rv:  reflect.ValueOf(dest).Elem().Field(4),
				val: "true",
			},
			wantErr: false,
		},
		{
			name:    "enum",
			args:    args{
				sf:  reflect.TypeOf(dest).Elem().Field(5),
				rv:  reflect.ValueOf(dest).Elem().Field(5),
				val: "ENABLED",
			},
			wantErr: false,
		},
		{
			name:    "bool error",
			args:    args{
				sf:  reflect.TypeOf(dest).Elem().Field(4),
				rv:  reflect.ValueOf(dest).Elem().Field(4),
				val: "ab",
			},
			wantErr: true,
		},
		{
			name:    "enum error",
			args:    args{
				sf:  reflect.TypeOf(dest).Elem().Field(5),
				rv:  reflect.ValueOf(dest).Elem().Field(5),
				val: "UNSET",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := buildGoStructField(tt.args.sf, tt.args.rv, tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("buildGoStructField() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	if !reflect.DeepEqual(tgt, dest) {
		t.Errorf("buildGoStructField() target = %v, processed %v", tgt, dest)
	}
}

func TestApp_convertRequestBodyToInternal(t *testing.T) {
	type args struct {
		gs interface{}
	}

	//var input interface{}
	description := "ODU4-1-1-C12"
	index := uint32(10)
	testSignal := false
	input := &ocbinds.OpenconfigTerminalDevice_TerminalDevice_LogicalChannels_Channel_Config{
		AdminState:         ocbinds.OpenconfigTerminalDevice_AdminStateType_DISABLED,
		ClientMappingMode:  ocbinds.OpenconfigTransportTypes_CLIENT_MAPPING_MODE_MODE_1X100G,
		Description:        &description,
		Index:              &index,
		LogicalChannelType: ocbinds.OpenconfigTransportTypes_LOGICAL_ELEMENT_PROTOCOL_TYPE_PROT_ETHERNET,
		LoopbackMode:       ocbinds.OpenconfigTerminalDevice_LoopbackModeType_NONE,
		RateClass:          ocbinds.OpenconfigTransportTypes_TRIBUTARY_RATE_CLASS_TYPE_TRIB_RATE_1000G,
		TestSignal:         &testSignal,
		TribProtocol:       ocbinds.OpenconfigTransportTypes_TRIBUTARY_PROTOCOL_TYPE_PROT_100GE,
	}

	output := db.Value{Field: map[string]string{
		"admin-state" : "DISABLED",
		"client-mapping-mode" : "MODE_1X100G",
		"description" : "ODU4-1-1-C12",
		"index" : "10",
		"logical-channel-type" : "PROT_ETHERNET",
		"loopback-mode" : "NONE",
		"rate-class" : "TRIB_RATE_1000G",
		"test-signal" : "false",
		"trib-protocol" : "PROT_100GE",
	}}

	tests := []struct {
		name string
		args args
		want db.Value
	}{
		// case1 : body has only leaf nodes
		{
			name: "leaf nodes",
			args: args {
				gs: input,
			},
			want: output,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertRequestBodyToInternal(tt.args.gs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertRequestBodyToInternal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApp_getFieldStringValue(t *testing.T) {
	type args struct {
		t reflect.Type
		v reflect.Value
	}

	testFloat := float32(0.123)
	testInt := int32(-123)
	testUint := uint(123)
	testString := "hello world"
	testBool := true
	testEnum := ocbinds.OpenconfigTerminalDevice_AdminStateType_ENABLED

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// case1 float type
		{
			name:    "float",
			args:    args{
				t: reflect.TypeOf(testFloat),
				v: reflect.ValueOf(testFloat),
			},
			want:    "0.123",
			wantErr: false,
		},
		// case2 int32 type
		{
			name:    "int32",
			args:    args{
				t: reflect.TypeOf(testInt),
				v: reflect.ValueOf(testInt),
			},
			want:    "-123",
			wantErr: false,
		},
		// case3 uint32 type
		{
			name:    "uint32",
			args:    args{
				t: reflect.TypeOf(testUint),
				v: reflect.ValueOf(testUint),
			},
			want:    "123",
			wantErr: false,
		},
		// case4 string type
		{
			name:    "string",
			args:    args{
				t: reflect.TypeOf(testString),
				v: reflect.ValueOf(testString),
			},
			want:    "hello world",
			wantErr: false,
		},
		// case5 bool type
		{
			name:    "bool",
			args:    args{
				t: reflect.TypeOf(testBool),
				v: reflect.ValueOf(testBool),
			},
			want:    "true",
			wantErr: false,
		},
		// case6 enum type
		{
			name:    "enum",
			args:    args{
				t: reflect.TypeOf(testEnum),
				v: reflect.ValueOf(testEnum),
			},
			want:    "ENABLED",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getFieldStringValue(tt.args.t, tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("getFieldStringValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got != tt.want {
				t.Errorf("getFieldStringValue() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApp_getRedisKey(t *testing.T) {
	type args struct {
		prefix string
		mdlKey interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// case1: uint32 key
		{
			name: "uint32 key",
			args: args{
				prefix: "CH",
				mdlKey: uint32(10),
			},
			want: "CH10",
		},
		// case2: int32 key
		{
			name: "int32 key",
			args: args{
				prefix: "CH",
				mdlKey: int32(-10),
			},
			want: "CH-10",
		},
		// case3: string key
		{
			name: "string key",
			args: args{
				prefix: "CH",
				mdlKey: "100",
			},
			want: "CH100",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getRedisKey(tt.args.prefix, tt.args.mdlKey); got != tt.want {
				t.Errorf("getRedisKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApp_getYangMdlKey(t *testing.T) {
	type args struct {
		prefix  string
		tblKey  string
		keyType reflect.Type
	}

	wantUint32 := uint32(10)
	wantInt32 := int32(-10)
	wantString := "10"

	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name:    "uin32",
			args:    args{
				prefix:  "CH",
				tblKey:  "CH10",
				keyType: reflect.TypeOf(wantUint32),
			},
			want:    wantUint32,
			wantErr: false,
		},
		{
			name:    "in32",
			args:    args{
				prefix:  "CH",
				tblKey:  "CH-10",
				keyType: reflect.TypeOf(wantInt32),
			},
			want:    wantInt32,
			wantErr: false,
		},
		{
			name:    "string",
			args:    args{
				prefix:  "CH",
				tblKey:  "CH10",
				keyType: reflect.TypeOf(wantString),
			},
			want:    wantString,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getYangMdlKey(tt.args.prefix, tt.args.tblKey, tt.args.keyType)
			if (err != nil) != tt.wantErr {
				t.Errorf("getYangMdlKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getYangMdlKey() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApp_requestBodyHasField(t *testing.T) {
	type args struct {
		t reflect.Type
		v reflect.Value
	}

	var TestUint   uint32  = 123
	var existPtr *uint32   = &TestUint
	var nilPtr *int32
	var enumUnset   ocbinds.E_OpenconfigTerminalDevice_AdminStateType = ocbinds.OpenconfigTerminalDevice_AdminStateType_UNSET
	var enumValid   ocbinds.E_OpenconfigTerminalDevice_AdminStateType = ocbinds.OpenconfigTerminalDevice_AdminStateType_ENABLED

	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "ptr exist",
			args: args{
				t: reflect.TypeOf(existPtr),
				v: reflect.ValueOf(existPtr),
			},
			want: true,
		},
		{
			name: "ptr null",
			args: args{
				t: reflect.TypeOf(nilPtr),
				v: reflect.ValueOf(nilPtr),
			},
			want: false,
		},
		{
			name: "enum unset",
			args: args{
				t: reflect.TypeOf(enumUnset),
				v: reflect.ValueOf(enumUnset),
			},
			want: false,
		},
		{
			name: "enum valid",
			args: args{
				t: reflect.TypeOf(enumValid),
				v: reflect.ValueOf(enumValid),
			},
			want: true,
		},
		{
			name: "enum valid",
			args: args{
				t: reflect.TypeOf(enumValid),
				v: reflect.ValueOf(enumValid),
			},
			want: true,
		},
		{
			name: "uint32",
			args: args{
				t: reflect.TypeOf(TestUint),
				v: reflect.ValueOf(TestUint),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := requestBodyHasField(tt.args.t, tt.args.v); got != tt.want {
				t.Errorf("requestBodyHasField() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApp_constructRegexPathWithKey(t *testing.T) {
	type args struct {
		mdb    db.MDB
		num    db.DBNum
		path   string
		params *regexPathKeyParams
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "one key with prefix",
			args: args{
				mdb:   db.MDB{
					"host" : [db.MaxDB]*db.DB{},
				},
				num:   db.StateDB,
				path:   "/openconfig-terminal-device:terminal-device/logical-channels/channel/state",
				params: &regexPathKeyParams{
					tableName:    "LOGICAL_CHANNEL",
					listNodeName: []string{"channel"},
					keyName:      []string{"index"},
					redisPrefix:  []string{"CH"},
				},
			},
			want: []string{"/openconfig-terminal-device:terminal-device/logical-channels/channel[index=101]/state"},
		},
		{
			name: "one key without prefix",
			args: args{
				mdb:   db.MDB{
					"host" : [db.MaxDB]*db.DB{},
				},
				num:   db.StateDB,
				path:   "/openconfig-platform:components/component/state",
				params: &regexPathKeyParams{
					tableName:    "CHASSIS",
					listNodeName: []string{"component"},
					keyName:      []string{"name"},
					redisPrefix:  []string{""},
				},
			},
			want: []string{"/openconfig-platform:components/component[name=CHASSIS-1]/state"},
		},
		{
			name: "two keys",
			args: args{
				mdb:   db.MDB{
					"host" : [db.MaxDB]*db.DB{},
				},
				num:   db.StateDB,
				path:   "/openconfig-platform:components/component/openconfig-platform-transceiver:transceiver/physical-channels/channel/state",
				params: &regexPathKeyParams{
					tableName:    "TRANSCEIVER",
					listNodeName: []string{"component", "channel"},
					keyName:      []string{"name", "index"},
					redisPrefix:  []string{"", "CHANNEL-"},
				},
			},
			want: []string{"/openconfig-platform:components/component[name=TRANSCEIVER-1-1-C3]/openconfig-platform-transceiver:transceiver/physical-channels/channel[index=1]/state"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patches := gomonkey.ApplyMethod(reflect.TypeOf(&db.DB{}), "GetKeysPattern", func(_ *db.DB, ts *db.TableSpec, _ db.Key) ([]db.Key, error) {
				var keys []db.Key
				if ts.Name == "LOGICAL_CHANNEL" {
					keys = append(keys, asKey("CH101"))
				} else if ts.Name == "CHASSIS" {
					keys = append(keys, asKey("CHASSIS-1"))
				} else if ts.Name == "TRANSCEIVER" {
					keys = append(keys, asKey("TRANSCEIVER-1-1-C3", "CHANNEL-1"))
				}
				return keys, nil
			})
			defer patches.Reset()
			if got := constructRegexPathWithKey(tt.args.mdb, tt.args.num, tt.args.path, tt.args.params); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("constructRegexPathWithKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApp_needQuery(t *testing.T) {
	type args struct {
		targetPath string
		curNode    interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "parent need query",
			args: args{
				targetPath: "/openconfig-terminal-device:terminal-device/logical-channels/channel[index=101]/ethernet/lldp",
				curNode:    &ocbinds.OpenconfigTerminalDevice_TerminalDevice_LogicalChannels_Channel_Ethernet{},
			},
			want: true,
		},
		{
			name: "grandfather need query",
			args: args{
				targetPath: "/openconfig-terminal-device:terminal-device/logical-channels/channel[index=101]/ethernet/lldp",
				curNode:    &ocbinds.OpenconfigTerminalDevice_TerminalDevice_LogicalChannels_Channel{},
			},
			want: true,
		},
		{
			name: "parent does not need query",
			args: args{
				targetPath: "/openconfig-terminal-device:terminal-device/logical-channels/channel[index=101]/ethernet/lldp",
				curNode:    &ocbinds.OpenconfigTerminalDevice_TerminalDevice_LogicalChannels_Channel_State{},
			},
			want: false,
		},
		{
			name: "child need query",
			args: args{
				targetPath: "/openconfig-terminal-device:terminal-device/logical-channels/channel[index=101]/ethernet/lldp",
				curNode:    &ocbinds.OpenconfigTerminalDevice_TerminalDevice_LogicalChannels_Channel_Ethernet_Lldp_State{},
			},
			want: true,
		},
		{
			name: "sibling does not need query",
			args: args{
				targetPath: "/openconfig-terminal-device:terminal-device/logical-channels/channel[index=101]/ethernet/lldp",
				curNode:    &ocbinds.OpenconfigTerminalDevice_TerminalDevice_LogicalChannels_Channel_Ethernet_State{},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := needQuery(tt.args.targetPath, tt.args.curNode); got != tt.want {
				t.Errorf("needQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}