package service

import (
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"github.com/rs/zerolog"
	"log"
)

// Сервис для кеширования в Redis
type CacheService struct {
	redisConn redis.Conn
	logger    zerolog.Logger
}

// Конструктор
func NewCacheService(
	logger zerolog.Logger,
	redisConn redis.Conn,
) *CacheService {
	return &CacheService{
		logger:    logger,
		redisConn: redisConn,
	}
}

// Set Chache Value for Tag
func (s CacheService) SetByTag(tag string, value interface{}) {
	jsonValue, err := json.Marshal(value)
	_, err = s.redisConn.Do("HMSET", tag, "value", jsonValue)
	if err != nil {
		log.Fatal(err)
	}
}

// Get Chache Value for Tag
func (s CacheService) GetByTag(tag string, v interface{}) (result bool) {
	exists, err := redis.Bool(s.redisConn.Do("EXISTS", tag))
	if err != nil {
		log.Println(err)
		return false
	}
	if exists == false {
		return false
	}

	value, err := redis.String(s.redisConn.Do("HGET", tag, "value"))
	if err != nil {
		log.Println(err)
		return false
	}
	err = json.Unmarshal([]byte(value), v)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

// Delete Chache Value for Tag
func (s CacheService) DeleteByTag(tag string) {
	_, err := redis.Bool(s.redisConn.Do("DEL", tag))
	if err != nil {
		log.Println(err)
	}
}
