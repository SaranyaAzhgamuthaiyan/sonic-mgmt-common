package snmp

import (
	"fmt"
	"github.com/Azure/sonic-mgmt-common/translib/db"
	"github.com/golang/glog"
	g "github.com/gosnmp/gosnmp"
	"strconv"
	"strings"
	"time"
)

type server struct {
	target string
	targetPort uint16
	localAddr string
	localAddrPort uint16
}

type snmpClient struct {
	client g.GoSNMP
}

var alarm2Trap = []struct{
	trapField string
	alarmField string
}{
	{".1.3.6.1.4.1.57874.1.8.1.0","resource"},
	{".1.3.6.1.4.1.57874.1.8.2.0","type-id"},
	{".1.3.6.1.4.1.57874.1.8.3.0","severity"},
	{".1.3.6.1.4.1.57874.1.8.4.0","text"},
	{".1.3.6.1.4.1.57874.1.8.5.0","time-created"},
	{".1.3.6.1.4.1.57874.1.8.6.0","time-send"},
	{".1.3.6.1.4.1.57874.1.8.7.0","alarm-number"},
	{".1.3.6.1.4.1.57874.1.8.8.0","reserved"},
}

func newSnmpClientV2(s server) *snmpClient {
	c := &snmpClient{
		g.GoSNMP{
			Target:             s.target,
			Port:               s.targetPort,
			Community:          "public",
			Version:            g.Version2c,
			Timeout:            time.Duration(2) * time.Second,
			Retries:            3,
			ExponentialTimeout: true,
			MaxOids:            g.MaxOids,
		},
	}

	if len(s.localAddr) != 0 {
		c.client.LocalAddr = fmt.Sprintf("%s:%d", s.localAddr, s.localAddrPort)
	}

	return c
}

func constructTrap(alarmInfo db.Value, isDel bool) (g.SnmpTrap, error) {
	var trap g.SnmpTrap
	var pdu g.SnmpPDU

	pdu.Name = ".1.3.6.1.6.3.1.1.4.1.0"
	pdu.Type = g.ObjectIdentifier
	pdu.Value = ".1.3.6.1.4.1.57874.1.8.7"
	trap.Variables = append(trap.Variables, pdu)

	// alarm info -> trap info
	alarmInfo.Set("time-send", time.Now().Format("2006-01-02 15:04:05.000"))

	timeCreatedField := "time-created"
	if isDel {
		alarmInfo.Set("severity", "CLEAR")
		timeCreatedField = "time-cleared"
	}
	timeInt64, err := strconv.ParseInt(alarmInfo.Get(timeCreatedField), 10, 64)
	if err != nil {
		glog.Errorf("get alarm time-created field failed")
		return trap, err
	}
	timeFmt := time.Unix(0, timeInt64).Format("2006-01-02 15:04:05.000")
	alarmInfo.Set("time-created", timeFmt)

	for _, m := range alarm2Trap {
		v := alarmInfo.Get(m.alarmField)
		if len(v) == 0 {
			continue
		}

		pdu.Name = m.trapField
		pdu.Type = g.OctetString
		pdu.Value = v
		trap.Variables = append(trap.Variables, pdu)
	}

	return trap, nil
}

func sendTrap(s server, trap g.SnmpTrap) {
	c := newSnmpClientV2(s)
	err := c.client.Connect()
	if err != nil {
		glog.Errorf("Connect() err: %v", err)
		return
	}
	defer c.client.Conn.Close()

	_, err = c.client.SendTrap(trap)
	if err != nil {
		glog.Errorf("sendTrap() err: %v", err)
		return
	}

	glog.V(1).Infof("send %s:%d the trap %v", s.target, s.targetPort, trap)
}

func getTrapServer() []server {
	var servers []server
	var loopbackIfIp string

	d, err := db.NewDBForMultiAsic(db.GetDBOptions(db.ConfigDB, true), "host")
	if err != nil {
		glog.V(1).Infof("connect host config_db failed when to get snmp trap server")
		return servers
	}
	defer d.DeleteDB()

	ts := &db.TableSpec{Name: "LOOPBACK_INTERFACE"}
	keys, _ := d.GetKeys(ts)
	for _, key := range keys {
		if key.Get(0) != "Loopback0" || key.Len() != 2 {
			continue
		}

		elmts := strings.Split(key.Get(1), "/")
		loopbackIfIp = elmts[0]
		glog.V(1).Infof("snmp trap has loopback-interface %s", loopbackIfIp)
	}

	ts = &db.TableSpec{Name: "SNMP_TRAP_SERVER"}
	keys, _ = d.GetKeys(ts)
	if len(keys) == 0 {
		glog.V(1).Infof("no snmp trap servers")
		return servers
	}

	for _, iter := range keys {
		var s server
		data, _ := d.GetEntry(ts, iter)
		s.target = data.Get("address")
		p, _ := strconv.ParseUint(data.Get("port"), 10, 16)
		s.targetPort = uint16(p)

		if len(loopbackIfIp) != 0 {
			s.localAddr = loopbackIfIp
		}

		servers = append(servers, s)
	}

	return servers
}

func convertAlarmAndSendTrap(alarm db.Value, isDel bool) error {
	trap, err := constructTrap(alarm, isDel)
	if err != nil {
		return err
	}
	servers := getTrapServer()
	for _, server := range servers {
		sendTrap(server, trap)
		if isDel {
			glog.Infof("send %s:%d the cleared alarm %s", server.target, server.targetPort, alarm.Get("id"))
		} else {
			glog.Infof("send %s:%d the created alarm %s", server.target, server.targetPort, alarm.Get("id"))
		}
	}

	return nil
}

func TrapHandler(d *db.DB, sKey *db.SKey, key *db.Key, event db.SEvent) error {
	if event != db.SEventHSet && event != db.SEventHDel {
		return nil
	}

	glog.V(1).Infof("trap handler start process")

	isDel := false
	if sKey.Ts.Name == "HISALARM" {
		isDel = true
	}

	alarm, err := d.GetEntry(sKey.Ts, *key)
	if err != nil {
		return err
	}
	err = convertAlarmAndSendTrap(alarm, isDel)
	if err != nil {
		return err
	}

	return nil
}

func SubscribeDbForSnmpTrap(dbNo db.DBNum, ts *db.TableSpec, key db.Key) error {
	dbs, err := db.NewNumberDB(dbNo, true)
	if err != nil {
		return err
	}

	sKeyList := []*db.SKey {
		{Ts: ts, Key: &key},
	}

	for name, d := range dbs {
		err = db.SubscribeDB(d, sKeyList, TrapHandler)
		if err != nil {
			return err
		}

		if ts.Name != "CURALARM" {
			continue
		}

		oneDbKeys, err := d.GetKeys(ts)
		if err != nil {
			glog.Errorf("get ALARM keys from %s failed as %v", name, err)
			return err
		}
		if len(oneDbKeys) == 0 {
			continue
		}

		for _, k := range oneDbKeys {
			alarm, err1 := d.GetEntry(ts, k)
			if err1 != nil {
				continue
			}
			err1 = convertAlarmAndSendTrap(alarm, false)
			if err1 != nil {
				return err1
			}
		}
	}

	return nil
}