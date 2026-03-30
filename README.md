> 基于 Chirpy 主题的 Jekyll 博客

发布新的文章，遵循以下步骤：

- 在`_posts`目录中添加新的Markdown文件。
- 执行命令生成静态网页：

```
bundle exec jekyll s 
```

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

