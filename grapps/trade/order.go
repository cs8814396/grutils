package trade

import (
	"database/sql"
	"encoding/json"
	"fmt"

	//"grutils/grdatabase"
	"github.com/gdgrc/grutils/grmath"
	//"lottery_admin/src/apps/user"
	"github.com/gdgrc/grutils/grapps/config"
	//"github.com/garyburd/redigo/redis"
)

type Order struct {
	Appid         string  `json:"appid" orm:"pk"`
	OrderCode     string  `json:"order_code"`
	FrontUid      int64   `json:"front_uid"`
	BenifitUid    int64   `json:"benifit_uid"`
	ProductSaleId int64   `json:"product_sale_id"`
	ProductNum    int     `json:"product_num"`
	Uid           int64   `json:"uid"`
	Money         float64 `json:"money"`
	BaseMoney     float64 `json:"base_money"`
	GiftMoney     float64 `json:"gift_money"`

	CustomData string `json:"custom_data"`
	PayType    int    `json:"pay_type"`

	CreateTime      string         `json:"create_time"`
	PayStatus       int            `json:"pay_status"`
	PayTime         string         `json:"pay_time"`
	PayMoney        float64        `json:"pay_money"`
	PayBaseMoney    float64        `json:"pay_base_money"`
	PayGiftMoney    float64        `json:"pay_gift_money" orm:"column(pay_gift_money)"`
	Status          int            `json:"status" orm:"column(status)"`
	ExtraInfoString string         `json:"-" orm:"type(json);column(extra_info)"`
	ExtraInfo       OrderExtraInfo `json:"extra_info" orm:"-"`
}

type OrderExtraInfo struct {
	DetailString   *string                  `json:"detail,omitempty"`
	DetailSliceMap []map[string]interface{} `json:"detail_slice_map"`
	//============
	PayTotalMoney     *float64 `json:"pay_total_money,omitempty"`
	PtOrderCode       *string  `json:"pt_order_code,omitempty"`
	RefundTime        *string  `json:"refund_time,omitempty"`
	RefundOrderCode   *string  `json:"refund_order_code,omitempty"`
	RefundMsg         *string  `json:"refund_msg,omitempty"`
	RefundSuccessTime *string  `json:"refund_success_time,omitempty"`
	ChnId             *int64   `json:"chn_id,omitempty"`
	IsHidden          *int     `json:"is_hidden,omitempty"`
}

const (
	ORDER_CODE_PREFIX = "P"
	BASE_MONEY        = 0x01
	GIFT_MONEY        = 0x02
	WECHAT_PAY        = 0x04
	ALIPAY            = 0x08
	ELIFE_PAY         = 0x10
	MIXED_WECHAT_PAY  = 18
	MIXED_ALIPAY_PAY  = 19

	ORDER_STATUS_NOPAY            = 0
	ORDER_STATUS_NOTIFYING        = 1
	ORDER_STATUS_RENOTIFYING      = 2
	ORDER_STATUS_SUCCESS          = 3
	ORDER_STATUS_FAILURE          = 4
	ORDER_STATUS_REFUND_SUCCESS   = 5
	ORDER_STATUS_REFUND_FAILURE   = 6
	ORDER_STATUS_REFUNDING        = 7
	ORDER_STATUS_REFUND_MONEY_ERR = 8

	ORDER_PAY_STATUS_NOPAY   = 0
	ORDER_PAY_STATUS_PAYING  = 1
	ORDER_PAY_STATUS_SUCCESS = 2
	ORDER_PAY_STATUS_FAILURE = 3

/*
   (0x01, '赠送余额'),
   (0x02, '余额'),
   (0x04, '微信'),
   (0x08, '支付宝'),
   (0x10, '工银e支付'),*/

)

func GetOrderForUpdate(tx *sql.Tx, appid string, orderCode string, payStatus int) (cOrder Order, err error) {

	err, createTime := grmath.GetCreateTimeFromOrder(ORDER_CODE_PREFIX, orderCode)
	if err != nil {
		return
	}

	sql := "SELECT "
	sql += "`appid`, `order_code`, `front_uid`, `benifit_uid`,`product_sale_id`,`product_num`,`uid`,"
	sql += "`money`, `base_money`, `gift_money`, `custom_data`, date_format(`create_time`, '%Y-%m-%d %H:%i:%s'), `pay_type`,"
	sql += "`pay_status`,date_format(`pay_time`, '%Y-%m-%d %H:%i:%s'),`pay_money`,`pay_base_money`,`pay_gift_money`,`status`,`extra_info`"
	sql += " FROM `order`"
	sql += " WHERE `appid` = ? AND `create_time`= ? AND `order_code` = ? "

	if payStatus >= 0 {
		sql += fmt.Sprintf(" AND `pay_status` = %d ", payStatus)

	}

	sql += " LIMIT 1 FOR UPDATE"

	err = tx.QueryRow(sql, appid, createTime, orderCode).Scan(&cOrder.Appid, &cOrder.OrderCode, &cOrder.FrontUid, &cOrder.BenifitUid, &cOrder.ProductSaleId, &cOrder.ProductNum, &cOrder.Uid,
		&cOrder.Money, &cOrder.BaseMoney, &cOrder.GiftMoney, &cOrder.CustomData, &cOrder.CreateTime, &cOrder.PayType,
		&cOrder.PayStatus, &cOrder.PayTime, &cOrder.PayMoney, &cOrder.PayBaseMoney, &cOrder.PayGiftMoney, &cOrder.Status, &cOrder.ExtraInfoString)

	if err != nil {
		config.DefaultLogger.Error("admin query  fail, err: %s", err)
		return
	}

	if err = json.Unmarshal([]byte(cOrder.ExtraInfoString), &cOrder.ExtraInfo); err != nil {
		return
	}
	return
}
