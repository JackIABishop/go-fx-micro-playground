

# Go FX Micro-Playground 💱

A tiny Go project designed to explore microservices architecture, with a working example of currency conversion.

## 🏗️ Project Structure

```
go-fx-micro-playground/
├── services/
│   ├── rates/       # Provides FX rates
│   └── gateway/     # Handles /convert requests
├── internal/        # Shared code
├── go.mod
├── go.sum
└── README.md
```

## 🚀 Running the Services

In two terminals:

```
go run services/rates/main.go
go run services/gateway/main.go
```

## 🌐 Endpoints

### Health Checks
- `GET /health` → Returns `✅ <service> is up`

### Currency Conversion
- `GET /convert?from=GBP&to=EUR&amount=200`

Example response:
```json
{
  "from": "GBP",
  "to": "EUR",
  "amount": 200,
  "rate": 1.17,
  "converted": 234
}
```

## 💡 Future Plans

- Dockerize with `docker-compose`
- Add live FX rate updating
- Unit tests
- Request logging
- Swagger/OpenAPI docs

## 🧠 Why?

This is a hands-on playground to explore service separation, HTTP comms, Go idioms, and system design with a real-world use case.