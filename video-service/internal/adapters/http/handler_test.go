package httpadapter_test

import (
	"fmt"
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

func TestGetVideoById(t *testing.T) {
	ctrl := gomock.NewController(t)
	s := mock_app.NewMockVideoService(ctrl)
	h := httpadapter.NewVideoHandler(s, nil, log.New(io.Discard, "", 0))
	r := mux.NewRouter()
	httpadapter.SetupRouter(r, h)

	vidID := uuid.New()

	pubID := uuid.New()
	expTop := "topic"
	expDec := "description"

	cases := []struct {
		name    string
		wantErr bool
		vidID   string
	}{
		{"ok", false, vidID.String()},
		{"ivalid id", true, "1"},
		{"empty id", true, ""},
	}

	for _, c := range cases {
		if c.wantErr {
			s.EXPECT().GetByID(gomock.Any(), gomock.Any()).MaxTimes(0)
		} else {
			s.EXPECT().GetByID(gomock.Any(), gomock.Eq(uuid.MustParse(c.vidID))).Return(domain.Video{
				PublisherID: pubID,
				Topic:       expTop,
				Description: &expDec,
			}, nil).MaxTimes(1)
		}

		req := httptest.NewRequest(
			http.MethodGet,
			strings.Replace(
				httpadapter.RouteVideo,
				"{"+httpadapter.PathVarVideoID+"}",
				c.vidID,
				1,
			),
			nil,
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

func TestGetVideosByPublisher(t *testing.T) {
	ctrl := gomock.NewController(t)
	s := mock_app.NewMockVideoService(ctrl)
	h := httpadapter.NewVideoHandler(s, nil, log.New(io.Discard, "", 0))
	r := mux.NewRouter()
	httpadapter.SetupRouter(r, h)

	pubID := uuid.New()
	expTop := "topic"
	expDec := "description"

	cases := []struct {
		name    string
		wantErr bool
		pubID   string
		offset  string
		limit   string
	}{
		{"ok", false, pubID.String(), "0", "5"},
		{"negative offset", false, pubID.String(), "-10", "5"},
		{"invalid offset", false, pubID.String(), "A", "5"},
		{"negative limit", false, pubID.String(), "0", "-5"},
		{"zero limit", false, pubID.String(), "0", "0"},
		{"invalid limit", false, pubID.String(), "0", "B"},
		{"ivalid id", true, "1", "0", "5"},
		{"empty id", true, "", "0", "5"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			if c.wantErr {
				s.EXPECT().
					GetByPublisher(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					MaxTimes(0)
				s.EXPECT().
					SearchPublisher(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					MaxTimes(0)
			} else {
				s.EXPECT().GetByPublisher(
					gomock.Any(),
					gomock.Eq(uuid.MustParse(c.pubID)),
					gomock.Any(),
					gomock.Any(),
				).Return(
					[]domain.Video{{
						PublisherID: pubID,
						Topic:       expTop,
						Description: &expDec,
					}}, nil,
				).MaxTimes(1)
				s.EXPECT().
					SearchPublisher(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					MaxTimes(0)
			}

			req := httptest.NewRequest(
				http.MethodGet,
				strings.Replace(
					httpadapter.RoutePublisherVideos,
					"{"+httpadapter.PathVarPublisherID+"}",
					fmt.Sprintf(
						"%s?%s=%s&%s=%s",
						c.pubID,
						httpadapter.URLParamOffset,
						c.offset,
						httpadapter.URLParamLimit,
						c.limit),
					1),
				nil)

			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			if c.wantErr {
				require.NotEqual(t, http.StatusOK, rec.Code)
			} else {
				require.Equal(t, http.StatusOK, rec.Code)
			}
		})
	}
}

