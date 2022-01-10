package main

import (
	"SampleTraceScylla/pb"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gocql/gocql"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Server struct {
	db   *gocql.Session
	gRPC pb.MessageMethodsClient
}

func NewServer(scylla *gocql.Session) Server {
	return Server{db: scylla}
}

func (s Server) ContentJson(handler http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		handler.ServeHTTP(w, r)
	}
}

// Recebe todos os utilizadores existentes
// se não existirem []
func (s Server) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := s.db.Query("SELECT * FROM users").Iter().SliceMap()
	if err != nil {
		errMsg := map[string]string{"message": fmt.Sprintf("Erro: %s", err.Error())}
		resp, _ := json.Marshal(errMsg)
		w.Write(resp)
		return
	}

	resp, _ := json.Marshal(users)
	w.Write(resp)
}

// Cria um novo utilizador se receber informação do request {username: "..."}
func (s Server) InsertUsers(w http.ResponseWriter, r *http.Request) {
	var user map[string]string
	json.NewDecoder(r.Body).Decode(&user)

	uuid, _ := gocql.RandomUUID()
	username := user["username"]

	err := s.db.Query("INSERT INTO users (id, username) VALUES (?, ?)", uuid, username).Exec()
	if err != nil {
		errMsg := map[string]string{"message": fmt.Sprintf("Erro: %s", err.Error())}
		resp, _ := json.Marshal(errMsg)
		w.Write(resp)
		return
	}

	user["id"] = uuid.String()

	resp, _ := json.Marshal(user)
	w.Write(resp)
}

// Send a msg over gRPC if a content and user_id "id" is passed over
func (s Server) SendMessage(w http.ResponseWriter, r *http.Request) {
	var body = map[string]string{}
	var resp []byte
	json.NewDecoder(r.Body).Decode(&body)
	if body["id"] == "" {
		errMsg := map[string]string{"message": "Must contain UserID"}
		resp, _ = json.Marshal(errMsg)
		w.Write(resp)
		return
	}

	var user_db = map[string]interface{}{}
	err := s.db.Query("SELECT id, username FROM users WHERE id=?", body["id"]).MapScan(user_db)
	if err != nil {
		if err == gocql.ErrNotFound {
			errMsg := map[string]string{"message": "User does not exist"}
			resp, _ = json.Marshal(errMsg)
			w.Write(resp)
			return
		}
		errMsg := map[string]string{"message": "Err: " + err.Error()}
		resp, _ = json.Marshal(errMsg)
		w.Write(resp)
		return
	}
	id, _ := gocql.ParseUUID(
		fmt.Sprintf(
			"%v",
			user_db["id"],
		),
	)
	msg := pb.Message{User: id.String(), Content: body["content"]}
	_, err = s.gRPC.Send(context.Background(), &msg)
	if err != nil {
		errMsg := map[string]string{"message": "Err: " + err.Error()}
		resp, _ = json.Marshal(errMsg)
		w.Write(resp)
		return
	}

	// msgSuc := map[string]string{"message": "Created "}
	resp, _ = json.Marshal(user_db)
	w.Write(resp)
}

func (s Server) HealthCheck(w http.ResponseWriter, r *http.Request) {
	message := map[string]string{
		"message": "Server is running",
	}
	resp, _ := json.Marshal(message)
	w.WriteHeader(200)
	w.Write(resp)
}

func Routes(s Server) {
	http.HandleFunc("/get-all", otelhttp.NewHandler(s.ContentJson(s.GetUsers), "Get All").ServeHTTP)

	http.HandleFunc("/insert-one", otelhttp.NewHandler(s.ContentJson(s.InsertUsers), "Insert One").ServeHTTP)

	http.HandleFunc("/send-msg", otelhttp.NewHandler(s.ContentJson(s.SendMessage), "Send Message").ServeHTTP)

	http.HandleFunc("/healthcheck", otelhttp.NewHandler(s.ContentJson(s.HealthCheck), "healthcheck").ServeHTTP)

	fmt.Println("Server is running on port 8000")
	http.ListenAndServe(":8000", nil)
}
