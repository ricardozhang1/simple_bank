package api

import (
	"bytes"
	"database/sql"
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

	testCases := []struct{
		name string
		accountID int64
		buildStubs func(store *mockdb.MockStore)
	}{
		{
			name: "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
		},
		{
			name: "NotFound",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
		},
		{
			name: "InternalError",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
		},
		{
			name: "InvalidID",
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			// 生成一个gomock的controlor
			ctrl := gomock.NewController(t)
			// 关闭controlor
			defer ctrl.Finish()

			// 生成一个MockStore对象
			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			// start test server and request and recorder.
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", account.ID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			// check response
			require.Equal(t, http.StatusOK, recorder.Code)
			requireBodyMatchAccount(t, recorder.Body, account)
		})
	}
}

// TODO: TestCreateAccountAPI 完成登陆模块后完成
//func TestCreateAccountAPI(t *testing.T) {
//	account := randomAccount()
//
//}


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
