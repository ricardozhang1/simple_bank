package api

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
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

func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser(t)

	// 测试样例
	testCase := []struct{
		name			string
		body			gin.H
		buildStubs		func(store *mockdb.MockStore)
		checkResponse	func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username":	user.Username,
				"password":	password,
				"full_name": user.FullName,
				"email": user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email: user.Email,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
	}

	for i := range testCase {
		tc := testCase[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/users"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

// randomUser 随机生成一个User
func randomUser(t *testing.T) (user db.User, password string) {
	password = util.RandomString(6)
	hashPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		Username: 		util.RandomOwnerName(),
		HashedPassword: hashPassword,
		FullName: 		util.RandomFullName(),
		Email: 			util.RandomEmail(),
	}
	return
}

// requireBodyMatchUser 对返回的数据与请求数据进行check
func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)

	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.FullName, gotUser.FullName)
	require.Equal(t, user.Email, gotUser.Email)
	require.NotEmpty(t, gotUser.HashedPassword)
}
