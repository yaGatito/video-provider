package httpadapter_test

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	httpadapter "video-service/internal/adapters/http"
	mock_app "video-service/internal/app/mock"
	"video-service/internal/domain"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestCreateVideo(t *testing.T) {
	ctrl := gomock.NewController(t)
	s := mock_app.NewMockVideoService(ctrl)
	h := httpadapter.NewVideoHandler(s, nil, log.New(io.Discard, "", 0))
	r := mux.NewRouter()
	httpadapter.SetupRouter(r, h)

	pubID := uuid.New()
	expTop := "topic"
	expDec := "description"
	reqBody := `{
		"topic": "` + expTop + `",aa
		"description": "` + expDec + `"
	}`

	cases := []struct {
		name    string
		wantErr bool
		pubID   string
		reqBody string
	}{
		{"ok", false, pubID.String(), getReqBody(expTop, expDec)},
		{"ivalid id", true, "1", getReqBody(expTop, expDec)},
		{"ivalid req body", true, pubID.String(), "budyyy"},
		{"ivalid req body 1", true, pubID.String(), reqBody + " 1"},
		{"ivalid req body 2", true, pubID.String(), "1 " + reqBody},
	}

	for _, c := range cases {
		if c.wantErr {
			s.EXPECT().Create(gomock.Any(), gomock.Any()).MaxTimes(0)
		} else {
			s.EXPECT().Create(gomock.Any(), gomock.Eq(domain.Video{
				PublisherID: pubID,
				Topic:       expTop,
				Description: &expDec,
			})).MaxTimes(1)
		}

		req := httptest.NewRequest(
			http.MethodPost,
			strings.Replace(
				httpadapter.RoutePublisherVideos,
				"{"+httpadapter.PathVarPublisherID+"}",
				c.pubID,
				1,
			),
			strings.NewReader(c.reqBody),
		)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		if c.wantErr {
			require.NotEqual(t, http.StatusOK, rec.Code)
		} else {
			require.Equal(t, http.StatusOK, rec.Code)
		}
	}

}

func getReqBody(topic, desc string) string {
	return `{
		"topic": "` + topic + `",
		"description": "` + desc + `"
	}`
}
