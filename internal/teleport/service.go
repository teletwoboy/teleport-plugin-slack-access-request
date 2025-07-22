package teleport

import (
	"context"
	"fmt"
	"sync"
	"teleport-plugin-slack-access-request/internal/teleport/builder/accessrequest"
	"teleport-plugin-slack-access-request/internal/teleport/models"
	teleporttypes "teleport-plugin-slack-access-request/internal/teleport/types"
	"time"

	"github.com/gravitational/teleport/api/types"
	events "github.com/gravitational/teleport/api/types/events"
)

const (
	MaxProcessedEvents = 1000
	CleanupThreshold   = 500
)

type Service interface {
	CreateAccessRequest(ctx context.Context, accessRequest *models.AccessRequest) (*models.AccessRequest, error)
	CreateAccessReview(ctx context.Context, accessReview *models.AccessReview) (*models.AccessReview, error)
	CreateUser(ctx context.Context, user models.User) (*models.User, error)
	ExistsAccessRequestByName(ctx context.Context, name string) (bool, error)
	FetchAccessRequests(ctx context.Context, builder accessrequest.FilterBuilder) ([]types.AccessRequest, error)
	FetchUsersWithoutSecrets(ctx context.Context) ([]models.User, error)
	FetchUserAccessInfo(ctx context.Context, user *models.User) (*teleporttypes.UserAccessInfo, error)
	GetAccessRequestByName(ctx context.Context, name string) (*models.AccessRequest, error)
	GetAccessRequestStateByName(ctx context.Context, name string) (string, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	SubmitAccessRequest(ctx context.Context, builder accessrequest.CreateBuilder) (types.AccessRequest, error)
	SubmitAccessRequestState(ctx context.Context, builder accessrequest.UpdateBuilder) error
	UpdateAccessRequestStateByName(ctx context.Context, accessRequest *models.AccessRequest) (*models.AccessRequest, error)
	StartMFAEventListener(ctx context.Context) error
}

type API interface {
	CreateAccessRequestV2(ctx context.Context, req types.AccessRequest) (types.AccessRequest, error)
	GetAccessCapabilities(ctx context.Context, req types.AccessCapabilitiesRequest) (*types.AccessCapabilities, error)
	GetAccessRequests(ctx context.Context, filter types.AccessRequestFilter) ([]types.AccessRequest, error)
	GetUsers(ctx context.Context, withSecrets bool) ([]types.User, error)
	SetAccessRequestState(ctx context.Context, params types.AccessRequestUpdate) error
	SearchEvents(ctx context.Context, fromUTC, toUTC time.Time, namespace string, eventTypes []string, limit int, order types.EventOrder, startKey string) ([]events.AuditEvent, string, error)
}

type Repository interface {
	CreateAccessRequest(ctx context.Context, accessRequest *models.AccessRequest) (*models.AccessRequest, error)
	CreateAccessReview(ctx context.Context, accessReview *models.AccessReview) (*models.AccessReview, error)
	CreateUser(ctx context.Context, user models.User) (*models.User, error)
	ExistsAccessRequestByName(ctx context.Context, name string) (bool, error)
	GetAccessRequestByName(ctx context.Context, name string) (*models.AccessRequest, error)
	GetAccessRequestStateByName(ctx context.Context, name string) (string, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	UpdateAccessRequestStateByName(ctx context.Context, accessRequest *models.AccessRequest) (*models.AccessRequest, error)
}

type service struct {
	api             API
	repo            Repository
	processedEvents map[string]bool
	mutex           sync.RWMutex
}

func NewService(api API, repo Repository) Service {
	return &service{api: api, repo: repo, processedEvents: make(map[string]bool)}
}

func (s *service) CreateAccessRequest(ctx context.Context, accessRequest *models.AccessRequest) (*models.AccessRequest, error) {
	return s.repo.CreateAccessRequest(ctx, accessRequest)
}

func (s *service) CreateAccessReview(ctx context.Context, accessReview *models.AccessReview) (*models.AccessReview, error) {
	return s.repo.CreateAccessReview(ctx, accessReview)
}

func (s *service) CreateUser(ctx context.Context, user models.User) (*models.User, error) {
	createdUser, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed tp create Teleport user: %w", err)
	}
	return createdUser, nil
}

func (s *service) ExistsAccessRequestByName(ctx context.Context, name string) (bool, error) {
	exists, err := s.repo.ExistsAccessRequestByName(ctx, name)
	if err != nil {
		return false, fmt.Errorf("failed to check if access request exists: %w", err)
	}
	return exists, nil
}

func (s *service) FetchAccessRequests(ctx context.Context, builder accessrequest.FilterBuilder) ([]types.AccessRequest, error) {
	accessRequestFilter := builder.Build()
	return s.api.GetAccessRequests(ctx, accessRequestFilter)
}

func (s *service) FetchUsersWithoutSecrets(ctx context.Context) ([]models.User, error) {
	rawUsers, err := s.api.GetUsers(ctx, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	humanUsers := filterHumanUsers(rawUsers)
	return convertToUsers(humanUsers), nil
}

func (s *service) FetchUserAccessInfo(ctx context.Context, user *models.User) (*teleporttypes.UserAccessInfo, error) {
	req := types.AccessCapabilitiesRequest{
		User:             user.Username,
		RequestableRoles: true,
	}

	resp, err := s.api.GetAccessCapabilities(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get access capabilities: %w", err)
	}
	return &teleporttypes.UserAccessInfo{
		Roles:         resp.RequestableRoles,
		RequireReason: resp.RequireReason,
	}, nil
}

func (s *service) GetAccessRequestByName(ctx context.Context, name string) (*models.AccessRequest, error) {
	accessRequest, err := s.repo.GetAccessRequestByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get access request state: %w", err)
	}
	return accessRequest, nil
}

func (s *service) GetAccessRequestStateByName(ctx context.Context, name string) (string, error) {
	accessRequest, err := s.repo.GetAccessRequestByName(ctx, name)
	if err != nil {
		return "", fmt.Errorf("failed to get access request state: %w", err)
	}
	return accessRequest.State, nil
}

func (s *service) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	u, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get teleport user by username: %w", err)
	}
	return u, nil
}

