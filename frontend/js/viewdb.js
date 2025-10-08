// Wait until DOM is loaded
document.addEventListener("DOMContentLoaded", () => {
  const searchBtn = document.getElementById("searchBtn");
  const rollInput = document.getElementById("rollInput");
  const resultBox = document.getElementById("result");
  const errorMsg = document.getElementById("errorMsg");

  // 🔹 Token Handling — moved inside DOMContentLoaded to avoid reference issues
  const urlParams = new URLSearchParams(window.location.search);
  const tokenFromUrl = urlParams.get("token");

  if (tokenFromUrl) {
    localStorage.setItem("authToken", tokenFromUrl);
  }

  const token = localStorage.getItem("authToken");

  if (!token) {
    alert("Not logged in. Redirecting...");
    window.location.href = "index.html";
    return;
  }

  // Sidebar toggle (for responsive menu)
  const hamburger = document.getElementById("hamburger");
  if (hamburger) {
    hamburger.addEventListener("click", () => {
      document.getElementById("sidebar").classList.toggle("active");
    });
  }

  // Search button click event
  searchBtn.addEventListener("click", async () => {
    const roll = rollInput.value.trim();

    // Clear old data
    resultBox.classList.add("hidden");
    errorMsg.classList.add("hidden");

    if (!roll) {
      showError("⚠️ Please enter a roll number!");
      return;
    }

    // Show loading message
    resultBox.innerHTML = `<p class="loading">Loading student details...</p>`;
    resultBox.classList.remove("hidden");

    try {
      const student = await fetchStudentDetails(roll, token);
      if (!student) throw new Error("No student found");

      displayStudent(student);
    } catch (error) {
      showError("❌ Student not found or server error. Please try again.");
    }
  });
});


// Fetch student details from backend API
async function fetchStudentDetails(roll, token) {
  const endpoint = `https://placement-profiling-system-production.up.railway.app/admin/student/rollno/${roll}`;
  const response = await fetch(endpoint, {
    method: "GET",
    headers: { Authorization: "Bearer " + token },
  });

  if (!response.ok) throw new Error("Network error");
  const data = await response.json();
  return data;
}


// Display the fetched student details
function displayStudent(student) {
  const resultBox = document.getElementById("result");

  // Generate initials for profile avatar
  const initials = student.name
    ? student.name.split(" ").map(n => n[0]).join("").toUpperCase()
    : "NA";

  resultBox.innerHTML = `
    <div class="profile-card">
      <div class="profile-header">
        <div class="profile-avatar">${initials}</div>
        <div class="profile-name">${student.name || "Name Not Provided"}</div>
        <div class="profile-roll">Roll No: ${student.roll_no || "N/A"}</div>
      </div>

      <div class="profile-content">
        ${createPersonalSection(student)}
        ${createAcademicSection(student)}
        ${createContactSection(student)}
      </div>
    </div>
  `;

  resultBox.classList.remove("hidden");
}


// Personal details section
function createPersonalSection(student) {
  return `
    <div class="section">
      <div class="section-header">
        <h3>👤 Personal Information</h3>
      </div>
      <div class="section-content">
        <div class="info-grid">
          ${infoItem("Full Name", student.name)}
          ${infoItem("Roll Number", student.roll_no)}
          ${infoItem("Department", student.department)}
          ${infoItem("Year", student.year)}
          ${infoItem("Gender", student.gender)}
          ${infoItem("Date of Birth", formatDate(student.date_of_birth))}
        </div>
      </div>
    </div>
  `;
}


// Academic details section
function createAcademicSection(student) {
  return `
    <div class="section">
      <div class="section-header">
        <h3>🎓 Academic Information</h3>
      </div>
      <div class="section-content">
        <div class="info-grid">
          ${infoItem("10th Percentage", formatPercentage(student.tenth_percentage))}
          ${infoItem("12th Percentage", formatPercentage(student.twelth_percentage))}
          ${infoItem("CGPA", student.cgpa_overall)}
          ${infoItem("Current Backlogs", student.current_backlogs || "0")}
          ${infoItem("Backlog History", student.has_backlog_history || "No")}
        </div>
      </div>
    </div>
  `;
}


// Contact details section
function createContactSection(student) {
  return `
    <div class="section">
      <div class="section-header">
        <h3>📞 Contact Information</h3>
      </div>
      <div class="section-content">
        <div class="info-grid">
          ${infoItem("Email", student.email)}
          ${infoItem("Phone", student.phone)}
          ${infoItem("Address", student.address)}
          ${infoItem("City", student.city)}
          ${infoItem("Pincode", student.pincode)}
          ${infoItem("LinkedIn", formatLinkedIn(student.linkedin_url))}
        </div>
      </div>
    </div>
  `;
}


// Reusable info item generator
function infoItem(label, value) {
  return `
    <div class="info-item">
      <div class="info-label">${label}</div>
      <div class="info-value">${value || "Not Provided"}</div>
    </div>
  `;
}


// Format date for display
function formatDate(dateString) {
  if (!dateString) return "Not Provided";
  try {
    const date = new Date(dateString);
    return date.toLocaleDateString("en-IN", {
      year: "numeric",
      month: "long",
      day: "numeric"
    });
  } catch {
    return dateString;
  }
}


// Format percentage
function formatPercentage(value) {
  if (!value) return "Not Provided";
  return value.toString().includes("%") ? value : `${value}%`;
}


// Format LinkedIn URL
function formatLinkedIn(url) {
  if (!url) return "Not Provided";
  return `<a href="${url}" target="_blank" rel="noopener noreferrer">View Profile</a>`;
}


// Show error messages
function showError(message) {
  const errorMsg = document.getElementById("errorMsg");
  const resultBox = document.getElementById("result");
  resultBox.classList.add("hidden");

  errorMsg.textContent = message;
  errorMsg.classList.remove("hidden");
}
