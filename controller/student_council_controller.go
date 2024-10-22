package controller

import (
	"GOMS-BACKEND-GO/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StudentCouncilController struct {
	studentCouncilUseCase model.StudentCouncilUseCase
}

func NewStudentCouncilController(studentCouncilUseCase model.StudentCouncilUseCase) *StudentCouncilController {
	return &StudentCouncilController{
		studentCouncilUseCase: studentCouncilUseCase,
	}
}

func (controller *StudentCouncilController) CreateOuting(ctx *gin.Context) {

	newUUID, err := controller.studentCouncilUseCase.CreateOuting(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"outingUUID": newUUID.String()})

}
