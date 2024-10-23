package controller

import (
	"GOMS-BACKEND-GO/model"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OutingController struct {
	outingUseCase model.OutingUseCase
}

func NewOutingController(outingUseCase model.OutingUseCase) *OutingController {
	return &OutingController{
		outingUseCase: outingUseCase,
	}
}

func (controller *OutingController) OutingStudent(ctx *gin.Context) {
	outingUUIDStr := ctx.Param("outingUUID")

	outingUUID, err := uuid.Parse(outingUUIDStr)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid outing UUID"})
		return
	}

	if err := controller.outingUseCase.OutingStudent(ctx, context.Background(), outingUUID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)

}

func (controller *OutingController) ListOutingStudent(ctx *gin.Context) {
	outings, err := controller.outingUseCase.FindAllOutingStudent(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve outing-students"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"outing-students": outings})

}
