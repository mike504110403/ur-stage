package apollo

import (
	// "net/http"
	// "net/http/httptest"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	member "gitlab.com/gogogo2712128/common_moduals/dbModel/Member"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// TODO: 這裡的測試是一個範例，請依照實際情況修改
// 內有兩種 打外部的可以直接拿來測對外API，
// 打mock的可以拿來測 對內API或是程式邏輯(這邊是後者)

func TestGetAccountAPI(t *testing.T) {
	// // 外部API測試
	// httpmock.Activate()
	// defer httpmock.DeactivateAndReset()

	// url := Httpdomain + "/user/exist"
	// // Register a mock responder
	// httpmock.RegisterResponder(
	// 	"POST",
	// 	url,
	// 	httpmock.NewBytesResponder(200, []byte(`{"result": "Y"}`)),
	// )

	// data, err := getAccountExistAPI(ExistReq{
	// 	Agid:     1,
	// 	Username: "test",
	// })
	// assert.NoError(t, err)
	// assert.Equal(t, "Y", data.Result)

	// httpmock 測試
	mockResponse := ExistRes{
		Result: "Y",
	}
	mockResponseBytes, _ := json.Marshal(mockResponse)
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(mockResponseBytes)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	data, err := getAccountExistAPI(
		ExistReq{
			Agid:     1,
			Username: "test",
		},
		server.URL,
	)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// 逐一檢查每個字段是否正確
	if data.Result != mockResponse.Result {
		t.Errorf("Expected Status %v, got %v", mockResponse.Result, data)
	}
}

// TODO: 這裡的測試是一個範例，請依照實際情況修改
// 測DB的部分 可以透過 sqlmock.NewRows() 去定義假DB內的資料
// 透過 mock.ExpectQuery() 去模擬查詢
// func TestGetAccountExistByUserName(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	assert.NoError(t, err)
// 	defer db.Close()

// 	// 模擬查詢結果
// 	rows := sqlmock.NewRows([]string{"username"}).AddRow("testuser")
// 	mock.ExpectQuery("SELECT username FROM MEMBER WHERE member_id = ?").
// 		WithArgs("test").
// 		WillReturnRows(rows)

// 	// 呼叫被測試的函式
// 	username, err := GameAccountExist(db, "test")
// 	assert.NoError(t, err)
// 	assert.Equal(t, "testuser", username)

// 	// 檢查所有預期的操作是否都被執行
// 	err = mock.ExpectationsWereMet()
// 	assert.NoError(t, err)
// }

func TestGetAccountRegisterAPI(t *testing.T) {
	// // 外部API測試
	// httpmock.Activate()
	// defer httpmock.DeactivateAndReset()

	// url := Httpdomain + "/user/register"
	// // Register a mock responder
	// httpmock.RegisterResponder(
	// 	"POST",
	// 	url,
	// 	httpmock.NewBytesResponder(200, []byte(`{"status": 1,"error": ""}`)),
	// )

	// data, err := getAccountRegisterAPI(RegisterReq{
	// 	Agid:     1,
	// 	Username: "test",
	// 	Password: "test",
	// 	Nickname: "test",
	// })
	// assert.NoError(t, err)
	// assert.Equal(t, 1, data.Status)
	// assert.Equal(t, "", data.Error)

	// httpmock 測試
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": 1,"error": ""}`))
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	data, err := getAccountRegisterAPI(
		member.Member{
			Username: "test",
			Password: "test",
			//NickName: "test",
		},
		server.URL,
	)
	assert.NoError(t, err)
	assert.Equal(t, 1, data.Status)
	assert.Equal(t, "", data.Error)
}

func TestRegisterUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	mock.ExpectBegin()
	tx, err := db.Begin()
	assert.NoError(t, err)

	defer db.Close()

	// 模擬插入新帳號
	mock.ExpectExec("INSERT INTO Member_Game_Account \\(username, password, nickname\\) VALUES \\(\\?, \\?, \\?\\)").
		WithArgs("apollo", sqlmock.AnyArg(), "nick").
		WillReturnResult(sqlmock.NewResult(1, 1))

	member := MemberGameAccount{
		MemberId: 1,
		Username: "apollo",
		Password: "password",
		//Nickname: "nick",
	}
	mock.ExpectCommit()
	err = registerUser(tx, member)
	tx.Commit()
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
	assert.Equal(t, "apollo", member.Username)
	assert.Equal(t, "nick", member.NickName)
	assert.NotEmpty(t, member.Password) // 確保密碼已生成
}

