package main

import (
	"flag"

	"go-eth1.16.5-evm/replay/replay_geth"
	"go-eth1.16.5-evm/replay/replay_gethcopy"
)

func main() {
	mode := flag.String("mode", "", "replay mode: geth or copy")
	flag.Parse()

	if *mode == "geth" {
		replay_geth.ReplayGeth()
	} else {
		replay_gethcopy.ReplayCopy()
	}
}
