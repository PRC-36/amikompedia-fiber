# amikompedia-fiber
Migration from Native

## Requirements/dependencies
- Docker
- Docker-compose

## Getting Started

- Run Postgres Images in Docker

```sh
make postgres
```

- Create Database on Postgres

```sh
make createdb
```

- Run Database Migrations

```sh
make migrateup
```


## API Request

| Endpoint                | HTTP Method |    Description    |
|-------------------------|:-----------:|:-----------------:|
| `api/v1/auth/_register` |   `POST`    | `Create accounts` |
| `api/v1/auth/_login`    |   `POST`    | `Login accounts`  |
| `api/v1/users`          |   `PATCH`   |   `User update`   |
| `api/v1/users/profile`  |    `GET`    |  `User Profile`   |
| `api/v1/surveys`        |   `POST`    | `Create Surveys`  |

