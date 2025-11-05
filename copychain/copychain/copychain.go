package copychain

import (
	"fmt"
	"math/big"
	"path/filepath"

	"github.com/ethereum/go-ethereum/rlp"
	"go-eth1.16.5-evm/copychain/copychain_common"
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
		File:      "/data/ethereum/state_snapshot/chaindata_2100/geth/chaindata",
		Cache:     21462, // 如果内存较小，请修改
		Handles:   524288,
		Namespace: "eth/db/chaindata/",
		Readonly:  false,
	}
	rawConfig := &database.RawConfig{
		Ancient:          filepath.Join(pebbleConfig.File, "ancient"),
		Era:              "",
		MetricsNamespace: "eth/db/chaindata/",
		ReadOnly:         false,
	}
	newFrdb, err := database.OpenDatabaseWithFreezer(pebbleConfig, rawConfig)
	if err != nil {
		fmt.Printf("OpenDatabaseWithFreezer err:" + err.Error())
		return
	}
	defer newFrdb.Close()
	var (
		blocks   types.Blocks
		receipts []rlp.RawValue
	)
	block, receipt, err := readBlock(oldFrdb, new(big.Int).SetInt64(0))
	if err != nil {
		fmt.Printf("wrtieBlock err:" + err.Error())
		return
	}
	blocks = append(blocks, block)
	receipts = append(receipts, receipt)

	for blockNumber := copychain_common.StartBlockNumber; blockNumber.Cmp(copychain_common.FinishBlockNumber) == -1; blockNumber = blockNumber.Add(blockNumber, copychain_common.AddSpan) {
		block, receipt, err = readBlock(oldFrdb, blockNumber)
		if err != nil {
			fmt.Printf("wrtieBlock err:" + err.Error())
			return
		}
		blocks = append(blocks, block)
		receipts = append(receipts, receipt)
		if blockNumber.Uint64()%100000 == 0 {
			size, err := writeBlocks(newFrdb, blocks, receipts)
			if err != nil {
				fmt.Printf("writeBlocks err:" + err.Error())
				return
			}
			fmt.Printf("finish copy block, number=%d, size=%fMB\n", blockNumber.Uint64(), size)
			blocks = blocks[:0]
			receipts = receipts[:0]
		}
	}
}
