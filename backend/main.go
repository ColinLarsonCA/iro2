package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/ColinLarsonCA/iro2/backend/collabcafe"
	"github.com/ColinLarsonCA/iro2/backend/greeting"
	"github.com/ColinLarsonCA/iro2/backend/pb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

var (
	db       *sql.DB
	httpPort = 8090
	grpcPort = 9090
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var err error
	connStr := "host=postgres user=postgres password=postgres dbname=iro2 sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	grpcServer := grpc.NewServer()
	pb.RegisterGreetingServiceServer(grpcServer, greeting.NewService(db))
	pb.RegisterCollabCafeServiceServer(grpcServer, collabcafe.NewService(db))
	reflection.Register(grpcServer)
	err = pb.RegisterGreetingServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", grpcPort), opts)
	if err != nil {
		log.Fatalf("failed to register GreetingServiceHandler: %+v\n", err)
	}
	err = pb.RegisterCollabCafeServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", grpcPort), opts)
	if err != nil {
		log.Fatalf("failed to register CollabCafeServiceHandler: %+v\n", err)
	}
	go listenAndServe(grpcServer, grpcPort)
	log.Println("starting http server on port", httpPort)
	withCors := cors.New(cors.Options{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"ACCEPT", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler(mux)
	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", httpPort), withCors)
}

type ServerInterface interface {
	Serve(net.Listener) error
}

func listenAndServe(server ServerInterface, port int) {
	conn, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen on port %d: %v\n", port, err)
	}
	log.Println("starting grpc server on port", grpcPort)
	err = server.Serve(conn)
	if err != nil {
		log.Fatalf("failed to serve on port %d: %v\n", port, err)
	}
}
