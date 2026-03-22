# Sync to GitHub - 技能使用说明（优化版）

## ✨ 优化特性

1. **使用 GitHub API 直接推送** - 避免 git 命令权限问题和网络超时
2. **智能文件更新** - 自动检测文件是否存在，支持覆盖更新
3. **不破坏远程历史** - 不会像 `git push -f` 那样覆盖他人提交
4. **支持 Fine-grained Token** - 只需要 `contents:write` 权限即可
5. **可配置目标路径** - 支持推送到 `_posts` 等子目录

## 快速开始

### 1. 安装技能

```bash
# 在项目根目录执行
./install-skill.sh sync-to-github ./sync-to-github
```

### 2. 配置环境变量

在 `.env` 文件中添加：

```bash
SKILLS_CONFIG_JSON={
    "sync-to-github": {
        "enabled": true,
        "xProjectPath": "/home/node/.openclaw/workspace/x-project",
        "githubRepo": "your-username/your-repo",
        "githubToken": "ghp_xxxxxxxxxxxxx",
        "githubBranch": "main",
        "targetPath": "_posts",
        "categories": ["技术", "博客"],
        "tags": ["OpenClaw", "自动化"]
    }
}
```

### 3. 重启服务

```bash
docker compose restart openclaw-gateway
```

### 4. 验证安装

```bash
docker exec -it openclaw-gateway bash -c "su node -c 'openclaw skills list'"
```

## 使用方式

### 方式一：通过自然语言调用

对 AI 机器人说：

```
请将 a.md 的文章内容同步到 x 项目并推送到 GitHub
```

或者：

```
帮我把这篇文章转换格式后推到 GitHub 仓库
```

### 方式二：在代码中调用

```javascript
const SyncToGitHub = require('./sync-to-github');

const sync = new SyncToGitHub({
  xProjectPath: '/path/to/x-project',
  githubRepo: 'username/repo',
  githubToken: 'ghp_xxx',
  targetPath: '_posts'  // 推送到 _posts 目录
});

// 从文件读取
await sync.execute({
  sourcePath: './a.md',
  commitMessage: '同步文章：我的标题'
});

// 或直接提供内容
await sync.execute({
  content: '# 我的文章\n\n这是正文内容...',
  title: '我的文章标题',
  categories: ['分类 1', '分类 2'],
  tags: ['标签 1', '标签 2']
});
```

## 输出示例

### 输入 (a.md)

```markdown
# 我的技术文章

这是一篇关于 OpenClaw 的技术文章。

## 正文内容

详细介绍如何使用 OpenClaw...
```

### 输出 (2025-11-04-我的技术文章.md)

```markdown
---
title: 我的技术文章
date: 2025-11-04 11:00:00 +0800
categories: [技术，博客]
tags: [OpenClaw, 自动化，同步]
---
# 我的技术文章

这是一篇关于 OpenClaw 的技术文章。

## 正文内容

详细介绍如何使用 OpenClaw...
```

## 配置参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `xProjectPath` | string | 否 | 本地项目路径（用于缓存） |
| `githubRepo` | string | **是** | GitHub 仓库名，格式：username/repo |
| `githubToken` | string | **是** | GitHub Personal Access Token |
| `githubBranch` | string | 否 | 目标分支，默认 main |
| `targetPath` | string | 否 | 目标子目录，如 `_posts` |
| `categories` | Array | 否 | 默认分类 |
| `tags` | Array | 否 | 默认标签 |

## 执行参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `sourcePath` | string | 否 | 源 markdown 文件路径 |
| `content` | string | 否 | markdown 内容（与 sourcePath 二选一） |
| `title` | string | 否 | 文章标题（不填则自动提取） |
| `date` | Date/string | 否 | 文章日期（不填则使用当前时间） |
| `categories` | Array | 否 | 分类数组（不填则使用默认值） |
| `tags` | Array | 否 | 标签数组（不填则使用默认值） |
| `commitMessage` | string | 否 | Git 提交信息（不填则自动生成） |

## GitHub Token 配置

### Fine-grained Token（推荐）

1. 访问 https://github.com/settings/tokens?type=beta
2. 创建新的 Fine-grained token
3. 选择你的仓库
4. 权限设置：
   - **Contents**: Read and write

### Classic Token

1. 访问 https://github.com/settings/tokens
2. 生成新的 classic token
3. 勾选 `repo` 权限

## 注意事项

1. **文件命名**：自动使用 `年-月-日-标题.md` 格式
2. **标题提取**：优先从 `# 标题` 提取，如果没有则使用前几行文字
3. **文件更新**：如果文件已存在，会自动更新而不创建重复文件
4. **本地缓存**：`xProjectPath` 是可选的，用于本地缓存转换后的文件

## 故障排查

### 问题 1：GitHub API 返回 404

**原因**：文件路径错误或仓库不存在

**解决**：
- 检查 `githubRepo` 格式是否正确（username/repo）
- 确认仓库是否存在

### 问题 2：GitHub API 返回 403

**原因**：Token 权限不足

**解决**：
- 检查 Token 是否正确
- 确保 Token 有 `contents:write` 权限（Fine-grained）或 `repo` 权限（Classic）

### 问题 3：文件格式转换失败

**原因**：缺少依赖包

**解决**：
```bash
docker exec -it openclaw-gateway bash -c "su node -c 'cd ~/.openclaw/workspace/skills/sync-to-github && npm install'"
```

## 高级用法

### 自定义分类和标签

可以在调用时指定：

```
请将 a.md 同步到 GitHub，分类设为 ["AI","技术"]，标签设为 ["LLM","自动化"]
```

### 指定日期

```
请同步 a.md 到 GitHub，日期设为 2025-11-04
```

### 推送到不同目录

```javascript
const sync = new SyncToGitHub({
  githubRepo: 'username/blog',
  githubToken: 'ghp_xxx',
  targetPath: '_posts'  // 推送到 _posts 子目录
});
```

### 批量处理

```javascript
const files = ['a.md', 'b.md', 'c.md'];
for (const file of files) {
  await sync.execute({ sourcePath: file });
}
```

## 返回结果

```javascript
{
  success: true,
  message: '同步完成',
  filename: '2025-11-04-我的文章.md',
  title: '我的文章',
  savedPath: '/path/to/x-project/2025-11-04-我的文章.md',
  pushResult: {
    success: true,
    message: '成功推送到 GitHub: username/repo',
    pushed: true,
    branch: 'main',
    commit: 'Sync: 我的文章',
    filePath: '_posts/2025-11-04-我的文章.md',
    htmlUrl: 'https://github.com/username/repo/blob/main/_posts/2025-11-04-我的文章.md',
    sha: '86eb0af75d6bb2c1942f55b2ad5d666a28828ef4'
  }
}
```

## 相关文件

- [`index.js`](./index.js) - 主要实现代码
- [`package.json`](./package.json) - 依赖配置
- [`README.md`](./README.md) - 技能说明
- [`skill.config.example.json`](./skill.config.example.json) - 配置示例

## 技术支持

如有问题，请查看 OpenClaw 日志：

```bash
docker logs -f openclaw-gateway | grep sync-to-github
```
