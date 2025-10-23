package api

import (
	"github.com/cybercoder/ik8s-ipam/pkg/api/v1alpha1"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	v1alpha1Group := router.Group("/apis/ovn.ik8s.ir/v1alpha1")
	v1alpha1.SetupRoutes(v1alpha1Group)
}
