package httpadp_test

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	httpadp "video-provider/video-service/adapters/http"
	mock_app "video-provider/video-service/app/mock"
	"video-provider/video-service/domain"
	"video-provider/video-service/policy"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"github.com/yaGatito/slicex"
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
	expVideoRes := httpadp.VideoResponseBody{
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
		{"ok", 1, http.StatusCreated, pubID.String(),
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
			val, _ := httpadp.NewVideoValidator()
			h := httpadp.NewVideoHandler(s, log.New(io.Discard, "", 0), val)
			r := mux.NewRouter()
			mockMiddleware := func(next http.Handler) http.Handler {
				return next
			}
			httpadp.SetupRouter(r, h, mockMiddleware, mockMiddleware, mockMiddleware)

			s.EXPECT().
				Create(gomock.Any(), gomock.Any()).
				Return(testVideo, nil).
				MaxTimes(c.expCallCnt)

			req := httptest.NewRequest(
				http.MethodPost,
				strings.Replace(
					httpadp.RoutePublisherVideos,
					"{"+httpadp.PathVarPublisherID+"}",
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
				var actualRes httpadp.VideoResponseBody
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
	val, _ := httpadp.NewVideoValidator()
	h := httpadp.NewVideoHandler(s, log.New(io.Discard, "", 0), val)
	r := mux.NewRouter()
	mockMiddleware := func(next http.Handler) http.Handler {
		return next
	}
	httpadp.SetupRouter(r, h, mockMiddleware, mockMiddleware, mockMiddleware)

	vidID := uuid.Must(uuid.NewRandom())

	video := domain.Video{
		ID:          vidID,
		PublisherID: uuid.Must(uuid.NewRandom()),
		Topic:       "topic",
		Description: "description",
		CreatedAt:   time.Now(),
	}

	expectedResponseBody := httpadp.DtoVideo(video)

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
					httpadp.RouteVideo,
					"{"+httpadp.PathVarVideoID+"}",
					c.reqVidID,
					1,
				),
				nil,
			)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			if c.expCallCnt == 1 {
				var actualRes httpadp.VideoResponseBody
				err := json.Unmarshal(rec.Body.Bytes(), &actualRes)
				require.Equal(t, expectedResponseBody, actualRes)
				require.NoError(t, err)
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
		PublisherID: uuid.Must(uuid.NewRandom()),
		Topic:       expTop,
		Description: expDec,
		CreatedAt:   time.Now(),
	}}

	expectedResponseBody := httpadp.VideosResponseBody{
		Videos: slicex.Map(expectedRes, httpadp.DtoVideo),
	}

	cases := []struct {
		name       string
		expCallCnt int
		pubID      string
		urlParams  string
	}{
		{"ok", 1,
			pubIDStr, "?offset=0&limit=5&sort=date&order=t"},

		{"large URL query", 0,
			pubIDStr, "?offset=5&limit=5&sort=date&asc=" + strings.Repeat("a", policy.UrlMaxLen)},

		{"invalid offset", 0,
			pubIDStr, "?offset=H&limit=5&sort=date&order=t"},

		{"no offset in URL", 0,
			pubIDStr, "?limit=5&sort=date&order=t"},

		{"invalid limit", 0,
			pubIDStr, "?offset=0&limit=H&sort=date&order=t"},

		{"no limit in URL", 0,
			pubIDStr, "?offset=0&sort=date&order=t"},

		{"invalid asc", 0,
			pubIDStr, "?offset=0&limit=5&sort=date&asc=!@#%"},

		{"no asc in URL", 0,
			pubIDStr, "?offset=0&limit=5&sort=date"},

		{"invalid orderBy", 0,
			pubIDStr, "?offset=0&limit=5&order=!@#%&order=t"},

		{"no orderBy in URL", 0,
			pubIDStr, "?offset=0&limit=5&order=t"},

		{"wrong url params", 0,
			pubIDStr, "?affset=0&leemeet=5&sort=date&order=t&kueree=search"},

		{"partial wrong url params", 0,
			pubIDStr, "?affset=0&leemeet=5&sort=date&order=t&search="},
	}

	t.Parallel()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			s := mock_app.NewMockVideoService(ctrl)
			val, _ := httpadp.NewVideoValidator()
			h := httpadp.NewVideoHandler(s, log.New(io.Discard, "", 0), val)
			r := mux.NewRouter()
			mockMiddleware := func(next http.Handler) http.Handler {
				return next
			}
			httpadp.SetupRouter(r, h, mockMiddleware, mockMiddleware, mockMiddleware)

			s.EXPECT().GetByPublisher(
				gomock.Any(),
				gomock.Eq(uuid.MustParse(c.pubID)),
				gomock.Any(),
			).Return(
				expectedRes, nil,
			).MaxTimes(c.expCallCnt)

			url := strings.Replace(
				httpadp.RoutePublisherVideos,
				"{"+httpadp.PathVarPublisherID+"}",
				c.pubID,
				1,
			) + c.urlParams

			req := httptest.NewRequest(http.MethodGet, url, nil)
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			switch c.expCallCnt {
			case 1:
				require.Equal(t, http.StatusOK, rec.Code)
				var actualRes httpadp.VideosResponseBody
				json.Unmarshal(rec.Body.Bytes(), &actualRes)
				require.Equal(t, expectedResponseBody, actualRes)
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

	expectedResponseBody := httpadp.VideosResponseBody{
		Videos: slicex.Map(expectedRes, httpadp.DtoVideo),
	}

	cases := []struct {
		name          string
		expCallCnt    int
		pubID         string
		urlParams     string
		expStatusCode int
	}{
		// SearchPublisher
		{"ok search", 1,
			pubIDStr, "?offset=0&limit=5&sort=date&order=t&query=search", http.StatusOK},

		{"few words search", 1,
			pubIDStr, "?offset=0&limit=5&sort=date&order=t&query=search+lorem+ipsum+kitty+dolor", http.StatusOK},

		{"invalid offset search", 0,
			pubIDStr, "?offset=A&limit=5&sort=date&order=t&query=search", http.StatusBadRequest},

		{"invalid limit search", 0,
			pubIDStr, "?offset=0&limit=A&sort=date&order=t&query=search", http.StatusBadRequest},

		{"ivalid id", 0,
			"1", "?offset=0&limit=5&sort=date&order=t&query=search", http.StatusBadRequest},

		{"empty id", 0,
			"", "?offset=0&limit=5&sort=date&order=t&query=search", http.StatusNotFound},

		{"invalid search", 0,
			pubIDStr, "?offset=0&limit=5&sort=date&order=t&query=,.<>?/;", http.StatusBadRequest},
	}

	t.Parallel()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			s := mock_app.NewMockVideoService(ctrl)
			val, _ := httpadp.NewVideoValidator()
			h := httpadp.NewVideoHandler(s, log.New(io.Discard, "", 0), val)
			r := mux.NewRouter()
			mockMiddleware := func(next http.Handler) http.Handler {
				return next
			}
			httpadp.SetupRouter(r, h, mockMiddleware, mockMiddleware, mockMiddleware)

			s.EXPECT().SearchPublisher(
				gomock.Any(),
				gomock.Any(),
				gomock.Any(),
				gomock.Any()).Return(expectedRes, nil).MaxTimes(c.expCallCnt)

			url := strings.Replace(
				httpadp.RoutePublisherVideos,
				"{"+httpadp.PathVarPublisherID+"}",
				c.pubID,
				1,
			) + c.urlParams

			req := httptest.NewRequest(http.MethodGet, url, nil)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			require.Equal(t, c.expStatusCode, rec.Code)
			if c.expCallCnt == 1 {
				var actualRes httpadp.VideosResponseBody
				err := json.Unmarshal(rec.Body.Bytes(), &actualRes)
				require.Equal(t, expectedResponseBody, actualRes)
				require.NoError(t, err)
			}
		})
	}
}

