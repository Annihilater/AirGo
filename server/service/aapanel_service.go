package service

import (
	"AirGo/global"
	"AirGo/model"
	"encoding/base64"
	"encoding/json"
	"gorm.io/gorm"
	"net/url"
	"strconv"
	"strings"

	//"gopkg.in/yaml.v2"
	"gopkg.in/yaml.v2"
)

// 节点信息
func SSNodeInfo(nodeID int64) (model.SSNodeInfo, error) {
	var node model.Node
	err := global.DB.Where("id = ? and enabled = true", nodeID).First(&node).Error //节点号 是否启用
	if err != nil {
		return model.SSNodeInfo{}, err
	}
	var nodeInfo model.SSNodeInfo
	//nodeInfo.NodeGroup = 0
	//nodeInfo.NodeClass = 0
	//nodeInfo.MuOnly = 1
	nodeInfo.NodeSpeedlimit = node.NodeSpeedlimit
	nodeInfo.TrafficRate = node.TrafficRate

	switch node.NodeType {
	case "vless":
		nodeInfo.Sort = 15
	case "vmess":
		nodeInfo.Sort = 11
	case "trojan":
		nodeInfo.Sort = 14

	}
	switch node.NodeType {
	case "vmess": //vmess
		if node.Type == "http" && node.Network == "tcp" {
			nodeInfo.Server = node.Address + ";" + strconv.FormatInt(node.Port, 10) + ";" + strconv.FormatInt(node.Aid, 10) + ";" + node.Network + ";" + node.Security + ";path=" + node.Path + "|host=" + node.Host + ";headertype=http"
		} else if node.Network == "grpc" && node.Security != "" {
			nodeInfo.Server = node.Address + ";" + strconv.FormatInt(node.Port, 10) + ";" + strconv.FormatInt(node.Aid, 10) + ";" + node.Network + ";" + node.Security + ";path=" + node.Path + "|host=" + node.Host + "|servicename=mygrpc"
		}
		nodeInfo.Server = node.Address + ";" + strconv.FormatInt(node.Port, 10) + ";" + strconv.FormatInt(node.Aid, 10) + ";" + node.Network + ";" + node.Security + ";path=" + node.Path + "|host=" + node.Host
	case "vless": //vless
		if node.Type == "http" && node.Network == "tcp" {
			nodeInfo.Server = node.Address + ";" + strconv.FormatInt(node.Port, 10) + ";" + strconv.FormatInt(node.Aid, 10) + ";" + node.Network + ";" + node.Server + ";path=" + node.Path + "|host=" + node.Host + ";headertype=http" + "|enable_vless=true"
		} else if node.Network == "grpc" && node.Security != "" {
			nodeInfo.Server = node.Address + ";" + strconv.FormatInt(node.Port, 10) + ";" + strconv.FormatInt(node.Aid, 10) + ";" + node.Network + ";" + node.Security + ";path=" + node.Path + "|host=" + node.Host + "|servicename=mygrpc" + "|enable_vless=true"
		}
		nodeInfo.Server = node.Address + ";" + strconv.FormatInt(node.Port, 10) + ";" + strconv.FormatInt(node.Aid, 10) + ";" + node.Network + ";" + node.Security + ";path=" + node.Path + "|host=" + node.Host + "|enable_vless=true"

	case "trojan": //trojan
		if node.Network == "grpc" {
			nodeInfo.Server = node.Address + ":" + strconv.FormatInt(node.Port, 10) + "|host=" + node.Host + "|grpc=1|servicename=mygrpc"
		}
		nodeInfo.Server = node.Address + ":" + strconv.FormatInt(node.Port, 10) + "|host=" + node.Host
	}
	return nodeInfo, nil

}

