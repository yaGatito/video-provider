package testdb

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"
	"video-service/internal/domain"
	"video-service/internal/ports"

	"github.com/google/uuid"
)

const TAG = "[TEST-DB]"

type VideoRepoTestDB struct {
	log   *log.Logger
	store map[domain.UUID]domain.Video
	mu    sync.Mutex
}

var _ ports.VideoRepository = (*VideoRepoTestDB)(nil)

func NewVideoRepoTestDB(str map[domain.UUID]domain.Video, logger *log.Logger) *VideoRepoTestDB {
	repo := &VideoRepoTestDB{
		log:   logger,
		store: str,
	}

	desc := "asdasdsadasd"
	publisherID, _ := uuid.Parse("d9fa522f-0006-464f-8d68-356ba1d6ad7d")

	// Init value
	err := repo.CreateVideo(context.Background(), domain.Video{
		PublisherID: publisherID,
		Topic:       "sadas",
		Description: &desc,
	})

	if err != nil {
		return nil
	}

	return repo
}

func (r *VideoRepoTestDB) CreateVideo(ctx context.Context, video domain.Video) error {
	videoID := uuid.New()
	video.ID = domain.UUID(videoID)
	video.CreatedAt = time.Now()

	r.mu.Lock()
	r.log.Printf("%s video created: %v\n", TAG, video)
	r.store[videoID] = video
	r.mu.Unlock()
	return nil
}

func (r *VideoRepoTestDB) GetVideoByID(
	ctx context.Context,
	videoID domain.UUID,
) (domain.Video, error) {
	r.mu.Lock()
	video := r.store[videoID]
	r.log.Printf("%s retrieved video: %v\n", TAG, video)
	r.mu.Unlock()
	return video, nil
}

func (r *VideoRepoTestDB) GetPublisherVideos(
	ctx context.Context,
	publisherID domain.UUID,
	pagination ports.PageRequest,
) ([]domain.Video, error) {

	res := make([]domain.Video, 1)

	r.mu.Lock()
	for _, v := range r.store {
		if publisherID == v.PublisherID {
			res = append(res, v)
		}
	}
	r.log.Printf("retrieved videos: %v\n", res)
	r.mu.Unlock()
	return res, nil
}

func (r *VideoRepoTestDB) SearchPublisher(
	ctx context.Context,
	publisherID domain.UUID,
	search ports.VideoSearch,
) ([]domain.Video, error) {

	res := make([]domain.Video, 1)

	r.mu.Lock()
	for _, v := range r.store {
		if publisherID == v.PublisherID &&
			(strings.Contains(v.Topic, string(search.Query)) ||
				strings.Contains(*v.Description, string(search.Query))) {
			res = append(res, v)
		}
	}
	r.log.Printf("retrieved videos: %v\n", res)
	r.mu.Unlock()

	return res, nil
}

func (r *VideoRepoTestDB) SearchGlobal(
	ctx context.Context,
	search ports.VideoSearch,
) ([]domain.Video, error) {
	res := make([]domain.Video, 1)

	r.mu.Lock()
	for _, v := range r.store {
		if strings.Contains(v.Topic, string(search.Query)) ||
			strings.Contains(*v.Description, string(search.Query)) {
			res = append(res, v)
		}
	}
	r.log.Printf("retrieved videos: %v\n", res)
	r.mu.Unlock()

	return res, nil
}
