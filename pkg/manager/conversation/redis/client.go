package conn_redis

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"math/big"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type Config struct {
	Addr       string `yaml:"addr" env:"REDIS_ADDR" env-default:"localhost:6379"`
	Password   string `yaml:"password" env:"REDIS_PASSWORD" env-default:""`
	GroupCount int    `yaml:"groupCount" env:"REDIS_GROUP_COUNT" env-default:"1"`
	Env        string `yaml:"env" env:"ENV" env-default:"local"`
}
type RedisClient interface {
	AddMessage(ctx context.Context, groupID int, conversationID uuid.UUID, message interface{}) error
	Handle(ctx context.Context, callback RedisMessageHandler, handlerName string) error
	GetGroupID(conversationID uuid.UUID) int
}

type redisClient struct {
	client   *redis.Client
	cfg      Config
	groupIds []int
}

func Init(ctx context.Context, cfg Config) (RedisClient, error) {
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	if cfg.Env == "local" {
		tlsConfig = nil
	}

	conn := redis.NewClient(&redis.Options{
		Addr:      cfg.Addr,
		Password:  cfg.Password,
		TLSConfig: tlsConfig,
	})

	pong, err := conn.Ping(ctx).Result()
	if err != nil {
		log.Printf("failed to connect to redis: %v\n", err)
		return nil, err
	}
	log.Printf("connected to redis: %v\n", pong)

	var groupIds []int = make([]int, cfg.GroupCount)
	for i := 0; i < cfg.GroupCount; i++ {
		groupIds[i] = i
	}

	redisClient := redisClient{
		client:   conn,
		cfg:      cfg,
		groupIds: groupIds,
	}
	return &redisClient, nil
}

func (c *redisClient) AddMessage(ctx context.Context, groupID int, conversationID uuid.UUID, message interface{}) error {
	// check  if group exists
	if !contains(c.groupIds, groupID) {
		return fmt.Errorf("group with id: %+v does not exist", groupID)
	}

	// publish message
	key := fmt.Sprintf("group%+v", groupID)
	fmt.Printf("publishing message: %+v\n", message)
	if err := c.client.Publish(ctx, key, message).Err(); err != nil {
		return err
	}

	return nil
}

func (c *redisClient) Handle(ctx context.Context, callback RedisMessageHandler, handlerName string) error {
	fmt.Printf("subs for handling redis messages is running\n")

	var wg sync.WaitGroup

	var subs []*redis.PubSub = make([]*redis.PubSub, c.cfg.GroupCount)

	for i, groupId := range c.groupIds {
		key := fmt.Sprintf("group%+v", groupId)
		sub := c.client.Subscribe(ctx, key)
		subs[i] = sub

		wg.Add(1)
		fmt.Printf("running sub with groupId: %+v\n", key)

		go func(sub *redis.PubSub) {
			defer wg.Done()
			ch := sub.Channel()
			for msg := range ch {
				fmt.Printf("received message from redis: %+v, name: %+v\n", msg, handlerName)
				if err := callback(ctx, msg.Payload); err != nil {
					// TODO: handle error properly
					log.Printf("failed to handle message: %v\n", err)

					return
				}
			}

		}(sub)
	}

	wg.Wait()

	for _, sub := range subs {
		if err := sub.Close(); err != nil {
			fmt.Printf("failed to close sub: %+v\n", err)
		}
	}
	return nil
}

func (c *redisClient) GetGroupID(conversationID uuid.UUID) int {
	intValue := new(big.Int)
	intValue.SetBytes(conversationID[:])

	group := new(big.Int)
	group.Mod(intValue, big.NewInt(int64(c.cfg.GroupCount)))

	return int(group.Int64())
}

type RedisMessageHandler func(ctx context.Context, mes string) error

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
