package jwt_auth

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/asyauqi15/payslip-system/internal"
	"github.com/asyauqi15/payslip-system/internal/constant"
	"github.com/asyauqi15/payslip-system/internal/entity"
	v1 "github.com/asyauqi15/payslip-system/pkg/openapi/v1"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/spf13/cast"
)

const (
	jwtUserIDKey       = "user_id"
	jwtUserUsernameKey = "username"
	jwtUserRoleKey     = "user_role"
)

type JWTAuthentication struct {
	accessTokenAuth      *jwtauth.JWTAuth
	refreshTokenAuth     *jwtauth.JWTAuth
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

type TokenClaims struct {
	UserID    string `json:"user_id"`
	UserEmail string `json:"user_email"`
	UserRole  string `json:"user_role"`
}

func NewJWTAuthentication(config internal.HTTPServerConfig) (*JWTAuthentication, error) {
	at, err := config.GetAccessTokenSecret()
	if err != nil {
		return nil, fmt.Errorf("failed to decode access token secret: %w", err)
	}

	rt, err := config.GetRefreshTokenSecret()
	if err != nil {
		return nil, fmt.Errorf("failed to decode refresh token secret: %w", err)
	}

	return &JWTAuthentication{
		accessTokenAuth:      jwtauth.New("HS256", at, nil),
		refreshTokenAuth:     jwtauth.New("HS256", rt, nil),
		accessTokenDuration:  config.AccessTokenDuration,
		refreshTokenDuration: config.RefreshTokenDuration,
	}, nil
}

func (ja *JWTAuthentication) Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := jwtauth.VerifyRequest(ja.accessTokenAuth, r, jwtauth.TokenFromHeader)
		if err != nil {
			handleErr(w, r, err.Error(), http.StatusUnauthorized)
			return
		}

		claims, err := token.AsMap(r.Context())
		if err != nil {
			handleErr(w, r, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		{
			ctx = context.WithValue(ctx, constant.ContextKeyUserID, cast.ToString(claims[jwtUserIDKey]))
			ctx = context.WithValue(ctx, constant.ContextKeyUsername, cast.ToString(claims[jwtUserUsernameKey]))
			ctx = context.WithValue(ctx, constant.ContextKeyUserRole, cast.ToString(claims[jwtUserRoleKey]))
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (ja *JWTAuthentication) GenerateAccessToken(ctx context.Context, user *entity.User) (string, time.Time, error) {
	expiredAt := time.Now().Add(ja.accessTokenDuration)
	token, err := generateToken(ctx, ja.accessTokenAuth, user, expiredAt)
	if err != nil {
		return "", time.Time{}, err
	}
	return token, expiredAt, err
}

func (ja *JWTAuthentication) GenerateRefreshToken(ctx context.Context, user *entity.User) (string, time.Time, error) {
	expiredAt := time.Now().Add(ja.refreshTokenDuration)
	token, err := generateToken(ctx, ja.refreshTokenAuth, user, expiredAt)
	if err != nil {
		return "", time.Time{}, err
	}
	return token, expiredAt, err
}

func (ja *JWTAuthentication) ParseAccessToken(ctx context.Context, token string) (TokenClaims, error) {
	return parseToken(ctx, ja.accessTokenAuth, token)
}

func (ja *JWTAuthentication) ParseRefreshToken(ctx context.Context, token string) (TokenClaims, error) {
	return parseToken(ctx, ja.refreshTokenAuth, token)
}

func parseToken(ctx context.Context, jwt *jwtauth.JWTAuth, token string) (TokenClaims, error) {
	refreshToken, err := jwtauth.VerifyToken(jwt, token)
	if err != nil {
		slog.ErrorContext(ctx, "error when verify token", "error", err)
		return TokenClaims{}, err
	}

	claims, err := refreshToken.AsMap(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "error when get token claims", "error", err)
		return TokenClaims{}, err
	}

	return TokenClaims{
		UserID:    cast.ToString(claims[jwtUserIDKey]),
		UserEmail: cast.ToString(claims[jwtUserUsernameKey]),
		UserRole:  cast.ToString(claims[jwtUserRoleKey]),
	}, nil
}

func generateToken(ctx context.Context, jwt *jwtauth.JWTAuth, user *entity.User, expiredAt time.Time) (string, error) {
	claims := map[string]interface{}{
		jwtUserIDKey:       user.ID,
		jwtUserUsernameKey: user.Username,
		jwtUserRoleKey:     user.Role,
	}
	jwtauth.SetExpiry(claims, expiredAt)
	jwtauth.SetIssuedNow(claims)

	_, tokenString, err := jwt.Encode(claims)
	if err != nil {
		slog.ErrorContext(ctx, "error when encode token claims", "error", err)
		return "", err
	}
	return tokenString, nil
}

func handleErr(w http.ResponseWriter, r *http.Request, errMessage string, httpCode int) {
	resp := &v1.DefaultErrorResponse{}
	resp.Error.Message = errMessage

	render.Status(r, httpCode)
	render.JSON(w, r, resp)
}
