# Sproot
Engine for store data in graph database and send notifiy.


```docker
docker pull dgraph/dgraph

# Directory to store data in. This would be passed to `-v` flag.
mkdir -p /tmp/data

# Run Dgraph Zero
docker run -it -p 8080:8080 -p 9080:9080 -v /tmp/data:/dgraph --name diggy dgraph/dgraph dgraph zero --port_offset -2000

# Run Dgraph Server
docker exec -it diggy dgraph server --memory_mb 2048 --zero localhost:5080

```

```
// For test
go get -u
go test ./...
```

## With DockerCompose for use with docker-compose.yaml
```
docker-compose up -d

docker-compose stop
```