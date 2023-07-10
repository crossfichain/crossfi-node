package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	"math/big"
	"strconv"

	"github.com/mineplexio/mineplex-2-node/x/gravity/types"
)

// Check that distKeeper implements the expected type
var _ types.DistributionKeeper = (*distrkeeper.Keeper)(nil)

// AttestationHandler processes `observed` Attestations
type AttestationHandler struct {
	// NOTE: If you add anything to this struct, add a nil check to ValidateMembers below!
	keeper *Keeper
}

// Check for nil members
func (a AttestationHandler) ValidateMembers() {
	if a.keeper == nil {
		panic("Nil keeper!")
	}
}

// Handle is the entry point for Attestation processing, only attestations with sufficient validator submissions
// should be processed through this function, solidifying their effect in chain state
func (a AttestationHandler) Handle(ctx sdk.Context, att types.Attestation, claim types.EthereumClaim) error {
	switch claim := claim.(type) {

	case *types.MsgSendToCosmosClaim:
		return a.handleSendToCosmos(ctx, *claim)

	case *types.MsgBatchSendToEthClaim:
		return a.handleBatchSendToEth(ctx, *claim)

	case *types.MsgERC20DeployedClaim:
		return a.handleErc20Deployed(ctx, *claim)

	case *types.MsgValsetUpdatedClaim:
		return a.handleValsetUpdated(ctx, *claim)

	default:
		panic(fmt.Sprintf("Invalid event type for attestations %s", claim.GetType()))
	}
}

