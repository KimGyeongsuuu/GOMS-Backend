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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"outing-students": outings})

}

func (controller *OutingController) CountOutingStudent(ctx *gin.Context) {
	outingsCount, err := controller.outingUseCase.CountAllOutingStudent(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"outing-students-count": outingsCount})

}

func (controller *OutingController) SearchOutingStudent(ctx *gin.Context) {
	name := ctx.Query("name")

	outings, err := controller.outingUseCase.SearchOutingStudent(ctx, name)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"outing-students": outings})

}
