package iscn

import (
	"encoding/base64"
	"encoding/binary"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

const (
	DefaultParamspace = ModuleName
)

// TODO: move to individual file
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

func (keeper Keeper) Codespace() sdk.CodespaceType {
	return keeper.codespace
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

func (k Keeper) SetEntity(ctx sdk.Context, entity *Entity) CID {
	// TODO: call SetCid and record only the type is entity
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(entity) // TODO: cbor?
	cid := tmhash.Sum(bz)                               // TODO: cid
	key := GetEntityKey(cid)
	ctx.KVStore(k.storeKey).Set(key, bz)
	return cid
}

func (k Keeper) GetEntity(ctx sdk.Context, cid CID) *Entity {
	key := GetEntityKey(cid)
	bz := ctx.KVStore(k.storeKey).Get(key)
	if bz == nil {
		return nil
	}
	entity := Entity{}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &entity)
	return &entity
}

func (k Keeper) IterateEntitys(ctx sdk.Context, f func(cid CID, entity *Entity) bool) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, EntityKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		entityCid := iterator.Key()[len(EntityKey):]
		entity := Entity{}
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &entity)
		stop := f(entityCid, &entity)
		if stop {
			break
		}
	}
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

func (k Keeper) SetIscnRecord(ctx sdk.Context, iscnId []byte, record *IscnRecord) {
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(record) // TODO: cbor?
	key := GetIscnRecordKey(iscnId)
	ctx.KVStore(k.storeKey).Set(key, bz)
}

func (k Keeper) GetIscnRecord(ctx sdk.Context, iscnId []byte) *IscnRecord {
	key := GetIscnRecordKey(iscnId)
	bz := ctx.KVStore(k.storeKey).Get(key)
	if bz == nil {
		return nil
	}
	record := IscnRecord{}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &record)
	return &record
}

func (k Keeper) IterateIscnRecords(ctx sdk.Context, f func(iscnId []byte, iscnRecord *IscnRecord) bool) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, IscnRecordKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		id := iterator.Key()[len(IscnRecordKey):]
		record := IscnRecord{}
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &record)
		stop := f(id, &record)
		if stop {
			break
		}
	}
}

var allowedStakeholderTypes = map[string]bool{
	"author": true,
	// TODO: more types
}

func (k Keeper) AddIscnRecord(ctx sdk.Context, feePayer sdk.AccAddress, record *IscnRecord) (iscnId []byte, err sdk.Error) {
	acc := k.accountKeeper.GetAccount(ctx, feePayer)
	if acc == nil {
		return nil, sdk.NewError(DefaultCodespace, 1, "No account") // TODO: proper error
	}
	// TODO: checkings (in handler?)
	// 1. len(stakeholder) > 0?
	// 2. one entity only?
	feePerByte := k.GetParams(ctx).FeePerByte
	feeAmount := feePerByte.Amount.MulInt64(int64(len(ctx.TxBytes())))
	fees := sdk.NewCoins(sdk.NewCoin(feePerByte.Denom, feeAmount.Ceil().RoundInt()))
	result := auth.DeductFees(k.supplyKeeper, ctx, acc, fees)
	if !result.IsOK() {
		// TODO: proper error
		return nil, sdk.NewError(DefaultCodespace, 2, "Not enough fee, %s %s needed", feeAmount.String(), feePerByte.Denom)
	}
	for i := range record.Stakeholders {
		stakeholder := &record.Stakeholders[i]
		_, ok := allowedStakeholderTypes[stakeholder.Type]
		if !ok {
			// TODO: proper error
			return nil, sdk.NewError(DefaultCodespace, 3, "Unknown stakeholder type: %s", stakeholder.Type)
		}
		if stakeholder.Type == "entity" {
			entityPrefix := "likecoin-chain://entities/" // TODO: better format
			if strings.HasPrefix(stakeholder.Entity, entityPrefix) {
				cidStr := stakeholder.Entity[len(entityPrefix):]
				cid, err := base64.URLEncoding.DecodeString(cidStr)
				if err != nil {
					// TODO: proper error
					return nil, sdk.NewError(DefaultCodespace, 4, "Invalid entity string: %s", err.Error())
				}
				if k.GetEntity(ctx, cid) == nil {
					// TODO: proper error
					return nil, sdk.NewError(DefaultCodespace, 5, "Unknown entity: %s", cidStr)
				}
				stakeholder.Entity = cidStr
			} else {
				entity := Entity{}
				err := k.cdc.UnmarshalJSON([]byte(stakeholder.Entity), &entity)
				if err != nil {
					// TODO: proper error
					return nil, sdk.NewError(DefaultCodespace, 5, "Cannot decode entity: %s", err.Error())
				}
				cid := k.SetEntity(ctx, &entity)
				cidStr := base64.URLEncoding.EncodeToString(cid)
				stakeholder.Entity = cidStr
				// TODO: return entity CID for tx event, or return event
			}
		}
	}
	hasher := tmhash.New()
	hasher.Write(ctx.BlockHeader().LastBlockId.Hash)
	iscnCount := k.GetIscnCount(ctx)
	k.SetIscnCount(ctx, iscnCount+1)
	binary.Write(hasher, binary.BigEndian, iscnCount)
	id := hasher.Sum(nil)
	k.SetIscnRecord(ctx, id, record)
	return id, nil
}

func (k Keeper) SetCidBlock(ctx sdk.Context, cid []byte, block []byte) {
	key := GetCidBlockKey(cid)
	ctx.KVStore(k.storeKey).Set(key, block)
}

func (k Keeper) GetCidBlock(ctx sdk.Context, cid []byte) []byte {
	key := GetCidBlockKey(cid)
	return ctx.KVStore(k.storeKey).Get(key)
}

func (k Keeper) HasCidBlock(ctx sdk.Context, cid []byte) bool {
	key := GetCidBlockKey(cid)
	return ctx.KVStore(k.storeKey).Has(key)
}

func (k Keeper) IterateCidBlocks(ctx sdk.Context, f func(cid []byte, block []byte) bool) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, CidBlockKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		cid := iterator.Key()[len(CidBlockKey):]
		stop := f(cid, iterator.Value())
		if stop {
			break
		}
	}
}
