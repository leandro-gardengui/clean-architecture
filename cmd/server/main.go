package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"clean-architecture/internal/infra/database"
	graph "clean-architecture/internal/infra/graphql"

	// "clean-architecture/internal/infra/grpc/pb"
	// "clean-architecture/internal/infra/grpc/service"
	"clean-architecture/internal/infra/web"
	"clean-architecture/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	graphqlLib "github.com/graphql-go/graphql"
	_ "github.com/lib/pq"
	// "google.golang.org/grpc"
	// "google.golang.org/grpc/reflection"
)

const (
	webPort = ":8000"
	// grpcPort    = ":50051"
	graphqlPort = ":8080"
)

func main() {
	// Database connection
	db, err := sql.Open("postgres", "postgres://user:password@db:5432/orders?sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Run migrations
	runMigrations(db)

	// Setup dependencies
	orderRepository := database.NewOrderRepository(db)
	createOrderUseCase := usecase.NewCreateOrderUseCase(orderRepository)
	listOrdersUseCase := usecase.NewListOrdersUseCase(orderRepository)

	// Start services
	go startWebServer(createOrderUseCase, listOrdersUseCase)
	// go startGRPCServer(createOrderUseCase, listOrdersUseCase)
	startGraphQLServer(createOrderUseCase, listOrdersUseCase)
}

func runMigrations(db *sql.DB) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal("Could not create database driver:", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		log.Fatal("Could not create migration instance:", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Could not run migrations:", err)
	}

	log.Println("Migrations completed successfully")
}

func startWebServer(createOrderUseCase *usecase.CreateOrderUseCase, listOrdersUseCase *usecase.ListOrdersUseCase) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	orderHandler := web.NewWebOrderHandler(createOrderUseCase, listOrdersUseCase)

	router.POST("/order", orderHandler.CreateOrder)
	router.GET("/order", orderHandler.ListOrders)

	log.Printf("Web server starting on port %s", webPort)
	log.Fatal(http.ListenAndServe(webPort, router))
}

// func startGRPCServer(createOrderUseCase *usecase.CreateOrderUseCase, listOrdersUseCase *usecase.ListOrdersUseCase) {
//	lis, err := net.Listen("tcp", grpcPort)
//	if err != nil {
//		log.Fatal("Failed to listen on port", grpcPort, err)
//	}
//
//	grpcServer := grpc.NewServer()
//	orderService := service.NewOrderService(createOrderUseCase, listOrdersUseCase)
//	pb.RegisterOrderServiceServer(grpcServer, orderService)
//	reflection.Register(grpcServer)
//
//	log.Printf("gRPC server starting on port %s", grpcPort)
//	if err := grpcServer.Serve(lis); err != nil {
//		log.Fatal("Failed to serve gRPC server:", err)
//	}
// }

func startGraphQLServer(createOrderUseCase *usecase.CreateOrderUseCase, listOrdersUseCase *usecase.ListOrdersUseCase) {
	graphService := graph.NewGraph(createOrderUseCase, listOrdersUseCase)
	schema, err := graphService.Schema()
	if err != nil {
		log.Fatal("Failed to create GraphQL schema:", err)
	}

	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		if r.Method != "POST" {
			http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
			return
		}

		var query struct {
			Query string `json:"query"`
		}

		if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result := graphqlLib.Do(graphqlLib.Params{
			Schema:        schema,
			RequestString: query.Query,
		})

		if err := json.NewEncoder(w).Encode(result); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// GraphQL Playground
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <title>GraphQL Playground</title>
  <style>
    body { font-family: Arial, sans-serif; margin: 40px; }
    .container { max-width: 800px; margin: 0 auto; }
    pre { background: #f4f4f4; padding: 20px; border-radius: 5px; }
  </style>
</head>
<body>
  <div class="container">
    <h1>GraphQL Playground</h1>
    <p>Send POST requests to <strong>/graphql</strong></p>
    
    <h3>Example Queries:</h3>
    
    <h4>List Orders:</h4>
    <pre>
{
  listOrders {
    id
    price
    tax
    final_price
  }
}
    </pre>
    
    <h4>Create Order:</h4>
    <pre>
mutation {
  createOrder(price: 100.0, tax: 10.0) {
    id
    price
    tax
    final_price
  }
}
    </pre>
  </div>
</body>
</html>
		`)))
	})

	log.Printf("GraphQL server starting on port %s", graphqlPort)
	log.Fatal(http.ListenAndServe(graphqlPort, nil))
}
