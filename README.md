```
█   █ █   █ ███ █   █ ████  █   █      ████   ████  ███  █       
█░  █░█░  █░ █░░██ ██░█░░░█  █ █ ░     █░░░█ █ ░░░░█ ░░█ █░      
█░█ █░█████░░█░░█░█ █░████░░  █ ░ ████ ████░░ ███░░█░ ░█░█░░     
██░██░█░░░█░░█░░█░░░█░█░░░░ ░ █░ ░ ░░░░█░░░░ ░ ░░█ █░░█ ░█░░     
█░░ █░█░░░█░███░█░░ █░█░░░░░  █░░   ░░░█░░░░░████░░ ██ █ █████   
 ░░░ ░░░░  ░░░░░ ░░  ░░░░      ░░       ░░    ░░░░ ░ ░░ ░ ░░░░░  
```
# Whimpy API Server (PostgreSQL Database)

A Go-based REST API server for managing users with PostgreSQL database. The API supports HR, Mentor, and Intern user roles with creation, retrieval, and authentication functionality.

## Project Overview

This project provides a backend API server for user management with three distinct roles:
- **HR** (Role ID: 2) - User management, mentor assignment
- **Mentor** (Role ID: 3) - Participant control, task assignment  
- **Intern/Participant** (Role ID: 4) - Task management, propose assignment

The API handles user creation, retrieval, password hashing with bcrypt, and role-based sign-in with session management via cookies.

## Prerequisites

- Go 1.26.1 or later
- PostgreSQL database
- Git (optional, for cloning)

## Installation

1. **Clone the repository** (or navigate to the project directory):
   ```bash
   cd golang-postgresql-api
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Set up PostgreSQL**:
   - Create a database named `whimpy-db` (or update the connection string)
   - Ensure the tables exist:
     - `wp_hr_list`
     - `wp_mentor_list`
     - `wp_intern_list`
     - `wp_user_session`

4. **Run the server**:
   ```bash
   go run app.go
   ```

   The server will start on `http://localhost:5724`

## Environment Variables

The project uses a PostgreSQL connection string. You can configure it by setting the `DATABASE_URL` environment variable or modifying the `DATABASE_URL` constant in `app.go`.

Example connection string format:
```
postgresql://username:password@host:port/dbname?sslmode=require&channel_binding=require
```

Default configuration reads from the `DATABASE_URL` environment variable. Example (dummy values):
```
postgresql://example_user:example_password@ep-example-host-123456.ap-southeast-1.aws.neon.tech/example_db?sslmode=require&channel_binding=require
```

> **Note:** Never commit real database credentials to this public repository. Set your own `DATABASE_URL` locally (e.g. in a git-ignored `.env` file).

## API Endpoints

### User Creation

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/hr/add` | POST | Create a new HR user |
| `/mentor/add` | POST | Create a new Mentor user |
| `/intern/add` | POST | Create a new Intern user |

### User Retrieval

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/hr/get_all` | GET | Get all HR users |
| `/mentor/get_all` | GET | Get all Mentor users |
| `/intern/get_all` | GET | Get all Intern users |

### User Update (PUT)

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/hr/update/{id}` | PUT | Update HR user info (username, fullname, email, phone) |
| `/mentor/update/{id}` | PUT | Update Mentor user info (username, fullname, email, phone) |
| `/intern/update/{id}` | PUT | Update Intern user info (username, fullname, email, phone) |

### User Password Reset

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/hr/set_new_pass/{id}` | PUT | Set new password for HR user (hashes password with bcrypt) |
| `/mentor/set_new_pass/{id}` | PUT | Set new password for Mentor user (hashes password with bcrypt) |
| `/intern/set_new_pass/{id}` | PUT | Set new password for Intern user (hashes password with bcrypt) |

#### Request Body Format

```json
{
  "userPassword": "newpassword123"
}
```

- **userPassword**: Required. New password string (min 8, max 72 characters)
- The password will be hashed using bcrypt before storage (14 rounds)

#### Response Formats

**Success (200 OK):**
```json
{
  "message": "Password has been Changed."
}
```

**User Not Found (404 Not Found):**
```json
{
  "message": "User Not Found (Or no changes affected)"
}
```

**Bad Request (400 Bad Request):**
```json
{
  "message": "Invalid JSON"
}
```

#### Important Notes

