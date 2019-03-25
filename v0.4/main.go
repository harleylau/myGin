package main

import (
	"github.com/harleylau/myGin/v0.4/gin"
)

// LoginJSON .
type LoginJSON struct {
	User     string `json:"user" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func v1PasswordLoginfunc(c *gin.Context) {
	var json LoginJSON
	// If EnsureBody returns false, it will write automatically the error
	// in the HTTP stream and return a 400 error. If you want custom error
	// handling you should use: c.ParseBody(interface{}) error
	if c.EnsureBody(&json) {
		if json.User == "harleylau" && json.Password == "password" {
			c.JSON(200, gin.H{"status": "you are logged in"})
		} else {
			c.JSON(401, gin.H{"status": "unauthorized"})
		}
	}
}

func v1IndexLoginfunc(c *gin.Context) {
	name := c.Params.ByName("name")
	c.Set("innerName", name)
	message := getInfo(c)
	c.String(200, message)
}

func getInfo(c *gin.Context) string {
	name := c.Get("innerName")
	message := "welcome " + name.(string) // 前提是知道这个是string类型
	return message
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
		v1.GET("/login/:name", v1IndexLoginfunc)
		v1.POST("/pwlogin", v1PasswordLoginfunc)
		v1.GET("/submit", v1IndexSubmitfunc)
	}

	// Listen and server on 0.0.0.0:8080
	r.Run(":8080")
}