func TestGetDepositAPI(t *testing.T) {
	// // 外部API測試
	// httpmock.Activate()
	// defer httpmock.DeactivateAndReset()

	// url := Httpdomain + "/user/deposit"
	// // Register a mock responder
	// httpmock.RegisterResponder(
	// 	"POST",
	// 	url,
	// 	httpmock.NewBytesResponder(200, []byte(`{"status": 1,"error": ""}`)),
	// )

	// data, err := getDepositAPI(DepositReq{
	// 	Agid: 1,
	// 	User: "test",
	// 	Amount: 100,
	// })
	// assert.NoError(t, err)
	// assert.Equal(t, 1, data.Status)
	// assert.Equal(t, "", data.Error)

	mockResponse := DepositRes{
		Status: 1,
		Error:  "",
		Result: DepositResult{
			Refno:     "test123",
			PaymentId: "test123456",
			Balance:   100.00,
			Amt:       100.00,
		},
	}
	// httpmock 測試
	// 透過物件轉換成 JSON byte
	mockResponseBytes, _ := json.Marshal(mockResponse)
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(mockResponseBytes)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	data, err := getDepositAPI(
		DepositReq{
			Agid:     1,
			Username: "test",
			Gametype: "pe",
			Amt:      100.00,
			Refno:    "test123",
			Type:     "IN",
		},
		server.URL,
	)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// 逐一檢查每個字段是否正確
	if data.Status != mockResponse.Status {
		t.Errorf("Expected Status %v, got %v", mockResponse.Status, data.Status)
	}
	if data.Error != mockResponse.Error {
		t.Errorf("Expected Error %v, got %v", mockResponse.Error, data.Error)
	}

	// 比較 Result 的 DepositResult
	expected := mockResponse.Result
	actual := data.Result

	if expected != actual {
		t.Errorf("Expected DepositResult %+v, got %+v", expected, actual)
	}
}

// // TODO: DBTEST
// func TestGetDepositUser(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	assert.NoError(t, err)
// 	defer db.Close()

// 	// 模擬檢查帳號是否存在
// 	rows := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0)
// 	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM MEMBER WHERE username = ?").
// 		WithArgs("test").
// 		WillReturnRows(rows)

// 	// 模擬插入新帳號
// 	mock.ExpectExec("INSERT INTO MEMBER \\(username, password, nickname\\) VALUES \\(\\?, \\?, \\?\\)").
// 		WithArgs("test", "password", "nickname").
// 		WillReturnResult(sqlmock.NewResult(1, 1))

// 	input := RegisterInput{
// 		Username: "test",
// 		Password: "password",
// 		Nickname: "nickname",
// 	}
// 	err = getAccountRegisterUser(db, input)
// 	assert.NoError(t, err)

// 	err = mock.ExpectationsWereMet()
// 	assert.NoError(t, err)
// }

func TestGetQuotaAPI(t *testing.T) {
	// // 外部API測試
	// httpmock.Activate()
	// defer httpmock.DeactivateAndReset()

	// url := Httpdomain + "/user/quota"
	// // Register a mock responder
	// httpmock.RegisterResponder(
	// 	"POST",
	// 	url,
	// 	httpmock.NewBytesResponder(200, []byte(`{"status": 1,"error": "","result": {"balance": 300}}`)),
	// )

	// data, err := getQuotaAPI(GetquotaReq{
	// Agid:     1,
	// Username: "test",
	// Gametype: "pe",
	// })
	// assert.NoError(t, err)
	// assert.Equal(t, 1, data.Status)
	// assert.Equal(t, "", data.Error)
	// assert.Equal(t, 300, data.Result.Balance)

	// httpmock 測試
	mockResponse := GetquotaRes{
		Status: 1,
		Error:  "",
		Result: GetquotaResult{
			Balance: 300.00,
		},
	}
	mockResponseBytes, _ := json.Marshal(mockResponse)
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(mockResponseBytes)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	data, err := getQuotaAPI(
		GetquotaReq{
			Agid:     1,
			Username: "test",
			Gametype: "pe",
		},
		server.URL,
	)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// 逐一檢查每個字段是否正確
	if data.Status != mockResponse.Status {
		t.Errorf("Expected Status %v, got %v", mockResponse.Status, data.Status)
	}
	if data.Error != mockResponse.Error {
		t.Errorf("Expected Error %v, got %v", mockResponse.Error, data.Error)
	}

	// 比較 Result 的 Balance
	expected := mockResponse.Result.Balance
	actual := data.Result.Balance

	if expected != actual {
		t.Errorf("Expected Result.Balance %+v, got %+v", expected, actual)
	}
}

// TODO DBTEST
// func TestGetQuotaUser(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	assert.NoError(t, err)
// 	defer db.Close()

// 	// 模擬檢查帳號是否存在
// 	rows := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0)
// 	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM MEMBER WHERE username = ?").
// 		WithArgs("test").
// 		WillReturnRows(rows)

// 	// 模擬插入新帳號
// 	mock.ExpectExec("INSERT INTO MEMBER \\(username, password, nickname\\) VALUES \\(\\?, \\?, \\?\\)").
// 		WithArgs("test", "pe").
// 		WillReturnResult(sqlmock.NewResult(1, 1))

