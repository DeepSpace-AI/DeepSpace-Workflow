package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"deepspace-worker/internal/config"
	"deepspace-worker/internal/service/email"

	"github.com/redis/go-redis/v9"
)

func main() {
	log.Println("启动邮件队列 Worker...")

	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("配置校验失败: %v", err)
	}

	emailService, err := email.New(cfg)
	if err != nil {
		log.Fatalf("初始化邮件服务失败: %v", err)
	}

	client := emailService.RedisClient()
	if client == nil {
		log.Fatal("Redis 客户端不可用")
	}

	ctx := context.Background()
	for {
		result, err := client.BRPop(ctx, cfg.PollTimeout, cfg.RedisQueueKey).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				continue
			}
			log.Printf("队列拉取失败: %v", err)
			continue
		}
		if len(result) < 2 {
			continue
		}

		payload := result[1]
		if err := handleQueueItem(ctx, emailService, client, cfg, payload); err != nil {
			log.Printf("处理邮件任务失败: %v", err)
		}
	}
}

func handleQueueItem(ctx context.Context, svc *email.Service, client *redis.Client, cfg *config.Config, payload string) error {
	var item email.QueueItem
	if err := json.Unmarshal([]byte(payload), &item); err != nil {
		log.Printf("解析邮件任务失败: %v", err)
		return nil
	}

	retryMax := cfg.RetryMax
	if retryMax <= 0 {
		retryMax = 1
	}

	var lastErr error
	for attempt := 1; attempt <= retryMax; attempt++ {
		if err := svc.Send(ctx, item.Email); err != nil {
			lastErr = err
			log.Printf("发送失败(第%d次): %v", attempt, err)
			time.Sleep(time.Duration(attempt) * time.Second)
			continue
		}
		return nil
	}

	if lastErr != nil {
		_, err := client.LPush(ctx, cfg.RedisDeadKey, payload).Result()
		if err != nil {
			return err
		}
	}

	return lastErr
}
