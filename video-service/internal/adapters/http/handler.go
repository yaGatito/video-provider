package httpadp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"video-service/internal/app"
	"video-service/internal/domain"
	"video-service/internal/policy"
	"video-service/internal/ports"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const VIDEO_ID_PATH_VAR = "videoId"
const PUBLISHER_ID_PATH_VAR = "publisherId"

const SEARCH_URL_PARAM = "query"
const LIMIT_URL_PARAM = "limit"
const OFFSET_URL_PARAM = "offset"

type VideoHandler struct {
	VideoInteractor app.VideoService
	IDGen           ports.IDGen
	log             *log.Logger
}

func NewVideoHandler(userInteractor app.VideoService, idGen ports.IDGen, log *log.Logger) VideoHandler {
	return VideoHandler{VideoInteractor: userInteractor, IDGen: idGen, log: log}
}

// curl.exe -X POST "http://localhost:8081/v1/videos/pub/d9fa522f-0006-464f-8d68-356ba1d6ad7d" -H "Content-Type: application/json" -d '{"topic":"huy sosal","description":"sadasd"}'
func (h *VideoHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Required path variable
	publisherID, err := h.extractUUIDFromPathVar(r, PUBLISHER_ID_PATH_VAR)
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, err)
		return
	}

	// Required request body
	var createVideoRequestData createVideoRequestBody
	if err := json.NewDecoder(r.Body).Decode(&createVideoRequestData); !errors.Is(err, io.EOF) && err != nil {
		h.log.Printf("Error decoding request body: %v", err)
		h.writeJSON(w, http.StatusBadRequest, err)
		return
	}
	if err := createVideoRequestData.validate(); err != nil {
		h.log.Printf("Error validating request body: %v", err)
		h.writeJSON(w, http.StatusBadRequest, err)
		return
	}

	// Calling the interactor
	err = h.VideoInteractor.Create(r.Context(), domain.Video{
		PublisherID: domain.UUID(publisherID),
		Topic:       createVideoRequestData.Topic,
		Description: &createVideoRequestData.Description,
	})

	if err != nil {
		h.log.Printf("Error creating video: %v", err)
		h.writeJSON(w, http.StatusInternalServerError, err)
		return
	}

	h.writeJSON(w, http.StatusOK, nil)
}

func (h *VideoHandler) GetById(w http.ResponseWriter, r *http.Request) {
	// Required path variable
	videoID, err := h.extractUUIDFromPathVar(r, VIDEO_ID_PATH_VAR)
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, err)
		return
	}

	// Calling the interactor
	video, err := h.VideoInteractor.GetByID(r.Context(), domain.UUID(videoID))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving video: %v", err), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(h.toDtoVideo(video))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding video response: %v", err), http.StatusInternalServerError)
		return
	}
	h.log.Println("Response were written successfully")
}

func (h *VideoHandler) GetByPublisher(w http.ResponseWriter, r *http.Request) {
	// Required path variable
	publisherID, err := h.extractUUIDFromPathVar(r, PUBLISHER_ID_PATH_VAR)
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, err)
		return
	}

	// Optional url parameters
	offset := h.extractOptionalIntFromURLVars(r.URL, OFFSET_URL_PARAM)
	limit := h.extractOptionalIntFromURLVars(r.URL, LIMIT_URL_PARAM)
	offset, limit = app.ValidatePagination(offset, limit)
	// Not length exceeded search string or empty string
	search := h.extractOptionalStringFromURLVars(r.URL, SEARCH_URL_PARAM, policy.MAX_SEARCH_BYTES_SIZE)

	var videos []domain.Video

	// Calling the interactor
	if search != "" {
		videos, err = h.VideoInteractor.SearchPublisher(r.Context(), publisherID, search, limit, offset)
		if err != nil {
			h.writeJSON(w, http.StatusInternalServerError, err)
			return
		}
	} else {
		videos, err = h.VideoInteractor.GetByPublisher(r.Context(), publisherID, limit, offset)
		if err != nil {
			h.writeJSON(w, http.StatusInternalServerError, err)
			return
		}
	}

	err = json.NewEncoder(w).Encode(h.toDtoVideos(videos))
	if err != nil {
		h.writeJSON(w, http.StatusInternalServerError, err)
		return
	}
	h.log.Println("Response were written successfully")
}

func (h *VideoHandler) SearchGlobal(w http.ResponseWriter, r *http.Request) {
	// Required url parameters
	search, err := h.extractStringFromURLVars(r.URL, SEARCH_URL_PARAM, policy.MAX_SEARCH_BYTES_SIZE)
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
	offset := h.extractOptionalIntFromURLVars(r.URL, OFFSET_URL_PARAM)
	limit := h.extractOptionalIntFromURLVars(r.URL, LIMIT_URL_PARAM)
	offset, limit = app.ValidatePagination(offset, limit)

	// Calling the interactor
	videos, err := h.VideoInteractor.SearchGlobal(r.Context(), search, limit, offset)
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

func (h *VideoHandler) writeJSON(w http.ResponseWriter, status int, v any) {
	h.log.Println(v)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func (h VideoHandler) extractUUIDFromPathVar(r *http.Request, varName string) (uuid.UUID, error) {
	id, ok := mux.Vars(r)[varName]
	if !ok {
		return uuid.UUID{}, fmt.Errorf("no %s provided ", varName)
	}
	idSize := len([]byte(id))
	if idSize == 0 {
		return uuid.UUID{}, ValidationError{
			ErrorCode: ID_EMPTY, ErrorMessage: varName + " is empty",
		}
	}
	if idSize > policy.MAX_ID_BYTES_SIZE {
		return uuid.UUID{}, ValidationError{
			ErrorCode: ID_SIZE_EXCEEDED, ErrorMessage: varName + " len is more then expected",
		}
	}
	res, err := uuid.Parse(string(id))
	if err != nil {
		return uuid.UUID{}, ValidationError{
			ErrorCode: ID_SIZE_EXCEEDED, ErrorMessage: varName + " len is more then expected",
		}
	}

	return res, nil
}

func (h VideoHandler) extractOptionalIntFromURLVars(u *url.URL, paramName string) int32 {
	res, _ := strconv.Atoi(u.Query().Get(paramName))
	return int32(res)
}

func (h VideoHandler) extractOptionalStringFromURLVars(u *url.URL, paramName string, maxBytesLimit int) string {
	queryStr := u.Query().Get(paramName)

	if len(queryStr) > maxBytesLimit {
		return ""
	}

	return queryStr
}

func (h VideoHandler) extractStringFromURLVars(u *url.URL, paramName string, maxBytesLimit int) (string, error) {
	queryStr := u.Query().Get(paramName)
	if queryStr == "" {
		return "", fmt.Errorf("%s empty", paramName)
	}
	if len(queryStr) > maxBytesLimit {
		return "", fmt.Errorf("%s size exceeded", paramName)
	}
	return queryStr, nil
}

func (h VideoHandler) toDtoVideo(v domain.Video) VideoResponseBody {
	return VideoResponseBody{
		ID:          v.ID.String(),
		PublisherID: v.ID.String(),
		Topic:       v.Topic,
		Description: *v.Description,
		CreatedAt:   v.CreatedAt,
	}
}

func (h VideoHandler) toDtoVideos(videos []domain.Video) []VideoResponseBody {
	res := make([]VideoResponseBody, len(videos))
	for _, v := range videos {
		res = append(res, h.toDtoVideo(v))
	}
	return res
}
