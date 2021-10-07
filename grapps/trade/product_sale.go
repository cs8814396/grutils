package trade

import (
	"database/sql"
	"encoding/json"

	//"errors"
	"fmt"
	//"grutils/grdatabase"
	"github.com/gdgrc/grutils/grmath"
	//"lottery_admin/src/apps/user"
	"github.com/gdgrc/grutils/grapps/config"
	//"github.com/garyburd/redigo/redis"
)

type ProductSale struct {
	Appid         string `json:"appid" orm:"pk"`
	ProductSaleId int64  `json:"product_sale_id" orm:"column(product_sale_id)"`

	ProductId       int                  `json:"-" orm:"column(product_id)`
	Type            int                  `json:"type" orm:"column(type)"`
	Status          int                  `json:"status" orm:"column(status)"`
	Money           float64              `json:"money,omitempty"`
	Count           int                  `json:"count" orm:"column(count)"`
	FreezeCount     int                  `json:"freeze_count" orm:"column(freeze_count)"`
	CreateTime      string               `json:"create_time" orm:"column(create_time)`
	SaleStartTime   string               `json:"sale_start_time" orm:"column(sale_start_time)`
	SaleEndTime     string               `json:"sale_end_time" orm:"column(sale_end_time)`
	ExtraInfo       ProductSaleExtraInfo `json:"extra_info,omitempty" orm:"-"`
	ExtraInfoString string               `json:"-" orm:"column(extra_info);type(json)"`
}
type ProductSaleExtraInfo struct {
	RewardAmount *float64 `json:"reward_amount,omitempty"`
	Name         *string  `json:"name"`
	Desc         *string  `json:"desc"`
	LogoUrl      *string  `json:"logo_url"`
	QrCodeUrl    *string  `json:"qrcode_url"`

	WxMpAppid *string `json:"wx_mp_appid"`
	WxMpPath  *string `json:"wx_mp_path"`
	///===========
	ChainReward       *int      `json:"chain_reward,omitempty"`
	ChainRewardHeight *int      `json:"chain_reward_height,omitempty"`
	ChainRewardAmount []float64 `json:"chain_reward_amount,omitempty"`

	ChainRewardInvIndex *int `json:"chain_reward_inv_index"` // inv invalid
	ChainRewardInvPr    *int `json:"chain_reward_inv_pr"`    // pr probability 1000 means 100.0%
	//==========
	//RewardLottery    *int `json:"reward_lottery,omitempty"`
	//RewardLotteryMax *int `json:"reward_lottery_max,omitempty"`
	//RewardLotteryMin *int `json:"reward_lottery_min,omitempty"`

	//==========

	EveryLife *int `json:"every_life,omitempty"`

	//======== stat
	Sales int64 `json:"sales,omitempty"`
}

const (
	PRODUCT_SALE_TYPE_NO_SPECIFIC_MONEY_MAX = 20000
	PRODUCT_SALE_TYPE_INNER_PRODUCT_MAX     = 10000
	// 1 - 10000 no secific
	PRODUCT_SALE_TYPE_RECHARGE              = 1
	PRODUCT_SALE_TYPE_REWARD_TASK           = 2
	PRODUCT_SALE_TYPE_LOTTERY_PRODUCT       = 3
	PRODUCT_SALE_TYPE_LOTTERY_TIMES_PRODUCT = 4
	// 10000- ?
	PRODUCT_SALE_TYPE_BASE_MONEY_CONSUME = 10001
	//20000 - ?
	//PRODUCT_SALE_TYPE_PHONE_CHARGE_CARD = 20001
)

const (
	ACHIEVEMENT_PRODUCT_SALE = "ACPS%d" // achievement_product_sale
)

func (ps *ProductSale) IsReward() (isFlag bool) {
	if ps.ExtraInfo.RewardAmount != nil && *ps.ExtraInfo.RewardAmount > 0.0 {
		return true
	}
	return false
}
func (ps *ProductSale) IsChainReward() (isFlag bool) {
	if ps.ExtraInfo.ChainReward != nil && *ps.ExtraInfo.ChainReward == 1 &&
		ps.ExtraInfo.ChainRewardHeight != nil && *ps.ExtraInfo.ChainRewardHeight > 0 &&
		ps.ExtraInfo.ChainRewardAmount != nil && len(ps.ExtraInfo.ChainRewardAmount) == *ps.ExtraInfo.ChainRewardHeight {
		return true

	}
	return false
}

