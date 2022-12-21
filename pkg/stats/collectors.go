package stats

import (
	"github.com/grippenet/user-stats-service/pkg/db"
	"github.com/grippenet/user-stats-service/pkg/types"
)

type StatCollector interface {
	Fetch(db *db.UserDBService, instanceID string, filter types.StatFilter) (types.Counter, error)
}

type UserStatCollector struct {
}

func (u *UserStatCollector) Fetch(dbService *db.UserDBService, instanceID string, filter types.StatFilter) (types.Counter, error) {
	counter := types.Counter{Name: "users_count"}
	count, err := dbService.CountUser(instanceID, filter, db.UserOptions{})
	if err != nil {
		return counter, err
	}
	counter.Count = count
	return counter, nil
}

type UserActiveCollector struct {
}

func (u *UserActiveCollector) Fetch(dbService *db.UserDBService, instanceID string, filter types.StatFilter) (types.Counter, error) {
	counter := types.Counter{Name: "users_active"}
	count, err := dbService.CountUser(instanceID, filter, db.UserOptions{ActiveAccount: true})
	if err != nil {
		return counter, err
	}
	counter.Count = count
	return counter, nil
}
