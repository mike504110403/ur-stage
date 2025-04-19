package htpay

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	mlog "github.com/mike504110403/goutils/log"
)

func TestInit() {
	mlog.Init(mlog.Config{
		EnvMode: mlog.EnvMode("dev"),
		LogType: mlog.LogType("console"),
	})
}

func TestToQueryString(t *testing.T) {
	testOrigin := PlaceOrderReq{
		Amount:     10,
		ClientIp:   "220.132.92.168",
		MhtOrderNo: uuid.New().String(),
		MhtUserId:  "1",
		OpmhtId:    "testmerchant002",
		Random:     uuid.New().String(),
		ReturnUrl:  "https://www.google.es/",
	}

	result, err := ToQueryString(testOrigin)
	if err != nil {
		t.Errorf("ToQueryString() error = %v", err)
		return
	}
	t.Log(result)
}

func TestPlaceOrder(t *testing.T) {
	TestInit()
	r, err := PlaceOrder(PlaceOrderReq{
		Amount:     10,
		ClientIp:   "220.132.92.168",
		MhtOrderNo: uuid.New().String(),
		MhtUserId:  "1",
		OpmhtId:    "testmerchant002",
		Random:     uuid.New().String(),
		ReturnUrl:  "https://www.google.es/",
		Extra:      `{"channel":"TRC20"}`,
	})
	if err != nil {
		t.Errorf("PlaceOrder() error = %v", err)
		return
	} else {
		mlog.Info(fmt.Sprintf("url: %v", r.Result.PayUrl))
	}
}

func TestGetInfo(t *testing.T) {
	TestInit()
	_, err := GetInfo(PayInfoReq{
		MhtOrderNo: "e223ef05-34a3-45ed-b5f4-cca5019c4072",
	})
	if err != nil {
		t.Errorf("GetInfo() error = %v", err)
		return
	}
}
