// JavaScript for form validation and handling

// Validate registration form
function validateRegister() {
    // Get form inputs
    const username = document.getElementById('username').value.trim();
    const email = document.getElementById('email').value.trim();
    const password = document.getElementById('password').value;
    const errorMessage = document.getElementById('error-message');

    // Clear previous errors
    errorMessage.style.display = 'none';
    errorMessage.textContent = '';

    // Validate inputs
    if (!username || !email || !password) {
        errorMessage.textContent = 'All fields are required';
        errorMessage.style.display = 'block';
        return false;
    }

    // Validate email format
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(email)) {
        errorMessage.textContent = 'Invalid email format';
        errorMessage.style.display = 'block';
        return false;
    }

    // Validate password length
    if (password.length < 6) {
        errorMessage.textContent = 'Password must be at least 6 characters';
        errorMessage.style.display = 'block';
        return false;
    }

    return true;
}

// Validate login form
function validateLogin() {
    // Get form inputs
    const email = document.getElementById('email').value.trim();
    const password = document.getElementById('password').value;
    const errorMessage = document.getElementById('error-message');

    // Clear previous errors
    errorMessage.style.display = 'none';
    errorMessage.textContent = '';

    // Validate inputs
    if (!email || !password) {
        errorMessage.textContent = 'Email and password are required';
        errorMessage.style.display = 'block';
        return false;
    }

    return true;
}

// Validate post creation form and show preview
function validatePost() {
    // Get form inputs
    const title = document.getElementById('title').value.trim();
    const content = document.getElementById('content').value.trim();
    const errorMessage = document.getElementById('error-message');
    const preview = document.getElementById('preview');
    const previewTitle = document.getElementById('preview-title');
    const previewContent = document.getElementById('preview-content');

    // Clear previous errors
    errorMessage.style.display = 'none';
    errorMessage.textContent = '';

    // Validate inputs
    if (!title || !content) {
        errorMessage.textContent = 'Title and content are required';
        errorMessage.style.display = 'block';
        return false;
    }

    // Show preview
    previewTitle.textContent = title;
    previewContent.textContent = content;
    preview.style.display = 'block';

    return true;
}

// Validate comment form
function validateComment() {
    // Get form input
    const content = document.getElementById('comment-content').value.trim();
    const errorMessage = document.getElementById('error-message');

    // Clear previous errors
    if (errorMessage) {
        errorMessage.style.display = 'none';
        errorMessage.textContent = '';
    }

    // Validate input
    if (!content) {
        if (errorMessage) {
            errorMessage.textContent = 'Comment content is required';
            errorMessage.style.display = 'block';
        } else {
            alert('Comment content is required');
        }
        return false;
    }

    return true;
}