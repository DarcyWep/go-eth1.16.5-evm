package copychain

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/rlp"
	"go-eth1.16.5-evm/copychain/copychain_config"
	"go-eth1.16.5-evm/core/rawdb"
	"go-eth1.16.5-evm/core/types"
	"go-eth1.16.5-evm/database"
	"go-eth1.16.5-evm/ethdb"
)

func readBlock(oldFrdb ethdb.Database, blockNumber *big.Int) (*types.Block, rlp.RawValue, error) {
	block, err := database.GetBlockByNumber(oldFrdb, blockNumber)
	if err != nil {
		return nil, nil, fmt.Errorf("GetBlockByNumber err:" + err.Error())
	}
	return block, []byte(""), nil
}

func writeBlocks(newFrdb ethdb.Database, blocks types.Blocks, receipts []rlp.RawValue) (float64, error) {
	var size float64 = 0
	// Write all chain data to ancients.
	writeSize, err := rawdb.WriteAncientBlocks(newFrdb, blocks, receipts)
	if err != nil {
		return 0, fmt.Errorf("Error importing chain data to ancients " + err.Error())
	}
	size = size + (float64(writeSize) / 1024 / 1024) // 取MB

	// Sync the ancient store explicitly to ensure all data has been flushed to disk.
	if err := newFrdb.SyncAncient(); err != nil {
		return 0, err
	}
	// Write hash to number mappings
	batch := newFrdb.NewBatch()
	for _, block := range blocks {
		rawdb.WriteHeaderNumber(batch, block.Hash(), block.NumberU64())
	}
	if err := batch.Write(); err != nil {
		return 0, err
	}
	return size, nil
}

func CopyChain() {
	oldFrdb, err := database.OpenDatabaseWithFreezer(database.DefaultPebbleConfig, database.DefaultRawConfig)
	if err != nil {
		fmt.Printf("OpenDatabaseWithFreezer err:" + err.Error())
		return
	}
	defer oldFrdb.Close()
	pebbleConfig := &database.PebbleConfig{
		File:      "/data/ethereum/state_snapshot/chaindata_2100/",
		Cache:     21462, // 如果内存较小，请修改
		Handles:   524288,
		Namespace: "eth/db/chaindata/",
		Readonly:  false,
	}
	newFrdb, err := database.OpenDatabase(pebbleConfig)
	if err != nil {
		fmt.Printf("OpenDatabase err:" + err.Error())
		return
	}
	defer newFrdb.Close()

	// 复制创世块
	block, _, err := readBlock(oldFrdb, new(big.Int).SetInt64(0))
	if err != nil {
		fmt.Printf("wrtieBlock err:" + err.Error())
		return
	}
	rawdb.WriteBlock(newFrdb, block)
	rawdb.WriteCanonicalHash(newFrdb, block.Hash(), block.NumberU64())
	fmt.Printf("finish copy block, number=%d\n", 0)

	// 复制后续区块
	for blockNumber := copychain_config.StartBlockNumber; blockNumber.Cmp(copychain_config.FinishBlockNumber) == -1; blockNumber = blockNumber.Add(blockNumber, copychain_config.AddSpan) {
		block, _, err = readBlock(oldFrdb, blockNumber)
		if err != nil {
			fmt.Printf("readBlock err:" + err.Error())
			return
		}
		rawdb.WriteBlock(newFrdb, block)
		rawdb.WriteCanonicalHash(newFrdb, block.Hash(), block.NumberU64())

		// 每10万个区块写入一次
		if blockNumber.Uint64()%100000 == 0 {
			fmt.Printf("finish copy block, number=%d\n", blockNumber.Uint64())
		}
	}
}