- Passwords are automatically hashed using bcrypt with 14 rounds (same as during user creation)
- The `{id}` path parameter corresponds to the user's database ID (hrID, mtID, or internID)
- After changing a password, the existing session cookie remains valid until its 720-hour expiration
- For security, update sessions invalidated on next login or implement token rotation as needed

### User Deletion (DELETE)

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/hr/delete/{id}` | DELETE | Delete HR user by ID |
| `/mentor/delete/{id}` | DELETE | Delete Mentor user by ID |
| `/intern/delete/{id}` | DELETE | Delete Intern user by ID |

### Authentication (Sign In)

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/hr/signin` | POST | HR user sign in |
| `/mentor/signin` | POST | Mentor user sign in |
| `/intern/signin` | POST | Intern user sign in |

### Root Endpoint

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/` | GET | API welcome message with documentation note |

## Request/Response Formats

### User Creation Request (all roles)

```json
{
  "hrUsername"/"mentorUsername"/"internUsername": "string",
  "hrFullName"/"mentorFullName"/"internFullName": "string",
  "hrEmail"/"mentorEmail"/"internEmail": "string",
  "hrPassword"/"mentorPassword"/"internPassword": "string (min 8, max 72 chars)",
  "hrPhone"/"mentorPhone"/"internPhone": "string"
}
```

### Sign In Request

```json
{
  "userCredential": "username or email",
  "userPassword": "password"
}
```

### Sign In Response

```json
{
  "message": "SessionID created. Info : ",
  "userID": 123,
  "sessionID": "generated-session-token"
}
```

A session cookie named `sessionID` will be set in the response with:
- HttpOnly flag (JavaScript-inaccessible)
- 720-hour (30-day) expiration
- SameSite=Lax CSRF protection

### Success Responses

- **Creation**: `201 Created` with user details
- **Retrieval**: `200 OK` with user list
- **Authentication**: `200 OK` with credential acceptance message

### Error Responses

- `400 Bad Request`: Invalid JSON or validation failures
- `409 Conflict`: Duplicate username or email
- `401 Unauthorized`: Invalid credentials
- `500 Internal Server Error`: Database or server errors

## Database Tables

The application expects the following tables to exist in your PostgreSQL database:

```sql
-- HR Users
CREATE TABLE wp_hr_list (
    hrID SERIAL PRIMARY KEY,
    hrUsername VARCHAR(255) UNIQUE,
    hrFullname VARCHAR(255),
    hrEmail VARCHAR(255) UNIQUE,
    hrPassword VARCHAR(255),
    hrPhone VARCHAR(255),
    roleID INTEGER DEFAULT 2
);

-- Mentor Users
CREATE TABLE wp_mentor_list (
    mtID SERIAL PRIMARY KEY,
    mtUsername VARCHAR(255) UNIQUE,
    mtFullname VARCHAR(255),
    mtEmail VARCHAR(255) UNIQUE,
    mtPassword VARCHAR(255),
    mtPhone VARCHAR(255),
    roleID INTEGER DEFAULT 3
);

-- Intern Users
CREATE TABLE wp_intern_list (
    internID SERIAL PRIMARY KEY,
    internUsername VARCHAR(255) UNIQUE,
    internFullname VARCHAR(255),
    internEmail VARCHAR(255) UNIQUE,
    internPassword VARCHAR(255),
    internPhone VARCHAR(255),
    roleID INTEGER DEFAULT 4
);

-- User Sessions
CREATE TABLE wp_user_session (
    id SERIAL PRIMARY KEY,
    sessionid VARCHAR(255) UNIQUE,
    userID INTEGER,
    senssionExpire TIMESTAMP
);
```

## Development

- **Code structure**: All logic resides in `app.go` - a single-file Go application
- **Dependencies**: 
  - `github.com/gorilla/mux` - HTTP router
  - `github.com/lib/pq` - PostgreSQL driver
  - `golang.org/x/crypto` - bcrypt for password hashing
- **Testing**: Use Postman or curl to test the endpoints

## Deployment

The server listens on port `:5724` by default. For production:

1. Set `Secure: true` in the cookie configuration (requires HTTPS)
2. Use environment variables for database credentials instead of hardcoding
3. Consider adding rate limiting and additional security headers
4. Set up proper database connection pooling and migration strategy

## License

This project is open source. Please check the repository for license details.

## Postman Testing

You can test the API using Postman. Import the collection below or use the individual endpoint descriptions.

### Collection: Whimpy API

**Base URL:** `http://localhost:5724`

### Endpoints

