package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/likecoin/likechain/x/likenft/types"
)

func (k msgServer) CreateListing(goCtx context.Context, msg *types.MsgCreateListing) (*types.MsgCreateListingResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	userAddress, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf(err.Error())
	}

	// check user own the nft
	if !k.nftKeeper.GetOwner(ctx, msg.ClassId, msg.NftId).Equals(userAddress) {
		return nil, types.ErrFailedToCreateListing.Wrapf("User do not own the NFT")
	}

	// Validate expiration range
	if err := validateListingExpiration(ctx, msg.Expiration); err != nil {
		return nil, err
	}

	// Check if the value already exists
	_, isFound := k.GetListing(
		ctx,
		msg.ClassId,
		msg.NftId,
		userAddress,
	)
	if isFound {
		return nil, types.ErrListingAlreadyExists
	}

	// Prune listings by non-owners
	k.PruneInvalidListingsForNFT(ctx, msg.ClassId, msg.NftId)

	// Create new listing
	var listing = types.ListingStoreRecord{
		ClassId:    msg.ClassId,
		NftId:      msg.NftId,
		Seller:     userAddress,
		Price:      msg.Price,
		Expiration: msg.Expiration,
	}

	k.SetListing(
		ctx,
		listing,
	)

	k.SetListingExpireQueueEntry(ctx, types.ListingExpireQueueEntry{
		ExpireTime: listing.Expiration,
		ListingKey: types.ListingKey(listing.ClassId, listing.NftId, listing.Seller),
	})

	pubListing := listing.ToPublicRecord()

	ctx.EventManager().EmitTypedEvent(&types.EventCreateListing{
		ClassId: pubListing.ClassId,
		NftId:   pubListing.NftId,
		Seller:  pubListing.Seller,
	})

	return &types.MsgCreateListingResponse{
		Listing: pubListing,
	}, nil
}

func (k msgServer) UpdateListing(goCtx context.Context, msg *types.MsgUpdateListing) (*types.MsgUpdateListingResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	userAddress, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf(err.Error())
	}

	// Check if the value exists
	oldListing, isFound := k.GetListing(
		ctx,
		msg.ClassId,
		msg.NftId,
		userAddress,
	)
	if !isFound {
		return nil, types.ErrListingNotFound
	}

	// Validate expiration range
	if err := validateListingExpiration(ctx, msg.Expiration); err != nil {
		return nil, err
	}

	var newListing = types.ListingStoreRecord{
		ClassId:    msg.ClassId,
		NftId:      msg.NftId,
		Seller:     userAddress,
		Price:      msg.Price,
		Expiration: msg.Expiration,
	}

	k.SetListing(ctx, newListing)

	k.UpdateListingExpireQueueEntry(
		ctx,
		oldListing.Expiration,
		types.OfferKey(oldListing.ClassId, oldListing.NftId, oldListing.Seller),
		newListing.Expiration,
	)

	pubListing := newListing.ToPublicRecord()

	ctx.EventManager().EmitTypedEvent(&types.EventUpdateListing{
		ClassId: pubListing.ClassId,
		NftId:   pubListing.NftId,
		Seller:  pubListing.Seller,
	})

	return &types.MsgUpdateListingResponse{
		Listing: pubListing,
	}, nil
}

func (k msgServer) DeleteListing(goCtx context.Context, msg *types.MsgDeleteListing) (*types.MsgDeleteListingResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	userAddress, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf(err.Error())
	}

	// Check if the value exists
	listing, isFound := k.GetListing(
		ctx,
		msg.ClassId,
		msg.NftId,
		userAddress,
	)
	if !isFound {
		return nil, types.ErrListingNotFound
	}

	k.RemoveListing(
		ctx,
		listing.ClassId,
		listing.NftId,
		listing.Seller,
	)

	k.RemoveListingExpireQueueEntry(
		ctx,
		listing.Expiration,
		types.OfferKey(listing.ClassId, listing.NftId, listing.Seller),
	)

	pubListing := listing.ToPublicRecord()

	ctx.EventManager().EmitTypedEvent(&types.EventDeleteListing{
		ClassId: pubListing.ClassId,
		NftId:   pubListing.NftId,
		Seller:  pubListing.Seller,
	})

	return &types.MsgDeleteListingResponse{}, nil
}
