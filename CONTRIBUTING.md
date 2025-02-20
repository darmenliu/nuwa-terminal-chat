# 贡献指南

感谢您对 Nuwa Terminal Chat 项目的关注！我们非常欢迎您的贡献，无论是提交代码、报告问题还是改进文档。

## 开发流程

1. Fork 本仓库
2. 创建您的特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交您的改动 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建一个 Pull Request

## 代码提交规范

### Commit 信息格式

每个 commit 消息都包含一个标题和一个主体：

```
<type>(<scope>): <subject>

<body>
```

#### Type

- feat: 新功能
- fix: 修复问题
- docs: 文档修改
- style: 代码格式修改
- refactor: 代码重构
- test: 测试用例修改
- chore: 其他修改

### 代码风格

- 遵循 Go 标准代码风格
- 使用 `gofmt` 格式化代码
- 添加必要的注释
- 保持代码简洁明了

## 测试要求

- 为新功能编写单元测试
- 确保所有测试用例通过
- 运行 `make test` 进行测试

## 文档规范

- 保持文档的及时更新
- 使用清晰简洁的语言
- 提供必要的示例和说明
- 中英文文档保持同步

## 问题报告

提交问题时，请包含以下信息：

- 问题的详细描述
- 复现步骤
- 期望的结果
- 实际的结果
- 系统环境信息

## 开发环境设置

1. 确保已安装 Go 1.16 或更高版本
2. 克隆代码库
3. 安装依赖：`go mod download`
4. 运行测试：`make test`

## 许可证

通过提交 pull request，您同意您的贡献将按照项目的开源许可证进行授权。

## 联系我们

如果您有任何问题，欢迎通过 Issues 与我们联系。

感谢您的贡献！
