# Homework Microservice
A small go microservice that enables students to share their homeworks.
The frontend functionality is implemented as a part of [KAG-App](https://github.com/kagonlineteam/kag-app).

## Endpoints
|Endpoint|Description| Method | Permissions |
|-|-|-|-|
|/homework/v1/my| List homeworks for current user| GET | student |
|/homework/v1/homeworks | List (filtered) homeworks | GET | teacher,admin, `homework-show-all`|
|/homework/v1/homeworks| Create homework | POST | student,teacher,admin|
|/homework/v1/homework/id | Edit homework with given id | PUT | student,teacher,admin (reported only teacher,admin)|
|/homework/v1/homework/id | Edit homework with given id | PUT | `homework-allow-delete`|
|/homework/v1/report/id | Report homework with id | POST | student,teacher,admin|


## Environment variables
|Name|Description|
|-|-|
|HOMEWORK_JWT_PUB_KEY|Public key for jwt check. Newlines replaced by \n|
|HOMEWORK_PROXY_IP|Allowed proxy ip (typically docker host, or 127.0.0.1)|
|GIN_MODE|Should be "release" for production use|
|HOMEWORK_POSTGRES_DSN|DSN to Postgres database. Required when GIN_MODE=release|

## Authentication
Authentication is done via a JWT token that is passend as the `Authentication: Bearer` header.<br>
The token will be checked against the public key provided via `HOMEWORK_JWT_PUB_KEY`.<br>
The JWT needs the Claims:
- `roles` (included `ROLE_TEACHER` or  `ROLE_ADMINISTRATOR`)
- `stufe`
- `klasse`
- `sub` (needs to be `access_main`)

## Building
Manually building and pushing the docker container is possible by using
```bash
docker build -t ghcr.io/kagonlineteam/homework-microservice:latest .
docker push ghcr.io/kagonlineteam/homework-microservice:latest
```
