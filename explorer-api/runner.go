package explorer_api

import (
	"context"
	"github.com/cosmos/cosmos-sdk/server/api"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/gogo/gateway"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/ignite/cli/ignite/pkg/cosmosclient"
	"github.com/ignite/cli/ignite/pkg/cosmostxcollector/adapter/postgres"
	"github.com/mineplex/mineplex-chain/app"
	"github.com/mineplex/mineplex-chain/explorer-api/docs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
)

func RunRest(client cosmosclient.Client) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	router := mux.NewRouter()

	marshalerOption := &gateway.JSONPb{
		EmitDefaults: true,
		Indent:       "  ",
		OrigName:     true,
		AnyResolver:  client.Context().InterfaceRegistry,
	}

	grpcMux := runtime.NewServeMux(
		// Custom marshaler option is required for gogo proto
		runtime.WithMarshalerOption(runtime.MIMEWildcard, marshalerOption),

		// This is necessary to get error details properly
		// marshalled in unary requests.
		runtime.WithProtoErrorHandler(runtime.DefaultHTTPProtoErrorHandler),

		// Custom header matcher for mapping request headers to
		// GRPC metadata
		runtime.WithIncomingHeaderMatcher(api.CustomGRPCHeaderMatcher),
	)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := RegisterQueryHandlerFromEndpoint(ctx, grpcMux, "localhost:12201", opts)
	if err != nil {
		panic(err)
	}

	docs.RegisterOpenAPIService(app.Name, router)
	router.PathPrefix("/").Handler(grpcMux)

	log.Printf("server listening at 8081")
	if err := http.ListenAndServe(":8081", router); err != nil {
		panic(err)
	}
}

func RunGrpc(client cosmosclient.Client, db postgres.Adapter) {
	lis, err := net.Listen("tcp", ":12201")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	RegisterQueryServer(s, &server{
		client:                  client,
		accountRetriever:        authtypes.AccountRetriever{},
		bankQueryClient:         banktypes.NewQueryClient(client.Context()),
		stakingQueryClient:      stakingtypes.NewQueryClient(client.Context()),
		distributionQueryClient: distrtypes.NewQueryClient(client.Context()),
		txQueryClient:           txtypes.NewServiceClient(client.Context()),
		db:                      db,
	})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}
