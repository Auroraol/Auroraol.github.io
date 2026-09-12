> 基于 Chirpy 主题的 Jekyll 博客

发布新的文章，遵循以下步骤：

- 在`_posts`目录中添加新的Markdown文件。
- 执行命令生成静态网页：

```
bundle exec jekyll build
```

生成后的 `_site/index.html` 不要直接双击打开。页面使用站点根路径加载 CSS、JavaScript 和字体，`file://` 协议下这些资源会找不到，页面会显示成浏览器默认样式。请通过本地服务器预览：

```powershell
bundle exec jekyll serve --livereload
```

然后访问 <http://localhost:4000/>。

- 推送至GitHub仓库，更新后的网页会自动部署到GitHub Page上。

修改js:

```
# 1. 重新构建JavaScript
npx rollup -c --bundleConfigAsCjs

# 2. 重新构建网站
bundle exec jekyll build

# 3. 启动本地服务器（关键步骤）
bundle exec jekyll serve
```





```
   
```

