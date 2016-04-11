package main

import (
	"fmt"
	"github.com/alexstoick/wow-statistics/controllers"
	"github.com/gin-gonic/gin"
	"gopkg.in/redis.v3"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

var client *redis.Client

func RedisMapper(c *gin.Context) {
	c.Set("redis", client)

	c.Next()
}

func main() {
	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		//Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	router := gin.Default()

	router.Use(RedisMapper)
	router.Use(CORSMiddleware())

	v1 := router.Group("v1/")
	{
		v1.GET("/", controllers.NewData)
	}

	router.Run(":3001")
}
