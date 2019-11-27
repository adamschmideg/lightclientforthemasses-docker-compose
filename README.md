## Installation

In a directory of your choice, type:

```
$ git clone http://github.com/gballet/lightclientforthemasses-docker-compose
$ cd lightclientforthemasses-docker-compose
$ docker-compose up
```

Wait for all the containers to start up and find one another, and then you can access Grafana by visiting `localhost:3000` with your browser.

Grafana admin credentials: admin/admin
Influxdb data source:
  * URL: http://ze_db:8086
  * DB name: metrics
  * User name: admin
  * Password: password

## Run in development

1. `docker-compose -f docker-compose.dev.yml up`
2. `go run faucet.go`
3. `go run usefaucet.go`