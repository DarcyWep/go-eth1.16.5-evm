package database

import (
	"path/filepath"
	"runtime"

	"github.com/ethereum/go-ethereum/core"
	"go-eth1.16.5-evm/core/rawdb"
	"go-eth1.16.5-evm/triedb"
	"go-eth1.16.5-evm/triedb/hashdb"
	"go-eth1.16.5-evm/triedb/pathdb"
)

type PebbleConfig struct {
	File      string
	Cache     int // MB
	Handles   int
	Namespace string
	Readonly  bool
}

type RawConfig struct {
	Ancient          string // ancients directory
	Era              string // era files directory
	MetricsNamespace string // prefix added to freezer metric names
	ReadOnly         bool
}

type stateDBConfig struct {
	path    string
	cache   int
	handles int
}

var path string
var DefaultPebbleConfig *PebbleConfig
var DefaultRawConfig *RawConfig
var defaultStateDBConfig *stateDBConfig

func init() {
	if runtime.GOOS == "darwin" {
		path = "/Volumes/ETH_DATA/ethereum/geth/chaindata"
	} else {
		path = "/data/ethereum/execution/geth/chaindata"
		//path = "/root/ethereum/execution/geth/chaindata"
	}
	DefaultPebbleConfig = &PebbleConfig{
		File:      path,
		Cache:     21462, // 如果内存较小，请修改
		Handles:   524288,
		Namespace: "eth/db/chaindata/",
		Readonly:  true,
	}

	DefaultRawConfig = &RawConfig{
		Ancient:          filepath.Join(path, "ancient"),
		Era:              "",
		MetricsNamespace: "eth/db/chaindata/",
		ReadOnly:         true,
	}

	defaultStateDBConfig = &stateDBConfig{
		path:    "/data/ethereum/state_snapshot",
		cache:   32768,
		handles: 32768,
	}

}

func trieDBConfig(blockChainConfig *core.BlockChainConfig, isVerkle bool) *triedb.Config {
	config := &triedb.Config{
		Preimages: blockChainConfig.Preimages,
		IsVerkle:  isVerkle,
	}
	if blockChainConfig.StateScheme == rawdb.HashScheme {
		config.HashDB = &hashdb.Config{
			CleanCacheSize: blockChainConfig.TrieCleanLimit * 1024 * 1024,
		}
	}
	if blockChainConfig.StateScheme == rawdb.PathScheme {
		config.PathDB = &pathdb.Config{
			StateHistory:        blockChainConfig.StateHistory,
			EnableStateIndexing: blockChainConfig.ArchiveMode,
			TrieCleanSize:       blockChainConfig.TrieCleanLimit * 1024 * 1024,
			StateCleanSize:      blockChainConfig.SnapshotLimit * 1024 * 1024,
			JournalDirectory:    blockChainConfig.TrieJournalDirectory,

			// TODO(rjl493456442): The write buffer represents the memory limit used
			// for flushing both trie data and state data to disk. The config name
			// should be updated to eliminate the confusion.
			WriteBufferSize: blockChainConfig.TrieDirtyLimit * 1024 * 1024,
			NoAsyncFlush:    blockChainConfig.TrieNoAsyncFlush,
		}
	}
	return config
}
