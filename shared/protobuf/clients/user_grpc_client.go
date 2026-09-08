package clients

import (
	_dto "common/domain/dto"
	"context"
	"errors"
	authpb "pb/types/auth"
	sharepb "pb/types/shared"
	userpb "pb/types/user"

	"google.golang.org/grpc"
)

type UserGrpcClient struct {
	Client                       userpb.ProfileServiceClient
	InternalClient               userpb.InternalUserServiceClient
	OrganizationMembershipClient userpb.InternalOrganizationMembershipServiceClient
	AdminClient                  userpb.AdminUserProfileServiceClient
	AccessClient                 userpb.InternalAccessServiceClient
	PermissionClient             authpb.PermissionServiceClient
	RoleClient                   authpb.RoleServiceClient
}

func BindUserGrpcClient(conn grpc.ClientConnInterface) *UserGrpcClient {
	return &UserGrpcClient{
		Client:                       userpb.NewProfileServiceClient(conn),
		InternalClient:               userpb.NewInternalUserServiceClient(conn),
		OrganizationMembershipClient: userpb.NewInternalOrganizationMembershipServiceClient(conn),
		AdminClient:                  userpb.NewAdminUserProfileServiceClient(conn),
		AccessClient:                 userpb.NewInternalAccessServiceClient(conn),
		PermissionClient:             authpb.NewPermissionServiceClient(conn),
		RoleClient:                   authpb.NewRoleServiceClient(conn),
	}
}

func (c *UserGrpcClient) GetAllPermissions(ctx context.Context) ([]*sharepb.Permission, error) {
	if c.PermissionClient == nil {
		return nil, errors.New("user permission client not available")
	}
	response, err := c.PermissionClient.GetAllPermissions(ctx, &sharepb.IdRequest{})
	if err != nil {
		return nil, err
	}
	return response.GetData(), nil
}

func (c *UserGrpcClient) GetMapByIDs(ctx context.Context, profileIDSet map[uint64]struct{}) (map[uint64]*sharepb.ProfileItem, error) {
	if c.Client == nil {
		return nil, errors.New("user client not avaiable")
	}

	profileIDs := make([]uint64, 0, len(profileIDSet))
	for id := range profileIDSet {
		profileIDs = append(profileIDs, id)
	}

	profiles, err := c.Client.GetProfileByIds(ctx, &sharepb.GetProfileByIdsRequest{Ids: profileIDs})
	if err != nil {
		return nil, err
	}

	profileMap := make(map[uint64]*sharepb.ProfileItem, len(profiles.Profiles))
	for _, profile := range profiles.Profiles {
		profileMap[profile.Id] = profile
	}

	return profileMap, nil
}

func (c *UserGrpcClient) GetProfileByIds(ctx context.Context, profileIds []uint64) (*sharepb.GetProfileByIdsResponse, error) {
	if c.Client == nil {
		return nil, errors.New("user client not avaiable")
	}

	return c.Client.GetProfileByIds(ctx, &sharepb.GetProfileByIdsRequest{Ids: profileIds})
}

func (c *UserGrpcClient) GetRoleByID(ctx context.Context, id uint64) (*sharepb.Role, error) {
	if c.InternalClient == nil {
		return nil, errors.New("user internal client not avaiable")
	}
	return c.InternalClient.GetRoleById(ctx, &sharepb.IdRequest{Id: id})
}

func (c *UserGrpcClient) GetRolesByIds(ctx context.Context, roleIds []uint64) (*authpb.RoleListResponse, error) {
	if c.InternalClient == nil {
		return nil, errors.New("user internal client not avaiable")
	}
	return c.InternalClient.GetRolesByIds(ctx, &sharepb.IdRequest{Ids: roleIds})
}

func (c *UserGrpcClient) GetRolesByRoleKeys(ctx context.Context, roleKeys []uint32) (*authpb.RoleListResponse, error) {
	if c.InternalClient == nil {
		return nil, errors.New("user internal client not avaiable")
	}
	return c.InternalClient.GetRolesByRoleKeys(ctx, &userpb.GetRolesByRoleKeysRequest{RoleKeys: roleKeys})
}

func (c *UserGrpcClient) GetRolePermissions(ctx context.Context, req *sharepb.IdRequest) (*userpb.GetRolePermissionsResponse, error) {
	if c.InternalClient == nil {
		return nil, errors.New("user internal client not avaiable")
	}
	return c.InternalClient.GetRolePermissions(ctx, req)
}

func (c *UserGrpcClient) GetRoleIdsByProfileId(ctx context.Context, profileId uint64) ([]uint64, error) {
	if c.InternalClient == nil {
		return nil, errors.New("user internal client not avaiable")
	}
	resp, err := c.InternalClient.GetRoleIdsByProfileId(ctx, &sharepb.IdRequest{Id: profileId})
	if err != nil {
		return nil, err
	}
	return resp.GetRoleIds(), nil
}

func (c *UserGrpcClient) GetAuthUsersByProfileIDs(ctx context.Context, profileIDs []uint64) (map[uint64]*_dto.AuthUserProfileDTO, error) {
	if c.InternalClient == nil {
		return nil, errors.New("user internal client not avaiable")
	}

	response, err := c.InternalClient.GetAuthUsersByProfileIDs(ctx, &userpb.GetAuthUsersByProfileIDsRequest{ProfileIds: profileIDs})
	if err != nil {
		return nil, err
	}

	roleMap := make(map[uint64]*_dto.AuthUserProfileDTO)
	for _, user := range response.Data {
		roleMap[user.ProfileId] = &_dto.AuthUserProfileDTO{
			ProfileID: user.ProfileId,
			RoleID:    user.RoleId,
		}

		if user.Role != nil {
			r := &_dto.RoleDTO{
				ID:       user.Role.Id,
				RoleName: user.Role.RoleName,
				RoleKey:  user.Role.RoleKey,
			}
			if user.Role.Color != nil {
				r.Color = &_dto.ColorDTO{
					ID:              user.Role.Color.Id,
					Code:            user.Role.Color.ColorKey,
					Color:           user.Role.Color.ContentColor,
					ContentColor:    user.Role.Color.ContentColor,
					BackgroundColor: user.Role.Color.BackgroundColor,
					ColorKey:        user.Role.Color.ColorKey,
					HexCode:         user.Role.Color.HexCode,
				}
			}

			roleMap[user.ProfileId].Role = r
		}
	}
	return roleMap, nil
}

func (c *UserGrpcClient) GetAuthAdminByIds(ctx context.Context, ids []uint64) (*userpb.GetAuthAdminByIdsResponse, error) {
	if c.InternalClient == nil {
		return nil, errors.New("user internal client not avaiable")
	}
	return c.InternalClient.GetAuthAdminByIds(ctx, &userpb.GetAuthAdminByIdsRequest{Ids: ids})
}

func (c *UserGrpcClient) GetPushTokensByProfileId(ctx context.Context, profileID uint64) (*userpb.GetPushTokensResponse, error) {
	if c.InternalClient == nil {
		return nil, errors.New("user internal client not avaiable")
	}
	return c.InternalClient.GetPushTokensByProfileId(ctx, &sharepb.IdRequest{Id: profileID})
}
