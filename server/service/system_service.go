package service

import (
	"AirGo/global"
	"AirGo/model"
	"AirGo/utils/mail_plugin"
)

// 修改系统配置
func UpdateSetting(setting *model.Server) error {
	// 修改系统配置
	err := global.DB.Save(&setting).Error
	if err != nil {
		return err
	}
	//重新加载系统配置
	global.Server = *setting
	//重新加载email
	d := mail_plugin.InitEmailDialer()
	if d != nil {
		global.EmailDialer = d
	}
	return nil
}
