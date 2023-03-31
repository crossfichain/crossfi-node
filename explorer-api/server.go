package explorer_api

import (
	"context"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ignite/cli/ignite/pkg/cosmosclient"
	"github.com/ignite/cli/ignite/pkg/cosmostxcollector/adapter/postgres"
	"github.com/ignite/cli/ignite/pkg/cosmostxcollector/query"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

type server struct {
	client                  cosmosclient.Client
	accountRetriever        client.AccountRetriever
	bankQueryClient         banktypes.QueryClient
	stakingQueryClient      stakingtypes.QueryClient
	db                      postgres.Adapter
	distributionQueryClient distrtypes.QueryClient
}

func (s *server) Txs(ctx context.Context, request *TxsRequest) (*TxsResponse, error) {
	var txs []TxResponse

	panic("unimplemented")
	filters := []query.Filter{postgres.NewStringSliceFilter("from", []string{})} // todo

	q := query.NewEventQuery(query.AtPage(1), query.WithPageSize(100), query.WithFilters(filters...))

	events, err := s.db.QueryEvents(ctx, q)
	if err != nil {
		return nil, err
	}

	for _, event := range events {
		hash := common.HexToHash(event.TXHash).Bytes()
		tx, err := s.client.RPC.Tx(ctx, hash, false)
		if err != nil {
			return nil, err
		}

		decodedTx, err := s.client.Context().TxConfig.TxDecoder()(tx.Tx)
		jsonTx, err := s.client.Context().TxConfig.TxJSONEncoder()(decodedTx)

		txs = append(txs, TxResponse{
			Tx:     string(jsonTx),
			Height: tx.Height,
			Result: tx.TxResult,
		})
	}

	return &TxsResponse{Txs: txs}, nil
}

func (s *server) Validators(ctx context.Context, _ *ValidatorsRequest) (*ValidatorsResponse, error) {
	resp, err := s.stakingQueryClient.Validators(ctx, &stakingtypes.QueryValidatorsRequest{})
	if err != nil {
		return nil, err
	}

	return &ValidatorsResponse{
		Validators: resp.Validators,
	}, nil
}

func (s *server) Validator(ctx context.Context, request *ValidatorRequest) (*ValidatorResponse, error) {
	resp, err := s.stakingQueryClient.Validator(ctx, &stakingtypes.QueryValidatorRequest{
		ValidatorAddr: request.Address,
	})
	if err != nil {
		return nil, err
	}

	return &ValidatorResponse{
		Validator: resp.Validator,
	}, nil
}

func (s *server) Tx(ctx context.Context, request *TxRequest) (*TxResponse, error) {
	hash := common.HexToHash(request.Hash).Bytes()
	tx, err := s.client.RPC.Tx(ctx, hash, false)
	if err != nil {
		return nil, err
	}

	decodedTx, err := s.client.Context().TxConfig.TxDecoder()(tx.Tx)
	jsonTx, err := s.client.Context().TxConfig.TxJSONEncoder()(decodedTx)

	return &TxResponse{
		Tx:     string(jsonTx),
		Height: tx.Height,
		Result: tx.TxResult,
	}, nil
}

func (s *server) Status(ctx context.Context, _ *StatusRequest) (*StatusResponse, error) {
	status, err := s.client.Status(ctx)
	if err != nil {
		return nil, err
	}

	return &StatusResponse{
		LatestBlockHash:   status.SyncInfo.LatestBlockHash.String(),
		LatestBlockHeight: status.SyncInfo.LatestBlockHeight,
		LatestBlockTime:   status.SyncInfo.LatestBlockTime.String(),
	}, nil
}

func (s *server) Coins(ctx context.Context, _ *CoinsRequest) (*CoinsResponse, error) {
	resp, err := s.bankQueryClient.TotalSupply(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &CoinsResponse{
		Coins: resp.Supply,
	}, nil
}

func (s *server) Address(ctx context.Context, request *AddressRequest) (*AddressResponse, error) {
	coins, err := s.client.BankBalances(ctx, request.Address, nil)
	if err != nil {
		return nil, err
	}

	address, err := types.AccAddressFromBech32(request.Address)
	if err != nil {
		return nil, err
	}

	number, sequence, err := s.accountRetriever.GetAccountNumberSequence(s.client.Context(), address)
	if err != nil {
		return nil, err
	}

	delegations, err := s.stakingQueryClient.DelegatorDelegations(ctx, &stakingtypes.QueryDelegatorDelegationsRequest{
		DelegatorAddr: request.Address,
	})
	if err != nil {
		return nil, err
	}

	unbondingDelegations, err := s.stakingQueryClient.DelegatorUnbondingDelegations(ctx, &stakingtypes.QueryDelegatorUnbondingDelegationsRequest{
		DelegatorAddr: request.Address,
	})
	if err != nil {
		return nil, err
	}

	redelegations, err := s.stakingQueryClient.Redelegations(ctx, &stakingtypes.QueryRedelegationsRequest{
		DelegatorAddr: request.Address,
	})
	if err != nil {
		return nil, err
	}

	rewards, err := s.distributionQueryClient.DelegationTotalRewards(ctx, &distrtypes.QueryDelegationTotalRewardsRequest{
		DelegatorAddress: request.Address,
	})
	if err != nil {
		return nil, err
	}

	rewardsResponse := RewardsResponse{
		Total: rewards.Total,
	}

	for _, reward := range rewards.Rewards {
		rewardsResponse.Rewards = append(rewardsResponse.Rewards, DelegationDelegatorReward{
			ValidatorAddress: reward.ValidatorAddress,
			Reward:           reward.Reward,
		})
	}

	return &AddressResponse{
		Coins:                coins,
		Number:               number,
		Sequence:             sequence,
		Delegations:          delegations.DelegationResponses,
		UnbondingDelegations: unbondingDelegations.UnbondingResponses,
		Redelegations:        redelegations.RedelegationResponses,
		Rewards:              rewardsResponse,
	}, nil
}

func (s *server) Blocks(ctx context.Context, request *BlocksRequest) (*BlocksResponse, error) {
	response := &BlocksResponse{}

	for i := request.FromHeight; i < request.ToHeight; i++ {
		resp, err := s.client.RPC.Block(ctx, &i)
		if err != nil {
			return nil, err
		}

		block, err := resp.Block.ToProto()
		if err != nil {
			return nil, err
		}

		response.Blocks = append(response.Blocks, block)
	}

	return response, nil
}

func (s *server) Block(ctx context.Context, request *BlockRequest) (*tmproto.Block, error) {
	resp, err := s.client.RPC.Block(ctx, &request.Height)
	if err != nil {
		return nil, err
	}

	return resp.Block.ToProto()
}