func BenchmarkCreateVideo(b *testing.B) {
	pubID := uuid.Must(uuid.NewRandom())
	testVideo := domain.Video{
		PublisherID: pubID,
		Topic:       "TEST",
		Description: "TESTEE",
		CreatedAt:   time.Now(),
		Status:      domain.StatusPublished,
	}
	reqBody := `{"topic":"TESTE","description":"TESTEE"}`

	ctrl := gomock.NewController(b)
	s := mock_app.NewMockVideoService(ctrl)
	val, _ := httpadp.NewVideoValidator()
	h := httpadp.NewVideoHandler(s, httpadp.DefaultLogger, val)
	r := mux.NewRouter()
	mockMiddleware := func(next http.Handler) http.Handler {
		return next
	}
	httpadp.SetupRouter(r, h, mockMiddleware, mockMiddleware, mockMiddleware)

	s.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(testVideo, nil).
		MinTimes(b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(
			http.MethodPost,
			strings.Replace(
				httpadp.RoutePublisherVideos,
				"{"+httpadp.PathVarPublisherID+"}",
				pubID.String(),
				1,
			),
			bytes.NewBufferString(reqBody),
		)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)
	}
}

func BenchmarkGetVideoById(b *testing.B) {
	ctrl := gomock.NewController(b)
	s := mock_app.NewMockVideoService(ctrl)
	val, _ := httpadp.NewVideoValidator()
	h := httpadp.NewVideoHandler(s, httpadp.DefaultLogger, val)
	r := mux.NewRouter()
	mockMiddleware := func(next http.Handler) http.Handler {
		return next
	}
	httpadp.SetupRouter(r, h, mockMiddleware, mockMiddleware, mockMiddleware)

	vidID := uuid.Must(uuid.NewRandom())
	video := domain.Video{
		ID:          vidID,
		PublisherID: uuid.Must(uuid.NewRandom()),
		Topic:       "topic",
		Description: "description",
		CreatedAt:   time.Now(),
	}

	s.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(video, nil).MinTimes(b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(
			http.MethodGet,
			strings.Replace(
				httpadp.RouteVideo,
				"{"+httpadp.PathVarVideoID+"}",
				vidID.String(),
				1,
			),
			nil,
		)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
	}
}

