package orc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	FrameTickCount    int     `json:"FrameTickCount"`
	PhysicsTickCount  int     `json:"PhysicsTickCount"`
	Accel             float64 `json:"Accel"`
	MaxSpeed          float64 `json:"MaxSpeed"`
	PlayerRadius      float64 `json:"PlayerRadius"`
	ProjectileRadius  float64 `json:"ProjectileRadius"`
	ProjectileSpeed   float64 `json:"ProjectileSpeed"`
	ServerPort        int     `json:"ServerPort"`
	PlayerAttackRange int     `json:"PlayerAttackRange"`
}

var GlobalConfig Config

func init() {
	b, err := ioutil.ReadFile("./config.json")
	if err != nil {
		fmt.Println("cannot read json file ./config.json")
		return
	}
	err = json.Unmarshal(b, &GlobalConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
}
