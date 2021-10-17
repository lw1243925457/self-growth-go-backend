package main

import (
	"fmt"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"seltGrowth/internal/growth_record/controller"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

type Product struct {
	Username    string    `json:"username" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	Category    string    `json:"category" binding:"required"`
	Price       int       `json:"price" binding:"gte=-1"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}

type productHandler struct {
	sync.RWMutex
	products map[string]Product
}

func newProductHandler() *productHandler {
	return &productHandler{
		products: make(map[string]Product),
	}
}

func (u *productHandler) Create(c *gin.Context) {
	u.Lock()
	defer u.Unlock()

	// 0. 参数解析
	var product Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. 参数校验
	if _, ok := u.products[product.Name]; ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("product %s already exist", product.Name)})
		return
	}
	product.CreatedAt = time.Now()

	// 2. 逻辑处理
	u.products[product.Name] = product
	log.Printf("Register product %s success", product.Name)

	// 3. 返回结果
	c.JSON(http.StatusOK, product)
}

func (u *productHandler) Get(c *gin.Context) {
	u.Lock()
	defer u.Unlock()

	product, ok := u.products[c.Param("name")]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Errorf("can not found product %s", c.Param("name"))})
		return
	}

	c.JSON(http.StatusOK, product)
}

func router() http.Handler {
	router := gin.Default()
	productHandler := newProductHandler()
	helloHandler := controller.NewHelloHandler()
	phoneUseController := controller.NewPhoneUseController()
	// 路由分组、中间件、认证
	v1 := router.Group("/v1")
	{
		productv0 := v1.Group("/products")
		{
			// 路由匹配
			productv0.POST("", productHandler.Create)
			productv0.GET(":name", productHandler.Get)
		}

		hello := v1.Group("/hello")
		{
			hello.GET("", helloHandler.Hello)
		}

		phoneUser := v1.Group("/phone")
		{
			phoneUser.POST("/useRecord", phoneUseController.UploadRecord)
			phoneUser.GET("/overview", phoneUseController.Overview)
		}
	}

	return router
}


func initMongodb() {
	// Setup the mgm default config
	err := mgm.SetDefaultConfig(nil, "phone_record", options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	initMongodb()
	var eg errgroup.Group

	// 一进程多端口
	insecureServer := &http.Server{
		Addr:         ":8080",
		Handler:      router(),
		ReadTimeout:  4 * time.Second,
		WriteTimeout: 9 * time.Second,
	}

	//secureServer := &http.Server{
	//	Addr:         "192.168.1.3:8443",
	//	Handler:      router(),
	//	ReadTimeout:  4 * time.Second,
	//	WriteTimeout: 9 * time.Second,
	//}

	eg.Go(func() error {
		err := insecureServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
		return err
	})

	//eg.Go(func() error {
	//	err := secureServer.ListenAndServeTLS("D:\\temp\\https\\server.crt", "D:\\temp\\https\\server_no_passwd.key")
	//	if err != nil && err != http.ErrServerClosed {
	//		log.Fatal(err)
	//	}
	//	return err
	//})

	if err := eg.Wait(); err != nil {
		log.Fatal(err)
	}
}