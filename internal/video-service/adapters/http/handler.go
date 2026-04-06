package httpadp

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"video-provider/internal/pkg/shared"
	"video-provider/internal/video-service/app"
	"video-provider/internal/video-service/domain"
	"video-provider/internal/video-service/policy"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/yaGatito/slicex"
)

const (
	SearchUrlParam  = "query"
	LimitUrlParam   = "limit"
	OffsetUrlParam  = "offset"
	OrderByUrlParam = "order"
	IsAscUrlParam   = "asc"
)

type VideoHandler struct {
	VideoInteractor app.VideoService
	log             *log.Logger
	validate        *validator.Validate
}

func NewVideoHandler(
	userInteractor app.VideoService,
	log *log.Logger,
) VideoHandler {
	return VideoHandler{VideoInteractor: userInteractor, log: log, validate: NewVideoValidator()}
}

// Create godoc
// @Summary      Creates new video.
// @Description  Creates a new video record for the specified publisher.
// @Tags         videos
// @Accept       json
// @Produce      json
// @Param        publisherID  path      string                 true  "Publisher ID (UUID)"
// @Param        video        body      createVideoRequestBody true  "Video creation request body"
// @Success      201          {object}  nil
// @Failure      400          {object}  string "Invalid input"
// @Failure      500          {object}  string "Internal error"
// @Router       /v1/videos/pub/{publisherID} [post]
func (h *VideoHandler) Create(w http.ResponseWriter, r *http.Request) {
	publisherID, err := h.pathVarHandler(r, PathVarPublisherID)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	var createVideoRequestData createVideoRequestBody
	if err := json.NewDecoder(r.Body).Decode(&createVideoRequestData); err != nil {
		h.writeErrorResponse(w, shared.NewError(
			http.StatusBadRequest, "failed to decode request body", err))
		return
	}

	err = h.validate.Struct(createVideoRequestData)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	video, err := h.VideoInteractor.Create(r.Context(), domain.Video{
		PublisherID: publisherID,
		Topic:       createVideoRequestData.Topic,
		Description: createVideoRequestData.Description,
	})

	h.writeResponse(w, dtoVideo(video), http.StatusCreated)
}

// GetByID godoc
// @Summary      Get video by ID
// @Description  Returns details of a single video by its unique identifier
// @Tags         videos
// @Produce      json
// @Param        videoID  path      string  true  "video ID (UUID)"  Format(uuid)
// @Success      200      {object}  videoResponseBody
// @Failure      400      {object}  string  "Invalid video ID format"
// @Failure      500      {object}  string  "Internal server error"
// @Router       /v1/videos/id/{videoID} [get]
func (h *VideoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	videoID, err := h.pathVarHandler(r, PathVarVideoID)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	if videoID == uuid.Nil {
		h.writeErrorResponse(w, shared.NewError(http.StatusBadRequest, "empty video ID", nil))
		return
	}

	video, err := h.VideoInteractor.GetByID(r.Context(), domain.UUID(videoID))
	h.writeResponse(w, dtoVideo(video), http.StatusOK)
}

// GetByPublisher godoc
// @Summary      Get videos by publisher
// @Description  Returns a list of videos for a specific publisher with pagination and search support
// @Tags         videos
// @Produce      json
// @Param        publisherID  	path      string  true   "publisher ID (UUID)"
// @Param        limit   	  	query     int     false  "Limit (example: 10)"
// @Param        offset  		query     int     false  "Offset (example: 0)"
// @Param        sort    		query     string  false  "Sort (example: `date`)"
// @Param        order   		query     string  false  "Order (asc or desc, example: `t` for ascending, `f` for descending)"
// @Success      200          {array}   videoResponseBody
// @Router       /v1/videos/pub/{publisherID} [get]
func (h *VideoHandler) GetByPublisher(w http.ResponseWriter, r *http.Request) {
	// Required path variable
	publisherID, err := h.pathVarHandler(r, PathVarPublisherID)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	if publisherID == uuid.Nil {
		h.writeErrorResponse(w, shared.NewError(
			http.StatusBadRequest, "empty publisher ID", nil),
		)
		return
	}

	urlValues, err := h.parseUrlValues(r.URL.RawQuery)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	resInt, err := h.parseIntsUrlParams(urlValues, LimitUrlParam, OffsetUrlParam)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}
	limit := resInt[0]
	offset := resInt[1]

	resStr, err := h.parseStringsUrlParams(urlValues, OrderByUrlParam, IsAscUrlParam)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}
	orderBy := resStr[0]
	asc := resStr[1]

	search, err := h.extractUrlVarString(urlValues, SearchUrlParam)
	if err != nil && !errors.Is(err, shared.ErrEmptyValue) {
		h.writeErrorResponse(w, err)
		return
	}

	pageParams, err := NewVideoPageParams(orderBy, offset, limit, asc)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	var videos []domain.Video
	if search == "" {
		videos, err = h.VideoInteractor.GetByPublisher(r.Context(), publisherID, pageParams)
		if err != nil {
			h.writeErrorResponse(w, err)
			return
		}
	} else {
		search, err = validateSearchQuery(search)
		if err != nil {
			h.writeErrorResponse(w, err)
			return
		}

		videos, err = h.VideoInteractor.SearchPublisher(r.Context(), publisherID, search, pageParams)
		if err != nil {
			h.writeErrorResponse(w, err)
			return
		}
	}
	h.writeResponse(w, videosResponseBody{Videos: slicex.Map(videos, dtoVideo)}, http.StatusOK)
}