// Upon acceptance of sufficient validator SendToCosmos claims: transfer tokens to the appropriate cosmos account
// The cosmos receiver can be a native account (e.g. gravity1abc...) or a foreign account (e.g. cosmos1abc...)
// In the event of a native receiver, bank module handles the transfer, otherwise an IBC transfer is initiated
// Note: Previously SendToCosmos was referred to as a bridge "Deposit", as tokens are deposited into the gravity contract
func (a AttestationHandler) handleSendToCosmos(ctx sdk.Context, claim types.MsgSendToCosmosClaim) error {
	invalidAddress := false
	// Validate the receiver as a valid bech32 address
	receiverAddress, addressErr := types.IBCAddressFromBech32(claim.CosmosReceiver)

	if addressErr != nil {
		invalidAddress = true
		hash, er := claim.ClaimHash()
		if er != nil {
			return sdkerrors.Wrapf(er, "Unable to log error %v, could not compute ClaimHash for claim %v: %v", addressErr, claim, er)
		}

		a.keeper.logger(ctx).Error("Invalid SendToCosmos receiver",
			"address", receiverAddress,
			"cause", addressErr.Error(),
			"claim type", claim.GetType(),
			"id", types.GetAttestationKey(claim.ChainID(), claim.GetEventNonce(), hash),
			"nonce", fmt.Sprint(claim.GetEventNonce()),
		)
	}
	tokenAddress, errTokenAddress := types.NewEthAddress(claim.TokenContract)
	ethereumSender, errEthereumSender := types.NewEthAddress(claim.EthereumSender)
	// nil address is not possible unless the validators get together and submit
	// a bogus event, this would create lost tokens stuck in the bridge
	// and not accessible to anyone
	if errTokenAddress != nil {
		hash, er := claim.ClaimHash()
		if er != nil {
			return sdkerrors.Wrapf(er, "Unable to log error %v, could not compute ClaimHash for claim %v: %v", errTokenAddress, claim, er)
		}
		a.keeper.logger(ctx).Error("Invalid token contract",
			"cause", errTokenAddress.Error(),
			"claim type", claim.GetType(),
			"id", types.GetAttestationKey(claim.ChainID(), claim.GetEventNonce(), hash),
			"nonce", fmt.Sprint(claim.GetEventNonce()),
		)
		return sdkerrors.Wrap(errTokenAddress, "invalid token contract on claim")
	}
	// likewise nil sender would have to be caused by a bogus event
	if errEthereumSender != nil {
		hash, er := claim.ClaimHash()
		if er != nil {
			return sdkerrors.Wrapf(er, "Unable to log error %v, could not compute ClaimHash for claim %v: %v", errEthereumSender, claim, er)
		}
		a.keeper.logger(ctx).Error("Invalid ethereum sender",
			"cause", errEthereumSender.Error(),
			"claim type", claim.GetType(),
			"id", types.GetAttestationKey(claim.ChainID(), claim.GetEventNonce(), hash),
			"nonce", fmt.Sprint(claim.GetEventNonce()),
		)
		return sdkerrors.Wrap(errTokenAddress, "invalid ethereum sender on claim")
	}

	// Block blacklisted asset transfers
	// (these funds are unrecoverable for the blacklisted sender, they will instead be sent to community pool)
	if a.keeper.IsOnBlacklist(ctx, *ethereumSender) {
		hash, er := claim.ClaimHash()
		if er != nil {
			return sdkerrors.Wrapf(er, "Unable to log blacklisted error, could not compute ClaimHash for claim %v: %v", claim, er)
		}
		a.keeper.logger(ctx).Error("Invalid SendToCosmos: receiver is blacklisted",
			"address", receiverAddress,
			"claim type", claim.GetType(),
			"id", types.GetAttestationKey(claim.ChainID(), claim.GetEventNonce(), hash),
			"nonce", fmt.Sprint(claim.GetEventNonce()),
		)
		invalidAddress = true
	}

	// Check if coin is Cosmos-originated asset and get denom
	isCosmosOriginated, denom := a.keeper.ERC20ToDenomLookup(ctx, claim.ChainID(), *tokenAddress)
	coin := sdk.NewCoin(denom, claim.Amount)
	coins := sdk.Coins{coin}

	moduleAddr := a.keeper.accountKeeper.GetModuleAddress(types.ModuleName)
	if !isCosmosOriginated { // We need to mint eth-originated coins (aka vouchers)
		if err := a.mintEthereumOriginatedVouchers(ctx, moduleAddr, claim, coin); err != nil {
			// TODO: Evaluate closely, if we can't mint an ethereum voucher, what should we do?
			return err
		}
	}

	if !invalidAddress { // address appears valid, attempt to send minted/locked coins to receiver
		preSendBalance := a.keeper.bankKeeper.GetBalance(ctx, moduleAddr, denom)
		// Failure to send will result in funds transfer to community pool
		ibcForwardQueued, err := a.sendCoinToCosmosAccount(ctx, claim, receiverAddress, coin)

		// Perform module balance assertions
		if err != nil || ibcForwardQueued { // ibc forward enqueue and errors should not send tokens to anyone
			a.assertNothingSent(ctx, moduleAddr, preSendBalance, denom)
		} else { // No error, local send -> assert send had right amount
			a.assertSentAmount(ctx, moduleAddr, preSendBalance, denom, claim.Amount)
		}

		if err != nil { // trigger send to community pool

		}
	}

	// for whatever reason above, blacklisted, invalid string, etc this deposit is not valid
	// we can't send the tokens back on the Ethereum side, and if we don't put them somewhere on
	// the cosmos side they will be lost an inaccessible even though they are locked in the bridge.
	// so we deposit the tokens into the community pool for later use via governance vote
	if invalidAddress {
		if err := a.keeper.SendToCommunityPool(ctx, coins); err != nil {
			hash, er := claim.ClaimHash()
			if er != nil {
				return sdkerrors.Wrapf(er, "Unable to log error %v, could not compute ClaimHash for claim %v: %v", err, claim, er)
			}
			a.keeper.logger(ctx).Error("Failed community pool send",
				"cause", err.Error(),
				"claim type", claim.GetType(),
				"id", types.GetAttestationKey(claim.ChainID(), claim.GetEventNonce(), hash),
				"nonce", fmt.Sprint(claim.GetEventNonce()),
			)
			return sdkerrors.Wrap(err, "failed to send to Community pool")
		}

		if err := ctx.EventManager().EmitTypedEvent(
			&types.EventInvalidSendToCosmosReceiver{
				Amount: claim.Amount.String(),
				Nonce:  strconv.Itoa(int(claim.GetEventNonce())),
				Token:  tokenAddress.GetAddress().Hex(),
				Sender: claim.EthereumSender,
			},
		); err != nil {
			return err
		}

	} else {
		if err := ctx.EventManager().EmitTypedEvent(
			&types.EventSendToCosmos{
				Amount: claim.Amount.String(),
				Nonce:  strconv.Itoa(int(claim.GetEventNonce())),
				Token:  tokenAddress.GetAddress().Hex(),
			},
		); err != nil {
			return err
		}
	}

	return nil
}

