---
title: Chirpy主题的个性化自定义
date: 2023-12-20 11:00:00 +0800
categories: [jekyll, Chirpy使用]
tags: [Chirpy使用]
---

## 更改\_data文件夹中contact.yml文件

`contact.yml`直接关联着网站左下角的链接，只保留github和email。后面自行可以添加

```yaml
- type: github
  icon: "fab fa-github"

#- type: twitter
#  icon: "fa-brands fa-x-twitter"

- type: email
  icon: "fas fa-envelope"
  
- type: bluesky
  icon: 'fa-solid fa-screwdriver-wrench'
  url: 'https://jinhanlei.github.io/posts/Jekyll+%E6%A8%A1%E6%9D%BF-%E6%90%AD%E5%BB%BA%E4%B8%AA%E4%BA%BA%E5%8D%9A%E5%AE%A2-%E4%BA%8C-%E6%A8%A1%E6%9D%BF%E4%BF%AE%E6%94%B9/'
```

## 更改\_data文件夹中的share.yml

`share.yml`是在点进文章之后最底下的分享链接，也注释掉一些不认识的网站，只留下weibo。

```yaml
platforms:
  - type: Weibo
    icon: "fab fa-weibo"
    link: "http://service.weibo.com/share/share.php?title=TITLE&url=URL"
```

## 更改\_tabs/about.md文件

`about.md`关联着侧边栏的about，改成自己的导航。

```yaml
---
icon: fas fa-info-circle
order: 4
---

| 内容                                               | 网址                                                                                                | 说明                                                                               |
| :------------------------------------------------- | :-------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------- |
| **在线PS**                                         | [Photopea](https://www.photopea.com/ "Photopea")                                                    | 非常实用的一款在线Photoshop工具，无需注册登录                                      |
|                                                    | [uupoop](https://www.uupoop.com "uupoop")                                                           | 简单看作国内版的「Photopea」。下载需要注册                                         |
| **一键抠图**                                       | [remove.bg](https://www.remove.bg/zh "remove.bg")                                                   | 普通分辨率图片下载免费（无需注册）高分辨率下载需注册（不免费）。中文版容易打不开。 |
| **OCR**                                            | [白描网页版](https://web.baimiaoapp.com/)                                             | 白描网页版                                                                         |
| **AI对话**                                         | [https://codenews.cc/chatgpt](https://codenews.cc/chatgpt )                                         | 下面有AI导航网站                                                                   |
|                                                    | [https://ai.qdymys.cn/#/](https://ai.qdymys.cn/#/ )                                                 | 有额度限制                                                                         |
|                                                    | [https://www.zaiwen.top/](https://www.zaiwen.top/ )                                                 | 高效问答无限制，非GPT模型，有多种模型                                              |
|                                                    | [https://chatgai.lovepor.cn/](https://chatgai.lovepor.cn/ )                                         | 有额度限制                                                                         |
|                                                    | [https://chat10.free2gpt.xyz/](https://chat10.free2gpt.xyz/ )                                       | 完全免费 无需魔法 无需登录 120次/天 免费                                           |
| **网页设计**                                       | [https://www.jyshare.com/front-end/61/](https://www.jyshare.com/front-end/61/ )                     | 菜鸟工具                                                                           |
|                                                    | [https://htmlpage.cn/builder/](https://htmlpage.cn/builder/ )                                       | 无需注册，自己设计，可导入导出源码                                                 |
|                                                    | [HTML编辑器](https://www.58html.com/html/)                                         | 可以QQ登录，右键编辑CSS属性，不可以导入源码                                        |
| **Markdown**                                       | [https://hintryer.github.io/tool/Markdown.html](https://hintryer.github.io/tool/Markdown.html )     | 自制 Markdown 在线编辑器                                                           |
|                                                    | [https://pandao.github.io/editor.md/](https://pandao.github.io/editor.md/ )                         | 开源Editor.md，有各种快捷按键，打开较慢                                            |
|                                                    | [https://www.jyshare.com/front-end/712/](https://www.jyshare.com/front-end/712/ )                   | 菜鸟工具，简洁无广告，打开较快，引用开源模块                                       |
|                                                    | [https://www.lddgo.net/string/markdown-editor](https://www.lddgo.net/string/markdown-editor )       | 有其他东西，不影响使用                                                             |
| **HTML2MD**                                        | [https://www.helloworld.net/html2md](https://www.helloworld.net/html2md )                           | 开源，也能通过网址直接生成                                                         |
|                                                    | [https://devtool.tech/html-md](https://devtool.tech/html-md )                                       | HTML转Markdown工具，也能通过网址直接生成                                           |
|                                                    | [https://md.weiyan.cc/](https://md.weiyan.cc/ )                                                     | 支持自定义样式的 Markdown 编辑器，支持导出 pdf 和 markdown                         |
| **源阅读**                                         | [https://hintryer.github.io/tool/yuedu.html](https://hintryer.github.io/tool/yuedu.html )           | 自制 源阅读编辑器                                                                  |
|                                                    | [https://rc.real9.cn/](https://rc.real9.cn/ )                                                       | 源阅读书源转换工具                                                                 |
|                                                    | [https://yuedu.juhua.fun/tool](https://yuedu.juhua.fun/tool )                                       | 美化发现工具                                                                       |
|                                                    | [https://yuedu.miaogongzi.net/gx.html](https://yuedu.miaogongzi.net/gx.html )                       | 猫公子源                                                                           |
|                                                    | [https://www.yckceo.com/yuedu/shuyuan/index.html](https://www.yckceo.com/yuedu/shuyuan/index.html ) | 源仓库                                                                             |
| **音乐下载**                                       | [https://iui.su/2217/](https://iui.su/2217/)                                                        | 来源不死鸟 • 2024 年 09 月 20 日                                                   |
|                                                    | [https://freemp3cn.com](https://freemp3cn.com )                                                     | 无限制，搜索下载简单，只一首一首地下载                                             |
|                                                    | [https://tools.liumingye.cn/music](https://tools.liumingye.cn/music )                               | MyFreeMP3，已关闭新用户注册。无限制，搜索下载简单                                  |
|                                                    | [https://www.toomic.com](https://www.toomic.com )                                                   | 提供歌词下载，可下低品质音乐                                                       |
| [测试1](https://hintryer.github.io/ceshi/ "测试1") | [测试2](https://hintryer.github.io/ceshi2/ "测试2")                                                 | [测试3](https://hintryer.github.io/ceshi3/ "测试3")                                |

|                                                                |                                                                                      |                                               |                                                                        |                                                                                                   |
| :------------------------------------------------------------: | :----------------------------------------------------------------------------------: | :-------------------------------------------: | :--------------------------------------------------------------------: | :-----------------------------------------------------------------------------------------------: |
| [短网址u3v](https://www.sojson.com/dwz/batch.html "短网址u3v") | [自制HTML代码格式化](https://hintryer.github.io/tool/code.html "自制HTML代码格式化") | [正则测试](https://regexr-cn.com/ "正则测试") | [HTML运行与美化](https://codebeautify.org/htmlviewer "HTML运行与美化") | [jsbeautify库属性](https://blog.csdn.net/u012732909/article/details/121791340 "jsbeautify库属性") |
|                                                                |                                                                                      |                                               |                                                                        |                                                                                                   |
```

