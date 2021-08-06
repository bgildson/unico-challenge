# UNICO-CHALLENGE

[![Test and Send Coverage Report](https://github.com/bgildson/unico-challenge/actions/workflows/config.yml/badge.svg)](https://github.com/bgildson/unico-challenge/actions/workflows/config.yml)
[![Coverage Status](https://coveralls.io/repos/github/bgildson/unico-challenge/badge.svg)](https://coveralls.io/github/bgildson/unico-challenge)
[![Go Report Card](https://goreportcard.com/badge/github.com/bgildson/unico-challenge)](https://goreportcard.com/report/github.com/bgildson/unico-challenge)

This repository contains the solution to the [Unico Challenge](./challenge.pdf).

## Running the solution

_To follow the steps bellow, you must have installed [Docker](https://docs.docker.com/get-docker/) and [docker-compose](https://docs.docker.com/compose/install/)._

To run locally as production, execute the command bellow

```sh
docker-compose -f docker-compose-prod.yml up --build
```

Run the command bellow to apply the database migrations

```sh
docker run --rm -v $(pwd)/migrations:/migrations --network host migrate/migrate:v4.11.0 -path=/migrations -database "postgres://postgres:postgres@localhost:5432/unico_challenge?sslmode=disable" -verbose up
```

Run the command bellow to import the registers from the file [DEINFO_AB_FEIRASLIVRES_2014.csv](./DEINFO_AB_FEIRASLIVRES_2014.csv)

```sh
docker-compose -f docker-compose-prod.yml exec app /unico-challenge import -f /app/DEINFO_AB_FEIRASLIVRES_2014.csv
```

The file "[unico-challenge.postman_collection.json](./unico-challenge.postman_collection.json)" contains a **Postman Collection** to interact with the challenge solution.
