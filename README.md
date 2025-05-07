# Pack Size Calculator

This application calculates the optimal pack sizes for shipping items based on order quantities. It follows these rules:
1. Only whole packs can be sent (packs cannot be broken open)
2. Within the constraints of Rule 1, send out the least amount of items to fulfil the order
3. Within the constraints of Rules 1 & 2, send out as few packs as possible to fulfil each order

## Architecture

The application consists of:
- Go backend (REST API)
- React frontend
- Docker containers for both services

## Prerequisites

- Docker and Docker Compose
- Go 1.21 or later (for local development)
- Node.js 18 or later (for local development)

## Running the Application

### Using Docker (Recommended)

1. Clone the repository
2. Run the following command in the root directory:
   ```bash
   docker-compose up --build
   ```
3. Access the application at http://localhost

### Local Development

#### Backend
1. Navigate to the root directory
2. Run:
   ```bash
   go mod download
   go run main.go
   ```

#### Frontend
1. Navigate to the frontend directory
2. Run:
   ```bash
   npm install
   npm run dev
   ```
3. Access the frontend at http://localhost:3000

## API Endpoints

### Pack Size Management
- `GET /pack-sizes`
  - Response: `{ "sizes": number[] }`
  - Returns the current list of available pack sizes

- `PUT /pack-sizes`
  - Request body: `{ "sizes": number[] }`
  - Response: `{ "sizes": number[] }`
  - Updates the entire list of pack sizes

- `POST /pack-sizes`
  - Request body: `{ "size": number }`
  - Response: `{ "sizes": number[] }`
  - Adds a new pack size to the list

- `DELETE /pack-sizes/:size`
  - Response: `{ "sizes": number[] }`
  - Removes a pack size from the list

### Pack Calculation
- `POST /calculate`
  - Request body: `{ "orderQuantity": number }`
  - Response: `{ "orderQuantity": number, "packs": { "size": quantity }, "totalItems": number }`
  - Calculates the optimal pack distribution for the given order quantity using the configured pack sizes

## Testing

### Backend Tests
Run the following command in the root directory:
```bash
go test
```

### Frontend Tests
Run the following command in the frontend directory:
```bash
npm test
``` 