package txcollector

import (
	"context"
	"github.com/ignite/cli/ignite/pkg/cosmosclient"
	"github.com/ignite/cli/ignite/pkg/cosmostxcollector"
	"github.com/ignite/cli/ignite/pkg/cosmostxcollector/adapter/postgres"
)

func Collect(ctx context.Context, db postgres.Adapter, client cosmosclient.Client) error {
	// Make sure that the data backend schema is up-to-date
	if err := db.Init(ctx); err != nil {
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
