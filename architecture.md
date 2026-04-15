Project Architecture: Gopher Forum
1. Overview

This project is a web-based forum designed for categorized communication. It follows a Domain-Driven Design (DDD) approach with a Vertical Slice architecture to ensure modularity and ease of collaboration for a 5-member team.
2. Directory Structure

We use a domain-based structure. Each feature (User, Post, Comment) contains its own logic, making it easier to assign tasks without causing merge conflicts.
Plaintext
```
forum/
├── cmd/
│   └── server/
│       └── main.go           # Entry point: Initializes DB and starts the server
├── internal/
│   ├── user/                 # Domain: Registration, Login, Profiles
│   │   ├── handler.go        # HTTP logic (parsing forms)
│   │   ├── service.go        # Business logic (Bcrypt, validation)
│   │   └── repository.go     # SQL queries for 'users' table
│   ├── post/                 # Domain: Threads, Categories, Filtering
│   │   ├── handler.go
│   │   ├── service.go        # Category validation and filtering logic
│   │   └── repository.go     # SQL queries for 'posts' & 'post_categories'
│   ├── comment/              # Domain: Discussion replies
│   │   ├── handler.go
│   │   ├── service.go
│   │   └── repository.go
│   ├── interaction/          # Domain: Likes and Dislikes
│   │   ├── service.go        # Toggle logic (prevents double-voting)
│   │   └── repository.go
│   ├── session/              # Domain: Security and Cookies
│   │   ├── service.go        # UUID generation and expiry check
│   │   └── repository.go
│   ├── models/               # Shared Data Structs (User, Post, etc.)
│   └── database/             # SQLite Driver setup & Migration runner
├── migrations/               # SQL files (001_initial_schema.sql)
├── web/
│   ├── static/               # CSS, Images, Vanilla JS
│   └── templates/            # HTML (Layouts and Pages)
├── .gitignore                # Ignore forum.db, binaries, and .env
├── Dockerfile                # Build instructions
├── docker-compose.yml        # Orchestration
├── go.mod                    # Dependency management
└── README.md                 # Setup instructions
```
3. The "Service" Pattern

The Service Layer acts as the "Brain" of each domain. It sits between the Web Handler and the Database Repository.

    Handler: Receives the HTTP request → passes data to Service.

    Service: Validates data, hashes passwords, applies rules → passes to Repository.

    Repository: Executes the SQL command → saves to forum.db.

4. Database Schema

The database uses SQLite3 and is organized to support many-to-many relationships (for categories) and parent-child relationships (for comments).
Core Tables:

    Users: id, email, username, password_hash.

    Posts: id, user_id, title, content, created_at.

    Categories: id, name.

    Post_Categories: post_id, category_id (Join table).

    Comments: id, post_id, user_id, content.

    Reactions: id, user_id, post_id/comment_id, type (1 or -1).

    Sessions: id (UUID), user_id, expires_at.

5. Security Protocols

    Passwords: Encrypted using Bcrypt (cost factor 10). Never stored as plain text.

    Sessions: Uses UUID v4 stored in HttpOnly cookies to prevent XSS.

    Integrity: Uses ON DELETE CASCADE to ensure that if a post is deleted, its comments and likes are also removed.

    SQL Safety: All database interactions use parameterized queries to prevent SQL injection.

6. Team Collaboration Workflow

    Branching: All features are developed on branches prefixed with feat/ (e.g., feat/auth).

    Commits: Follow the format: type(scope): message (e.g., feat(post): add category filtering).

    PRs: Every Pull Request must be reviewed by at least one teammate before merging to main.

    Environment: All members run the project via Docker to ensure consistency.