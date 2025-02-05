package initialize

import (
	"AirGo/global"
	"AirGo/model"
	utils "AirGo/utils/encrypt_plugin"
	"errors"
	"gorm.io/driver/sqlite"
	//"go-admin/initialize"
	//github.com/satori/go.uuid
	gormadapter "github.com/casbin/gorm-adapter/v3"
	uuid "github.com/satori/go.uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Gorm 初始化数据库并产生数据库全局变量
// Author SliverHorn
func Gorm() *gorm.DB {

	switch global.Config.SystemParams.DbType {
	case "mysql":
		return GormMysql()
	case "sqlite":
		return GormSqlite()
	default:
		return GormMysql()
	}
}

// 初始化sqlite数据库
func GormSqlite() *gorm.DB {

	if db, err := gorm.Open(sqlite.Open(global.Config.Sqlite.Path), &gorm.Config{
		SkipDefaultTransaction: true, //关闭事务，将获得大约 30%+ 性能提升
		NamingStrategy: schema.NamingStrategy{
			//TablePrefix: "gormv2_",
			SingularTable: true, //单数表名
		},
	}); err != nil {
		global.Logrus.Error("gorm.Open error:", err)
		panic(err)
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(int(global.Config.Mysql.MaxIdleConns))
		sqlDB.SetMaxOpenConns(int(global.Config.Mysql.MaxOpenConns))
		return db
	}
}

// 初始化Mysql数据库
func GormMysql() *gorm.DB {
	mysqlConfig := mysql.Config{
		DSN:                       global.Config.Mysql.Username + ":" + global.Config.Mysql.Password + "@tcp(" + global.Config.Mysql.Address + ":" + global.Config.Mysql.Port + ")/" + global.Config.Mysql.Dbname + "?" + global.Config.Mysql.Config,
		DefaultStringSize:         191, // string 类型字段的默认长度
		SkipInitializeWithVersion: false,
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{
		SkipDefaultTransaction: true, //关闭事务，将获得大约 30%+ 性能提升
		NamingStrategy: schema.NamingStrategy{
			//TablePrefix: "gormv2_",
			SingularTable: true, //单数表名
		},
	}); err != nil {
		global.Logrus.Error("gorm.Open error:", err)
		panic(err)
	} else {
		db.InstanceSet("gorm:table_options", "ENGINE="+global.Config.Mysql.Engine)
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(int(global.Config.Mysql.MaxIdleConns))
		sqlDB.SetMaxOpenConns(int(global.Config.Mysql.MaxOpenConns))
		return db
	}
}

// RegisterTables 注册数据库表专用
func RegisterTables() {
	err := global.DB.AutoMigrate(
		// 用户表
		model.User{},
		//动态路由表
		model.DynamicRoute{},
		//角色表
		model.Role{},
		//商品
		model.Goods{},
		//订单
		model.Orders{},
		//流量统计表
		model.TrafficLog{},
		//主题
		model.Theme{},
		//系统设置参数
		model.Server{},
		//图库
		model.Gallery{},
		//文章
		model.Article{},
		//折扣
		model.Coupon{},
		//isp
		model.ISP{},
		//共享节点
		model.NodeShared{},
	)
	if err != nil {
		//os.Exit(0)
		global.Logrus.Error("table AutoMigrate error:", err.Error())
		return
	}
	global.Logrus.Info("table AutoMigrate success")
}

// 导入数据
func InsertInto(db *gorm.DB) error {
	uuid1 := uuid.NewV4()
	uuid2 := uuid.NewV4()
	sysUserData := []model.User{
		{
			UUID:           uuid1,
			UserName:       global.Config.SystemParams.AdminEmail,
			Password:       utils.BcryptEncode(global.Config.SystemParams.AdminPassword),
			NickName:       "admin",
			InvitationCode: utils.RandomString(8),
		},
		{
			UUID:           uuid2,
			UserName:       "测试1@oicq.com",
			Password:       utils.BcryptEncode("123456"),
			NickName:       "测试1",
			InvitationCode: utils.RandomString(8),
		},
	}
	if err := db.Create(&sysUserData).Error; err != nil {
		return errors.New("db.Create(&userData) Error")
	}
	//插入sys_dynamic-router_data表
	DynamicRouteData := []model.DynamicRoute{
		{ParentID: 0, Path: "/admin", Name: "admin", Component: "/layout/routerView/parent.vue", Meta: model.Meta{Title: "超级管理员", Icon: "iconfont icon-shouye_dongtaihui"}},     //id==1
		{ParentID: 1, Path: "/admin/menu", Name: "adminMenu", Component: "/admin/menu/index.vue", Meta: model.Meta{Title: "菜单管理", Icon: "iconfont icon-caidan"}},                //id==2
		{ParentID: 1, Path: "/admin/role", Name: "adminRole", Component: "/admin/role/index.vue", Meta: model.Meta{Title: "角色管理", Icon: "iconfont icon-icon-"}},                 //id==3
		{ParentID: 1, Path: "/admin/user", Name: "adminUser", Component: "/admin/user/index.vue", Meta: model.Meta{Title: "用户管理", Icon: "iconfont icon-gerenzhongxin"}},         //id==4
		{ParentID: 1, Path: "/admin/order", Name: "adminOrder", Component: "/admin/order/index.vue", Meta: model.Meta{Title: "订单管理", Icon: "iconfont icon--chaifenhang"}},       //id==5
		{ParentID: 1, Path: "/admin/node", Name: "adminNode", Component: "/admin/node/index.vue", Meta: model.Meta{Title: "节点管理", Icon: "iconfont icon-shuxingtu"}},             //id==6
		{ParentID: 1, Path: "/admin/shop", Name: "adminShop", Component: "/admin/shop/index.vue", Meta: model.Meta{Title: "商品管理", Icon: "iconfont icon-zhongduancanshuchaxun"}}, //id==7
		{ParentID: 1, Path: "/admin/system", Name: "system", Component: "/admin/system/index.vue", Meta: model.Meta{Title: "系统设置", Icon: "iconfont icon-xitongshezhi"}},         //id==8
		{ParentID: 1, Path: "/admin/article", Name: "article", Component: "/admin/article/index.vue", Meta: model.Meta{Title: "文章设置", Icon: "iconfont icon-huanjingxingqiu"}},   //id==9
		{ParentID: 1, Path: "/admin/coupon", Name: "coupon", Component: "/admin/coupon/index.vue", Meta: model.Meta{Title: "折扣码管理", Icon: "ele-ShoppingBag"}},                   //id==10

		{ParentID: 0, Path: "/home", Name: "home", Component: "/home/index.vue", Meta: model.Meta{Title: "首页", Icon: "iconfont icon-shouye"}},                           //11
		{ParentID: 0, Path: "/shop", Name: "shop", Component: "/shop/index.vue", Meta: model.Meta{Title: "商店", Icon: "iconfont icon-zidingyibuju"}},                     //12
		{ParentID: 0, Path: "/myOrder", Name: "myOrder", Component: "/myOrder/index.vue", Meta: model.Meta{Title: "我的订单", Icon: "iconfont icon--chaifenhang"}},          //13
		{ParentID: 0, Path: "/personal", Name: "personal", Component: "/personal/index.vue", Meta: model.Meta{Title: "个人信息", Icon: "iconfont icon-gerenzhongxin"}},      //14
		{ParentID: 0, Path: "/serverStatus", Name: "serverStatus", Component: "/serverStatus/index.vue", Meta: model.Meta{Title: "节点状态", Icon: "iconfont icon-putong"}}, //15
		{ParentID: 0, Path: "/gallery", Name: "gallery", Component: "/gallery/index.vue", Meta: model.Meta{Title: "无限图库", Icon: "iconfont icon-step"}},                  //16
		{ParentID: 0, Path: "/income", Name: "income", Component: "/income/index.vue", Meta: model.Meta{Title: "营收概览", Icon: "iconfont icon-xingqiu"}},                  //17
		{ParentID: 0, Path: "/isp", Name: "isp", Component: "/isp/index.vue", Meta: model.Meta{Title: "套餐监控", Icon: "iconfont icon-xingqiu"}},                           //18
	}
	if err := db.Create(&DynamicRouteData).Error; err != nil {
		return errors.New("sys_dynamic-router_data表数据初始化失败!")
	}
	//插入user_role表
	sysRoleData := []model.Role{
		{ID: 1, RoleName: "admin", Description: "超级管理员"},
		{ID: 2, RoleName: "客服", Description: "客服"},
		{ID: 3, RoleName: "合作伙伴", Description: "合作伙伴"},
		{ID: 4, RoleName: "普通用户", Description: "普通用户"},
	}
	if err := db.Create(&sysRoleData).Error; err != nil {
		return errors.New("user_role表数据初始化失败!")
	}
	//插入user_and_role表
	userAndRoleData := []model.UserAndRole{
		{UserID: 1, RoleID: 1},
		{UserID: 1, RoleID: 2},
		{UserID: 2, RoleID: 2},
	}
	if err := db.Create(&userAndRoleData).Error; err != nil {
		return errors.New("user_and_role_data表数据初始化失败!")
	}
	//插入role_and_menu
	roleAndMenuData := []model.RoleAndMenu{
		//管理员的权限
		{RoleID: 1, DynamicRouteID: 1},  //超级管理员
		{RoleID: 1, DynamicRouteID: 2},  //菜单管理
		{RoleID: 1, DynamicRouteID: 3},  //角色管理
		{RoleID: 1, DynamicRouteID: 4},  //用户管理
		{RoleID: 1, DynamicRouteID: 5},  //订单管理
		{RoleID: 1, DynamicRouteID: 6},  //节点管理
		{RoleID: 1, DynamicRouteID: 7},  //商品管理
		{RoleID: 1, DynamicRouteID: 8},  //系统设置
		{RoleID: 1, DynamicRouteID: 9},  //文章设置
		{RoleID: 1, DynamicRouteID: 10}, //折扣码管理
		{RoleID: 1, DynamicRouteID: 11},
		{RoleID: 1, DynamicRouteID: 12},
		{RoleID: 1, DynamicRouteID: 13},
		{RoleID: 1, DynamicRouteID: 14},
		{RoleID: 1, DynamicRouteID: 15},
		{RoleID: 1, DynamicRouteID: 16},
		{RoleID: 1, DynamicRouteID: 17}, //营收概览
		{RoleID: 1, DynamicRouteID: 18},

		//客服的权限
		{RoleID: 2, DynamicRouteID: 1},
		{RoleID: 2, DynamicRouteID: 4},
		{RoleID: 2, DynamicRouteID: 5},
		{RoleID: 2, DynamicRouteID: 7},
		{RoleID: 2, DynamicRouteID: 9},
		{RoleID: 2, DynamicRouteID: 10},

		{RoleID: 2, DynamicRouteID: 11},
		{RoleID: 2, DynamicRouteID: 12},
		{RoleID: 2, DynamicRouteID: 13},
		{RoleID: 2, DynamicRouteID: 14},
		{RoleID: 2, DynamicRouteID: 15},
		{RoleID: 2, DynamicRouteID: 16},
		{RoleID: 2, DynamicRouteID: 17},
		{RoleID: 2, DynamicRouteID: 18},

		//合作伙伴的权限
		{RoleID: 3, DynamicRouteID: 11},
		{RoleID: 3, DynamicRouteID: 12},
		{RoleID: 3, DynamicRouteID: 13},
		{RoleID: 3, DynamicRouteID: 14},
		{RoleID: 3, DynamicRouteID: 15},
		{RoleID: 3, DynamicRouteID: 16},
		{RoleID: 3, DynamicRouteID: 17},
		{RoleID: 3, DynamicRouteID: 18},

		//普通用户的权限
		{RoleID: 4, DynamicRouteID: 11},
		{RoleID: 4, DynamicRouteID: 12},
		{RoleID: 4, DynamicRouteID: 13},
		{RoleID: 4, DynamicRouteID: 14},
		{RoleID: 4, DynamicRouteID: 15},
		{RoleID: 4, DynamicRouteID: 16},
		{RoleID: 4, DynamicRouteID: 18},
	}
	if err := global.DB.Create(&roleAndMenuData).Error; err != nil {
		return errors.New("role_and_menu表数据初始化失败!")
	}
	//插入货物 goods
	goodsData := []model.Goods{
		{Subject: "10G|30天", TotalBandwidth: 10, ExpirationDate: 30, TotalAmount: "0.01", Des: text2},
		{Subject: "20G|180天", TotalBandwidth: 20, ExpirationDate: 180, TotalAmount: "0", Des: text2},
	}
	if err := global.DB.Create(&goodsData).Error; err != nil {
		return errors.New("goods表数据初始化失败!")
	}
	//插入node
	nodeData := []model.Node{
		{Remarks: "测试节点1", Address: "www.10010.com", Path: "/path", Port: 5566, NodeType: "vless", Enabled: true},
		{Remarks: "测试节点2", Address: "www.10086.com", Path: "/path", Port: 5566, NodeType: "vless", Enabled: true},
	}
	if err := global.DB.Create(&nodeData).Error; err != nil {
		return errors.New("node表数据初始化失败!")
	}
	//插入goods_and_nodes
	goodsAndNodesData := []model.GoodsAndNodes{
		{GoodsID: 1, NodeID: 1},
		{GoodsID: 1, NodeID: 2},
		{GoodsID: 2, NodeID: 1},
		{GoodsID: 2, NodeID: 2},
	}
	if err := global.DB.Create(&goodsAndNodesData).Error; err != nil {
		return errors.New("goods_and_nodes表数据初始化失败!")
	}
	// 插入casbin_rule
	casbinRuleData := []gormadapter.CasbinRule{
		//{Ptype: "p", V0: "1", V1: "/public/getThemeConfig", V2: "GET"},

		{Ptype: "p", V0: "1", V1: apiPrefix + "/user/register", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/user/login", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/user/getSub", V2: "GET"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/user/changeSubHost", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/user/getUserInfo", V2: "GET"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/user/changeUserPassword", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/user/resetUserPassword", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/user/resetSub", V2: "GET"},

		{Ptype: "p", V0: "1", V1: apiPrefix + "/user/getUserList", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/user/newUser", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/user/updateUser", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/user/deleteUser", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/user/findUser", V2: "POST"},

		{Ptype: "p", V0: "1", V1: apiPrefix + "/role/getRoleList", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/role/modifyRoleInfo", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/role/addRole", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/role/delRole", V2: "DELETE"},

		{Ptype: "p", V0: "1", V1: apiPrefix + "/menu/getAllRouteList", V2: "GET"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/menu/getAllRouteTree", V2: "GET"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/menu/newDynamicRoute", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/menu/delDynamicRoute", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/menu/updateDynamicRoute", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/menu/findDynamicRoute", V2: "POST"},

		{Ptype: "p", V0: "1", V1: apiPrefix + "/menu/getRouteList", V2: "GET"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/menu/getRouteTree", V2: "GET"},

		{Ptype: "p", V0: "1", V1: apiPrefix + "/shop/preCreatePay", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/shop/purchase", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/shop/getAllEnabledGoods", V2: "GET"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/shop/getAllGoods", V2: "GET"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/shop/findGoods", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/shop/newGoods", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/shop/deleteGoods", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/shop/updateGoods", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/shop/alipayNotify", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/shop/goodsSort", V2: "POST"},

		{Ptype: "p", V0: "1", V1: apiPrefix + "/node/getAllNode", V2: "GET"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/node/newNode", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/node/deleteNode", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/node/updateNode", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/node/getTraffic", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/node/nodeSort", V2: "POST"},

		{Ptype: "p", V0: "1", V1: apiPrefix + "/node/newNodeShared", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/node/getNodeSharedList", V2: "GET"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/node/deleteNodeShared", V2: "POST"},

		{Ptype: "p", V0: "1", V1: apiPrefix + "/mod_mu/nodes/:nodeID/info", V2: "GET"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/mod_mu/users", V2: "GET"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/mod_mu/users/traffic", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/mod_mu/users/aliveip", V2: "POST"},

		{Ptype: "p", V0: "1", V1: apiPrefix + "/casbin/getPolicyByRoleIds", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/casbin/updateCasbinPolicy", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/casbin/updateCasbinPolicyNew", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/casbin/getAllPolicy", V2: "GET"},

		{Ptype: "p", V0: "1", V1: apiPrefix + "/order/getOrderInfo", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/order/getAllOrder", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/order/getOrderByUserID", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/order/completedOrder", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/order/getMonthOrderStatistics", V2: "POST"},

		{Ptype: "p", V0: "1", V1: apiPrefix + "/system/updateThemeConfig", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/system/getSetting", V2: "GET"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/system/updateSetting", V2: "POST"},

		{Ptype: "p", V0: "1", V1: apiPrefix + "/upload/newPictureUrl", V2: "GET"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/upload/getPictureList", V2: "POST"},

		{Ptype: "p", V0: "1", V1: apiPrefix + "/websocket/msg", V2: "GET"},

		{Ptype: "p", V0: "1", V1: apiPrefix + "/article/newArticle", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/article/deleteArticle", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/article/updaterticle", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/article/getArticle", V2: "POST"},

		{Ptype: "p", V0: "1", V1: apiPrefix + "/report/getDB", V2: "GET"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/report/getTables", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/report/getColumn", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/report/reportSubmit", V2: "POST"},

		{Ptype: "p", V0: "1", V1: apiPrefix + "/coupon/newCoupon", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/coupon/deleteCoupon", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/coupon/updateCoupon", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/coupon/getCoupon", V2: "POST"},

		{Ptype: "p", V0: "1", V1: apiPrefix + "/isp/sendCode", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/isp/ispLogin", V2: "POST"},
		{Ptype: "p", V0: "1", V1: apiPrefix + "/isp/getMonitorByUserID", V2: "POST"},
		//{Ptype: "p", V0: "1", V1: apiPrefix+"/isp/queryPackage", V2: "POST"},

		//普通用户
		{Ptype: "p", V0: "2", V1: apiPrefix + "/user/login", V2: "POST"},
		{Ptype: "p", V0: "2", V1: apiPrefix + "/user/getSub", V2: "GET"},
		{Ptype: "p", V0: "2", V1: apiPrefix + "/user/changeUserPassword", V2: "POST"},
		{Ptype: "p", V0: "2", V1: apiPrefix + "/user/resetUserPassword", V2: "POST"},
		{Ptype: "p", V0: "2", V1: apiPrefix + "/user/getUserInfo", V2: "GET"},
		{Ptype: "p", V0: "2", V1: apiPrefix + "/user/resetSub", V2: "GET"},
		{Ptype: "p", V0: "2", V1: apiPrefix + "/user/changeSubHost", V2: "POST"},

		{Ptype: "p", V0: "2", V1: apiPrefix + "/menu/getRouteList", V2: "GET"},
		{Ptype: "p", V0: "2", V1: apiPrefix + "/menu/getRouteTree", V2: "GET"},

		{Ptype: "p", V0: "2", V1: apiPrefix + "/order/getOrderInfo", V2: "POST"},
		{Ptype: "p", V0: "2", V1: apiPrefix + "/order/getOrderByUserID", V2: "POST"},

		{Ptype: "p", V0: "2", V1: apiPrefix + "/shop/preCreatePay", V2: "POST"},
		{Ptype: "p", V0: "2", V1: apiPrefix + "/shop/purchase", V2: "POST"},
		{Ptype: "p", V0: "2", V1: apiPrefix + "/shop/getAllEnabledGoods", V2: "GET"},
		{Ptype: "p", V0: "2", V1: apiPrefix + "/shop/findGoods", V2: "POST"},

		{Ptype: "p", V0: "2", V1: apiPrefix + "/websocket/msg", V2: "GET"},
		{Ptype: "p", V0: "2", V1: apiPrefix + "/upload/newPictureUrl", V2: "GET"},
		{Ptype: "p", V0: "2", V1: apiPrefix + "/upload/getPictureList", V2: "POST"},

		{Ptype: "p", V0: "2", V1: apiPrefix + "/article/getArticle", V2: "POST"},

		{Ptype: "p", V0: "2", V1: apiPrefix + "/isp/sendCode", V2: "POST"},
		{Ptype: "p", V0: "2", V1: apiPrefix + "/isp/ispLogin", V2: "POST"},
		{Ptype: "p", V0: "2", V1: apiPrefix + "/isp/getMonitorByUserID", V2: "POST"},
	}
	if err := global.DB.Create(&casbinRuleData).Error; err != nil {
		return errors.New("casbin_rule表数据初始化失败!")
	}
	//主题配置
	themeData := model.Theme{
		ID: 1,
	}
	if err := global.DB.Create(&themeData).Error; err != nil {
		return errors.New("theme表数据初始化失败!")
	}

	//系统设置
	settingData := model.Server{
		ID: 1,
		Email: model.Email{
			EmailContent: text1,
		},
	}
	if err := global.DB.Create(&settingData).Error; err != nil {
		return errors.New("server表数据初始化失败!")
	}
	return nil
}

// 默认邮件验证码样式
const text1 = `
<style>
.cookieCard {
  margin:auto;
  width: 300px;
  height: 200px;
  background: linear-gradient(to right,rgb(137, 104, 255),rgb(175, 152, 255));
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  justify-content: center;
  gap: 20px;
  padding: 20px;
  position: relative;
  overflow: hidden;
}
.cookieCard::before {
  width: 150px;
  height: 150px;
  content: "";
  background: linear-gradient(to right,rgb(142, 110, 255),rgb(208, 195, 255));
  position: absolute;
  z-index: 1;
  border-radius: 50%;
  right: -25%;
  top: -25%;
}
.cookieHeading {
  font-size: 1.5em;
  font-weight: 600;
  z-index: 2;
}

.cookieDescription {
  font-size: 0.9em;
  z-index: 2;
}
</style>
</head>
<body>

<div class="cookieCard">
  <p class="cookieHeading">验证码</p>
  <p class="cookieDescription">欢迎使用，请及时输入验证码</p>
  <span style="font-size:30px">emailcode</span>
</div>
</body>
`

// 商品默认描述
const text2 = `
<h3 style="color:#00BFFF">究竟什么样的终点，才配得上这一路的颠沛流离---管泽元</h3>
<h3 style="color:#DDA0DD">世界聚焦于你---管泽元</h3>
`
