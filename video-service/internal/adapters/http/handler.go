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
	"video-service/internal/app"
	"video-service/internal/domain"
	"video-service/internal/policy"
	"video-service/internal/ports"

	"github.com/gorilla/mux"
)

const (
	PathVarVideoID     = "videoID"
	PathVarPublisherID = "publisherID"

	URLParamSearch = "query"
	URLParamLimit  = "limit"
	URLParamOffset = "offset"
)

type VideoHandler struct {
	VideoInteractor app.VideoService
	IDGen           ports.IDGen
	log             *log.Logger
}

func NewVideoHandler(
	userInteractor app.VideoService,
	idGen ports.IDGen,
	log *log.Logger,
) VideoHandler {
	return VideoHandler{VideoInteractor: userInteractor, IDGen: idGen, log: log}
}

// Create godoc
// @Summary      Створити нове відео
// @Description  Створює новий запис відео для вказаного видавця
// @Tags         videos
// @Accept       json
// @Produce      json
// @Param        publisherID  path      string                 true  "ID видавця (UUID)"
// @Param        video        body      createVideoRequestBody true  "Дані відео"
// @Success      200          {object}  nil
// @Failure      400          {object}  string "Invalid input"
// @Failure      500          {object}  string "Internal error"
// @Router       /v1/videos/pub/{publisherID} [post]
func (h *VideoHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Required path variable
	publisherID, err := h.extractValidUUIDFromPathVar(r, PathVarPublisherID)
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, fmt.Errorf("invalid pub id: %e", err))
		return
	}

	// Required request body
	var createVideoRequestData createVideoRequestBody
	if err := json.NewDecoder(r.Body).Decode(&createVideoRequestData); !errors.Is(err, io.EOF) &&
		err != nil {
		h.log.Printf("Error decoding request body: %v", err)
		h.writeJSON(w, http.StatusBadRequest, fmt.Errorf("error decoding request body: %e", err))
		return
	}
	if err := createVideoRequestData.validate(); err != nil {
		h.log.Printf("Error validating request body: %v", err)
		h.writeJSON(w, http.StatusBadRequest, fmt.Errorf("error validating request body: %e", err))
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
		h.writeJSON(w, http.StatusInternalServerError, fmt.Errorf("error creating video: %e", err))
		return
	}

	h.writeJSON(w, http.StatusOK, nil)
}

// GetByID godoc
// @Summary      Отримати відео за ID
// @Description  Повертає деталі одного відео за його унікальним ідентифікатором
// @Tags         videos
// @Produce      json
// @Param        videoID  path      string  true  "ID відео (UUID)"
// @Success      200      {object}  VideoResponseBody
// @Failure      400      {object}  string
// @Failure      500      {object}  string
// @Router       /v1/videos/{videoID} [get]
func (h *VideoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Required path variable
	videoID, err := h.extractValidUUIDFromPathVar(r, PathVarVideoID)
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, fmt.Errorf("parse vid id param: %e", err))
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
// @Summary      Отримати відео видавця
// @Description  Повертає список відео конкретного видавця з підтримкою пагінації та пошуку
// @Tags         videos
// @Produce      json
// @Param        publisherID  path      string  true   "ID видавця (UUID)"
// @Param        query        query     string  false  "Пошуковий запит"
// @Param        limit        query     int     false  "Ліміт (за замовчуванням 10)"
// @Param        offset       query     int     false  "Зміщення (за замовчуванням 0)"
// @Success      200          {array}   VideoResponseBody
// @Router       /v1/videos/pub/{publisherID} [get]
func (h *VideoHandler) GetByPublisher(w http.ResponseWriter, r *http.Request) {
	// Required path variable
	publisherID, err := h.extractValidUUIDFromPathVar(r, PathVarPublisherID)
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, fmt.Errorf("parse pub id param: %e", err))
		return
	}

	// Url parameters
	values, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, fmt.Errorf("parse query params: %e", err))
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
		h.writeJSON(w, http.StatusBadRequest, fmt.Errorf("parse search param: %e", err))
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
				fmt.Errorf("interactor get by publisher error: %e", err),
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
			h.writeJSON(w, http.StatusInternalServerError, fmt.Errorf("interactor search publisher videos error: %e", err))
			return
		}
	}

	err = json.NewEncoder(w).Encode(h.toDtoVideos(videos))
	if err != nil {
		h.writeJSON(
			w,
			http.StatusInternalServerError,
			fmt.Errorf("error encoding in response body: %e", err),
		)
		return
	}
	h.log.Println("Response were written successfully")
}

// SearchGlobal godoc
// @Summary      Глобальний пошук
// @Description  Пошук відео по всій базі за ключовим словом
// @Tags         videos
// @Produce      json
// @Param        query   query     string  true   "Рядок пошуку"
// @Param        limit   query     int     false  "Ліміт"
// @Param        offset  query     int     false  "Зміщення"
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

func (h *VideoHandler) writeJSON(w http.ResponseWriter, status int, v any) {
	h.log.Println(v)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func (h VideoHandler) extractValidUUIDFromPathVar(r *http.Request, varName string) (domain.UUID, error) {
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
		return "", fmt.Errorf("failed to unescape %s: %s; err: %e", paramName, query, err)
	}

	return query, nil
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
