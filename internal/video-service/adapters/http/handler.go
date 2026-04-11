package httpadp

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"pkg/shared"
	"strconv"
	"video-service/app"
	"video-service/domain"
	"video-service/policy"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/yaGatito/slicex"
)

const (
	SearchUrlParam = "query"
	LimitUrlParam  = "limit"
	OffsetUrlParam = "offset"
	SortByUrlParam = "sort"
	IsAscUrlParam  = "order"
)

// VideoHandler handles HTTP requests related to video operations.
type VideoHandler struct {
	videoInteractor app.VideoService
	log             *log.Logger
	validate        *validator.Validate
}

// NewVideoHandler creates a new VideoHandler.
func NewVideoHandler(
	userInteractor app.VideoService,
	log *log.Logger,
	validate *validator.Validate,
) *VideoHandler {
	return &VideoHandler{videoInteractor: userInteractor, log: log, validate: validate}
}

// Create godoc
//
//	@Summary		Creates new video.
//	@Description	Creates a new video record for the specified publisher.
//	@Tags			videos
//	@Accept			json
//	@Produce		json
//	@Param 			Authorization 	header 		string 					true "JWT token for authentication (e.g., Bearer <token>)"
//	@Param			publisherID		path		string					true	"Publisher ID (UUID)"
//	@Param			video			body		createVideoRequestBody	true	"Video creation request body"
//	@Success		201			{object}	nil
//	@Failure		400			{object}	string	"Invalid input"
//	@Failure		500			{object}	string	"Internal error"
//	@Router			/v1/videos/pub/{publisherID} [post]
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

	video, err := h.videoInteractor.Create(r.Context(), domain.Video{
		PublisherID: publisherID,
		Topic:       createVideoRequestData.Topic,
		Description: createVideoRequestData.Description,
	})
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeResponse(w, DtoVideo(video), http.StatusCreated)
}

// GetByID godoc
//
//	@Summary		Get video by ID
//	@Description	Returns details of a single video by its unique identifier
//	@Tags			videos
//	@Produce		json
//	@Param 			Authorization 		header 	string 	true "JWT token for authentication (e.g., Bearer <token>)"
//	@Param			videoID	path		string	true	"video ID (UUID)"	Format(uuid)
//	@Success		200		{object}	VideoResponseBody
//	@Failure		400		{object}	string	"Invalid video ID format"
//	@Failure		500		{object}	string	"Internal server error"
//	@Router			/v1/videos/id/{videoID} [get]
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

	video, err := h.videoInteractor.GetByID(r.Context(), videoID)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeResponse(w, DtoVideo(video), http.StatusOK)
}

// GetByPublisher godoc
//
//	@Summary		Get videos by publisher
//	@Description	Returns a list of videos for a specific publisher with pagination and search support
//	@Tags			videos
//	@Produce		json
//	@Param 			Authorization 		header 	string 	true "JWT token for authentication (e.g., Bearer <token>)"
//	@Param			publisherID	path	string	true	"publisher ID (UUID)"
//	@Param			limit		query	int		false	"Limit (example: 10)"
//	@Param			offset		query	int		false	"Offset (example: 0)"
//	@Param			sort		query	string	false	"Sort (example: `date`)"
//	@Param			order		query	string	false	"Order (`t` for ascending, `f` for descending)"
//	@Success		200			{array}	VideoResponseBody
//	@Router			/v1/videos/pub/{publisherID} [get]
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

	resStr, err := h.parseStringsUrlParams(urlValues, SortByUrlParam, IsAscUrlParam)
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

	pageParams, err := newVideoPageParams(orderBy, offset, limit, asc)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	var videos []domain.Video
	if search == "" {
		videos, err = h.videoInteractor.GetByPublisher(r.Context(), publisherID, pageParams)
		if err != nil {
			h.writeErrorResponse(w, err)
			return
		}
	} else {
		search, err = ValidateSearchQuery(search)
		if err != nil {
			h.writeErrorResponse(w, err)
			return
		}

		videos, err = h.videoInteractor.SearchPublisher(r.Context(), publisherID, search, pageParams)
		if err != nil {
			h.writeErrorResponse(w, err)
			return
		}
	}
	h.writeResponse(w, VideosResponseBody{Videos: slicex.Map(videos, DtoVideo)}, http.StatusOK)
}

// SearchGlobal godoc
//
//	@Summary		Global search
//	@Description	Search videos in the entire database by a keyword
//	@Tags			videos
//	@Produce		json
//	@Param			query	query	string	true	"Search query"
//	@Param			limit	query	int		false	"Limit (example: 10)"
//	@Param			offset	query	int		false	"Offset (example: 0)"
//	@Param			sort	query	string	false	"Sort (example: `date`)"
//	@Param			order	query	string	false	"Order (`t` for ascending, `f` for descending)"
//	@Success		200		{array}	VideoResponseBody
//	@Router			/v1/videos/search/ [get]
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

	resStr, err := h.parseStringsUrlParams(urlValues, SearchUrlParam, SortByUrlParam, IsAscUrlParam)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}
	search := resStr[0]
	orderBy := resStr[1]
	asc := resStr[2]

	pageParams, err := newVideoPageParams(orderBy, offset, limit, asc)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	videos, err := h.videoInteractor.SearchGlobal(r.Context(), search, pageParams)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeResponse(w, VideosResponseBody{Videos: slicex.Map(videos, DtoVideo)}, http.StatusOK)
}

// parseUrlValues parses URL query parameters.
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

// parseIntsUrlParams parses integer parameters from the URL query.
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

// parseStringsUrlParams parses string parameters from the URL query.
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

// pathVarHandler extracts a path variable and parses it as a UUID.
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

// extractUrlVarString extracts and unescapes a string parameter from the URL query.
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

// writeResponse writes a JSON response with the specified HTTP status code.
func (h *VideoHandler) writeResponse(w http.ResponseWriter, v any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		h.log.Println("Error encoding response body:", err)
	}
}

// writeErrorResponse writes an error response in JSON format.
func (h *VideoHandler) writeErrorResponse(w http.ResponseWriter, vErr error) {
	w.Header().Set("Content-Type", "application/json")

	switch vErr := vErr.(type) {
	case shared.Error:
		h.log.Printf("Error: %s\n", vErr.Message)
		if vErr.Err != nil {
			h.log.Printf("Details: %s\n", vErr.Err.Error())
		}

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
			Message: "video-provider error",
		})
		if err != nil {
			h.log.Println("Error encoding error response body:", err)
		}
	}
}
