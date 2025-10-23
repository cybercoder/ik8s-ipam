package main

import (
	"github.com/cybercoder/ik8s-ipam/pkg/api"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	api.SetupRoutes(router)
	router.Run("0.0.0.0:8000")
}
