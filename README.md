

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

### Rates Endpoint
- `GET /rates` → Returns a JSON object mapping base currencies to target currency rates, for example:
  ```json
  {
    "USD": {"EUR": 0.92, "GBP": 0.78, "JPY": 135.33},
    "EUR": {"USD": 1.09, "GBP": 0.85},
    "GBP": {"USD": 1.29, "EUR": 1.17}
  }
  ```

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

- Dockerize services with `docker-compose`  
- Add live FX rate updates from an external API  
- Set up CI/CD pipeline (GitHub Actions) for automated testing  
- Generate Swagger/OpenAPI documentation  

## 🧪 Testing

**Unit tests**  
Run all unit tests for both services with:
```bash
go test ./...
```

Or individually:
```bash
go test services/rates
go test services/gateway
```

**Integration tests**  
After making `test.sh` executable (`chmod +x test.sh`), run:
```bash
./test.sh
```
This script starts both services, probes health and conversion endpoints, then tears them down.

## 🧠 Why?

This is a hands-on playground to explore service separation, HTTP comms, Go idioms, and system design with a real-world use case.