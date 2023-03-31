package main

import (
	"context"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ignite/cli/ignite/pkg/clictx"
	"github.com/ignite/cli/ignite/pkg/cosmosclient"
	"github.com/ignite/cli/ignite/pkg/cosmostxcollector"
	"github.com/ignite/cli/ignite/pkg/cosmostxcollector/adapter/postgres"
	"github.com/mineplex/mineplex-chain/app"
	explorer_api "github.com/mineplex/mineplex-chain/explorer-api"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// todo: read from config
const (
	dbName  = "mineplex"
	rpcAddr = "http://0.0.0.0:26657"
)

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(0)
	}()

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

	// Make sure that the data backend schema is up-to-date
	if err := db.Init(ctx); err != nil {
		panic(err)
	}

	// Collect transactions and events starting from a block height.
	// The collector stops at the latest height available at the time of the call.
	collector := cosmostxcollector.New(db, client)

	// Get the latest block height
	latestHeight, err := client.LatestBlockHeight(ctx)
	if err != nil {
		panic(err)
	}

	for {
		if err := collector.Collect(ctx, latestHeight); err != nil {
			log.Error("Error while collecting new txs")
			continue
		}

		// Get the latest block height
		for {
			currentHeight, err := client.LatestBlockHeight(ctx)
			if err != nil {
				panic(err)
			}

			if currentHeight == latestHeight {
				time.Sleep(time.Millisecond * 500)
				continue
			}

			latestHeight = currentHeight
			break
		}
	}
}
