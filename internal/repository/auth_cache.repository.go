package repository

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type AuthCacheRepository struct {
	rdb    *redis.Client
	prefix string
}

func NewAuthCacheRepository(rdb *redis.Client) *AuthCacheRepository {
	prefix := os.Getenv("RDB_PREFIX")

	return &AuthCacheRepository{
		rdb:    rdb,
		prefix: prefix,
	}
}

func (acr *AuthCacheRepository) SaveToken(ctx context.Context, tokenHash string, userId int, expiresAt time.Time) error {
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return nil
	}

	return acr.rdb.Set(ctx, acr.userTokenKey(userId, tokenHash), strconv.Itoa(userId), ttl).Err()
}

func (acr *AuthCacheRepository) DeleteToken(ctx context.Context, tokenHash string, userId int) error {
	return acr.rdb.Del(ctx, acr.userTokenKey(userId, tokenHash)).Err()
}

func (acr *AuthCacheRepository) IsTokenActive(ctx context.Context, tokenHash string, userId int) (bool, error) {
	exists, err := acr.rdb.Exists(ctx, acr.userTokenKey(userId, tokenHash)).Result()
	if err != nil {
		return false, err
	}

	return exists == 1, nil
}

// InvalidateUserTokens deletes all tokens associated with the given userId, effectively logging out the user from all sessions.
// This method uses the SCAN command to find all keys matching the pattern for the user's tokens and deletes them.
func (acr *AuthCacheRepository) InvalidateUserTokens(ctx context.Context, userId int) error {
	var cursor uint64
	pattern := acr.userTokenKey(userId, "*")

	for {
		keys, nextCursor, err := acr.rdb.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return err
		}

		for _, key := range keys {
			if err := acr.rdb.Del(ctx, key).Err(); err != nil {
				return err
			}
		}

		if nextCursor == 0 {
			break
		}
		cursor = nextCursor
	}

	return nil
}

func (acr *AuthCacheRepository) SavePasswordResetToken(ctx context.Context, userId int, tokenHash string, ttl time.Duration) error {
	if ttl <= 0 {
		return nil
	}

	return acr.rdb.Set(ctx, acr.passwordResetTokenKey(tokenHash), strconv.Itoa(userId), ttl).Err()
}

func (acr *AuthCacheRepository) ConsumePasswordResetToken(ctx context.Context, tokenHash string) (int, error) {
	userId, err := acr.rdb.GetDel(ctx, acr.passwordResetTokenKey(tokenHash)).Int()
	if err != nil {
		return 0, err
	}

	return userId, nil
}

func (acr *AuthCacheRepository) userTokenKey(userId int, tokenHash string) string {
	return fmt.Sprintf("%s:auth:user:%d:token:%s", acr.prefix, userId, tokenHash)
}

func (acr *AuthCacheRepository) passwordResetTokenKey(tokenHash string) string {
	return fmt.Sprintf("%s:auth:password-reset:%s", acr.prefix, tokenHash)
}
