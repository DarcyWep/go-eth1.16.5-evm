package config

import (
	"go-eth1.16.5-evm/core/vm"
)

//{Tracer:<nil> NoBaseFee:false EnablePreimageRecording:false ExtraEips:[] StatelessSelfValidation:false EnableWitnessStats:false

var DefaultVmConfig = vm.Config{
	Tracer:                  nil,
	NoBaseFee:               false,
	EnablePreimageRecording: false,
	ExtraEips:               make([]int, 0),
	StatelessSelfValidation: false,
	EnableWitnessStats:      false,
}
