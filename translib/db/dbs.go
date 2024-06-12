package db

import (
	"bufio"
	//"errors"
	//cmn "github.com/Azure/sonic-mgmt-common/cvl/common"
	//"github.com/go-redis/redis/v7"
	"github.com/golang/glog"
	"os"
	"strconv"
	"strings"
)

func isMultiAsic() bool {
	return NumAsic > 1
}

func getNumAsic() int {
	file, err := os.Open(DefaultAsicConfFilePath)
	if err != nil {
		glog.Warning("Cannot find the asic.conf file, set num_asic to 1 by default")
		return 1
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var text []string
	for scanner.Scan() {
		text = append(text, scanner.Text())
	}

	var numAsics = 1
	for _, line := range text {
		tokens := strings.Split(line, "=")
		if len(tokens) != 2 {
			continue
		}

		if strings.ToLower(tokens[0]) == "num_asic" {
			num, err := strconv.Atoi(tokens[1])
			if err != nil {
				glog.Warning("invalid line in asic.config")
				continue
			}
			numAsics = num
		}
	}

	return numAsics
}

func GetMultiDbNames() []string {
	dbNames := []string{"host"}
	if isMultiAsic() {
		for num := 0; num < NumAsic; num++ {
			dbNames = append(dbNames, "asic"+strconv.Itoa(num))
		}
	}

	return dbNames
}

func getDbNameBySlotNum(slotNum int) string {
	if NumAsic == 1 || slotNum > NumAsic || slotNum < 1 {
		return "host"
	}

	return "asic" + strconv.Itoa(slotNum-1)
}

func GetMDBNameFromEntity(entity interface{}) string {
	var slotNum int
	var slotNumStr string
	switch t := entity.(type) {
	case uint16:
		valUint16, _ := entity.(uint16)
		slotNum = int(valUint16 / 100)
	case uint32:
		//ch115
		valUint32, _ := entity.(uint32)
		slotNum = int(valUint32 / 100)
	case string:
		valString, _ := entity.(string)
		elmts := strings.Split(valString, "-")
		if len(elmts) == 2 && elmts[0] == "SLOT" {
			// used for reboot entity-name
			slotNumStr = elmts[1]
		} else if len(elmts) >= 3 {
			slotNumStr = elmts[2]
		} else {
			goto error
		}

		tmp, err := strconv.ParseInt(slotNumStr, 10, 64)
		if err != nil {
			goto error
		}
		slotNum = int(tmp)
	default:
		glog.Errorf("unexpected type %T", t)
	}

error:
	return getDbNameBySlotNum(slotNum)
}
