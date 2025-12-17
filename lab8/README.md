# Todo API (Go + MongoDB + Kubernetes)

Простое REST API для списка дел на Go с хранением в MongoDB. В репозитории есть Dockerfile, манифесты Kubernetes и пайплайн GitHub Actions для сборки и деплоя.

## API
- `GET /todos` — список задач
- `POST /todos` — создать задачу, тело: `{"title":"write docs","notes":"optional"}`
- `GET /todos/{id}` — получить по id
- `PUT /todos/{id}` — частичное обновление (`title`, `notes`, `completed`)
- `DELETE /todos/{id}` — удалить задачу
- `GET /healthz` — проверка живости

Формат задачи:
```json
{
  "id": "ObjectID hex",
  "title": "string",
  "notes": "string",
  "completed": false,
  "createdAt": "RFC3339 timestamp",
  "updatedAt": "RFC3339 timestamp"
}
```

## Конфигурация
Переменные окружения:
- `MONGO_URI` (по умолчанию `mongodb://localhost:27017`)
- `MONGO_DB` (по умолчанию `todos`)
- `PORT` (по умолчанию `8080`)

## Локальный запуск
```bash
go mod tidy      # скачает go.mongodb.org/mongo-driver, нужен интернет
go run ./cmd/server
```

## Docker
Сборка и запуск:
```bash
docker build -t todo-api:latest .
docker run --rm -p 8080:8080 -e MONGO_URI=mongodb://host.docker.internal:27017 todo-api:latest
```

## Kubernetes (пример с Minikube)
```bash
kubectl apply -f k8s/mongo.yaml
kubectl apply -f k8s/app.yaml
kubectl rollout status deployment/todo-api
minikube service todo-api --url   # либо NodePort 30080
```
В `k8s/app.yaml` указан образ `ghcr.io/example/todo-api:latest`; замените на свой тег из реестра/CI.

## GitHub Actions
`.github/workflows/ci.yml`:
- `go test`
- сборка и пуш Docker-образа в GHCR (`ghcr.io/<owner>/todo-api`)
- деплой манифестов при пуше в `main` с использованием секрета `KUBE_CONFIG` (base64 kubeconfig)

## Проверка через curl
```bash
curl -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"test task","notes":"demo"}'
curl http://localhost:8080/todos
```
