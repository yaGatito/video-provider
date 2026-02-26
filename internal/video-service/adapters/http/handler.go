package httpadapter

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"video-provider/internal/video-service/app"
	"video-provider/internal/video-service/domain"
	"video-provider/internal/video-service/policy"
	"video-provider/internal/video-service/ports"

	"github.com/gorilla/mux"
)

const (
	SearchUrlParam  = "query"
	LimitUrlParam   = "limit"
	OffsetUrlParam  = "offset"
	OrderByUrlParam = "orderBy"
	IsAscUrlParam   = "asc"
)

// VideoHandler handles HTTP requests for video operations.
// It provides endpoints for creating, retrieving, and searching videos.
type VideoHandler struct {
	VideoInteractor app.VideoService
	IDGen           ports.IDGen
	log             *log.Logger
}

// NewVideoHandler creates and returns a new VideoHandler instance.
// Parameters:
//   - userInteractor: service for video operations
//   - idGen: ID generator for parsing UUIDs
//   - log: logger instance for recording events
//
// Returns a configured VideoHandler ready to handle HTTP requests.
func NewVideoHandler(
	userInteractor app.VideoService,
	idGen ports.IDGen,
	log *log.Logger,
) VideoHandler {
	return VideoHandler{VideoInteractor: userInteractor, IDGen: idGen, log: log}
}

// Create godoc
// @Summary      Creates new video. Publisher id example: d9fa522f-0006-464f-8d68-356ba1d6ad7d
// @Description  Creates a new video record for the specified publisher
// @Tags         videos
// @Accept       json
// @Produce      json
// @Param        publisherID  path      string                 true  "publisher ID (UUID)"
// @Param        video        body      createVideoRequestBody true  "Video creation request body"
// @Success      200          {object}  nil
// @Failure      400          {object}  string "Invalid input"
// @Failure      500          {object}  string "Internal error"
// @Router       /v1/videos/pub/{publisherID} [post]
func (h *VideoHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Required path variable
	publisherID, err := h.pathVarHandler(w, r, PathVarPublisherID)
	if err != nil {
		return
	}

	var createVideoRequestData createVideoRequestBody
	if err := json.NewDecoder(r.Body).Decode(&createVideoRequestData); !errors.Is(err, io.EOF) && err != nil {
		h.writeResponse(w, nil, fmt.Errorf("error decoding request body: %w", err), http.StatusBadRequest)
		return
	}
	if err := createVideoRequestData.validate(); err != nil {
		h.writeResponse(w, nil, fmt.Errorf("error validating request body: %w", err), http.StatusBadRequest)
		return
	}

	video, err := h.VideoInteractor.Create(r.Context(), domain.Video{
		PublisherID: publisherID,
		Topic:       createVideoRequestData.Topic,
		Description: createVideoRequestData.Description,
	})
	h.writeResponse(w, video, err, http.StatusInternalServerError)
}

// GetByID godoc
// @Summary      Get video by ID
// @Description  Returns details of a single video by its unique identifier
// @Tags         videos
// @Produce      json
// @Param        videoID  path      string  true  "video ID (UUID)"
// @Success      200      {object}  VideoResponseBody
// @Failure      400      {object}  string
// @Failure      500      {object}  string
// @Router       /v1/videos/id/{videoID} [get]
func (h *VideoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	videoID, err := h.pathVarHandler(w, r, PathVarVideoID)
	if err != nil {
		return
	}

	video, err := h.VideoInteractor.GetByID(r.Context(), domain.UUID(videoID))
	h.writeResponse(w, video, err, http.StatusInternalServerError)
}

// GetByPublisher godoc
// @Summary      Get videos by publisher
// @Description  Returns a list of videos for a specific publisher with pagination and search support
// @Tags         videos
// @Produce      json
// @Param        publisherID  path      string  true   "publisher ID (UUID)"
// @Param        query        query     string  false  "Search query"
// @Param        limit        query     int     false  "Limit (default 10)"
// @Param        offset       query     int     false  "Offset (default 0)"
// @Success      200          {array}   VideoResponseBody
// @Router       /v1/videos/pub/{publisherID} [get]
func (h *VideoHandler) GetByPublisher(w http.ResponseWriter, r *http.Request) {
	// Required path variable
	publisherID, err := h.pathVarHandler(w, r, PathVarPublisherID)
	if err != nil {
		h.writeResponse(w, nil, fmt.Errorf("failed to parse `publisherId` path var: %w", err), http.StatusBadRequest)
		return
	}

	urlValues, err := h.parseUrlValues(w, r.URL.RawQuery)
	if err != nil {
		return
	}

	resInt, err := h.parseIntsUrlParams(w, urlValues, LimitUrlParam, OffsetUrlParam)
	if err != nil {
		return
	}
	limit := resInt[0]
	offset := resInt[1]

	resStr, err := h.parseStringsUrlParams(w, urlValues, OrderByUrlParam, IsAscUrlParam)
	if err != nil {
		return
	}
	orderBy := resStr[0]
	asc := resStr[1]

	search, err := h.extractUrlVarString(urlValues, SearchUrlParam)
	if err != nil {
		// TODO: Check `empty value` error here and skip, in other case -> write error response
		search = ""
	}

	var videos []domain.Video
	if search == "" {
		videos, err = h.VideoInteractor.GetByPublisher(r.Context(), publisherID, orderBy, offset, limit, asc)
	} else {
		videos, err = h.VideoInteractor.SearchPublisher(r.Context(), publisherID, search, orderBy, offset, limit, asc)
	}
	h.writeResponse(w, videos, err, http.StatusInternalServerError)
}

