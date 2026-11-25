package eventCategory;

import (
	"context"
	"time"
	"ticketBooking/internal/cache"
	"ticketBooking/internal/eventCategory/dto"
	"ticketBooking/internal/httpresponse"
	"ticketBooking/internal/utils/i18n"
	"ticketBooking/internal/utils/query"
)

type Service struct {
	repo Repository
}

func NewEventCategoryService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateEventCategory(req *dto.CreateEventCategoryRequest) (*dto.RawResponse, error) {
	eventCategory := &EventCategory{
		Name:            i18n.LocalizedString(req.Name),
		Description:      i18n.LocalizedString(req.Description),
		EventCategoryImageURL:    req.EventCategoryImageURL,
		EventCategoryImageId:     req.EventCategoryImageId,

	}

	if err := s.repo.Create(eventCategory); err != nil {
		return nil, err
	}
	_ = cache.Get().DeleteByPrefix(context.Background(), "categories:")
	return eventCategory.ToRawResponse(), nil
}

func (s *Service) GetAllEventCategories(params query.Params, lang string) (*httpresponse.PaginatedData, error) {
	// search এর জন্য language লাগে, খালি হলে en
	searchLang := lang
	if searchLang == "" {
		searchLang = "en"
	}

	cacheKey := cache.GenerateKey("categories:all", lang, searchLang, params.Page, params.Limit, params.Search, params.SortBy, params.SortDir)
	var cachedData httpresponse.PaginatedData
	if found, _ := cache.Get().Get(context.Background(), cacheKey, &cachedData); found {
		return &cachedData, nil
	}

	eventCategories, total, err := s.repo.GetAll(params, searchLang)
	if err != nil {
		return nil, err
	}

	docs := make([]any, 0, len(eventCategories))
	for _, e := range eventCategories {
		if lang == "" {
			docs = append(docs, e.ToRawResponse())
		} else {
			docs = append(docs, e.ToResponse(lang))
		}
	}

	meta := httpresponse.BuildPaginationMeta(total, params.Page, params.Limit)
	result := &httpresponse.PaginatedData{Docs: docs, PaginationMeta: meta}
	_ = cache.Get().Set(context.Background(), cacheKey, result, 30*time.Minute)
	return result, nil
}

func (s *Service) GetEventCategoryByID(eventId uint, lang string) (any, error) {
	cacheKey := cache.GenerateKey("categories:detail", eventId, lang)
	var cachedDoc any
	if found, _ := cache.Get().Get(context.Background(), cacheKey, &cachedDoc); found {
		return cachedDoc, nil
	}

	eventCategory, err := s.repo.GetByID(eventId)
	if err != nil {
		return nil, err
	}
	var res any
	if lang == "" {
		res = eventCategory.ToRawResponse()
	} else {
		res = eventCategory.ToResponse(lang)
	}

	_ = cache.Get().Set(context.Background(), cacheKey, res, 30*time.Minute)
	return res, nil
}

func (s *Service) UpdateEventCategory(eventId uint, req *dto.UpdateEventCategoryRequest) (*dto.RawResponse, error) {
	eventCategory, err := s.repo.GetByID(eventId)
	if err != nil {
		return nil, err
	}

	eventCategory.Name = i18n.LocalizedString(req.Name)
	eventCategory.Description = i18n.LocalizedString(req.Description)
	eventCategory.EventCategoryImageURL = req.EventCategoryImageURL
	eventCategory.EventCategoryImageId = req.EventCategoryImageId

	if err := s.repo.Update(eventCategory); err != nil {
		return nil, err
	}

	_ = cache.Get().DeleteByPrefix(context.Background(), "categories:")
	return eventCategory.ToRawResponse(), nil
}

func (s *Service) DeleteEventCategory(eventId uint) error {
	if err := s.repo.Delete(eventId); err != nil {
		return err
	}
	_ = cache.Get().DeleteByPrefix(context.Background(), "categories:")
	return nil
}