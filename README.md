# web-text-manger

## Recommend

Use `docker-compose`

## Local Dependency

* Redis server
  need change `main.go` redis-server code 
* Mysql server
  need change `conf.json` DBHost

## Usage

`go run .`

If error occurred, run `go get github.com/kataras/iris/v12@master`

## Sample login data

* username: `admin`
* password: `password`

## Deploy

```bash

docker-compose up --build
```
