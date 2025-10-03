// ---------------- Check Token and Fetch Student Info ----------------
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

    // Fill in profile details
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

  // ---------------- Collapsible Submenu ----------------
  const collapsibles = document.getElementsByClassName("collapsible");
  for (let i = 0; i < collapsibles.length; i++) {
    collapsibles[i].addEventListener("click", function() {
      const submenu = this.nextElementSibling;
      submenu.style.display = (submenu.style.display === "flex") ? "none" : "flex";
    });
  }

  // ---------------- Button Navigation ----------------
  const navMap = {
    btnPersonal: "personal.html",
    btnAcademic: "acadamic.html",
    btnSkills: "skills.html",
    btnAdditional: "additional.html",
    btnViewPortfolio: "portfolio-view.html"
  };

  Object.keys(navMap).forEach(btnId => {
    const btn = document.getElementById(btnId);
    if (btn) {
      btn.onclick = () => { window.location.href = navMap[btnId]; };
    }
  });

  // ---------------- Download Portfolio Placeholder ----------------
  const btnDownload = document.getElementById("btnDownloadPortfolio");
  if (btnDownload) {
    btnDownload.onclick = () => {
      alert("Download feature to be implemented");
      // TODO: Implement API call for PDF download
    };
  }

  // ---------------- Logout ----------------
  const btnLogout = document.getElementById("btnLogout");
  if (btnLogout) {
    btnLogout.onclick = () => {
      localStorage.removeItem("authToken");
      window.location.href = "index.html";
    };
  }
};
