package types

import (
	"github.com/multiformats/go-multibase"
)

var (
	// TODO: plug real codec
	IscnKernelCodecType   = uint64(0xc0de1337)
	IscnContentCodecType  = uint64(0x1337c0de)
	EntityCodecType       = uint64(0x13371337)
	RightTermsCodecType   = uint64(0xc0dec0de)
	RightsCodecType       = uint64(0xdeadbeef)
	StakeholdersCodecType = uint64(0xbeefdead)
)

var CidMbaseEncoder multibase.Encoder

func init() {
	var err error
	CidMbaseEncoder, err = multibase.EncoderByName("base58btc")
	if err != nil {
		panic(err)
	}
}
