document.addEventListener("DOMContentLoaded", () => {
  const personalContainer = document.getElementById("personal-details");
  const academicContainer = document.getElementById("academic-details");
  const skillsContainer = document.getElementById("skills-details");
  const additionalContainer = document.getElementById("additional-details");
  const moreInfoBtn = document.getElementById("moreInfoBtn");
  const moreInfoSection = document.getElementById("more-info");

  const token = localStorage.getItem("authToken");

  if (!token) {
    alert("Session expired. Please login again.");
    window.location.href = "index.html";
    return;
  }

  fetch("https://placement-profiling-system-production.up.railway.app/profile/complete", {
    method: "GET",
    headers: {
      "Authorization": `Bearer ${token}`,
      "Content-Type": "application/json"
    }
  })
    .then(res => {
      if (!res.ok) throw new Error("Failed to fetch profile");
      return res.json();
    })
    .then(data => {
      const student = data.student || data; // fallback if backend wraps in student object

      // Personal Details
      personalContainer.innerHTML = `
        <p><strong>Name:</strong> ${student.name || "N/A"}</p>
        <p><strong>Roll Number:</strong> ${student.rollNumber || "N/A"}</p>
        <p><strong>Degree:</strong> ${student.degree || "N/A"}</p>
        <p><strong>Department:</strong> ${student.department || "N/A"}</p>
        <p><strong>DOB:</strong> ${student.dob || "N/A"}</p>
      `;

      // Academic
      if (data.academic) {
        academicContainer.innerHTML = `
          <p><strong>Year:</strong> ${data.academic.year || "N/A"}</p>
          <p><strong>GPA:</strong> ${data.academic.gpa || "N/A"}</p>
          <p><strong>Achievements:</strong> ${data.academic.achievements || "N/A"}</p>
        `;
      }

      // Skills
      if (data.skills && data.skills.length) {
        skillsContainer.innerHTML = data.skills.map(skill => `<p>• ${skill}</p>`).join("");
      }

      // Additional
      if (data.additional) {
        additionalContainer.innerHTML = `
          <p><strong>Hobbies:</strong> ${data.additional.hobbies || "N/A"}</p>
          <p><strong>Projects:</strong> ${data.additional.projects || "N/A"}</p>
        `;
      }
    })
    .catch(err => {
      personalContainer.innerHTML = `<p style="color:red;">Error loading profile</p>`;
      console.error("Error fetching details:", err);
    });

  // More Info toggle
  moreInfoBtn.addEventListener("click", () => {
    moreInfoSection.classList.toggle("hidden");
    moreInfoBtn.innerText = moreInfoSection.classList.contains("hidden")
      ? "More Info"
      : "Hide Info";
  });

  // Hamburger Menu
  const hamburger = document.querySelector(".hamburger");
  const sideMenu = document.querySelector(".side-menu");

  hamburger.addEventListener("click", () => {
    sideMenu.classList.toggle("active");
  });
});
