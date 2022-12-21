package db

import (
	"context"
	"fmt"
	"time"

	"github.com/grippenet/user-stats-service/pkg/types"
	"github.com/influenzanet/user-management-service/pkg/dbs/userdb"
	"github.com/influenzanet/user-management-service/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserDBService struct {
	*userdb.UserDBService
	timeout int
}

func NewUserDBService(configs models.DBConfig) *UserDBService {

	return &UserDBService{
		UserDBService: userdb.NewUserDBService(configs),
		timeout:       configs.Timeout,
	}
}

func (dbService *UserDBService) collectionRefUsers(instanceID string) *mongo.Collection {
	return dbService.DBClient.Database(dbService.DBNamePrefix + instanceID + "_users").Collection("users")
}

// DB utils
func (dbService *UserDBService) getContext() (ctx context.Context, cancel context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(dbService.timeout)*time.Second)
}

func filterField(field string, filter types.StatFilter) interface{} {

	if filter.From == 0 && filter.Until == 0 {
		return nil
	}

	var criteria interface{}

	if filter.From > 0 && filter.Until == 0 {
		criteria = bson.D{{"$gt", filter.From}}
	}

	if filter.Until > 0 && filter.From == 0 {
		criteria = bson.D{{"$lt", filter.Until}}
	}

	if filter.Until > 0 && filter.From > 0 {
		criteria = bson.M{"$and": bson.A{
			bson.D{{"$gt", filter.From}},
			bson.D{{"$lt", filter.Until}},
		},
		}
	}
	return bson.D{{field, criteria}}
}

type UserOptions struct {
	ActiveAccount      bool
	SubscribedToWeekly bool
}

func combineCriteria(cc []interface{}) interface{} {
	if len(cc) == 0 {
		return bson.D{}
	}
	if len(cc) == 1 {
		return cc[0]
	}
	a := make(bson.A, 0, len(cc))
	for _, c := range cc {
		a = append(a, c)
	}
	return bson.M{"$and": a}
}

func (svc *UserDBService) CountUser(instanceID string, filter types.StatFilter, opts UserOptions) (int64, error) {
	ctx, cancel := svc.getContext()
	defer cancel()

	users := svc.collectionRefUsers(instanceID)

	criteria := make([]interface{}, 0, 1)

	filters := filterField("timestamps.createdAt", filter)
	if filters != nil {
		criteria = append(criteria, filters)
	}

	if opts.ActiveAccount {
		criteria = append(criteria, bson.M{"account.accountConfirmedAt": bson.M{"$gt": 0}})
	}
	if opts.SubscribedToWeekly {
		criteria = append(criteria, bson.M{"contactPreferences.subscribedToWeekly": true})
	}

	cc := combineCriteria(criteria)

	count, err := users.CountDocuments(ctx, cc)

	if err != nil {
		return 0, err
	}
	return count, nil
}

func (svc *UserDBService) WeekDayReminders(instanceID string, filter types.StatFilter, opts UserOptions) (map[string]int64, error) {
	ctx, cancel := svc.getContext()
	defer cancel()

	users := svc.collectionRefUsers(instanceID)

	criteria := make([]interface{}, 0, 1)

	filters := filterField("timestamps.createdAt", filter)
	if filters != nil {
		criteria = append(criteria, filters)
	}

	criteria = append(criteria, bson.M{"account.accountConfirmedAt": bson.M{"$gt": 0}})
	criteria = append(criteria, bson.M{"contactPreferences.subscribedToWeekly": true})

	cc := combineCriteria(criteria)

	agg := bson.A{
		bson.M{"$match": cc},
		bson.M{"$group": bson.M{
			"_id":   "$contactPreferences.receiveWeeklyMessageDayOfWeek",
			"count": bson.M{"$sum": 1},
		},
		},
	}
	//[{'$group': { '_id': , 'count': {'$sum': 1}}}]

	cursor, err := users.Aggregate(ctx, agg)

	if err != nil {
		return nil, err
	}

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	counts := make(map[string]int64, len(results))

	for _, result := range results {
		fmt.Println(result)
		id := fmt.Sprint(result["_id"])
		count := result["count"].(int32)
		counts[id] = int64(count)
	}

	return counts, nil
}
