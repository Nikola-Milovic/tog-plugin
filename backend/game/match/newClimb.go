package match

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/startup"
	"github.com/heroiclabs/nakama-common/runtime"
)

func RegisterMatchRPC(initializer runtime.Initializer) error {
	if err := initializer.RegisterRpc("start_new_climb", StartNewClimb); err != nil {
		fmt.Printf("Unable to register start_new_climb : %v", err)
		return err
	}

	if err := initializer.RegisterRpc("save_user_draft", SaveUserDraft); err != nil {
		fmt.Printf("Unable to register save_user_draft : %v", err)
		return err
	}

	return nil
}

type ClimbScore struct {
	Losses int `json:"losses"`
	Wins   int `json:"wins"`
}

type CurrentRunData struct {
	Army    []map[string][]engine.Vector `json:"army"`
	Reserve map[string]int               `json:"reserve"`
	Floor   int                          `json:"floor"`
	Score   ClimbScore                   `json:"score"`
}

func StartNewClimb(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	logger.Info("Payload: %s", payload)

	userpayload := make(map[string]string, 1)
	err := json.Unmarshal([]byte(payload), &userpayload)

	if err != nil {
		return "err", err
	}

	runData := CurrentRunData{Floor: 0, Score: ClimbScore{Wins: 0, Losses: 0}}

	runDataJSON, err := json.Marshal(runData)

	if err != nil {
		return "err", err
	}

	userID := userpayload["user_id"]

	objects := []*runtime.StorageWrite{
		&runtime.StorageWrite{
			Collection:      "runs",
			Key:             "current_run",
			UserID:          userID,
			Value:           string(runDataJSON),
			PermissionRead:  1,
			PermissionWrite: 0,
		},
	}

	if _, err := nk.StorageWrite(ctx, objects); err != nil {
		return "err", err
	}

	draftData := GetNewClimbDraft()

	dataToSend, err := json.Marshal(&draftData)

	if err != nil {
		return "err", nil
	}

	return string(dataToSend), nil
}

type DraftUnitMessage struct {
	ID     string `json:"unit_id"`
	Amount int    `json:"amount"`
}

func GetNewClimbDraft() []DraftUnitMessage {
	unitMessage := make([]DraftUnitMessage, 0, 10)
	//	s1 := rand.NewSource(time.Now().UnixNano())
	//r1 := rand.New(s1)
	for id, _ := range startup.UnitDataMap {
		//_ := r1.Intn(5)
		unitMessage = append(unitMessage, DraftUnitMessage{ID: id, Amount: 5})
	}
	return unitMessage
}

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

	objectIds := []*runtime.StorageRead{
		&runtime.StorageRead{
			Collection: "runs",
			Key:        "current_run",
			UserID:     userID,
		},
	}

	objects, err := nk.StorageRead(ctx, objectIds)
	if err != nil {
		logger.Error("error reading from runs current_run", err.Error())
	}

	runData := CurrentRunData{}

	err = json.Unmarshal([]byte(objects[0].Value), &runData)
	if err != nil {
		logger.Error("error unmarshaling run data", err.Error())
	}

	runData.Reserve = saveUserDraftData.Reserve

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
