package httpadp

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	mock_app "video-provider/internal/user-service/app/mock"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestValidCreateUserRequest(t *testing.T) {
	userId := "123e4567-e89b-12d3-a456-426614174000"
	expResponse := "\"" + userId + "\"\n"
	expUserID := uuid.MustParse(userId)

	cases := []struct {
		name    string
		reqBody string
	}{
		{
			"Valid user creation request",
			`{"email":"test@example.com","name":"John","lastname":"Doe","password":"Password123!!"}`,
		},
	}

	t.Parallel()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			s := mock_app.NewMockUserInteractor(ctrl)
			h := NewUserHandler(s, log.New(io.Discard, "", 0))
			r := mux.NewRouter()
			SetupRouter(r, h)

			s.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(expUserID, nil).MaxTimes(1)

			rec := httptest.NewRecorder()

			r.ServeHTTP(rec,
				httptest.NewRequest(
					http.MethodPost,
					routeUsers,
					strings.NewReader(c.reqBody)))

			require.Equal(t, http.StatusCreated, rec.Code)
			require.Equal(t, expResponse, rec.Body.String())
		})
	}
}

func TestInvalidCreateUserRequest(t *testing.T) {
	cases := []struct {
		name          string
		expStatusCode int
		reqBody       string
	}{
		{
			"Invalid email format", http.StatusBadRequest,
			`{"email":"invalid@email","name":"John","lastname":"Doe","password":"Password123!!"}`,
		},
		{
			"Short email", http.StatusBadRequest,
			`{"email":"email@","name":"John","lastname":"Doe","password":"Password123!!"}`,
		},
		{
			"Long email", http.StatusBadRequest,
			`{"email":"emaillllllllllllllllllllllllllllllllllllllllllllllllllllllll@longdomainname.com","name":"John","lastname":"Doe","password":"Password123!!"}`,
		},
		{
			"Invalid name format", http.StatusBadRequest,
			`{"email":"test@example.com","name":"John!!","lastname":"Doe","password":"Password123!!"}`,
		},
		{
			"Short name", http.StatusBadRequest,
			`{"email":"test@example.com","name":"J","lastname":"Doe","password":"Password123!!"}`,
		},
		{
			"Too long name", http.StatusBadRequest,
			`{"email":"test@example.com","name":"JOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOHN","lastname":"Doe","password":"Password123!!"}`,
		},
		{
			"Invalid lastname format", http.StatusBadRequest,
			`{"email":"test@example.com","name":"John","lastname":"Doe!!!","password":"Password123!!"}`,
		},
		{
			"Short lastname", http.StatusBadRequest,
			`{"email":"test@example.com","name":"John","lastname":"D","password":"Password123!!"}`,
		},
		{
			"Too long lastname", http.StatusBadRequest,
			`{"email":"test@example.com","name":"John","lastname":"DOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOE","password":"Password123!!"}`,
		},
		{
			"No digits in password", http.StatusBadRequest,
			`{"email":"test@example.com","name":"John","lastname":"Doe","password":"pEEsword!!"}`,
		},
		{
			"No spec chars password", http.StatusBadRequest,
			`{"email":"test@example.com","name":"John","lastname":"Doe","password":"pEEsword11"}`,
		},
		{
			"No cap chars password", http.StatusBadRequest,
			`{"email":"test@example.com","name":"John","lastname":"Doe","password":"peesword!!11"}`,
		},
		{
			"No reg chars password", http.StatusBadRequest,
			`{"email":"test@example.com","name":"John","lastname":"Doe","password":"PEESWORD!!11"}`,
		},
	}

	t.Parallel()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			s := mock_app.NewMockUserInteractor(ctrl)
			h := NewUserHandler(s, log.New(io.Discard, "", 0))
			r := mux.NewRouter()
			SetupRouter(r, h)

			s.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).MaxTimes(0)

			rec := httptest.NewRecorder()

			r.ServeHTTP(rec,
				httptest.NewRequest(
					http.MethodPost,
					routeUsers,
					strings.NewReader(c.reqBody)))

			require.Equal(t, c.expStatusCode, rec.Code)
			require.NotEmpty(t, rec.Body.String())
		})
	}
}
