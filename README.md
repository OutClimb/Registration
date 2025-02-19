# OutClimb Registration Manager

This application manages registration for OutClimb events. It allows for us to have forms that open and closes at a specific time, while also limiting the amount of submissions to an event.

## Structure

This project user N-Layer architecture that uses three consecutive layers:

- HTTP Layer for handling requests and responses
- App Layer for handling business logic
- Store Layer for handling database transactions

This allows for a seperation of concerns throughout the web application making sure only the necessary data is passed through to each layer.

## Running locally

Make sure you have Go v1.23.4 or newer installed, then run the following from the root of the project directory to install dependencies:

```
go mod download
```

You will also need MariaDB running and configured as well with a empty database. You will then need the following environment variables to configure the database:

- `DB_HOST` - The hostname of your database. (Ex localhost)
- `DB_USERNAME` - The username to login with your database.
- `DB_PASSWORD` - The password to authenticate your user to your database.
- `DB_NAME` - The name of the empty database.

With the environment variables all configured on your system you should then be able to run the following to start the service:

```
go run ./cmd/main.go
```

## Deploying

To deploy you should use the docker image, along with a docker compose file to configure everything. For example:

```
services:
  registration:
    restart: unless-stopped
    build:
      context: .
    environment:
      - DB_HOST=mariadb
      - DB_NAME=${DB_NAME}
      - DB_USERNAME=${DB_USERNAME}
      - DB_PASSWORD=${DB_PASSWORD}

  mariadb:
    restart: unless-stopped
    image: mariadb
    environment:
      - MARIADB_ROOT_PASSWORD=${DB_ROOT_PASSWORD}
      - MARIADB_DATABASE=${DB_NAME}
      - MARIADB_USER=${DB_USERNAME}
      - MARIADB_PASSWORD=${DB_PASSWORD}
    volumes:
      - mariadb_data:/var/lib/mysql

  nginx:
    restart: always
    image: nginx:latest
    ports:
      - 80:80
      - 443:443
    volumes:
      - ./nginx/conf/:/etc/nginx/conf.d/:ro

volumes:
  mariadb_data:
    driver: local
```

This is not the exact config we use, but just an example of what you can do. You will then want a `.env` file for all of your secrets, and to configure Nginx appropriately to proxy the registration application. You may also want to add `TRUSTED_PROXIES` to the registration environment variables.
