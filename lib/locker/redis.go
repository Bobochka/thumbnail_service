package locker

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	goRedis "github.com/garyburd/redigo/redis"
)

func newPool(server string) *goRedis.Pool {
	return &goRedis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (goRedis.Conn, error) {
			c, err := goRedis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c goRedis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func cleanupHook(pool *goRedis.Pool) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGKILL)
	go func() {
		if pool != nil {
			<-c
			pool.Close()
			os.Exit(0)
		}
	}()
}
