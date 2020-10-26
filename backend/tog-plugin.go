package main

import (
	"backend/src/control"
	"backend/src/rpc"
	"context"
	"database/sql"

	"github.com/heroiclabs/nakama-common/runtime"
)

// InitModule used when first starting the server
func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {

	logger.Info("Loaded test plugin!")

	if err := initializer.RegisterMatch("control", func(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule) (runtime.Match, error) {
		return &control.Match{}, nil
	}); err != nil {
		return err
	}

	if err := initializer.RegisterRpc("get_world_id", rpc.GetWorldId); err != nil {
		logger.Error("Unable to register: %v", err)
		return err
	}

	return nil
}
