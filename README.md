# SRATest - Loan Billing Engine

A simple **Loan Billing Engine** built in Go, supporting:

- Create a loan with weekly schedule
- Track **outstanding** amount
- Track **delinquent** borrowers (missed 2 consecutive payments)
- Make payments (single or multiple weeks at once)

---

## Prerequisites

- Go >= 1.23
- Git
- `curl` (or any HTTP client like Postman)

---

## How to Run

1. Clone repository:

```bash
  git clone https://github.com/kaes1412/SRATest.git
cd SRATest
```
2. Run the server
```bash
  go run main.go
```

## API Endpoints
### 1. POST `/loan`
**Response**
```json
{
  "code": 200,
  "message": "loan created successfully",
  "data": {
    "ID": "100",
    "Principal": 5000000,
    "InterestRate": 0.1,
    "TotalWeeks": 50,
    "WeeklyPayment": 110000,
    "Payments": [
      {"Week": 1, "Paid": false},
      {"Week": 2, "Paid": false},
      {"Week": 3, "Paid": false},
      ...
      {"Week": 50, "Paid": false}
    ]
  }
}
```
**Curl**
```bash
  curl -X POST http://localhost:8080/loan \
-H "Content-Type: application/json" \
-d '{"id":"100","principal":5000000}'
```

### 2. POST `/loan/{id}/pay`
**Request**
```json
{
  "code": 200,
  "message": "payment success",
  "data": {}
}
```
**Response**
```json
{
  "code": 200,
  "message": "payment success",
  "data": {}
}
```
**Curl**
```bash
  curl -X POST http://localhost:8080/loan/100/pay \
-H "Content-Type: application/json" \
-d '{"amount":110000}'
```

### 3. GET `/loan/{id}/outstanding`
**Response**
```json
{
  "code": 200,
  "message": "outstanding fetched successfully",
  "data": {
    "outstanding": 5500000
  }
}
```
**Curl**
```bash
  curl http://localhost:8080/loan/100/outstanding
```

### 4. GET `/loan/{id}/delinquent`
**Response**
```json
{
  "code": 200,
  "message": "Delinquent status fetched successfully",
  "data": {
    "delinquent": true
  }
}
```
**Curl**
```bash
  curl http://localhost:8080/loan/100/delinquent
```


## Notes
- All responses are JSON with consistent structure:
```json
{
  "code": 200,
  "message": "string",
  "data": {}
}
```