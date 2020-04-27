package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateIscn{}, "likechain/MsgCreateISCN", nil)
	cdc.RegisterConcrete(MsgAddEntity{}, "likechain/MsgAddEntity", nil)
	cdc.RegisterInterface((*EntityInput)(nil), nil)
	cdc.RegisterConcrete(Entity{}, "likechain/iscn/Entity", nil)
	cdc.RegisterConcrete(EntityByCid{}, "likechain/iscn/EntityByCid", nil)
	cdc.RegisterInterface((*RightTermsInput)(nil), nil)
	cdc.RegisterConcrete(RightTerms{}, "likechain/iscn/RightTerms", nil)
	cdc.RegisterConcrete(RightTermsByCid{}, "likechain/iscn/RightTermsByCid", nil)
}

var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
