package main

import (
	"context"
	"github.com/ignite/cli/ignite/pkg/clictx"
	"github.com/ignite/cli/ignite/pkg/cosmosclient"
	"github.com/ignite/cli/ignite/pkg/cosmostxcollector"
	"github.com/ignite/cli/ignite/pkg/cosmostxcollector/adapter/postgres"
)

// todo: read from config
const (
	// Name of a local PostgreSQL database
	dbName = "mineplex"

	// Cosmos RPC address
	rpcAddr = "http://0.0.0.0:26657"
)

func collect(ctx context.Context, db postgres.Adapter) error {
	// Make sure that the data backend schema is up-to-date
	if err := db.Init(ctx); err != nil {
		return err
	}

	// Init the Cosmos client
	client, err := cosmosclient.New(ctx, cosmosclient.WithNodeAddress(rpcAddr))
	if err != nil {
		return err
	}

	// Get the latest block height
	latestHeight, err := client.LatestBlockHeight(ctx)
	if err != nil {
		return err
	}

	// Collect transactions and events starting from a block height.
	// The collector stops at the latest height available at the time of the call.
	collector := cosmostxcollector.New(db, client)
	if err := collector.Collect(ctx, latestHeight); err != nil {
		return err
	}

	return nil
}

func main() {
	ctx := clictx.From(context.Background())

	// Init an adapter for a local PostgreSQL database running with the default values
	params := map[string]string{"sslmode": "disable"}
	db, err := postgres.NewAdapter(dbName, postgres.WithParams(params))
	if err != nil {
		panic(err)
	}

	if err := collect(ctx, db); err != nil {
		panic(err)
	}
}
