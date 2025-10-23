package controllers

import (
	"net/http"

	"github.com/cybercoder/ik8s-ipam/pkg/api/v1alpha1/services"
	types "github.com/cybercoder/ik8s-ipam/pkg/api/v1alpha1/types"
	"github.com/gin-gonic/gin"
)

type IpamController struct {
	ipamService *services.IpamService
}

func NewIpamController() *IpamController {
	return &IpamController{
		ipamService: services.NewIpamService(),
	}
}

func (ic *IpamController) AssignIP(ctx *gin.Context) {
	body := types.IpAssignmentRequestBody{}
	if err := ctx.ShouldBind(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := ic.ipamService.FindOrCreateIPAssignment(body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, result.Spec)
}
