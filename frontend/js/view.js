document.addEventListener("DOMContentLoaded", () => {
  const personalContainer = document.getElementById("personal-details");
  const academicContainer = document.getElementById("academic-details");
  const skillsContainer = document.getElementById("skills-details");
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
      const student = data.student || data; // fallback if backend wraps in "student"

      // Personal
      personalContainer.innerHTML = `
        <p><strong>Name:</strong> ${student.name || "Not Provided"}</p>
        <p><strong>Roll Number:</strong> ${student.rollNumber || "Not Provided"}</p>
        <p><strong>Degree:</strong> ${student.degree || "Not Provided"}</p>
        <p><strong>Department:</strong> ${student.department || "Not Provided"}</p>
        <p><strong>DOB:</strong> ${student.dob || "Not Provided"}</p>
      `;

      // Academic (hidden initially, shown only when More Info clicked)
      if (data.academic) {
        academicContainer.innerHTML = `
          <h3>Academic Details</h3>
          <p><strong>Year:</strong> ${data.academic.year || "Not Provided"}</p>
          <p><strong>GPA:</strong> ${data.academic.gpa || "Not Provided"}</p>
          <p><strong>Achievements:</strong> ${data.academic.achievements || "Not Provided"}</p>
        `;
      }

      // Skills
      if (data.skills && data.skills.length) {
        skillsContainer.innerHTML = `
          <h3>Skills</h3>
          ${data.skills.map(skill => `<p>• ${skill}</p>`).join("")}
        `;
      }
    })
    .catch(err => {
      personalContainer.innerHTML = `<p style="color:red;">Error loading profile</p>`;
      console.error("Error fetching details:", err);
    });

  // More Info toggle (academic + skills)
  moreInfoBtn.addEventListener("click", () => {
    moreInfoSection.classList.toggle("hidden");
    moreInfoBtn.innerText = moreInfoSection.classList.contains("hidden")
      ? "More Info"
      : "Hide Info";
  });
});
