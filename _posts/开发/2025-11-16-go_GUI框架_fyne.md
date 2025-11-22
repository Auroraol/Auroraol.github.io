---
title: go-GUI框架-fyne
date: 2025-11-16 00:40:00 +0800
categories: [开发]
tags: [go-GUI框架-fyne]
---

### 1.安装C编译器

https://sourceforge.net/projects/mingw-w64/files/mingw-w64/

安装后会自动加入环境变量path，自行确认

cmd下输入指令gcc -v，确定安装成功和版本

![img](https://github.com/Auroraol/Drawing-bed/raw/main/img/3202470-20230818155754679-1679994735.png)

###  2.安装fyne

使用标准的go工具，安装Fyne的核心库使用:

```
go get fyne.io/fyne/v2@latest  #安装fyne框架库
go get fyne.io/fyne/v2/cmd/fyne  #安装fyne工具，go版本<1.16
go install fyne.io/fyne/v2/cmd/fyne@latest  #安装fyne工具,go高版本
```

### 3.检查安装版本

在编写应用程序代码或运行示例之前，您可以使用 [Fyne 安装](https://geoffrey-artefacts.fynelabs.com/github/andydotxyz/fyne-io/setup/latest/)工具检查您的安装。只需从链接下载适合您计算机的应用程序并运行它，您应该会看到类似以下屏幕的内容

![img](https://github.com/Auroraol/Drawing-bed/raw/main/img/3202470-20230818160232624-1457660590.png)

### 4.运行演示

请注意，第一次运行必须编译一些 C 代码，因此可能需要比平时更长的时间。后续构建会重用缓存，速度会快很多。

```
go run fyne.io/fyne/v2/cmd/fyne_demo@latest
```

运行会自动下载相应的依赖：

![img](https://github.com/Auroraol/Drawing-bed/raw/main/img/3202470-20230818160443847-505802920.jpg)

运行结果， 看起来还不错：

![img](https://github.com/Auroraol/Drawing-bed/raw/main/img/3202470-20230818160519303-366697541.jpg)
