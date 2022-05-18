-- Example postgres commands for creating database and user
-- You can use postgres database and user of your choice.
-- Please put these values in .env file.
create database wager_app;

create user wager_app_user with encrypted password 'wagerAppPass';

grant all privileges on database wager_app to wager_app_user;



