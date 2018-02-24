package locker

import (
	"time"

	"github.com/Bobochka/thumbnail_service/lib"
	"gopkg.in/redsync.v1"
)

type RedisLocker struct {
	*redsync.Redsync
}

func New(host string) (*RedisLocker, error) {
	pool := newPool(host)

	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("PING")
	if err != nil {
		return nil, err
	}

	cleanupHook(pool)

	return &RedisLocker{
		redsync.New([]redsync.Pool{pool}),
	}, nil
}

func (r *RedisLocker) NewMutex(name string) lib.Mutex {
	return r.Redsync.NewMutex(
		name,
		redsync.SetTries(3),
		redsync.SetExpiry(5*time.Second),
		redsync.SetRetryDelay(200*time.Millisecond),
	)
}
