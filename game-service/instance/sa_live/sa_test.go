package sa_live

import (
	"fmt"
	"game_service/internal/service"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	user := "test003"
	data, err := createUserAPI(user)
	fmt.Println(data, err)
	if err != nil {
		t.Error(err)
	}
}

func TestGameAccountExist(t *testing.T) {
	user := "test004"
	data, err := GameAccountExist(user)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(data, err)
}

func TestLogin(t *testing.T) {
	user := "test003"
	data, err := loginAPI(user)
	fmt.Println(data, err)
	if err != nil {
		t.Error(err)
	}
}

func TestVerify(t *testing.T) {
	user := "test003"
	data, err := verifyUserNameAPI(user)
	fmt.Println(data, err)
	if err != nil {
		t.Error(err)
	}
}

func TestKickUser(t *testing.T) {
	user := "test003"
	data, err := kickUserAPI(user)
	fmt.Println(data, err)
	if err != nil {
		t.Error(err)
	}
}

func TestPointIn(t *testing.T) {
	user := "test003"
	amount := 100.01
	data, err := CreditBalanceDVAPI(user, amount)
	fmt.Println(data, err)
}

func TestGetMemberStatus(t *testing.T) {

	user := service.MemberGameAccountInfo{
		UserName: "test003",
	}
	data, err := GetBalance(user)
	fmt.Println(data, err)
}

func TestPoinOut(t *testing.T) {
	user := "test003"
	amount := 80.00
	data, err := DebitBalanceDVAPI(user, amount)
	fmt.Println(data, err)
}

func TestBetRecord(t *testing.T) {
	data, err := getBetRecordAPI(3 * time.Minute)
	fmt.Println(data, err)
}
