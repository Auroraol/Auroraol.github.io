---
title: jekyll主题使用
date: 2023-12-20 11:00:00 +0800
categories: [jekyll, jekyll使用]
tags: [jekyll使用]
---

结合了一下几篇文章，做个总结：
1.修改gem的source ：

gem sources --remove https://rubygems.org/
AI写代码
c
运行
1
这是本身的gem的网站，但是它是国外的，一般都被墙，所以要修改，先将其移除。
2.改为可以使用的：

gem sources -a 'https://gems.ruby-china.com'
AI写代码
c
运行
1
我看过有地址写为http://ruby.taobao.org/的，但是我设置后没有用，应该是用不了了，所以，用上面那个。
3.查看当前有的source

 gem sources -l
AI写代码
c
运行
1
会显示：

*** CURRENT SOURCES ***

https://gems.ruby-china.com
AI写代码
1
2
3
4.最后：
————————————————
版权声明：本文为CSDN博主「哥兜兜里有泡泡糖」的原创文章，遵循CC 4.0 BY-SA版权协议，转载请附上原文出处链接及本声明。
原文链接：https://blog.csdn.net/weixin_44512194/article/details/107053421

### 本地部署环境和调试

首先调试环境非必须，因为后续大部分使用md文件来创建文章，一般的编辑器都带有实时预览功能，所见即所得，一般写完之后直接上传md和相关资源文件即可。如果上面使用Using the Chirpy Starter方式创建，也完全可以不用该调试环境，专注于内容创作即可。如果是fork方式，建议一定要配置好自己的本地环境，非常方便开发调试。

不同的操作系统可能会不太一样，这里以Windows举例，配置本地调试环境：

1. 安装RUBY

   https://rubyinstaller.org/downloads/

   记得下载带devkit的版本

   安装过程中，如果有类似MSYS2 and MINGW development tool chain的选项，记得勾选

   安装成功后，使用命令行验证安装结果：

   ` ruby -v gem -v `

   展开

   未出现报错即安装成功

2. 安装Jekyll

   ` gem install jekyll bundler `

   展开

   安装成功后，使用命令行验证安装结果：

   ` jekyll -v `

   展开

3. 运行项目

   首先使用git clone拉到项目所有文件，然后进入项目根目录文件夹中，运行命令：

   ` bundle `

   展开

   如果安装过程过于缓慢，建议配置国内镜像：

   ` bundle config mirror.https://rubygems.org https://gems.ruby-china.com `

   展开

   然后再运行bundle

   运行成功后，记得配置项目目录下的_config.yml中的变量，url，avatar，timezone，lang

   还有其他变量可以多探索下。

   最后使用命令

   ` bundle exec jekyll serve `

   展开

   成功的话命令行中会有提醒对应的网址，一般是 http://127.0.0.1:4000/

   访问该网址即是对应自己网站的首页。
