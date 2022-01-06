# Simple example to practice Tracing with jaeger

### Run instructions

Start containers

```shell
docker-compose up -d
```

After the containers start we need to run the keyspace script to create a new keyspace test and a table users
```shell
chmod +x keyspace.sh
./keyspace.sh
```

Now we can run the program
```go
go run main.go
```

With this steps we can access the Jaeger dashboard and start using the tracing

```
http://localhost:16686
```

