package service

import (
	"context"
	"github.com/rs/zerolog"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func TestCacheService(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	ctx := context.Background()

	cacheService := NewCacheService(client, zerolog.Nop(), ctx)

	t.Run("Set and Get by Tag", func(t *testing.T) {
		err := cacheService.SetByTag("test", "value", 10)
		assert.Nil(t, err)

		var value string
		exists, err := cacheService.GetByTag("test", &value)
		assert.Nil(t, err)
		assert.True(t, exists)
		assert.Equal(t, "value", value)
	})

	t.Run("Delete by Tag", func(t *testing.T) {
		err := cacheService.DeleteByTag("test")
		assert.Nil(t, err)

		var value string
		exists, err := cacheService.GetByTag("test", &value)
		assert.Nil(t, err)
		assert.False(t, exists)
	})

	t.Run("Delete Tags by Pattern", func(t *testing.T) {
		err := cacheService.SetByTag("test1", "value", 10)
		assert.Nil(t, err)
		err = cacheService.SetByTag("test2", "value", 10)
		assert.Nil(t, err)

		err = cacheService.DeleteTagsByPattern("test*")
		assert.Nil(t, err)

		var value string
		exists, err := cacheService.GetByTag("test1", &value)
		assert.Nil(t, err)
		assert.False(t, exists)
		exists, err = cacheService.GetByTag("test2", &value)
		assert.Nil(t, err)
		assert.False(t, exists)
	})

	t.Run("Clear Cache", func(t *testing.T) {
		err := cacheService.SetByTag("test", "value", 10)
		assert.Nil(t, err)

		err = cacheService.ClearCache()
		assert.Nil(t, err)

		var value string
		exists, err := cacheService.GetByTag("test", &value)
		assert.Nil(t, err)
		assert.False(t, exists)
	})

	t.Run("Clear Cache by Tags", func(t *testing.T) {
		err := cacheService.SetByTag("test1", "value", 10)
		assert.Nil(t, err)
		err = cacheService.SetByTag("test2", "value", 10)
		assert.Nil(t, err)

		err = cacheService.ClearCacheByTags([]string{"test1"})
		assert.Nil(t, err)

		var value string
		exists, err := cacheService.GetByTag("test1", &value)
		assert.Nil(t, err)
		assert.False(t, exists)
		exists, err = cacheService.GetByTag("test2", &value)
		assert.Nil(t, err)
		assert.True(t, exists)
	})
}
