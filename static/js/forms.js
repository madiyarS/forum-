// JavaScript for form validation and handling

// Validate registration form
function validateRegister() {
    const username = document.getElementById('username').value.trim();
    const email = document.getElementById('email').value.trim();
    const password = document.getElementById('password').value;
    const errorMessage = document.getElementById('error-message');

    errorMessage.style.display = 'none';
    errorMessage.textContent = '';

    if (!username || !email || !password) {
        errorMessage.textContent = 'All fields are required';
        errorMessage.style.display = 'block';
        return false;
    }

    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(email)) {
        errorMessage.textContent = 'Invalid email format';
        errorMessage.style.display = 'block';
        return false;
    }

    if (password.length < 6) {
        errorMessage.textContent = 'Password must be at least 6 characters';
        errorMessage.style.display = 'block';
        return false;
    }

    return true;
}

// Validate login form
function validateLogin() {
    const email = document.getElementById('email').value.trim();
    const password = document.getElementById('password').value;
    const errorMessage = document.getElementById('error-message');

    errorMessage.style.display = 'none';
    errorMessage.textContent = '';

    if (!email || !password) {
        errorMessage.textContent = 'Email and password are required';
        errorMessage.style.display = 'block';
        return false;
    }

    return true;
}

// Validate post creation form and show preview
function validatePost() {
    const title = document.getElementById('title').value.trim();
    const content = document.getElementById('content').value.trim();
    const errorMessage = document.getElementById('error-message');
    const preview = document.getElementById('preview');
    const previewTitle = document.getElementById('preview-title');
    const previewContent = document.getElementById('preview-content');

    errorMessage.style.display = 'none';
    errorMessage.textContent = '';

    if (!title || !content) {
        errorMessage.textContent = 'Title and content are required';
        errorMessage.style.display = 'block';
        return false;
    }

    previewTitle.textContent = title;
    previewContent.textContent = content;
    preview.style.display = 'block';

    return true;
}

// Validate comment form
function validateComment() {
    const content = document.getElementById('comment-content').value.trim();
    const errorMessage = document.getElementById('error-message');

    if (errorMessage) {
        errorMessage.style.display = 'none';
        errorMessage.textContent = '';
    }

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