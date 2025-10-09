document.addEventListener("DOMContentLoaded", async () => {
  const urlParams = new URLSearchParams(window.location.search);
  const tokenFromUrl = urlParams.get("token");
  const roleFromUrl = urlParams.get("role");

  if (tokenFromUrl) {
    localStorage.setItem("authToken", tokenFromUrl);
    localStorage.setItem("role", roleFromUrl);
    window.history.replaceState({}, document.title, window.location.pathname);
  }

  const token = localStorage.getItem("authToken");
  if (!token) {
    alert("Not logged in. Redirecting...");
    window.location.href = "index.html";
    return;
  }

  // === Fetch Admin Profile ===
  try {
    const res = await fetch(
      "https://placement-profiling-system-production.up.railway.app/admin/profile",
      {
        method: "GET",
        headers: { Authorization: "Bearer " + token },
      },
    );

    if (!res.ok) throw new Error("Failed to fetch admin info");
    const data = await res.json();
    const admin = data.admin;
    const analytics = data.analytics;

    document.getElementById("userName").innerText =
      admin.name || "Unknown User";
    document.getElementById("userEmail").innerText =
      admin.email || "No email found";
    document.getElementById("userPhoto").src =
      "https://via.placeholder.com/120";
    document.getElementById("totalStudent").innerText =
      analytics.total_students;
    document.getElementById("filledFormCount").innerText =
      analytics.profile_completed;
  } catch (err) {
    console.error("Error:", err);
    alert("Session expired or invalid token. Please login again.");
    localStorage.removeItem("authToken");
    window.location.href = "index.html";
  }

  // === Sidebar Navigation ===
  const navMap = {
    btnDashboard: "admin-profile.html",
    btnViewDatabase: "viewdb.html",
  };

  Object.keys(navMap).forEach((btnId) => {
    const btn = document.getElementById(btnId);
    if (btn) {
      btn.addEventListener("click", () => {
        window.location.href = navMap[btnId];
      });
    }
  });

  const btnLogout = document.getElementById("btnLogout");
  if (btnLogout) {
    btnLogout.addEventListener("click", () => {
      localStorage.removeItem("authToken");
      window.location.href = "index.html";
    });
  }

  // === Sidebar Toggle ===
  const hamburger = document.getElementById("hamburger");
  const sidebar = document.getElementById("sidebar");
  const mainContent = document.querySelector(".main-content");

  if (hamburger && sidebar) {
    hamburger.addEventListener("click", () => {
      sidebar.classList.toggle("active");
      if (sidebar.classList.contains("active")) {
        mainContent.style.filter = "blur(4px)";
        mainContent.style.pointerEvents = "none";
      } else {
        mainContent.style.filter = "none";
        mainContent.style.pointerEvents = "auto";
      }
    });
  }
});
