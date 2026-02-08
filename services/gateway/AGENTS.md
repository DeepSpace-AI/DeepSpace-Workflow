1. 每当实现一个新的API服务，你需要在 `/docs/gateway-openapi.yaml` 中新增一个接口描述文档
2. 每当新增一个数据库模型，你需要在 `/pkg/db/migrate.go`中补充 AutoMigrate 的自动迁移模型
3. 注意 Golang 语言格式，完成后进行一次 format