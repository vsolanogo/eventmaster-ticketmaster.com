@echo off
echo ========================================
echo     SCRIPT FOR CREATING THE DATABASE
echo ========================================
echo.

:: --- НАСТРОЙКИ ПОДКЛЮЧЕНИЯ ---
:: Замените значения на свои (должны совпадать с drop_db.bat)
set DB_NAME=events
set DB_USER=postgres
set DB_PASSWORD=passwordSuperUser1111
set DB_HOST=localhost
set DB_PORT=5433

:: --- ПОЛНЫЙ ПУТЬ К PSQL ---
:: Замените "16" на вашу версию PostgreSQL!
set PSQL_PATH="C:\Program Files\PostgreSQL\18\bin\psql.exe"
:: -----------------------------

echo Connecting to PostgreSQL...
echo Attempting to create database "%DB_NAME%"...

:: Устанавливаем пароль в переменную окружения для текущей сессии
set PGPASSWORD=%DB_PASSWORD%

:: Выполняем SQL-команду для создания базы данных.
:: Также подключаемся к служебной базе "postgres".
%PSQL_PATH% -h %DB_HOST% -p %DB_PORT% -U %DB_USER% -d postgres -c "CREATE DATABASE %DB_NAME%;"

:: Проверяем результат
if %ERRORLEVEL% EQU 0 (
    echo.
    echo SUCCESS: Database "%DB_NAME%" has been created successfully.
) else (
    echo.
    echo ERROR: Failed to create database. Check connection settings, permissions, or if the database already exists.
)

echo.
pause
