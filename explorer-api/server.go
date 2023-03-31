package explorer_api

import (
	"context"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ignite/cli/ignite/pkg/cosmosclient"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type server struct {
	client           cosmosclient.Client
	accountRetriever client.AccountRetriever
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

func RunRest() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := RegisterQueryHandlerFromEndpoint(ctx, mux, "localhost:12201", opts)
	if err != nil {
		panic(err)
	}
	log.Printf("server listening at 8081")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		panic(err)
	}
}

func RunGrpc(client cosmosclient.Client) {
	lis, err := net.Listen("tcp", ":12201")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	RegisterQueryServer(s, &server{
		client:           client,
		accountRetriever: authtypes.AccountRetriever{},
	})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}
