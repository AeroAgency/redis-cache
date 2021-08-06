package service

// Интерфейс для работы с redis cache
type RedisInterface interface {
	// Set Cache Value for Tag
	SetByTag(tag string, value interface{}, expire int)
	// Get Cache Value for Tag
	GetByTag(tag string, v interface{}) (result bool)
	// Delete Cache Value for Tag
	DeleteByTag(tag string)
}
