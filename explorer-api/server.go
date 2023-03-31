package explorer_api

import (
	"context"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ignite/cli/ignite/pkg/cosmosclient"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

type server struct {
	client           cosmosclient.Client
	accountRetriever client.AccountRetriever
	bankQueryClient  banktypes.QueryClient
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
		Result: &tx.TxResult,
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

	return &AddressResponse{
		Coins:    coins,
		Number:   number,
		Sequence: sequence,
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
