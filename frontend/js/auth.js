// ============================================
// AUTHENTICATION FUNCTIONS
// ============================================

// Handle OAuth callback parameters
function handleAuthCallback() {
    const params = new URLSearchParams(window.location.search);

    // Check for successful auth callback
    if (params.has('token')) {
        const token = params.get('token');
        const role = params.get('role') || 'student';
        const name = decodeURIComponent(params.get('name') || '');
        const email = decodeURIComponent(params.get('email') || '');
        const userId = params.get('user_id');

        // Store credentials
        localStorage.setItem('token', token);
        localStorage.setItem('user', JSON.stringify({
            id: userId,
            name: name,
            email: email,
            role: role
        }));

        // Clear URL params
        window.history.replaceState({}, document.title, window.location.pathname);

        // Redirect based on role
        if (role === 'admin') {
            window.location.href = 'admin-dashboard.html';
        } else {
            window.location.href = 'profile.html';
        }
        return true;
    }

    // Check for error
    if (params.has('error')) {
        const error = params.get('error');
        const errorDesc = decodeURIComponent(params.get('error_description') || error);
        window.history.replaceState({}, document.title, window.location.pathname);
        return { error: errorDesc };
    }

    return false;
}

// Initiate login flow
function login(role = 'student') {
    window.location.href = `${API_BASE_URL}/auth/login?role=${role}`;
}

// Logout user
function logout() {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    window.location.href = 'index.html';
}

// Verify token is still valid
async function verifyAuth() {
    const token = localStorage.getItem('token');
    if (!token) return null;

    try {
        const response = await fetch(`${API_BASE_URL}/auth/me`, {
            headers: { 'Authorization': `Bearer ${token}` }
        });

        if (response.ok) {
            const data = await response.json();
            return data.user;
        }

        // Token invalid
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        return null;
    } catch (error) {
        console.error('Auth verification failed:', error);
        return null;
    }
}

// Get current user from storage
function getCurrentUser() {
    const userStr = localStorage.getItem('user');
    return userStr ? JSON.parse(userStr) : null;
}

// Check if user is logged in
function isLoggedIn() {
    return !!localStorage.getItem('token');
}

// Require authentication for protected pages
function requireAuth(requiredRole = null) {
    if (!isLoggedIn()) {
        window.location.href = 'index.html';
        return false;
    }

    if (requiredRole) {
        const user = getCurrentUser();
        if (!user || user.role !== requiredRole) {
            window.location.href = 'index.html';
            return false;
        }
    }

    return true;
}

// Check if already logged in and redirect
function checkExistingAuth() {
    if (isLoggedIn()) {
        const user = getCurrentUser();
        if (user?.role === 'admin') {
            window.location.href = 'admin-dashboard.html';
        } else {
            window.location.href = 'profile.html';
        }
        return true;
    }
    return false;
}
