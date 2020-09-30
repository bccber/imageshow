package conf

import (
	"encoding/json"
	"io/ioutil"
)

var Config *conf

type conf struct {
	SpiderDB string   `json:"spider_db"`
	MasterDB []string `json:"master_db"`
	SlaveDB  []string `json:"slave_db"`
}

func init() {
	buf, err := ioutil.ReadFile("./conf/conf.json")
	if err != nil {
		panic("conf.json 有误")
	}

	if err = json.Unmarshal(buf, &Config); err != nil {
		panic(err)
	}
}
