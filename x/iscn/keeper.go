package iscn

import (
	"encoding/binary"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/tendermint/tendermint/crypto/tmhash"

	gocid "github.com/ipfs/go-cid"
	cbornode "github.com/ipfs/go-ipld-cbor"
	mh "github.com/multiformats/go-multihash"
)

const (
	DefaultParamspace = ModuleName
)

type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) auth.Account
}

type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           *codec.Codec
	paramstore    params.Subspace
	codespace     sdk.CodespaceType
	accountKeeper AccountKeeper
	supplyKeeper  authTypes.SupplyKeeper
}

func NewKeeper(
	cdc *codec.Codec, key sdk.StoreKey, accountKeeper AccountKeeper,
	supplyKeeper authTypes.SupplyKeeper, paramstore params.Subspace,
	codespace sdk.CodespaceType,
) Keeper {
	return Keeper{
		storeKey:      key,
		cdc:           cdc,
		paramstore:    paramstore.WithKeyTable(ParamKeyTable()),
		codespace:     codespace,
		accountKeeper: accountKeeper,
		supplyKeeper:  supplyKeeper,
	}
}

func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}

func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

func (k Keeper) FeePerByte(ctx sdk.Context) (res sdk.DecCoin) {
	k.paramstore.Get(ctx, KeyFeePerByte, &res)
	return
}

func (k Keeper) GetParams(ctx sdk.Context) Params {
	return Params{
		FeePerByte: k.FeePerByte(ctx),
	}
}

func (k Keeper) SetParams(ctx sdk.Context, params Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

func (k Keeper) GetCidBlock(ctx sdk.Context, cid CID) []byte {
	key := GetCidBlockKey(cid.Bytes())
	return ctx.KVStore(k.storeKey).Get(key)
}

func (k Keeper) HasCidBlock(ctx sdk.Context, cid CID) bool {
	key := GetCidBlockKey(cid.Bytes())
	return ctx.KVStore(k.storeKey).Has(key)
}

func (k Keeper) SetCidBlock(ctx sdk.Context, cid CID, bz []byte) {
	key := GetCidBlockKey(cid.Bytes())
	ctx.KVStore(k.storeKey).Set(key, bz)
}

func (k Keeper) GetCidIscnData(ctx sdk.Context, cid CID) *IscnData {
	bz := k.GetCidBlock(ctx, cid)
	if bz == nil {
		return nil
	}
	data := IscnData{}
	// TODO: deserialize by IPFS block
	err := cbornode.DecodeInto(bz, &data)
	if err != nil {
		return nil
	}
	return &data
}

func (k Keeper) SetCidIscnData(ctx sdk.Context, data interface{}, codec uint64) (*CID, error) {
	// TODO: serialize by IPFS block
	bz, err := cbornode.DumpObject(data)
	if err != nil {
		return nil, err
	}
	// TODO: compute CID by IPFS block
	cid, err := gocid.V1Builder{
		Codec:  codec,
		MhType: mh.SHA2_256,
	}.Sum(bz)
	if err != nil {
		return nil, err
	}

	k.SetCidBlock(ctx, cid, bz)
	return &cid, nil
}

func (k Keeper) checkCodecAndGetIscnData(ctx sdk.Context, cid CID, codec uint64) *IscnData {
	if cid.Prefix().GetCodec() != codec {
		return nil
	}
	return k.GetCidIscnData(ctx, cid)
}

func (k Keeper) setIscnDataAndEmitEvent(
	ctx sdk.Context, data interface{}, codec uint64,
	eventType string, attr string,
) (*CID, error) {
	cid, err := k.SetCidIscnData(ctx, data, codec)
	if err != nil {
		return nil, err
	}
	cidStr, err := cid.StringOfBase(CidMbaseEncoder.Encoding())
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			eventType,
			sdk.NewAttribute(attr, cidStr),
		),
	)
	return cid, nil
}

func (k Keeper) GetEntity(ctx sdk.Context, cid CID) *IscnData {
	return k.checkCodecAndGetIscnData(ctx, cid, EntityCodecType)
}

func (k Keeper) SetEntity(ctx sdk.Context, data IscnData) (*CID, error) {
	return k.setIscnDataAndEmitEvent(
		ctx, data, EntityCodecType, EventTypeAddEntity, AttributeKeyEntityCid,
	)
}

func (k Keeper) GetRightTerms(ctx sdk.Context, cid CID) *string {
	bz := k.GetCidBlock(ctx, cid)
	if bz == nil {
		return nil
	}
	terms := string(bz)
	return &terms
}

