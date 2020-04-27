package iscn

import (
	"github.com/likecoin/likechain/x/iscn/types"
)

const (
	ModuleName           = types.ModuleName
	StoreKey             = types.StoreKey
	QuerierRoute         = types.QuerierRoute
	RouterKey            = types.RouterKey
	QueryEntity          = types.QueryEntity
	QueryIscnRecord      = types.QueryIscnRecord
	QueryParams          = types.QueryParams
	QueryCidBlockGet     = types.QueryCidBlockGet
	QueryCidBlockGetSize = types.QueryCidBlockGetSize
	QueryCidBlockHas     = types.QueryCidBlockHas
)

var (
	ModuleCdc                   = types.ModuleCdc
	NewMsgCreateIscn            = types.NewMsgCreateIscn
	ErrInvalidApprover          = types.ErrInvalidApprover
	ErrValidatorNotInWEhitelist = types.ErrValidatorNotInWEhitelist
	KeyFeePerByte               = types.KeyFeePerByte
	DefaultParams               = types.DefaultParams
	DefaultGenesisState         = types.DefaultGenesisState
	DefaultCodespace            = types.DefaultCodespace
	ValidateGenesis             = types.ValidateGenesis
	IscnRecordKey               = types.IscnRecordKey
	IscnCountKey                = types.IscnCountKey
	EntityKey                   = types.EntityKey
	RightTermsKey               = types.RightTermsKey
	CidBlockKey                 = types.CidBlockKey
	GetIscnRecordKey            = types.GetIscnRecordKey
	GetEntityKey                = types.GetEntityKey
	GetRightTermsKey            = types.GetRightTermsKey
	GetCidBlockKey              = types.GetCidBlockKey
	EventTypeCreateIscn         = types.EventTypeCreateIscn
	EventTypeAddEntity          = types.EventTypeAddEntity
	AttributeKeyIscnId          = types.AttributeKeyIscnId
	AttributeKeyEntityCid       = types.AttributeKeyEntityCid
	AttributeValueCategory      = types.AttributeValueCategory
	RegisterCodec               = types.RegisterCodec
)

type (
	MsgCreateIscn   = types.MsgCreateIscn
	MsgAddEntity    = types.MsgAddEntity
	EntityByCid     = types.EntityByCid
	EntityInput     = types.EntityInput
	RightTermsByCid = types.RightTermsByCid
	RightTermsInput = types.RightTermsInput
	IscnRecord      = types.IscnRecord
	Entity          = types.Entity
	Stakeholder     = types.Stakeholder
	Right           = types.Right
	RightTerms      = types.RightTerms
	IscnContent     = types.IscnContent
	CID             = types.CID
	Params          = types.Params
	GenesisState    = types.GenesisState
)
