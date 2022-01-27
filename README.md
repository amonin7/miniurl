# miniurl

## Short description

This is a service with which the user can get a shortened url. The service supports two functions:
- shorten the url provided by the user;
- redirect the user from the shortened url to the real one.

## Installation guide

1. Start [Docker](https://www.docker.com/) on your machine
2. Run this command in terminal
```shell
docker-compose up --build
```

## Prerequisites

- [Golang (version 1.17)](https://go.dev/)
- [Docker (version 3)](https://www.docker.com/)
- [MongoDB (version 4.4)](https://www.mongodb.com/) - as main storage
- [RedisDB (version 6.2.6)](https://redis.io/) - as a cache storage
- [Driver for Redis](https://github.com/go-redis/redis) - to connect GoLang and Redis
