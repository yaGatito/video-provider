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
	URLParamSearch = "query"
	URLParamLimit  = "limit"
	URLParamOffset = "offset"
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
	publisherID, err := h.extractValidUUIDFromPathVar(r, PathVarPublisherID)
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, fmt.Errorf("invalid pub id: %w", err))
		return
	}

	// Required request body
	var createVideoRequestData createVideoRequestBody
	if err := json.NewDecoder(r.Body).Decode(&createVideoRequestData); !errors.Is(err, io.EOF) &&
		err != nil {
		h.log.Printf("Error decoding request body: %v", err)
		h.writeJSON(w, http.StatusBadRequest, fmt.Errorf("error decoding request body: %w", err))
		return
	}
	if err := createVideoRequestData.validate(); err != nil {
		h.log.Printf("Error validating request body: %v", err)
		h.writeJSON(w, http.StatusBadRequest, fmt.Errorf("error validating request body: %w", err))
		return
	}

	// Calling the interactor
	video := domain.Video{
		PublisherID: domain.UUID(publisherID),
		Topic:       createVideoRequestData.Topic,
		Description: createVideoRequestData.Description,
	}
	video, err = h.VideoInteractor.Create(r.Context(), video)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Error creating video: %v", err),
			http.StatusInternalServerError,
		)
		return
	}

	err = json.NewEncoder(w).Encode(h.toDtoVideo(video))
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Error encoding video response: %v", err),
			http.StatusInternalServerError,
		)
		return
	}
	h.log.Println("Response were written successfully")

	h.writeJSON(w, http.StatusOK, nil)
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
// @Router       /v1/videos/{videoID} [get]
func (h *VideoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Required path variable
	videoID, err := h.extractValidUUIDFromPathVar(r, PathVarVideoID)
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, fmt.Errorf("parse vid id param: %w", err))
		return
	}

	// Calling the interactor
	video, err := h.VideoInteractor.GetByID(r.Context(), domain.UUID(videoID))
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Error retrieving video: %v", err),
			http.StatusInternalServerError,
		)
		return
	}

	err = json.NewEncoder(w).Encode(h.toDtoVideo(video))
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Error encoding video response: %v", err),
			http.StatusInternalServerError,
		)
		return
	}
	h.log.Println("Response were written successfully")
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
	publisherID, err := h.extractValidUUIDFromPathVar(r, PathVarPublisherID)
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, fmt.Errorf("parse pub id param: %w", err))
		return
	}

	// Url parameters
	values, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, fmt.Errorf("parse query params: %w", err))
		return
	}
	offset, limit := app.ValidatePagination(
		h.extractOptionalIntFromURLVars(values, URLParamOffset),
		h.extractOptionalIntFromURLVars(values, URLParamLimit))

	search, err := h.extractOptionalStringFromURLVars(
		values,
		URLParamSearch,
		policy.MaxSearchBytesSize,
	)
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, fmt.Errorf("parse search param: %w", err))
		return
	}

	var videos []domain.Video

	// Calling the interactor
	if search == "" {
		videos, err = h.VideoInteractor.GetByPublisher(r.Context(), publisherID, offset, limit)
		if err != nil {
			h.writeJSON(
				w,
				http.StatusInternalServerError,
				fmt.Errorf("interactor get by publisher error: %w", err),
			)
			return
		}
	} else {
		videos, err = h.VideoInteractor.SearchPublisher(
			r.Context(),
			publisherID,
			search,
			offset,
			limit,
		)
		if err != nil {
			h.writeJSON(
				w,
				http.StatusInternalServerError,
				fmt.Errorf("interactor search publisher videos error: %w", err))
			return
		}
	}

	err = json.NewEncoder(w).Encode(h.toDtoVideos(videos))
	if err != nil {
		h.writeJSON(
			w,
			http.StatusInternalServerError,
			fmt.Errorf("error encoding in response body: %w", err),
		)
		return
	}
	h.log.Println("Response were written successfully")
}

// SearchGlobal godoc
// @Summary      Global search
// @Description  Search videos in the entire database by a keyword
// @Tags         videos
// @Produce      json
// @Param        query   query     string  true   "Search query"
// @Param        limit   query     int     false  "Limit (default 10)"
// @Param        offset  query     int     false  "Offset (default 0)"
// @Success      200     {array}   VideoResponseBody
// @Router       /v1/videos/search [get]
func (h *VideoHandler) SearchGlobal(w http.ResponseWriter, r *http.Request) {
	// Required url parameters
	values, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, err)
		return
	}
	search, err := h.extractStringFromURLVars(values, URLParamSearch, policy.MaxSearchBytesSize)
	if err != nil {
		h.writeJSON(w, http.StatusInternalServerError, err)
		return
	}
	search, err = app.ValidateSearchQuery(search)
	if err != nil {
		h.writeJSON(w, http.StatusInternalServerError, err)
		return
	}

	// Optional url parameters
	offset, limit := app.ValidatePagination(
		h.extractOptionalIntFromURLVars(values, URLParamOffset),
		h.extractOptionalIntFromURLVars(values, URLParamLimit))

	// Calling the interactor
	videos, err := h.VideoInteractor.SearchGlobal(r.Context(), search, offset, limit)
	if err != nil {
		h.writeJSON(w, http.StatusInternalServerError, err)
		return
	}

	err = json.NewEncoder(w).Encode(h.toDtoVideos(videos))
	if err != nil {
		h.writeJSON(w, http.StatusInternalServerError, err)
		return
	}
	h.log.Println("Response were written successfully")
}

