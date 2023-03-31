package main

import (
	"context"
	"github.com/ignite/cli/ignite/pkg/clictx"
	"github.com/ignite/cli/ignite/pkg/cosmosclient"
	"github.com/ignite/cli/ignite/pkg/cosmostxcollector/adapter/postgres"
	"github.com/mineplex/mineplex-chain/app"
	explorer_api "github.com/mineplex/mineplex-chain/explorer-api"
	"github.com/mineplex/mineplex-chain/txcollector"
	"os"
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
	db, err := postgres.NewAdapter(dbName, postgres.WithParams(params), postgres.WithUser(os.Getenv("PG_USER")), postgres.WithPassword(os.Getenv("PG_PASSWORD")))
	if err != nil {
		panic(err)
	}

	// Init the Cosmos client
	client, err := cosmosclient.New(ctx, cosmosclient.WithNodeAddress(rpcAddr), cosmosclient.WithAddressPrefix(app.AccountAddressPrefix))
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
