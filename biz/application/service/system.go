package service

import (
	"context"
	"github.com/CloudStriver/cloudmind-system/biz/infrastructure/convertor"
	notificationmapper "github.com/CloudStriver/cloudmind-system/biz/infrastructure/mapper/notification"
	slidermapper "github.com/CloudStriver/cloudmind-system/biz/infrastructure/mapper/slider"
	"github.com/CloudStriver/go-pkg/utils/pagination/mongop"
	"github.com/CloudStriver/go-pkg/utils/pconvertor"
	gensystem "github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/system"
	"github.com/google/wire"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SystemService interface {
	GetNotifications(ctx context.Context, req *gensystem.GetNotificationsReq) (resp *gensystem.GetNotificationsResp, err error)
	GetNotificationCount(ctx context.Context, req *gensystem.GetNotificationCountReq) (resp *gensystem.GetNotificationCountResp, err error)
	CreateNotifications(ctx context.Context, req *gensystem.CreateNotificationsReq) (resp *gensystem.CreateNotificationsResp, err error)
	DeleteSlider(ctx context.Context, req *gensystem.DeleteSliderReq) (resp *gensystem.DeleteSliderResp, err error)
	UpdateSlider(ctx context.Context, req *gensystem.UpdateSliderReq) (resp *gensystem.UpdateSliderResp, err error)
	CreateSlider(ctx context.Context, req *gensystem.CreateSliderReq) (resp *gensystem.CreateSliderResp, err error)
	GetSliders(ctx context.Context, req *gensystem.GetSlidersReq) (resp *gensystem.GetSlidersResp, err error)
	DeleteNotifications(ctx context.Context, req *gensystem.DeleteNotificationsReq) (resp *gensystem.DeleteNotificationsResp, err error)
	UpdateNotifications(ctx context.Context, req *gensystem.UpdateNotificationsReq) (resp *gensystem.UpdateNotificationsResp, err error)
}

type SystemServiceImpl struct {
	NotificationMongoMapper notificationmapper.INotificationMongoMapper
	SliderMongoMapper       slidermapper.ISliderMongoMapper
}

func (s *SystemServiceImpl) DeleteNotifications(ctx context.Context, req *gensystem.DeleteNotificationsReq) (resp *gensystem.DeleteNotificationsResp, err error) {
	if err = s.NotificationMongoMapper.DeleteNotifications(ctx, &notificationmapper.FilterOptions{
		OnlyUserId:          req.OnlyUserId,
		OnlyType:            req.OnlyType,
		OnlyIsRead:          req.OnlyIsRead,
		OnlyNotificationIds: req.OnlyNotificationIds,
	}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *SystemServiceImpl) UpdateNotifications(ctx context.Context, req *gensystem.UpdateNotificationsReq) (resp *gensystem.UpdateNotificationsResp, err error) {
	if err = s.NotificationMongoMapper.UpdateNotifications(ctx, &notificationmapper.FilterOptions{
		OnlyUserId:          req.OnlyUserId,
		OnlyType:            req.OnlyType,
		OnlyIsRead:          req.OnlyIsRead,
		OnlyNotificationIds: req.OnlyNotificationIds,
	}, req.IsRead); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *SystemServiceImpl) DeleteSlider(ctx context.Context, req *gensystem.DeleteSliderReq) (resp *gensystem.DeleteSliderResp, err error) {
	if err = s.SliderMongoMapper.DeleteOne(ctx, req.SliderId); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *SystemServiceImpl) UpdateSlider(ctx context.Context, req *gensystem.UpdateSliderReq) (resp *gensystem.UpdateSliderResp, err error) {
	oid, _ := primitive.ObjectIDFromHex(req.SliderId)
	if err = s.SliderMongoMapper.UpdateOne(ctx, &slidermapper.Slider{
		ID:       oid,
		ImageUrl: req.ImageUrl,
		LinkUrl:  req.LinkUrl,
		Type:     req.Type,
		IsPublic: req.IsPublic,
	}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *SystemServiceImpl) CreateSlider(ctx context.Context, req *gensystem.CreateSliderReq) (resp *gensystem.CreateSliderResp, err error) {
	if err = s.SliderMongoMapper.InsertOne(ctx, &slidermapper.Slider{
		ImageUrl: req.ImageUrl,
		LinkUrl:  req.LinkUrl,
		Type:     req.Type,
		IsPublic: req.IsPublic,
	}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *SystemServiceImpl) GetSliders(ctx context.Context, req *gensystem.GetSlidersReq) (resp *gensystem.GetSlidersResp, err error) {
	resp = new(gensystem.GetSlidersResp)
	p := pconvertor.PaginationOptionsToModelPaginationOptions(req.PaginationOptions)
	sliders, total, err := s.SliderMongoMapper.GetSlidersAndCount(ctx, &slidermapper.FilterOptions{
		OnlyType:     req.OnlyType,
		OnlyIsPublic: req.OnlyIsPublic,
	}, p, mongop.IdCursorType)
	if err != nil {
		return resp, err
	}

	resp.Sliders = lo.Map[*slidermapper.Slider, *gensystem.Slider](sliders,
		func(item *slidermapper.Slider, _ int) *gensystem.Slider {
			return convertor.SliderMapperToSlider(item)
		})

	if p.LastToken != nil {
		resp.Token = *p.LastToken
	}
	resp.Total = total
	return resp, nil
}

func (s *SystemServiceImpl) GetNotifications(ctx context.Context, req *gensystem.GetNotificationsReq) (resp *gensystem.GetNotificationsResp, err error) {
	resp = new(gensystem.GetNotificationsResp)
	p := pconvertor.PaginationOptionsToModelPaginationOptions(req.PaginationOptions)
	notifications, total, err := s.NotificationMongoMapper.GetNotificationsAndCount(ctx, &notificationmapper.FilterOptions{
		OnlyUserId: req.OnlyUserId,
		OnlyType:   req.OnlyType,
		OnlyIsRead: req.OnlyIsRead,
	}, p, mongop.IdCursorType)
	if err != nil {
		return resp, err
	}
	resp.Notifications = lo.Map[*notificationmapper.Notification, *gensystem.Notification](notifications,
		func(item *notificationmapper.Notification, index int) *gensystem.Notification {
			return convertor.NotificationMapperToNotification(item)
		})
	if p.LastToken != nil {
		resp.Token = *p.LastToken
	}
	resp.Total = total
	return resp, nil
}

func (s *SystemServiceImpl) GetNotificationCount(ctx context.Context, req *gensystem.GetNotificationCountReq) (resp *gensystem.GetNotificationCountResp, err error) {
	resp = new(gensystem.GetNotificationCountResp)
	if resp.Total, err = s.NotificationMongoMapper.Count(ctx, &notificationmapper.FilterOptions{
		OnlyUserId: req.OnlyUserId,
		OnlyType:   req.OnlyType,
		OnlyIsRead: req.OnlyIsRead,
	}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *SystemServiceImpl) CreateNotifications(ctx context.Context, req *gensystem.CreateNotificationsReq) (resp *gensystem.CreateNotificationsResp, err error) {
	notifications := lo.Map[*gensystem.NotificationInfo, *notificationmapper.Notification](req.Notifications, func(item *gensystem.NotificationInfo, _ int) *notificationmapper.Notification {
		return convertor.NotificationInfoToNotificationMapper(item)
	})
	if err = s.NotificationMongoMapper.InsertMany(ctx, notifications); err != nil {
		return resp, err
	}
	return resp, nil
}

var SystemSet = wire.NewSet(
	wire.Struct(new(SystemServiceImpl), "*"),
	wire.Bind(new(SystemService), new(*SystemServiceImpl)),
)