// 	err = getQuotaUser(db, "test", "pe")
// 	assert.NoError(t, err)

// 	err = mock.ExpectationsWereMet()
// 	assert.NoError(t, err)
// }

func TestAccountLoginAPI(t *testing.T) {
	// // 外部API測試
	// httpmock.Activate()
	// defer httpmock.DeactivateAndReset()

	// url := Httpdomain + "/user/login"
	// // Register a mock responder
	// httpmock.RegisterResponder(
	// 	"POST",
	// 	url,
	// 	httpmock.NewBytesResponder(200, []byte(`{"status": 1,"error": ""}`)),
	// )

	// data, err := getAccountLoginAPI(LoginReq{
	// 	Agid:     1,
	// 	Username: "test",
	// 	Password: "test",
	// })
	// assert.NoError(t, err)
	// assert.Equal(t, 1, data.Status)
	// assert.Equal(t, "", data.Error)

	mockResponse := LoginRes{
		Status: 1,
		Error:  "",
		Result: LoginResult{
			Site:  "A",
			Url:   "http://test.com",
			Query: "test",
		},
	}
	mockResponseBytes, _ := json.Marshal(mockResponse)
	// httpmock 測試
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(mockResponseBytes)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()
	data, err := getAccountLoginAPI(
		Login{
			Username: "test",
			Password: "2313",
			Gametype: "pe",
		},
		server.URL,
	)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// 逐一檢查每個字段是否正確
	if data.Status != mockResponse.Status {
		t.Errorf("Expected Status %v, got %v", mockResponse.Status, data.Status)
	}
	if data.Error != mockResponse.Error {
		t.Errorf("Expected Error %v, got %v", mockResponse.Error, data.Error)
	}

	// 比較 Result 的 LoginResult
	expected := mockResponse.Result
	actual := data.Result

	if expected != actual {
		t.Errorf("Expected Result LoginResult %+v, got %+v", expected, actual)
	}
}

// TODO: API測試
func TestGetBetReportAPI(t *testing.T) {
	// // 外部API測試
	// httpmock.Activate()
	// defer httpmock.DeactivateAndReset()

	// url := Httpdomain + "/report/betreport"
	// // Register a mock responder
	// httpmock.RegisterResponder(
	// 	"POST",
	// 	url,
	// 	httpmock.NewBytesResponder(200, []byte(`{"status": 1,"error": ""}`)),
	// )

	// data, err := getBetReportAPI(BetReportReq{
	// 	Agid:     1,
	// 	Username: "test",
	// })
	// assert.NoError(t, err)
	// assert.Equal(t, 1, data.Status)
	// assert.Equal(t, "", data.Error)

	// 測試物件
	mockResponse := BetReportRes{
		Status: 1,
		Error:  "",
		Result: BetReportResult{
			Lottory: LotttoryResult{
				BetList: []BetListResult{
					{
						CurId:      "NTD",
						GameId:     "test",
						GameType:   "pe",
						GameName:   "Basketball",
						Username:   "test",
						BetId:      "test123",
						EventId:    "1",
						Status:     1,
						BetTime:    time.Date(2014, time.January, 18, 0, 0, 0, 0, time.UTC),
						Result:     "YES",
						ResultCode: "PS",
						BillTime:   time.Date(2014, time.January, 18, 0, 0, 0, 0, time.UTC),
						Amt:        100,
						Payout:     100,
						PlayId:     "PS",
						Number:     "1",
						Selection:  "Basketball",
						Odds:       2,
						Rollback:   0.5,
						Ip:         "127.0.0.1",
						Device:     "iphone",
						Line:       "A",
					},
				},
			},
		},
	}

	// 透過物件轉換成 JSON byte
	mockResponseBytes, _ := json.Marshal(mockResponse)
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(mockResponseBytes)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	data, err := getBetReportAPI(
		BetReportReq{
			Agid:        1,
			Username:    "",
			Startdate:   time.Date(2014, time.January, 18, 0, 0, 0, 0, time.UTC),
			Enddate:     time.Date(2014, time.January, 18, 0, 0, 0, 0, time.UTC),
			Gametype:    "pe",
			Inquirymode: "",
		},
		server.URL,
	)
	assert.NoError(t, err)

	// 逐一檢查每個字段是否正確
	assert.Equal(t, mockResponse.Status, data.Status)
	assert.Equal(t, mockResponse.Error, data.Error)

	// 比較 Lottory 的 BetList
	expected := mockResponse.Result.Lottory.BetList
	actual := data.Result.Lottory.BetList

	assert.Equal(t, len(expected), len(actual))
	for i := range expected {
		assert.Equal(t, expected[i], actual[i])
	}
}

