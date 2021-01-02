package match

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/heroiclabs/nakama-common/runtime"
)

type SaveUserDraftMessage struct {
	UserID  string                     `json:"user_id"`
	Army    map[string][]engine.Vector `json:"army"`
	Reserve map[string]int             `json:"reserve"`
}

func SaveUserDraft(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	logger.Info("Payload for save user draft is : %s", payload)

	// Get userpayload with user_id, reserve and army
	// reserve is string - int key value pairs of Knight 3 (how many of that unit)
	// army is UnitID - []Vector of where that unit is positioned
	saveUserDraftData := SaveUserDraftMessage{}
	err := json.Unmarshal([]byte(payload), &saveUserDraftData)

	if err != nil {
		logger.Error("Couldnt unmarshal SaveUserDraft Payload", err.Error())
		return "err", err
	}

	// Get the current run and update army and reserve
	userID := saveUserDraftData.UserID

	runData, _ := getUserRun(userID, ctx, logger, nk)

	runData.Reserve = saveUserDraftData.Reserve
	runData.Army = saveUserDraftData.Army

	//Write the user run data back
	runDataJSON, err := json.Marshal(runData)

	if err != nil {
		return "err", err
	}

	writeObjects := []*runtime.StorageWrite{
		&runtime.StorageWrite{
			Collection:      "runs",
			Key:             "current_run",
			UserID:          userID,
			Value:           string(runDataJSON),
			PermissionRead:  1,
			PermissionWrite: 0,
		},
	}

	if _, err := nk.StorageWrite(ctx, writeObjects); err != nil {
		logger.Error("Error writing the run data after altering it in save user draft", err.Error())
		return "err", err
	}

	return "ok", nil
}
