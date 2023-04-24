package translib

import (
	"github.com/Azure/sonic-mgmt-common/translib/db"
	"github.com/agiledragon/gomonkey"
	"testing"
)

func TestApp_getOnChangePrecisePath(t *testing.T) {
	type args struct {
		path string
		key   *db.Key
	}
	tests := []struct {
		name string
		args args
		want *notificationInfo
	}{
		{
			name: "system alarm changed",
			args: args{
				path: "/openconfig-system:system/alarms/alarm/state",
				key: &db.Key{Comp: []string{"PORT-1-1-C3#RXRFI"}},
			},
			want: &notificationInfo{
				path: "/openconfig-system:system/alarms/alarm[id=PORT-1-1-C3#RXRFI]/state",
			},
		},
		{
			name: "terminal-device channel changed",
			args: args{
				path: "/openconfig-terminal-device:terminal-device/logical-channels/channel/state",
				key: &db.Key{Comp: []string{"CH101"}},
			},
			want: &notificationInfo{
				path: "/openconfig-terminal-device:terminal-device/logical-channels/channel[index=101]/state",
			},
		},
		{
			name: "platform component changed",
			args: args{
				path: "/openconfig-platform:components/component/state",
				key: &db.Key{Comp: []string{"FAN-1-9"}},
			},
			want: &notificationInfo{
				path: "/openconfig-platform:components/component[name=FAN-1-9]/state",
			},
		},
		{
			name: "platform transceiver channel changed",
			args: args{
				path: "/openconfig-platform:components/component/openconfig-platform-transceiver:transceiver/physical-channels/channel/state",
				key: &db.Key{Comp: []string{"TRANSCEIVER-1-1-C3", "CH-1"}},
			},
			want: &notificationInfo{
				path: "/openconfig-platform:components/component[name=TRANSCEIVER-1-1-C3]/openconfig-platform-transceiver:transceiver/physical-channels/channel[index=1]/state",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patches := gomonkey.ApplyFunc(getJson, func(nInfo *notificationInfo) ([]byte, error) {
				return nil, nil
			})
			defer patches.Reset()
			if got := getOnChangePrecisePath(tt.args.path, tt.args.key); len(got) != 0 {
				if got != tt.want.path {
					t.Errorf("getOnChangePrecisePath() = %v, want.path %v", got, tt.want.path)
				}
			}
		})
	}
}
