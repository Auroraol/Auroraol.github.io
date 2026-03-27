#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import os
import re
import shutil

def extract_title_from_frontmatter(content):
    """从 front matter 中提取 title"""
    match = re.search(r'^title:\s*(.+)$', content, re.MULTILINE)
    if match:
        return match.group(1).strip()
    return None

def extract_title_from_content(content):
    """从内容中提取第一个 Tip 标题"""
    match = re.search(r'^#\s*(Tip\s*#\d+[^\n]+)$', content, re.MULTILINE)
    if match:
        return match.group(1).strip()
    return None

def sanitize_filename(title):
    """将标题转换为安全的文件名"""
    # 替换特殊字符
    filename = title.replace('#', '_')
    filename = filename.replace(' ', '_')
    filename = filename.replace('/', '_')
    filename = filename.replace('\\', '_')
    filename = filename.replace(':', '_')
    filename = filename.replace('*', '_')
    filename = filename.replace('?', '_')
    filename = filename.replace('"', '_')
    filename = filename.replace('<', '_')
    filename = filename.replace('>', '_')
    filename = filename.replace('|', '_')
    # 移除多个连续的下划线
    filename = re.sub(r'_+', '_', filename)
    # 移除开头和结尾的下划线
    filename = filename.strip('_')
    return filename + '.md'

def add_reference(content, original_filename):
    """在 front matter 后添加引用"""
    # 检查是否已经有引用
    if '@' + original_filename in content:
        return content
    
    reference = f"@{original_filename} (1-6)\n\n"
    
    # 查找 front matter 结束位置
    frontmatter_end = content.find('---\n', 4)  # 从第4个字符开始查找，跳过开头的 ---
    if frontmatter_end != -1:
        # 找到第二个 ---，在它后面添加引用
        insert_pos = frontmatter_end + 4
        content = content[:insert_pos] + reference + content[insert_pos:]
    else:
        # 没有 front matter，在文件开头添加
        content = reference + content
    return content

def process_file(filename):
    """处理单个文件"""
    if not filename.endswith('.md') or filename == 'process_files.py':
        return
    
    # 跳过已经处理过的文件（以 Tip_ 开头的）
    if filename.startswith('Tip_'):
        return
    
    original_filename = filename
    print(f"处理文件: {original_filename}")
    
    # 读取文件内容
    with open(original_filename, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # 提取标题
    title = extract_title_from_frontmatter(content)
    if not title:
        title = extract_title_from_content(content)
    
    if not title:
        print(f"  警告: 无法从 {original_filename} 提取标题，跳过")
        return
    
    print(f"  提取的标题: {title}")
    
    # 生成新文件名
    new_filename = sanitize_filename(title)
    print(f"  新文件名: {new_filename}")
    
    # 添加引用
    content = add_reference(content, original_filename)
    
    # 写入新文件
    with open(new_filename, 'w', encoding='utf-8') as f:
        f.write(content)
    
    # 删除原文件
    if new_filename != original_filename:
        os.remove(original_filename)
        print(f"  已重命名: {original_filename} -> {new_filename}")
    else:
        print(f"  文件已更新: {original_filename}")

def main():
    """主函数"""
    # 获取当前目录所有 .md 文件
    files = [f for f in os.listdir('.') if f.endswith('.md') and not f.startswith('Tip_')]
    files.sort()
    
    print(f"找到 {len(files)} 个文件需要处理\n")
    
    for filename in files:
        try:
            process_file(filename)
        except Exception as e:
            print(f"  错误处理 {filename}: {e}")
        print()

if __name__ == '__main__':
    main()

