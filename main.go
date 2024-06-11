package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/roopeshkp34/restaurant-management/database"
	"github.com/roopeshkp34/restaurant-management/middleware"
	"github.com/roopeshkp34/restaurant-management/routes"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

func main() {
	// for stopping server
	stopchan := make(chan os.Signal)
	signal.Notify(stopchan, os.Interrupt)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	router.Use(middleware.Authentication)

	routes.FoodRoutes(router)
	routes.MenuRoutes(router)
	routes.TableRoutes(router)
	routes.OrderRoutes(router)
	routes.OrderItemRoutes(router)
	routes.InvoiceRoutes(router)
	router.Run(":" + port)

	<-stopchan
	log.Println("Shutting Down server..")
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	log.Println("Server gracefully stopped")

}
