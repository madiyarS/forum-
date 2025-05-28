Forum Application
A web forum built with Go, SQLite, and Docker, featuring user authentication, post/comment creation, category filtering, and a like/dislike system.
Prerequisites

Docker
Docker Compose

Setup

Clone the repository:git clone <repository-url>
cd forum


Make the build script executable:chmod +x build.sh


Build and run the application:./build.sh


Access the application at http://localhost:8080.

Testing

Database: Use SQLite CLI to verify data:sqlite3 data/forum.db "SELECT * FROM users;"
sqlite3 data/forum.db "SELECT * FROM posts;"
sqlite3 data/forum.db "SELECT * FROM comments;"


Docker: Verify build and run:docker images
docker ps -a


Unit Tests: Run Go tests:go test ./tests/...



Features

User registration and login with single-session enforcement
Post and comment creation with category assignment
Like/dislike system with persistent counts
Filtering by category, user-created posts, and liked posts
Responsive frontend with client-side validation
Dockerized deployment

