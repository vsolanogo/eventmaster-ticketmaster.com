@echo off
echo ========================================
echo   SCRIPT FOR LISTING DATABASE TABLES
echo ========================================
echo.

:: --- НАСТРОЙКИ ПОДКЛЮЧЕНИЯ ---
:: Убедитесь, что они совпадают с вашими другими скриптами
set DB_NAME=events
set DB_USER=postgres
set DB_PASSWORD=passwordSuperUser1111
set DB_HOST=localhost
set DB_PORT=5433

:: --- ПОЛНЫЙ ПУТЬ К PSQL ---
:: Замените "18" на вашу версию PostgreSQL!
set PSQL_PATH="C:\Program Files\PostgreSQL\18\bin\psql.exe"
:: -----------------------------

echo Connecting to database "%DB_NAME%"...
echo.

:: Устанавливаем пароль в переменную окружения для текущей сессии
set PGPASSWORD=%DB_PASSWORD%

:: --- КОМАНДА 1: Показать список всех таблиц с деталями ---
echo ========================================================
echo 1. Overview of all tables (with size and owner):
echo ========================================================
echo.
:: --pset pager=off отключает пейджер (more/less), чтобы вывод был непрерывным
%PSQL_PATH% --pset pager=off -h %DB_HOST% -p %DB_PORT% -U %DB_USER% -d %DB_NAME% -c "\dt+"

echo.
echo.

:: --- КОМАНДА 2: Показать детальную информацию по всем колонкам ---
echo ========================================================
echo 2. Detailed information for all columns in all tables:
echo ========================================================
echo.
%PSQL_PATH% --pset pager=off -h %DB_HOST% -p %DB_PORT% -U %DB_USER% -d %DB_NAME% -c "SELECT t.table_name, c.column_name, c.data_type, c.is_nullable FROM information_schema.tables AS t JOIN information_schema.columns AS c ON t.table_name = c.table_name WHERE t.table_schema = 'public' ORDER BY t.table_name, c.ordinal_position;"

echo.
echo ========================================================
echo Script finished.
echo ========================================================

echo.
pause