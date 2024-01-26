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

func (s *SystemServerImpl) ReadNotification(ctx context.Context, req *system.ReadNotificationReq) (res *system.ReadNotificationResp, err error) {
	return s.SystemService.ReadNotification(ctx, req)
}

func (s *SystemServerImpl) CreateNotification(ctx context.Context, req *system.CreateNotificationReq) (res *system.CreateNotificationResp, err error) {
	return s.SystemService.CreateNotification(ctx, req)
}
