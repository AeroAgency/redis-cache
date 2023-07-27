package service

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/rs/zerolog"
)

// CacheService Сервис для кеширования в Redis
type CacheService struct {
	redisPool *redis.Pool
	logger    zerolog.Logger
}

// NewCacheService Конструктор
func NewCacheService(
	logger zerolog.Logger,
	redisPool *redis.Pool,
) *CacheService {
	return &CacheService{
		logger:    logger,
		redisPool: redisPool,
	}
}

// SetByTag Set Cache Value for Tag
func (s CacheService) SetByTag(tag string, value interface{}, expire int) error {
	c := s.redisPool.Get()
	defer c.Close()

	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	_, err = c.Do("HMSET", tag, "value", jsonValue)
	if err != nil {
		s.logger.Error().Err(err)
		return err
	}
	if expire != 0 {
		_, err = c.Do("EXPIRE", tag, expire)
		if err != nil {
			s.logger.Error().Err(err)
			return err
		}
	}
	return nil
}

// GetByTag Get Cache Value for Tag
func (s CacheService) GetByTag(tag string, v interface{}) (bool, error) {
	c := s.redisPool.Get()
	defer c.Close()

	exists, err := redis.Bool(c.Do("EXISTS", tag))
	if err != nil {
		s.logger.Error().Err(err)
		return false, err
	}
	if exists == false {
		return false, nil
	}

	value, err := redis.String(c.Do("HGET", tag, "value"))
	if err != nil {
		s.logger.Error().Err(err)
		return false, err
	}
	err = json.Unmarshal([]byte(value), v)
	if err != nil {
		s.logger.Error().Err(err)
		return false, err
	}
	return true, nil
}

// DeleteByTag Delete Cache Value for Tag
func (s CacheService) DeleteByTag(tag string) error {
	c := s.redisPool.Get()
	defer c.Close()
	_, err := redis.Bool(c.Do("DEL", tag))
	if err != nil {
		s.logger.Error().Err(err)
		return err
	}
	return nil
}

// DeleteTagsByPattern Delete Cache Value for Tag
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

// ClearCache Очистить весь кеш
func (s CacheService) ClearCache() error {
	err := s.DeleteTagsByPattern("cache:*")
	if err != nil {
		return err
	}
	return nil
}

// ClearCacheByTags Очистить кеш по тегу
func (s CacheService) ClearCacheByTags(tags []string) error {
	for _, tag := range tags {
		err := s.DeleteTagsByPattern("cache:" + tag)
		if err != nil {
			return err
		}
	}
	return nil
}
