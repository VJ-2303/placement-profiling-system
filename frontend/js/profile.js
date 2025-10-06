// ---------------- Check Token and Fetch Student Info ----------------
window.onload = function () {
  const urlParams = new URLSearchParams(window.location.search);
  const tokenFromUrl = urlParams.get("token");

  if (tokenFromUrl) {
    localStorage.setItem("authToken", tokenFromUrl);
    // Clean the URL by removing the token parameter
    window.history.replaceState({}, document.title, window.location.pathname);
  }

  const token = localStorage.getItem("authToken");

  if (!token) {
    alert("Not logged in. Redirecting...");
    window.location.href = "index.html";
    return;
  }

  fetch(
    "https://placement-profiling-system-production.up.railway.app/profile",
    {
      method: "GET",
      headers: { Authorization: "Bearer " + token },
    },
  )
    .then((res) => {
      if (!res.ok) throw new Error("Failed to fetch user info");
      return res.json();
    })
    .then((data) => {
      const student = data.student;

      document.getElementById("userName").innerText =
        student.name || "Unknown User";
      document.getElementById("userEmail").innerText =
        student.official_email || "No email found";
      document.getElementById("userRoll").innerText = student.id || "N/A";

      const userPhotoEl = document.getElementById("userPhoto");
      userPhotoEl.src = student.photo || "https://via.placeholder.com/120"; // fallback
    })

    .catch((err) => {
      console.error(err);
      alert("Session expired or invalid token. Please login again.");
      localStorage.removeItem("authToken");
      window.location.href = "index.html";
    });

  // ---------------- Collapsible Submenu ----------------
  const collapsibles = document.getElementsByClassName("collapsible");
  for (let i = 0; i < collapsibles.length; i++) {
    collapsibles[i].addEventListener("click", function () {
      const submenu = this.nextElementSibling;
      submenu.style.display =
        submenu.style.display === "flex" ? "none" : "flex";
    });
  }

  // ---------------- Button Navigation ----------------
  const navMap = {
    btnPersonal: "personal.html",
    btnAcademic: "acadamic.html",
    btnSkills: "skills.html",
    btnViewPortfolio: "view.html",
  };

  Object.keys(navMap).forEach((btnId) => {
    const btn = document.getElementById(btnId);
    if (btn) {
      btn.onclick = () => {
        window.location.href = navMap[btnId];
      };
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

  // ---------------- Hamburger Toggle (New) ----------------
  const hamburger = document.getElementById("hamburger");
  const sidebar = document.getElementById("sidebar");

  if (hamburger && sidebar) {
    hamburger.addEventListener("click", () => {
      sidebar.classList.toggle("active");
    });
  }
};