func GetProductSaleSelectMysql(valid_flag bool, update_flag bool) (sql string) {
	sql = "SELECT `appid`, `product_sale_id`, `product_id`,`type`,`status`,`money`,`count`,`freeze_count`,date_format(`create_time`, '%Y-%m-%d %H:%i:%s'),`sale_start_time`,`sale_end_time`,`extra_info` FROM `product_sale` where `appid` = ? AND `product_sale_id`= ? "
	if valid_flag {
		nowStr := grmath.GetFormatCurTime()
		sql += fmt.Sprintf(" AND `status`=1 AND `sale_end_time` >= '%s' AND `sale_start_time` <= '%s'", nowStr, nowStr)
	}
	if update_flag {
		sql += " FOR UPDATE"
	}
	return
}

func GetProductSaleForUpdate(tx *sql.Tx, appid string, productSaleId int64, valid_flag bool) (ps ProductSale, err error) {
	/*
		if productSaleId == 1 {
			ps.ProductSaleId = productSaleId
			ps.Type = PRODUCT_SALE_TYPE_REWARD_TASK
			ps.Status = 1
			ps.Money = 0.0

			chainReward := 1
			chainRewardHeight := 2
			rewardAmount := 3.3

			ps.ExtraInfo.RewardAmount = &rewardAmount
			ps.ExtraInfo.ChainReward = &chainReward
			ps.ExtraInfo.ChainRewardHeight = &chainRewardHeight

			amountSlice := make([]float64, 0)
			amountSlice = append(amountSlice, 2.2, 1.1)
			ps.ExtraInfo.ChainRewardAmount = amountSlice

		} else {
			err = fmt.Errorf("no product_sale")
			return
		}*/

	sql := GetProductSaleSelectMysql(valid_flag, true)

	err = tx.QueryRow(sql, appid, productSaleId).Scan(&ps.Appid, &ps.ProductSaleId, &ps.ProductId, &ps.Type, &ps.Status, &ps.Money, &ps.Count, &ps.FreezeCount, &ps.CreateTime, &ps.SaleStartTime, &ps.SaleEndTime, &ps.ExtraInfoString)
	if err != nil {
		config.DefaultLogger.Error("admin query  fail, err: %s", err)
		return
	}
	if err = json.Unmarshal([]byte(ps.ExtraInfoString), &ps.ExtraInfo); err != nil {
		return
	}

	return

}
func (p *ProductSale) TxProductSaleIncrSales(tx *sql.Tx, sales int64) (err error) {

	p.ExtraInfo.Sales = p.ExtraInfo.Sales + sales

	sql := fmt.Sprintf("UPDATE `product_sale` SET `extra_info`=JSON_SET(`extra_info`,'$.sales', ?) WHERE `appid` = ? AND `product_sale_id` = ? LIMIT 1")
	//sql += "`extra_info`= JSON_SET(`extra_info`,'$.update_time', ? , $.head_url', ? ,'$.sex', ? )"

	paramSlice := make([]interface{}, 0)

	paramSlice = append(paramSlice, p.ExtraInfo.Sales, p.Appid, p.ProductSaleId)

	_, err = tx.Exec(sql, paramSlice...)
	config.DefaultLogger.Debug(" TxProductSaleIncrSales t_sql: %s,paramSlice: %+v, len: %d,err: %s", sql, paramSlice, len(paramSlice), err)
	if err != nil {
		return
	}
	return

}
func GetProductSaleFromDBorCache(appid string, productSaleId int64, valid_flag bool, slave bool) (ps ProductSale, err error) {

	sql := GetProductSaleSelectMysql(valid_flag, false)

	mspAdmin, err := config.DataAdminGet(false)
	if err != nil {
		return
	}

	err = mspAdmin.QueryRow(sql, appid, productSaleId).Scan(&ps.Appid, &ps.ProductSaleId, &ps.ProductId, &ps.Type, &ps.Status, &ps.Money, &ps.Count, &ps.FreezeCount, &ps.CreateTime, &ps.SaleStartTime, &ps.SaleEndTime, &ps.ExtraInfoString)
	if err != nil {
		config.DefaultLogger.Error("admin query  fail, err: %s", err)
		return
	}
	if err = json.Unmarshal([]byte(ps.ExtraInfoString), &ps.ExtraInfo); err != nil {
		return
	}

	return

}
func GetValidProductSaleWxMpAppidSliceFromDBorCache(appid string, iType int, beginIndex int, count int, sortType int, slave bool) (wxMpAppidSlice []string, err error) {
	if slave {
		//get from cache and return if exists
	}

	paramSlice := make([]interface{}, 0)
	paramSlice = append(paramSlice, appid)

	sql := "SELECT `appid`,`product_sale_id`, `product_id`,`type`,`status`,`money`,`count`,`freeze_count`,"
	sql += "`create_time`,`sale_start_time`,`sale_end_time`,`extra_info` FROM `product_sale` where `appid`= ? "

	nowStr := grmath.GetFormatCurTime()
	sql += fmt.Sprintf(" AND `status`=1 AND `sale_end_time` >= '%s' AND `sale_start_time` <= '%s' AND `extra_info`->>'$.wx_mp_appid' is NOT NULL ", nowStr, nowStr)
	if iType > 0 {
		sql += " AND `type` = ? "

		paramSlice = append(paramSlice, iType)
	}
	if sortType == 1 {
		sql += " ORDER BY `money`"
	} else {
		sql += " ORDER BY `product_sale_id`"
	}

	sql += fmt.Sprintf(" limit %d, %d ", beginIndex, count)

	adminConn, err := config.DataAdminGet(false)
	if err != nil {
		return
	}

	config.DefaultLogger.Debug("sql: %s, appid: %s, type: %d", sql, appid, iType)

	rows, err := adminConn.Query(sql, paramSlice...)
	if err != nil {
		config.DefaultLogger.Error("admin query  fail, err: %s", err)
		return
	}

	defer rows.Close()

	for rows.Next() {

		var ps ProductSale

		err = rows.Scan(&ps.Appid, &ps.ProductSaleId, &ps.ProductId, &ps.Type, &ps.Status, &ps.Money, &ps.Count, &ps.FreezeCount, &ps.CreateTime, &ps.SaleStartTime, &ps.SaleEndTime, &ps.ExtraInfoString)
		if err != nil {
			return
		}

		if err = json.Unmarshal([]byte(ps.ExtraInfoString), &ps.ExtraInfo); err != nil {
			return
		}

		wxMpAppidSlice = append(wxMpAppidSlice, *ps.ExtraInfo.WxMpAppid)
		//ps0Slice = append(ps0Slice, ps)

		//psSlice = append(psSlice, ps)

	}

	config.DefaultLogger.Debug("wxMpAppidSlice: %+v", wxMpAppidSlice)

	//psSlice = &ps0Slice

	return

}

