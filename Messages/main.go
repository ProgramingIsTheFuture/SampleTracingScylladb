package main

func main() {
	session := Scylladb("host.docker.internal:9042", "messages_service")
	initJaeger("Messages", "http://localhost:14268")

	go StartGrpc(session)

	Routes(session)

}
