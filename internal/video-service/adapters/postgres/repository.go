package postgres

import (
	"context"
	"time"
	postgres "video-provider/internal/video-service/adapters/postgres/db"
	"video-provider/internal/video-service/domain"
	"video-provider/internal/video-service/ports"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yaGatito/slicex"
)

const (
	OrderByCreatedAt string = "createdAt"
	OrderAsc         string = " ASC"
	OrderDesc        string = " DESC"
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
	params domain.VideoPageParams,
) ([]domain.Video, error) {

	args := postgres.GetVideosByPublisherParams{
		Publisherid: publisherID,
		Offset:      params.Offset,
		Limit:       params.Limit,
	}
	videos, err := r.queries.GetVideosByPublisher(ctx, args)
	if err != nil {
		return nil, err
	}

	return slicex.Map(videos, toDomainVideo), nil
}

func (r *VideoRepoPostgreSQL) SearchPublisher(
	ctx context.Context,
	publisherID domain.UUID,
	query string,
	params domain.VideoPageParams,
) ([]domain.Video, error) {

	args := postgres.SearchPublisherParams{
		Publisherid: publisherID,
		Column2:     query,
		Column3:     getOrderBy(params.OrderBy, params.Asc),
		Offset:      params.Offset,
		Limit:       params.Limit,
	}
	videos, err := r.queries.SearchPublisher(ctx, args)
	if err != nil {
		return nil, err
	}

	return slicex.Map(videos, toDomainVideo), nil
}

func (r *VideoRepoPostgreSQL) SearchGlobal(
	ctx context.Context,
	query string,
	params domain.VideoPageParams,
) ([]domain.Video, error) {

	args := postgres.SearchGlobalParams{
		Column1: query,
		Column2: getOrderBy(params.OrderBy, params.Asc),
		Offset:  params.Offset,
		Limit:   params.Limit,
	}
	videos, err := r.queries.SearchGlobal(ctx, args)
	if err != nil {
		return nil, err
	}

	return slicex.Map(videos, toDomainVideo), nil
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

func getOrderBy(order string, asc string) string {
	switch order {
	case domain.OrderByDate:
		order = OrderByCreatedAt
	}

	switch asc {
	case domain.AscOrder:
		asc = OrderAsc
	case domain.DescOrder:
		asc = domain.DescOrder
	}

	return order + asc
}
