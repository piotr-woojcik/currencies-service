# currencies-service
Interview assignment for Go developer role built with Gin framework.


## Docker build
To build the Docker image, you have to provide argument `exchange_app_id` with APP ID from https://openexchangerates.org/

```shell
docker build --build-arg exchange_app_id=<YOUR_APP_ID> -t currencies-service .
```