package snmp

import (
	"github.com/Azure/sonic-mgmt-common/translib/db"
	"testing"
)

func TestApp_SendTrap(t *testing.T) {
	type args struct {
		alarmInfo db.Value
		isDel     bool
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "alarm created",
			args: args{
				isDel:     false,
			},
		},
		{
			name: "alarm cleared",
			args: args{
				isDel:     true,
			},
		},
	}
	for _, tt := range tests {
		alarmInfoMock := db.Value{
			Field: map[string]string{
				"id": "PORT-1-3-C3#OTN_LOFLOM",
				"resource": "PORT-1-3-C3",
				"severity": "MAJOR",
				"text": "OTUFlex Loss of frame/multi-frame",
				"time-created": "1652930670000000000",
				"type-id": "OTN_LOFLOM",
			},
		}
		t.Run(tt.name, func(t *testing.T) {
			sendTrap(alarmInfoMock, tt.args.isDel)
		})
	}
}
