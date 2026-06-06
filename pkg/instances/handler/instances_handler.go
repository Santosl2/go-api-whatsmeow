package instances

import (
	"github.com/Santosl2/go-api-whatsmeow/pkg/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type InstancesHandler struct {
	whatsmeowService *services.WhatsmeowService
}

func NewInstancesHandler(whatsmeowService *services.WhatsmeowService) *InstancesHandler {
	return &InstancesHandler{
		whatsmeowService: whatsmeowService,
	}
}

func (h *InstancesHandler) GetInstances(ctx *gin.Context) {
	schema := struct {
		Name string `json:"name"`
	}{}

	if err := ctx.ShouldBindJSON(&schema); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	clients := h.whatsmeowService.GetClients()
	if len(clients) == 0 {
		ctx.JSON(200, gin.H{"message": "No instances found"})
		return
	}

	ctx.JSON(200, gin.H{"message": "GetInstances endpoint hit", "result": clients})
}

func (h *InstancesHandler) CreateInstance(ctx *gin.Context) {
	schema := struct {
		Name string `json:"name"`
	}{}

	if err := ctx.ShouldBindJSON(&schema); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	randomID := uuid.New()
	h.whatsmeowService.StartNewConnection(randomID.String())
	ctx.JSON(200, gin.H{"message": "CreateInstance endpoint hit", "result": schema})
}
