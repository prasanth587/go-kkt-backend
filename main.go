package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/rs/cors"

	"go-transport-hub/dbconn/mssqlcon"
	"go-transport-hub/routes"

	_ "github.com/swaggo/http-swagger/example/go-chi/docs"
)

// @title Go Transport API
// @version 1.0
// @description logistics application
// @host localhost:9005
// @BasePath /
func main() {
	log.Println("-----------------------------------------------------------")
	log.Println(time.Now().In(timeLoc()))
	//mssqlurl := "3.8.31.220:cbAdmin2018@/cyberliver_platform?parseTime=true"
	//MsSqlurl := "root:cbAdmin2018@tcp(localhost:3306)/cyberliver_platform?charset=utf8"
	mssqlurl := os.Getenv("DB_CONNECTION_CONN")
	log.Printf("mssqlurl-con %s", mssqlurl)
	if mssqlurl == "" {
		log.Fatalln("no db connection found!, exiting")
	}

	mssqlcon.MSSqlInit(mssqlurl)
	log.Println("-----------------------------------------------------------")

	//router := routes.RouterConfig()
	router := routes.RouterConfig()
	//r := chi.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
		MaxAge:           86400, // 24 hours
		Debug:            false,
	})

	server := http.Server{
		Addr:         fmt.Sprintf(":%d", 9005),
		ReadTimeout:  90 * time.Second,
		WriteTimeout: 90 * time.Second,
		Handler:      c.Handler(router),
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	//Graceful shut down
	go func() {
		<-quit
		log.Println("Server is shutting down...")

		//Close resources before shut down
		//mssqlcon.MSSqlConnClose()

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		//Shutdown server
		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Unable to gracefully shutdown the server: %v\n", err)
		}

		//Close channels
		close(quit)
		close(done)
	}()

	log.Printf("Listening on: %d", 9005)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Error in listening server: %s", err.Error())
	}
	<-done
	log.Fatal("Server stopped")
}

func timeLoc() *time.Location {
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return nil
	}
	return loc
}