func (s *service) StartMFAEventListener(ctx context.Context) error {
	go s.mfaEventPolling(ctx)
	return nil
}

func (s *service) SubmitAccessRequestState(ctx context.Context, builder accessrequest.UpdateBuilder) error {
	accessRequestState := builder.Build()
	return s.api.SetAccessRequestState(ctx, accessRequestState)
}

func (s *service) SubmitAccessRequest(ctx context.Context, builder accessrequest.CreateBuilder) (types.AccessRequest, error) {
	accessRequest := builder.Build()
	return s.api.CreateAccessRequestV2(ctx, accessRequest)
}

func (s *service) UpdateAccessRequestStateByName(ctx context.Context, accessRequest *models.AccessRequest) (*models.AccessRequest, error) {
	return s.repo.UpdateAccessRequestStateByName(ctx, accessRequest)
}

func filterHumanUsers(users []types.User) []types.User {
	var humanUsers []types.User
	for _, user := range users {
		copiedUser := user
		if !copiedUser.IsBot() {
			humanUsers = append(humanUsers, copiedUser)
		}
	}
	return humanUsers
}

func convertToUsers(users []types.User) []models.User {
	var result []models.User
	for _, user := range users {
		copiedUser := user
		result = append(result, models.User{
			Username: copiedUser.GetName(),
		})
	}
	return result
}

func (s *service) mfaEventPolling(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	lastEventTime := time.Now().UTC()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			events, _, err := s.api.SearchEvents(ctx, lastEventTime, time.Now().UTC(), "", nil, 100, types.EventOrderAscending, "")
			if err != nil {
				fmt.Printf("Error searching events: %v\n", err)
				continue
			}
			for _, event := range events {
				if s.isMFAAddEvent(event) && !s.isProcessed(event.GetID()) {
					_, err := s.handleMFAAddEvent(ctx, event)
					if err != nil {
						fmt.Printf("Error handling MFA add event: %v\n", err)
						continue
					}
					s.markAsProcessed(event.GetID())
				}
				if event.GetTime().After(lastEventTime) {
					lastEventTime = event.GetTime()
				}
			}
		}
	}
}

func (s *service) handleMFAAddEvent(ctx context.Context, event events.AuditEvent) (bool, error) {
	_ = ctx
	switch e := event.(type) {
	case *events.MFADeviceAdd:
		fmt.Printf("MFA Device Added: %s by %s\n", e.DeviceName, e.User)
		return true, nil
	default:
		return false, fmt.Errorf("unhandled event type: %T", e)
	}
}

func (s *service) isMFAAddEvent(event events.AuditEvent) bool {
	switch e := event.(type) {
	case *events.MFADeviceAdd:
		return true
	case *events.UserTokenCreate:
		return e.Name == "mfa"
	default:
		return event.GetType() == "mfa.add"
	}
}

func (s *service) isProcessed(eventID string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.processedEvents[eventID]
}

func (s *service) markAsProcessed(eventID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.processedEvents[eventID] = true

	if len(s.processedEvents) > MaxProcessedEvents {
		for id := range s.processedEvents {
			delete(s.processedEvents, id)
			if len(s.processedEvents) <= CleanupThreshold {
				break
			}
		}
	}
}
