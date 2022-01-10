package main

import "github.com/gocql/gocql"

func Scylladb(host, keyspace string) *gocql.Session {
	// Initialize Scylladb
	cluster := gocql.NewCluster(host)
	// Define the keyspace create from "keyspace.sh"
	cluster.Keyspace = keyspace

	session, err := cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	return session
}
