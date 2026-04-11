package httpadp

import (
	"fmt"
	"pkg/shared"
	"strings"
	"video-service/domain"
	"video-service/policy"

	"github.com/go-playground/validator/v10"
)

func NewVideoValidator() (*validator.Validate, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())

	err := validate.RegisterValidation("maxTopic", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) <= policy.TopicMaxLen
	})
	if err != nil {
		return nil, err
	}
	err = validate.RegisterValidation("minTopic", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) >= policy.TopicMinLen
	})
	if err != nil {
		return nil, err
	}
	err = validate.RegisterValidation("maxDescription", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) <= policy.DescriptionMaxLen
	})
	if err != nil {
		return nil, err
	}

	return validate, nil
}

func newVideoPageParams(
	orderByStr string,
	offset int32,
	limit int32,
	asc string,
) (domain.VideoPageParams, error) {
	offset, err := ValidateOffset(offset)
	if err != nil {
		return domain.VideoPageParams{}, err
	}
	limit, err = ValidateLimit(limit)
	if err != nil {
		return domain.VideoPageParams{}, err
	}
	orderBy, err := ValidateOrderBy(orderByStr)
	if err != nil {
		return domain.VideoPageParams{}, err
	}
	asc, err = ValidateIsAsc(asc)
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

func ValidateSearchQuery(query string) (string, error) {
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

func ValidateLimit(limit int32) (int32, error) {
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

func ValidateOffset(offset int32) (int32, error) {
	if offset < 0 {
		return 0, shared.NewError(
			shared.ErrInvalidInput,
			fmt.Sprintf("offset is zero or less: %d", offset),
			nil,
		)
	}
	return offset, nil
}

func ValidateOrderBy(orderBy string) (string, error) {
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

func ValidateIsAsc(asc string) (string, error) {
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
