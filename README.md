> 基于 [Kirameku](https://github.com/Xinghongia/Kirameku) 后端 Golang 重写

---

## 项目结构

.
├── Kira_frontend/                 # 前端（Next.js 16 + React 19）
│   ├── app/                       # 页面路由
│   │   ├── posts/                 # 文章系统
│   │   ├── photowall/             # 照片墙
│   │   ├── timeline/              # 归档（时间河流）
│   │   ├── music/                 # 音乐播放器
│   │   ├── projects/              # 项目展示
│   │   ├── about/                 # 关于页面
│   │   ├── garden/                # 管理后台入口
│   │   └── ...
│   ├── components/                # UI 组件
│   │   ├── layout/                # 导航栏、背景、页脚
│   │   ├── providers/             # 主题、音乐、背景等 Provider
│   │   └── ui/                    # 通用组件（毛玻璃、动画等）
│   └── siteConfig.ts              # 站点全局配置
│
└── Kira_backend/                  # 后端（Golang + Gin）
    ├── cmd/                       # 程序入口
    ├── internal/
    │   ├── api/                   # RESTful API 控制器
    │   ├── models/                # GORM 数据模型
    │   ├── router/                # 路由注册
    │   ├── service/               # 业务逻辑层
    │   ├── middleware/            # 中间件（CORS、JWT 等）
    │   └── utils/                 # 工具函数（OSS 上传等）
    ├── admin/                     # 管理后台（Vue 3，打包后内嵌）
    ├── static/                    # 静态资源
    └── init_db_go.sql             # 数据库初始化脚本

## 快速开始

### 1. 后端

```bash
cd Kira_backend

# 配置环境变量
cp .env.example .env
# 编辑 .env，填入数据库、密钥、OSS 等配置

# 初始化数据库
psql -U postgres -d your_db -f init_db_go.sql

# 打包管理后台（首次运行）
cd admin
pnpm install
pnpm build
cd ..

# 启动服务
go run cmd/main.go
```

- API 文档：http://localhost:8000/docs（Swagger）
- 管理后台：http://localhost:8000/admin
2. 前端
```bash
cd Kira_frontend

# 安装依赖
pnpm install

# 开发模式
pnpm dev
# 访问 http://localhost:3000

# 生产构建
pnpm build
pnpm start
```
## 环境变量
``后端 .env``
```
# 数据库配置
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=kirameku

# JWT 配置
JWT_SECRET=your-secret-key-here

# 阿里云 OSS 配置
OSS_ACCESS_KEY_ID=your-access-key-id
OSS_ACCESS_KEY_SECRET=your-access-key-secret
OSS_ENDPOINT=oss-cn-beijing.aliyuncs.com
OSS_BUCKET=your-bucket-name

# 其他配置
APP_SECRET=your-app-secret
前端 .env.local
NEXT_PUBLIC_API_URL=http://localhost:8000
```
## License
MIT