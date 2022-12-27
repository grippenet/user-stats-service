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
	counter := types.Counter{Name: name, Type: types.COUNTER_TYPE_COUNT}
	count, err := dbService.CountUser(instanceID, filter, userOptions)
	if err != nil {
		return counter, err
	}
	counter.Value = count
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

type UserNewsletterSubscribersCollector struct {
}

func (u *UserNewsletterSubscribersCollector) Fetch(dbService *db.UserDBService, instanceID string, filter types.StatFilter) (types.Counter, error) {
	return userCounter("users_newsletter", dbService, instanceID, filter, db.UserOptions{ActiveAccount: true, SubscribedToNewsletter: true})
}

type UserWeekDayCollector struct {
}

func (u *UserWeekDayCollector) Fetch(dbService *db.UserDBService, instanceID string, filter types.StatFilter) (types.Counter, error) {
	counter := types.Counter{Name: "users_weekday", Type: types.COUNTER_TYPE_MAP}

	counts, err := dbService.WeekDayReminders(instanceID, filter, db.UserOptions{})
	if err != nil {
		return counter, err
	}
	counter.Value = counts
	return counter, nil
}

func profilesCounters(name string, dbService *db.UserDBService, instanceID string, filter types.StatFilter, userOptions db.UserOptions) (types.Counter, error) {
	counter := types.Counter{Name: name, Type: types.COUNTER_TYPE_COUNT}
	counts, err := dbService.CountProfiles(instanceID, filter, userOptions)
	if err != nil {
		return counter, err
	}
	counter.Value = counts
	return counter, nil
}

type UserProfilesCollector struct {
}

func (u *UserProfilesCollector) Fetch(dbService *db.UserDBService, instanceID string, filter types.StatFilter) (types.Counter, error) {
	return profilesCounters("profiles_count", dbService, instanceID, filter, db.UserOptions{})
}

type ActiveProfilesCollector struct {
}

func (u *ActiveProfilesCollector) Fetch(dbService *db.UserDBService, instanceID string, filter types.StatFilter) (types.Counter, error) {
	return profilesCounters("profiles_active", dbService, instanceID, filter, db.UserOptions{ActiveAccount: true})
}

type ProfilesHistogramCollector struct {
}

func (u *ProfilesHistogramCollector) Fetch(dbService *db.UserDBService, instanceID string, filter types.StatFilter) (types.Counter, error) {
	counter := types.Counter{Name: "profiles_histogram", Type: types.COUNTER_TYPE_MAP}

	counts, err := dbService.ProfilesCount(instanceID, filter, db.UserOptions{ActiveAccount: true})
	if err != nil {
		return counter, err
	}
	counter.Value = counts
	return counter, nil
}
