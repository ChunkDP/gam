 # 配置文档

## 环境配置
系统支持多环境配置，配置文件位于 `config/environments/` 目录：
- `config.yaml`: 基础配置
- `config.dev.yaml`: 开发环境
- `config.test.yaml`: 测试环境
- `config.prod.yaml`: 生产环境

## 配置项说明

### Server 服务器配置
- `port`: 服务器端口号
- `mode`: 运行模式 (debug/release)

### Database 数据库配置
- `driver`: 数据库驱动 (mysql)
- `host`: 数据库主机地址
- `port`: 数据库端口
- `username`: 数据库用户名
- `password`: 数据库密码
- `dbname`: 数据库名称
- `max_idle_conns`: 最大空闲连接数
- `max_open_conns`: 最大打开连接数
- `conn_max_lifetime`: 连接最大生命周期(秒)

### JWT 配置
- `secret_key`: JWT密钥
- `expire_time`: 过期时间(小时)
- `issuer`: 签发者

### Redis 配置
- `host`: Redis主机地址
- `port`: Redis端口
- `password`: Redis密码
- `db`: 数据库编号

### CORS 跨域配置
- `allowed_origins`: 允许的源
- `allowed_methods`: 允许的HTTP方法
- `allowed_headers`: 允许的请求头


### Log 日志配置
- `level`: 日志级别(debug/info/warn/error)
- `filename`: 日志文件路径
- `max_size`: 单个日志文件最大大小(MB)
- `max_age`: 日志保留天数
- `max_backups`: 保留的旧日志文件数量
- `compress`: 是否压缩旧日志

## 环境变量覆盖
所有配置项都可以通过环境变量覆盖，环境变量名使用大写，以下划线连接，例如：
- `SERVER_PORT`
- `DATABASE_PASSWORD`
- `JWT_SECRET_KEY`