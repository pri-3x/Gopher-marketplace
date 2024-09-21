package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gopher-market/services/auth-service/internal/models"
	"time"

	"github.com/go-redis/redis/v8"
)

type AuthRepository struct {
	redisClient *redis.Client
}

func NewAuthRepository(redisClient *redis.Client) *AuthRepository {
	return &AuthRepository{redisClient: redisClient}
}

func (r *AuthRepository) SaveUser(ctx context.Context, user *models.User) error {
	userJSON, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}

	err = r.redisClient.Set(ctx, fmt.Sprintf("user:%s", user.Username), userJSON, 0).Err()
	if err != nil {
		return fmt.Errorf("failed to save user to Redis: %w", err)
	}

	return nil
}

func (r *AuthRepository) GetUser(ctx context.Context, username string) (*models.User, error) {
	userJSON, err := r.redisClient.Get(ctx, "user:"+username).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	var user models.User
	err = json.Unmarshal([]byte(userJSON), &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *AuthRepository) SaveToken(ctx context.Context, username, token string) error {
	return r.redisClient.Set(ctx, "token:"+username, token, 24*time.Hour).Err()
}

func (r *AuthRepository) GetToken(ctx context.Context, username string) (string, error) {
	return r.redisClient.Get(ctx, "token:"+username).Result()
}
