# 可选博客标题配置指南

## 功能说明

此功能允许您灵活配置博客的主标题显示方式。您可以选择：
1. 显示完整的标题（页面标题 | 站点标题）
2. 只显示页面标题（隐藏站点标题）
3. 使用自定义的站点标题

## 配置选项

### 1. 完全隐藏站点标题（推荐）
使用零宽空格字符实现视觉上的空标题：

```yaml
title: "&#8203;"  # 零宽空格字符 - 视觉上为空但技术上非空，避免SEO插件错误
```

### 2. 使用默认标题
保持现有的配置，显示完整的标题格式：

```yaml
title: My Personal Blog  # 显示为 "页面标题 | My Personal Blog"
```

### 3. 简洁模式
使用简短的站点标题：

```yaml
title: Blog  # 显示为 "页面标题 | Blog"
```

### 4. 自定义模式
使用个性化的站点标题：

```yaml
title: 我的技术笔记  # 显示为 "页面标题 | 我的技术笔记"
```

## 技术实现

### 核心修改
修改了 `_includes/head.html` 文件中的标题生成逻辑：

```html
<title>
  {%- unless page.layout == 'home' -%}
    {{ page.title }}
  {%- else -%}
    {{ site.title }}
  {%- endunless -%}
</title>
```

### 零宽空格解决方案
由于Jekyll SEO插件不能处理真正的空标题，我们使用了零宽空格字符（`&#8203;`）：
- 技术上非空，避免SEO插件错误
- 视觉上为空，达到隐藏效果
- 不会影响页面的实际显示

## 使用场景

### 场景1：极简主义
```yaml
title: "&#8203;"
```
效果：只显示页面标题，没有任何站点标识

### 场景2：专业博客
```yaml
title: John's Tech Blog
```
效果：所有页面标题都会加上站点名称后缀

### 场景3：个人品牌
```yaml
title: 张三的笔记
```
效果：体现个人特色

## SEO 影响

### 正面影响
- ✅ 搜索引擎仍能正确索引页面内容
- ✅ 页面标题结构清晰，有利于SEO
- ✅ 用户在搜索结果中能清楚识别页面来源

### 注意事项
- ⚠️ 首页标题会显示为零宽空格字符
- ⚠️ 建议在重要页面手动设置合适的标题
- ⚠️ 一致的标题格式有助于用户体验

## 测试验证

修改配置后，请按以下步骤验证：

```bash
# 1. 清理旧构建
bundle exec jekyll clean

# 2. 重新构建
bundle exec jekyll build

# 3. 启动本地服务器测试
bundle exec jekyll serve

# 4. 访问 http://localhost:4000 查看效果
```

## 故障排除

### 问题1：构建失败
**现象**：出现 "no implicit conversion of nil into String" 错误
**解决**：使用零宽空格字符 `&#8203;` 而不是真正的空字符串

### 问题2：标题显示异常
**现象**：页面标题格式不符合预期
**解决**：
1. 检查 `_includes/head.html` 文件是否正确修改
2. 确认 `_config.yml` 中的 `title` 配置
3. 重新构建项目

### 问题3：SEO标签异常
**现象**：页面源码中SEO相关信息缺失
**解决**：这是正常现象，因为我们禁用了SEO插件的标题生成功能，改用自己的逻辑

## 最佳实践

1. **保持一致性**：选择一种标题格式后尽量保持不变
2. **考虑长度**：站点标题建议控制在2-4个词范围内
3. **品牌化**：站点标题应该体现博客的主题或个人品牌
4. **测试充分**：修改后要在不同页面和设备上测试显示效果

## 示例配置

### 技术博客示例
```yaml
title: Developer Notes
tagline: Sharing coding experiences and technical insights
```

### 个人日记示例
```yaml
title: My Daily Thoughts
tagline: A personal space for reflection and learning
```

### 专业知识库示例
```yaml
title: Knowledge Base
tagline: Curated collection of useful information and tutorials
```

通过这种灵活的配置方式，您可以根据自己的需求和偏好来定制博客的标题显示效果。