function login() {
  const loginUrl = "https://placement-profiling-system-production.up.railway.app/auth/login";
  window.location.href = loginUrl;
}

window.onload = function() {
  const params = new URLSearchParams(window.location.search);

  
  if (params.has("token")) {
    const token = params.get("token");
    const role = params.get("role") || "user"; 

    localStorage.setItem("authToken", token);
    localStorage.setItem("userRole", role);

    window.history.replaceState({}, document.title, window.location.pathname);

    if (role.toLowerCase() === "admin") {
      window.location.href = "admin-profile.html";
    } else {
      window.location.href = "profile.html";
    }
  }
};
