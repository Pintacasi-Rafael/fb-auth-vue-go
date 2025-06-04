# fb-auth-vue-go
## Features
- Facebook OAuth 2.0 login (OAuth Redirect Flow)
- Backend: Go + net/http
- Frontend: Vue + Vite
- JWT-based authentication
- MySQL user database
- No Facebook SDK used

## Folder Structure
- `/back`: Go backend (OAuth, JWT, MySQL)
- `/front`: Vue frontend (Facebook login UI)

## Getting Started

1. Clone the repo:
   ```bash
   git clone https://github.com/Pintacasi-Rafael/fb-auth-vue-go.git
   cd fb-auth-vue-go


## Configure the .env file
  ```bash
  DB_USER=[your db user]
  DB_PASS=[your db password]
  DB_NAME=[your db name]
  DB_HOST=127.0.0.1
  DB_PORT=3306
  JWT_SECRET=supersecretkey123
  
  
  FB_APP_ID=[your app id]
  FB_APP_SECRET=[your app secret]
  ```

## Start the backend
  ```bash
  cd back
  go run main.go
  ```

## Start the frontend
  ```bash
  cd ../front
  npm install
  npm run dev
  ```
## Create the database in MySQL
  ```bash
  CREATE DATABASE facebook_auth_db;
  
  USE facebook_auth_db;
  
  CREATE TABLE users (
    facebook_id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255),
    email VARCHAR(255),
    gender VARCHAR(50),
    locale VARCHAR(10),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
  );
  ```
