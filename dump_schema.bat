@echo off
echo ========================================
echo   SCRIPT FOR DUMPING DATABASE SCHEMA
echo ========================================
echo.

:: --- НАСТРОЙКИ ПОДКЛЮЧЕНИЯ ---
:: Убедитесь, что они совпадают с вашими другими скриптами
set DB_NAME=events
set DB_USER=postgres
set DB_PASSWORD=passwordSuperUser1111
set DB_HOST=localhost
set DB_PORT=5433

:: --- ПОЛНЫЕ ПУТИ К УТИЛИТАМ POSTGRESQL ---
:: Замените "18" на вашу версию PostgreSQL!
:: pg_dump.exe обычно находится в той же папке, что и psql.exe
set PSQL_PATH="C:\Program Files\PostgreSQL\18\bin\psql.exe"
set PG_DUMP_PATH="C:\Program Files\PostgreSQL\18\bin\pg_dump.exe"
:: -----------------------------

echo Checking for pg_dump.exe...
IF NOT EXIST %PG_DUMP_PATH% (
    echo ERROR: pg_dump.exe not found at the specified path:
    echo %PG_DUMP_PATH%
    echo Please check the PG_DUMP_PATH variable in this script.
    echo.
    pause
    exit /b 1
)
echo pg_dump.exe found.
echo.

:: Устанавливаем пароль в переменную окружения для текущей сессии
set PGPASSWORD=%DB_PASSWORD%

:: --- СОЗДАНИЕ ИМЕНИ ФАЙЛА С ДАТОЙ И ВРЕМЕНЕМ (ИСПРАВЛЕННЫЙ МЕТОД) ---
:: Используем стандартные переменные %DATE% и %TIME%
:: и заменяем символы, недопустимые в имени файла (/ : , пробел).
:: Этот метод более надежен, чем wmic, и работает на всех системах Windows.
set "FILE_DATE=%date:/=_%"
set "FILE_TIME=%time::=-%"
set "FILE_TIME=%FILE_TIME: =_%"
set "FILE_TIME=%FILE_TIME:,=_%"
set "TIMESTAMP=%FILE_DATE%_%FILE_TIME%"
set "OUTPUT_FILE=schema_%TIMESTAMP%.sql"

echo ========================================================
echo Dumping schema for database "%DB_NAME%"...
echo Output will be saved to: %OUTPUT_FILE%
echo ========================================================
echo.

:: --- ОСНОВНАЯ КОМАНДА ---
:: --schema-only: выгружает только схему (CREATE TABLE, INDEX и т.д.), без данных
:: --no-owner: не включает команды смены владельца
:: --no-privileges: не включает команды GRANT/REVOKE
%PG_DUMP_PATH% -h %DB_HOST% -p %DB_PORT% -U %DB_USER% -d %DB_NAME% --schema-only --no-owner --no-privileges > %OUTPUT_FILE%

echo.
echo ========================================================
echo Schema dump complete!
echo File saved as: %OUTPUT_FILE%
echo ========================================================

echo.
pause