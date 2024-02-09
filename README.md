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

| Endpoint                                 | HTTP Method |      Description       |
|------------------------------------------|:-----------:|:----------------------:|
| `api/v1/users`                           |   `POST`    |   `Create accounts`    |
| `api/v1/user/_login`                     |   `POST`    |    `Login accounts`    |