// TODO DBTEST　TestGetBetReportUser
func TestGetBetReport(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// 模擬插入
	mock.ExpectExec("INSERT INTO WALLET \\( username, betId, betTime, payOut, data \\) VALUES \\(\\?, \\?, \\?, \\?, \\?\\)").
		WithArgs("test", "bet123", sqlmock.AnyArg(), 100.0, "some data").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = getBetReport(db, "test", "bet123", time.Now(), 100.0, "some data")
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetChangePasswordAPI(t *testing.T) {
	// // 外部API測試
	// httpmock.Activate()
	// defer httpmock.DeactivateAndReset()

	// url := Httpdomain + "/report/betreport"
	// // Register a mock responder
	// httpmock.RegisterResponder(
	// 	"POST",
	// 	url,
	// 	httpmock.NewBytesResponder(200, []byte(`{"status": 1,"error": ""}`)),
	// )

	// data, err := getBetReportAPI(BetReportReq{
	// 	Agid:     1,
	// 	Username: "test",
	// })
	// assert.NoError(t, err)
	// assert.Equal(t, 1, data.Status)
	// assert.Equal(t, "", data.Error)

	// 測試物件
	mockResponse := AccountChangePasswordRes{
		Status: 1,
		Error:  "",
	}
	// httpmock 測試
	// 透過物件轉換成 JSON byte
	mockResponseBytes, _ := json.Marshal(mockResponse)
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(mockResponseBytes)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	data, err := getAccountChangePasswordAPI(
		AccountChangePasswordReq{
			Agid:     1,
			Username: "test",
			Password: "test123",
		},
		server.URL,
	)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// 逐一檢查每個字段是否正確
	if data.Status != mockResponse.Status {
		t.Errorf("Expected Status %v, got %v", mockResponse.Status, data.Status)
	}
	if data.Error != mockResponse.Error {
		t.Errorf("Expected Error %v, got %v", mockResponse.Error, data.Error)
	}
}

// TODO: DBTEST TESTGETAccountChangePasswordDB

func TestGetAccountBookAPI(t *testing.T) {
	// // 外部API測試
	// httpmock.Activate()
	// defer httpmock.DeactivateAndReset()

	// url := Httpdomain + "/report/accountbook"
	// // Register a mock responder
	// httpmock.RegisterResponder(
	// 	"POST",
	// 	url,
	// 	httpmock.NewBytesResponder(200, []byte(`{"status": 1,"error": ""}`)),
	// )

	// data, err := getAccountBookAPI(AccountBookReq{
	// 	Agid:     1,
	// 	Username: "test",
	// })
	// assert.NoError(t, err)
	// assert.Equal(t, 1, data.Status)
	// assert.Equal(t, "", data.Error)

	// 測試物件
	mockResponse := GetAccountBookRes{
		Status: 1,
		Error:  "",
		Result: CashListResult{
			Cashlist: CashListDTO{
				Amt:      100.00,
				DoMan:    "mike",
				BetId:    "test123",
				Remark:   "test",
				UpTime:   time.Date(2014, time.January, 18, 0, 0, 0, 0, time.UTC),
				GameType: "pe",
			},
		},
	}
	mockResponseBytes, _ := json.Marshal(mockResponse)
	// httpmock 測試
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(mockResponseBytes)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	data, err := getAccountBookAPI(
		GetAccountBookReq{
			Agid:     1,
			Username: "test",
			Gametype: "pe",
		},
		server.URL,
	)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if data.Status != mockResponse.Status {
		t.Errorf("Expected Status %v, got %v", mockResponse.Status, data.Status)
	}
	if data.Error != mockResponse.Error {
		t.Errorf("Expected Error %v, got %v", mockResponse.Error, data.Error)
	}

	// 比較 Result 的 CashList
	expected := mockResponse.Result.Cashlist
	actual := data.Result.Cashlist

	if expected != actual {
		t.Errorf("Expected Result CashList %+v, got %+v", expected, actual)
	}
}

// TODO: DBTEST TESTGetAccountBookDB

func TestGetPeriodResultAPI(t *testing.T) {
	// // 外部API測試
	// httpmock.Activate()
	// defer httpmock.DeactivateAndReset()

	// url := Httpdomain + "/report/periodresult"
	// // Register a mock responder
	// httpmock.RegisterResponder(
	// 	"POST",
	// 	url,
	// 	httpmock.NewBytesResponder(200, []byte(`{"status": 1,"error": ""}`)),
	// )

	// data, err := getPeriodResultAPI(PeriodResultReq{
	// 	Agid:     1,
	// 	Username: "test",
	// })
	// assert.NoError(t, err)
	// assert.Equal(t, 1, data.Status)
	// assert.Equal(t, "", data.Error)

	// 測試物件
	mockResponse := PeriodResultRes{
		Status: 1,
		Error:  "",
		Result: ResultResult{
			Data: ResultDTO{
				GameNum:    "1",
				GameResult: "369",
			},
		},
	}
	mockResponseBytes, _ := json.Marshal(mockResponse)
	// httpmock 測試
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(mockResponseBytes)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	data, err := getPeriodResultAPI(
		PeriodResultReq{
			Gametype: "pe",
			Gamedate: time.Date(2014, time.January, 18, 0, 0, 0, 0, time.UTC),
			Gamenum:  "1",
		},
		server.URL,
	)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if data.Status != mockResponse.Status {
		t.Errorf("Expected Status %v, got %v", mockResponse.Status, data.Status)
	}

	if data.Error != mockResponse.Error {
		t.Errorf("Expected Error %v, got %v", mockResponse.Error, data.Error)
	}

	// 比較 Result 的 Data
	expected := mockResponse.Result.Data
	actual := data.Result.Data

	if expected != actual {
		t.Errorf("Expected Result Data %+v, got %+v", expected, actual)
	}
}

// TODO: DBTEST TESTGetPeriodDB

func TestGetRealReportAPI(t *testing.T) {
	// // 外部API測試
	// httpmock.Activate()
	// defer httpmock.DeactivateAndReset()

	// url := Httpdomain + "/report/realreport"
	// // Register a mock responder
	// httpmock.RegisterResponder(
	// 	"POST",
	// 	url,
	// 	httpmock.NewBytesResponder(200, []byte(`{"status": 1,"error": ""}`)),
	// )

	// data, err := getRealReportAPI(RealReportReq{
	// 	Agid:     1,
	// 	Username: "test",
	// })
	// assert.NoError(t, err)
	// assert.Equal(t, 1, data.Status)
	// assert.Equal(t, "", data.Error)

	// 測試物件
	mockResponse := RealReportRes{
		Status: "1",
		Error:  "",
		Result: RealReportResult{
			Lottery: RealReportDTO{
				Username: "mike",
				Ordergold: RealReportOrderGold{
					GameName: "Basketball",
					Realgold: 300.00,
					Win:      100.00,
				},
			},
		},
	}
	mockResponseBytes, _ := json.Marshal(mockResponse)
	// httpmock 測試
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(mockResponseBytes)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	data, err := getRealReportAPI(
		RealReportReq{
			Agid:      1,
			Username:  "mike",
			Startdate: time.Date(2014, time.January, 18, 0, 0, 0, 0, time.UTC),
			Enddate:   time.Date(2014, time.January, 18, 0, 0, 0, 0, time.UTC),
			Gametype:  "pe",
		},
		server.URL,
	)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if data.Status != mockResponse.Status {
		t.Errorf("Expected Status %v, got %v", mockResponse.Status, data.Status)
	}

	if data.Error != mockResponse.Error {
		t.Errorf("Expected Error %v, got %v", mockResponse.Error, data.Error)
	}

	// 比較 Result 的 Data
	expected := mockResponse.Result.Lottery
	actual := data.Result.Lottery

	if expected != actual {
		t.Errorf("Expected Result Data %+v, got %+v", expected, actual)
	}
}

// TODO: DBTEST TESTGetRealReportDB

func TestGetCheckRefnoAPI(t *testing.T) {
	// // 外部API測試
	// httpmock.Activate()
	// defer httpmock.DeactivateAndReset()

	// url := Httpdomain + "/report/checkrefno"
	// // Register a mock responder
	// httpmock.RegisterResponder(
	// 	"POST",
	// 	url,
	// 	httpmock.NewBytesResponder(200, []byte(`{"status": 1,"error": ""}`)),
	// )

	// data, err := getCheckRefnoAPI(CheckRefnoReq{
	// 	Agid: 1,
	// 	Refno: "test",
	// })
	// assert.NoError(t, err)
	// assert.Equal(t, 1, data.Status)
	// assert.Equal(t, "", data.Error)

	// 測試物件
	mockResponse := RefnoRes{
		Status: "1",
		Error:  "",
		Result: RefnoResult{
			Type: "IN",
			Amt:  100.00,
		},
	}
	mockResponseBytes, _ := json.Marshal(mockResponse)
	// httpmock 測試
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(mockResponseBytes)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	data, err := getCheckRefnoAPI(
		RefnoReq{
			Agid:     1,
			Username: "mick",
			Refno:    "test20240812",
		},
		server.URL,
	)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if data.Status != mockResponse.Status {
		t.Errorf("Expected Status %v, got %v", mockResponse.Status, data.Status)
	}

	if data.Error != mockResponse.Error {
		t.Errorf("Expected Error %v, got %v", mockResponse.Error, data.Error)
	}
	expected := mockResponse.Result
	actual := data.Result

	if expected != actual {
		t.Errorf("Expected Result Data %+v, got %+v", expected, actual)
	}
}

// TODO: DBTEST TESTGetCheckRefnoDB

func TestGetTotalRealGoldAPI(t *testing.T) {
	// // 外部API測試
	// httpmock.Activate()
	// defer httpmock.DeactivateAndReset()

	// url := Httpdomain + "/report/totalrealgold"
	// // Register a mock responder
	// httpmock.RegisterResponder(
	// 	"POST",
	// 	url,
	// 	httpmock.NewBytesResponder(200, []byte(`{"status": 1,"error": ""}`)),
	// )

	// data, err := getTotalRealGoldAPI(TotalRealGoldReq{
	// 	Agid: 1,
	// 	Username: "test",
	// })
	// assert.NoError(t, err)
	// assert.Equal(t, 1, data.Status)
	// assert.Equal(t, "", data.Error)

	// 測試物件
	mockResponse := TotalRealGoldRes{
		Status: "1",
		Error:  "",
		Result: TotalRealGoldResult{
			Gametype:      "pe",
			Num:           "50",
			Totalgold:     100.00,
			Totalrealgold: 100.00,
			Totalwingold:  100.00,
		},
	}
	mockResponseBytes, _ := json.Marshal(mockResponse)
	// httpmock 測試
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(mockResponseBytes)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	data, err := getTotalRealGoldAPI(
		TotalRealGoldReq{
			Agid:      1,
			Gametype:  "pe",
			Startdate: time.Date(2014, time.January, 18, 0, 0, 0, 0, time.UTC),
			Enddate:   time.Date(2014, time.January, 18, 0, 0, 0, 0, time.UTC),
		},
		server.URL,
	)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if data.Status != mockResponse.Status {
		t.Errorf("Expected Status %v, got %v", mockResponse.Status, data.Status)
	}

	if data.Error != mockResponse.Error {
		t.Errorf("Expected Error %v, got %v", mockResponse.Error, data.Error)
	}
	expected := mockResponse.Result
	actual := data.Result

	if expected != actual {
		t.Errorf("Expected Result Data %+v, got %+v", expected, actual)
	}
}

// TODO: DBTEST TESTGetTotalRealGoldDB

func TestGetFhDetailsAPI(t *testing.T) {
	// // 外部API測試
	// httpmock.Activate()
	// defer httpmock.DeactivateAndReset()

	// url := Httpdomain + "/report/totalrealgold"
	// // Register a mock responder
	// httpmock.RegisterResponder(
	// 	"POST",
	// 	url,
	// 	httpmock.NewBytesResponder(200, []byte(`{"status": 1,"error": ""}`)),
	// )

	// data, err := getTotalRealGoldAPI(TotalRealGoldReq{
	// 	Agid: 1,
	// 	Username: "test",
	// })
	// assert.NoError(t, err)
	// assert.Equal(t, 1, data.Status)
	// assert.Equal(t, "", data.Error)

	// 測試物件
	mockResponse := FhDetailsRes{
		Status: "1",
		Error:  "",
		Result: FhDetailsResult{
			Id:            "1",
			TableId:       "1",
			CreateTime:    time.Date(2014, time.January, 18, 0, 0, 0, 0, time.UTC),
			BeforeBalance: 150.00,
			AfterBalance:  100.00,
			Bet:           50.00,
			BetWin:        50.00,
			WinLoss:       50.00,
			ProcessStatus: "shark",
			FishSpecies:   "34",
		},
	}
	mockResponseBytes, _ := json.Marshal(mockResponse)
	// httpmock 測試
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(mockResponseBytes)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	data, err := getFhDetailsAPI(
		FhDetailsReq{
			Agid:  1,
			Betid: "1",
		},
		server.URL,
	)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if data.Status != mockResponse.Status {
		t.Errorf("Expected Status %v, got %v", mockResponse.Status, data.Status)
	}

	if data.Error != mockResponse.Error {
		t.Errorf("Expected Error %v, got %v", mockResponse.Error, data.Error)
	}
	expected := mockResponse.Result
	actual := data.Result

	if expected != actual {
		t.Errorf("Expected Result Data %+v, got %+v", expected, actual)
	}
}

// TODO: DBTEST TESTGetFhDetailsDB

func TestGetRecReportAPI(t *testing.T) {
	// // 外部API測試
	// httpmock.Activate()
	// defer httpmock.DeactivateAndReset()

	// url := Httpdomain + "/report/totalrealgold"
	// // Register a mock responder
	// httpmock.RegisterResponder(
	// 	"POST",
	// 	url,
	// 	httpmock.NewBytesResponder(200, []byte(`{"status": 1,"error": ""}`)),
	// )

	// data, err := getTotalRealGoldAPI(TotalRealGoldReq{
	// 	Agid: 1,
	// 	Username: "test",
	// })
	// assert.NoError(t, err)
	// assert.Equal(t, 1, data.Status)
	// assert.Equal(t, "", data.Error)

	// 測試物件
	mockResponse := RecReportRes{
		Status: "1",
		Error:  "",
		Result: RecReportResult{
			TotalNum:       "5",
			TotalBetAmount: 100.00,
			TotalWinGold:   100.00,
			TotalRealGold:  100.00,
		},
	}
	mockResponseBytes, _ := json.Marshal(mockResponse)
	// httpmock 測試
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(mockResponseBytes)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	data, err := getRecReportAPI(
		RecReportReq{
			Agid:      1,
			Username:  "Lebron",
			Startdate: time.Date(2014, time.January, 18, 0, 0, 0, 0, time.UTC),
			Enddate:   time.Date(2014, time.January, 18, 0, 0, 0, 0, time.UTC),
		},
		server.URL,
	)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if data.Status != mockResponse.Status {
		t.Errorf("Expected Status %v, got %v", mockResponse.Status, data.Status)
	}

	if data.Error != mockResponse.Error {
		t.Errorf("Expected Error %v, got %v", mockResponse.Error, data.Error)
	}

	expected := mockResponse.Result
	actual := data.Result

	if expected != actual {
		t.Errorf("Expected Result Data %+v, got %+v", expected, actual)
	}
}

//TODO: DBTEST TESTGetRecReportDB

func TestGetAccountLogoutAPI(t *testing.T) {
	// // 外部API測試
	// httpmock.Activate()
	// defer httpmock.DeactivateAndReset()

	// url := Httpdomain + "/report/totalrealgold"
	// // Register a mock responder
	// httpmock.RegisterResponder(
	// 	"POST",
	// 	url,
	// 	httpmock.NewBytesResponder(200, []byte(`{"status": 1,"error": ""}`)),
	// )

	// data, err := getTotalRealGoldAPI(TotalRealGoldReq{
	// 	Agid: 1,
	// 	Username: "test",
	// })
	// assert.NoError(t, err)
	// assert.Equal(t, 1, data.Status)
	// assert.Equal(t, "", data.Error)

	// 測試物件
	mockResponse := AccountLogoutRes{
		Status: "1",
		Error:  "",
	}
	mockResponseBytes, _ := json.Marshal(mockResponse)
	// httpmock 測試
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(mockResponseBytes)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	data, err := getAccountLogoutAPI(
		AccountLogoutReq{
			Agid:     1,
			Username: "Lebron",
		},
		server.URL,
	)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if data.Status != mockResponse.Status {
		t.Errorf("Expected Status %v, got %v", mockResponse.Status, data.Status)
	}

	if data.Error != mockResponse.Error {
		t.Errorf("Expected Error %v, got %v", mockResponse.Error, data.Error)
	}
}

// TODO: DBTEST TESTGetAccountLogoutDB

func TestGetCheckOnlineStatusAPI(t *testing.T) {
	// // 外部API測試
	// httpmock.Activate()
	// defer httpmock.DeactivateAndReset()

	// url := Httpdomain + "/report/totalrealgold"
	// // Register a mock responder
	// httpmock.RegisterResponder(
	// 	"POST",
	// 	url,
	// 	httpmock.NewBytesResponder(200, []byte(`{"status": 1,"error": ""}`)),
	// )

	// data, err := getTotalRealGoldAPI(TotalRealGoldReq{
	// 	Agid: 1,
	// 	Username: "test",
	// })
	// assert.NoError(t, err)
	// assert.Equal(t, 1, data.Status)
	// assert.Equal(t, "", data.Error)

	// 測試物件
	mockResponse := OnlineStatusRes{
		Status: "1",
		Error:  "",
		Result: OnlineStatusResult{
			OnlineStatus: "1",
		},
	}
	mockResponseBytes, _ := json.Marshal(mockResponse)
	// httpmock 測試
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(mockResponseBytes)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	data, err := getCheckOnlineStatusAPI(
		OnlineStatusReq{
			Agid:     1,
			Username: "Lebron",
			Gametype: "pe",
		},
		server.URL,
	)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if data.Status != mockResponse.Status {
		t.Errorf("Expected Status %v, got %v", mockResponse.Status, data.Status)
	}

	if data.Error != mockResponse.Error {
		t.Errorf("Expected Error %v, got %v", mockResponse.Error, data.Error)
	}

	expected := mockResponse.Result
	actual := data.Result

	if expected != actual {
		t.Errorf("Expected Result Data %+v, got %+v", expected, actual)
	}
}

// TODO: DBTEST TESTGetCheckOnlineStatusDB

func TestGetLineStatusAPI(t *testing.T) {
	// // 外部API測試
	// httpmock.Activate()
	// defer httpmock.DeactivateAndReset()

	// url := Httpdomain + "/report/totalrealgold"
	// // Register a mock responder
	// httpmock.RegisterResponder(
	// 	"POST",
	// 	url,
	// 	httpmock.NewBytesResponder(200, []byte(`{"status": 1,"error": ""}`)),
	// )

	// data, err := getTotalRealGoldAPI(TotalRealGoldReq{
	// 	Agid: 1,
	// 	Username: "test",
	// })
	// assert.NoError(t, err)
	// assert.Equal(t, 1, data.Status)
	// assert.Equal(t, "", data.Error)

	// 測試物件
	mockResponse := GetLineRes{
		Status: "1",
		Error:  "",
		Result: GetLineResult{
			Account: "test123",
			ALine:   "Y",
			BLine:   "Y",
			CLine:   "N",
			DLine:   "N",
		},
	}
	mockResponseBytes, _ := json.Marshal(mockResponse)
	// httpmock 測試
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(mockResponseBytes)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	data, err := getLineStatusAPI(
		GetLineReq{
			Agid:     1,
			Username: "Lebron",
		},
		server.URL,
	)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if data.Status != mockResponse.Status {
		t.Errorf("Expected Status %v, got %v", mockResponse.Status, data.Status)
	}

	if data.Error != mockResponse.Error {
		t.Errorf("Expected Error %v, got %v", mockResponse.Error, data.Error)
	}

	expected := mockResponse.Result
	actual := data.Result

	if expected != actual {
		t.Errorf("Expected Result Data %+v, got %+v", expected, actual)
	}
}

// TODO: DBTEST TESTGetLineStatusDB

func TestGetPreDepositAPI(t *testing.T) {
	// // 外部API測試
	// httpmock.Activate()
	// defer httpmock.DeactivateAndReset()

	// url := Httpdomain + "/report/totalrealgold"
	// // Register a mock responder
	// httpmock.RegisterResponder(
	// 	"POST",
	// 	url,
	// 	httpmock.NewBytesResponder(200, []byte(`{"status": 1,"error": ""}`)),
	// )

	// data, err := getTotalRealGoldAPI(TotalRealGoldReq{
	// 	Agid: 1,
	// 	Username: "test",
	// })
	// assert.NoError(t, err)
	// assert.Equal(t, 1, data.Status)
	// assert.Equal(t, "", data.Error)

	// 測試物件
	mockResponse := PreDepositRes{
		Status: "1",
		Error:  "",
		Result: PreDepositResult{
			Preno: "test123",
		},
	}
	mockResponseBytes, _ := json.Marshal(mockResponse)
	// httpmock 測試
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(mockResponseBytes)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	data, err := getPreDepositAPI(
		PreDepositReq{
			Agid:     1,
			Username: "Lebron",
			Gametype: "pe",
			Amt:      100.00,
			Refno:    "test123",
			Type:     "IN",
		},
		server.URL,
	)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if data.Status != mockResponse.Status {
		t.Errorf("Expected Status %v, got %v", mockResponse.Status, data.Status)
	}

	if data.Error != mockResponse.Error {
		t.Errorf("Expected Error %v, got %v", mockResponse.Error, data.Error)
	}

	expected := mockResponse.Result.Preno
	actual := data.Result.Preno

	if expected != actual {
		t.Errorf("Expected Result Data %+v, got %+v", expected, actual)
	}
}

// TODO: DBTEST TESTGetPreDepositDB

func TestGetCheckDepositAPI(t *testing.T) {
	// // 外部API測試
	// httpmock.Activate()
	// defer httpmock.DeactivateAndReset()

	// // Register a mock responder
	// httpmock.RegisterResponder(
	// 	"POST",
	// 	Httpdomain+"/report/checkdeposit",
	// 	httpmock.NewBytesResponder(200, []byte(`{"status": 1,"error": ""}`)),
	// )

	// data, err := getCheckDepositAPI(CheckDepositReq{
	// 	Agid: 1,
	// 	Username: "test",
	// })
	// assert.NoError(t, err)
	// assert.Equal(t, 1, data.Status)
	// assert.Equal(t, "", data.Error)

	// 測試物件
	mockResponse := CheckDepositRes{
		Status: "1",
		Error:  "",
		Result: CheckDepositResult{
			Refno:     "test123",
			PayMentId: "testtest123",
			Balance:   100.00,
			Amt:       100.00,
		},
	}
	mockResponseBytes, _ := json.Marshal(mockResponse)
	// httpmock 測試
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(mockResponseBytes)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	data, err := getCheckDepositAPI(
		CheckDepositReq{
			Agid:  1,
			Preno: "test123",
		},
		server.URL,
	)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if data.Status != mockResponse.Status {
		t.Errorf("Expected Status %v, got %v", mockResponse.Status, data.Status)
	}

	if data.Error != mockResponse.Error {
		t.Errorf("Expected Error %v, got %v", mockResponse.Error, data.Error)
	}

	expected := mockResponse.Result
	actual := data.Result

	if expected != actual {
		t.Errorf("Expected Result Data %+v, got %+v", expected, actual)
	}
}

// TODO: DBTEST TESTGetCheckDepositDB
