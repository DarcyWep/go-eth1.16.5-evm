package main

import (
	"fmt"

	ethereumethdb "github.com/ethereum/go-ethereum/ethdb"
	"go-eth1.16.5-evm/config"
	"go-eth1.16.5-evm/database"
	//"go-eth1.16.5-evm/config"
	//"go-eth1.16.5-evm/core"
	//"go-eth1.16.5-evm/database"
)

func main() {
	frdb, err := database.OpenDatabaseWithFreezer(&config.DefaultsEthConfig)
	if err != nil {
		fmt.Println("OpenDatabaseWithFreezer err", err.Error())
		return
	}
	if _, ok := frdb.(ethereumethdb.Database); ok {
		fmt.Println("assert success")
	} else {
		fmt.Println("assert fail")
	}

	//bcHc, err := core.NewHeaderChain(frdb.(ethereumethdb.Database), chainConfig, engine, bc.insertStopped)
	//bc.processor = NewStateProcessor(bc.hc)
	//fmt.Printf("Hello and welcome, %s!\n", s)
	//
	//for i := 1; i <= 5; i++ {
	//	fmt.Println("i =", 100/i)
	//}
}