// Upon acceptance of sufficient validator BatchSendToEth claims: burn ethereum originated vouchers, invalidate pending
// batches with lower claim.BatchNonce, and clean up state
// Note: Previously SendToEth was referred to as a bridge "Withdrawal", as tokens are withdrawn from the gravity contract
func (a AttestationHandler) handleBatchSendToEth(ctx sdk.Context, claim types.MsgBatchSendToEthClaim) error {
	contract, err := types.NewEthAddress(claim.TokenContract)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid token contract on batch")
	}
	a.keeper.OutgoingTxBatchExecuted(ctx, claim.ChainID(), *contract, claim)

	err = ctx.EventManager().EmitTypedEvent(
		&types.EventBatchSendToEthClaim{
			Nonce: strconv.Itoa(int(claim.BatchNonce)),
		},
	)

	return err
}

// Upon acceptance of sufficient ERC20 Deployed claims, register claim.TokenContract as the canonical ethereum
// representation of the metadata governance previously voted for
func (a AttestationHandler) handleErc20Deployed(ctx sdk.Context, claim types.MsgERC20DeployedClaim) error {
	tokenAddress, err := types.NewEthAddress(claim.TokenContract)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid token contract on claim")
	}
	// Disallow re-registration when a token already has a canonical representation
	existingERC20, exists := a.keeper.GetCosmosOriginatedERC20(ctx, claim.ChainID(), claim.CosmosDenom)
	if exists {
		return sdkerrors.Wrap(
			types.ErrInvalid,
			fmt.Sprintf("ERC20 %s already exists for denom %s", existingERC20.GetAddress().Hex(), claim.CosmosDenom))
	}

	// Check if denom metadata has been accepted by governance
	metadata, ok := a.keeper.bankKeeper.GetDenomMetaData(ctx, claim.CosmosDenom)
	if !ok || metadata.Base == "" {
		return sdkerrors.Wrap(types.ErrUnknown, fmt.Sprintf("denom not found %s", claim.CosmosDenom))
	}

	// Check if attributes of ERC20 match Cosmos denom
	if claim.Name != metadata.Name {
		return sdkerrors.Wrap(
			types.ErrInvalid,
			fmt.Sprintf("ERC20 name %s does not match denom name %s", claim.Name, metadata.Description))
	}

	if claim.Symbol != metadata.Symbol {
		return sdkerrors.Wrap(
			types.ErrInvalid,
			fmt.Sprintf("ERC20 symbol %s does not match denom symbol %s", claim.Symbol, metadata.Display))
	}

	// ERC20 tokens use a very simple mechanism to tell you where to display the decimal point.
	// The "decimals" field simply tells you how many decimal places there will be.
	// Cosmos denoms have a system that is much more full featured, with enterprise-ready token denominations.
	// There is a DenomUnits array that tells you what the name of each denomination of the
	// token is.
	// To correlate this with an ERC20 "decimals" field, we have to search through the DenomUnits array
	// to find the DenomUnit which matches up to the main token "display" value. Then we take the
	// "exponent" from this DenomUnit.
	// If the correct DenomUnit is not found, it will default to 0. This will result in there being no decimal places
	// in the token's ERC20 on Ethereum. So, for example, if this happened with Atom, 1 Atom would appear on Ethereum
	// as 1 million Atoms, having 6 extra places before the decimal point.
	// This will only happen with a Denom Metadata which is for all intents and purposes invalid, but I am not sure
	// this is checked for at any other point.
	decimals := uint32(0)
	for _, denomUnit := range metadata.DenomUnits {
		if denomUnit.Denom == metadata.Display {
			decimals = denomUnit.Exponent
			break
		}
	}

	if decimals != uint32(claim.Decimals) {
		return sdkerrors.Wrap(
			types.ErrInvalid,
			fmt.Sprintf("ERC20 decimals %d does not match denom decimals %d", claim.Decimals, decimals))
	}

	// Add to denom-erc20 mapping
	a.keeper.setCosmosOriginatedDenomToERC20(ctx, claim.ChainID(), claim.CosmosDenom, *tokenAddress)

	err = ctx.EventManager().EmitTypedEvent(
		&types.EventERC20DeployedClaim{
			Token: tokenAddress.GetAddress().Hex(),
			Nonce: strconv.Itoa(int(claim.GetEventNonce())),
		},
	)
	return err
}

