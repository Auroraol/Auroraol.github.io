#!/bin/bash

# 处理所有 .md 文件（排除已处理的 Tip_ 开头的文件）
for file in *.md; do
    # 跳过已处理的文件和脚本文件
    if [[ "$file" == Tip_* ]] || [[ "$file" == process_* ]]; then
        continue
    fi
    
    echo "处理文件: $file"
    
    # 提取标题
    title=$(grep -m1 "^title:" "$file" 2>/dev/null | sed 's/^title:[[:space:]]*//' | sed 's/[[:space:]]*$//')
    
    # 如果没有 front matter，从内容中提取
    if [ -z "$title" ]; then
        title=$(grep -m1 "^# Tip #" "$file" 2>/dev/null | sed 's/^#[[:space:]]*//' | sed 's/[[:space:]]*$//')
    fi
    
    if [ -z "$title" ]; then
        echo "  警告: 无法提取标题，跳过"
        continue
    fi
    
    echo "  提取的标题: $title"
    
    # 生成新文件名（替换特殊字符）
    newname=$(echo "$title" | sed 's/#/_/g' | sed 's/ /_/g' | sed 's/\//_/g' | sed 's/\\/_/g' | sed 's/:/_/g' | sed 's/*/_/g' | sed 's/?/_/g' | sed 's/"/_/g' | sed 's/</_/g' | sed 's/>/_/g' | sed 's/|/_/g' | sed 's/__*/_/g' | sed 's/^_\|_$//g')
    newname="${newname}.md"
    
    echo "  新文件名: $newname"
    
    # 检查是否已有引用
    if ! grep -q "@${file}" "$file" 2>/dev/null; then
        # 查找 front matter 结束位置
        if grep -q "^---$" "$file"; then
            # 找到第二个 --- 的位置
            line_num=$(awk '/^---$/ {count++; if(count==2) {print NR; exit}}' "$file")
            if [ -n "$line_num" ]; then
                # 在第二个 --- 后插入引用
                reference="@${file} (1-6)\n\n"
                awk -v ref="$reference" -v line="$line_num" 'NR==line {print; print ref; next} 1' "$file" > "${file}.tmp"
                mv "${file}.tmp" "$file"
            else
                # 没有找到第二个 ---，在文件开头添加
                reference="@${file} (1-6)\n\n"
                echo -e "$reference$(cat "$file")" > "${file}.tmp"
                mv "${file}.tmp" "$file"
            fi
        else
            # 没有 front matter，在文件开头添加
            reference="@${file} (1-6)\n\n"
            echo -e "$reference$(cat "$file")" > "${file}.tmp"
            mv "${file}.tmp" "$file"
        fi
    fi
    
    # 重命名文件
    if [ "$file" != "$newname" ]; then
        mv "$file" "$newname"
        echo "  已重命名: $file -> $newname"
    else
        echo "  文件已更新: $file"
    fi
    
    echo ""
done

echo "处理完成！"

