# Payment system

## Introduction
This application uses a classic MVC architecture along with some popular Golang libraries ([GORM](https://gorm.io/), [Ginkgo](https://onsi.github.io/ginkgo/#top), [Gomega](https://onsi.github.io/gomega/#top), [gorilla/mux](https://github.com/gorilla/mux)) and [PostgreSQL](https://www.postgresql.org/) as a DB for an implementation of a simple payment system.  

Its directory strcuture is created accordingly to the MVC pattern:
You can find: 
- **Models** in [models](internal/models) package
- **Views** in [views](internal/views) package
- **Controllers** in [controllers](internal/controllers) package

The application starts an HTTP server, which also serves static HTML that presents the merchants (and their transactions). 

You can access it by opening your browser and following thi next link for a local setup:
- http://localhost:8080/views/merchant

## Configuration

Most of the configurations for this application can be made using environment variables. 
Check out the Go structs defined in the [config](internal/config) package and pay attention to the `envconfig` tags to find out which are the environment variables that can be used.  

## Security (dummy)

The following API endpoints are secured by JWT tokens:
- **POST** /transaction (Create a transaction)
- **PUT** /merchant (Update a merchant)
- **DELETE** /merchant (Delete a merchant)

The security mechanism implemented is really simple. Your JWT token's `sub` claim needs to match the merchant's email. Since no passwords are stored in the DB and no login endpoints in exposed, you can craft the tokens yourself by using the [JWT debugger](https://jwt.io/) and **secretKey** as the signature's secret.

## CSV import

You can set the following environment variables to a .csv file path :
- **APP_MERCHANTS_IMPORT_PATH** for merchants
- **APP_ADMINS_IMPORT_PATH** for admins

Having set those, the application will load these users on startup by reading the specified .svc file.

For the file format, check out [assets/csv](assets/csv).

You can comment/uncomment lines 11 & 12 in [docker-compose.yml](docker-compose.yml) to control the behavior in a Dockerized environment.

## Running locally

### Docker compose

The simplest way to run this project is using [Docker Compose](https://docs.docker.com/compose/).

If you have docker compose installed, you can run the following commands in the project root to start the application in a container:
```bash
docker-compose build
docker-compose up
```

#### Tested using 
- docker-compose **v2.13.0** 
- Docker Desktop **4.15.0**

### Local startup
If you want to start this project locally (not in a container environment), you will need to make sure that you have a postgresql instance running with the proper schema migrations executed against it.

You can execute the migrations using [sqlmigrate](https://github.com/rubenv/sql-migrate) by running the following command in **project root**:

```bash 
migrate -path /sql -database  "postgres://YOUR_DB_USER:YOUR_DB_PASS@YOUR_DB_HOST:YOUR_DB_PORT/YOUR_DB_NAME?sslmode=disable" up
```

> **Note:** this works for Postgres instances which are configured to run in disabled SSL mode
> 

Then, once you've made sure that your environment is correctly configured (see **Configuration** above) run the application by running the following command in **project root**:
```bash 
go run cmd/main.go
```
#### Tested using
- go **1.19.4**
- migrate **4.14.0**
- postgres **v15.1**

## Example API requests

Since data modification is only possible by making HTTP Requests (not enough spare time to implement a React frontend :/ ) you can find a [Postman](https://www.postman.com/) collection with sample requests in [assets/http](assets/http).