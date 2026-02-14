---
title: Colorbox 和 Details 功能测试
date: 2026-02-13 14:12:00 +0800
categories: [jekyll, Chirpy使用]
pin: false
---

本文档用于测试新增的 colorbox 提示框和 details 折叠块功能。

## Colorbox 提示框测试

> 测试

### Info 提示框

{: .box-info }
> #### Info 提示框标题
>
> 这是一个 info 类型的提示框，通常用于提供一般性信息。

### Tip 提示框

{: .box-tip }
> #### Tip 提示框标题
>
> 这是一个 tip 类型的提示框，通常用于提供技巧或建议。

### Warning 提示框

{: .box-warning }
> #### Warning 提示框标题
>
> 这是一个 warning 类型的提示框，通常用于警告重要信息。

### Danger 提示框

{: .box-danger }
> #### Danger 提示框标题
>
> 这是一个 danger 类型的提示框，通常用于显示错误或危险信息。

{: .box-danger }

> 这是一个 danger 类型的提示框，通常用于显示错误或危险信息。

## Details 折叠块测试

### 基础折叠块

<details class="details-block" markdown="1">
<summary>点击查看详细信息</summary>

这里是一些详细的解释内容：

- 第一点说明
- 第二点说明
- 第三点说明

还可以包含代码块：
```go
func main() {
    fmt.Println("Hello, World!")
}
```

数学公式：
$$
E = mc^2
$$

</details>

### 默认展开的折叠块

<details class="details-block" markdown="1" open>
<summary>默认展开的详细信息</summary>

这个折叠块默认是展开状态的。

包含一些列表内容：
1. 有序列表项一
2. 有序列表项二
3. 有序列表项三

引用内容：
> 这是一段引用的文字内容

表格：
| 姓名 | 年龄 | 职业 |
|------|------|------|
| 张三 | 25   | 工程师 |
| 李四 | 30   | 设计师 |

</details>

### 多层嵌套测试

<details class="details-block" markdown="1">
<summary>外层折叠块</summary>

这是外层的内容。

<details class="details-block" markdown="1">
<summary>内层折叠块</summary>

这是内层的详细内容。

- 内层要点一
- 内层要点二

代码示例：
```javascript
console.log("Nested details test");
```

</details>

</details>

## 使用说明

### Colorbox 语法

```markdown
{: .box-info }
> #### 标题
>
> 内容文本
```

支持的类型：
- `.box-info` - 蓝色信息框
- `.box-tip` - 绿色提示框  
- `.box-warning` - 黄色警告框
- `.box-danger` - 红色危险框

### Details 语法

```html
<details class="details-block" markdown="1">
<summary>标题内容</summary>

详细内容，支持 Markdown 语法

</details>
```

添加 `open` 属性可以让折叠块默认展开：
```html
<details class="details-block" markdown="1" open>
```

## 注意事项

1. Colorbox 需要使用特定的类名才能生效
2. Details 折叠块需要添加 `class="details-block"` 和 `markdown="1"` 属性
3. 样式会自动适配深色/浅色主题
4. 支持在内部使用完整的 Markdown 语法