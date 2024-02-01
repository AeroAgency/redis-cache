package service

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"time"
)

// CacheService Сервис для кеширования в Redis
type CacheService struct {
	client *redis.Client
	logger zerolog.Logger
	ctx    context.Context
}

// NewCacheService Конструктор
func NewCacheService(
	client *redis.Client,
	logger zerolog.Logger,
	ctx context.Context,
) *CacheService {
	return &CacheService{
		logger: logger,
		client: client,
		ctx:    ctx,
	}
}

// SetByTag Set Cache Value for Tag
func (s *CacheService) SetByTag(tag string, value interface{}, expire int) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		s.logger.Error().Err(err)
		return err
	}
	st := s.client.HMSet(s.ctx, tag, "cache", jsonValue)
	if st.Err() != nil {
		s.logger.Error().Err(st.Err())
		return err
	}
	if expire != 0 {
		exSt := s.client.Expire(s.ctx, tag, time.Duration(expire)*time.Second)
		if exSt.Err() != nil {
			s.logger.Error().Err(exSt.Err())
			return err
		}
	}
	return st.Err()

}

// GetByTag Get Cache Value for Tag
func (s *CacheService) GetByTag(tag string, v interface{}) (bool, error) {
	exSt := s.client.Exists(s.ctx, tag)
	r, err := exSt.Result()
	if err != nil {
		s.logger.Error().Err(err)
		return false, err
	}
	if r == 0 {
		return false, nil
	}
	hGet := s.client.HGet(s.ctx, tag, "cache")
	if err != nil {
		s.logger.Error().Err(err)
		return false, err
	}
	err = json.Unmarshal([]byte(hGet.Val()), v)
	if err != nil {
		s.logger.Error().Err(err)
		return false, err
	}
	return true, nil
}

// DeleteByTag Delete Cache Value for Tag
func (s *CacheService) DeleteByTag(tag string) error {
	st := s.client.Del(s.ctx, tag)
	return st.Err()
}

// DeleteTagsByPattern Delete Cache Value for Tag
func (s *CacheService) DeleteTagsByPattern(pattern string) error {
	var cursor uint64
	for {
		var keys []string
		var err error
		keys, cursor, err = s.client.Scan(s.ctx, cursor, pattern, 0).Result()
		if err != nil {
			s.logger.Error().Err(err)
			return err
		}
		for _, key := range keys {
			err = s.DeleteByTag(key)
			if err != nil {
				s.logger.Error().Err(err)
				return err
			}
		}
		if cursor == 0 { // no more keys
			break
		}
	}
	return nil
}

// ClearCache Очистить весь кеш
func (s *CacheService) ClearCache() error {
	err := s.DeleteTagsByPattern("*")
	if err != nil {
		s.logger.Error().Err(err)
		return err
	}
	return nil
}

// ClearCacheByTags Очистить кеш по тегу
func (s *CacheService) ClearCacheByTags(tags []string) error {
	for _, tag := range tags {
		err := s.DeleteTagsByPattern(tag + "*")
		if err != nil {
			s.logger.Error().Err(err)
			return err
		}
	}
	return nil
}
