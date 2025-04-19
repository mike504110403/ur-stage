package order

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"wallet_service/internal/cachedata"
	"wallet_service/internal/database"
	"wallet_service/internal/transection"

	order "gitlab.com/gogogo2712128/common_moduals/dbModel/Order"
	"gitlab.com/gogogo2712128/common_moduals/typeparam"

	htpay "wallet_service/instance/ht_pay"

	"github.com/google/uuid"
	mlog "github.com/mike504110403/goutils/log"
)

// thirdPayOrderInit : 訂單起單	- 第三方支付
func orderInit(orderReq OrderInitReq, mid int, cip string) (string, error) {
	url, trade := "", uuid.New().String()
	// 呼叫第三方起單
	thirdRes, err := thirdPayInitAPI(
		thirdPayInitAPIReq{
			ThirdPayType: orderReq.PayType,
			Amount:       orderReq.Amount,
			TradeNo:      trade,
		},
		mid,
		cip,
	)
	if err != nil {
		return url, err
	} else {
		url = thirdRes.Url
	}

	db, err := database.ORDER.DB()
	if err != nil {
		return url, err
	}

	// 三方 Id
	thirdPayIdStr, ok := cachedata.ThirdPayIdMap()[orderReq.PayType]
	if !ok {
		return url, errors.New("無對應的支付方式")
	}
	thirdPayId, err := strconv.Atoi(thirdPayIdStr)
	if err != nil {
		return url, err
	}

	// 點數商品參數
	itemIdStr, ok := cachedata.ItemIdMap()[orderReq.ItemType]
	if !ok {
		return url, errors.New("無對應的商品")
	}
	itemId, err := strconv.Atoi(itemIdStr)
	if err != nil {
		return url, err
	}
	// 訂單狀態 - 起單
	initOrder := typeparam.TypeParam{
		MainType: "order_status",
		SubType:  "init",
	}
	initState, err := initOrder.Get()
	if err != nil {
		return url, err
	}

	insertStr := `
		INSERT INTO OrderPurchase
		(
			member_id,
			trade_no,
			third_pay_id,
			third_pay_trade_no,
			amount,
			order_status,
			item_type,
			item_count
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	if _, err = db.Exec(
		insertStr,
		mid,
		trade,
		thirdPayId,
		thirdRes.ThirdPayTradeNo,
		orderReq.Amount,
		initState,
		itemId,
		orderReq.ItemCount,
	); err != nil {
		return url, err
	}
	return url, nil
}

// 呼叫第三方起單
func thirdPayInitAPI(req thirdPayInitAPIReq, mid int, cip string) (thirdPayInitAPIRes, error) {
	resp := thirdPayInitAPIRes{}
	resp.TradeNo = req.TradeNo

	switch req.ThirdPayType {
	case "ht_trc20", "ht_erc20":
		var extra string
		if req.ThirdPayType == "ht_trc20" {
			extra = `{"channel":"TRC20"}`
		} else {
			extra = `{"channel":"ERC20"}`
		}
		usdtReq := htpay.PlaceOrderReq{
			Amount:     req.Amount,
			MhtOrderNo: req.TradeNo,
			MhtUserId:  strconv.Itoa(mid),
			ClientIp:   cip,
			Extra:      extra,
		}
		usdtRes, err := htpay.PlaceOrder(usdtReq)
		if err != nil {
			return resp, err
		}
		resp.ThirdPayTradeNo = usdtRes.Result.PfOrderNo
		resp.Url = usdtRes.Result.PayUrl
	default:
		return thirdPayInitAPIRes{}, errors.New("無對應支付方式")
	}
	return resp, nil
}

// paymentOrder 收款
func paymentOrder(req PayCallbackReq) error {
	mlog.Info(fmt.Sprintf("訂單收款: %s", req.MhtOrderNo))
	orderDb, err := database.ORDER.DB()
	if err != nil {
		mlog.Error(err.Error())
		return err
	}
	order, err := getOrder(orderDb, req.MhtOrderNo)
	if err != nil {
		mlog.Error(err.Error())
		return err
	}
	if err := checkOrder(order, req); err != nil {
		mlog.Error(err.Error())
		return err
	}
	// 轉點來源類型 - 訂單
	transOrderType := typeparam.TypeParam{
		MainType: "tx_src_type",
		SubType:  "order",
	}
	orderType, err := transOrderType.Get()
	if err != nil {
		return err
	}
	// 轉點單
	waleetDb, err := database.WALLET.DB()
	if err != nil {
		return err
	}
	if err = transectionSet(waleetDb, TransectionSetDto{
		MemberId:           order.MemberId,
		Amount:             order.Amount,
		TransectionSrcType: orderType,
		TransectionRelate:  order.ID,
	}); err != nil {
		return err
	}
	return nil
}

// getOrder 取得訂單
func getOrder(db *sql.DB, tradeNo string) (order.Order, error) {
	order := order.Order{}
	orderDone := typeparam.TypeParam{
		MainType: "order_status",
		SubType:  "done",
	}
	doneState, err := orderDone.Get()
	if err != nil {
		return order, err
	}

	queryStr := `
		SELECT * FROM OrderPurchase
		WHERE trade_no = ? AND order_status != ?
		FOR UPDATE
		SKIP LOCKED
	`
	err = db.QueryRow(queryStr, tradeNo, doneState).Scan(
		&order.ID,
		&order.MemberId,
		&order.TradeNo,
		&order.ThirdPayId,
		&order.ThirdPayTradeNo,
		&order.Amount,
		&order.OrderStatus,
		&order.ItmeType,
		&order.ItemCount,
		&order.Expiration,
		&order.CreatedAt,
		&order.UpdatedAt,
		&order.OrderDoneAt,
	)
	if err != nil {
		return order, err
	}
	return order, nil
}

// checkOrder 檢查訂單
func checkOrder(order order.Order, req PayCallbackReq) error {
	if req.Status != int(PayCallbackReqStatusSuccess) {
		return errors.New("訂單狀態不符")
	}
	if order.Amount != req.PaidAmount {
		return errors.New("訂單金額不符")
	}
	return nil
}

// transectionSet 設定交易
func transectionSet(db *sql.DB, trans TransectionSetDto) error {
	err := transection.SetPointTranection(db, transection.TransectionSetReq{
		MemberId:           trans.MemberId,
		Amount:             trans.Amount,
		TransectionSrcType: trans.TransectionSrcType,
		TransectionRelate:  trans.TransectionRelate,
	})
	if err != nil {
		return err
	}
	return nil
}

func checkAuth(mid int) bool {
	// 檢查權限
	db, err := database.MEMBER.DB()
	if err != nil {
		mlog.Error(err.Error())
		return false
	}
	memberState := typeparam.TypeParam{
		MainType: "member_state",
		SubType:  "admin",
	}
	state, err := memberState.Get()
	if err != nil {
		mlog.Error(err.Error())
		return false
	}
	queryStr := `
		SELECT COUNT(*) FROM Member
		WHERE id = ? AND member_state = ?
	`
	var count int
	err = db.QueryRow(queryStr, mid, state).Scan(&count)
	if err != nil {
		mlog.Error(err.Error())
		return false
	}
	if count == 0 {
		return false
	}
	return true
}

// adminPayment 管理員轉點
func adminPayment(req AdminPaymentReq, admin int) error {
	memberDb, err := database.MEMBER.DB()
	if err != nil {
		mlog.Error(err.Error())
		return err
	}

	// 取得會員 ID
	memberState := typeparam.TypeParam{
		MainType: "member_state",
		SubType:  "enable",
	}
	state, err := memberState.Get()
	if err != nil {
		mlog.Error(err.Error())
		return err
	}
	var transMemberId int
	queryStr := `
		SELECT id 
		FROM Member
		WHERE username = ? AND member_state = ?
	`
	err = memberDb.QueryRow(queryStr, req.UserName, state).Scan(&transMemberId)
	if err != nil {
		return errors.New("無此會員")
	}

	walletDb, err := database.WALLET.DB()
	if err != nil {
		mlog.Error(err.Error())
		return err
	}

	// 轉點來源類型 - 訂單
	transOrderType := typeparam.TypeParam{
		MainType: "tx_src_type",
		SubType:  "admin",
	}
	orderType, err := transOrderType.Get()
	if err != nil {
		return err
	}

	// 寫轉點單
	return transectionSet(walletDb, TransectionSetDto{
		MemberId:           transMemberId,
		Amount:             req.Amount,
		TransectionSrcType: orderType,
		TransectionRelate:  admin,
	})
}
