# Sproot
Engine for store data in graph database and send notifiy.

Sproot use [Cayley](https://cayley.io/) for store data

## Setup
Need [Cockroach Database](https://www.cockroachlabs.com/docs/install-cockroachdb.html)

```docker
docker pull cockroachdb/cockroach

// For test
docker run -p 8080:8080 -p 26257:26257 cockroachdb/cockroach start --insecure
```