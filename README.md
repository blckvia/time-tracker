# Time Tracker

## Getting start



### 1. Run apps

#### Run backend app

```
$> task docker
```

```
$> task compose-up
```

**Enable migrations**

```
$> task migrate-up
```


### Env

### DO NOT FORGET CREATE .env file

| Name                               | Description                                                             | Default                                                             |
|------------------------------------|-------------------------------------------------------------------------|---------------------------------------------------------------------|
| PORT                               | app port                                                                | 8080                                                                |
| DB_USERNAME                        | fill test data                                                          | postgres                                                            |
| DB_PASSWORD                        | fill test data                                                          | postgres                                                            |
| DB_NAME                            | fill test data                                                          | postgres                                                            |
| DB_HOST                            | fill test data                                                          | postgres                                                            |  
| DB_SSLMODE                         | ssl mode                                                                | disable                                                             |
| DB_URL                             | fill test data url                                                      | postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable |
| DB_PORT                            | bd port                                                                 | 5432                                                                |





### API Docs


http://127.0.10.5:8080/swagger/index.html


#### Swagger specs:
- [api group](./docs/swagger.yaml)

