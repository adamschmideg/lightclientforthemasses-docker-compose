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
- You learn about a light server
- Visit the faucet connected to it
- Get some tokens from the faucet
- Use the light server

There is no "official" list of light servers. (TBD the unofficial ones).
The plan is to run a light server with a faucet whose parameters will be
included here as soon as it gets deployed.

Caveat: the faucet gives tokens for a single light server only. It cannot be
used at another light server and it can't be transferred to another client
than it was originally given to (the token is actually associated with a
client nodeID).

Getting tokens does not mean your client will automically connect to it. You
have to `admin.addPeer()` its enode manually. Even then, you will not
necessarily get data right away, it depends on the capacity of the server and
the number of clients it is serving.

## The private faucet
Running a faucet is a risky business, because it is a kind of free money
dispenser some may want to abuse. The 
[Rinkeby testnet faucet][rinkeby faucet] uses a number of tools to
counter-measure attacks. That includes captcha, proof of authority, and
limiting requests. It's cumbersome to implement and maintain all that
functionality (or steal from [the code][faucet code]). The approach I am
taking now is to run the faucet as a web service in a way that only I can
access it (via ssh tunneling). You request your token by writing your nodeID
to a twitter thread and I will manually add it. The light server itself is
publicly available, any light client can connect to it, but those that were
manually added will have a higher priority.

[full node sync]: https://medium.com/@mswezey/2019ethereumfullnode-ba6e05ebf363
[rinkeby faucet]: https://faucet.rinkeby.io/
[faucet code]: https://github.com/ethereum/go-ethereum/tree/master/cmd/faucet

## Running your own light server and faucet
If you want to see how it works, including monitoring what is going on on the
server, just clone this repository and run `docker-compose up`. It will run a
number of services, such as
- A light server on the standard 8545 port, fully dedicated to light clients
- A Grafana instance on the standard 3000 port with dashboards already
  configured to show the internals of the server. The admin credentials are
  admin/admin that you will be requested to change upon your first login
- The light faucet on port 8088. Just visit http://localhost:8088, enter a
  nodeID of a client, and submit it. 

It takes a few minutes until Geth manages to send data to InfluxDB, so
monitoring data might be not immediately visible.

Running it:
- Make sure you have a `recaptcha_v2_public.txt` and a `recaptcha_v2_secret.txt` in your folder
- Run `docker swarm init`
- Run `docker-compose -f docker-compose.test.yml up`