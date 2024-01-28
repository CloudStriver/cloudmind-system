package adaptor

import (
	"context"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/system"

	"github.com/CloudStriver/cloudmind-system/biz/application/service"
	"github.com/CloudStriver/cloudmind-system/biz/infrastructure/config"
)

type SystemServerImpl struct {
	*config.Config
	SystemService service.SystemService
}

func (s *SystemServerImpl) GetSliders(ctx context.Context, req *system.GetSlidersReq) (resp *system.GetSlidersResp, err error) {
	return s.SystemService.GetSliders(ctx, req)
}

func (s *SystemServerImpl) CreateSlider(ctx context.Context, req *system.CreateSliderReq) (resp *system.CreateSliderResp, err error) {
	return s.SystemService.CreateSlider(ctx, req)
}

func (s *SystemServerImpl) UpdateSlider(ctx context.Context, req *system.UpdateSliderReq) (resp *system.UpdateSliderResp, err error) {
	return s.SystemService.UpdateSlider(ctx, req)
}

func (s *SystemServerImpl) DeleteSlider(ctx context.Context, req *system.DeleteSliderReq) (resp *system.DeleteSliderResp, err error) {
	return s.SystemService.DeleteSlider(ctx, req)
}

func (s *SystemServerImpl) ReadNotifications(ctx context.Context, req *system.ReadNotificationsReq) (res *system.ReadNotificationsResp, err error) {
	return s.SystemService.ReadNotifications(ctx, req)
}

func (s *SystemServerImpl) GetNotifications(ctx context.Context, req *system.GetNotificationsReq) (res *system.GetNotificationsResp, err error) {
	return s.SystemService.GetNotifications(ctx, req)
}

func (s *SystemServerImpl) CleanNotification(ctx context.Context, req *system.CleanNotificationReq) (res *system.CleanNotificationResp, err error) {
	return s.SystemService.CleanNotification(ctx, req)
}

func (s *SystemServerImpl) GetNotificationCount(ctx context.Context, req *system.GetNotificationCountReq) (res *system.GetNotificationCountResp, err error) {
	return s.SystemService.GetNotificationCount(ctx, req)
}

func (s *SystemServerImpl) CreateNotifications(ctx context.Context, req *system.CreateNotificationsReq) (res *system.CreateNotificationsResp, err error) {
	return s.SystemService.CreateNotifications(ctx, req)
}
