package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type healthResponse struct {
	Status    string `json:"status"`
	Service   string `json:"service"`
	Timestamp string `json:"timestamp"`
}

type infoResponse struct {
	Service     string `json:"service"`
	Version     string `json:"version"`
	Environment string `json:"environment"`
	Timestamp   string `json:"timestamp"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	version := os.Getenv("APP_VERSION")
	if version == "" {
		version = "dev"
	}
	environment := os.Getenv("APP_ENV")
	if environment == "" {
		environment = "local"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(responseWriter http.ResponseWriter, _ *http.Request) {
		responseWriter.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = fmt.Fprintf(responseWriter, `<!doctype html>
<html lang="ru">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>example-service</title>
  <style>
    :root {
      color-scheme: dark;
      font-family: "Segoe UI", sans-serif;
    }
    body {
      margin: 0;
      min-height: 100vh;
      display: grid;
      place-items: center;
      background:
        radial-gradient(circle at top left, rgba(255, 166, 0, 0.25), transparent 30%%),
        linear-gradient(135deg, #20120a, #3d1f0d 45%%, #121212);
      color: #f8efe9;
    }
    main {
      width: min(720px, calc(100%% - 32px));
      padding: 32px;
      border-radius: 24px;
      background: rgba(22, 16, 14, 0.82);
      border: 1px solid rgba(255, 255, 255, 0.08);
      box-shadow: 0 24px 60px rgba(0, 0, 0, 0.35);
    }
    h1 {
      margin: 0 0 12px;
      font-size: clamp(2rem, 6vw, 3.5rem);
      line-height: 1;
    }
    p {
      margin: 0 0 24px;
      font-size: 1.05rem;
      color: #e8d8cf;
    }
    .grid {
      display: grid;
      gap: 12px;
      grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
    }
    .card {
      padding: 16px;
      border-radius: 18px;
      background: rgba(255, 255, 255, 0.05);
    }
    .label {
      font-size: 0.85rem;
      text-transform: uppercase;
      letter-spacing: 0.08em;
      color: #f2af88;
    }
    .value {
      margin-top: 8px;
      font-size: 1.1rem;
      font-weight: 700;
    }
    code {
      color: #ffd0a5;
    }
  </style>
</head>
<body>
  <main>
    <h1>example-service работает</h1>
    <p>Это тестовое приложение для проверки деплоя через deploy-service. Для машинной проверки доступны <code>/health</code>, <code>/ready</code> и <code>/api/info</code>.</p>
    <section class="grid">
      <article class="card">
        <div class="label">Сервис</div>
        <div class="value">example-service</div>
      </article>
      <article class="card">
        <div class="label">Версия</div>
        <div class="value">%s</div>
      </article>
      <article class="card">
        <div class="label">Окружение</div>
        <div class="value">%s</div>
      </article>
      <article class="card">
        <div class="label">Порт</div>
        <div class="value">%s</div>
      </article>
    </section>
  </main>
</body>
</html>`, version, environment, port)
	})
	mux.HandleFunc("/health", func(responseWriter http.ResponseWriter, _ *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(responseWriter).Encode(healthResponse{
			Status:    "ok",
			Service:   "example-service",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
	})
	mux.HandleFunc("/ready", func(responseWriter http.ResponseWriter, _ *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(responseWriter).Encode(map[string]string{
			"status": "ready",
		})
	})
	mux.HandleFunc("/api/info", func(responseWriter http.ResponseWriter, _ *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(responseWriter).Encode(infoResponse{
			Service:     "example-service",
			Version:     version,
			Environment: environment,
			Timestamp:   time.Now().UTC().Format(time.RFC3339),
		})
	})

	address := ":" + port
	log.Printf("example-service started on %s", address)
	if err := http.ListenAndServe(address, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
