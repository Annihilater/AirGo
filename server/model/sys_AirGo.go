package model

import (
	uuid "github.com/satori/go.uuid"
)

type AGNodeStatus struct {
	ID int64 `json:"id"`
	AGNodeStatusItem
}

// NodeStatus Node status
type AGNodeStatusItem struct {
	CPU    float64
	Mem    float64
	Disk   float64
	Uptime uint64
}

//type AGNodeInfo struct {
//	NodeType          string // Must be V2ray, Trojan, and Shadowsocks
//	NodeID            int
//	Port              uint32
//	SpeedLimit        uint64 // Bps
//	AlterID           uint16
//	TransportProtocol string
//	FakeType          string
//	Host              string
//	Path              string
//	EnableTLS         bool
//	EnableVless       bool
//	VlessFlow         string
//	CypherMethod      string
//	ServerKey         string
//	ServiceName       string
//	Header            json.RawMessage
//	//NameServerConfig  []*conf.NameServerConfig
//	EnableREALITY bool
//	REALITYConfig *AGREALITYConfig
//}

type AGNodeInfo struct {
	ID             int64  `json:"id"`
	NodeSpeedlimit int64  `json:"node_speedlimit"` //节点限速/Mbps
	TrafficRate    int64  `json:"traffic_rate"`    //倍率
	NodeType       string `json:"node_type"`       //节点类型 vless,vmess,trojan
	Remarks        string `json:"remarks"`         //别名
	Address        string `json:"address"`         //地址
	Port           int64  `json:"port"`            //端口

	//vmess参数
	Scy string `json:"scy"` //加密方式 auto,none,chacha20-poly1305,aes-128-gcm,zero，vless选择none，否则v2rayng无法启动
	Aid int64  `json:"aid"` //额外ID
	//vless参数
	VlessFlow string `json:"flow"` //流控 none,xtls-rprx-vision,xtls-rprx-vision-udp443

	//传输参数
	Network     string `json:"network"`      //传输协议 tcp,kcp,ws,h2,quic,grpc
	Type        string `json:"type"`         //伪装类型 ws,h2：无    tcp,kcp：none，http    quic：none，srtp，utp，wechat-video，dtls，wireguard
	Host        string `json:"host"`         //伪装域名
	Path        string `json:"path"`         //path
	GrpcMode    string `json:"mode"`         //grpc传输模式 gun，multi
	ServiceName string `json:"service_name"` //

	//传输层安全
	Security    string `json:"security"` //传输层安全类型 none,tls,reality
	Sni         string `json:"sni"`      //
	Fingerprint string `json:"fp"`       //
	Alpn        string `json:"alpn"`     //
	Dest        string `json:"dest"`
	PrivateKey  string `json:"private_key"`
	PublicKey   string `json:"pbk"`
	ShortId     string `json:"sid"`
	SpiderX     string `json:"spx"`
}

//	type UserInfo struct {
//		UID         int
//		Email       string
//		UUID        string
//		Passwd      string
//		Port        uint32
//		AlterID     uint16
//		Method      string
//		SpeedLimit  uint64 // Bps
//		DeviceLimit int
//	}
type AGUserInfo struct {
	ID       int64     `json:"id"`
	UUID     uuid.UUID `json:"uuid"`
	UserName string    `json:"user_name"`
}

type AGUserTraffic struct {
	ID          int64               `json:"id"`
	UserTraffic []AGUserTrafficItem `json:"user_traffic"`
}

type AGUserTrafficItem struct {
	UID      int64
	Email    string
	Upload   int64
	Download int64
}

type AGREALITYConfig struct {
	Dest             string
	ProxyProtocolVer uint64
	ServerNames      []string
	PrivateKey       string
	MinClientVer     string
	MaxClientVer     string
	MaxTimeDiff      uint64
	ShortIds         []string
}

type AGREALITYx25519 struct {
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}
