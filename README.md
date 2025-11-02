Test Connection Go and Database

$env:DB_USER = "root"; $env:DB_PASSWORD = "password"; $env:DB_HOST = "localhost"; $env:DB_PORT = "3306"; $env:DB_NAME = "icecream_truck"; go run cmd/dbtest/main.go
