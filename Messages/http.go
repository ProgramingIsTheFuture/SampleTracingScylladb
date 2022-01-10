package main

import (
	"encoding/json"
	"net/http"

	"github.com/gocql/gocql"
)

type Server struct {
	db *gocql.Session
}

func (s *Server) GetAll(w http.ResponseWriter, r *http.Request) {
	messages, err := s.db.Query(`SELECT * FROM messages`).Iter().SliceMap()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	resp, _ := json.Marshal(messages)
	w.Write(resp)
}

func Routes(session *gocql.Session) {
	StartConsul("Messages-HTTP", 8001)
	s := Server{db: session}
	http.HandleFunc("/get-all", s.GetAll)

	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("Server is still alive"))

	})
	http.ListenAndServe(":8001", nil)
}
