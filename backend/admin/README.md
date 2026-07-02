<h1>vue-pure-admin</h1>

![GitHub license](https://img.shields.io/github/license/pure-admin/vue-pure-admin?style=flat)
![GitHub stars](https://img.shields.io/github/stars/pure-admin/vue-pure-admin?color=fa6470&style=flat)
![GitHub forks](https://img.shields.io/github/forks/pure-admin/vue-pure-admin?style=flat)

**中文** | [English](./README.en-US.md)

## 简介

`vue-pure-admin` 是一款开源免费且开箱即用的中后台管理系统模版。完全采用 `ECMAScript` 模块（`ESM`）规范来编写和组织代码，使用了最新的 `Vue3`、
`Vite`、`Element-Plus`、`TypeScript`、`Pinia`、`Tailwindcss` 等主流技术开发

## 研发理念

稳定中求创新，技术中见未来

## 精简版本（实际项目开发请用精简版本，提供 `非国际化` 、`国际化` 两个版本选择）

精简版本是基于 [vue-pure-admin](https://github.com/pure-admin/vue-pure-admin) 提炼出的架子，包含主体功能，更适合实际项目开发，打包后的大小在全局引入 [element-plus](https://element-plus.org) 的情况下仍然低于 `2.3MB`，并且会永久同步完整版的代码。开启 `brotli` 压缩和 `cdn` 替换本地库模式后，打包大小低于 `350kb`

[点我查看非国际化精简版本](https://github.com/pure-admin/pure-admin-thin)  
[点我查看国际化精简版本](https://github.com/pure-admin/pure-admin-thin/tree/i18n)

## 配套视频

[点我查看 UI 设计](https://www.bilibili.com/video/BV17g411T7rq)  
[点我查看快速开发教程](https://www.bilibili.com/video/BV1kg411v7QT)

## 配套保姆级文档

[点我查看 vue-pure-admin 文档](https://pure-admin.cn/)  
[点我查看 @pureadmin/utils 文档](https://pure-admin-utils.netlify.app)

## 高级服务

[点我查看详情](https://pure-admin.cn/pages/service/)

## `Tauri` 版本

[点我查看 Tauri 版本](https://github.com/pure-admin/tauri-pure-admin)

## `Electron` 版本

[点我查看 Electron 版本](https://github.com/pure-admin/electron-pure-admin)

## 预览

[点我查看预览](https://pure-admin.github.io/vue-pure-admin)

`PC` 端

<p align="center">
  <img alt="PureAdmin" src="https://xiaoxian521.github.io/hyperlink/img/vue-pure-admin/1.jpg">
  <br />
  <img alt="PureAdmin" src="https://xiaoxian521.github.io/hyperlink/img/vue-pure-admin/2.jpg">
</p>

暗色风格

<p align="center">
  <img alt="PureAdmin" src="https://xiaoxian521.github.io/hyperlink/img/vue-pure-admin/3.jpg">
  <br />
  <img alt="PureAdmin" src="https://xiaoxian521.github.io/hyperlink/img/vue-pure-admin/4.jpg">
</p>

移动端

<p align="center">
  <img alt="PureAdmin" src="https://xiaoxian521.github.io/hyperlink/img/vue-pure-admin/5.jpg">
</p>

### 使用 `Gitpod`

在 `Gitpod`（适用于 `GitHub` 的免费在线开发环境）中打开项目，并立即开始编码.

[![Open in Gitpod](https://gitpod.io/button/open-in-gitpod.svg)](https://gitpod.io/#https://github.com/pure-admin/vue-pure-admin)

## 安装使用

### 拉取代码

#### 推荐使用 `@pureadmin/cli` 脚手架

<br/>
<img src="https://xiaoxian521.github.io/hyperlink/gif/pure-admin-cli.gif" alt="pure-admin-cli" />

1. 全局安装

```bash
npm install -g @pureadmin/cli
```

2. 交互式选择模板并创建项目

```bash
pure create
```

[点我查看 @pureadmin/cli 脚手架详细用法](https://github.com/pure-admin/pure-admin-cli#pureadmincli)

#### 从 `GitHub` 上拉取

```bash
git clone https://github.com/pure-admin/vue-pure-admin.git
```

#### 从 `Gitee` 上拉取

```bash
git clone https://gitee.com/yiming_chang/vue-pure-admin.git
```

### 安装依赖

```bash
cd vue-pure-admin

pnpm install
```

### 启动平台

```bash
pnpm dev
```

### 项目打包

```bash
pnpm build
```

## Docker 支持

1. 自定义镜像名为 `vue-pure-admin` 的镜像（请注意下面命令末尾有一个点 `.` 表示使用当前路径下的 `Dockerfile` 文件，可根据实际情况指定路径）

```bash
docker build -t vue-pure-admin .
```

2. 端口映射并启动 `docker` 容器（`8080:80`：表示在容器中使用 `80` 端口，并将该端口转发到主机的 `8080` 端口；`pure-admin`：表示自定义容器名；`vue-pure-admin`：表示自定义镜像名）

```bash
docker run -dp 8080:80  --name pure-admin vue-pure-admin
```

操作完上面两个命令后，在浏览器打开 `http://localhost:8080` 即可预览

当然也可以通过 [Docker Desktop](https://www.docker.com/products/docker-desktop/) 可视化界面去操作 `docker` 项目，如下图

<p align="center">
  <img alt="docker-desktop" width="100%" src="https://xiaoxian521.github.io/hyperlink/img/docker-desktop.jpg">
</p>

## 更新日志

[CHANGELOG](./CHANGELOG.zh_CN.md)

## 如何贡献

非常欢迎您的加入！[提一个 Issue](https://github.com/pure-admin/vue-pure-admin/issues/new/choose) 或者提交一个 `Pull Request`

**Pull Request:**

1. Fork 代码!
2. 创建自己的分支: `git checkout -b feat/xxxx`
3. 提交您的修改: `git commit -am 'feat(function): add xxxxx'`
4. 推送您的分支: `git push origin feat/xxxx`
5. 提交`pull request`

## 特别代码贡献

非常感谢你们能深入了解源码并对 `pure-admin` 组织作出优秀贡献 ❤️

|                   **贡献人**                    |                                   **具体代码**                                   |
| :---------------------------------------------: | :------------------------------------------------------------------------------: |
|       [hb0730](https://github.com/hb0730)       |    [代码](https://github.com/pure-admin/vue-pure-admin/commits?author=hb0730)    |
|         [o-cc](https://github.com/o-cc)         |     [代码](https://github.com/pure-admin/vue-pure-admin/commits?author=o-cc)     |
| [yj-liuzepeng](https://github.com/yj-liuzepeng) | [代码](https://github.com/pure-admin/vue-pure-admin/commits?author=yj-liuzepeng) |
|   [skyline523](https://github.com/skyline523)   |  [代码](https://github.com/pure-admin/vue-pure-admin/commits?author=skyline523)  |
| [shark-lajiao](https://github.com/shark-lajiao) | [代码](https://github.com/pure-admin/vue-pure-admin/commits?author=shark-lajiao) |
|      [WitMiao](https://github.com/WitMiao)      |   [代码](https://github.com/pure-admin/vue-pure-admin/commits?author=WitMiao)    |
|     [QFifteen](https://github.com/QFifteen)     |   [代码](https://github.com/pure-admin/vue-pure-admin/commits?author=QFifteen)   |
|      [edgexie](https://github.com/edgexie)      |   [代码](https://github.com/pure-admin/vue-pure-admin/commits?author=edgexie)    |
|       [way-jm](https://github.com/way-jm)       |    [代码](https://github.com/pure-admin/vue-pure-admin/commits?author=way-jm)    |
|   [simple-hui](https://github.com/simple-hui)   |  [代码](https://github.com/pure-admin/vue-pure-admin/commits?author=simple-hui)  |
|   [tinysimple](https://github.com/tinysimple)   |  [代码](https://github.com/pure-admin/vue-pure-admin/commits?author=tinysimple)  |

## `Git` 贡献提交规范

参考 [vue](https://github.com/vuejs/vue/blob/dev/.github/COMMIT_CONVENTION.md) 规范 ([Angular](https://github.com/conventional-changelog/conventional-changelog/tree/master/packages/conventional-changelog-angular))

- `feat` 增加新功能
- `fix` 修复问题/BUG
- `style` 代码风格相关无影响运行结果的
- `perf` 优化/性能提升
- `refactor` 重构
- `revert` 撤销修改
- `test` 测试相关
- `docs` 文档/注释
- `chore` 依赖更新/脚手架配置修改等
- `workflow` 工作流改进
- `ci` 持续集成
- `types` 类型定义文件更改
- `wip` 开发中

## 浏览器支持

本地开发推荐使用 `Chrome`、`Edge`、`Firefox` 浏览器，作者常用的是最新版 `Chrome` 浏览器  
实际使用中感觉 `Firefox` 在动画上要比别的浏览器更加丝滑，只是作者用 `Chrome` 已经习惯了，看个人爱好选择吧  
更详细的浏览器兼容性支持请看 [Vue 支持哪些浏览器？](https://cn.vuejs.org/about/faq.html#what-browsers-does-vue-support) 和 [Vite 浏览器兼容性](https://cn.vitejs.dev/guide/build#browser-compatibility)

| [<img src="https://raw.githubusercontent.com/alrra/browser-logos/master/src/edge/edge_48x48.png" alt=" Edge" width="24px" height="24px" />](http://godban.github.io/browsers-support-badges/)</br>IE | [<img src="https://raw.githubusercontent.com/alrra/browser-logos/master/src/edge/edge_48x48.png" alt=" Edge" width="24px" height="24px" />](http://godban.github.io/browsers-support-badges/)</br>Edge | [<img src="https://raw.githubusercontent.com/alrra/browser-logos/master/src/firefox/firefox_48x48.png" alt="Firefox" width="24px" height="24px" />](http://godban.github.io/browsers-support-badges/)</br>Firefox | [<img src="https://raw.githubusercontent.com/alrra/browser-logos/master/src/chrome/chrome_48x48.png" alt="Chrome" width="24px" height="24px" />](http://godban.github.io/browsers-support-badges/)</br>Chrome | [<img src="https://raw.githubusercontent.com/alrra/browser-logos/master/src/safari/safari_48x48.png" alt="Safari" width="24px" height="24px" />](http://godban.github.io/browsers-support-badges/)</br>Safari |
| :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------: | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------: | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------: | :-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------: | :-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------: |
|                                                                                                不支持                                                                                                |                                                                                              最后两个版本                                                                                              |                                                                                                   最后两个版本                                                                                                    |                                                                                                 最后两个版本                                                                                                  |                                                                                                 最后两个版本                                                                                                  |

## 项目结构
admin/
├── build/                    # 构建配置目录（Vite 相关）
│   ├── cdn.ts               # CDN 资源配置
│   ├── compress.ts          # 压缩配置
│   ├── optimize.ts          # 构建优化配置
│   ├── plugins.ts           # Vite 插件配置
│   ├── static/              # 静态资源（构建时复制）
│   └── utils.ts             # 构建工具函数
│
├── locales/                  # 国际化语言包
│   ├── en.yaml              # 英文
│   └── zh-CN.yaml           # 简体中文
│
├── mock/                     # Mock 数据（开发环境模拟接口）
│   ├── login.ts             # 登录 Mock
│   ├── asyncRoutes.ts       # 异步路由 Mock
│   ├── system.ts            # 系统相关 Mock
│   └── refreshToken.ts      # Token 刷新 Mock
│
├── public/                   # 静态资源（直接复制到 dist，不经过编译）
│
├── src/                      # 源代码目录
│   ├── api/                  # API 接口层（与后端对接）
│   │   ├── album.ts         # 相册管理接口
│   │   ├── bookmark.ts      # 书签接口
│   │   ├── category.ts      # 分类接口
│   │   ├── chatter.ts       # 动态/聊天接口
│   │   ├── comment.ts       # 评论接口
│   │   ├── dashboard.ts     # 仪表盘数据接口
│   │   ├── friend-link.ts   # 友链接口
│   │   ├── message.ts       # 消息接口
│   │   ├── post.ts          # 文章接口
│   │   ├── project.ts       # 项目接口
│   │   ├── siteConfig.ts    # 站点配置接口
│   │   ├── system.ts        # 系统管理接口
│   │   ├── tag.ts           # 标签接口
│   │   ├── user.ts          # 用户接口
│   │   ├── visitor.ts       # 访客接口
│   │   └── routes.ts        # 路由接口
│
│   ├── assets/               # 静态资源（图片、字体等）
│   │   ├── iconfont/        # 字体图标
│   │   ├── login/           # 登录页相关图片
│   │   ├── status/          # 状态图标
│   │   ├── svg/             # SVG 图标
│   │   ├── table-bar/       # 表格工具栏图标
│   │   └── *.png/*.jpg      # 其他图片资源
│
│   ├── components/           # 全局复用组件
│   │   ├── ReAnimateSelector/ # 动画选择器组件
│   │   ├── ReAuth/            # 权限指令相关组件
│   │   ├── ReCol/             # 栅格布局组件
│   │   ├── ReCountTo/         # 数字滚动动画组件
│   │   ├── ReCropper/         # 图片裁剪组件
│   │   ├── ReCropperPreview/  # 裁剪预览组件
│   │   ├── ReDialog/          # 对话框组件封装
│   │   ├── ReDrawer/          # 抽屉组件封装
│   │   ├── ReFlicker/         # 闪烁动画组件
│   │   ├── ReIcon/            # 图标组件封装
│   │   ├── ReImageVerify/     # 图片验证码组件
│   │   ├── RePerms/           # 权限控制组件
│   │   ├── RePureTableBar/    # 表格工具栏组件
│   │   ├── ReQrcode/          # 二维码生成组件
│   │   ├── ReSegmented/       # 分段控制器组件
│   │   ├── ReText/            # 文本处理组件
│   │   ├── ReTypeit/          # 打字机效果组件
│   │   └── ReVxeTableBar/     # VxeTable 表格工具栏
│
│   ├── config/               # 全局配置文件
│   │   └── index.ts          # 环境配置、API 基础路径等
│
│   ├── directives/           # Vue 自定义指令
│   │   ├── auth/             # 权限指令
│   │   ├── copy/             # 复制指令
│   │   ├── longpress/        # 长按指令
│   │   ├── optimize/         # 性能优化指令
│   │   ├── perms/            # 权限指令（复数形式）
│   │   └── ripple/           # 水波纹效果指令
│
│   ├── layout/               # 布局组件（页面框架）
│   │   ├── components/       # 布局子组件
│   │   │   ├── lay-content/  # 主内容区域
│   │   │   ├── lay-footer/   # 页脚
│   │   │   ├── lay-frame/    # iframe 容器
│   │   │   ├── lay-navbar/   # 顶部导航栏
│   │   │   ├── lay-notice/   # 消息通知
│   │   │   ├── lay-panel/    # 面板
│   │   │   ├── lay-search/   # 搜索框
│   │   │   ├── lay-setting/  # 系统设置
│   │   │   ├── lay-sidebar/  # 侧边栏菜单
│   │   │   └── lay-tag/      # 标签页（Tab）
│   │   ├── hooks/            # 布局相关 Hooks
│   │   ├── index.vue         # 布局主入口
│   │   ├── frame.vue         # iframe 布局
│   │   ├── redirect.vue      # 重定向组件
│   │   └── types.ts          # 布局类型定义
│
│   ├── plugins/              # 第三方插件配置
│   │   ├── echarts.ts        # ECharts 图表库配置
│   │   ├── elementPlus.ts    # Element Plus UI 库配置
│   │   ├── i18n.ts           # 国际化插件配置
│   │   └── vxeTable.ts       # VxeTable 表格组件配置
│
│   ├── router/               # 路由配置（Vue Router）
│   │   ├── modules/          # 路由模块（懒加载路由）
│   │   ├── index.ts          # 路由入口
│   │   ├── enums.ts          # 路由枚举
│   │   └── utils.ts          # 路由工具函数
│
│   ├── store/                # 状态管理（Pinia）
│   │   ├── modules/          # Store 模块（各业务状态）
│   │   ├── index.ts          # Store 入口
│   │   ├── types.ts          # Store 类型定义
│   │   └── utils.ts          # Store 工具函数
│
│   ├── style/                # 全局样式
│   │   ├── reset.css         # CSS 重置
│   │   ├── variables.css     # CSS 变量定义
│   │   └── 其他样式文件
│
│   ├── utils/                # 工具函数库
│   │   ├── http/             # Axios 请求封装
│   │   ├── localforage/      # 离线存储封装
│   │   ├── progress/         # 请求进度条
│   │   ├── auth.ts           # 认证工具（Token 管理）
│   │   ├── chinaArea.ts      # 中国区域数据
│   │   ├── message.ts        # 消息提示封装
│   │   ├── print.ts          # 打印功能
│   │   ├── tree.ts           # 树形数据处理
│   │   ├── responsive.ts     # 响应式工具
│   │   └── sso.ts            # 单点登录相关
│
│   ├── views/                # 页面视图（业务页面）
│   │   ├── account-settings/ # 账号设置页面
│   │   ├── album/            # 相册管理
│   │   ├── bookmark/         # 书签管理
│   │   ├── category/         # 分类管理
│   │   ├── chatter/          # 动态/聊天管理
│   │   ├── comment/          # 评论管理
│   │   ├── dashboard/        # 仪表盘
│   │   ├── editor/           # 富文本编辑器
│   │   ├── empty/            # 空页面模板
│   │   ├── error/            # 错误页面（404/500）
│   │   ├── friend-link/      # 友链管理
│   │   ├── login/            # 登录/注册页面
│   │   ├── markdown/         # Markdown 编辑器
│   │   ├── message/          # 消息管理
│   │   ├── monitor/          # 监控页面
│   │   ├── permission/       # 权限管理
│   │   ├── post/             # 文章管理
│   │   ├── project/          # 项目管理
│   │   ├── site-config/      # 站点配置
│   │   ├── system/           # 系统管理（用户/角色/部门/菜单）
│   │   │   ├── user/         # 用户管理
│   │   │   ├── role/         # 角色管理
│   │   │   ├── dept/         # 部门管理
│   │   │   └── menu/         # 菜单管理
│   │   ├── tag/              # 标签管理
│   │   ├── visitor/          # 访客管理
│   │   └── welcome/          # 欢迎页
│
│   ├── App.vue               # 根组件
│   └── main.ts               # 应用入口文件
│
├── types/                    # TypeScript 类型定义
│   ├── global.d.ts          # 全局类型声明
│   ├── router.d.ts          # 路由类型扩展
│   ├── directives.d.ts      # 指令类型声明
│   └── global-components.d.ts # 全局组件类型声明
│
├── index.html                # HTML 模板
├── package.json              # 项目依赖配置
├── tsconfig.json             # TypeScript 配置
├── vite.config.ts            # Vite 构建配置
├── eslint.config.js          # ESLint 代码检查配置
├── stylelint.config.js       # Stylelint 样式检查配置
└── postcss.config.js         # PostCSS 配置
