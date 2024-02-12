package controller

import (
	"net/http"
	"rinha-backend-2024/api/internal/model"
	"rinha-backend-2024/api/internal/services/balance"
	"rinha-backend-2024/api/internal/services/util"

	"github.com/gin-gonic/gin"
)

func Init() {

}

func HandlerTransaction(c *gin.Context) {
	clientId := c.Param("id")

	id, err := util.StringToInt(clientId)

	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid client ID")
		c.Abort()
	}

	transaction := model.Transaction{}

	err = c.BindJSON(&transaction)

	validateError := transaction.Validate()

	if err != nil || validateError != nil {
		c.JSON(http.StatusBadRequest, "Invalid request body: "+err.Error())
		c.Abort()
		return
	}

	client, err := balance.InsertTransaction(id, transaction)

	if err != nil {
		c.JSON(500, err.Error())
		c.Abort()
		return
	}

	c.JSON(200, client)

}

func HandlerExtract(c *gin.Context) {
	clientId := c.Param("id")

	id, err := util.StringToInt(clientId)

	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid client ID")
		c.Abort()
	}

	extract, err := balance.GetExtractByUserId(id)

	if err != nil {
		c.JSON(500, err.Error())
		c.Abort()
		return
	}

	c.JSON(200, extract)
}
