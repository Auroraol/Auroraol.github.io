---
title: 使用giscus配置GitHub Pages评论功能
date: 2025-11-02 11:00:00 +0800
categories: [jekyll, jekyll使用]
tags: [jekyll使用]
---

# 评论能力开通

## GitHub仓库配置

首先需要在自己的仓库中开通Discussion功能，一般新仓库默认不会打开，需要手动开通然后再添加对应的模块即可

1. 开通Discussion功能

    进入自己的GitHub pages所在仓库，找到setting界面和里面的General选项

    ![主题](https://github.com/Auroraol/Drawing-bed/raw/main/img/20240104154923.png)下滑找到Feature中勾选Discussion即可打开

    ![主题](https://github.com/Auroraol/Drawing-bed/raw/main/img/20240104155100.png)创建comment

    找到Discussion界面，点Categories右侧的🖊按钮，进入后点New category，准备新建一个category

    ![主题](https://github.com/Auroraol/Drawing-bed/raw/main/img/20240104155447.png)

    这里name就填Comments，Format勾选Open-ended discussion，完事确认就配置ok

    ![主题](https://github.com/Auroraol/Drawing-bed/raw/main/img/20240104155820.png)


## giscus配置

首先打开gitscu官网，需要使用当前仓库的GitHub账号登录：https://giscus.app/

1. 配置仓库：

    这里配置的仓库就是自己现在想要配置的pages页面的仓库地址，用户名+仓库名

2. 配置Discussion 分类

    这里就是步骤1中创建的comment选项，勾选即可

    ![主题](https://github.com/Auroraol/Drawing-bed/raw/main/img/20240104160029.png)

3. 其余保持默认，但是不要关闭该也页面，下一步配置中需要使用到这里的信息

    ![image-20251102192236662](https://github.com/Auroraol/Drawing-bed/raw/main/img/image-20251102192236662.png)

1. 仓库config配置

    ![image-20251102192301298](https://github.com/Auroraol/Drawing-bed/raw/main/img/image-20251102192301298.png)

    

    