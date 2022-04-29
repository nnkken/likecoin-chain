package staking

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"sync"

	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"
	dbstore "github.com/cosmos/cosmos-sdk/store/dbadapter"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"

	"github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

var ValidatorDelegationsIndexPrefix = GetTableKey("ValidatorDelegationsIndexPrefix")
var ValidatorDelegationsIndexHeightKey = GetTableKey("ValidatorDelegationsIndexHeightKey")

func GetTableKey(name string) []byte {
	hash := sha256.Sum256([]byte(name))
	return hash[:8]
}

// `append` has concurrency issue so can't be used here
func getIndexValidatorPrefix(valAddr sdk.ValAddress) []byte {
	buf := &bytes.Buffer{}
	buf.Write(ValidatorDelegationsIndexPrefix)
	buf.Write(address.MustLengthPrefix(valAddr))
	return buf.Bytes()
}

// `append` has concurrency issue so can't be used here
func getIndexKey(valAddr sdk.ValAddress, delAddr sdk.AccAddress) []byte {
	buf := &bytes.Buffer{}
	buf.Write(ValidatorDelegationsIndexPrefix)
	buf.Write(address.MustLengthPrefix(valAddr))
	buf.Write(address.MustLengthPrefix(delAddr))
	return buf.Bytes()
}

func encodeUint64(n uint64) []byte {
	heightBz := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBz, n)
	return heightBz
}

type Querier struct {
	types.QueryServer
	keeper            *keeper.Keeper
	cdc               codec.BinaryCodec
	indexingDB        dbm.DB
	batch             dbm.Batch
	batchFlushQueue   chan dbm.Batch
	flushSwitchSignal chan bool
	buildingIndex     bool
	lock              sync.Mutex
}

func NewQuerier(k *keeper.Keeper, cdc codec.BinaryCodec, indexingDB dbm.DB) *Querier {
	q := Querier{
		QueryServer:     keeper.Querier{Keeper: *k},
		keeper:          k,
		cdc:             cdc,
		indexingDB:      indexingDB,
		batch:           indexingDB.NewBatch(),
		batchFlushQueue: make(chan dbm.Batch, 100),
		buildingIndex:   false,
		lock:            sync.Mutex{},
	}
	go q.batchFlusher()
	return &q
}

func (q *Querier) getStore() dbstore.Store {
	return dbstore.Store{DB: q.indexingDB}
}

func (q *Querier) GetHeight() uint64 {
	bz := q.getStore().Get(ValidatorDelegationsIndexHeightKey)
	if bz == nil {
		return 0
	}
	return binary.BigEndian.Uint64(bz)
}

func (q *Querier) setHeight(newHeight uint64) {
	q.batchSet(ValidatorDelegationsIndexHeightKey, encodeUint64(newHeight))
}

// TODO: need to make this async so it won't block the server (and drop connections)
func (q *Querier) BuildIndex(ctx sdk.Context) {
	blockHeight := ctx.BlockHeight()
	ctx.Logger().Debug("Rebuilding index for ValidatorDelegations", "block_height", blockHeight)
	cachedStore, err := ctx.MultiStore().CacheMultiStoreWithVersion(blockHeight)
	if err != nil {
		panic(err)
	}
	newCtx := ctx.WithMultiStore(cachedStore)
	b := q.queueBatch()
	go func() {
		// clear existing index
		it, err := q.indexingDB.Iterator(nil, nil)
		if err != nil {
			panic(err)
		}
		defer it.Close()
		for ; it.Valid(); it.Next() {
			err := b.Delete(it.Key())
			if err != nil {
				panic(err)
			}
		}

		q.keeper.IterateAllDelegations(newCtx, func(delegation types.Delegation) bool {
			key := getIndexKey(delegation.GetValidatorAddr(), delegation.GetDelegatorAddr())
			err := b.Set(key, types.MustMarshalDelegation(q.cdc, delegation))
			if err != nil {
				panic(err)
			}
			return false
		})
		q.batchSet(ValidatorDelegationsIndexHeightKey, encodeUint64(uint64(ctx.BlockHeight())))
	}()
}

func (q *Querier) BeginWriteIndex(ctx sdk.Context) {
	currentHeight := q.GetHeight()
	blockHeight := uint64(ctx.BlockHeight())
	ctx.Logger().Debug(
		"Beginning write index",
		"index_height", currentHeight,
		"block_height", blockHeight,
	)
	if currentHeight == 0 || blockHeight != currentHeight+1 {
		q.BuildIndex(ctx)
	}
}

func (q *Querier) RemoveIndex(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	ctx.Logger().Debug(
		"Removing index",
		"block_height", ctx.BlockHeight(),
		"delegator", delAddr.String(),
		"validator", valAddr.String(),
	)
	key := getIndexKey(valAddr, delAddr)
	q.batchDelete(key)
}

func (q *Querier) UpdateIndex(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	ctx.Logger().Debug(
		"Updating index",
		"block_height", ctx.BlockHeight(),
		"delegator", delAddr.String(),
		"validator", valAddr.String(),
	)
	key := getIndexKey(valAddr, delAddr)
	delegation, found := q.keeper.GetDelegation(ctx, delAddr, valAddr)
	if !found {
		ctx.Logger().Debug("Delegation not found, deleting index")
		q.batchDelete(key)
	} else {
		ctx.Logger().Debug("Delegation found, updating index", "delegation", delegation)
		q.batchSet(key, types.MustMarshalDelegation(q.cdc, delegation))
	}
}

func (q *Querier) CommitWriteIndex(ctx sdk.Context) {
	blockHeight := uint64(ctx.BlockHeight())
	ctx.Logger().Debug(
		"Committing index writes",
		"block_height", blockHeight,
	)
	q.setHeight(blockHeight)
	q.queueBatch()
}
