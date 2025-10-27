package database

import (
	"path/filepath"
	"runtime"
)

var path string

func init() {
	if runtime.GOOS == "darwin" {
		path = "/Volumes/ETH_DATA/ethereum/geth/chaindata"
	} else {
		path = "/data/ethereum/execution/geth/chaindata"
	}

}

type pebbleConfig struct {
	file      string
	cache     int // MB
	handles   int
	namespace string
	readonly  bool
}

var defaultPebbleConfig = &pebbleConfig{
	file:      path,
	cache:     21462, // 如果内存较小，请修改
	handles:   524288,
	namespace: "eth/db/chaindata/",
	readonly:  true,
}

type rawConfig struct {
	ancient          string // ancients directory
	era              string // era files directory
	metricsNamespace string // prefix added to freezer metric names
	readOnly         bool
}

var defaultRawConfig = &rawConfig{
	ancient:          filepath.Join(path, "ancient"),
	era:              "",
	metricsNamespace: "eth/db/chaindata/",
	readOnly:         true,
}

type StateDBConfig struct {
	Cache     int
	Journal   string
	Preimages bool
}

func defaultStateDBConfig() *StateDBConfig {
	if runtime.GOOS == "darwin" { // MacOS
		return &StateDBConfig{
			Cache: 614,
			//Journal:   "/Users/darcywep/Projects/ethereum/execution/geth/triecache",
			Journal:   "/Volumes/ETH_DATA/ethereum/geth/triecache",
			Preimages: false,
		}
	} else {
		return &StateDBConfig{
			Cache:     614,
			Journal:   "/experiment/ethereum/geth/triecache",
			Preimages: false,
		}
	}
}
