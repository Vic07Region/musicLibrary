#MusicLibrary service
в рамках тестовой задачи

#Пример dotenv

`.env` - для продакшен
`.env.local` - для локальной разработки

```dotenv
#APP HOST PARAM
APP_HOST=:8080

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

DB_DRIVER=postgres
MIGRATION_DIRS=./internal/database/migrations
```

