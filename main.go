package main

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"go-eth1.16.5-evm/config"
	"go-eth1.16.5-evm/core"
	"go-eth1.16.5-evm/database"
	"go-eth1.16.5-evm/ethdb"
)

var (
	rootBlockNumber   *big.Int = new(big.Int).SetInt64(21_000_000)
	startBlockNumber  *big.Int = new(big.Int).SetInt64(21_000_001)
	finishBlockNumber *big.Int = new(big.Int).SetInt64(21_000_005)
	addSpan           *big.Int = new(big.Int).SetInt64(1)
	parentStateRoot   common.Hash
)

//var (
//	rootBlockNumber   *big.Int = new(big.Int).SetInt64(2380000)
//	startBlockNumber  *big.Int = new(big.Int).SetInt64(2380001)
//	finishBlockNumber *big.Int = new(big.Int).SetInt64(2380005)
//	addSpan           *big.Int = new(big.Int).SetInt64(1)
//	parentStateRoot   common.Hash
//)

func newProcessor() (*core.StateProcessor, ethdb.Database, error) {
	frdb, err := database.OpenDatabaseWithFreezer(database.DefaultPebbleConfig, database.DefaultRawConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("OpenDatabaseWithFreezer err:" + err.Error())
	}
	engine, err := config.CreateConsensusEngine(config.MainnetChainConfig)
	if err != nil {
		_ = frdb.Close()
		return nil, nil, fmt.Errorf("CreateConsensusEngine err:" + err.Error())
	}

	bcHc, err := core.NewHeaderChain(frdb, config.MainnetChainConfig, engine)
	processor := core.NewStateProcessor(bcHc)
	return processor, frdb, nil
}

func main() {
	processor, frdb, err := newProcessor()
	if err != nil {
		panic(err)
		return
	}
	defer frdb.Close()
	blockPre, err := database.GetBlockByNumber(frdb, rootBlockNumber)
	if err != nil {
		panic(err)
		return
	}
	parentStateRoot = blockPre.Root()
	alldbForState, err := database.NewAllDBForState(blockPre.Number(), blockPre.Root(), false, false)
	defer alldbForState.Close()

	for blockNumber := startBlockNumber; blockNumber.Cmp(finishBlockNumber) == -1; blockNumber = blockNumber.Add(blockNumber, addSpan) {
		err := alldbForState.UpdateStateDB(parentStateRoot)
		if err != nil {
			panic(err)
			return
		}
		block, err := database.GetBlockByNumber(frdb, blockNumber)
		if err != nil {
			panic(err)
			return
		}
		fmt.Println(block.Number(), block.Transactions())
		_, err = processor.Process(block, alldbForState.StateDB, config.DefaultVmConfig)
		if err != nil {
			fmt.Println(err)
		}

		// Commit all cached state changes into underlying memory database.
		root, _, err := alldbForState.StateDB.CommitWithUpdate(block.NumberU64(), config.MainnetChainConfig.IsEIP158(block.Number()), config.MainnetChainConfig.IsCancun(block.Number(), block.Time()))
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("blockNumber="+blockNumber.String()+"\t process state root:", block.Root())
		fmt.Println("blockNumber="+blockNumber.String()+"\t block state root:", root)
		parentStateRoot = root
	}

	//fmt.Printf("Hello and welcome, %s!\n", s)
	//
	//for i := 1; i <= 5; i++ {
	//	fmt.Println("i =", 100/i)
	//}
}
