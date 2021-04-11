package match

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

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

	if err := initializer.RegisterRpc("draft", StartDraft); err != nil {
		fmt.Printf("Unable to register draft : %v", err)
		return err
	}

	return nil
}

type ClimbScore struct {
	Losses int `json:"losses"`
	Wins   int `json:"wins"`
}

type CurrentRunData struct {
	Army    map[string]int `json:"army"`
	Reserve map[string]int `json:"reserve"`
	Floor   int            `json:"floor"`
	Score   ClimbScore     `json:"score"`
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

	draftUnits := GetNewClimbDraft()

	draftData := DraftMessage{}
	draftData.Units = draftUnits
	draftData.MaxUnits = 5
	draftData.MaxArmy = 4

	dataToSend, err := json.Marshal(&draftData)

	if err != nil {
		return "err", nil
	}

	return string(dataToSend), nil
}

func GetNewClimbDraft() []DraftUnitMessage {
	unitMessage := make([]DraftUnitMessage, 0, 20)
	//	s1 := rand.NewSource(time.Now().UnixNano())
	//r1 := rand.New(s1)
	for id, _ := range startup.UnitDataMap {
		//_ := r1.Intn(5)
		unitMessage = append(unitMessage, DraftUnitMessage{ID: id, Amount: 5})
	}
	return unitMessage
}
