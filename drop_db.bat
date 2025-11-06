@echo off
echo ========================================
echo     SCRIPT FOR DELETING THE DATABASE
echo ========================================
echo.

:: --- НАСТРОЙКИ ПОДКЛЮЧЕНИЯ ---
:: Замените значения на свои (должны совпадать с create_db.bat)
set DB_NAME=events
set DB_USER=postgres
set DB_PASSWORD=passwordSuperUser1111
set DB_HOST=localhost
set DB_PORT=5433

:: --- ПОЛНЫЙ ПУТЬ К PSQL ---
:: Замените "18" на вашу версию PostgreSQL!
set PSQL_PATH="C:\Program Files\PostgreSQL\18\bin\psql.exe"
:: -----------------------------

echo Connecting to PostgreSQL...
echo Attempting to delete database "%DB_NAME%"...

:: Устанавливаем пароль в переменную окружения для текущей сессии
set PGPASSWORD=%DB_PASSWORD%

:: Выполняем SQL-команду для удаления базы данных.
:: Используем "IF EXISTS" чтобы избежать ошибки, если база уже удалена.
%PSQL_PATH% -h %DB_HOST% -p %DB_PORT% -U %DB_USER% -d postgres -c "DROP DATABASE IF EXISTS %DB_NAME%;"

:: Проверяем результат
if %ERRORLEVEL% EQU 0 (
    echo.
    echo SUCCESS: Database "%DB_NAME%" has been deleted successfully.
) else (
    echo.
    echo ERROR: Failed to delete database. Check connection settings, permissions, or if the database is in use.
)

echo.
pause