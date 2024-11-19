package controller

import (
	"GOMS-BACKEND-GO/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LateController struct {
	lateUseCase model.LateUseCase
}

func NewLateController(lateUseCase model.LateUseCase) *LateController {
	return &LateController{
		lateUseCase: lateUseCase,
	}
}

func (contoller *LateController) GetLateStudentTop3(ctx *gin.Context) {
	lates, err := contoller.lateUseCase.GetTop3LateStudent(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"outing-students": lates})

}
