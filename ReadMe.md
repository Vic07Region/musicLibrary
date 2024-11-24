# MusicLibrary service
в рамках тестовой задачи

описание задания лежит в `tz` [readme_tz](tz/readme.md) 

# Пример dotenv

* `.env` - для продакшен
* `.env.local` - для локальной разработки

```dotenv
#APP HOST PARAM
APP_HOST=:8080

DEBUG=TRUE
PRODUCTION=TRUE

#THIRD API SERVICE BASE URL
API_BASEURL=http://example.com/api


#database env param
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=dbname
#DB_HOST=localhost
#DB_PORT=5432
#DB_SSLMODE=disable
#DB_ROOTSERT=./path/to/root/cert

DB_MAX_CONN=80
DB_MAX_IDLE=5
DB_MAX_LIFETIME=10
#minute

DB_DRIVER=postgres
MIGRATION_DIRS=./migrations

## закомоентированные поля не обязательны к заполнению
```
присутствует [makefile](makefile)

## Доступные команды:
	* install-deps   Устанавливает зависимости проекта"
	* get-deps       Загружает зависимости проекта"
	* build          Собирает проект"
	* run           Запускает проект"
	* clean          Очищает сгенерированные файлы"
	* swag-docs      Генерирует документацию swagger"

# Endpoints
* `/api/v1` - root api
* `/api/v1/songs` *GET* список песен
* `/api/v1/songs/{id}` *GET* получение текста песни
* `/api/v1/songs/{id}` *DELETE* удаление песни
* `/api/v1/songs/{id}` *PATCH* изменение песни
* `/api/v1/songs/{id}/verse` *PATCH* изменение куплета песни
* `/api/v1/songs/new` *POST* создание песни
* `/info` *GET* демо ручка для тестирования NewSong

# Swagger info
[swagger_UI](http://localhost:8080/swagger/index.html) 
[swagger_json](http://localhost:8080/swagger/doc.json) 

# Карта проекта:

* `cmd/main.go` точка входа

* `internal/connector/songinfo` запрос к другому api для получения подробной информации о песне

* `internal/database` слой бд для выполнения запросов к базе

* `internal/database/migrator` мигратор бд

* `internal/database/migration` миграции

* `internal/service` сервисный слой 

* `internal/app` слой gin, endpoint,mw и др

* `internal/lib/logger` логгер 

* `internal/pkg/app` инициализцаия

* `docs` сгенерированные сваггером документы
