# PowerShell链接检查脚本
Write-Host "检查内部链接..." -ForegroundColor Green

# 检查是否存在指向/tag/的链接
$tagLinks = Get-ChildItem "_site" -Recurse -Filter "*.html" | Select-String "/tag/"

if ($tagLinks) {
    Write-Host "发现指向/tag/的链接:" -ForegroundColor Red
    $tagLinks | ForEach-Object { Write-Host $_.Path ":" $_.Line }
    exit 1
} else {
    Write-Host "未发现指向/tag/的链接 - 修复成功!" -ForegroundColor Green
}

# 检查图片链接
Write-Host "检查图片链接..." -ForegroundColor Green
$imageLinks = Get-ChildItem "_site" -Recurse -Filter "*.html" | Select-String "image-20260214174940676.png"

if ($imageLinks) {
    Write-Host "发现损坏的图片链接:" -ForegroundColor Red
    $imageLinks | ForEach-Object { Write-Host $_.Path ":" $_.Line }
    exit 1
} else {
    Write-Host "未发现损坏的图片链接" -ForegroundColor Green
}

Write-Host "所有检查通过!" -ForegroundColor Green