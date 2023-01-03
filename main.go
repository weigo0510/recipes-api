package main

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
	"os"
	"recipes-api/handlers"
)

var recipesHandler *handlers.RecipesHandler
var authHandler *handlers.AuthHandler

func init() {
	//recipes := make([]models.Recipe, 0)
	//file, _ := ioutil.ReadFile("recipes.json")
	//_ = json.Unmarshal([]byte(file), &recipes)
	ctx := context.Background()
	client, err := mongo.Connect(ctx,
		options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	recipesHandler = handlers.NewRecipesHandler(ctx, collection, redisClient)
	authHandler = &handlers.AuthHandler{}
	//var lRecipes []interface{}
	//for _, r := range recipes {
	//	lRecipes = append(lRecipes, r)
	//}
	////collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")
	//result, err := collection.InsertMany(ctx, lRecipes)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Println("Inserted recipes: ", len(result.InsertedIDs))
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenVal := c.GetHeader("Authorization")
		claims := &handlers.Claims{}
		tkn, err := jwt.ParseWithClaims(tokenVal, claims,
			func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_SECRET")), nil
			})
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		if tkn == nil || !tkn.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		c.Next()
	}
}

func main() {
	router := gin.Default()
	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	//4router.GET("/recipes/search", SearchRecipesHandler)
	router.POST("/signin", authHandler.SignInHandler)
	authorized := router.Group("/")
	authorized.Use(AuthMiddleware())
	{
		authorized.POST("/recipes", recipesHandler.NewRecipeHandler)
		authorized.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
		authorized.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)
	}
	router.Run()
}
