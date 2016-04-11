package controllers

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/redis.v3"
)

func FetchRedisFromContext(c *gin.Context) *redis.Client {
	fake_redis, _ := c.Get("redis")
	db := fake_redis.(*redis.Client)
	return db
}
