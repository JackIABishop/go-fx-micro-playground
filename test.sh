#!/usr/bin/env bash
set -e

# 1️⃣ Start services
go run services/rates/main.go & 
PID_RATES=$!
go run services/gateway/main.go &
PID_GATEWAY=$!

# Give them a sec to initialize
sleep 1

# 2️⃣ Hit health checks
echo "🔍 Checking health endpoints..."
curl -f http://localhost:8081/health
curl -f http://localhost:8080/health

# 3️⃣ Test conversion endpoint
echo "🔍 Testing conversion..."
curl -f "http://localhost:8080/convert?from=USD&to=EUR&amount=42"

# 4️⃣ Teardown
kill $PID_GATEWAY $PID_RATES
echo "✅ Integration tests passed!"