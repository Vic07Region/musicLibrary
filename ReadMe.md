#MusicLibrary service
в рамках тестовой задачи

описание задания лежит в `tz`

#Пример dotenv

`.env` - для продакшен
`.env.local` - для локальной разработки

```dotenv
#APP HOST PARAM
APP_HOST=:8080

#THIRD API SERVICE BASE URL
API_BASEURL=https://example.com/api

#database env param
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=dbname
#DB_HOST=localhost
#DB_PORT=5432
#DB_SSLMODE=disable
#DB_ROOTSERT=./path/to/root/cert

DB_DRIVER=postgres
MIGRATION_DIRS=./internal/database/migrations

## закомоентированные поля не обязательны к заполнению
```


Карта проекта:

* `cmd/main.go` точка входа

* `internal/connector/songinfo` запрос к другому api для получения подробной информации о песне

* `internal/database` слой бд для выполнения запросов к базе

* `internal/database/migrator` мигратор бд

* `internal/database/migration` миграции

* `internal/service` сервисный слой 

* `internal/app` слой gin, endpoint,mw и др

* `internal/lib/logger` логгер 

* `internal/pkg/app` инициализцаия
