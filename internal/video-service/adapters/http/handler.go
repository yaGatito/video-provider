package httpadp

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"video-provider/pkg/common"
	"video-provider/video-service/app"
	"video-provider/video-service/domain"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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
	log             *common.Logger
	validate        *validator.Validate
}

var DefaultLogger = log.New(os.Stdout, "[VIDSVC]", log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC)

// NewVideoHandler creates a new VideoHandler.
func NewVideoHandler(
	userInteractor app.VideoService,
	log *common.Logger,
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
	publisherID, err := common.PathVarHandler(r, PathVarPublisherID)
	if err != nil {
		common.WriteErrorResponse(w, h.log, err)
		return
	}

	var createVideoRequestData createVideoRequestBody
	if err := json.NewDecoder(r.Body).Decode(&createVideoRequestData); err != nil {
		common.WriteErrorResponse(w, h.log, &common.Error{
			Err:     err,
			Code:    http.StatusBadRequest,
			Message: "failed to decode request body",
		})
		return
	}

	err = h.validate.Struct(createVideoRequestData)
	if err != nil {
		common.WriteErrorResponse(w, h.log, err)
		return
	}

	video, err := h.videoInteractor.Create(r.Context(), domain.Video{
		PublisherID: publisherID,
		Topic:       createVideoRequestData.Topic,
		Description: createVideoRequestData.Description,
	})
	if err != nil {
		common.WriteErrorResponse(w, h.log, err)
		return
	}

	common.WriteResponse(w, h.log, DtoVideo(video), http.StatusCreated)
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
	videoID, err := common.PathVarHandler(r, PathVarVideoID)
	if err != nil {
		common.WriteErrorResponse(w, h.log, err)
		return
	}

	if videoID == uuid.Nil {
		common.WriteErrorResponse(w, h.log,
			&common.Error{Code: http.StatusBadRequest, Message: "empty video ID"})
		return
	}

	video, err := h.videoInteractor.GetByID(r.Context(), videoID)
	if err != nil {
		common.WriteErrorResponse(w, h.log, err)
		return
	}

	common.WriteResponse(w, h.log, DtoVideo(video), http.StatusOK)
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
	publisherID, err := common.PathVarHandler(r, PathVarPublisherID)
	if err != nil {
		common.WriteErrorResponse(w, h.log, err)
		return
	}

	if publisherID == uuid.Nil {
		common.WriteErrorResponse(w, h.log, &common.Error{
			Code: http.StatusBadRequest, Message: "empty publisher ID",
		})
		return
	}

	urlValues, err := common.ParseUrlValues(r.URL.RawQuery)
	if err != nil {
		common.WriteErrorResponse(w, h.log, err)
		return
	}

	intParams, err := common.ParseIntsUrlParams(urlValues, LimitUrlParam, OffsetUrlParam)
	if err != nil {
		common.WriteErrorResponse(w, h.log, err)
		return
	}
	limit := intParams[0]
	offset := intParams[1]

	strParams, err := common.ParseStringsUrlParams(urlValues,
		SortByUrlParam,
		IsAscUrlParam,
		SearchUrlParam)

	var searchEmpty bool
	var search string
	switch err := err.(type) {
	case nil:
		// success
	case *common.Error:
		searchEmpty = SearchUrlParam == err.GetDetails()
		// if error is because of SearchUrlParam - we can go on with other data
		if !searchEmpty {
			common.WriteErrorResponse(w, h.log, err)
			return
		}
	default:
		common.WriteErrorResponse(w, h.log, err)
		return
	}
	orderBy := strParams[0]
	asc := strParams[1]
	if !searchEmpty {
		search, err = ValidateSearchQuery(strParams[2])
		if err != nil {
			common.WriteErrorResponse(w, h.log, err)
			return
		}
	}

	pageParams, err := newVideoPageParams(orderBy, offset, limit, asc)
	if err != nil {
		common.WriteErrorResponse(w, h.log, err)
		return
	}

	var videos []domain.Video
	if search == "" {
		videos, err = h.videoInteractor.GetByPublisher(r.Context(), publisherID, pageParams)
		if err != nil {
			common.WriteErrorResponse(w, h.log, err)
			return
		}
	} else {
		videos, err = h.videoInteractor.SearchPublisher(r.Context(),
			publisherID, search, pageParams)
		if err != nil {
			common.WriteErrorResponse(w, h.log, err)
			return
		}
	}
	common.WriteResponse(w, h.log,
		VideosResponseBody{Videos: slicex.Map(videos, DtoVideo)}, http.StatusOK)
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
	urlValues, err := common.ParseUrlValues(r.URL.RawQuery)
	if err != nil {
		common.WriteErrorResponse(w, h.log, err)
		return
	}

	resInt, err := common.ParseIntsUrlParams(urlValues, LimitUrlParam, OffsetUrlParam)
	if err != nil {
		common.WriteErrorResponse(w, h.log, err)
		return
	}
	limit := resInt[0]
	offset := resInt[1]

	resStr, err := common.ParseStringsUrlParams(urlValues,
		SearchUrlParam,
		SortByUrlParam,
		IsAscUrlParam)

	if err != nil {
		common.WriteErrorResponse(w, h.log, err)
		return
	}
	search := resStr[0]
	orderBy := resStr[1]
	asc := resStr[2]

	pageParams, err := newVideoPageParams(orderBy, offset, limit, asc)
	if err != nil {
		common.WriteErrorResponse(w, h.log, err)
		return
	}

	videos, err := h.videoInteractor.SearchGlobal(r.Context(), search, pageParams)
	if err != nil {
		common.WriteErrorResponse(w, h.log, err)
		return
	}

	common.WriteResponse(w, h.log,
		VideosResponseBody{Videos: slicex.Map(videos, DtoVideo)}, http.StatusOK)
}
