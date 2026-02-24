@echo off
echo Starting Photography Website...
echo.

echo [1/2] Starting Backend Server...
start cmd /k "cd backend && go run cmd/server/main.go"
timeout /t 3 /nobreak >nul

echo [2/2] Starting Frontend Dev Server...
start cmd /k "cd frontend && npm run dev"

echo.
echo Servers are starting...
echo Backend: http://localhost:8080
echo Frontend: http://localhost:5173
echo.
pause
