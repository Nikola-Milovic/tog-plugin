package match

import (
	"context"
	"encoding/json"

	"github.com/heroiclabs/nakama-common/runtime"
)

func getUserRun(userID string, ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule) (CurrentRunData, error) {
	objectIds := []*runtime.StorageRead{
		{
			Collection: "runs",
			Key:        "current_run",
			UserID:     userID,
				},
	}

	objects, err := nk.StorageRead(ctx, objectIds)
	if err != nil {
		logger.Error("error reading from runs current_run", err.Error())
		return CurrentRunData{}, err
	}

	runData := CurrentRunData{}

	err = json.Unmarshal([]byte(objects[0].Value), &runData)
	if err != nil {
		logger.Error("error unmarshaling run data", err.Error())
		return CurrentRunData{}, err
	}
	return runData, nil
}
