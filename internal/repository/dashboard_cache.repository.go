package repository

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

const dashboardCacheTTL = 5 * time.Minute

type DashboardCacheRepository struct {
	rdb    *redis.Client
	prefix string
}

func NewDashboardCacheRepository(rdb *redis.Client) *DashboardCacheRepository {
	return &DashboardCacheRepository{
		rdb:    rdb,
		prefix: os.Getenv("RDB_PREFIX"),
	}
}

func (dcr *DashboardCacheRepository) GetBalance(ctx context.Context, userId int) (float64, bool, error) {
	return dcr.getFloat(ctx, dcr.balanceKey(userId))
}

func (dcr *DashboardCacheRepository) SetBalance(ctx context.Context, userId int, balance float64) error {
	return dcr.setFloat(ctx, dcr.balanceKey(userId), balance)
}

func (dcr *DashboardCacheRepository) InvalidateBalance(ctx context.Context, userId int) error {
	return dcr.rdb.Del(ctx, dcr.balanceKey(userId)).Err()
}

func (dcr *DashboardCacheRepository) GetIncome(ctx context.Context, userId int) (float64, bool, error) {
	return dcr.getFloat(ctx, dcr.incomeKey(userId))
}

func (dcr *DashboardCacheRepository) SetIncome(ctx context.Context, userId int, income float64) error {
	return dcr.setFloat(ctx, dcr.incomeKey(userId), income)
}

func (dcr *DashboardCacheRepository) InvalidateIncome(ctx context.Context, userId int) error {
	return dcr.rdb.Del(ctx, dcr.incomeKey(userId)).Err()
}

func (dcr *DashboardCacheRepository) GetExpense(ctx context.Context, userId int) (float64, bool, error) {
	return dcr.getFloat(ctx, dcr.expenseKey(userId))
}

func (dcr *DashboardCacheRepository) SetExpense(ctx context.Context, userId int, expense float64) error {
	return dcr.setFloat(ctx, dcr.expenseKey(userId), expense)
}

func (dcr *DashboardCacheRepository) InvalidateExpense(ctx context.Context, userId int) error {
	return dcr.rdb.Del(ctx, dcr.expenseKey(userId)).Err()
}

func (dcr *DashboardCacheRepository) getFloat(ctx context.Context, key string) (float64, bool, error) {
	value, err := dcr.rdb.Get(ctx, key).Float64()
	if errors.Is(err, redis.Nil) {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}

	return value, true, nil
}

func (dcr *DashboardCacheRepository) setFloat(ctx context.Context, key string, value float64) error {
	return dcr.rdb.Set(ctx, key, value, dashboardCacheTTL).Err()
}

func (dcr *DashboardCacheRepository) balanceKey(userId int) string {
	return fmt.Sprintf("%s:dashboard:user:%d:balance", dcr.prefix, userId)
}

func (dcr *DashboardCacheRepository) incomeKey(userId int) string {
	return fmt.Sprintf("%s:dashboard:user:%d:income", dcr.prefix, userId)
}

func (dcr *DashboardCacheRepository) expenseKey(userId int) string {
	return fmt.Sprintf("%s:dashboard:user:%d:expense", dcr.prefix, userId)
}
