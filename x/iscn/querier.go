package iscn

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gocid "github.com/ipfs/go-cid"
	cbornode "github.com/ipfs/go-ipld-cbor"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryParams:
			return queryParams(ctx, req, k)
		case QueryIscnKernel:
			return queryKernel(ctx, req, k)
		case QueryCID:
			return queryCID(ctx, req, k)
		case QueryCidBlockGet:
			return queryCidBlockGet(ctx, req, k)
		case QueryCidBlockGetSize:
			return queryCidBlockGetSize(ctx, req, k)
		case QueryCidBlockHas:
			return queryCidBlockHas(ctx, req, k)
		// TODO: QueryIscnContent, QueryStakeholders, QueryRights, QueryRightTerms
		default:
			return nil, sdk.ErrUnknownRequest("unknown iscn query endpoint")
		}
	}
}

func queryParams(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	params := k.GetParams(ctx)
	res, err := codec.MarshalJSONIndent(ModuleCdc, params)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}
	return res, nil
}

func queryKernel(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	kernelCID := k.GetIscnKernelCIDByIscnID(ctx, req.Data)
	if kernelCID == nil {
		return nil, nil
	}
	return kernelCID.Bytes(), nil
}

func queryCID(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	_, cid, err := gocid.CidFromBytes(req.Data)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("invalid CID", err.Error()))
	}
	bz := k.GetCidBlock(ctx, cid)
	if bz == nil {
		return nil, nil
	}
	var v interface{}
	switch cid.Prefix().GetCodec() {
	case RightTermsCodecType:
		v = string(bz)
	default:
		v = IscnData{}
		err = cbornode.DecodeInto(bz, &v)
		if err != nil {
			return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to cbor unmarshal ISCN data result", err.Error()))
		}
	}
	res, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result", err.Error()))
	}

	return res, nil
}

func queryCidBlockGet(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	_, cid, err := gocid.CidFromBytes(req.Data)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("invalid CID", err.Error()))
	}
	return k.GetCidBlock(ctx, cid), nil
}

func queryCidBlockGetSize(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	_, cid, err := gocid.CidFromBytes(req.Data)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("invalid CID", err.Error()))
	}
	block := k.GetCidBlock(ctx, cid)

	res, err := codec.MarshalJSONIndent(ModuleCdc, len(block))
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result", err.Error()))
	}

	return res, nil
}

func queryCidBlockHas(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	_, cid, err := gocid.CidFromBytes(req.Data)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("invalid CID: %s", err.Error()))
	}
	has := k.HasCidBlock(ctx, cid)

	res, err := codec.MarshalJSONIndent(ModuleCdc, has)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}

	return res, nil
}
