---
title: 测试代码折叠功能
date: 2026-02-12 23:00:00 +0800
categories: [测试]
tags: [代码折叠, 功能测试]
---

这是一个测试代码折叠功能的文章。

## 短代码块测试

```python
print("Hello World")
```

## 长代码块测试

```python
def fibonacci(n):
    """计算斐波那契数列的第n项"""
    if n <= 0:
        return 0
    elif n == 1:
        return 1
    else:
        a, b = 0, 1
        for i in range(2, n + 1):
            a, b = b, a + b
        return b

def main():
    # 计算前20个斐波那契数
    for i in range(20):
        print(f"F({i}) = {fibonacci(i)}")
    
    # 一些额外的代码来增加长度
    x = 10
    y = 20
    z = x + y
    
    if z > 25:
        print("z is greater than 25")
    else:
        print("z is less than or equal to 25")
    
    # 循环测试
    for i in range(100):
        if i % 10 == 0:
            print(f"Processing item {i}")
    
    # 列表操作
    my_list = [1, 2, 3, 4, 5]
    squared = [x**2 for x in my_list]
    print(f"Squared values: {squared}")
    
    # 字典操作
    my_dict = {"a": 1, "b": 2, "c": 3}
    for key, value in my_dict.items():
        print(f"{key}: {value}")
    
    # 函数调用
    result = fibonacci(15)
    print(f"Fibonacci(15) = {result}")

if __name__ == "__main__":
    main()
```

## 更长的代码块测试

```javascript
// 这是一个更长的JavaScript代码示例
class Calculator {
    constructor() {
        this.history = [];
    }
    
    add(a, b) {
        const result = a + b;
        this.history.push(`${a} + ${b} = ${result}`);
        return result;
    }
    
    subtract(a, b) {
        const result = a - b;
        this.history.push(`${a} - ${b} = ${result}`);
        return result;
    }
    
    multiply(a, b) {
        const result = a * b;
        this.history.push(`${a} * ${b} = ${result}`);
        return result;
    }
    
    divide(a, b) {
        if (b === 0) {
            throw new Error("Division by zero");
        }
        const result = a / b;
        this.history.push(`${a} / ${b} = ${result}`);
        return result;
    }
    
    getHistory() {
        return this.history;
    }
    
    clearHistory() {
        this.history = [];
    }
}

// 使用示例
const calc = new Calculator();

console.log(calc.add(10, 5));      // 15
console.log(calc.subtract(20, 8)); // 12
console.log(calc.multiply(6, 7));  // 42
console.log(calc.divide(100, 4));  // 25

console.log("Calculation history:");
calc.getHistory().forEach(entry => {
    console.log(entry);
});

// 更多复杂操作
function complexCalculation(x, y, z) {
    const step1 = calc.add(x, y);
    const step2 = calc.multiply(step1, z);
    const step3 = calc.subtract(step2, 10);
    return calc.divide(step3, 2);
}

const result = complexCalculation(5, 3, 4);
console.log(`Complex calculation result: ${result}`);

// 异步操作示例
async function asyncCalculation() {
    return new Promise((resolve) => {
        setTimeout(() => {
            const result = calc.add(100, 200);
            resolve(result);
        }, 1000);
    });
}

asyncCalculation().then(result => {
    console.log(`Async result: ${result}`);
});

// 错误处理
try {
    calc.divide(10, 0);
} catch (error) {
    console.error("Error occurred:", error.message);
}

// 最终历史记录
console.log("Final calculation history:");
console.log(calc.getHistory());
```