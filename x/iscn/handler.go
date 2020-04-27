package iscn

import (
	"encoding/base64"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgCreateIscn:
			return handleMsgCreateIscn(ctx, msg, keeper)
		case MsgAddEntity:
			return handleMsgAddEntity(ctx, msg, keeper)
		default:
			errMsg := fmt.Sprintf("unrecognized iscn message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleEntityInput(ctx sdk.Context, entity EntityInput, keeper Keeper) (CID, sdk.Error) {
	var cid CID
	switch entity.(type) {
	case Entity:
		e, _ := entity.(Entity)
		return keeper.SetEntity(ctx, &e), nil
	case EntityByCid:
		e, _ := entity.(EntityByCid)
		// TODO: string to CID
		cid = []byte(e.Cid)
		if keeper.GetEntity(ctx, cid) == nil {
			return nil, sdk.NewError(DefaultCodespace, 123, "Unknown entity CID: %s", e.Cid)
		}
	default:
		// TODO: better handling?
		panic("Unknown EntityInput type")
	}
	return cid, nil
}

func handleRightTermsInput(ctx sdk.Context, rightTerms RightTermsInput, keeper Keeper) (CID, sdk.Error) {
	var cid CID
	switch rightTerms.(type) {
	case RightTerms:
		// TODO: need to implement SetRightTerms
		// e, _ := rightTerms.(RightTerms)
		// return keeper.SetRightTerms(ctx, &e), nil
		return nil, nil
	case RightTermsByCid:
		e, _ := rightTerms.(RightTermsByCid)
		// TODO: string to CID
		cid = []byte(e.Cid)
		// TODO: need to implement GetRightTerms
		// if keeper.GetRightTerms(ctx, cid) == nil {
		// 	return nil, sdk.NewError(DefaultCodespace, 123, "Unknown RightTerms CID: %s", e.Cid)
		// }
	default:
		// TODO: better handling?
		panic("Unknown RightTermsInput type")
	}
	return cid, nil
}

func handleMsgCreateIscn(ctx sdk.Context, msg MsgCreateIscn, keeper Keeper) sdk.Result {
	// TODO:
	// 1. store nested fields and construct CIDs
	// 2. validate fields from IscnRecordInput
	// 3. construct IscnRecord
	iscnInput := msg.IscnRecord
	iscnRecord := IscnRecord{
		// TODO: parse iscnInput.Timestamp
		Timestamp: 0,
		// TODO: parse CID
		Parent:  []byte(iscnInput.Parent),
		Version: iscnInput.Version,
		Content: iscnInput.Content,
	}
	for i := range msg.IscnRecord.Stakeholders {
		input := msg.IscnRecord.Stakeholders[i]
		stakeholder := Stakeholder{
			Type:  input.Type,
			Stake: input.Stake,
		}
		entity := input.Entity
		cid, err := handleEntityInput(ctx, entity, keeper)
		if err != nil {
			return sdk.Result{
				/* TODO: proper error*/
				Code:      err.Code(),
				Codespace: err.Codespace(),
				Log:       err.Error(),
			}
		}
		stakeholder.Entity = cid
		iscnRecord.Stakeholders = append(iscnRecord.Stakeholders, stakeholder)
	}
	for i := range msg.IscnRecord.Rights {
		input := msg.IscnRecord.Rights[i]
		right := Right{
			Type:      input.Type,
			Period:    input.Period,
			Territory: input.Territory,
		}
		entity := input.Holder
		cid, err := handleEntityInput(ctx, entity, keeper)
		if err != nil {
			return sdk.Result{
				/* TODO: proper error*/
				Code:      err.Code(),
				Codespace: err.Codespace(),
				Log:       err.Error(),
			}
		}
		right.Holder = cid
		rightTerms := input.Terms
		cid, err = handleRightTermsInput(ctx, rightTerms, keeper)
		if err != nil {
			return sdk.Result{
				/* TODO: proper error*/
				Code:      err.Code(),
				Codespace: err.Codespace(),
				Log:       err.Error(),
			}
		}
		right.Terms = cid
		iscnRecord.Rights = append(iscnRecord.Rights, right)
	}
	// TODO: validation
	id, err := keeper.AddIscnRecord(ctx, msg.From, &iscnRecord)
	if err != nil {
		return sdk.Result{
			/* TODO: proper error*/
			Code:      err.Code(),
			Codespace: err.Codespace(),
			Log:       err.Error(),
		}
	}
	idStr := base64.URLEncoding.EncodeToString(id) // TODO: formatting iscn
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From.String()),
		),
		sdk.NewEvent(
			EventTypeCreateIscn,
			sdk.NewAttribute(AttributeKeyIscnId, idStr),
		),
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgAddEntity(ctx sdk.Context, msg MsgAddEntity, keeper Keeper) sdk.Result {
	// TODO: extract so we can reuse logic in MsgCreateIscn
	entityCid := keeper.SetEntity(ctx, &msg.EntityInfo)
	cidStr := base64.URLEncoding.EncodeToString(entityCid) // TODO: formatting cid
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventTypeAddEntity,
			sdk.NewAttribute(AttributeKeyEntityCid, cidStr),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From.String()),
		),
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}