##  自定义左侧标签

在`_tabs/`中添加以下MD文件，注意`order: 5`中的数字对应

```yaml
---
# the default layout is 'page'
title: 工具
icon: fas fa-info-circle
order: 5
---

 **在线PS**  [Photopea](https://www.photopea.com/ "Photopea")

```

### assets {#assets}

assets/img/favicons/下存放有头像，需要替换成自己的。

事先准备头像图片，前往[Real Favicon Generator](https://realfavicongenerator.net/)上传，上传后点击`Generate your Favicons and HTML code`即可，解压后删除以下两个文件：

- `browserconfig.xml`
- `site.webmanifest`

只保留`.png`和`.ico`，将其复制到assets/img/favicons/即可。

自此，网站已经布置成自己的嘞，着手开始安心写博客咯。

## 更改config.yml

```yaml

theme: jekyll-theme-chirpy

# 全局语言，例如将侧边的“home”改成“主页”等，可选项在_data/locales/下
lang: zh-CN

# 时区
timezone: Asia/Shanghai
 
# 站点标题，显示在网站标签页和左上角
title: Hin小站 # the main title

# 站点副标题，在标题之下
tagline: 时间太瘦，指缝太宽。 

# 站点描述，用于搜索引擎优化，改不改问题不大
description: 

# 站点域名，注意结尾不包括 '/'。
url: "https://hintryer.github.io"

# github名，改成自己的github名字
github:
  username: hintryer
  
# twitter名，改成自己的twitter名字
#twitter:
#  username: twitter_username 

### 自己的社交帐号信息，在页脚中显示为文章的默认作者和版权所有者。
social:
  name: Hin所有
  email: example@domain.com 
  links:
    # 第一个元素用作版权所有者的链接，在下方取消注释以添加更多社交链接
    - https://github.com/hintryer/hintryer.github.io
    - https://h2o-ac-doc.lisz.me/basic/config.html
    # - https://www.facebook.com/username
    # - https://www.linkedin.com/in/username

# 站点验证设置，默认为空
webmaster_verifications:
  google: # fill in your Google verification code
  bing: # fill in your Bing verification code
  alexa: # fill in your Alexa verification code
  yandex: # fill in your Yandex verification code
  baidu: # fill in your Baidu verification code
  facebook: # fill in your Facebook verification code

# 网页分析设置，默认为空
analytics:
  google:
    id: # fill in your Google Analytics ID
  goatcounter:
    id: # fill in your GoatCounter ID
  umami:
    id: # fill in your Umami ID
    domain: # fill in your Umami domain
  matomo:
    id: # fill in your Matomo ID
    domain: # fill in your Matomo domain
  cloudflare:
    id: # fill in your Cloudflare Web Analytics token
  fathom:
    id: # fill in your Fathom Site ID

# 页面查看次数设置，现在只支持'goatcounter'，默认为空
pageviews:
  provider: 

# 首选颜色方案设置。留空将默认遵循系统首选颜色，并且在侧边栏的左下角会有一个开关来在深色和浅色之间切换主题

theme_mode: # [light | dark]

# 媒体资源的 CDN 地址。
cdn: ""

# 侧边栏头像，支持本地或 CORS 资源
avatar: "/assets/img/favicons/apple-touch-icon.png"

# 社交预览图像，没什么用
social_preview_image: # string, local or CORS resources

# 博客中目录的全局开关。true | false
toc: true

# 评论相关，默认禁用。
comments:
  # 帖子评论系统的全局开关。将其留空表示禁用。
  provider: # [disqus | utterances | giscus]
  # The provider options are as follows:
  disqus:
    shortname: # fill with the Disqus shortname. › https://help.disqus.com/en/articles/1717111-what-s-a-shortname
  # utterances settings › https://utteranc.es/
  utterances:
    repo: # <gh-username>/<repo>
    issue_term: # < url | pathname | title | ...>
  # Giscus options › https://giscus.app
  giscus:
    repo: # <gh-username>/<repo>
    repo_id:
    category:
    category_id:
    mapping: # optional, default to 'pathname'
    strict: # optional, default to '0'
    input_position: # optional, default to 'bottom'
    lang: # optional, default to the value of `site.lang`
    reactions_enabled: # optional, default to the value of `1`

# 自托管静态资产，没什么用。留空禁止
assets:
  self_host:
    enabled: # boolean, keep empty means false
    env: # [development | production]
#PWA有关
pwa:
  enabled: true # The option for PWA feature (installable)
  cache:
    enabled: true # The option for PWA offline cache
    # Paths defined here will be excluded from the PWA cache.
    # Usually its value is the `baseurl` of another website that
    # shares the same domain name as the current website.
    deny_paths:
      # - "/example"  # URLs match `<SITE_URL>/example/*` will not be cached by the PWA
      
# 分页
paginate: 10

# The base URL of your site
baseurl: ""

# 不建议修改以下选项

kramdown:
  footnote_backlink: "&#8617;&#xfe0e;"  #定义应用于脚注反向链接的文本
  syntax_highlighter: rouge #设置语法高亮器
  syntax_highlighter_opts: # Rouge Options › https://github.com/jneen/rouge#full-options
    css_class: highlight
    #default_lang: console
    span:
      line_numbers: false
    block:
      line_numbers: true
      start_line: 1

collections:
  tabs:
    output: true
    sort_by: order

defaults:
  - scope:
      path: "" # An empty string here means all files in the project
      type: posts
    values:
      layout: post
      comments: true # Enable comments in posts.
      toc: true # Display TOC column in posts.
      # DO NOT modify the following parameter unless you are confident enough
      # to update the code of all other post links in this project.
      permalink: /posts/:title/
  - scope:
      path: _drafts
    values:
      comments: false
  - scope:
      path: ""
      type: tabs # see `site.collections`
    values:
      layout: page
      permalink: /:title/

sass:
  style: compressed

compress_html:
  clippings: all
  comments: all
  endings: all
  profile: false
  blanklines: false
  ignore:
    envs: [development]

exclude:
  - "*.gem"
  - "*.gemspec"
  - docs
  - tools
  - README.md
  - LICENSE
  - "*.config.js"
  - package*.json

jekyll-archives:
  enabled: [categories, tags]
  layouts:
    category: category
    tag: tag
  permalinks:
    tag: /tags/:name/
    category: /categories/:name/
```
