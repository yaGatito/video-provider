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

	pubID := uuid.Must(uuid.NewRandom())
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
				Description: expDec,
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

	vidID := uuid.Must(uuid.NewRandom())

	pubID := uuid.Must(uuid.NewRandom())
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
				Description: expDec,
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

func TestSearchPublisherVideos(t *testing.T) {
	ctrl := gomock.NewController(t)
	s := mock_app.NewMockVideoService(ctrl)
	h := httpadapter.NewVideoHandler(s, nil, log.New(io.Discard, "", 0))
	r := mux.NewRouter()
	httpadapter.SetupRouter(r, h)

	pubID := uuid.Must(uuid.NewRandom())
	expTop := "topic"
	expDec := "description"

	cases := []struct {
		name             string
		wantErr          bool
		pubID            string
		urlParams        string
		expCallByPub     int
		expCallSearchPub int
	}{
		// SearchPublisher
		{"ok search", false,
			pubID.String(), "?offset=0&limit=5&query=search", 0, 1},
		{"few words search", false,
			pubID.String(), "?offset=0&limit=5&query=search+lorem+ipsum+kitty+dolor", 0, 1},
		{"negative offset search", false,
			pubID.String(), "?offset=-10&limit=5&query=search", 0, 1},
		{"invalid offset search", false,
			pubID.String(), "?offset=H&limit=5&query=search", 0, 1},
		{"negative limit search", false,
			pubID.String(), "?offset=0&limit=-10&query=search", 0, 1},
		{"zero limit search", false,
			pubID.String(), "?offset=0&limit=0&query=search", 0, 1},
		{"invalid limit search", false,
			pubID.String(), "?offset=0&limit=H&query=search", 0, 1},

		// GetByPublisher
		{"ok", false,
			pubID.String(), "?offset=0&limit=5", 1, 0},
		{"negative offset", false,
			pubID.String(), "?offset=-10&limit=5", 1, 0},
		{"invalid offset", false,
			pubID.String(), "?offset=H&limit=5", 1, 0},
		{"negative limit", false,
			pubID.String(), "?offset=0&limit=-10", 1, 0},
		{"zero limit", false,
			pubID.String(), "?offset=0&limit=0", 1, 0},
		{"invalid limit", false,
			pubID.String(), "?offset=0&limit=H", 1, 0},
		{"wrong url params", false,
			pubID.String(), "?affset=0&leemeet=5&kueree=search", 1, 0},
		{"partial wrong url params", false,
			pubID.String(), "?affset=0&leemeet=5&search=", 1, 0},

		// Errors. TODO: test in real time
		// {"invalid search", true,
		// 	pubID.String(), "?offset=0&limit=5&query=search with spaces", 0, 0},
		{"invalid search", true,
			pubID.String(), "?offset=0&limit=5&query=,.<>?/;", 0, 0},
		{"ivalid id", true,
			"1", "?offset=0&limit=5&query=search", 0, 0},
		{"empty id", true,
			"", "?offset=0&limit=5&query=search", 0, 0},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			if c.wantErr {
				s.EXPECT().
					GetByPublisher(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					MaxTimes(c.expCallByPub)
				s.EXPECT().
					SearchPublisher(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					MaxTimes(c.expCallSearchPub)
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
						Description: expDec,
					}}, nil,
				).MaxTimes(c.expCallByPub)

				s.EXPECT().SearchPublisher(
					gomock.Any(),
					gomock.Eq(uuid.MustParse(c.pubID)),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(
					[]domain.Video{{
						PublisherID: pubID,
						Topic:       expTop,
						Description: expDec,
					}}, nil,
				).MaxTimes(c.expCallSearchPub)
			}

			url := strings.Replace(
				httpadapter.RoutePublisherVideos,
				"{"+httpadapter.PathVarPublisherID+"}",
				c.pubID+c.urlParams,
				1,
			)
			req := httptest.NewRequest(http.MethodGet, url, nil)
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
