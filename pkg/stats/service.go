package stats

import (
	"fmt"

	"github.com/coneno/logger"
	"github.com/grippenet/user-stats-service/pkg/db"
	"github.com/grippenet/user-stats-service/pkg/types"
)

type StatsService struct {
	dbService  *db.UserDBService
	collectors []StatCollector
}

func NewStatService(dbService *db.UserDBService) *StatsService {

	collectors := []StatCollector{
		&UserStatCollector{},
		&UserActiveCollector{},
		&UserWeeklySubscribersCollector{},
		&UserWeekDayCollector{},
	}

	return &StatsService{dbService: dbService, collectors: collectors}
}

func (s *StatsService) Fetch(instanceID string, filter types.StatFilter) ([]types.Counter, error) {

	counters := make([]types.Counter, 0, len(s.collectors))

	for index, collector := range s.collectors {
		counter, err := collector.Fetch(s.dbService, instanceID, filter)
		if err != nil {
			logger.Error.Printf("Error for counter %d : %s", index, err)
			counter.Error = fmt.Sprint(err)
		}
		counters = append(counters, counter)
	}
	return counters, nil
}
