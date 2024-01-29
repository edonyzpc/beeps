package v2

import (
	"context"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	apiv2pb "github.com/usememos/memos/proto/gen/api/v2"
	"github.com/usememos/memos/store"
)

func (s *APIV2Service) GetWorkspaceProfile(_ context.Context, _ *apiv2pb.GetWorkspaceProfileRequest) (*apiv2pb.GetWorkspaceProfileResponse, error) {
	workspaceProfile := &apiv2pb.WorkspaceProfile{
		Version: s.Profile.Version,
		Mode:    s.Profile.Mode,
	}
	response := &apiv2pb.GetWorkspaceProfileResponse{
		WorkspaceProfile: workspaceProfile,
	}
	return response, nil
}

func (s *APIV2Service) UpdateWorkspaceProfile(ctx context.Context, request *apiv2pb.UpdateWorkspaceProfileRequest) (*apiv2pb.UpdateWorkspaceProfileResponse, error) {
	user, err := getCurrentUser(ctx, s.Store)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get current user: %v", err)
	}
	if user.Role != store.RoleHost {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}
	if request.UpdateMask == nil || len(request.UpdateMask.Paths) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "update mask is required")
	}

	// Update system settings.
	for _, field := range request.UpdateMask.Paths {
		if field == "allow_registration" {
			_, err := s.Store.UpsertSystemSetting(ctx, &store.SystemSetting{
				Name:  "allow-signup",
				Value: strconv.FormatBool(request.WorkspaceProfile.AllowRegistration),
			})
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to update allow_registration system setting: %v", err)
			}
		} else if field == "disable_password_login" {
			_, err := s.Store.UpsertSystemSetting(ctx, &store.SystemSetting{
				Name:  "disable-password-login",
				Value: strconv.FormatBool(request.WorkspaceProfile.DisablePasswordLogin),
			})
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to update disable_password_login system setting: %v", err)
			}
		} else if field == "additional_script" {
			_, err := s.Store.UpsertSystemSetting(ctx, &store.SystemSetting{
				Name:  "additional-script",
				Value: request.WorkspaceProfile.AdditionalScript,
			})
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to update additional_script system setting: %v", err)
			}
		} else if field == "additional_style" {
			_, err := s.Store.UpsertSystemSetting(ctx, &store.SystemSetting{
				Name:  "additional-style",
				Value: request.WorkspaceProfile.AdditionalStyle,
			})
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to update additional_style system setting: %v", err)
			}
		}
	}

	workspaceProfileMessage, err := s.GetWorkspaceProfile(ctx, &apiv2pb.GetWorkspaceProfileRequest{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get system info: %v", err)
	}
	return &apiv2pb.UpdateWorkspaceProfileResponse{
		WorkspaceProfile: workspaceProfileMessage.WorkspaceProfile,
	}, nil
}