// SearchGlobal godoc
// @Summary      Global search
// @Description  Search videos in the entire database by a keyword
// @Tags         videos
// @Produce      json
// @Param        query   query     string  true   "Search query"
// @Param        limit   query     int     false  "Limit (example: 10)"
// @Param        offset  query     int     false  "Offset (example: 0)"
// @Param        sort    query     string  false  "Sort (example: `date`)"
// @Param        order   query     string  false  "Order (asc or desc, example: `t` for ascending, `f` for descending)"
// @Success      200     {array}   videoResponseBody
// @Router       /v1/videos/search/ [get]
func (h *VideoHandler) SearchGlobal(w http.ResponseWriter, r *http.Request) {
	// Url int parameters
	urlValues, err := h.parseUrlValues(r.URL.RawQuery)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	resInt, err := h.parseIntsUrlParams(urlValues, LimitUrlParam, OffsetUrlParam)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}
	limit := resInt[0]
	offset := resInt[1]

	resStr, err := h.parseStringsUrlParams(urlValues, SearchUrlParam, OrderByUrlParam, IsAscUrlParam)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}
	search := resStr[0]
	orderBy := resStr[1]
	asc := resStr[2]

	pageParams, err := NewVideoPageParams(orderBy, offset, limit, asc)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	videos, err := h.VideoInteractor.SearchGlobal(r.Context(), search, pageParams)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeResponse(w, videosResponseBody{Videos: slicex.Map(videos, dtoVideo)}, http.StatusOK)
}

func (h *VideoHandler) parseUrlValues(query string) (url.Values, error) {
	if len(query) > policy.UrlMaxLen {
		return nil, shared.NewError(http.StatusBadRequest, "too large url", nil)
	}
	urlValues, err := url.ParseQuery(query)
	if err != nil {
		h.log.Println(err)
		return nil, shared.NewError(
			http.StatusBadRequest, "unparsable url query", err)
	}
	return urlValues, nil
}

func (h *VideoHandler) parseIntsUrlParams(
	values url.Values,
	params ...string,
) ([]int32, error) {
	res := make([]int32, len(params))

	for i, param := range params {
		val, err := strconv.ParseInt(values.Get(param), 10, 32)
		if err != nil {
			return nil, shared.NewError(
				http.StatusBadRequest, "unparsable url param (int): "+param, err,
			)
		}
		res[i] = int32(val)
	}

	return res, nil
}

func (h *VideoHandler) parseStringsUrlParams(
	values url.Values,
	params ...string,
) ([]string, error) {
	res := make([]string, len(params))

	for i, param := range params {
		val, err := h.extractUrlVarString(values, param)
		if err != nil {
			return nil, shared.NewError(
				http.StatusBadRequest, "unparsable url param (string): "+param, err)
		}
		res[i] = val
	}

	return res, nil
}

func (h *VideoHandler) pathVarHandler(
	r *http.Request,
	varName string,
) (domain.UUID, error) {
	val, ok := mux.Vars(r)[varName]
	if !ok {
		return domain.UUID{}, shared.NewError(
			http.StatusBadRequest, "path var not specified: "+varName, nil)
	}
	res, err := uuid.Parse(val)
	if err != nil {
		return domain.UUID{}, shared.NewError(
			http.StatusBadRequest, "unparsable ID: "+varName, err)
	}

	return res, nil
}

func (h *VideoHandler) extractUrlVarString(
	values url.Values,
	paramName string,
) (string, error) {

	value := values.Get(paramName)
	if len(value) == 0 {
		return "", shared.ErrEmptyValue
	}
	value, err := url.QueryUnescape(value)
	if err != nil {
		return "", shared.NewError(
			http.StatusBadRequest, "failed to unescape url param: "+paramName, err)
	}

	return value, nil
}

func (h *VideoHandler) writeResponse(w http.ResponseWriter, v any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		h.log.Println("Error encoding response body:", err)
	}
}

func (h *VideoHandler) writeErrorResponse(w http.ResponseWriter, vErr error) {
	w.Header().Set("Content-Type", "application/json")

	switch vErr := vErr.(type) {
	case shared.Error:
		h.log.Printf("Error: %s\n", vErr.Err.Error())

		w.WriteHeader(int(vErr.Code))
		err := json.NewEncoder(w).Encode(serviceErrorResponse{
			Message: vErr.Message,
		})
		if err != nil {
			h.log.Println("Error encoding error response body:", err)
		}

	case validator.ValidationErrors:
		h.log.Printf("Validation request body error: %s\n", vErr[0].Error())
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(serviceErrorResponse{
			Message: "invalid field: " + vErr[0].Field(),
		})
		if err != nil {
			h.log.Println("Error validating request body:", err)
		}

	case error:
		h.log.Printf("Fallback error: %s\n", vErr.Error())
		w.WriteHeader(http.StatusInternalServerError)
		err := json.NewEncoder(w).Encode(serviceErrorResponse{
			Message: "internal error",
		})
		if err != nil {
			h.log.Println("Error encoding error response body:", err)
		}
	}
}