// Upon acceptance of sufficient ValsetUpdated claims: update LastObservedValset, mint cosmos-originated relayer rewards
// so that the reward holder can send them to cosmos
func (a AttestationHandler) handleValsetUpdated(ctx sdk.Context, claim types.MsgValsetUpdatedClaim) error {
	rewardAddress, err := types.NewEthAddress(claim.RewardToken)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid reward token on claim")
	}

	claimSet := types.Valset{
		Nonce:        claim.ValsetNonce,
		Members:      claim.Members,
		Height:       0, // Fill out later when used
		RewardAmount: claim.RewardAmount,
		RewardToken:  claim.RewardToken,
	}
	// check the contents of the validator set against the store, if they differ we know that the bridge has been
	// highjacked
	if claim.ValsetNonce != 0 { // Handle regular valsets
		trustedValset := a.keeper.GetValset(ctx, claim.ChainID(), claim.ValsetNonce)
		if trustedValset == nil {
			ctx.Logger().Error("Received attestation for a valset which does not exist in store", "nonce", claim.ValsetNonce, "claim", claim)
			return sdkerrors.Wrapf(types.ErrInvalidValset, "attested valset (%v) does not exist in store", claim.ValsetNonce)
		}
		observedValset := claimSet
		observedValset.Height = trustedValset.Height // overwrite the height, since it's not part of the claim

		if _, err := trustedValset.Equal(observedValset); err != nil {
			panic(fmt.Sprintf("Potential bridge highjacking: observed valset (%+v) does not match stored valset (%+v)! %s", observedValset, trustedValset, err.Error()))
		}

		a.keeper.SetLastObservedValset(ctx, claim.ChainID(), observedValset)
	} else { // The 0th valset is not stored on chain init, but we need to set it as the last one
		// Do not update Height, it's the first valset
		a.keeper.SetLastObservedValset(ctx, claim.ChainID(), claimSet)
	}

	// if the reward is greater than zero and the reward token
	// is valid then some reward was issued by this validator set
	// and we need to either add to the total tokens for a Cosmos native
	// token, or burn non cosmos native tokens
	if claim.RewardAmount.GT(sdk.ZeroInt()) && claim.RewardToken != types.ZeroAddressString {
		// Check if coin is Cosmos-originated asset and get denom
		isCosmosOriginated, denom := a.keeper.ERC20ToDenomLookup(ctx, claim.ChainID(), *rewardAddress)
		if isCosmosOriginated {
			// If it is cosmos originated, mint some coins to account
			// for coins that now exist on Ethereum and may eventually come
			// back to Cosmos.
			//
			// Note the flow is
			// user relays valset and gets reward -> event relayed to cosmos mints tokens to module
			// -> user sends tokens to cosmos and gets the minted tokens from the module
			//
			// it is not possible for this to be a race condition thanks to the event nonces
			// no matter how long it takes to relay the valset updated event the deposit event
			// for the user will always come after.
			//
			// Note we are minting based on the claim! This is important as the reward value
			// could change between when this event occurred and the present
			coins := sdk.Coins{sdk.NewCoin(denom, claim.RewardAmount)}
			if err := a.keeper.bankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
				if err := ctx.EventManager().EmitTypedEvent(
					&types.EventValsetUpdatedClaim{
						Nonce: strconv.Itoa(int(claim.GetEventNonce())),
					},
				); err != nil {
					return err
				}
				return sdkerrors.Wrapf(err, "unable to mint cosmos originated coins %v", coins)
			}
		} else {
			// If it is not cosmos originated, burn the coins (aka Vouchers)
			// so that we don't think we have more in the bridge than we actually do
			// coins := sdk.Coins{sdk.NewCoin(denom, claim.RewardAmount)}
			// a.bankKeeper.BurnCoins(ctx, types.ModuleName, coins)

			// if you want to issue Ethereum originated tokens remove this panic and uncomment
			// the above code but note that you will have to constantly replenish the tokens in the
			// module or your chain will eventually halt.
			panic("Can not use Ethereum originated token as reward!")
		}
	}
	err = ctx.EventManager().EmitTypedEvent(
		&types.EventValsetUpdatedClaim{
			Nonce: strconv.Itoa(int(claim.GetEventNonce())),
		},
	)

	return err
}