func (k Keeper) SetRightTerms(ctx sdk.Context, terms string) (*CID, error) {
	bz := []byte(terms)
	cid, err := gocid.V1Builder{
		Codec:  RightTermsCodecType,
		MhType: mh.SHA2_256,
	}.Sum(bz)
	if err != nil {
		return nil, err
	}
	k.SetCidBlock(ctx, cid, bz)
	cidStr, err := cid.StringOfBase(CidMbaseEncoder.Encoding())
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			EventTypeAddRightTerms,
			sdk.NewAttribute(AttributeKeyRightTermsCid, cidStr),
		),
	)
	return &cid, nil
}

func (k Keeper) GetIscnContent(ctx sdk.Context, cid CID) *IscnData {
	return k.checkCodecAndGetIscnData(ctx, cid, IscnContentCodecType)
}

func (k Keeper) SetIscnContent(ctx sdk.Context, data IscnData) (*CID, error) {
	return k.setIscnDataAndEmitEvent(
		ctx, data, IscnContentCodecType, EventTypeAddIscnContent, AttributeKeyIscnContentCid,
	)
}

func (k Keeper) GetIscnKernelByCID(ctx sdk.Context, cid CID) *IscnData {
	return k.checkCodecAndGetIscnData(ctx, cid, IscnKernelCodecType)
}

func (k Keeper) GetIscnKernelCIDByIscnID(ctx sdk.Context, iscnID IscnID) *CID {
	key := GetIscnKernelKey(iscnID)
	cidBytes := ctx.KVStore(k.storeKey).Get(key)
	if cidBytes == nil {
		return nil
	}
	_, cid, err := gocid.CidFromBytes(cidBytes)
	if err != nil {
		// TODO: should panic or at least log
		return nil
	}
	return &cid
}

func (k Keeper) SetIscnKernel(ctx sdk.Context, iscnID IscnID, data IscnData) (*CID, error) {
	cid, err := k.setIscnDataAndEmitEvent(
		ctx, data, IscnKernelCodecType, EventTypeAddIscnKernel, AttributeKeyIscnKernelCid,
	)
	if err != nil {
		return nil, err
	}
	cidBytes := cid.Bytes()
	key := GetIscnKernelKey(iscnID)
	ctx.KVStore(k.storeKey).Set(key, cidBytes)
	key = GetCidToIscnIDKey(cidBytes)
	ctx.KVStore(k.storeKey).Set(key, iscnID)
	return cid, err
}

func (k Keeper) SetIscnCount(ctx sdk.Context, count uint64) {
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(count)
	ctx.KVStore(k.storeKey).Set(IscnCountKey, bz)
}

func (k Keeper) GetIscnCount(ctx sdk.Context) uint64 {
	count := uint64(0)
	bz := ctx.KVStore(k.storeKey).Get(IscnCountKey)
	if bz != nil {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &count)
	}
	return count
}

func (k Keeper) DeductFeeForIscn(ctx sdk.Context, feePayer sdk.AccAddress, tx []byte) error {
	acc := k.accountKeeper.GetAccount(ctx, feePayer)
	if acc == nil {
		return fmt.Errorf("No account") // TODO: proper error
	}
	feePerByte := k.GetParams(ctx).FeePerByte
	feeAmount := feePerByte.Amount.MulInt64(int64(len(tx)))
	fees := sdk.NewCoins(sdk.NewCoin(feePerByte.Denom, feeAmount.Ceil().RoundInt()))
	result := auth.DeductFees(k.supplyKeeper, ctx, acc, fees)
	if !result.IsOK() {
		// TODO: proper error
		return fmt.Errorf("Not enough fee, %s %s needed", feeAmount.String(), feePerByte.Denom)
	}
	return nil
}

func (k Keeper) AddIscnKernel(ctx sdk.Context, kernel IscnData) (iscnID IscnID, err error) {
	hasher := tmhash.New()
	hasher.Write(ctx.BlockHeader().LastBlockId.Hash)
	iscnCount := k.GetIscnCount(ctx)
	k.SetIscnCount(ctx, iscnCount+1)
	binary.Write(hasher, binary.BigEndian, iscnCount)
	iscnID = hasher.Sum(nil)
	kernel.Set("id", iscnID)
	_, err = k.SetIscnKernel(ctx, iscnID, kernel)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			EventTypeCreateIscn,
			sdk.NewAttribute(AttributeKeyIscnID, iscnID.String()),
		),
	)
	return iscnID, nil
}
