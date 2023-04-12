package main

import (
	"bytes"
	"fmt"
	"log"
	apigin "projectx/api"
	"projectx/core"
	"projectx/crypto"
	"projectx/network"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
func main() {
	router := gin.Default()
	router.Use(CORSMiddleware())
	// router.Use(middleware.Recover())
	// appCtx := component.NewAppCtx(db)

	v1 := router.Group("/v1")
	{
		tx := v1.Group("/tx")
		{
			tx.POST("", apigin.CreateTxHandler())
			tx.GET("", apigin.GetTxHandler())
		}
	}
	router.Run(":9001")

	// trLocal := network.NewLocalTransport("LOCAL")
	// trRemoteA := network.NewLocalTransport("REMOTE_A")
	// trRemoteB := network.NewLocalTransport("REMOTE_B")
	// trRemoteC := network.NewLocalTransport("REMOTE_C")
	// trRemoteD := network.NewLocalTransport("REMOTE_D")
	// fmt.Println(trRemoteD)
	// // Local <-> A -> B -> C
	// trLocal.Connect(trRemoteA)
	// trRemoteA.Connect(trRemoteB)
	// trRemoteB.Connect(trRemoteC)
	// trRemoteA.Connect(trLocal)
	//
	// // Local is validator and have privateKey -> A,B,C is normal node
	// initRemoteServers([]network.Transport{trRemoteA, trRemoteB, trRemoteC})
	//
	// go func() {
	// 	for {
	// 		if err := sendTransaction(trRemoteA, trLocal.Addr()); err != nil {
	// 			logrus.Error(err)
	// 		}
	// 		time.Sleep(2 * time.Second)
	// 	}
	// }()
	//
	// go func() {
	// 	time.Sleep(7 * time.Second)

	// 	trLate := network.NewLocalTransport("LATE_REMOTE")
	// 	trRemoteC.Connect(trLate)
	// 	lateServer := makeServer(string(trLate.Addr()), trLate, nil)

	// 	go lateServer.Start()
	// }()

	// privKey := crypto.GeneratePrivateKey()
	// localServer := makeServer("LOCAL", trLocal, &privKey)
	// localServer.Start()

}

func initRemoteServers(trs []network.Transport) {
	for i := 0; i < len(trs); i++ {
		id := fmt.Sprintf("REMOTE_%d", i)
		s := makeServer(id, trs[i], nil)
		go s.Start()
	}
}

func makeServer(id string, tr network.Transport, pk *crypto.PrivateKey) *network.Server {
	opts := network.ServerOpts{
		PrivateKey: pk,
		ID:         id,
		Transports: []network.Transport{tr},
	}

	s, err := network.NewServer(opts)
	if err != nil {
		log.Fatal(err)
	}

	return s
}

func sendTransaction(tr network.Transport, to network.NetAddr) error {
	privKey := crypto.GeneratePrivateKey()
	tx := core.NewTransaction(contract())
	tx.Sign(privKey)
	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
		return err
	}

	msg := network.NewMessage(network.MessageTypeTx, buf.Bytes())

	return tr.SendMessage(to, msg.Bytes())
}

// FOO 5 - > FOO =5 -> FOO = 5
func contract() []byte {
	data := []byte{0x02, 0x0a, 0x03, 0x0a, 0x0c, 0x4f, 0x0b, 0x4f, 0x0b, 0x46, 0x0b, 0x03, 0x0a, 0x0e, 0x0f}
	pushFoo := []byte{0x4f, 0x0b, 0x4f, 0x0b, 0x46, 0x0b, 0x03, 0x0a, 0x0e, 0x001}

	data = append(data, pushFoo...)
	return data
}

// 2 3 + - > 5 O O F -> 5 [FOO]
