package auth

import (
	"context"

	proto "github.com/gogo/protobuf/proto"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

var TypeUrlBaseAccount = "/" + proto.MessageName((*types.BaseAccount)(nil))
var TypeUrlModuleAccount = "/" + proto.MessageName((*types.ModuleAccount)(nil))

type Querier struct {
	types.QueryServer
	cdc codec.BinaryCodec
}

var _ types.QueryServer = Querier{}

func NewQuerier(queryServer types.QueryServer, cdc codec.BinaryCodec) Querier {
	return Querier{
		QueryServer: queryServer,
		cdc:         cdc,
	}
}

func reEncodeAccount(any *codectypes.Any, cdc codec.BinaryCodec) *codectypes.Any {
	switch any.TypeUrl {
	case TypeUrlBaseAccount:
		{
			var acc types.BaseAccount
			err := cdc.Unmarshal(any.Value, &acc)
			if err != nil {
				break
			}
			acc.Address = acc.GetAddress().String()
			newAny, err := codectypes.NewAnyWithValue(&acc)
			if err != nil {
				break
			}
			return newAny
		}
	case TypeUrlModuleAccount:
		{
			var acc types.ModuleAccount
			err := cdc.Unmarshal(any.Value, &acc)
			if err != nil {
				break
			}
			acc.Address = acc.GetAddress().String()
			newAny, err := codectypes.NewAnyWithValue(&acc)
			if err != nil {
				break
			}
			return newAny
		}
	}
	return any
}

func (q Querier) Accounts(c context.Context, req *types.QueryAccountsRequest) (*types.QueryAccountsResponse, error) {
	res, err := q.Accounts(c, req)
	if err != nil {
		return res, err
	}
	for i, any := range res.Accounts {
		res.Accounts[i] = reEncodeAccount(any, q.cdc)
	}
	return res, nil
}

func (q Querier) Account(c context.Context, req *types.QueryAccountRequest) (*types.QueryAccountResponse, error) {
	res, err := q.QueryServer.Account(c, req)
	if err != nil {
		return res, err
	}
	res.Account = reEncodeAccount(res.Account, q.cdc)
	return res, nil
}
