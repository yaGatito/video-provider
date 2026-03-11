package httpadp

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
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
	pubID := uuid.Must(uuid.NewRandom())
	testVideo := domain.Video{
		PublisherID: pubID,
		Topic:       "TEST",
		Description: "TESTEE",
		CreatedAt:   time.Now(),
		Status:      domain.StatusPublished,
	}
	expVideoRes := VideoResponseBody{
		ID:          testVideo.ID.String(),
		PublisherID: testVideo.PublisherID.String(),
		Topic:       testVideo.Topic,
		Description: testVideo.Description,
		CreatedAt:   testVideo.CreatedAt.Format(time.DateTime),
	}

	cases := []struct {
		name          string
		expCallCnt    int
		expStatusCode int
		pubID         string
		reqBody       string
	}{
		{"ok", 1, http.StatusOK, pubID.String(),
			`{"topic":"TESTE","description":"TESTEE"}`},

		{"ivalid id", 0, http.StatusBadRequest, "1",
			`{"topic":"TEST","description":"TESTEE"}`},

		{"ivalid req body", 0, http.StatusBadRequest, pubID.String(),
			`{"topic":"TEST","description":"TESTEE"}1`},

		{"ivalid req body 1", 0, http.StatusBadRequest, pubID.String(),
			`{"topic":"TEST","encryption":"TESTEE"}`},

		{"ivalid req body 2", 0, http.StatusBadRequest, pubID.String(),
			`{"tropics":"TEST","description":"TESTEE"}`},
	}

	t.Parallel()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			s := mock_app.NewMockVideoService(ctrl)
			h := NewVideoHandler(s, idgen.New(), log.New(io.Discard, "", 0))
			r := mux.NewRouter()
			SetupRouter(r, h)

			s.EXPECT().Create(gomock.Any(), gomock.Any()).Return(testVideo, nil).MaxTimes(c.expCallCnt)

			req := httptest.NewRequest(
				http.MethodPost,
				strings.Replace(
					RoutePublisherVideos,
					"{"+PathVarPublisherID+"}",
					c.pubID,
					1,
				),
				strings.NewReader(c.reqBody),
			)
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			require.Equal(t, c.expStatusCode, rec.Code)

			switch c.expCallCnt {
			case 1:
				var actualRes VideoResponseBody
				err := json.Unmarshal(rec.Body.Bytes(), &actualRes)
				require.NoError(t, err)

				require.Equal(t, expVideoRes, actualRes)
			}

		})
	}
}

