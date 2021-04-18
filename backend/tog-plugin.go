package main

import (
	"context"
	"database/sql"

	"github.com/Nikola-Milovic/tog-plugin/game/match"
	"github.com/Nikola-Milovic/tog-plugin/startup"
	"github.com/heroiclabs/nakama-common/runtime"
)

// InitModule used when first starting the server
func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {

	logger.Info("Loaded test plugin!")

	if err := initializer.RegisterMatch("tog", func(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule) (runtime.Match, error) {
		return &match.Match{}, nil
	}); err != nil {
		return err
	}

	// Register as matchmaker matched hook, this call should be in InitModule.
	if err := initializer.RegisterMatchmakerMatched(startup.MakeMatch); err != nil {
		logger.Error("Unable to register: %v", err)
		return err
	}

	match.RegisterMatchRPC(initializer)

	startup.StartUp(false)

	return nil
}

func main(){}
