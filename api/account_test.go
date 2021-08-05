package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockdb "github.com/techschool/simplebank/db/mock"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/util"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAccountAPI(t *testing.T) {
	// 返回一个用户账户
	account := randomAccount()

	//testCases := []struct{
	//
	//}

	// 生成一个gomock的controlor
	ctrl := gomock.NewController(t)
	// 关闭controlor
	defer ctrl.Finish()

	// 生成一个MockStore对象
	store := mockdb.NewMockStore(ctrl)
	// build stubs
	store.EXPECT().  // 返回的是一个MockStoreMockRecorder
		GetAccount(gomock.Any(), gomock.Eq(account.ID)).
		Times(1).
		Return(account, nil)
	// start test server and request
	server := newTestServer(t, store)
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/accounts/%d", account.ID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	server.router.ServeHTTP(recorder, request)
	// check response
	require.Equal(t, http.StatusOK, recorder.Code)
	requireBodyMatchAccount(t, recorder.Body, account)
}

func randomAccount() db.Account {
	return db.Account{
		ID: util.RandomInt(1, 1000),
		Owner: util.RandomOwnerName(),
		Balance: util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}
