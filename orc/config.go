package orc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	Tickcount int `json:"tickcount"`
}

var GlobalConfig Config

func init() {
	b, err := ioutil.ReadFile("./config/config.json")
	if err != nil {
		fmt.Println("cannot read json file ./config/config.json")
		return
	}
	err = json.Unmarshal(b, &GlobalConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
}
