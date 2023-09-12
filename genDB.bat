@echo off
setlocal enabledelayedexpansion

rem Load environment variables from .env file
for /f "usebackq tokens=1,2 delims==" %%G in (".env") do (
  set "%%G=%%H"
)


rem Run cwgo model command with environment variables
cwgo model --db_type mysql --dsn "%DB_USERNAME%:%DB_PASSWORD%@tcp(%DB_HOST%:%DB_PORT%)/%DB_DATABASE%?charset=utf8&parseTime=True&loc=Local" --only_model "true"
