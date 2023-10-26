package auth_manager

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"ketalk-api/jwt"
	"ketalk-api/pkg/manager/auth/repository"
	"ketalk-api/pkg/manager/port"
	"ketalk-api/pkg/provider"
	"ketalk-api/pkg/provider/model"
	"time"

	"github.com/google/uuid"
)

type Token struct {
	AuthToken    string
	RefreshToken string
}

type authManager struct {
	authRepository repository.Repository
	userPort       port.UserPort
	provider       provider.ProviderClient
	jwtConfig      jwt.Config
}

func NewAuthManager(authRepository repository.Repository, userPort port.UserPort, provider provider.ProviderClient, jwtConfig jwt.Config) AuthManager {
	return &authManager{
		authRepository,
		userPort,
		provider,
		jwtConfig,
	}
}

func (m *authManager) SignupOrLogin(ctx context.Context, req SignupOrLoginRequest) (*SignupOrLoginResponse, error) {
	// create user
	providerToken := model.ProviderToken{}

	userDetails, err := m.provider.VerifyIDToken(ctx, &providerToken)
	var username string
	var email string
	var image *string
	if err == nil {
		username = userDetails.FirstName
		email = userDetails.Email
		image = &userDetails.Image
	} else if req.SignUpDetails != nil {
		username = req.SignUpDetails.UserName
		email = req.SignUpDetails.Email
	} else {
		return nil, fmt.Errorf("invalid request")
	}

	createUserReq := port.CreateOrGetUserRequest{
		Username: username,
		Email:    email,
		Image:    image,
	}
	if req.SignUpDetails != nil {
		createUserReq.Password = &req.SignUpDetails.Password
	}
	user, err := m.userPort.CreateOrGetUser(ctx, createUserReq)
	if err != nil {
		return nil, err
	}

	if user.Password != nil && req.SignUpDetails != nil && *user.Password != req.SignUpDetails.Password {
		return nil, fmt.Errorf("invalid password")
	}

	// create token
	token, err := m.issueTokens(ctx, user.ID, nil, req.DeviceID, req.DeviceOS)
	if err != nil {
		return nil, err
	}

	// return
	return &SignupOrLoginResponse{
		Id:           user.ID,
		UserName:     user.Username,
		Email:        user.Email,
		Image:        user.Image,
		AuthToken:    token.AuthToken,
		RefreshToken: token.RefreshToken,
	}, nil
}

func (h *authManager) issueTokens(ctx context.Context, userId uuid.UUID, oldToken *string, deviceID string, deviceOS string) (Token, error) {
	tokens, err := h.createTokens(ctx, userId)
	if err != nil {
		return tokens, err
	}

	if oldToken != nil {
		err = h.authRepository.DeleteRefreshToken(ctx, *oldToken)
		if err != nil {
			return tokens, err
		}
	}
	input := repository.DeviceRefreshToken{
		DeviceID:             deviceID,
		DeviceOS:             deviceOS,
		UserId:               userId,
		RefreshToken:         tokens.RefreshToken,
		RefreshTokenExpiryAt: time.Now().Add(h.jwtConfig.RefreshTokenExpiryDuration),
	}
	err = h.authRepository.AddRefreshToken(ctx, &input)
	return tokens, err
}

func (h *authManager) createTokens(_ context.Context, userId uuid.UUID) (Token, error) {
	authToken, err := jwt.IssueToken(h.jwtConfig, userId, map[string]interface{}{
		"id": userId,
	})
	if err != nil {
		return Token{}, fmt.Errorf("could not issue JWT for %v, err: %+v", userId, err)
	}

	refreshToken, err := GenerateRandomHex(128)
	if err != nil {
		return Token{}, fmt.Errorf("could not generate refresh token for %v", userId)
	}

	return Token{
		AuthToken:    authToken,
		RefreshToken: refreshToken,
	}, nil
}

func GenerateRandomHex(length int) (string, error) {
	byteLength := length / 2
	randBytes := make([]byte, byteLength)
	_, err := rand.Read(randBytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(randBytes), nil
}
