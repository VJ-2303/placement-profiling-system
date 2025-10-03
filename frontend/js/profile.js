// Check token and fetch student info
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
  .then(res => {
    if (!res.ok) throw new Error("Failed to fetch user info");
    return res.json();
  })
  .then(data => {
    const student = data.student;
    document.getElementById("userName").innerText = student.name || "Unknown User";
    document.getElementById("userEmail").innerText = student.official_email || "No email found";
    document.getElementById("userRoll").innerText = student.id || "N/A";
    document.getElementById("userPhoto").src = student.profile_image_url || "https://via.placeholder.com/120";
  })
  .catch(err => {
    console.error(err);
    alert("Session expired or invalid token. Please login again.");
    localStorage.removeItem("authToken");
    window.location.href = "index.html";
  });

  // Collapsible submenu
  const collapsibles = document.getElementsByClassName("collapsible");
  for (let i = 0; i < collapsibles.length; i++) {
    collapsibles[i].addEventListener("click", function() {
      const submenu = this.nextElementSibling;
      submenu.style.display = (submenu.style.display === "flex") ? "none" : "flex";
    });
  }

  // Button navigation
  document.getElementById("btnPersonal").onclick = () => { window.location.href = "personal.html"; };
  document.getElementById("btnAcademic").onclick = () => { window.location.href = "academic.html"; };
  document.getElementById("btnSkills").onclick = () => { window.location.href = "skills.html"; };
  document.getElementById("btnAdditional").onclick = () => { window.location.href = "additional.html"; };
  document.getElementById("btnViewPortfolio").onclick = () => { window.location.href = "portfolio-view.html"; };
  document.getElementById("btnDownloadPortfolio").onclick = () => { alert("Download feature to be implemented"); };
  document.getElementById("btnLogout").onclick = () => { 
    localStorage.removeItem("authToken"); 
    window.location.href = "index.html"; 
  };
};
