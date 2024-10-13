# Authentication Service

A comprehensive authentication service built with Go, Gin framework, and MySQL. This service provides essential authentication features such as user registration, login, two-factor authentication (2FA), password recovery, and more.

## Table of Contents

- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [API Endpoints](#api-endpoints)
- [Testing](#testing)

## Features

- **User Registration**: Create a new user account with a unique email and username.
- **User Login**: Authenticate users with email and password.
- **Two-Factor Authentication (2FA)**:
  - Enable or disable 2FA for enhanced security.
  - Confirm 2FA codes sent via email.
- **Password Recovery**:
  - Initiate password recovery by sending a recovery code to the user's email.
  - Reset password using the recovery code.
- **User Management**:
  - Update user information.
  - Deactivate user accounts.
- **Security**:
  - Passwords are hashed using bcrypt.
  - Tokens are generated and validated using JWT.
  - Middleware for protected routes.
- **Email Service**:
  - Send emails for verification codes and notifications.

## Prerequisites

- **Go**: Version 1.22 or higher
- **MySQL**: Version 8.2.0 or higher
- **Environment Variables**:
  - Create a `.env` file in the project root.

## Installation

1. **Clone the Repository**

   ```bash
   git clone https://github.com/Renan-Parise/codium-auth.git
   cd codium-auth

## Configuration

1. **Environment Variables**

   Create a `.env` file in the project root and add the following environment variables:

   ```env
    DB_USER=
    DB_PASSWORD=
    DB_NAME=
    DB_HOST=
    DB_PORT=

    JWT_SECRET=

    ELASTIC_APM_SERVER_URL=
    ELASTIC_APM_SERVICE_NAME=
    ELASTIC_APM_ENVIRONMENT=
    ELASTIC_APM_TRANSACTION_SAMPLE_RATE=0

    MAIL_SERVICE_URL=
    ```

2. **Install Dependencies**

   ```bash
   go mod tidy
   ```

## API Endpoints

Public Routes
- `POST /auth/register`: Register a new user.
- `POST /auth/login`: Login with email and password.
- `POST /auth/2fa/confirm`: Confirm 2FA code during login.
- `POST /auth/password/recover`: Initiate password recovery.
- `POST /auth/password/reset`: Reset password using recovery code.

Protected Routes (Require Authentication)
- `PUT /auth/update`: Update user information.
- `DELETE /auth/deactivate`: Deactivate user account.
- `POST /auth/2fa/toggle`: Enable or disable 2FA.
- `POST /auth/2fa/confirm-toggle`: Confirm 2FA code to toggle 2FA setting.

Utility Routes
- `GET /ping`: Health check endpoint.

## Testing

1. **Run Tests**

   ```bash
   go test ./...
   ```