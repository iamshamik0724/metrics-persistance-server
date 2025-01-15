# Metrics Persistence Server

A Go-based service for collecting, processing, and storing API performance metrics via UDP, with real-time WebSocket streaming and REST APIs for querying metrics.

---

## Setup and Installation

1. **Clone the Repository**  
   - Clone the repository and install dependencies:  
     ```bash
     git clone https://github.com/your-username/metrics-persistence-server.git
     cd metrics-persistence-server
     go mod tidy
     ```
2. **Install PostgreSQL and TimescaleDB**  
   - This service uses TimescaleDB for persisting the metrics data received from UDP server
   - Install postgres timescaleDB (refer: https://docs.timescale.com/self-hosted/latest/install/)
   - Execute the sql files under the folder internal/database in this repo to ceate the required database setup.

3. **Update Configuration**  
   - Edit the `config.toml` file in the root directory and update your setup details if needed

4. **Run the Service**  
   - Start the service:  
     ```bash
     go run main.go
     ```
---

## API Specification

### WebSocket API
- **Endpoint**: `/ws`  
  Opens a WebSocket connection to stream real-time metrics.

### REST API
- **GET `/metrics`**  
  Fetches the last 10 mins metrics stored in the database.
  (Note the timestamps and responses are provided in seperate arrays to ease the process of plot in frontend)

  **Response Format**:
  ```json
  {
      "metrics": [
          {
              "route": "/health",
              "method": "GET",
              "timestamps": ["2025-01-15T12:00:00Z", "2025-01-15T12:01:00Z"],
              "responses": [
                 {
                     "time": 150.2,
                     "status": 200
                 },
                 {
                     "time": 120.5,
                     "status": 400
                 }
             ]
          },
          {
              "route": "/api/products",
              "method": "GET",
              "timestamps": ["2025-01-15T12:00:00Z", "2025-01-15T12:01:00Z"],
              "responses": [
                 {
                     "time": 10.2,
                     "status": 200
                 },
                 {
                     "time": 120.5,
                     "status": 200
                 }
             ]
          },
      ]
  }