// 获取订阅
func GetUserSub(url string, subType string) string {
	//查找用户
	var u model.User
	err := global.DB.Where("subscribe_url = ? and sub_status = 1 and d + u < t", url).First(&u).Error
	if err != nil {
		return ""
	}
	//根据goodsID 查找具体的节点
	var goods model.Goods
	err = global.DB.Where("id = ?", u.SubscribeInfo.GoodsID).Preload("Nodes", func(db *gorm.DB) *gorm.DB { return db.Order("node_order") }).Find(&goods).Error
	// 计算剩余天数，流量
	//fmt.Println("根据goodsID 查找具体的节点", goods)
	expiredTime := u.SubscribeInfo.ExpiredAt.Format("2006-01-02")
	expiredBd1 := (float64(u.SubscribeInfo.T - u.SubscribeInfo.U - u.SubscribeInfo.D)) / 1024 / 1024 / 1024
	expiredBd2 := strconv.FormatFloat(expiredBd1, 'f', 2, 64)
	name := "到期时间:" + expiredTime + "  |  剩余流量:" + expiredBd2 + "GB"
	var firstSubNode = model.Node{
		Remarks:  name,
		Address:  global.Server.System.SubName,
		Port:     6666,
		Aid:      0,
		Network:  "ws",
		Enabled:  true,
		NodeType: "vmess",
	}
	//插入计算剩余天数，流量的第一条节点
	goods.Nodes = append(goods.Nodes, model.Node{})
	copy(goods.Nodes[1:], goods.Nodes[0:])
	goods.Nodes[0] = firstSubNode
	//再插入共享的节点
	nodeList, err := CommonSqlFind[model.NodeShared, string, []model.Node](model.NodeShared{}, "")
	if err == nil {
		for _, v := range nodeList {
			goods.Nodes = append(goods.Nodes, v)
		}
	}
	//fmt.Println("nodes:", goods.Nodes)
	//根据subType生成不同客户端订阅 1:v2rayng 2:clash 3 shadowrocket 4 Quantumult X
	switch subType {
	case "1":
		return V2rayNGSubscribe(&goods.Nodes, u.UUID.String(), u.SubscribeInfo.Host)
	case "2":
		return ClashSubscribe(&goods.Nodes, u.UUID.String(), u.SubscribeInfo.Host)
	case "3":
		//return ShadowRocketSubscribe(&goods.Nodes, u.UUID.String(), u.SubscribeInfo.Host, name)
		return V2rayNGSubscribe(&goods.Nodes, u.UUID.String(), u.SubscribeInfo.Host)
	case "4":
		//return QxSubscribe(&goods.Nodes, u.UUID.String(), u.SubscribeInfo.Host)
		return V2rayNGSubscribe(&goods.Nodes, u.UUID.String(), u.SubscribeInfo.Host)
	}
	return ""
}

// v2rayNG 订阅
func V2rayNGSubscribe(nodes *[]model.Node, uuid, host string) string {
	var subArr []string
	for _, v := range *nodes {
		//剔除禁用节点
		if !v.Enabled {
			continue
		}
		if host == "" {
			host = v.Host
		}
		switch v.NodeType {
		case "vmess":
			if res := V2rayNGVmess(v, uuid, host); res != "" {
				subArr = append(subArr, res)
			}
			continue
		case "vless":
			if res := V2rayNGVlessTrojan(v, "vless", uuid, host); res != "" {
				subArr = append(subArr, res)
			}
			continue
		case "trojan":
			if res := V2rayNGVlessTrojan(v, "trojan", uuid, host); res != "" {
				subArr = append(subArr, res)
			}
			continue
		}
	}
	return base64.StdEncoding.EncodeToString([]byte(strings.Join(subArr, "\r\n")))
}

// clash 订阅
func ClashSubscribe(nodes *[]model.Node, uuid, host string) string {
	var proxiesArr []model.ClashProxy
	//所有节点名称数组
	var nameArr []string
	for _, v := range *nodes {
		//剔除禁用节点
		if !v.Enabled {
			continue
		}
		if host == "" {
			host = v.Host
		}
		//
		nameArr = append(nameArr, v.Remarks)

		proxy := ClashVmessVlessNew(v, uuid, host)
		proxiesArr = append(proxiesArr, proxy)

	}
	var proxyGroup model.ClashProxyGroup
	proxyGroup.Name = global.Server.System.SubName
	proxyGroup.Type = "select"
	proxyGroup.Proxies = nameArr

	var clashYaml model.ClashYaml
	clashYaml.Port = 7890
	clashYaml.SocksPort = 7891
	clashYaml.RedirPort = 7892
	clashYaml.AllowLan = false
	clashYaml.Mode = "rule"
	clashYaml.LogLevel = "silent"
	clashYaml.ExternalController = "0.0.0.0:9090"
	clashYaml.Secret = ""
	clashYaml.Proxies = proxiesArr
	clashYaml.ProxyGroups = append(clashYaml.ProxyGroups, proxyGroup)
	clashYaml.Rules = append(clashYaml.Rules, "MATCH,"+global.Server.System.SubName)
	res, err := yaml.Marshal(clashYaml)
	if err != nil {
		global.Logrus.Error("yaml.Marshal error:", err)
		return ""
	}
	return string(res)

}

// {"add":"AirGo","aid":"0","alpn":"h2,http/1.1","fp":"qq","host":"www.baidu.com","id":"e0d5fe65-a5d1-4b8a-8d40-ed92a6a35d8b","net":"ws","path":"/path","port":"6666","ps":"到期时间:2024-03-06  |  剩余流量:20.00GB","scy":"auto","sni":"www.baidu.com","tls":"tls","type":"","v":"2"}
// generate v2rayNG vmess
func V2rayNGVmess(node model.Node, uuid, host string) string {
	var vmess model.Vmess
	vmess.V = node.V
	vmess.Name = node.Remarks
	if node.EnableTransfer {
		vmess.Address = node.TransferAddress
		vmess.Port = strconv.FormatInt(node.TransferPort, 10)
	} else {
		vmess.Address = node.Address
		vmess.Port = strconv.FormatInt(node.Port, 10)
	}
	vmess.Uuid = uuid
	vmess.Aid = strconv.FormatInt(node.Aid, 10)
	vmess.Net = node.Network
	vmess.Disguisetype = node.Type //伪装类型
	vmess.Host = host
	vmess.Path = node.Path
	//传输层安全
	switch node.Security {
	case "tls":
		// alpn fp sni
		vmess.Tls = node.Security
		vmess.Alpn = node.Alpn
		vmess.Sni = node.Sni
	}

	vmessMarshal, err := json.Marshal(vmess)
	if err != nil {
		return ""
	}
	vmessStr := base64.StdEncoding.EncodeToString([]byte(vmessMarshal))
	return "vmess://" + vmessStr
}

