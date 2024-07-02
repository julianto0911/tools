package redismq

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Consumer are used for reading from queues
type Consumer struct {
	Name  string
	Queue *Queue

	cancel         context.CancelFunc
	contextCleared <-chan struct{}
}

// Get returns a single package from the queue (blocking)
func (consumer *Consumer) Get() (*Package, error) {
	if consumer.HasUnacked() {
		return nil, fmt.Errorf("unacked Packages found")
	}
	return consumer.unsafeGet()
}

// NoWaitGet returns a single package from the queue (returns nil, nil if no package in queue)
func (consumer *Consumer) NoWaitGet() (*Package, error) {
	if consumer.HasUnacked() {
		return nil, fmt.Errorf("unacked Packages found")
	}
	answer := consumer.Queue.redisClient.RPopLPush(ctx,
		queueInputKey(consumer.Queue.Name),
		consumerWorkingQueueKey(consumer.Queue.Name, consumer.Name),
	)
	if answer.Val() == "" {
		return nil, nil
	}
	consumer.Queue.incrRate(
		consumerWorkingRateKey(consumer.Queue.Name, consumer.Name),
		1,
	)
	return consumer.parseRedisAnswer(answer)
}

// GetUnacked returns a single packages from the working queue of this consumer
func (consumer *Consumer) GetUnacked() (*Package, error) {
	if !consumer.HasUnacked() {
		return nil, fmt.Errorf("no unacked Packages found")
	}
	answer := consumer.Queue.redisClient.LIndex(ctx, consumerWorkingQueueKey(consumer.Queue.Name, consumer.Name), -1)
	return consumer.parseRedisAnswer(answer)
}

// HasUnacked returns true if the consumers has unacked packages
func (consumer *Consumer) HasUnacked() bool {
	return consumer.GetUnackedLength() != 0
}

// GetUnackedLength returns the number of packages in the unacked queue
func (consumer *Consumer) GetUnackedLength() int64 {
	return consumer.Queue.redisClient.LLen(ctx, consumerWorkingQueueKey(consumer.Queue.Name, consumer.Name)).Val()
}

// GetFailed returns a single packages from the failed queue of this consumer
func (consumer *Consumer) GetFailed() (*Package, error) {
	answer := consumer.Queue.redisClient.RPopLPush(ctx, queueFailedKey(consumer.Queue.Name), consumerWorkingQueueKey(consumer.Queue.Name, consumer.Name))
	consumer.Queue.incrRate(consumerWorkingRateKey(consumer.Queue.Name, consumer.Name), 1)
	return consumer.parseRedisAnswer(answer)
}

// ResetWorking deletes! all messages in the working queue of this consumer
func (consumer *Consumer) ResetWorking() error {
	return consumer.Queue.redisClient.Del(ctx, consumerWorkingQueueKey(consumer.Queue.Name, consumer.Name)).Err()
}

// RequeueWorking requeues all packages from working to input
func (consumer *Consumer) RequeueWorking() error {
	for consumer.HasUnacked() {
		p, err := consumer.GetUnacked()
		if err != nil {
			return err
		}
		p.Requeue()
	}
	return nil
}

func (consumer *Consumer) ackPackage() error {
	return consumer.Queue.redisClient.RPop(ctx, consumerWorkingQueueKey(consumer.Queue.Name, consumer.Name)).Err()
}

func (consumer *Consumer) requeuePackage() error {
	answer := consumer.Queue.redisClient.RPopLPush(ctx, consumerWorkingQueueKey(consumer.Queue.Name, consumer.Name), queueInputKey(consumer.Queue.Name))
	consumer.Queue.incrRate(queueInputRateKey(consumer.Queue.Name), 1)
	return answer.Err()
}

func (consumer *Consumer) failPackage() error {
	return consumer.Queue.redisClient.RPopLPush(ctx, consumerWorkingQueueKey(consumer.Queue.Name, consumer.Name), queueFailedKey(consumer.Queue.Name)).Err()
}

func (consumer *Consumer) startHeartbeat() {
	firstWrite := make(chan struct{}, 1)

	ctx, cancel := context.WithCancel(context.Background())
	consumer.cancel = cancel

	waitForClear := make(chan struct{}, 1)
	consumer.contextCleared = waitForClear

	go func() {
		firstRun := true
		for {
			consumer.Queue.redisClient.Set(ctx,
				consumerHeartbeatKey(consumer.Queue.Name, consumer.Name),
				"ping",
				time.Second,
			)
			if firstRun {
				// use close instead
				close(firstWrite)
				firstRun = false
			}
			select {
			case <-time.After(500 * time.Millisecond):
			case <-ctx.Done():
				// remove heart beat immediately
				consumer.Queue.redisClient.Del(ctx, consumerHeartbeatKey(consumer.Queue.Name, consumer.Name))
				close(waitForClear)
				return
			}
		}
	}()
	<-firstWrite
}

func (consumer *Consumer) Quit() {
	if consumer.cancel == nil {
		return
	}

	consumer.cancel()
	// wait until heart beat mark is removed
	<-consumer.contextCleared

	consumer.cancel = nil
}

func (consumer *Consumer) parseRedisAnswer(answer *redis.StringCmd) (*Package, error) {
	if answer.Err() != nil {
		return nil, answer.Err()
	}
	p, err := unmarshalPackage(answer.Val(), consumer.Queue, consumer)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (consumer *Consumer) unsafeGet() (*Package, error) {
	answer := consumer.Queue.redisClient.BRPopLPush(ctx,
		queueInputKey(consumer.Queue.Name),
		consumerWorkingQueueKey(consumer.Queue.Name, consumer.Name),
		0,
	)
	consumer.Queue.incrRate(
		consumerWorkingRateKey(consumer.Queue.Name, consumer.Name),
		1,
	)
	return consumer.parseRedisAnswer(answer)
}
