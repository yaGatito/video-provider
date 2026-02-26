package httpadapter_test

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	httpadapter "video-provider/internal/video-service/adapters/http"
	"video-provider/internal/video-service/adapters/idgen"
	mock_app "video-provider/internal/video-service/app/mock"
	"video-provider/internal/video-service/domain"
	"video-provider/internal/video-service/policy"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestCreateVideo(t *testing.T) {
	ctrl := gomock.NewController(t)
	s := mock_app.NewMockVideoService(ctrl)
	h := httpadapter.NewVideoHandler(s, idgen.New(), log.New(io.Discard, "", 0))
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
	h := httpadapter.NewVideoHandler(s, idgen.New(), log.New(io.Discard, "", 0))
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

func TestGetByPublisherVideos(t *testing.T) {
	pubIDStr := "d9fa522f-0016-464f-8d68-356ba1d6ad7d"
	expTop := "topic"
	expDec := "description"

	expectedRes := []domain.Video{{
		PublisherID: uuid.Must(uuid.Parse(pubIDStr)),
		Topic:       expTop,
		Description: expDec,
	}}

	cases := []struct {
		name      string
		wantErr   bool
		pubID     string
		urlParams string
	}{
		{"ok", false,
			pubIDStr, "?offset=0&limit=5&orderBy=createdAt&asc=t"},

		{"large URL query", true,
			pubIDStr, "?offset=5&limit=5&orderBy=createdAt&asc=" + strings.Repeat("a", policy.UrlMaxLen)},

		{"invalid offset", true,
			pubIDStr, "?offset=H&limit=5&orderBy=createdAt&asc=t"},

		{"no offset in URL", true,
			pubIDStr, "?limit=5&orderBy=createdAt&asc=t"},

		{"invalid limit", true,
			pubIDStr, "?offset=0&limit=H&orderBy=createdAt&asc=t"},

		{"no limit in URL", true,
			pubIDStr, "?offset=0&orderBy=createdAt&asc=t"},

		{"invalid asc", true,
			pubIDStr, "?offset=0&limit=5&orderBy=createdAt&asc=!@#%"},

		{"no asc in URL", true,
			pubIDStr, "?offset=0&limit=5&orderBy=createdAt"},

		{"invalid orderBy", true,
			pubIDStr, "?offset=0&limit=5&orderBy=!@#%&asc=t"},

		{"no orderBy in URL", true,
			pubIDStr, "?offset=0&limit=5&asc=t"},

		{"wrong url params", true,
			pubIDStr, "?affset=0&leemeet=5&orderBy=createdAt&asc=t&kueree=search"},

		{"partial wrong url params", true,
			pubIDStr, "?affset=0&leemeet=5&orderBy=createdAt&asc=t&search="},
	}

	t.Parallel()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			s := mock_app.NewMockVideoService(ctrl)
			h := httpadapter.NewVideoHandler(s, idgen.New(), log.New(io.Discard, "", 0))
			r := mux.NewRouter()
			httpadapter.SetupRouter(r, h)

			if c.wantErr {
				s.EXPECT().
					GetByPublisher(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					MaxTimes(0)
			} else {
				s.EXPECT().GetByPublisher(
					gomock.Any(),
					gomock.Eq(uuid.MustParse(c.pubID)),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(
					expectedRes, nil,
				).MaxTimes(1)
			}

			url := strings.Replace(
				httpadapter.RoutePublisherVideos,
				"{"+httpadapter.PathVarPublisherID+"}",
				c.pubID,
				1,
			) + c.urlParams

			req := httptest.NewRequest(http.MethodGet, url, nil)
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			if c.wantErr {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			} else {
				require.Equal(t, http.StatusOK, rec.Code)
			}
		})
	}
}

func TestSearchPublisherVideos(t *testing.T) {
	pubIDStr := "d9fa522f-0016-464f-8d68-356ba1d6ad7d"
	expTop := "topic"
	expDec := "description"

	expectedRes := []domain.Video{{
		PublisherID: uuid.Must(uuid.Parse(pubIDStr)),
		Topic:       expTop,
		Description: expDec,
	}}

	cases := []struct {
		name          string
		wantErr       bool
		pubID         string
		urlParams     string
		expStatusCode int
	}{
		// SearchPublisher
		{"ok search", false,
			pubIDStr, "?offset=0&limit=5&orderBy=createdAt&asc=t&query=search", http.StatusOK},

		{"few words search", false,
			pubIDStr, "?offset=0&limit=5&orderBy=createdAt&asc=t&query=search+lorem+ipsum+kitty+dolor", http.StatusOK},

		{"invalid offset search", true,
			pubIDStr, "?offset=A&limit=5&orderBy=createdAt&asc=t&query=search", http.StatusBadRequest},

		{"invalid limit search", true,
			pubIDStr, "?offset=0&limit=A&orderBy=createdAt&asc=t&query=search", http.StatusBadRequest},

		{"ivalid id", true,
			"1", "?offset=0&limit=5&orderBy=createdAt&asc=t&query=search", http.StatusBadRequest},

		{"empty id", true,
			"", "?offset=0&limit=5&orderBy=createdAt&asc=t&query=search", http.StatusNotFound},

		//TODO: temp fix, till error handling will be implemented
		{"invalid search", true,
			pubIDStr, "?offset=0&limit=5&orderBy=createdAt&asc=t&query=,.<>?/;", http.StatusBadRequest},
	}

	t.Parallel()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			s := mock_app.NewMockVideoService(ctrl)
			h := httpadapter.NewVideoHandler(s, idgen.New(), log.New(io.Discard, "", 0))
			r := mux.NewRouter()
			httpadapter.SetupRouter(r, h)

			if c.wantErr {
				s.EXPECT().
					SearchPublisher(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					MaxTimes(0)
			} else {
				s.EXPECT().SearchPublisher(
					gomock.Any(),
					gomock.Eq(uuid.MustParse(c.pubID)),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(
					expectedRes, nil,
				).MaxTimes(1)
			}

			url := strings.Replace(
				httpadapter.RoutePublisherVideos,
				"{"+httpadapter.PathVarPublisherID+"}",
				c.pubID,
				1,
			) + c.urlParams

			req := httptest.NewRequest(http.MethodGet, url, nil)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			require.Equal(t, c.expStatusCode, rec.Code)

		})
	}
}
