# example-service

Простой тестовый HTTP-сервис для проверки полного взаимодействия с `deploy-service`:

1. загрузка репозитория в GitHub,
2. создание проекта в `deploy-service`,
3. запуск GitHub bootstrap flow,
4. merge Pull Request,
5. автосборка и деплой в Kubernetes.

## Что внутри
- `cmd/server/main.go` — HTTP сервис с понятными маршрутами для браузера и probe-проверок
- `Dockerfile` — контейнеризация сервиса

## HTTP маршруты
- `/` — человекочитаемая HTML-страница, чтобы быстро проверить деплой в браузере
- `/health` — JSON healthcheck
- `/ready` — JSON readiness endpoint
- `/api/info` — JSON с версией, окружением и timestamp

## Локальный запуск

```bash
go run ./cmd/server
curl http://localhost:8080/
curl http://localhost:8080/health
curl http://localhost:8080/ready
curl http://localhost:8080/api/info
```

## Подготовка GitHub репозитория

1. Создать пустой репозиторий, например `example-service`.
2. В каталоге `example-service` выполнить:

```bash
git init
git add .
git commit -m "initial example service"
git branch -M main
git remote add origin git@github.com:<OWNER>/example-service.git
git push -u origin main
```

## Проверка интеграции с deploy-service

### 1) Запустить deploy-service

```bash
cd ../курсач_2/deploy-service
KUBERNETES_PROVISIONER=kubectl \
GITHUB_AUTOMATION_MODE=real \
GITHUB_TOKEN=<GITHUB_TOKEN> \
HTTP_ADDRESS=:8080 \
go run ./cmd/server
```

### 2) Создать проект

```bash
curl -X POST http://localhost:8080/projects \
  -H "Content-Type: application/json" \
  -d '{"name":"example-service","ownerId":"demo-user"}'
```

Сохранить `id` проекта в переменную:

```bash
PROJECT_ID=<project-id>
```

### 3) Запустить bootstrap flow для GitHub

```bash
curl -X POST "http://localhost:8080/projects/$PROJECT_ID/github/bootstrap" \
  -H "Content-Type: application/json" \
  -d '{
    "repositoryOwner":"<OWNER>",
    "repositoryName":"example-service",
    "baseBranch":"main",
    "serviceName":"example-service",
    "dockerfilePath":"Dockerfile",
    "servicePort":80,
    "containerPort":8080
  }'
```

В ответе будет `pullRequestUrl`.

### 4) Настроить секрет для GitHub Actions

Получить kubeconfig в base64:

```bash
kubectl config view --raw | base64
```

Добавить в GitHub repository secret:
- `KUBECONFIG_BASE64` = результат команды выше (одной строкой).

### 5) Merge Pull Request

После merge в репозитории появятся:
- `.github/workflows/deploy-service.yml`
- `k8s/example-service/deployment.yaml`
- `k8s/example-service/service.yaml`

И автоматически запустится workflow на ветке `main`.

## Что проверять после запуска workflow

```bash
kubectl get ns | grep "project-$PROJECT_ID"
kubectl get pods -n "project-$PROJECT_ID"
kubectl get svc -n "project-$PROJECT_ID"
```

Если pod в статусе `Running`, то ожидаемое взаимодействие подтверждено.

Если у сервиса есть внешний адрес или `port-forward`, то можно дополнительно проверить:

```bash
curl http://<SERVICE_URL>/
curl http://<SERVICE_URL>/health
curl http://<SERVICE_URL>/api/info
```