// GetAllVideos godoc
// @Summary      Get all videos
// @Description  Returns all videos with mock data for frontend
// @Tags         videos
// @Produce      json
// @Success      200      {object}  map[string]interface{}
// @Router       /api/videos [get]
func (h *VideoHandler) GetAllVideos(w http.ResponseWriter, r *http.Request) {
	videos := map[string]interface{}{
		"videos": []map[string]interface{}{
			{
				"id":           1,
				"title":        "Learn Go Programming",
				"description":  "Complete guide to learning Go from basics to advanced concepts",
				"likes":        1250,
				"dislikes":     15,
				"comments":     89,
				"previewImage": "https://via.placeholder.com/300x200?text=Go+Programming",
			},
			{
				"id":           2,
				"title":        "Web Development Tips",
				"description":  "Essential tips and tricks for building modern web applications",
				"likes":        2100,
				"dislikes":     32,
				"comments":     156,
				"previewImage": "https://via.placeholder.com/300x200?text=Web+Development",
			},
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(videos)
}

// GetVideoByID godoc
// @Summary      Get video by ID
// @Description  Returns a single video by ID with mock data for frontend
// @Tags         videos
// @Produce      json
// @Param        id  path      int  true  "Video ID"
// @Success      200      {object}  map[string]interface{}
// @Router       /api/videos/{id} [get]
func (h *VideoHandler) GetVideoByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	videoID := vars["id"]

	// Mock data for videos
	videoData := map[int]map[string]interface{}{
		1: {
			"id":           1,
			"title":        "Learn Go Programming",
			"description":  "Complete guide to learning Go from basics to advanced concepts. In this comprehensive video series, you'll learn everything you need to know about Go programming language, including:\n\n- Go basics and syntax\n- Goroutines and concurrency\n- Channels and synchronization\n- Error handling best practices\n- Building production-ready applications\n\nWhether you're a beginner or experienced programmer, this course will help you master Go programming.",
			"likes":        1250,
			"dislikes":     15,
			"comments":     89,
			"previewImage": "https://via.placeholder.com/600x400?text=Go+Programming",
		},
		2: {
			"id":           2,
			"title":        "Web Development Tips",
			"description":  "Essential tips and tricks for building modern web applications. This video covers:\n\n- Frontend best practices\n- React and Vue.js frameworks\n- State management strategies\n- Performance optimization\n- Security considerations\n- Testing and debugging\n\nLearn from industry experts and improve your web development skills.",
			"likes":        2100,
			"dislikes":     32,
			"comments":     156,
			"previewImage": "https://via.placeholder.com/600x400?text=Web+Development",
		},
	}

	id := 0
	fmt.Sscanf(videoID, "%d", &id)

	if video, exists := videoData[id]; exists {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode(video)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Video not found"})
	}
}

func (h *VideoHandler) writeJSON(w http.ResponseWriter, status int, v any) {
	h.log.Println(v)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func (h VideoHandler) extractValidUUIDFromPathVar(
	r *http.Request,
	varName string,
) (domain.UUID, error) {
	id, ok := mux.Vars(r)[varName]
	if !ok {
		return domain.UUID{}, fmt.Errorf("no %s provided ", varName)
	}
	idSize := len([]byte(id))
	if idSize == 0 {
		return domain.UUID{}, ValidationError{
			ErrorCode: IDEmpty, ErrorMessage: varName + " is empty",
		}
	}
	if idSize > policy.MaxIDBytesSize {
		return domain.UUID{}, ValidationError{
			ErrorCode: IDSizeExceeded, ErrorMessage: varName + " len is more then expected",
		}
	}
	res, err := h.IDGen.Parse(string(id))
	if err != nil {
		return domain.UUID{}, ValidationError{
			ErrorCode: IDSizeExceeded, ErrorMessage: varName + " len is more then expected",
		}
	}

	return res, nil
}

func (h VideoHandler) extractOptionalIntFromURLVars(values url.Values, paramName string) int32 {
	res, _ := strconv.Atoi(values.Get(paramName))
	return int32(res)
}

func (h VideoHandler) extractStringFromURLVars(
	values url.Values,
	paramName string,
	maxBytesLimit int,
) (string, error) {
	query, err := h.extractOptionalStringFromURLVars(values, paramName, maxBytesLimit)
	if len(query) == 0 && err == nil {
		return "", fmt.Errorf("%s empty", paramName)
	}
	if err != nil {
		return "", err
	}
	return query, nil
}

func (h VideoHandler) extractOptionalStringFromURLVars(
	values url.Values,
	paramName string,
	maxBytesLimit int,
) (string, error) {
	query := values.Get(paramName)
	// Letting query to be returned without error while being empty.
	if len(query) == 0 {
		return "", nil
	}
	if len(query) > maxBytesLimit {
		return "", fmt.Errorf("query search is too large")
	}
	query, err := url.QueryUnescape(query)
	if err != nil {
		return "", fmt.Errorf("failed to unescape %s: %s; err: %w", paramName, query, err)
	}

	return query, nil
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
