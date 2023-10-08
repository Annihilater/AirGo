package initialize

import (
	"AirGo/global"
	"AirGo/service"
)

func InitAlipayClient() {
	client, err := service.InitAlipayClient()
	if err != nil {
		global.Logrus.Error("init alipay client error:", err)
		return
	}
	global.AlipayClient = client
}
