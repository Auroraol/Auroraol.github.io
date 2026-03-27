# PowerShell script to process all markdown files
$files = Get-ChildItem -Filter "*.md" | Where-Object { $_.Name -notlike "Tip_*" -and $_.Name -ne "process_files.py" -and $_.Name -ne "process_files.ps1" } | Sort-Object Name

foreach ($file in $files) {
    $originalName = $file.Name
    Write-Host "处理文件: $originalName"
    
    $content = Get-Content $file.FullName -Raw -Encoding UTF8
    
    # 提取标题
    $title = $null
    
    # 尝试从 front matter 提取
    if ($content -match '(?m)^title:\s*(.+)$') {
        $title = $matches[1].Trim()
        Write-Host "  从 front matter 提取标题: $title"
    }
    
    # 如果没有，从内容中提取
    if (-not $title) {
        if ($content -match '(?m)^#\s*(Tip\s*#\d+[^\r\n]+)') {
            $title = $matches[1].Trim()
            Write-Host "  从内容提取标题: $title"
        }
    }
    
    if (-not $title) {
        Write-Host "  警告: 无法提取标题，跳过" -ForegroundColor Yellow
        continue
    }
    
    # 生成新文件名
    $newName = $title -replace '#', '_' -replace ' ', '_' -replace '/', '_' -replace '\\', '_' -replace ':', '_' -replace '\*', '_' -replace '\?', '_' -replace '"', '_' -replace '<', '_' -replace '>', '_' -replace '\|', '_'
    $newName = $newName -replace '_+', '_' -replace '^_|_$', ''
    $newName = $newName + '.md'
    
    Write-Host "  新文件名: $newName"
    
    # 添加引用
    $reference = "@${originalName} (1-6)`n`n"
    
    if ($content -notmatch [regex]::Escape("@$originalName")) {
        # 查找 front matter 结束位置
        $frontMatterEnd = $content.IndexOf("---`n", 4)
        if ($frontMatterEnd -ne -1) {
            $content = $content.Insert($frontMatterEnd + 4, $reference)
        } else {
            $content = $reference + $content
        }
    }
    
    # 写入新文件
    $newPath = Join-Path $file.DirectoryName $newName
    Set-Content -Path $newPath -Value $content -Encoding UTF8 -NoNewline
    
    # 删除原文件（如果名称不同）
    if ($newName -ne $originalName) {
        Remove-Item $file.FullName
        Write-Host "  已重命名: $originalName -> $newName" -ForegroundColor Green
    } else {
        Write-Host "  文件已更新: $originalName" -ForegroundColor Green
    }
    
    Write-Host ""
}

Write-Host "处理完成！" -ForegroundColor Cyan

