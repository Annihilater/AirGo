package api

import (
	"AirGo/global"
	"AirGo/model"
	"AirGo/service"
	"AirGo/utils/encrypt_plugin"
	"AirGo/utils/response"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func AGGetNodeInfo(ctx *gin.Context) {
	//验证key
	if global.Server.System.MuKey != ctx.Query("key") {
		return
	}
	id := ctx.Query("id")
	//idInt64, _ := strconv.ParseInt(id, 10, 64)
	fmt.Println("id:", id)
	node, _, err := service.CommonSqlFind[model.Node, string, model.AGNodeInfo]("id = " + id)
	if err != nil {
		global.Logrus.Error("AGGetNodeInfo error,id=" + id)
		return
	}
	b, _ := json.Marshal(node)

	fmt.Println("AGGetNodeInfo:", string(b))
	ctx.JSON(200, node)

}
func AGReportNodeStatus(ctx *gin.Context) {
	//验证key
	if global.Server.System.MuKey != ctx.Query("key") {
		return
	}
	var AGNodeStatus model.AGNodeStatus
	err := ctx.ShouldBind(&AGNodeStatus)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println("AGNodeStatus:", AGNodeStatus)
	cacheStatus, ok := global.LocalCache.Get(strconv.FormatInt(AGNodeStatus.ID, 10) + "status")
	if ok && cacheStatus != nil {
		oldStatus := cacheStatus.(model.NodeStatus)
		fmt.Println("old status:", oldStatus)
		oldStatus.CPU = AGNodeStatus.CPU
		oldStatus.Mem = AGNodeStatus.Mem
		oldStatus.Disk = AGNodeStatus.Disk
		//oldStatus.Uptime=AGNodeStatus.Uptime
		fmt.Println("new status:", oldStatus)
		global.LocalCache.Set(strconv.FormatInt(AGNodeStatus.ID, 10)+"status", oldStatus, 2*time.Minute) //2分钟后过期
	}
	ctx.String(200, "success")
}

func AGGetUserlist(ctx *gin.Context) {
	//验证key
	if global.Server.System.MuKey != ctx.Query("key") {
		return
	}
	id := ctx.Query("id")
	fmt.Println("id:", id)

	//节点号 是否启用？
	//节点属于哪些goods
	nodeIDInt, _ := strconv.ParseInt(id, 10, 64)
	goods, err := service.FindGoodsByNodeID(nodeIDInt)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	//goods属于哪些用户
	var goodsArr []int64
	for _, v := range goods {
		goodsArr = append(goodsArr, v.ID)
	}
	fmt.Println("goodsArr:", goodsArr)
	var users []model.AGUserInfo
	err = global.DB.Model(&model.User{}).Where("goods_id in (?) and sub_status = ?", goodsArr, true).Find(&users).Error
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println("AGGetUserlist:", users)
	ctx.JSON(200, users)

}
func AGReportUserTraffic(ctx *gin.Context) {
	//验证key
	if global.Server.System.MuKey != ctx.Query("key") {
		return
	}

	var AGUserTraffic model.AGUserTraffic
	err := ctx.ShouldBind(&AGUserTraffic)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println("AGUserTraffic:", AGUserTraffic)

	// 处理流量统计
	var userIds []int64
	var userArr []model.User
	var trafficLog = model.TrafficLog{
		NodeID: AGUserTraffic.ID,
	}
	for _, v := range AGUserTraffic.UserTraffic {
		//每个用户流量
		var user model.User
		userIds = append(userIds, v.UID)
		user.ID = v.UID
		user.SubscribeInfo.U = v.Upload
		user.SubscribeInfo.D = v.Download
		userArr = append(userArr, user)
		//该节点总流量
		trafficLog.D = trafficLog.U + v.Upload
		trafficLog.U = trafficLog.D + v.Download

	}
	// 统计status
	go func(id, userAmount, u, d int64) {
		var nodeStatus = model.NodeStatus{
			ID:         id,
			UserAmount: userAmount,
			U:          float64(u),
			D:          float64(d),
			LastTime:   time.Now(),
			Status:     true,
		}
		var duration float64 = 60
		cacheStatus, ok := global.LocalCache.Get(strconv.FormatInt(id, 10) + "status")
		if ok && cacheStatus != nil {
			oldStatus := cacheStatus.(model.NodeStatus)
			duration = nodeStatus.LastTime.Sub(oldStatus.LastTime).Seconds()
		}
		nodeStatus.D, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", nodeStatus.D/1024/1024/duration*8), 64) //Byte--->Mbps
		nodeStatus.U, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", nodeStatus.U/1024/1024/duration*8), 64)
		global.LocalCache.Set(strconv.FormatInt(id, 10)+"status", nodeStatus, 2*time.Minute)

	}(AGUserTraffic.ID, int64(len(userIds)), trafficLog.U, trafficLog.D)
	//插入流量统计统计
	err = service.CommonSqlCreate[model.TrafficLog](trafficLog)

	if err != nil {
		global.Logrus.Error("插入流量统计统计error:", err)
		return
	}
	//更新用户流量信息
	if len(userArr) == 0 {
		return
	}
	err = service.UpdateUserTrafficInfo(userArr, userIds)
	if err != nil {
		global.Logrus.Error("更新用户流量信息error:", err)
		return
	}
	ctx.String(200, "success")

}

// reality x25519
func Createx25519(ctx *gin.Context) {
	str := encrypt_plugin.RandomString(43)
	pub, pri, err := encrypt_plugin.ExecuteX25519(str)
	if err != nil {
		global.Logrus.Error(err.Error())
		return
	}
	response.OK("Createx25519 success", model.AGREALITYx25519{PublicKey: pub, PrivateKey: pri}, ctx)
}
