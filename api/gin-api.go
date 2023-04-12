package apigin

import (
	"bytes"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"projectx/biz"
	"projectx/core"
	"projectx/crypto"
	"projectx/model"

	"github.com/gin-gonic/gin"
)

func CreateKeyHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		privKey := crypto.GeneratePrivateKey()
		privKeyBytes, err := x509.MarshalECPrivateKey(privKey.Get())
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		}
		privKeyStr := base64.StdEncoding.EncodeToString(privKeyBytes)
		privKeyResp := model.PriKeyResponse{}
		privKeyResp.Key = privKeyStr
		c.JSON(http.StatusOK, privKeyResp)
	}
}

func CreateTxHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data model.TransactionCreate
		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		storage, err := core.NewMemoryStorage()
		defer storage.Close()
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

func GetTxHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		hash := c.Query("hash")
		storage, err := core.NewMemoryStorage()
		defer storage.Close()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		}

		biz := biz.NewGetTxBiz(storage)
		tx, err := biz.GetTx([]byte(hash))

		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		}
		buf := bytes.NewBuffer(tx)

		decodedTx := core.Transaction{}
		decodedTx.Decode(core.NewGobTxDecoder(buf))

		dataTx := model.TransactionCreate{}
		err = json.Unmarshal(decodedTx.Data, &dataTx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		}

		responseTx := model.TransactionResponse{}
		responseTx.Provider = dataTx.Provider
		responseTx.Track = dataTx.Track
		responseTx.Signature = decodedTx.Signature
		c.JSON(http.StatusOK, responseTx)

	}
}

func CreateBlockHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data model.BlockCreate
		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		storage, err := core.NewMemoryStorage()
		defer storage.Close()
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
