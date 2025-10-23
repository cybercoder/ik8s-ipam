package v1alpha1

import (
	controllers "github.com/cybercoder/ik8s-ipam/pkg/api/v1alpha1/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.RouterGroup) {
	ipamController := controllers.NewIpamController()
	ipam := r.Group("/")
	ipam.POST("/assignip", ipamController.AssignIP)
}
