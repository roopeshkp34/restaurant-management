package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/roopeshkp34/restaurant-management/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")
// var validate = validator.New()

func GetMenus() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		result, err := menuCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var allMenus []bson.M
		if err = result.All(ctx, &allMenus); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allMenus)

	}
}

func GetMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		menuId := c.Param("menu_id")
		var menu models.Menu

		err := menuCollection.FindOne(ctx, bson.M{"menu_id": menuId}).Decode(&menu)

		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch menu"})

		}
		c.JSON(http.StatusOK, menu)
	}
}

func CreateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var menu models.Menu
		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}
		validationError := validate.Struct(menu)
		if validationError != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationError.Error()})
			return
		}
		menu.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.ID = primitive.NewObjectID()
		menu.Menu_id = menu.ID.Hex()

		result, insertError := menuCollection.InsertOne(ctx, menu)
		if insertError != nil {
			msg := fmt.Sprintf("Menu item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		defer cancel()
		c.JSON(http.StatusOK, result)
		defer cancel()
	}
}

func UpdateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var menu models.Menu
		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}
		menuId := c.Param("menu_id")
		filter := bson.M{"menu_id": menuId}

		var updateObject primitive.D

		if menu.Start_date != menu.End_date != nil {
			if !inTimeSpan(*menu.Start_date, *menu.End_date, time.Now()) {
				msg := "Kindly retype the time"
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				defer cancel()
				return
			}
		}
		updateObject = append(updateObject, bson.E{"start_date": menu.Start_date})
		updateObject = append(updateObject, bson.E{"end_date": menu.End_date})

		if menu.Name != "" {
			updateObject = append(updateObject, bson.E{"name": menu.Name})
		}

		if menu.Category != "" {
			updateObject = append(updateObject, bson.E{"category": menu.Category})
		}

		menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObject = append(updateObject, bson.E{"updated_at": menu.Updated_at})

		upsert := true

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := menuCollection.UpdateOne(ctx, filter, bson.D{{"$set", updateObject}}, &opt)
		if err != nil {
			msg := " Failed to update menu"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		}

		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}
