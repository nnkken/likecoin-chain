package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgCreateIscn{}
var _ sdk.Msg = &MsgAddEntity{}

type EntityInput interface {
	EntityInput()
}

func (e Entity) EntityInput() {}

type EntityByCid struct {
	Cid string `json:"cid" yaml:"cid"`
}

func (e EntityByCid) EntityInput() {}

type RightTermsInput interface {
	RightTermsInput()
}

func (rt RightTerms) RightTermsInput() {}

type RightTermsByCid struct {
	Cid string `json:"cid" yaml:"cid"`
}

func (rt RightTermsByCid) RightTermsInput() {}

type StakeholderInput struct {
	Type   string      `json:"type" yaml:"type"`
	Entity EntityInput `json:"entity" yaml:"entity"`
	Stake  uint32      `json:"stake" yaml:"stake"`
}

type RightInput struct {
	Type      string          `json:"type" yaml:"type"`
	Holder    EntityInput     `json:"holder" yaml:"holder"`
	Terms     RightTermsInput `json:"terms" yaml:"terms"`
	Period    Period          `json:"period" yaml:"period"`
	Territory string          `json:"territory" yaml:"territory"`
}

type IscnRecordInput struct {
	Stakeholders []StakeholderInput `json:"stakeholders" yaml:"stakeholders"`
	Timestamp    string             `json:"timestamp" yaml:"timestamp"`
	Parent       string             `json:"parent" yaml:"parent"` // CID string to parent
	Version      uint32             `json:"version" yaml:"version"`
	Rights       []RightInput       `json:"rights" yaml:"rights"`
	Content      IscnContent        `json:"content" yaml:"content"`
}

type MsgCreateIscn struct {
	From       sdk.AccAddress  `json:"from" yaml:"from"`
	IscnRecord IscnRecordInput `json:"iscnRecord" yaml:"iscnRecord"`
}

func NewMsgCreateIscn(from sdk.AccAddress, iscnRecord IscnRecordInput) MsgCreateIscn {
	return MsgCreateIscn{
		From:       from,
		IscnRecord: iscnRecord,
	}
}

func (msg MsgCreateIscn) Route() string { return RouterKey }
func (msg MsgCreateIscn) Type() string  { return "create_iscn" }

func (msg MsgCreateIscn) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

func (msg MsgCreateIscn) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgCreateIscn) ValidateBasic() sdk.Error {
	if msg.From.Empty() {
		return ErrInvalidApprover(DefaultCodespace)
	}
	// TODO: validate IscnRecord
	return nil
}

type MsgAddEntity struct {
	From       sdk.AccAddress `json:"from" yaml:"from"`
	EntityInfo Entity         `json:"entityInfo" yaml:"entityInfo"`
}

func NewMsgAddEntity(from sdk.AccAddress, entityInfo Entity) MsgAddEntity {
	return MsgAddEntity{
		From:       from,
		EntityInfo: entityInfo,
	}
}

func (msg MsgAddEntity) Route() string { return RouterKey }
func (msg MsgAddEntity) Type() string  { return "add_entity" }

func (msg MsgAddEntity) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

func (msg MsgAddEntity) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgAddEntity) ValidateBasic() sdk.Error {
	if msg.From.Empty() {
		return ErrInvalidApprover(DefaultCodespace)
	}
	// TODO: validate IscnRecord
	return nil
}
