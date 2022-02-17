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

- This is made during [Design system course](http://wiki.cs.hse.ru/%D0%94%D0%B8%D0%B7%D0%B0%D0%B9%D0%BD_%D1%81%D0%B8%D1%81%D1%82%D0%B5%D0%BC_21/22)
taught by [HSE CS Faculty](https://cs.hse.ru/en/). 
[Youtube playlist](https://youtube.com/playlist?list=PLEwK9wdS5g0riA4Q_fqcjkv0zYf6HgRGJ), [Course Github link](https://github.com/hse-system-design-2021/public)