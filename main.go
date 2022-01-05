package main

import (
	"encoding/json"
	"net/http"
)

type Server struct {
}

func NewServer() Server {
	return Server{}
}

func (s *Server) ContentJson(handler http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		handler.ServeHTTP(w, r)
	}

}

func (s *Server) GetUsers(w http.ResponseWriter, r *http.Request) {
}

func (s *Server) InsertUsers(w http.ResponseWriter, r *http.Request) {
}

func (s *Server) HealthCheck(w http.ResponseWriter, r *http.Request) {
	message := map[string]string{
		"message": "Server is running",
	}
	resp, _ := json.Marshal(message)
	w.Write(resp)
}

func main() {

	s := NewServer()

	http.HandleFunc("/get-all", s.ContentJson(s.GetUsers))

	http.HandleFunc("/insert-one", s.ContentJson(s.InsertUsers))

	http.HandleFunc("/healthcheck", s.ContentJson(s.HealthCheck))

	http.ListenAndServe(":8000", nil)
}
