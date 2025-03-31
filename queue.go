package tools

import (
	"fmt"
	"time"

	"github.com/julianto0911/redismq"
)

type StatsQueue struct {
	BQueue       *redismq.Queue
	BConsumer    *redismq.Consumer
	DataQueue    *redismq.Queue
	DataConsumer *redismq.Consumer
}

func NewStatsQueue(rds RedisConfiguration) (*StatsQueue, error) {
	obj := StatsQueue{}

	var err error

	//add first statistics queue components
	obj.BQueue, err = redismq.SelectQueue(rds.Host, rds.Port, rds.Password, 9, rds.Prefix+"_b_queue")
	if err != nil {
		obj.BQueue = redismq.CreateQueue(rds.Host, rds.Port, rds.Password, 9, rds.Prefix+"_b_queue")
	}

	//sleep 1 second to avoid test fail
	time.Sleep(time.Second)
	name := "_b_reader_" + ShortUUID()

	obj.BConsumer, err = obj.BQueue.AddConsumer(rds.Prefix + name)
	if err != nil {
		return nil, fmt.Errorf("fail add consumer for b queue : %w", err)
	}

	//add second statistics queue components
	obj.DataQueue, err = redismq.SelectQueue(rds.Host, rds.Port, rds.Password, 9, rds.Prefix+"_data_queue")
	if err != nil {
		obj.DataQueue = redismq.CreateQueue(rds.Host, rds.Port, rds.Password, 9, rds.Prefix+"_data_queue")
	}

	//sleep 1 second to avoid test fail
	time.Sleep(time.Second)
	name = "_data_reader_" + ShortUUID()

	obj.DataConsumer, err = obj.DataQueue.AddConsumer(rds.Prefix + name)
	if err != nil {
		return nil, fmt.Errorf("fail add consumer for data feeder : %w", err)
	}

	return &obj, nil
}
func (s *StatsQueue) Close() {
	defer s.BConsumer.Quit()
	defer s.DataConsumer.Quit()
}
