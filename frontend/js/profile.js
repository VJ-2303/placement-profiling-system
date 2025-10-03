// Run on page load
window.onload = function() {
  const token = localStorage.getItem("authToken");
  if (!token) {
    alert("Not logged in. Redirecting...");
    window.location.href = "index.html";
    return;
  }

  fetch("https://placement-profiling-system-production.up.railway.app/profile", {
    method: "GET",
    headers: { "Authorization": "Bearer " + token }
  })
  .then(res => res.ok ? res.json() : Promise.reject("Failed to fetch"))
  .then(data => {
    const student = data.student;
    document.getElementById("userName").innerText = student.name || "Unknown";
    document.getElementById("userEmail").innerText = student.official_email || "N/A";
    document.getElementById("userRoll").innerText = student.id || "N/A";
    document.getElementById("userPhoto").src = student.profile_image_url || "https://via.placeholder.com/120";
  })
  .catch(err => {
    console.error(err);
    localStorage.removeItem("authToken");
    alert("Session expired. Redirecting to login.");
    window.location.href = "index.html";
  });
};

// Navigate to page
function goToPage(page) {
  window.location.href = page;
}

// Logout
function logout() {
  localStorage.removeItem("authToken");
  window.location.href = "index.html";
}

// Placeholder for portfolio download
function downloadPortfolio() {
  alert("Downloading portfolio (API integration required)");
}