// assertNothingSent performs a runtime assertion that the actual sent amount of `denom` is zero
func (a AttestationHandler) assertNothingSent(ctx sdk.Context, moduleAddr sdk.AccAddress, preSendBalance sdk.Coin, denom string) {
	postSendBalance := a.keeper.bankKeeper.GetBalance(ctx, moduleAddr, denom)
	if !preSendBalance.Equal(postSendBalance) {
		panic(fmt.Sprintf(
			"SendToCosmos somehow sent tokens in an error case! Previous balance %v Post-send balance %v",
			preSendBalance.String(), postSendBalance.String()),
		)
	}
}

// assertSentAmount performs a runtime assertion that the actual sent amount of `denom` equals the MsgSendToCosmos
// claim's amount to send
func (a AttestationHandler) assertSentAmount(ctx sdk.Context, moduleAddr sdk.AccAddress, preSendBalance sdk.Coin, denom string, amount sdk.Int) {
	postSendBalance := a.keeper.bankKeeper.GetBalance(ctx, moduleAddr, denom)
	if !preSendBalance.Sub(postSendBalance).Amount.Equal(amount) {
		panic(fmt.Sprintf(
			"SendToCosmos somehow sent incorrect amount! Previous balance %v Post-send balance %v claim amount %v",
			preSendBalance.String(), postSendBalance.String(), amount.String()),
		)
	}
}

