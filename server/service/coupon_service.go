package service

import (
	"AirGo/global"
	"AirGo/model"
	"errors"
	"gorm.io/gorm"
	"strconv"
	"time"
)

func VerifyCoupon(coupon string, userId int64) (model.Coupon, error) {
	var c model.Coupon
	err := global.DB.Where(&model.Coupon{Name: coupon}).First(&c).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return model.Coupon{}, errors.New("优惠码不存在")
		} else {
			return model.Coupon{}, errors.New("优惠码错误")
		}
	}
	if time.Now().After(c.ExpiredAt) {
		return model.Coupon{}, errors.New("优惠码已过期")
	}
	//判断使用次数
	orderArr, err := CommonSqlFind[model.Orders, string, []model.Orders](model.Orders{}, "user_id = "+strconv.FormatInt(userId, 10)+" AND coupon_id = "+strconv.FormatInt(c.ID, 10))
	if err != nil {
		return model.Coupon{}, errors.New("优惠码错误")
	}
	if int64(len(orderArr)) >= c.Limit {
		return model.Coupon{}, errors.New("优惠码次数用尽")
	}
	//返回
	return c, nil
}
