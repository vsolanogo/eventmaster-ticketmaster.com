Visit hosted app at https://mulhollandexpert.site/

# Docker Setup

1. Make sure you have Node.js 18.0+ and npm installed

2. Ensure that Docker and docker-compose is installed on your system.

3. Navigate to the root of this repository and run the app in Docker using the following commands in your terminal:

`docker-compose build`

`docker-compose up`

4. Once the NestJS app is running in the background, navigate to `/frontend` and execute the following command in your terminal:

`npm install`

5. Run the front-end part of the app in context of react folder with the command:

`npm run dev`

6. Explore the app at [http://localhost:3001/](http://localhost:3001/) in your preferred web browser.

7. Admin user credentials are described in the .env file: Email - `ROOT_ADMIN_EMAIL`, Password - `ROOT_ADMIN_PASSWORD`. Admin profile will be created on initial start. Only user with admin rights can create new events manually.

8. The Ticketmaster.com integration has been implemented, and the script is set to run every 6 hours.

If you encounter issues while running the project in Docker, try using the following commands to clean up. 

`docker system prune -af`

`docker stop $(docker ps -q)`

`docker rmi $(docker images -q)`

`docker rm $(docker ps -a -q)`

Alternatively, you can run the backend part without Docker. For this, ensure you have Mysql or Postgres locally installed. Change type of database in `nestjs/src/app.module.ts` if needed. Make sure that the database is created and has the name configured in `.env.development`. Also, verify that the password and username of your database match. You have to specify your `dbname`, `user`, and `password` in the `.env.development` file.

1. Install NestJS globally on your system using the following command: `npm install -g @nestjs/cli`

2. run `npm run start` in contect of nestjs folder

CORS in the browser is solved with a proxy server described in `vite.config.ts`, making it easier to switch the project from dev to prod mode.



![db](db-shema.png)
