package postgres

import (
	"context"
	"time"
	postgres "video-provider/internal/video-service/adapters/postgres/db"
	"video-provider/internal/video-service/domain"
	"video-provider/internal/video-service/ports"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VideoRepoPostgreSQL struct {
	queries *postgres.Queries
}

var _ ports.VideoRepository = (*VideoRepoPostgreSQL)(nil)

func NewVideoRepoPostgreSQL(db *pgxpool.Pool) ports.VideoRepository {
	v := VideoRepoPostgreSQL{}
	v.queries = postgres.New(db)
	return &v
}

func (r *VideoRepoPostgreSQL) CreateVideo(
	ctx context.Context,
	video domain.Video,
) (domain.Video, error) {
	arg := postgres.CreateVideoParams{
		Publisherid: video.PublisherID,
		Topic:       video.Topic,
		Description: pgtype.Text{String: video.Description, Valid: true},
	}
	res, err := r.queries.CreateVideo(ctx, arg)
	if err != nil {
		return domain.Video{}, err
	}
	return toDomainVideo(res), nil
}

func (r *VideoRepoPostgreSQL) GetVideoByID(
	ctx context.Context,
	id domain.UUID,
) (domain.Video, error) {
	video, err := r.queries.GetVideoByID(ctx, id)
	if err != nil {
		return domain.Video{}, err
	}

	return toDomainVideo(video), nil
}

func (r *VideoRepoPostgreSQL) GetPublisherVideos(
	ctx context.Context,
	publisherID domain.UUID,
	args ports.PageRequest,
) ([]domain.Video, error) {

	params := postgres.GetVideosByPublisherParams{
		Publisherid: publisherID,
		Offset:      args.Offset,
		Limit:       args.Limit,
	}
	videos, err := r.queries.GetVideosByPublisher(ctx, params)
	if err != nil {
		return nil, err
	}

	return toDomainVideos(videos), nil
}

func (r *VideoRepoPostgreSQL) SearchPublisher(
	ctx context.Context,
	publisherID domain.UUID,
	search ports.VideoSearch,
) ([]domain.Video, error) {

	params := postgres.SearchPublisherParams{
		Publisherid: publisherID,
		Concat:      search.Query,
		Offset:      search.Offset,
		Limit:       search.Limit,
	}
	videos, err := r.queries.SearchPublisher(ctx, params)
	if err != nil {
		return nil, err
	}

	return toDomainVideos(videos), nil
}

func (r *VideoRepoPostgreSQL) SearchGlobal(
	ctx context.Context,
	search ports.VideoSearch,
) ([]domain.Video, error) {
	params := postgres.SearchGlobalParams{
		Concat: search.Query,
		Offset: search.Offset,
		Limit:  search.Limit,
	}
	videos, err := r.queries.SearchGlobal(ctx, params)
	if err != nil {
		return nil, err
	}

	return toDomainVideos(videos), nil
}

func toDomainVideos(videos []postgres.Video) []domain.Video {
	res := make([]domain.Video, len(videos))
	for _, v := range videos {
		res = append(res, toDomainVideo(v))
	}
	return res

}

func toDomainVideo(video postgres.Video) domain.Video {
	return domain.Video{
		ID:          domain.UUID(video.ID),
		PublisherID: domain.UUID(video.Publisherid),
		Topic:       video.Topic,
		Description: video.Description.String,
		CreatedAt:   time.UnixMicro(video.Createdat.Microseconds),
		Status:      domain.Status(video.Status.String),
	}
}
