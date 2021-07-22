# Оглавление

### Назначение:

Сервис для работы с кешированием в Redis

#### Примеры использования

* Установить кеш по ключу tag - Название ключа, например cache:productId_123:store_all cacheContent - кешируемое
  значение expire - срок жизни ключа (можно передать 0 для остутсвия срока жизни)
  `CacheService.SetByTag(tag, cacheContent, expire)`
* Получить кеш по ключу
  ` CacheService.GetByTag(tag, cacheContent)`
* Удалить кеш по ключу
  `CacheService.DeleteByTag(tag)`

## Параметры окружения

* RedisHost
* RedisPort

## Пример сборки контейнеров

* //Собирает контейнер для работы с сервисом кеширования
  `func buildCacheServiceContainer(ctn di.Container) (interface{}, error) { redisContainer := ctn.Get("redis").(redis.Conn)
  cacheService := redisCache.NewCacheService(zerologger.Logger, redisContainer)
  return cacheService, nil }
  `
* //Собирает контейнер с Redis
  `func buildRedisContainer(ctn di.Container) (interface{}, error) { c, err := redis.Dial("tcp", env.RedisHost+":"+env.RedisPort)
  if err != nil { log.Fatal(err)
  } return c, err }
  `
