@echo off
echo Starting integration test...
echo.

REM 检查服务器是否运行
curl -s http://localhost:8080/api/config >nul 2>&1
if %errorlevel% neq 0 (
    echo [ERROR] Server not running at http://localhost:8080
    echo Please start the server first:
    echo   go run cmd/server/main.go -config config.yaml
    echo.
    pause
    exit /b 1
)

echo Server is running. Starting tests...
echo.

"D:\Program Files\Go\bin\go.exe" run test/integration/main.go

pause
