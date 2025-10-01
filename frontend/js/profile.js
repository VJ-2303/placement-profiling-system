window.onload = function() {
  const urlParams = new URLSearchParams(window.location.search);
  const tokenFromUrl = urlParams.get('token');
  
  if (tokenFromUrl) {
    localStorage.setItem("authToken", tokenFromUrl);
    window.history.replaceState({}, document.title, window.location.pathname);
  }
  
  // Now check for token in localStorage
  const token = localStorage.getItem("authToken");
  if (!token) {
    alert("Not logged in. Redirecting...");
    window.location.href = "index.html";
    return;
  }

  fetch("https://placement-profiling-system-production.up.railway.app/profile", {
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
    // Access the student object from the response
    const student = data.student;
    
    // Fill in profile details
    document.getElementById("userName").innerText = student.name || "Unknown User";
    document.getElementById("userEmail").innerText = student.official_email || "No email found";
    document.getElementById("userRoll").innerText = student.id || "N/A"; // Using ID as roll number for now

    // Use the profile image URL from the student data
    document.getElementById("userPhoto").src = student.profile_image_url || "https://via.placeholder.com/120";
  })
  .catch(err => {
    console.error("Error fetching user info:", err);
    alert("Session expired or invalid token. Please login again.");
    localStorage.removeItem("authToken");
    window.location.href = "index.html";
  });
}