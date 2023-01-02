package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"recipes-api/handlers"
)

var recipesHandler *handlers.RecipesHandler

func init() {
	/*recipes = make([]Recipe, 0)
	file, _ := ioutil.ReadFile("recipes.json")
	_ = json.Unmarshal([]byte(file), &recipes)*/
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
	//var lRecipes []interface{}
	//for _, r := range recipes {
	//	lRecipes = append(lRecipes, r)
	//}
	//collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")
	//result, err := collection.InsertMany(ctx, lRecipes)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Println("Inserted recipes: ", len(result.InsertedIDs))
}

func main() {
	router := gin.Default()
	router.POST("/recipes", recipesHandler.NewRecipeHandler)
	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	router.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
	router.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)
	//4router.GET("/recipes/search", SearchRecipesHandler)
	router.Run()
}
