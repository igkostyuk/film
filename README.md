# film cache api

### Run
```console
# Running from within docker compose
$ docker-compose up
```
### Load data
[load postgresql sample data](https://www.postgresqltutorial.com/postgresql-getting-started/load-postgresql-sample-database/)

### Test
#### Available Endpoints
```
GET /film/:title
GET /film/metric
```

[vegeta](https://github.com/tsenart/vegeta)
```console
# Running vegeta with target list
$ vegeta attack -duration=5s -rate=5 -targets=targets.list | vegeta report
```
