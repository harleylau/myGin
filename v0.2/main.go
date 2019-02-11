package main

import (
	"github.com/harleylau/myGin/v0.2/gin"
)

func v1IndexLoginfunc(c *gin.Context) {
	c.String(200, "login")
}

func v1IndexSubmitfunc(c *gin.Context) {
	c.String(200, "submit")
}

func main() {
	r := gin.New()
	r.Use(gin.Logger())

	// Simple group: v1
	v1 := r.Group("/v1")
	{
		v1.GET("/login", v1IndexLoginfunc)
		v1.GET("/submit", v1IndexSubmitfunc)
	}

	// Listen and server on 0.0.0.0:8080
	r.Run(":8080")
}
