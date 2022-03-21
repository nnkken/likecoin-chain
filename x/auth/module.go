package auth

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

type AppModule struct {
	auth.AppModule
	accountKeeper keeper.AccountKeeper
	querier       Querier
}

func NewAppModule(cdc codec.Codec, accountKeeper keeper.AccountKeeper, randGenAccountsFn types.RandomGenesisAccountsFn) AppModule {
	return AppModule{
		AppModule:     auth.NewAppModule(cdc, accountKeeper, randGenAccountsFn),
		accountKeeper: accountKeeper,
		querier:       NewQuerier(accountKeeper, cdc),
	}
}

func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterQueryServer(cfg.QueryServer(), am.querier)
	m := keeper.NewMigrator(am.accountKeeper, cfg.QueryServer())
	err := cfg.RegisterMigration(types.ModuleName, 1, m.Migrate1to2)
	if err != nil {
		panic(err)
	}
}
