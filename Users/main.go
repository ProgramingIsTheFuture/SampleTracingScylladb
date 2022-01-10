package main

func main() {
	s := NewServer(Scylladb("host.docker.internal:9042", "users_service"))
	s.gRPC = ConnectGrpc()

	StartConsul("users", 8000)
	initJaeger("Users", "http://host.docker.internal:14268")

	Routes(s)

}