func GetValidProductSaleSliceFromDBorCache(appid string, iType int, beginIndex int, count int, sortType int, slave bool) (psSlice *[]ProductSale, err error) {
	if slave {
		//get from cache and return if exists
	}

	paramSlice := make([]interface{}, 0)
	paramSlice = append(paramSlice, appid)

	sql := "SELECT `appid`,`product_sale_id`, `product_id`,`type`,`status`,`money`,`count`,`freeze_count`,"
	sql += "`create_time`,`sale_start_time`,`sale_end_time`,`extra_info` FROM `product_sale` where `appid`= ? "

	nowStr := grmath.GetFormatCurTime()
	sql += fmt.Sprintf(" AND `status`=1 AND `sale_end_time` >= '%s' AND `sale_start_time` <= '%s'", nowStr, nowStr)
	if iType > 0 {
		sql += " AND `type` = ? "

		paramSlice = append(paramSlice, iType)
	}
	if sortType == 0 {
		sql += " ORDER BY `money`"
	}

	sql += fmt.Sprintf(" limit %d, %d ", beginIndex, count)

	adminConn, err := config.DataAdminGet(false)
	if err != nil {
		return
	}

	config.DefaultLogger.Debug("sql: %s, appid: %s, type: %d", sql, appid, iType)

	rows, err := adminConn.Query(sql, paramSlice...)
	if err != nil {
		config.DefaultLogger.Error("admin query  fail, err: %s", err)
		return
	}

	ps0Slice := make([]ProductSale, 0)

	defer rows.Close()

	for rows.Next() {

		var ps ProductSale

		err = rows.Scan(&ps.Appid, &ps.ProductSaleId, &ps.ProductId, &ps.Type, &ps.Status, &ps.Money, &ps.Count, &ps.FreezeCount, &ps.CreateTime, &ps.SaleStartTime, &ps.SaleEndTime, &ps.ExtraInfoString)
		if err != nil {
			return
		}

		if err = json.Unmarshal([]byte(ps.ExtraInfoString), &ps.ExtraInfo); err != nil {
			return
		}

		ps0Slice = append(ps0Slice, ps)

		//psSlice = append(psSlice, ps)

	}

	config.DefaultLogger.Debug("ps0Slice: %+v", ps0Slice)

	psSlice = &ps0Slice

	return

}
