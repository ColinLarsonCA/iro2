package main

import (
	"database/sql"
	"io"
	"log"
	"net/http"

	"github.com/ColinLarsonCA/iro2/backend/greeting"
	"github.com/ColinLarsonCA/iro2/backend/pb"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

var db *sql.DB

func main() {
	var err error
	connStr := "host=postgres user=postgres password=postgres dbname=iro2 sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	greetingService := greeting.NewService(db)
	pb.RegisterGreetingServiceServer(grpcServer, greetingService)

	http.HandleFunc("/ping", ping)
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/greeting", func(w http.ResponseWriter, r *http.Request) {
		allowCors(w)
		res, err := greetingService.GetGreeting(r.Context(), &pb.GreetingRequest{})
		if err != nil {
			internalError(w, err)
			return
		}
		plainResponse(w, res.Message)
	})
	http.ListenAndServe(":8090", nil)
}

func ping(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "pong!!\n")
}

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello!!\n")
}

func allowCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func internalError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func plainResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "text/plain")
	io.WriteString(w, message)
}
