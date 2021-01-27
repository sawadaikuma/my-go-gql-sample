package main

import (
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"log"
	"my-go-gql-sample/graph"
	"my-go-gql-sample/graph/dataloader"
	"my-go-gql-sample/graph/generated"
	"my-go-gql-sample/util/middleware/auth"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/rs/cors"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const defaultPort = "8080"

func main() {
	// TODO: DBをリファクタリング
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	user := os.Getenv("MYSQL_USER")
	pass := os.Getenv("MYSQL_PASSWORD")
	protocol := os.Getenv("MYSQL_PROTOCOL")
	dbname := os.Getenv("MYSQL_DATABASE")
	dsn := user + ":" + pass + "@" + protocol + "/" + dbname + "?parseTime=true&charset=utf8"

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      false,       // Disable color
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		fmt.Printf("DB Open Error :%v", err)
		panic(err.Error())
	}

	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			fmt.Printf("DB Close Error :%v", err)
			panic(err.Error())
		}
		sqlDB.Close()
	}()

	fmt.Println(dsn)

	// NOTE: [complexityは加重計算を設定可能](https://gqlgen.com/reference/complexity/)
	complexity := generated.ComplexityRoot{}
	complexity.Todo.Tags = func(childComplexity int) int {
		fmt.Println("[childComplexity]", childComplexity)
		return childComplexity * 6
	}

	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: &graph.Resolver{
					DB: db,
				},
				Complexity: complexity,
			},
		),
	)

	// NOTE: X クエリのネスト制限
	// NOTE: X ネスト制限ではなく、クエリのフィールド数の制限
	// NOTE: O フィールド数ではなく、複雑さ(complexity)で制限する
	srv.Use(extension.FixedComplexityLimit(50))

	rowDBdriver, err := db.DB()
	if err != nil {
		fmt.Printf("DB Open Error :%v", err)
		panic(err.Error())
	}

	// TODO: Echoで動かす
	router := chi.NewRouter()

	authHandler := auth.Middleware(rowDBdriver)
	dlHandler := dataloader.Middleware(rowDBdriver)
	corsHandler := cors.New(cors.Options{
		// AllowedOrigins:   []string{"http://localhost:8080"},
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		Debug:            true,
	}).Handler
	router.Handle("/query", srv)
	router.Handle("/playground", playground.Handler("GraphQL playground", "/query"))

	router.Route("/", func(r chi.Router) {
		// fmt.Println(authHandler)
		// fmt.Println(dlHandler)
		r.Use(corsHandler)
		r.Use(authHandler)
		r.Use(dlHandler)
		r.Handle("/", srv)
	})

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
	log.Println(router)
}
