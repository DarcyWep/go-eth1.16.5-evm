package database

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"go-eth1.16.5-evm/config"
	"go-eth1.16.5-evm/core/rawdb"
	"go-eth1.16.5-evm/core/types"
	"go-eth1.16.5-evm/ethdb"
	"go-eth1.16.5-evm/ethdb/pebble"
)

func OpenDatabaseWithFreezer(ethConfig *config.EthConfig) (ethdb.Database, error) {
	db, err := pebble.New(defaultPebbleConfig.file, defaultPebbleConfig.cache, defaultPebbleConfig.handles, defaultPebbleConfig.namespace, defaultPebbleConfig.readonly)
	if err != nil {
		return nil, err
	}

	frdb, err := rawdb.Open(db, rawdb.OpenOptions{
		Ancient:          defaultRawConfig.ancient,
		Era:              defaultRawConfig.era,
		MetricsNamespace: defaultRawConfig.metricsNamespace,
		ReadOnly:         defaultRawConfig.readOnly,
	})
	return frdb, err
}

func GetBlockByNumber(db ethdb.Database, number *big.Int) (*types.Block, error) {
	var (
		block *types.Block
		err   error
	)
	hash := rawdb.ReadCanonicalHash(db, number.Uint64()) // 获取区块hash
	if (hash != common.Hash{}) {
		block = rawdb.ReadBlock(db, hash, number.Uint64())
		if block == nil {
			err = fmt.Errorf("read block(" + number.String() + ") error! block is nil")
		}
	} else {
		err = fmt.Errorf("read block(" + number.String() + ") error! hash is nil")
	}
	return block, err
}

func GetHeaderByNumber(db ethdb.Database, number uint64) (*types.Header, error) {
	var (
		header *types.Header = nil
		err    error         = nil
	)
	hash := rawdb.ReadCanonicalHash(db, number) // 创建StateDB
	if (hash != common.Hash{}) {
		if h := rawdb.ReadHeader(db, hash, number); header != nil {
			header = h
		} else {
			err = fmt.Errorf("create stateDB error! header is nil")
		}
	}
	return header, err
}
