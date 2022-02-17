package main

import (
	reg "traffic-subscribe/nacos_go"

	"github.com/gin-gonic/gin"
)

func main() {

	var nacosIPs = []string{"127.0.0.1"}

	// 发布配置
	var cfg = reg.ConfigParam{
		DataId:  "userview-uv.local",
		Group:   "Userview-Dev",
		Content: "hello2222wwww",
	}
	reg.ClientPublishConfig(nacosIPs, cfg)

	// 注册实例
	var svc = reg.ServiceParam{
		Ip:          "172.16.100.30",
		Port:        8080,
		ServiceName: "traffic-subscriber",
		GroupName:   "Userview",
		Weight:      30,
	}
	reg.ClientRegisterInstance(nacosIPs, svc)

	r := gin.Default()
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Pong",
		})
	})
	r.Run("0.0.0.0:8080")
}