#### 1. Root Endpoint
- **Method:** GET
- **URL:** `http://localhost:5724/`
- **Description:** Welcome message and API documentation note
- **Expected Response:** 200 OK with text output showing ASCII art and available functions

#### 2. Create HR User
- **Method:** POST
- **URL:** `http://localhost:5724/hr/add`
- **Description:** Create a new HR user with username, full name, email, password, and phone number
- **Headers:** `Content-Type: application/json`
- **Body (raw JSON):**
  ```json
  {
    "hrUsername": "testuser",
    "hrFullName": "Test User",
    "hrEmail": "test@example.com",
    "hrPassword": "password123",
    "hrPhone": "1234567890"
  }
  ```
- **Expected Response:** 201 Created with user details, or 409 Conflict if username/email already exists

#### 3. Create Mentor User
- **Method:** POST
- **URL:** `http://localhost:5724/mentor/add`
- **Description:** Create a new Mentor user with username, full name, email, password, and phone number
- **Headers:** `Content-Type: application/json`
- **Body (raw JSON):**
  ```json
  {
    "mentorUsername": "testmentor",
    "mentorFullName": "Test Mentor",
    "mentorEmail": "mentor@example.com",
    "mentorPassword": "password123",
    "mentorPhone": "1234567890"
  }
  ```
- **Expected Response:** 201 Created with user details, or 409 Conflict if username/email already exists

#### 4. Create Intern User
- **Method:** POST
- **URL:** `http://localhost:5724/intern/add`
- **Description:** Create a new Intern user with username, full name, email, password, and phone number
- **Headers:** `Content-Type: application/json`
- **Body (raw JSON):**
  ```json
  {
    "internUsername": "testintern",
    "internFullName": "Test Intern",
    "internEmail": "intern@example.com",
    "internPassword": "password123",
    "internPhone": "1234567890"
  }
  ```
- **Expected Response:** 201 Created with user details, or 409 Conflict if username/email already exists

#### 5. Get All HR Users
- **Method:** GET
- **URL:** `http://localhost:5724/hr/get_all`
- **Description:** Retrieve all HR users from the database
- **Expected Response:** 200 OK with array of HR user objects including ID, username, full name, email, phone, and role name

#### 6. Get All Mentor Users
- **Method:** GET
- **URL:** `http://localhost:5724/mentor/get_all`
- **Description:** Retrieve all Mentor users from the database
- **Expected Response:** 200 OK with array of Mentor user objects including ID, username, full name, email, phone, and role name

#### 7. Get All Intern Users
- **Method:** GET
- **URL:** `http://localhost:5724/intern/get_all`
- **Description:** Retrieve all Intern users from the database
- **Expected Response:** 200 OK with array of Intern user objects including ID, username, full name, email, phone, and role name

### Authentication (Sign In)

#### HR Sign In
- **Method:** POST
- **URL:** `http://localhost:5724/hr/signin`
- **Headers:** `Content-Type: application/json`
- **Body (raw JSON):**
  ```json
  {
    "userCredential": "hrUsername or hrEmail",
    "userPassword": "password"
  }
  ```
- **Expected Response:** 200 OK with session ID set as a cookie, or 401 Unauthorized for invalid credentials

#### Mentor Sign In
- **Method:** POST
- **URL:** `http://localhost:5724/mentor/signin`
- **Headers:** `Content-Type: application/json`
- **Body (raw JSON):**
  ```json
  {
    "userCredential": "mtUsername or mtEmail",
    "userPassword": "password"
  }
  ```
- **Expected Response:** 200 OK with session ID set as a cookie, or 401 Unauthorized for invalid credentials

#### Intern Sign In
- **Method:** POST
- **URL:** `http://localhost:5724/intern/signin`
- **Headers:** `Content-Type: application/json`
- **Body (raw JSON):**
  ```json
  {
    "userCredential": "internUsername or internEmail",
    "userPassword": "password"
  }
  ```
- **Expected Response:** 200 OK with session ID set as a cookie, or 401 Unauthorized for invalid credentials

### Cookie Information

Upon successful sign-in, Postman will receive a cookie named `sessionID` with the following properties:
- **Name:** `sessionID`
- **Value:** Generated session token
- **Expires:** 720 hours (30 days) from creation
- **HttpOnly:** true (JavaScript-inaccessible, prevents XSS theft)
- **Secure:** false (set to true in production with HTTPS)
- **SameSite:** Lax (CSRF protection)
- **Path:** /