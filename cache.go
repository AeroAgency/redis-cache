package service

import (
	"encoding/json"
	"fmt"
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

// Set Cache Value for Tag
func (s CacheService) SetByTag(tag string, value interface{}, expire int) {
	jsonValue, err := json.Marshal(value)
	_, err = s.redisConn.Do("HMSET", tag, "value", jsonValue)
	if err != nil {
		log.Fatal(err)
	}
	if expire != 0 {
		_, err = s.redisConn.Do("EXPIRE", tag, expire)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// Get Cache Value for Tag
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

// Delete Cache Value for Tag
func (s CacheService) DeleteByTag(tag string) {
	_, err := redis.Bool(s.redisConn.Do("DEL", tag))
	if err != nil {
		log.Println(err)
	}
}

// Delete Cache Value for Tag
func (s CacheService) DeleteTagsByPattern(pattern string) error {
	iter := 0
	keys := []string{}
	for {
		arr, err := redis.Values(s.redisConn.Do("SCAN", iter, "MATCH", pattern))
		if err != nil {
			return fmt.Errorf("error retrieving '%s' keys", pattern)
		}

		iter, _ = redis.Int(arr[0], nil)
		k, _ := redis.Strings(arr[1], nil)
		keys = append(keys, k...)

		if iter == 0 {
			break
		}
	}
	for _, key := range keys {
		s.DeleteByTag(key)
	}
	return nil
}