// mintEthereumOriginatedVouchers creates new "gravity0x..." vouchers for ethereum tokens and asserts both that the
// supply of that voucher does not exceed Uint256 max value, and the minted balance is correct
func (a AttestationHandler) mintEthereumOriginatedVouchers(
	ctx sdk.Context, moduleAddr sdk.AccAddress, claim types.MsgSendToCosmosClaim, coin sdk.Coin,
) error {
	preMintBalance := a.keeper.bankKeeper.GetBalance(ctx, moduleAddr, coin.Denom)
	// Ensure that users are not bridging an impossible amount, only 2^256 - 1 tokens can exist on ethereum
	prevSupply := a.keeper.bankKeeper.GetSupply(ctx, coin.Denom)
	newSupply := new(big.Int).Add(prevSupply.Amount.BigInt(), claim.Amount.BigInt())
	if newSupply.BitLen() > 256 { // new supply overflows uint256
		a.keeper.logger(ctx).Error("Deposit Overflow",
			"claim type", claim.GetType(),
			"nonce", fmt.Sprint(claim.GetEventNonce()),
		)
		return sdkerrors.Wrap(types.ErrIntOverflowAttestation, "invalid supply after SendToCosmos attestation")
	}

	coins := sdk.NewCoins(coin)
	if err := a.keeper.bankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
		// in this case we have lost tokens! They are in the bridge, but not
		// in the community pool or out in some users balance, every instance of this
		// error needs to be detected and resolved
		hash, er := claim.ClaimHash()
		if er != nil {
			return sdkerrors.Wrapf(er, "Unable to log error %v, could not compute ClaimHash for claim %v: %v", err, claim, er)
		}

		a.keeper.logger(ctx).Error("Failed minting",
			"cause", err.Error(),
			"claim type", claim.GetType(),
			"id", types.GetAttestationKey(claim.ChainID(), claim.GetEventNonce(), hash),
			"nonce", fmt.Sprint(claim.GetEventNonce()),
		)
		return sdkerrors.Wrapf(err, "mint vouchers coins: %s", coins)
	}

	postMintBalance := a.keeper.bankKeeper.GetBalance(ctx, moduleAddr, coin.Denom)
	if !postMintBalance.Sub(preMintBalance).Amount.Equal(claim.Amount) {
		panic(fmt.Sprintf(
			"Somehow minted incorrect amount! Previous balance %v Post-mint balance %v claim amount %v",
			preMintBalance.String(), postMintBalance.String(), claim.Amount.String()),
		)
	}
	return nil
}

// Send tokens via bank keeper to a native gravity address, re-prefixing receiver to a gravity native address if necessary
// Note: This should only be used as part of SendToCosmos attestation handling and is not a good solution for general use
func (a AttestationHandler) sendCoinToLocalAddress(
	ctx sdk.Context, claim types.MsgSendToCosmosClaim, receiver sdk.AccAddress, coin sdk.Coin) (err error) {

	err = a.keeper.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver, sdk.NewCoins(coin))
	if err != nil {
		// someone attempted to send tokens to a blacklisted user from Ethereum, log and send to Community pool
		hash, er := claim.ClaimHash()
		if er != nil {
			return sdkerrors.Wrapf(er, "Unable to log error %v, could not compute ClaimHash for claim %v: %v", err, claim, er)
		}
		a.keeper.logger(ctx).Error("Blacklisted deposit",
			"cause", err.Error(),
			"claim type", claim.GetType(),
			"id", types.GetAttestationKey(claim.ChainID(), claim.GetEventNonce(), hash),
			"nonce", fmt.Sprint(claim.GetEventNonce()),
		)
	} else { // no error
		a.keeper.logger(ctx).Info("SendToCosmos to local gravity receiver", "ethSender", claim.EthereumSender,
			"receiver", receiver, "denom", coin.Denom, "amount", coin.Amount.String(), "nonce", claim.EventNonce,
			"ethContract", claim.TokenContract, "ethBlockHeight", claim.EthBlockHeight,
			"cosmosBlockHeight", ctx.BlockHeight(),
		)
		if err := ctx.EventManager().EmitTypedEvent(&types.EventSendToCosmosLocal{
			Nonce:    fmt.Sprint(claim.EventNonce),
			Receiver: receiver.String(),
			Token:    coin.Denom,
			Amount:   coin.Amount.String(),
		}); err != nil {
			return err
		}
	}

	return err // returns nil if no error
}

func (a AttestationHandler) sendCoinToCosmosAccount(
	ctx sdk.Context, claim types.MsgSendToCosmosClaim, receiver sdk.AccAddress, coin sdk.Coin,
) (ibcForwardQueued bool, err error) {
	return false, a.sendCoinToLocalAddress(ctx, claim, receiver, coin)
}