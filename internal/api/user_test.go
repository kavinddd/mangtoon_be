package api

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/kavinddd/mangtoon_be/internal/db"
	mockdb "github.com/kavinddd/mangtoon_be/internal/db/mock"
	"github.com/kavinddd/mangtoon_be/pkg/util"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func randomUser() db.User {
	return db.User{
		ID:        uuid.New(),
		Username:  util.RandomUsername(),
		Email:     util.RandomEmail(),
		Password:  util.RandomString(10),
		IsActive:  false,
		CreatedAt: time.Time{},
	}
}
func randomUsers(n int) []db.User {
	var users []db.User
	for i := 0; i < n; i++ {
		users = append(users, randomUser())
	}
	return users
}

func requireBodyMatchUserResponse(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var userFromBody db.User
	err = json.Unmarshal(data, &userFromBody)
	require.NoError(t, err)
	require.Equal(t, user, userFromBody)
}

func TestFindUserById(t *testing.T) {
	user := randomUser()

	testCases := []struct {
		name          string
		userId        string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			userId: user.ID.String(),
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUserResponse(t, recorder.Body, user)
			},
		},
		{
			name:   "NotFound",
			userId: user.ID.String(),
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:   "InternalError",
			userId: user.ID.String(),
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(db.User{}, sql.ErrConnDone) // anything except ErrNoRow since we are handling it
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},

		{
			name:   "InvalidId",
			userId: "thisisnotauuid",
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
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
			url := fmt.Sprintf("/users/%s", tc.userId)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// send request (response is stored in recorder)
			server := NewServer(store)
			recorder := httptest.NewRecorder()
			server.router.ServeHTTP(recorder, request)

			// check response
			tc.checkResponse(t, recorder)
		})
	}
}

func requireBodyMatchListUsersResponse(t *testing.T, body *bytes.Buffer, users []db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)
	var usersFromBody []db.User
	err = json.Unmarshal(data, &usersFromBody)
	require.NoError(t, err)
	require.Len(t, usersFromBody, len(users))
	require.Equal(t, users, usersFromBody)
}

func TestListUsers(t *testing.T) {

	users := randomUsers(15)

	testCases := []struct {
		name          string
		pageNo        int32
		pageSize      int32
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			pageNo:   1,
			pageSize: 10,
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					ListUsers(gomock.Any(), gomock.Eq(db.ListUsersParams{
						Limit:  10,
						Offset: 0,
					})).
					Times(1).
					Return(users, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchListUsersResponse(t, recorder.Body, users)
			},
		},
		{
			name:     "Empty",
			pageNo:   999,
			pageSize: 10,
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					ListUsers(gomock.Any(), gomock.Eq(db.ListUsersParams{
						Limit:  10,
						Offset: 9980,
					})).
					Times(1).
					Return([]db.User{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchListUsersResponse(t, recorder.Body, []db.User{})
			},
		},
		{
			name:     "InternalError",
			pageNo:   1,
			pageSize: 10,
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					ListUsers(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:     "BadRequest",
			pageNo:   1,
			pageSize: 20, // exceed maximum
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					ListUsers(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
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
			url := fmt.Sprintf("/users?page_size=%d&page_no=%d", tc.pageSize, tc.pageNo)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// send request (response is stored in recorder)
			server := NewServer(store)
			recorder := httptest.NewRecorder()
			server.router.ServeHTTP(recorder, request)

			// check response
			tc.checkResponse(t, recorder)
		})
	}

}
