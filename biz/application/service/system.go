package service

import (
	"context"
	"fmt"
	"github.com/CloudStriver/cloudmind-system/biz/infrastructure/consts"
	"github.com/CloudStriver/cloudmind-system/biz/infrastructure/convertor"
	notificationmapper "github.com/CloudStriver/cloudmind-system/biz/infrastructure/mapper/notification"
	notificationcountmapper "github.com/CloudStriver/cloudmind-system/biz/infrastructure/mapper/notificationCount"
	slidermapper "github.com/CloudStriver/cloudmind-system/biz/infrastructure/mapper/slider"
	"github.com/CloudStriver/go-pkg/utils/pagination/mongop"
	"github.com/CloudStriver/go-pkg/utils/pconvertor"
	gensystem "github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/system"
	"github.com/google/wire"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SystemService interface {
	DeleteSlider(ctx context.Context, req *gensystem.DeleteSliderReq) (resp *gensystem.DeleteSliderResp, err error)
	UpdateSlider(ctx context.Context, req *gensystem.UpdateSliderReq) (resp *gensystem.UpdateSliderResp, err error)
	CreateSlider(ctx context.Context, req *gensystem.CreateSliderReq) (resp *gensystem.CreateSliderResp, err error)
	GetSliders(ctx context.Context, req *gensystem.GetSlidersReq) (resp *gensystem.GetSlidersResp, err error)
	GetNotifications(ctx context.Context, req *gensystem.GetNotificationsReq) (resp *gensystem.GetNotificationsResp, err error)
	GetNotificationCount(ctx context.Context, req *gensystem.GetNotificationCountReq) (resp *gensystem.GetNotificationCountResp, err error)
	CreateNotifications(ctx context.Context, req *gensystem.CreateNotificationsReq) (resp *gensystem.CreateNotificationsResp, err error)
	CreateNotificationCount(ctx context.Context, req *gensystem.CreateNotificationCountReq) (resp *gensystem.CreateNotificationCountResp, err error)
	DeleteNotifications(ctx context.Context, req *gensystem.DeleteNotificationsReq) (resp *gensystem.DeleteNotificationsResp, err error)
}

type SystemServiceImpl struct {
	NotificationMongoMapper      notificationmapper.INotificationMongoMapper
	NotificationCountMongoMapper notificationcountmapper.INotificationCountMongoMapper
	SliderMongoMapper            slidermapper.ISliderMongoMapper
	Redis                        *redis.Redis
}

func (s *SystemServiceImpl) DeleteNotifications(ctx context.Context, req *gensystem.DeleteNotificationsReq) (resp *gensystem.DeleteNotificationsResp, err error) {
	if err = s.NotificationMongoMapper.DeleteNotifications(ctx, &notificationmapper.FilterOptions{
		OnlyUserId:          lo.ToPtr(req.UserId),
		OnlyNotificationIds: req.NotificationIds,
		OnlyType:            req.OnlyType,
	}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *SystemServiceImpl) CreateNotificationCount(ctx context.Context, req *gensystem.CreateNotificationCountReq) (resp *gensystem.CreateNotificationCountResp, err error) {
	uid, _ := primitive.ObjectIDFromHex(req.UserId)
	err = s.NotificationCountMongoMapper.CreateNotificationCount(ctx, &notificationcountmapper.NotificationCount{
		ID:   uid,
		Read: 0,
	})
	return resp, err
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
	notifications, cnt, err := s.NotificationMongoMapper.GetNotificationsAndCount(ctx, &notificationmapper.FilterOptions{
		OnlyUserIds: []string{req.UserId, consts.NotificationSystemKey},
		OnlyType:    req.OnlyType,
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

	if req.OnlyType == nil {
		uid, _ := primitive.ObjectIDFromHex(req.UserId)
		if err = s.NotificationCountMongoMapper.UpdateNotificationCount(ctx, &notificationcountmapper.NotificationCount{
			ID:   uid,
			Read: cnt,
		}); err != nil {
			return resp, err
		}
	}
	return resp, nil
}

func (s *SystemServiceImpl) GetNotificationCount(ctx context.Context, req *gensystem.GetNotificationCountReq) (resp *gensystem.GetNotificationCountResp, err error) {
	cnt, err := s.NotificationMongoMapper.Count(ctx, &notificationmapper.FilterOptions{
		OnlyUserIds: []string{consts.NotificationSystemKey, req.UserId},
	})
	if err != nil {
		return resp, err
	}

	read, err := s.NotificationCountMongoMapper.GetNotificationCount(ctx, req.UserId)
	if err != nil {
		fmt.Println(err)

		return resp, err
	}
	return &gensystem.GetNotificationCountResp{
		Total: cnt - read,
	}, nil
}

func (s *SystemServiceImpl) CreateNotifications(ctx context.Context, req *gensystem.CreateNotificationsReq) (resp *gensystem.CreateNotificationsResp, err error) {
	if err = s.NotificationMongoMapper.InsertOne(ctx, &notificationmapper.Notification{
		TargetUserId:    req.TargetUserId,
		SourceUserId:    req.SourceUserId,
		SourceContentId: req.SourceContentId,
		Type:            req.Type,
		TargetType:      req.TargetType,
		Text:            req.Text,
	}); err != nil {
		return resp, err
	}

	return resp, nil
}

var SystemSet = wire.NewSet(
	wire.Struct(new(SystemServiceImpl), "*"),
	wire.Bind(new(SystemService), new(*SystemServiceImpl)),
)