// SearchGlobal godoc
// @Summary      Global search
// @Description  Search videos in the entire database by a keyword
// @Tags         videos
// @Produce      json
// @Param        query   query     string  true   "Search query"
// @Param        limit   query     int     false  "Limit (default 10)"
// @Param        offset  query     int     false  "Offset (default 0)"
// @Param        sort    query     string  false  "Sort (default 'createdAt')"
// @Param        order   query     string  false  "Order (default 'asc')"
// @Success      200     {array}   VideoResponseBody
// @Router       /v1/videos/search/ [get]
func (h *VideoHandler) SearchGlobal(w http.ResponseWriter, r *http.Request) {
	// Url int parameters
	urlValues, err := h.parseUrlValues(w, r.URL.RawQuery)
	if err != nil {
		return
	}

	resInt, err := h.parseIntsUrlParams(w, urlValues, LimitUrlParam, OffsetUrlParam)
	if err != nil {
		return
	}
	limit := resInt[0]
	offset := resInt[1]

	resStr, err := h.parseStringsUrlParams(w, urlValues, SearchUrlParam, OrderByUrlParam, IsAscUrlParam)
	if err != nil {
		return
	}
	search := resStr[0]
	orderBy := resStr[1]
	asc := resStr[2]

	videos, err := h.VideoInteractor.SearchGlobal(r.Context(), search, orderBy, offset, limit, asc)
	h.writeResponse(w, videos, err, http.StatusInternalServerError)
}

func (h VideoHandler) writeResponse(w http.ResponseWriter, v any, err error, errStatusCode int) {
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(errStatusCode)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	switch val := v.(type) {
	case domain.Video:
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(h.toDtoVideo(val))
		if err != nil {
			h.log.Println("Error encoding response body:", err)
		}

	case []domain.Video:
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(h.toDtoVideos(val))
		if err != nil {
			h.log.Println("Error encoding response body:", err)
		}
	}
}

func (h VideoHandler) writeJSON(w http.ResponseWriter, status int, v any) {
}

func (h VideoHandler) parseUrlValues(w http.ResponseWriter, query string) (url.Values, error) {
	if len(query) > policy.UrlMaxLen {
		err := fmt.Errorf("too large query")
		h.writeResponse(w, nil, err, http.StatusBadRequest)
		return nil, err
	}
	urlValues, err := url.ParseQuery(query)
	if err != nil {
		h.writeResponse(w, nil, fmt.Errorf("failed to parse params: %w", err), http.StatusBadRequest)
		return nil, err
	}
	return urlValues, nil
}

func (h VideoHandler) parseIntsUrlParams(
	w http.ResponseWriter,
	values url.Values,
	params ...string,
) ([]int32, error) {
	res := make([]int32, len(params))

	for i, param := range params {
		val, err := strconv.ParseInt(values.Get(param), 10, 32)
		if err != nil {
			h.writeResponse(w, nil, fmt.Errorf("failed to parse query param `%s`: %w", param, err), http.StatusBadRequest)
			return nil, err
		}
		res[i] = int32(val)
	}

	return res, nil
}

func (h VideoHandler) parseStringsUrlParams(
	w http.ResponseWriter,
	values url.Values,
	params ...string,
) ([]string, error) {
	res := make([]string, len(params))

	for i, param := range params {
		val, err := h.extractUrlVarString(values, param)
		if err != nil {
			h.writeResponse(w, nil, fmt.Errorf("failed to parse param `%s`: %w", param, err), http.StatusBadRequest)
			return nil, err
		}
		res[i] = val
	}

	return res, nil
}

func (h VideoHandler) pathVarHandler(
	w http.ResponseWriter,
	r *http.Request,
	varName string,
) (domain.UUID, error) {
	val, ok := mux.Vars(r)[varName]
	if !ok {
		err := fmt.Errorf("path var `%s` not found", varName)
		h.writeResponse(w, nil, err, http.StatusBadRequest)
		return domain.UUID{}, err
	}
	res, err := h.IDGen.Parse(val)
	if err != nil {
		h.writeResponse(w, nil, fmt.Errorf("failed to parse path var `%s`: %w", varName, err), http.StatusBadRequest)
		return domain.UUID{}, fmt.Errorf("failed to parse %s: %w", varName, err)
	}

	return res, nil
}

func (h VideoHandler) extractUrlVarString(
	values url.Values,
	paramName string,
) (string, error) {

	value := values.Get(paramName)
	if len(value) == 0 {
		return "", fmt.Errorf("%s empty", paramName)
	}
	value, err := url.QueryUnescape(value)
	if err != nil {
		return "", fmt.Errorf("failed to unescape %s: %s; err: %w", paramName, value, err)
	}

	return value, nil
}

func (h VideoHandler) toDtoVideo(v domain.Video) VideoResponseBody {
	return VideoResponseBody{
		ID:          v.ID.String(),
		PublisherID: v.PublisherID.String(),
		Topic:       v.Topic,
		Description: v.Description,
		CreatedAt:   v.CreatedAt,
	}
}

func (h VideoHandler) toDtoVideos(videos []domain.Video) []VideoResponseBody {
	res := make([]VideoResponseBody, len(videos))
	for i, v := range videos {
		res[i] = h.toDtoVideo(v)
	}
	return res
}
