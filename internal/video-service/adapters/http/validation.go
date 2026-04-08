package httpadp

import (
	"fmt"
	"github.com/yaGatito/video-provider/internal/pkg/shared"
	"strings"
	"video-service/domain"
	"video-service/policy"

	"github.com/go-playground/validator/v10"
)

func newVideoValidator() *validator.Validate {
	validate := validator.New(validator.WithRequiredStructEnabled())

	validate.RegisterValidation("maxTopic", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) <= policy.TopicMaxLen
	})
	validate.RegisterValidation("minTopic", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) >= policy.TopicMinLen
	})
	validate.RegisterValidation("maxDescription", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) <= policy.DescriptionMaxLen
	})

	return validate
}

func newVideoPageParams(
	orderByStr string,
	offset int32,
	limit int32,
	asc string,
) (domain.VideoPageParams, error) {

	offset, err := validateOffset(offset)
	if err != nil {
		return domain.VideoPageParams{}, err
	}
	limit, err = validateLimit(limit)
	if err != nil {
		return domain.VideoPageParams{}, err
	}
	orderBy, err := validateOrderBy(orderByStr)
	if err != nil {
		return domain.VideoPageParams{}, err
	}
	asc, err = validateIsAsc(asc)
	if err != nil {
		return domain.VideoPageParams{}, err
	}

	return domain.VideoPageParams{
		Offset:  offset,
		Limit:   limit,
		OrderBy: orderBy,
		Asc:     asc,
	}, nil
}

func validateSearchQuery(query string) (string, error) {
	query = strings.TrimSpace(query)

	if len(query) < policy.SearchMinLen {
		return "", shared.NewError(
			shared.ErrInvalidInput,
			fmt.Sprintf("%s size is less than threshold: %d", query, policy.SearchMinLen),
			nil,
		)
	}
	if !policy.GetWordsFormatRE128().MatchString(query) {
		return "", shared.NewError(
			shared.ErrInvalidInput,
			"query string contains prohibited characters",
			nil,
		)
	}

	return query, nil
}

func validateLimit(limit int32) (int32, error) {
	if limit < policy.ThresholdVideosLimit {
		return 0, shared.NewError(
			shared.ErrInvalidInput,
			fmt.Sprintf("limit is less then threshold(%d): %d", policy.ThresholdVideosLimit, limit),
			nil,
		)
	}
	if limit > policy.VideosMaxLimit {
		return 0, shared.NewError(
			shared.ErrInvalidInput,
			fmt.Sprintf("limit reached maximum allowed value: %d", limit),
			nil,
		)
	}
	return limit, nil
}

func validateOffset(offset int32) (int32, error) {
	if offset < 0 {
		return 0, shared.NewError(
			shared.ErrInvalidInput,
			fmt.Sprintf("offset is zero or less: %d", offset),
			nil,
		)
	}
	return offset, nil
}

func validateOrderBy(orderBy string) (string, error) {
	switch orderBy {
	case domain.OrderByDate:
		return orderBy, nil
	default:
		return "", shared.NewError(
			shared.ErrInvalidInput,
			fmt.Sprintf("invalid orderBy argument: %s", orderBy),
			nil,
		)
	}
}

func validateIsAsc(asc string) (string, error) {
	switch asc {
	case domain.AscOrder:
		return domain.AscOrder, nil
	case domain.DescOrder:
		return domain.DescOrder, nil
	default:
		return "", shared.NewError(
			shared.ErrInvalidInput,
			fmt.Sprintf("invalid asc argument: %s; only `t` and `f` are allowed", asc),
			nil,
		)
	}
}
