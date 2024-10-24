package controller

import (
	"GOMS-BACKEND-GO/model"
	"GOMS-BACKEND-GO/model/data/constant"
	"GOMS-BACKEND-GO/model/data/input"
	"fmt"
	"net/http"
	"strconv"

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

	outingUUID, err := controller.studentCouncilUseCase.CreateOuting(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"outingUUID": outingUUID.String()})

}

func (controller *StudentCouncilController) FindOutingList(ctx *gin.Context) {

	accounts, err := controller.studentCouncilUseCase.FindAllAccount(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"accounts": accounts})
}

func (controller *StudentCouncilController) SearchAccountByInfo(ctx *gin.Context) {

	grade := ctx.Query("grade")
	gender := ctx.Query("gender")
	name := ctx.Query("name")
	isBlackList := ctx.Query("isBlackList")
	authority := ctx.Query("authority")
	major := ctx.Query("major")

	var input input.SearchAccountInput

	if grade != "" {
		grade, err := strconv.Atoi(grade)
		if err == nil {
			input.Grade = &grade
		}
	}

	if gender != "" {
		gender := constant.Gender(gender)
		input.Gender = &gender
	}

	if name != "" {
		input.Name = &name
	}

	if isBlackList != "" {
		isBlackList, err := strconv.ParseBool(isBlackList)
		if err == nil {
			input.IsBlackList = &isBlackList
		}
	}

	if authority != "" {
		authority := constant.Authority(authority)
		input.Authority = &authority
	}

	if major != "" {
		major := constant.Major(major)
		input.Major = &major
	}

	accounts, err := controller.studentCouncilUseCase.SearchAccount(ctx, &input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"eror": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"accounts": accounts})
}

func (controller *StudentCouncilController) UpdateAuthority(ctx *gin.Context) {
	var input input.UpdateAccountAuthorityInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := controller.studentCouncilUseCase.UpdateAccountAuthority(ctx, &input)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (controller *StudentCouncilController) AddBlackList(ctx *gin.Context) {
	accountIDParam := ctx.Param("accountID")

	accountID, err := strconv.ParseUint(accountIDParam, 10, 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid accountID"})
		return
	}

	fmt.Println(accountID)
	err = controller.studentCouncilUseCase.AddBlackList(ctx, accountID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
