package cache

import (
	"sync"
	"time"

	"github.com/uber-go/zap"

	"github.com/lomik/go-carbon/points"
)

type WriteoutQueue struct {
	sync.RWMutex
	cache *Cache

	// Writeout queue. Usage:
	// q := <- queue
	// p := cache.Pop(q.Metric)
	queue   chan *points.Points
	rebuild func(abort chan bool) chan bool // return chan waiting for complete
	logger  zap.Logger
}

func NewWriteoutQueue(cache *Cache) *WriteoutQueue {
	q := &WriteoutQueue{
		cache:  cache,
		queue:  nil,
		logger: cache.logger,
	}
	q.rebuild = q.makeRebuildCallback(time.Time{})
	return q
}

func (q *WriteoutQueue) makeRebuildCallback(nextRebuildTime time.Time) func(chan bool) chan bool {
	var nextRebuildOnce sync.Once
	nextRebuildComplete := make(chan bool)

	nextRebuild := func(abort chan bool) chan bool {
		// next rebuild
		nextRebuildOnce.Do(func() {
			now := time.Now()
			q.logger.Debug("WriteoutQueue.nextRebuildOnce.Do",
				zap.String("now", now.String()),
				zap.String("next", nextRebuildTime.String()),
			)
			if now.Before(nextRebuildTime) {
				sleepTime := nextRebuildTime.Sub(now)
				q.logger.Debug("WriteoutQueue sleep before rebuild",
					zap.String("sleepTime", sleepTime.String()),
				)

				select {
				case <-time.After(sleepTime):
					// pass
				case <-abort:
					// pass
				}
			}
			q.update()
			close(nextRebuildComplete)
		})

		return nextRebuildComplete
	}

	return nextRebuild
}

func (q *WriteoutQueue) update() {
	queue := q.cache.makeQueue()

	q.Lock()
	q.queue = queue
	q.rebuild = q.makeRebuildCallback(time.Now().Add(100 * time.Millisecond))
	q.Unlock()
}

func (q *WriteoutQueue) get(abort chan bool, pop func(key string) (p *points.Points, exists bool)) *points.Points {
QueueLoop:
	for {
		q.RLock()
		queue := q.queue
		rebuild := q.rebuild
		q.RUnlock()

	FetchLoop:
		for {
			select {
			case qp := <-queue:
				// pop from cache
				if p, exists := pop(qp.Metric); exists {
					return p
				}
				continue FetchLoop
			case <-abort:
				return nil
			default:
				// queue is empty, create new
				select {
				case <-rebuild(abort):
					// wait for rebuild
					continue QueueLoop
				case <-abort:
					return nil
				}
			}
		}
	}
}

func (q *WriteoutQueue) Get(abort chan bool) *points.Points {
	return q.get(abort, q.cache.Pop)
}

func (q *WriteoutQueue) GetNotConfirmed(abort chan bool) *points.Points {
	return q.get(abort, q.cache.PopNotConfirmed)
}
