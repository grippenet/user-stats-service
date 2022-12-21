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

func userCounter(name string, dbService *db.UserDBService, instanceID string, filter types.StatFilter, userOptions db.UserOptions) (types.Counter, error) {
	counter := types.Counter{Name: name}
	count, err := dbService.CountUser(instanceID, filter, userOptions)
	if err != nil {
		return counter, err
	}
	counter.Value = SimpleCounter{Count: count}
	return counter, nil
}

func (u *UserStatCollector) Fetch(dbService *db.UserDBService, instanceID string, filter types.StatFilter) (types.Counter, error) {
	return userCounter("users_count", dbService, instanceID, filter, db.UserOptions{})
}

type UserActiveCollector struct {
}

func (u *UserActiveCollector) Fetch(dbService *db.UserDBService, instanceID string, filter types.StatFilter) (types.Counter, error) {
	return userCounter("users_active", dbService, instanceID, filter, db.UserOptions{ActiveAccount: true})
}

// UserWeeklySubscribersCollector count Active users with weekly subscription
type UserWeeklySubscribersCollector struct {
}

func (u *UserWeeklySubscribersCollector) Fetch(dbService *db.UserDBService, instanceID string, filter types.StatFilter) (types.Counter, error) {
	return userCounter("users_weekly", dbService, instanceID, filter, db.UserOptions{ActiveAccount: true, SubscribedToWeekly: true})
}

type UserWeekDayCollector struct {
}

func (u *UserWeekDayCollector) Fetch(dbService *db.UserDBService, instanceID string, filter types.StatFilter) (types.Counter, error) {
	counter := types.Counter{Name: "users_weekday"}

	counts, err := dbService.WeekDayReminders(instanceID, filter, db.UserOptions{})
	if err != nil {
		return counter, err
	}
	counter.Value = MapCounter{Counts: counts}
	return counter, nil
}
