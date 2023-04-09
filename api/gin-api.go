package apigin

import (
	"net/http"
	"projectx/biz"
	"projectx/core"
	"projectx/model"

	"github.com/gin-gonic/gin"
)

func CreateTxHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data model.TransactionCreate
		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		storage, err := core.NewMemoryStorage()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		}
		biz := biz.NewCreateTxBiz(storage)

		if err := biz.CreateTx(&data); err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, "ok")

	}
}
