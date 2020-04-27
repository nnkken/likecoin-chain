package iscn

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/likechain/x/iscn/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func InitGenesis(ctx sdk.Context, keeper Keeper, genesisState GenesisState) []abci.ValidatorUpdate {
	keeper.SetParams(ctx, genesisState.Params)
	for _, record := range genesisState.IscnRecords {
		keeper.SetIscnRecord(ctx, record.Id, &record.Record)
	}
	for _, entity := range genesisState.Entities {
		keeper.SetEntity(ctx, &entity)
	}
	keeper.SetIscnCount(ctx, uint64(len(genesisState.IscnRecords)))
	return nil
}

func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	params := keeper.GetParams(ctx)
	entities := []Entity{}
	keeper.IterateEntitys(ctx, func(_ []byte, entity *Entity) bool {
		entities = append(entities, *entity)
		return false
	})
	records := []types.IscnPair{}
	keeper.IterateIscnRecords(ctx, func(id []byte, record *IscnRecord) bool {
		records = append(records, types.IscnPair{
			Id:     id,
			Record: *record,
		})
		return false
	})
	// TODO: CIDs
	return GenesisState{
		Params:      params,
		Entities:    entities,
		IscnRecords: records,
	}
}
