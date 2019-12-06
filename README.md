# Light faucet

This is a private faucet for an Ethereum light server. It needs a little
explanation what it does. 

## The problem
Syncing with the Ethereum mainnet takes days. If you're using fast sync, it
can go down to [a little more than one day][full node sync]. Light clients
are developed to address this issue. A light client can sync really fast, but
it needs a node willing to provide it with data. If you are running a full
node, you're in a symmetrical situation: you rely on other nodes to give you
data; at the same time you give data to other nodes that need it. Running a
light server is a different case. Light clients will request data from you,
but you will get nothing in return. In short, there is no economic incentive
to run a light server. Those who run them do so for the benefit of the
ecosystem. No wonder, there are very few of them. Which leads to a sad
situation: if you spin up a light client, chances are you will not find a
peer or the one you find will keep dropping the connection.

## Light protocol v4 
The idea to solve this problem is light servers charge light clients. When
you run a light client, you have some service tokens you can spend on
requesting data from light servers. Light servers have a strategy to serve
requests. They may choose to prioritize requests by the amount of tokens
offered, just like a miner may choose to rank transactions by the gas price
they are willing to pay. A light server may even offer some of its capacity
for free.

## The light faucet
Don't run to the nearest exchange yet, this token is still in the making.
What is on offer, though, is a limited variant. Here's how it works.

[full node sync]: https://medium.com/@mswezey/2019ethereumfullnode-ba6e05ebf363

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

