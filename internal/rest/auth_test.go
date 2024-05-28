package rest

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/golang/mock/gomock"
	mockdb "github.com/kavinddd/mangtoon_be/internal/db/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogin(t *testing.T) {
	user, password := randomUser()

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username": user.Username,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), user.Username).
					Times(1).
					Return(user, nil)
				store.EXPECT().
					ListRolesByUserId(gomock.Any(), user.ID).
					Times(1)
				store.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()). // TODO: change second param to CreateSessionParam struct
					Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// build stubs
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			// prepare to send request
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)
			url := fmt.Sprintf("/login")
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			// send request (response is stored in recorder)
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()
			server.router.ServeHTTP(recorder, request)

			// check response
			tc.checkResponse(t, recorder)
		})
	}

}
