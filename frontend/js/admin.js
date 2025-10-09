window.onload = function () {
  const urlParams = new URLSearchParams(window.location.search);
  const tokenFromUrl = urlParams.get("token");
  const roleFromUrl = urlParams.get("role");

  if (tokenFromUrl) {
    localStorage.setItem("authToken", tokenFromUrl);
    localStorage.setItem("role", roleFromUrl);
    // Clean the URL
    window.history.replaceState({}, document.title, window.location.pathname);
  }

  const token = localStorage.getItem("authToken");

  if (!token) {
    alert("Not logged in. Redirecting...");
    window.location.href = "index.html";
    return;
  }

  fetch(
    "https://placement-profiling-system-production.up.railway.app/admin/profile",
    {
      method: "GET",
      headers: { Authorization: "Bearer " + token },
    },
  )
    .then((res) => {
      if (!res.ok) throw new Error("Failed to fetch student info");
      return res.json();
    })
    .then((data) => {
      const admin = data.admin;

      document.getElementById("userName").innerText =
        admin.name || "Unknown User";
      document.getElementById("userEmail").innerText =
        admin.email || "No email found";
      document.getElementById("userPhoto").src =
        "https://via.placeholder.com/120";
    })
    .catch((err) => {
      console.error("Error:", err);
      alert("Session expired or invalid token. Please login again.");
      localStorage.removeItem("authToken");
      window.location.href = "index.html";
    });

  const navMap = {
    btnDashboard: "admin-profile.html",
    btnViewDatabase: "viewdb.html",
  };

  Object.keys(navMap).forEach((btnId) => {
    const btn = document.getElementById(btnId);
    if (btn) {
      btn.onclick = () => {
        window.location.href = navMap[btnId];
      };
    }
  });

  const btnLogout = document.getElementById("btnLogout");
  if (btnLogout) {
    btnLogout.onclick = () => {
      localStorage.removeItem("authToken");
      window.location.href = "index.html";
    };
  }

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
};
// Example placeholders — replace API URLs with your actual endpoints
document.addEventListener("DOMContentLoaded", async () => {
    try {
        // Fetch total students
        document.getElementById("totalStudents").textContent = 350; // Static example

        // Fetch students who filled the form
        const res = await fetch("/api/students/filled-form"); // <-- Your backend endpoint
        const data = await res.json();
        document.getElementById("filledFormCount").textContent = data.count || 0;

        // Render Performance Analytics Chart
        renderPerformanceChart(data.performance || []);
    } catch (error) {
        console.error("Error loading dashboard data:", error);
    }
});

// ===== Chart.js for Performance Analytics =====
function renderPerformanceChart(performanceData) {
    const ctx = document.getElementById("performanceChart").getContext("2d");

    new Chart(ctx, {
        type: "bar",
        data: {
            labels: performanceData.map(item => item.category),
            datasets: [{
                label: "Performance Score",
                data: performanceData.map(item => item.score),
                backgroundColor: "#4c8bf5"
            }]
        },
        options: {
            responsive: true,
            plugins: {
                legend: { display: false }
            },
            scales: {
                y: {
                    beginAtZero: true,
                    grid: { color: "rgba(255,255,255,0.1)" },
                    ticks: { color: "#fff" }
                },
                x: {
                    grid: { display: false },
                    ticks: { color: "#fff" }
                }
            }
        }
    });
}
