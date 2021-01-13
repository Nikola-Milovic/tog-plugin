package match

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/Nikola-Milovic/tog-plugin/startup"
	"github.com/heroiclabs/nakama-common/runtime"
)

type SaveUserDraftMessage struct {
	UserID  string         `json:"user_id"`
	Reserve map[string]int `json:"reserve"`
	Army    map[string]int `json:"army"`
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

type DraftMessage struct {
	Units    []DraftUnitMessage `json:"units"`
	MaxUnits int                `json:"max_units"`
	MaxArmy  int                `json:"max_army"`
}

type DraftUnitMessage struct {
	ID     string `json:"unit_id"`
	Amount int    `json:"amount"`
}

func StartDraft(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	userPayload := make(map[string]string)

	err := json.Unmarshal([]byte(payload), &userPayload)

	draftData := DraftMessage{}

	//userRun, err := getUserRun(userPayload["user_id"], ctx, logger, nk)

	unitMessage := make([]DraftUnitMessage, 0, 20)
	//	s1 := rand.NewSource(time.Now().UnixNano())
	//r1 := rand.New(s1)
	for id, _ := range startup.UnitDataMap {
		//_ := r1.Intn(5)
		unitMessage = append(unitMessage, DraftUnitMessage{ID: id, Amount: 3})
	}

	draftData.Units = unitMessage
	//draftData.MaxUnits = int(math.Max(float64(userRun.Floor), 1)) * 5
	//draftData.MaxArmy = int((float64(draftData.MaxUnits) * float64(80)) / float64(100))
	draftData.MaxArmy = 6
	draftData.MaxUnits = 8

	dataToSend, err := json.Marshal(&draftData)

	if err != nil {
		return "err", nil
	}

	return string(dataToSend), nil
}
