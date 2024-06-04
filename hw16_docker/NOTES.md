### Сборка приложения

1. Копировать Dockerfile в папку ../hw15_go_sql/
2. Запустить команду `docker buildx build --tag onlinestore:0.0.1 .`
3. Для размещения в репозитории Gitlab использовались команды:

```bash
docker pull postgres:15-alpine
docker tag postgres:15-alpine registry.gitlab.com/go-basic-2023-11/hw16-docker/postgres:15-alpine
docker push registry.gitlab.com/go-basic-2023-11/hw16-docker/postgres:15-alpine
docker buildx build --tag registry.gitlab.com/go-basic-2023-11/hw16-docker/onlinestore:0.0.1 --push .
```

### Запустить приложение

> из папки `hw16_docker`:

```bash
docker rm postgresql
docker volume rm otus_postgresql_data otus_postgresql_log -f
chmod go+rX -R .initdb
docker compose up
```

> если нет возможность получить образ из Gitlab Container Registry
> в директории `hw15_go_sql` docker compose выполняет сборку и запускается

```bash
cd hw15_go_sql
docker compose up postgresql onlinestore -d
```

Для проверки подключения приложения к БД, можно использовать другие настройки DSN, например:

```bash
go run . -dsn="postgres://manager:password@localhost:5432/onlinestore?sslmode=disable"
```

Команды для `curl` можно посмотреть в [заметках](/hw15_go_sql/NOTES.md)
