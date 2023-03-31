package main

import (
	"context"
	"github.com/ignite/cli/ignite/pkg/clictx"
	"github.com/ignite/cli/ignite/pkg/cosmosclient"
	"github.com/ignite/cli/ignite/pkg/cosmostxcollector/adapter/postgres"
	explorer_api "github.com/mineplex/mineplex-chain/explorer-api"
	"github.com/mineplex/mineplex-chain/txcollector"
)

// todo: read from config
const (
	dbName  = "mineplex"
	rpcAddr = "http://0.0.0.0:26657"
)

func main() {
	ctx := clictx.From(context.Background())

	// Init an adapter for a local PostgreSQL database running with the default values
	params := map[string]string{"sslmode": "disable"}
	db, err := postgres.NewAdapter(dbName, postgres.WithParams(params), postgres.WithUser("root"), postgres.WithPassword("123"))
	if err != nil {
		panic(err)
	}

	// Init the Cosmos client
	client, err := cosmosclient.New(ctx, cosmosclient.WithNodeAddress(rpcAddr))
	if err != nil {
		panic(err)
	}

	go explorer_api.RunGrpc(client, db)
	go explorer_api.RunRest(client)

	if err := txcollector.Collect(ctx, db, client); err != nil {
		panic(err)
	}

	select {}
}
