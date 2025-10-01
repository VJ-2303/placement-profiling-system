window.onload = function() {
  const token = localStorage.getItem("authToken");
  if (!token) {
    alert("Not logged in. Redirecting...");
    window.location.href = "index.html";
    return;
  }

  fetch("http://localhost:5000/userinfo", {
    method: "GET",
    headers: {
      "Authorization": "Bearer " + token
    }
  })
  .then(res => {
    if (!res.ok) {
      throw new Error("Failed to fetch user info");
    }
    return res.json();
  })
  .then(data => {
    // Fill in profile details
    document.getElementById("userName").innerText = data.name || "Unknown User";
    document.getElementById("userEmail").innerText = data.email || "No email found";
    document.getElementById("userRoll").innerText = data.rollNumber || "N/A";

    // If photo is available, use it; otherwise fallback image
    document.getElementById("userPhoto").src = data.photo || "https://via.placeholder.com/120";
  })
  .catch(err => {
    console.error("Error fetching user info:", err);
    alert("Session expired or invalid token. Please login again.");
    localStorage.removeItem("authToken");
    window.location.href = "index.html";
  });
}
