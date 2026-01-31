#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
爬取 Go Tips 网站内容
从 https://colobu.com/gotips/001.html 到 https://colobu.com/gotips/082.html
"""

import requests
from bs4 import BeautifulSoup
import os
import time
from pathlib import Path
import re
from urllib.parse import urljoin, urlparse

def scrape_page(page_num):
    """爬取单个页面"""
    url = f"https://colobu.com/gotips/{page_num:03d}.html"
    print(f"正在爬取: {url}")
    
    try:
        headers = {
            'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36'
        }
        response = requests.get(url, headers=headers, timeout=10)
        response.raise_for_status()
        response.encoding = 'utf-8'
        
        return response.text
    except requests.RequestException as e:
        print(f"错误: 无法获取 {url}: {e}")
        return None

def download_image(img_url, page_num, output_dir='gotips_content'):
    """下载图片到本地"""
    try:
        headers = {
            'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36'
        }
        
        # 处理相对URL和绝对URL
        original_url = img_url
        if not img_url.startswith('http'):
            # 如果是相对路径，补全为完整URL
            if img_url.startswith('/'):
                img_url = urljoin('https://colobu.com', img_url)
            else:
                # 相对路径，使用当前页面的base URL
                img_url = urljoin(f'https://colobu.com/gotips/{page_num:03d}.html', img_url)
        
        print(f"  下载图片: {img_url}")
        response = requests.get(img_url, headers=headers, timeout=10)
        response.raise_for_status()
        
        # 创建图片目录
        img_dir = os.path.join(output_dir, 'images', f'{page_num:03d}')
        os.makedirs(img_dir, exist_ok=True)
        
        # 获取图片文件名
        parsed_url = urlparse(img_url)
        img_filename = os.path.basename(parsed_url.path)
        
        # 如果URL路径中包含页面编号目录（如 /gotips/images/004/1.jpeg），直接使用文件名
        path_parts = parsed_url.path.split('/')
        if len(path_parts) > 1 and path_parts[-2].isdigit():
            # 路径格式为 /gotips/images/004/1.jpeg，使用最后的文件名
            img_filename = path_parts[-1]
        
        if not img_filename or '.' not in img_filename:
            # 如果没有扩展名，尝试从Content-Type获取
            content_type = response.headers.get('Content-Type', '')
            if 'image/png' in content_type:
                img_filename = 'image.png'
            elif 'image/jpeg' in content_type or 'image/jpg' in content_type:
                img_filename = 'image.jpg'
            elif 'image/gif' in content_type:
                img_filename = 'image.gif'
            elif 'image/webp' in content_type:
                img_filename = 'image.webp'
            else:
                img_filename = 'image.png'
        
        # 保存图片
        img_path = os.path.join(img_dir, img_filename)
        with open(img_path, 'wb') as f:
            f.write(response.content)
        
        print(f"  图片已保存: {img_path}")
        
        # 返回相对路径（用于Markdown）
        return f'images/{page_num:03d}/{img_filename}'
    except Exception as e:
        print(f"警告: 无法下载图片 {original_url}: {e}")
        return original_url  # 返回原始URL

def html_to_markdown(element, page_num, base_url, output_dir='gotips_content'):
    """将HTML元素转换为Markdown格式"""
    if element is None:
        return ""
    
    markdown = []
    
    for child in element.children:
        if isinstance(child, str):
            # 纯文本节点
            text = child.strip()
            if text:
                markdown.append(text)
        elif hasattr(child, 'name'):
            tag_name = child.name
            
            if tag_name == 'h1':
                markdown.append(f"# {child.get_text().strip()}\n")
            elif tag_name == 'h2':
                markdown.append(f"## {child.get_text().strip()}\n")
            elif tag_name == 'h3':
                markdown.append(f"### {child.get_text().strip()}\n")
            elif tag_name == 'h4':
                markdown.append(f"#### {child.get_text().strip()}\n")
            elif tag_name == 'p':
                p_text = html_to_markdown(child, page_num, base_url, output_dir).strip()
                if p_text:
                    markdown.append(f"{p_text}\n")
            elif tag_name == 'pre':
                # 代码块
                code_tag = child.find('code')
                if code_tag:
                    code_text = code_tag.get_text()
                    lang = code_tag.get('class', [])
                    lang_str = ''
                    if lang:
                        # 提取语言类型，如 'language-go' -> 'go'
                        for cls in lang:
                            if 'language-' in cls:
                                lang_str = cls.replace('language-', '')
                                break
                    markdown.append(f"```{lang_str}\n{code_text}\n```\n")
                else:
                    markdown.append(f"```\n{child.get_text()}\n```\n")
            elif tag_name == 'code':
                # 检查是否在pre标签内（代码块），如果是则跳过，因为pre标签会处理
                parent = child.find_parent('pre')
                if not parent:
                    # 行内代码
                    code_text = child.get_text()
                    markdown.append(f"`{code_text}`")
            elif tag_name == 'img':
                # 图片
                img_src = child.get('src') or child.get('data-src') or child.get('data-lazy-src')
                if img_src:
                    # 下载图片
                    local_path = download_image(img_src, page_num, output_dir)
                    alt_text = child.get('alt', '') or child.get('title', '')
                    markdown.append(f"![{alt_text}]({local_path})\n")
            elif tag_name == 'a':
                # 链接
                href = child.get('href', '')
                link_text = child.get_text().strip()
                markdown.append(f"[{link_text}]({href})")
            elif tag_name == 'strong' or tag_name == 'b':
                # 粗体
                text = child.get_text().strip()
                markdown.append(f"**{text}**")
            elif tag_name == 'em' or tag_name == 'i':
                # 斜体
                text = child.get_text().strip()
                markdown.append(f"*{text}*")
            elif tag_name == 'ul' or tag_name == 'ol':
                # 列表
                list_md = html_to_markdown(child, page_num, base_url, output_dir)
                markdown.append(list_md)
            elif tag_name == 'li':
                # 列表项
                li_text = html_to_markdown(child, page_num, base_url, output_dir).strip()
                markdown.append(f"- {li_text}\n")
            elif tag_name == 'br':
                markdown.append("\n")
            elif tag_name in ['div', 'span', 'section']:
                # 容器标签，递归处理
                content = html_to_markdown(child, page_num, base_url, output_dir)
                if content.strip():
                    markdown.append(content)
            else:
                # 其他标签，提取文本
                text = child.get_text().strip()
                if text:
                    markdown.append(text)
    
    return ''.join(markdown)

def extract_content(html_content, page_num, base_url, output_dir='gotips_content'):
    """提取页面主要内容"""
    soup = BeautifulSoup(html_content, 'html.parser')
    
    # 提取文章标题
    title = ""
    title_tag = soup.find('h1')
    if title_tag:
        title = title_tag.get_text().strip()
    
    # 提取文章内容 - 通常内容在 article 或 main 标签中
    content = ""
    article = None
    
    # 尝试多种方式找到主要内容
    article = soup.find('article')
    if not article:
        article = soup.find('main')
    if not article:
        article = soup.find('div', class_='content')
    if not article:
        article = soup.find('div', id='content')
    
    if not article:
        # 如果找不到特定容器，尝试找到包含h1的div
        if title_tag:
            article = title_tag.find_parent('div')
    
    if article:
        # 移除脚本和样式标签
        for script in article(["script", "style", "nav", "header", "footer"]):
            script.decompose()
        
        # 转换为Markdown
        content = html_to_markdown(article, page_num, base_url, output_dir)
    else:
        # 如果找不到特定容器，提取body内容
        body = soup.find('body')
        if body:
            for script in body(["script", "style", "nav", "header", "footer"]):
                script.decompose()
            content = html_to_markdown(body, page_num, base_url, output_dir)
    
    return title, content

def save_content(page_num, title, content, output_dir='gotips_content'):
    """保存内容到文件"""
    os.makedirs(output_dir, exist_ok=True)
    
    filename = f"{page_num:03d}_{title.replace('/', '_').replace('\\', '_')[:50]}.md"
    filepath = os.path.join(output_dir, filename)
    
    with open(filepath, 'w', encoding='utf-8') as f:
        f.write(f"# {title}\n\n")
        f.write(f"来源: https://colobu.com/gotips/{page_num:03d}.html\n\n")
        f.write("---\n\n")
        f.write(content)
    
    print(f"已保存: {filepath}")
    return filepath

def main():
    """主函数"""
    start_page = 1
    end_page = 82
    
    print(f"开始爬取 Go Tips 网站 (页面 {start_page} 到 {end_page})")
    print("=" * 60)
    
    success_count = 0
    fail_count = 0
    
    for page_num in range(start_page, end_page + 1):
        base_url = f"https://colobu.com/gotips/{page_num:03d}.html"
        html_content = scrape_page(page_num)
        
        if html_content:
            title, content = extract_content(html_content, page_num, base_url)
            
            if title and content:
                save_content(page_num, title, content)
                success_count += 1
            else:
                print(f"警告: 页面 {page_num:03d} 内容提取失败")
                fail_count += 1
        else:
            fail_count += 1
        
        # 避免请求过快，添加延迟
        time.sleep(0.5)
    
    print("=" * 60)
    print(f"爬取完成!")
    print(f"成功: {success_count} 页")
    print(f"失败: {fail_count} 页")

if __name__ == "__main__":
    main()

