function login() {
  // Backend login URL (given by your teammate)
  const loginUrl = "http://localhost:5000/login";  
  window.location.href = loginUrl;
}

// After successful login, backend may redirect back with ?token=...
window.onload = function() {
  const params = new URLSearchParams(window.location.search);
  if (params.has("token")) {
    const token = params.get("token");
    localStorage.setItem("authToken", token);
    window.location.href = "profile.html";  // go to profile page
  }
}
