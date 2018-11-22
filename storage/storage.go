package storage

import (
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/petaki/probe/config"
)

// Current instance.
var Current = Storage{}

// Storage type.
type Storage struct {
	Pool *redis.Pool
}

// Setup function.
func (s *Storage) Setup() {
	s.Pool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			options := []redis.DialOption{
				redis.DialDatabase(config.Current.RedisDatabase),
			}

			if config.Current.RedisPassword != "" {
				options = append(options, redis.DialPassword(config.Current.RedisPassword))
			}

			return redis.Dial("tcp", config.Current.RedisHost, options...)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}

			_, err := c.Do("PING")

			return err
		},
	}
}