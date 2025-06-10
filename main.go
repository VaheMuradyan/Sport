package main

import (
	"github.com/VaheMuradyan/Sport/db"
	"github.com/VaheMuradyan/Sport/router"
	"github.com/VaheMuradyan/Sport/user"
	"github.com/gin-gonic/gin"
)

func main() {
	database := db.ConnectDB()
	r := gin.Default()

	repo := user.NewUserRepo(database)
	service := user.NewUserService(repo)
	handler := user.NewUserHandler(service)

	router.RegisterRoutes(r, handler)

	r.Run(":8080")
}
