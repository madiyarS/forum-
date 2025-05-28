// Main JavaScript for general interactivity

// Function to like or dislike a post
function likePost(postID, isLike) {
    // Create form data for POST request
    const formData = new FormData();
    formData.append('post_id', postID);
    formData.append('is_like', isLike);

    // Send POST request to backend
    fetch(`/post/like`, {
        method: 'POST',
        body: formData
    })
    .then(response => {
        if (response.ok) {
            // Reload page to update like counts
            window.location.reload();
        } else {
            alert('Error updating like');
        }
    })
    .catch(error => {
        console.error('Error:', error);
        alert('Failed to update like');
    });
}

// Function to like or dislike a comment
function likeComment(commentID, isLike) {
    // Create form data for POST request
    const formData = new FormData();
    formData.append('comment_id', commentID);
    formData.append('is_like', isLike);

    // Send POST request to backend
    fetch(`/comment/like`, {
        method: 'POST',
        body: formData
    })
    .then(response => {
        if (response.ok) {
            // Reload page to update like counts
            window.location.reload();
        } else {
            alert('Error updating like');
        }
    })
    .catch(error => {
        console.error('Error:', error);
        alert('Failed to update like');
    });
}

// Show/hide elements based on login status
document.addEventListener('DOMContentLoaded', () => {
    const isLoggedIn = document.querySelector('body').dataset.loggedIn === 'true';
    const createPostLink = document.querySelector('a[href="/create-post"]');
    if (createPostLink && !isLoggedIn) {
        createPostLink.style.display = 'none';
    }
});