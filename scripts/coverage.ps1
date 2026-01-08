# Script to generate test coverage reports locally (PowerShell)
# Usage: .\scripts\coverage.ps1

Write-Host "🧪 Running tests with coverage..." -ForegroundColor Cyan
go test ./... -coverprofile=coverage.out -covermode=atomic

Write-Host ""
Write-Host "📊 Coverage summary:" -ForegroundColor Cyan
go tool cover -func=coverage.out

Write-Host ""
Write-Host "📈 Generating HTML coverage report..." -ForegroundColor Cyan
go tool cover -html=coverage.out -o coverage.html

Write-Host ""
Write-Host "✅ Coverage reports generated:" -ForegroundColor Green
Write-Host "   - coverage.out (for SonarCloud)" -ForegroundColor White
Write-Host "   - coverage.html (for local viewing)" -ForegroundColor White
Write-Host ""
Write-Host "💡 Open coverage.html in your browser to see detailed coverage" -ForegroundColor Yellow
