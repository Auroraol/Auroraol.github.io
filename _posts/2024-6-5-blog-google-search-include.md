---
title: Blog可以被Google搜索到
date: 2025-11-02 17:00:00 +0800
categories: [jekyll, jekyll使用]
tags: [jekyll使用]
---

# Google

1. 确认自己的blog确实未被Google搜索收录

    浏览器输入如下网址：

    ```
    site:https://your_github_id.github.io/
    ```

    如果没有被Google搜索收录，那会出现如下图所示，跟着笔者往下一步一步配置吧

    ![image-20251102193837802](https://github.com/Auroraol/Drawing-bed/raw/main/img/2deb8366d0f602668070c426b2e2b251.png)

    如果被收录了，那就会出现如下图所示，忽略本文章的后续步骤吧

2. 进入上图所提示的Google Search Console

    选择右侧的网址前缀，并输入自己的blog地址点继续

    ```
    https://your_github_id.github.io/
    ```

    ![img](https://github.com/Auroraol/Drawing-bed/raw/main/img/image-20251102194904658.png)

3. 下载并将对应的html文件上传到自己blog中

    ![截图](https://github.com/Auroraol/Drawing-bed/raw/main/img/20240605005315.png)

4. 添加站点地图

    上一步成功之后，点弹窗右下角前往资源页面，在跳转的网页中右侧会有网络地图的tag

    ![image-20251102194955914](https://github.com/Auroraol/Drawing-bed/raw/main/img/image-20251102193837802.png)网络地图的作用是给Google知道你网站的结构，方便它的爬虫来爬你的内容的，所以非常有必要设置

     [https://www.xml-sitemaps.com/](https://www.xml-sitemaps.com/)这个网站，输入自己网站地址

    ```
    https://your_github_id.github.io/
    ```
    
    等待一会，弹窗提示成功，点VIEW SITEMAP DETAILS，在跳转的网页点DOWNLOAD YOUR XML SITEMAP FILE，下载到sitemap.xml文件

    ![image-20251102194904658](https://github.com/Auroraol/Drawing-bed/raw/main/img/image-20251102194955914.png)

    将该sitemap.xml文件上传到自己网站的根目录

    等待1天后再回来看看数据是否发生变化，也可以在Google搜索引擎中搜索自己的blog查看

# Bing

由于GitHub同样也封禁了Bing搜索，在设置完Google搜索后，同时也要设置下Bing搜索。

[Bing Webmaster Tools](https://www.bing.com/webmasters) 

![image-20251102202817224](https://github.com/Auroraol/Drawing-bed/raw/main/img/image-20251102202817224.png)
