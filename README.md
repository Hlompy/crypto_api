[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/golang-migrate/migrate/ci.yaml?branch=master)](https://github.com/golang-migrate/migrate/actions/workflows/ci.yaml?query=branch%3Amaster)
[![GoDoc](https://pkg.go.dev/badge/github.com/golang-migrate/migrate)](https://pkg.go.dev/github.com/golang-migrate/migrate/v4)
[![Coverage Status](https://img.shields.io/coveralls/github/golang-migrate/migrate/master.svg)](https://coveralls.io/github/golang-migrate/migrate?branch=master)
[![packagecloud.io](https://img.shields.io/badge/deb-packagecloud.io-844fec.svg)](https://packagecloud.io/golang-migrate/migrate?filter=debs)
[![Docker Pulls](https://img.shields.io/docker/pulls/migrate/migrate.svg)](https://hub.docker.com/r/migrate/migrate/)
![Supported Go Versions](https://img.shields.io/badge/Go-1.19%2C%201.20-lightgrey.svg)
[![GitHub Release](https://img.shields.io/github/release/golang-migrate/migrate.svg)](https://github.com/golang-migrate/migrate/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/golang-migrate/migrate/v4)](https://goreportcard.com/report/github.com/golang-migrate/migrate/v4)

# Crypto-Api

**Crypto-Api** is an ultimate API and database service designed to handle cryptocurrency data. It supports all CRUD (Create, Read, Update, Delete) operations using **PostgreSQL** for data storage and **Docker** for continuous integration and deployment (CI/CD). This project is built with **Golang**, utilizing the **Gin Framework** for building the API endpoints.

---

## Technologies Used
- **Golang**: The backend language used to build the API.
- **PostgreSQL**: The database used for storing cryptocurrency data.
- **Docker**: Containerization and deployment of the application.
- **Gin Framework**: A high-performance web framework for Golang.

---

## Features
- **Full CRUD operations** for cryptocurrency data.
- **API endpoints** for retrieving, creating, updating, and deleting coins.
- **PostgreSQL** database integration for persistent data storage.
- **Docker** support for easy deployment and CI/CD.

---

## Project Setup

### 1. Clone the Repository

First, clone the repository to your local machine:

```bash
git clone https://github.com/Hlompy/crypto_api.git
cd crypto_api
```

### 2. Initialize Go Modules

Second, initialize go modules

```bash
go mod init crypto_api
go mod tidy
```

### 3. Setup Docker and PostgreSQL

The project uses Docker and Docker Compose to manage the PostgreSQL database and application container. 
To start everything, use the following commands:

```bash
docker-compose up --build
```

## API Endpoints

### 1. GET /coins

Description: Fetch a list of all stored cryptocurrencies.

Response: Returns a JSON array with all the top cryptocurrencies.

```json
[
    {
        "name": "Bitcoin",
        "symbol": "btc",
        "usd_price": 78862
    },
    {
        "name": "Ethereum",
        "symbol": "eth",
        "usd_price": 1556.72
    },
    {
        "name": "Tether",
        "symbol": "usdt",
        "usd_price": 0.999463
    },
]

```

### 2. GET /coins/users

Description: Fetch a list of all cryptocurrencies created by USER from DATABASE.

Response: Returns a JSON array with all the cryptocurrencies of user.

```json
[
    {
        "id": 1,
        "name": "Hlompushka",
        "symbol": "hlmp",
        "usd_price": 0.1489,
        "is_user_coin": true
    },
    {
        "id": 2,
        "name": "Vatrushka",
        "symbol": "vatr",
        "usd_price": 0.7777,
        "is_user_coin": true
    },
    {
        "id": 3,
        "name": "NASHA-PUSHKA",
        "symbol": "NPU",
        "usd_price": 0.1337,
        "is_user_coin": true
    }
]
```

### 3. GET /coins/symbol/{symbol}

Description: Retrieve a specific cryptocurrency by its SYMBOL from API - CoinGecko

Parameters: symbol (required) - The symbol of the cryptocurrency from top cryptocurrencies
from all world.

Response: Returns the data for a specific coin.

```json
[
  {
    "name": "Bitcoin",
    "symbol": "btc",
    "usd_price": 50000
  },
  ...
]
```


### 4. GET /coins/{id}

Description: Retrieve a specific cryptocurrency by its ID from DATABASE.

Parameters: id (required) - The ID of the cryptocurrency created by user.

Response: Returns the data for a specific coin.

```json
{
    "id": 2,
    "name": "Vatrushka",
    "symbol": "vatr",
    "usd_price": 0.7777
}
```

### 5. POST /coins

Description: Create a new cryptocurrency record in DATABASE.

Body: The request body should contain the details of the new coin, including its name, symbol, and USD price.

```json
{
    "name": "Monetka",
    "symbol": "MNT",
    "usd_price": 0.12345
}
```

### 6. PATCH /coins/{id}

Description: Update an existing cryptocurrency record in DATABASE by its ID.

Parameters: id (required) - The ID of the cryptocurrency.

Body: The fields to be updated, such as name, symbol, or USD price.

```json
{
  "usd_price": 55000
}
```

### 7. DELETE /coins/{id}

Description: Delete a cryptocurrency record by its ID from DATABASE.

Parameters: id (required) - The ID of the cryptocurrency to delete.

Response: Returns a success message.

```json
{
  "message": "Coin deleted successfully"
}
```

## Docker Configuration

Dockerfile: The project includes a Dockerfile to build the application container.

docker-compose.yml: The docker-compose.yml file is used to define and run multi-container Docker applications. It configures the API container and a PostgreSQL container.

### Environment Variables

The project may require setting environment variables for database credentials. By default, Docker Compose handles the environment variables defined in the docker-compose.yml file, but you can override them in a ``` .env file.```
