document.addEventListener("DOMContentLoaded", function () {
    let coll = document.getElementsByClassName("collapsible-container");
    let maxLines = 19; // 设置折叠显示的行数 Maximum number of lines to display without collapsing
    let defaultOpenLines = 30; // 默认展开的最大行数

    for (let i = 0; i < coll.length; i++) {
        // 创建展开/收起按钮
        let trigger = document.createElement('p');
        trigger.className = 'language-javascript collapsible-trigger collapsible-trigger-css';
        trigger.textContent = '展开';
        trigger.style.cssText = `
            display: block;
            width: auto;
            max-width: 80px;
            padding: 0.3rem 0.5rem;
            background-color: var(--highlight-bg-color);
            border: 1px solid var(--language-border-color);
            border-radius: 0 0 0.625rem 0.625rem;
            cursor: pointer;
            font-weight: 600;
            color: var(--code-header-text-color);
            text-align: center;
            transition: all 0.3s ease;
            margin: 0 auto;
            font-size: 0.85rem;
            line-height: 1.4;
        `;
        
        trigger.addEventListener('mouseenter', function() {
            this.style.backgroundColor = 'rgba(128, 128, 128, 0.37)';
            this.style.color = 'white';
        });
        
        trigger.addEventListener('mouseleave', function() {
            this.style.backgroundColor = '';
            this.style.color = '';
        });

        let content = coll[i].querySelector('.highlight');
        if (!content) continue;
        
        // 计算实际的代码行数 (考虑chirpy主题使用table结构)
        let codeLines = (content.textContent.split('\n').length - 1) / 2;
        
        // 插入按钮到容器底部（在代码内容之后）
        coll[i].appendChild(trigger);
        
        // 如果代码行数较少，隐藏按钮并展开内容
        if (codeLines <= defaultOpenLines) {
            trigger.style.display = 'none';
            content.style.maxHeight = content.scrollHeight + "px";
        } else {
            // 设置初始折叠状态
            content.style.maxHeight = (24 * maxLines) + "px"; // 每行约24px
            content.style.overflow = 'hidden';
            
            // 添加点击事件
            trigger.addEventListener("click", function () {
                // 切换按钮文字
                if (this.textContent.includes("展开")) {
                    this.textContent = "收起";
                    content.style.maxHeight = content.scrollHeight + "px";
                } else {
                    this.textContent = "展开";
                    content.style.maxHeight = (24 * maxLines) + "px";
                    // 滚动到按钮位置
                    this.scrollIntoView({
                        behavior: "smooth",
                        block: "center"
                    });
                }
            });
        }
    }
});