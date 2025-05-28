Forum Application
A simple, web-based forum application built with Go, SQLite, and a Clean Architecture design. Users can register, log in, create posts, comment, like/dislike content, and filter posts by category or user activity. The frontend uses pure HTML, CSS, and vanilla JavaScript for a lightweight, beginner-friendly interface.
Features

Authentication:
User registration and login with email, username, and password.
Secure session management with single-session enforcement.
Logout functionality.


Posts:
Create posts with titles, content, and multiple category tags.
View all posts or a single post with comments.
Filter posts by category, user-created posts, or liked posts.


Comments:
Add comments to posts.
View comments with like/dislike counts.


Likes/Dislikes:
Like or dislike posts and comments.
Display like/dislike counts.


Frontend:
Responsive design using CSS Grid and Flexbox.
Pure HTML, CSS, and vanilla JavaScript (no frameworks).
Simple, professional color scheme and system fonts.


Backend:
Clean Architecture with entities, use cases, interfaces, and frameworks.
SQLite database for persistent storage.
RESTful API endpoints for all features.



Project Structure
forum/
├── cmd/
│   └── server/
│       └── main.go           # Application entry point
├── internal/
│   ├── entities/            # Core business models (User, Post, etc.)
│   ├── usecases/            # Business logic (auth, posts, comments, likes)
│   ├── interfaces/          # Adapters (HTTP handlers, SQLite repositories)
│   │   ├── database/
│   │   │   └── sqlite/
│   │   ├── http/
│   │   │   ├── handlers/
│   │   │   └── middleware/
│   │   └── repositories/
├── static/
│   ├── css/
│   │   ├── main.css         # Main styles
│   │   └── responsive.css   # Responsive design
│   ├── js/
│   │   ├── main.js          # General interactivity
│   │   └── forms.js         # Form validation
│   └── favicon.ico          # Favicon for browser tabs
├── templates/
│   ├── base.html            # Base template
│   ├── index.html           # Homepage
│   ├── register.html        # Registration page
│   ├── login.html           # Login page
│   ├── create-post.html     # Post creation page
│   ├── post.html            # Post detail page
│   ├── profile.html         # User profile (My Posts, Liked Posts)
│   └── error.html           # Error page
├── go.mod                   # Go module dependencies
├── go.sum
└── README.md                # This file

Prerequisites

Go: Version 1.21 or higher (tested with 1.24.2).
SQLite: Included via github.com/mattn/go-sqlite3.
C Compiler: Xcode Command Line Tools (for cgo, required by go-sqlite3).
macOS: Instructions are tailored for macOS (tested with Apple Clang 17.0.0).

Setup Instructions (macOS)

Clone the Repository:
git clone <repository-url>
cd forum


Install Xcode Command Line Tools:
xcode-select --install
clang --version


Ensure output shows Apple Clang (e.g., version 17.0.0).


Set CGO Environment:
export CGO_ENABLED=1


Install Go Dependencies:
go mod tidy


Place Favicon:

Copy a favicon.ico file to static/:cp /path/to/favicon.ico static/favicon.ico


If you don’t have one, download a sample from Favicon.io or use a placeholder.


Run the Application:
go run cmd/server/main.go


Expected output:2025/05/29 ... Database schema applied successfully with 7 tables
2025/05/29 ... Templates parsed successfully
2025/05/29 ... Server starting on :8080




Access the Forum:

Open http://localhost:8080 in a browser (use Incognito mode to avoid cache issues).
Expect: Homepage with navigation, categories, and “No posts found” message.



Usage

Register: Navigate to /register, enter a username, email, and password (minimum 6 characters).
Login: Go to /login, enter your email and password.
Create Post: Click “New Post” (/create-post), add a title, content, and select categories.
View Post: Click a post title to view details and comments (/post/<id>).
Comment: On a post page, add a comment if logged in.
Like/Dislike: Click like/dislike buttons on posts or comments (disabled if already liked/disliked).
Filter Posts: Use category buttons or navigate to /my-posts or /liked-posts.
Logout: Click “Logout” to end the session.

Testing

Run Unit Tests:
go test ./...


Manual Testing:

Authentication:
Register with duplicate email (expect error).
Login with wrong credentials (expect 401).
Submit empty login form (expect 400).
Login on two browsers, refresh first session (expect redirect to /login).


Posts:
Create a post with multiple categories.
Filter by category (/posts?category=Technology).
View user’s posts (/my-posts) and liked posts (/liked-posts).


Comments: Add a comment, verify it appears.
Likes: Like/dislike a post/comment, check counts.
Database:sqlite3 data/forum.db
.tables
SELECT * FROM users;
SELECT * FROM posts;
SELECT * FROM comments;
SELECT * FROM likes;





Troubleshooting

Black Screen at http://localhost:8080:

Clear Browser Cache:
Chrome: Settings > Privacy > Clear browsing data > Cached images > Clear data.
Use Incognito mode (Command + Shift + N).


Test Server Response:curl http://localhost:8080
curl http://localhost:8080/static/css/main.css
curl http://localhost:8080/static/js/main.js
curl http://localhost:8080/static/favicon.ico


Expect HTML for /, CSS/JS content, and favicon binary.


Check Browser Logs:
Open http://localhost:8080 in Chrome.
Inspect > Console: Note JavaScript errors.
Network tab: Verify requests (/, /static/*, /static/favicon.ico) return 200.


Try Alternative Port:go run cmd/server/main.go -addr=:8082


Update cmd/server/main.go if needed:if err := http.ListenAndServe(":8082", mux); err != nil {
    log.Fatal("Server failed to start:", err)
}




Verify Files:ls -l static/
ls -l templates/


Ensure favicon.ico exists in static/.




Compilation Errors:

Verify imports:grep -r "import" cmd/server/main.go
ls -l internal/interfaces/http/middleware/


Re-run go mod tidy.


CGO Error:

If go run fails:clang --version
xcode-select --install
export CGO_ENABLED=1




Port Conflict:

Check:lsof -i :8080
lsof -i :8082


Kill processes:kill -9 <PID>




Database Issues:

If tables disappear:rm data/forum.db
go run cmd/server/main.go





Contributing

Fork the repository.
Create a feature branch (git checkout -b feature/your-feature).
Commit changes (git commit -m "Add your feature").
Push to the branch (git push origin feature/your-feature).
Open a pull request.

License
This project is licensed under the MIT License.
Contact
For issues or questions, open an issue on the repository or contact the maintainer.
