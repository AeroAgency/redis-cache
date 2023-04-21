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
	redisPool redis.Pool
	logger    zerolog.Logger
}

// Конструктор
func NewCacheService(
	logger zerolog.Logger,
	redisPool redis.Pool,
) *CacheService {
	return &CacheService{
		logger:    logger,
		redisPool: redisPool,
	}
}

// Set Cache Value for Tag
func (s CacheService) SetByTag(tag string, value interface{}, expire int) {
	c := s.redisPool.Get()
	defer c.Close()

	jsonValue, err := json.Marshal(value)

	_, err = c.Do("HMSET", tag, "value", jsonValue)
	if err != nil {
		log.Println(err)
	}
	if expire != 0 {
		_, err = c.Do("EXPIRE", tag, expire)
		if err != nil {
			log.Println(err)
		}
	}
}

// Get Cache Value for Tag
func (s CacheService) GetByTag(tag string, v interface{}) (result bool) {
	c := s.redisPool.Get()
	defer c.Close()

	exists, err := redis.Bool(c.Do("EXISTS", tag))
	if err != nil {
		log.Println(err)
		return false
	}
	if exists == false {
		return false
	}

	value, err := redis.String(c.Do("HGET", tag, "value"))
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
	c := s.redisPool.Get()
	defer c.Close()
	_, err := redis.Bool(c.Do("DEL", tag))
	if err != nil {
		log.Println(err)
	}
}

// Delete Cache Value for Tag
func (s CacheService) DeleteTagsByPattern(pattern string) error {
	c := s.redisPool.Get()
	defer c.Close()
	iter := 0
	keys := []string{}
	for {
		arr, err := redis.Values(c.Do("SCAN", iter, "MATCH", pattern))
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
