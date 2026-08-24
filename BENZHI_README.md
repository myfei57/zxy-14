# MaskHub

MaskHub 是数据脱敏与隐私保护网关：字段识别、敏感分类、脱敏策略发布与回滚、
脱敏执行与配额控制、结果审计，全部操作留审计。

## 功能

- 数据源接入与字段识别
- 敏感分类与分类模型版本
- 脱敏策略草稿、发布、回滚
- 脱敏执行与字段替换
- 租户配额控制
- 脱敏结果审计
- 控制台页面与 HTTP API

## 运行

```bash
go build -mod=vendor -o maskhub ./cmd/server
./maskhub -addr :8080 -data ./data
```

打开 http://localhost:8080/ 查看总览，/console/sources、/console/policies、
/console/masking、/console/audit 分别查看数据源、策略、脱敏与审计页面。

## 测试

```bash
go test -mod=vendor ./...
```