func TestGetVideoById(t *testing.T) {
	ctrl := gomock.NewController(t)
	s := mock_app.NewMockVideoService(ctrl)
	h := NewVideoHandler(s, idgen.New(), log.New(io.Discard, "", 0))
	r := mux.NewRouter()
	SetupRouter(r, h)

	vidID := uuid.Must(uuid.NewRandom())

	video := domain.Video{
		PublisherID: uuid.Must(uuid.NewRandom()),
		Topic:       "topic",
		Description: "description",
	}

	cases := []struct {
		name          string
		reqVidID      string
		expCallCnt    int
		expStatusCode int
	}{
		{"ok", vidID.String(), 1, http.StatusOK},
		{"ivalid id", "1", 0, http.StatusBadRequest},
		{"empty id", "", 0, http.StatusNotFound},
	}

	t.Parallel()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(video, nil).MaxTimes(c.expCallCnt)

			req := httptest.NewRequest(
				http.MethodGet,
				strings.Replace(
					RouteVideo,
					"{"+PathVarVideoID+"}",
					c.reqVidID,
					1,
				),
				nil,
			)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			switch c.expCallCnt {
			case 1:
				var actualRes domain.Video
				json.Unmarshal(rec.Body.Bytes(), &actualRes)
				require.Equal(t, video, actualRes)
			}

			require.Equal(t, c.expStatusCode, rec.Code)
		})
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
		name       string
		expCallCnt int
		pubID      string
		urlParams  string
	}{
		{"ok", 1,
			pubIDStr, "?offset=0&limit=5&orderBy=createdAt&asc=t"},

		{"large URL query", 0,
			pubIDStr, "?offset=5&limit=5&orderBy=createdAt&asc=" + strings.Repeat("a", policy.UrlMaxLen)},

		{"invalid offset", 0,
			pubIDStr, "?offset=H&limit=5&orderBy=createdAt&asc=t"},

		{"no offset in URL", 0,
			pubIDStr, "?limit=5&orderBy=createdAt&asc=t"},

		{"invalid limit", 0,
			pubIDStr, "?offset=0&limit=H&orderBy=createdAt&asc=t"},

		{"no limit in URL", 0,
			pubIDStr, "?offset=0&orderBy=createdAt&asc=t"},

		{"invalid asc", 0,
			pubIDStr, "?offset=0&limit=5&orderBy=createdAt&asc=!@#%"},

		{"no asc in URL", 0,
			pubIDStr, "?offset=0&limit=5&orderBy=createdAt"},

		{"invalid orderBy", 0,
			pubIDStr, "?offset=0&limit=5&orderBy=!@#%&asc=t"},

		{"no orderBy in URL", 0,
			pubIDStr, "?offset=0&limit=5&asc=t"},

		{"wrong url params", 0,
			pubIDStr, "?affset=0&leemeet=5&orderBy=createdAt&asc=t&kueree=search"},

		{"partial wrong url params", 0,
			pubIDStr, "?affset=0&leemeet=5&orderBy=createdAt&asc=t&search="},
	}

	t.Parallel()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			s := mock_app.NewMockVideoService(ctrl)
			h := NewVideoHandler(s, idgen.New(), log.New(io.Discard, "", 0))
			r := mux.NewRouter()
			SetupRouter(r, h)

			s.EXPECT().GetByPublisher(
				gomock.Any(),
				gomock.Eq(uuid.MustParse(c.pubID)),
				gomock.Any(),
				gomock.Any(),
				gomock.Any(),
				gomock.Any(),
			).Return(
				expectedRes, nil,
			).MaxTimes(c.expCallCnt)

			url := strings.Replace(
				RoutePublisherVideos,
				"{"+PathVarPublisherID+"}",
				c.pubID,
				1,
			) + c.urlParams

			req := httptest.NewRequest(http.MethodGet, url, nil)
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			switch c.expCallCnt {
			case 1:
				require.Equal(t, http.StatusOK, rec.Code)
				var actualRes []domain.Video
				json.Unmarshal(rec.Body.Bytes(), &actualRes)
				require.Equal(t, expectedRes, actualRes)
			case 0:
				require.Equal(t, http.StatusBadRequest, rec.Code)
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
		expCallCnt    int
		pubID         string
		urlParams     string
		expStatusCode int
	}{
		// SearchPublisher
		{"ok search", 1,
			pubIDStr, "?offset=0&limit=5&orderBy=createdAt&asc=t&query=search", http.StatusOK},

		{"few words search", 1,
			pubIDStr, "?offset=0&limit=5&orderBy=createdAt&asc=t&query=search+lorem+ipsum+kitty+dolor", http.StatusOK},

		{"invalid offset search", 0,
			pubIDStr, "?offset=A&limit=5&orderBy=createdAt&asc=t&query=search", http.StatusBadRequest},

		{"invalid limit search", 0,
			pubIDStr, "?offset=0&limit=A&orderBy=createdAt&asc=t&query=search", http.StatusBadRequest},

		{"ivalid id", 0,
			"1", "?offset=0&limit=5&orderBy=createdAt&asc=t&query=search", http.StatusBadRequest},

		{"empty id", 0,
			"", "?offset=0&limit=5&orderBy=createdAt&asc=t&query=search", http.StatusNotFound},

		//TODO: temp fix, till error handling will be implemented
		{"invalid search", 0,
			pubIDStr, "?offset=0&limit=5&orderBy=createdAt&asc=t&query=,.<>?/;", http.StatusBadRequest},
	}

	t.Parallel()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			s := mock_app.NewMockVideoService(ctrl)
			h := NewVideoHandler(s, idgen.New(), log.New(io.Discard, "", 0))
			r := mux.NewRouter()
			SetupRouter(r, h)

			s.EXPECT().SearchPublisher(
				gomock.Any(),
				gomock.Any(),
				gomock.Any(),
				gomock.Any(),
				gomock.Any(),
				gomock.Any(),
				gomock.Any()).Return(expectedRes, nil).MaxTimes(c.expCallCnt)

			url := strings.Replace(
				RoutePublisherVideos,
				"{"+PathVarPublisherID+"}",
				c.pubID,
				1,
			) + c.urlParams

			req := httptest.NewRequest(http.MethodGet, url, nil)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			require.Equal(t, c.expStatusCode, rec.Code)
			if c.expCallCnt == 1 {
				var actualRes []domain.Video
				json.Unmarshal(rec.Body.Bytes(), &actualRes)
				require.Equal(t, expectedRes, actualRes)
			}
		})
	}
}
