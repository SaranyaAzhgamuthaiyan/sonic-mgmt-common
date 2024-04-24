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
