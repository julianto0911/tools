package tools

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/julianto0911/redismq"
)

type StatsQueue struct {
	BetQueue     *redismq.Queue
	BetConsumer  *redismq.Consumer
	DataQueue    *redismq.Queue
	DataConsumer *redismq.Consumer
}

func NewStatsQueue(rds RedisConfiguration) (*StatsQueue, error) {
	obj := StatsQueue{}

	var err error

	//add first statistics queue components
	obj.BetQueue, err = redismq.SelectQueue(rds.Host, rds.Port, rds.Password, 9, rds.Prefix+"_bet_queue")
	if err != nil {
		obj.BetQueue = redismq.CreateQueue(rds.Host, rds.Port, rds.Password, 9, rds.Prefix+"_bet_queue")
	}

	//sleep 1 second to avoid test fail
	time.Sleep(time.Second)
	id := uuid.New()
	name := "_bet_reader_" + id.String()

	obj.BetConsumer, err = obj.BetQueue.AddConsumer(rds.Prefix + name)
	if err != nil {
		return nil, fmt.Errorf("fail add consumer for bet queue : %w", err)
	}

	//add second statistics queue components
	obj.DataQueue, err = redismq.SelectQueue(rds.Host, rds.Port, rds.Password, 9, rds.Prefix+"_data_queue")
	if err != nil {
		obj.DataQueue = redismq.CreateQueue(rds.Host, rds.Port, rds.Password, 9, rds.Prefix+"_data_queue")
	}

	//sleep 1 second to avoid test fail
	time.Sleep(time.Second)
	id = uuid.New()
	name = "_data_reader_" + id.String()

	obj.DataConsumer, err = obj.DataQueue.AddConsumer(rds.Prefix + name)
	if err != nil {
		return nil, fmt.Errorf("fail add consumer for data feeder : %w", err)
	}

	return &obj, nil
}
func (s *StatsQueue) Close() {
	defer s.BetConsumer.Quit()
	defer s.DataConsumer.Quit()
}
