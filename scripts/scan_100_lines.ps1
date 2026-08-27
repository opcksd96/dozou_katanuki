$exts = @('*.go', '*.py', '*.ts', '*.js', '*.vue')
$files = Get-ChildItem -Recurse -File -Include $exts | Where-Object {
    $_.FullName -notmatch '(node_modules|\.git|dist|wailsjs|__pycache__|\.venv|release|build|bin[\\/]|docs[\\/]|tests[\\/]|scripts[\\/]docs|\.trash)'
}


$over100 = @()
$compliantCount = 0

foreach ($f in $files) {
    $lines = (Get-Content $f.FullName | Measure-Object -Line).Lines
    $relPath = $f.FullName.Replace((Get-Location).Path + '\', '')
    if ($lines -gt 100) {
        $over100 += [PSCustomObject]@{
            Lines = $lines
            Path  = $relPath
        }
    } else {
        $compliantCount++
    }
}

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host " 100-Line Rule Compliance Audit Report" -ForegroundColor Cyan
Write-Host " Total Source Files Audited: $($files.Count)" -ForegroundColor White
Write-Host " Compliant Files (<= 100):  $compliantCount" -ForegroundColor Green
$statusColor = if ($over100.Count -eq 0) { 'Green' } else { 'Red' }
Write-Host " Non-Compliant Files (> 100): $($over100.Count)" -ForegroundColor $statusColor
Write-Host "==========================================" -ForegroundColor Cyan

if ($over100.Count -eq 0) {
    Write-Host "🎉 PERFECT: All $($files.Count) production source files strictly comply with SPEC-PRINCIPLE-001 (100-line rule)!" -ForegroundColor Green
} else {
    Write-Host "WARNING: $($over100.Count) files exceed 100 lines:" -ForegroundColor Yellow
    $over100 | Sort-Object Lines -Descending | Format-Table -AutoSize
}
