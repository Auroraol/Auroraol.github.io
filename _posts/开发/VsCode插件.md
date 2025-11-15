---
title: VsCode插件开发
date: 2025-11-15 00:40:00 +0800
categories: [开发, VsCode插件开发]
tags: [插件开发]
---



[◼️VS Code插件创作中文开发文档](https://liiked.github.io/VS-Code-Extension-Doc-ZH/#/working-with-extensions/testing-extension)

选择想要创建的插件类型。根据选择的类型，下载对应的模板。[更多模板](https://github.com/microsoft/vscode-extension-samples)

项目结构一般包含:

```
.vscode`: 里面的文件是用来测试插件或者测试代码的一些文件。
`node_modules`: 第三方依赖。
`src/test`：测试文件。
`src/extension.ts`：插件的主文件。
```