// generate  v2rayng vless
// vless例子 vless://d342d11e-d424-4583-b36e-524ab1f0afa7@1.6.1.1:443?path=%2F%3Fed%3D2048&security=tls&encryption=none&alpn=h2,http/1.1&host=v2.airgoo.link&fp=randomized&flow=xtls-rprx-vision-udp443&type=ws&sni=v2.airgoo.link#v2.airgoo.link
// vless例子 vless://d342d11e-d424-4583-b36e-524ab1f0afa7@1.6.1.4:443?path=%2F%3Fed%3D2048&security=reality&encryption=none&pbk=ppkk&host=v2.airgoo.link&fp=randomized&spx=ssxx&flow=xtls-rprx-vision-udp443&type=ws&sni=v2.airgoo.link&sid=ssdd#v2.airgoo.link
// [scheme:][//[userinfo@]host][/]path[?query][#fragment]
func V2rayNGVlessTrojan(node model.Node, scheme, uuid, host string) string {
	var vlessUrl url.URL
	vlessUrl.Scheme = scheme
	vlessUrl.User = url.UserPassword(uuid, "")
	vlessUrl.Host = node.Address + ":" + strconv.FormatInt(node.Port, 10)
	values := url.Values{}

	switch scheme {
	case "vless":
		values.Add("encryption", node.Scy)
		values.Add("type", node.Network)
		values.Add("host", host)
		values.Add("path", node.Path)
		values.Add("flow", node.VlessFlow)
	case "trojan":
		values.Add("headerType", node.Type)
		values.Add("type", node.Network)
		values.Add("host", host)
		values.Add("path", node.Path)

	}
	switch node.Security {
	case "tls":
		values.Add("security", node.Security)
		values.Add("alpn", node.Alpn)
		values.Add("fp", node.Fingerprint)
		values.Add("sni", node.Sni)
	case "reality":
		values.Add("security", node.Security)
		values.Add("pbk", node.PublicKey)
		values.Add("fp", node.Fingerprint)
		values.Add("spx", node.SpiderX)
		values.Add("sni", node.Sni)
		values.Add("sid", node.ShortId)
	}

	vlessUrl.RawQuery = values.Encode()
	vlessUrl.Fragment = node.Remarks

	return vlessUrl.String()
}

// generate  Clash vmess vless trojan
func ClashVmessVlessNew(v model.Node, uuid, host string) model.ClashProxy {
	var proxy model.ClashProxy
	switch v.NodeType {
	case "vmess":
		proxy.Type = "vmess"
		proxy.Uuid = uuid
		proxy.Alterid = strconv.FormatInt(v.Aid, 10)
		proxy.Cipher = "auto"
	case "vless":
		proxy.Type = "vless"
		proxy.Uuid = uuid
		proxy.Flow = v.VlessFlow
	case "trojan":
		proxy.Type = "trojan"
		proxy.Password = uuid
		proxy.Sni = v.Sni
	}
	if v.EnableTransfer {
		proxy.Server = v.TransferAddress
		proxy.Port = int(v.TransferPort)
	} else {
		proxy.Server = v.Address
		proxy.Port = int(v.Port)
	}
	proxy.Name = v.Remarks
	proxy.Udp = true
	proxy.Network = v.Network
	proxy.SkipCertVerify = v.AllowInsecure

	switch proxy.Network {
	case "ws":
		proxy.WsOpts.Path = v.Path
		proxy.WsOpts.Headers = make(map[string]string, 1)
		proxy.WsOpts.Headers["Host"] = host
	case "grpc":
		proxy.GrpcOpts.GrpcServiceName = "grpc"
	case "tcp":
	case "h2":
		proxy.H2Opts.Path = v.Path
		proxy.H2Opts.Host = append(proxy.H2Opts.Host, v.Host)
	}
	switch v.Security {
	case "tls":
		proxy.Tls = true
		proxy.Servername = v.Sni
		proxy.ClientFingerprint = v.Fingerprint
		proxy.Alpn = append(proxy.Alpn, v.Alpn)
	case "reality":
		proxy.Tls = true
		proxy.Servername = v.Sni
		proxy.RealityOpts.PublicKey = v.PublicKey
		proxy.RealityOpts.ShortID = v.ShortId
		proxy.ClientFingerprint = v.Fingerprint
		proxy.Alpn = append(proxy.Alpn, v.Alpn)
	}
	return proxy
}
