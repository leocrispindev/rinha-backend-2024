package controller

import (
	"net/http"
	"rinha-backend-2024/api/internal/model"
	"rinha-backend-2024/api/internal/model/exception"
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
		c.JSON(http.StatusBadRequest, "Invalid client Id")
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

	client, errInterface := balance.InsertTransaction(id, transaction)

	if errInterface != nil {
		response := convertInterfaceToError(errInterface)

		c.JSON(response.Status, response.Data)
		c.Abort()
		return
	}

	c.JSON(200, client)

}

func HandlerExtract(c *gin.Context) {
	clientId := c.Param("id")

	id, err := util.StringToInt(clientId)

	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid client Id")
		c.Abort()
	}

	extract, errInterface := balance.GetExtractByUserId(id)

	if errInterface != nil {
		response := convertInterfaceToError(errInterface)

		c.JSON(response.Status, response.Data)
		c.Abort()
		return
	}

	c.JSON(200, extract)
}

func convertInterfaceToError(err interface{}) model.Response {

	switch e := err.(type) {
	case exception.TransactionError:
		return model.Response{
			Status: http.StatusInternalServerError,
			Data:   e.Error(),
		}

	case exception.UserNotFound:
		return model.Response{
			Status: http.StatusNotFound,
			Data:   e.Error(),
		}

	case exception.UnprocessableEntity:
		return model.Response{
			Status: http.StatusUnprocessableEntity,
			Data:   e.Error(),
		}
	default:
		return model.Response{
			Status: http.StatusInternalServerError,
			Data:   "unexpected error",
		}
	}
}
