## Installation

In a directory of your choice, type:

```
$ git clone http://github.com/adamschmideg/lightclientforthemasses-docker-compose
$ cd lightclientforthemasses-docker-compose
$ git checkout faucet
$ docker-compose up
```

Wait for all the containers to start up and find one another, and then you
can access Grafana by visiting `localhost:3000` with your browser. Grafana
admin credentials: admin/admin. You'll be requested to change it on your
first login. It takes a few minutes until Geth manages to send data to
InfluxDB, so monitoring data might be not immediately visible.

## Use it

- Visit http://localhost:8088
- Enter a nodeID and submit it

