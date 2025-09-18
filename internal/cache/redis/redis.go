package redis

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Komilov31/delayed-notifier/internal/config"
	"github.com/Komilov31/delayed-notifier/internal/repository"
	"github.com/wb-go/wbf/redis"
)

type Redis struct {
	client redis.Client
}

func New() *Redis {
	password := os.Getenv("REDIS_PASSWORD")

	client := redis.New(
		config.Cfg.Redis.Host+config.Cfg.Redis.Port,
		password,
		0,
	)

	return &Redis{
		client: *client,
	}
}

func (r *Redis) Get(key string) (string, error) {
	return r.client.Get(context.Background(), key)
}

func (r *Redis) Set(key int, value interface{}) error {
	id := strconv.Itoa(key)
	return r.client.SetEX(context.Background(), id, value, time.Hour*24).Err()
}

func (r *Redis) LoadNotifications(repo *repository.Repository) error {
	notifications, err := repo.GetAllNotifications()
	if err != nil {
		log.Fatal("could not get notifications from db: ", err)
	}

	for _, notif := range notifications {
		if err := r.Set(notif.Id, notif.Status); err != nil {
			return err
		}
	}

	return nil
}
