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

### 🐳 With Docker Compose (Recommended)

Spin up both services with:

```bash
docker compose up --build
```

Then access:
- `http://localhost:8081/health` for the Rates service
- `http://localhost:8080/convert?from=USD&to=EUR&amount=100` for currency conversion

To stop and clean up:
```bash
docker compose down --volumes --remove-orphans
```

> **Note:** Use a `.env` file to set your `API_KEY` for authentication or set `DISABLE_AUTH=true` to disable auth for local development.

---

### 🧪 Dev Without Docker

In two terminals:

```bash
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

### Update Rates Endpoint
- `POST /rates` → Accepts a JSON body mapping base currencies to target currency rates, e.g.: 

  ```json
  {
    "USD": {"EUR": 0.95, "GBP": 0.82},
    "EUR": {"USD": 1.05}
  }
  ```

- Validates that each currency code is non-empty and all rates are positive.
- Returns:
  - `200 OK` with body `{"message":"rates updated"}` on success
  - `400 Bad Request` for malformed JSON or invalid data
  - `405 Method Not Allowed` for other HTTP methods

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

### Authentication

The `/convert` endpoint requires an `Authorization` header with a Bearer token matching the configured `API_KEY`:

```
Authorization: Bearer <API_KEY>
```

For local development, authentication can be disabled by setting the environment variable `DISABLE_AUTH=true`.


## 🧪 Testing

**CI Workflows**

- Unit tests are run automatically on every PR and push to `main` via GitHub Actions.
- Docker-based integration tests validate the containerised services and `/convert` endpoint.

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

## 🎉 Wrap-up
Thank you for exploring the Go FX Micro-Playground! This project demonstrates core microservice patterns in Go and provides a foundation for further experimentation.

## 🚀 Stretch Goals

If you’d like to fork this project and take it further, consider:

- Add live FX rate updates from an external API
- Generate Swagger/OpenAPI documentation