package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gocql/gocql"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

type Server struct {
	db *gocql.Session
}

func NewServer(scylla *gocql.Session) Server {
	return Server{db: scylla}
}

func (s *Server) ContentJson(handler http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		handler.ServeHTTP(w, r)
	}
}

// Recebe todos os utilizadores existentes
// se não existirem []
func (s *Server) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := s.db.Query("SELECT * FROM test.users").Iter().SliceMap()
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
func (s *Server) InsertUsers(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) HealthCheck(w http.ResponseWriter, r *http.Request) {
	message := map[string]string{
		"message": "Server is running",
	}
	resp, _ := json.Marshal(message)
	w.Write(resp)
}

func initJaeger(serviceName, jaegerAgentEndpoint string) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(fmt.Sprintf("%s/api/traces", jaegerAgentEndpoint))))
	if err != nil {
		return
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(serviceName),
				attribute.String("environment", "production"),
				attribute.String("ID", "1"),
			),
		),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

}

func main() {

	// Initialize Scylladb
	cluster := gocql.NewCluster("localhost:9042")
	// Define the keyspace create from "keyspace.sh"
	cluster.Keyspace = "test"

	session, err := cluster.CreateSession()
	if err != nil {
		panic(err)
	}

	initJaeger("Sample Scylla", "http://localhost:14268")

	s := NewServer(session)

	http.HandleFunc("/get-all", otelhttp.NewHandler(s.ContentJson(s.GetUsers), "Get All").ServeHTTP)

	http.HandleFunc("/insert-one", otelhttp.NewHandler(s.ContentJson(s.InsertUsers), "Insert One").ServeHTTP)

	http.HandleFunc("/healthcheck", otelhttp.NewHandler(s.ContentJson(s.HealthCheck), "healthcheck").ServeHTTP)

	fmt.Println("Server is running on port 8000")
	http.ListenAndServe(":8000", nil)
}
