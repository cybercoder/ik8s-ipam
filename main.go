package main

import (
	"github.com/cybercoder/ik8s-ipam/pkg/api"
	"github.com/cybercoder/ik8s-ipam/pkg/k8s"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	api.SetupRoutes(router)

	publicIpPoolInformer := k8s.CreatePublicIpPoolInformer()
	stopCh := make(chan struct{})
	// publicIpPoolInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
	// 	AddFunc:    events.OnAddPublicIpPool,
	// 	UpdateFunc: events.OnUpdatePublicIpPool,
	// 	DeleteFunc: events.OnDeletePublicIpPool,
	// })
	defer close(stopCh)
	go publicIpPoolInformer.Run(stopCh)

	router.Run("0.0.0.0:8000")
}