func BenchmarkGetByPublisherVideos(b *testing.B) {
	pubIDStr := "d9fa522f-0016-464f-8d68-356ba1d6ad7d"
	urlParams := "?offset=0&limit=5&sort=date&order=t"

	ctrl := gomock.NewController(b)
	s := mock_app.NewMockVideoService(ctrl)
	val, _ := httpadp.NewVideoValidator()
	h := httpadp.NewVideoHandler(s, httpadp.DefaultLogger, val)
	r := mux.NewRouter()
	mockMiddleware := func(next http.Handler) http.Handler {
		return next
	}
	httpadp.SetupRouter(r, h, mockMiddleware, mockMiddleware, mockMiddleware)

	expectedRes := []domain.Video{{
		PublisherID: uuid.Must(uuid.NewRandom()),
		Topic:       "topic",
		Description: "description",
		CreatedAt:   time.Now(),
	}}

	s.EXPECT().GetByPublisher(
		gomock.Any(),
		gomock.Eq(uuid.MustParse(pubIDStr)),
		gomock.Any(),
	).Return(
		expectedRes, nil,
	).MinTimes(b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		url := strings.Replace(
			httpadp.RoutePublisherVideos,
			"{"+httpadp.PathVarPublisherID+"}",
			pubIDStr,
			1,
		) + urlParams

		req := httptest.NewRequest(http.MethodGet, url, nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
	}
}

func BenchmarkSearchGlobalVideos(b *testing.B) {
	urlParams := "?offset=0&limit=5&sort=date&order=t&query=search"

	ctrl := gomock.NewController(b)
	s := mock_app.NewMockVideoService(ctrl)
	val, _ := httpadp.NewVideoValidator()
	h := httpadp.NewVideoHandler(s, httpadp.DefaultLogger, val)
	r := mux.NewRouter()
	mockMiddleware := func(next http.Handler) http.Handler {
		return next
	}
	httpadp.SetupRouter(r, h, mockMiddleware, mockMiddleware, mockMiddleware)

	expectedRes := []domain.Video{{
		PublisherID: uuid.Must(uuid.NewRandom()),
		Topic:       "topic",
		Description: "description",
		CreatedAt:   time.Now(),
	}}

	s.EXPECT().SearchGlobal(
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).Return(
		expectedRes, nil,
	).MinTimes(b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, urlParams, nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
	}
}
