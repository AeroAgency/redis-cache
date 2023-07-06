package service

// RedisInterface Интерфейс для работы с redis cache
type RedisInterface interface {
	// SetByTag Set Cache Value for Tag
	SetByTag(tag string, value interface{}, expire int) error
	// GetByTag Get Cache Value for Tag
	GetByTag(tag string, v interface{}) (bool, error)
	// DeleteByTag Delete Cache Value for Tag
	DeleteByTag(tag string) error
	// DeleteTagsByPattern Delete Cache Value for Tag
	DeleteTagsByPattern(pattern string) error
	// ClearCache Очистить весь кеш
	ClearCache() error
	// ClearCacheByTags Очистить кеш по тегу
	ClearCacheByTags(tags []string) error
}
