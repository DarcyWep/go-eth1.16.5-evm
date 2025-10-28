package main

import (
	"fmt"

	"github.com/prometheus/procfs/blockdevice"
	"go-eth1.16.5-evm/config"
	"go-eth1.16.5-evm/core"
	"go-eth1.16.5-evm/database"
)

func main() {
	frdb, err := database.OpenDatabaseWithFreezer(&config.DefaultsEthConfig)
	if err != nil {
		fmt.Println("OpenDatabaseWithFreezer err", err.Error())
		return
	}
	engine, err := config.CreateConsensusEngine(config.MainnetChainConfig)
	if err != nil {
		fmt.Println("CreateConsensusEngine err", err.Error())
		return
	}

	bcHc, err := core.NewHeaderChain(frdb, config.MainnetChainConfig, engine)
	processor := core.NewStateProcessor(bcHc)
	processor.Process()
	//fmt.Printf("Hello and welcome, %s!\n", s)
	//
	//for i := 1; i <= 5; i++ {
	//	fmt.Println("i =", 100/i)
	//}
}
