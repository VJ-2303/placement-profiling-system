const STORAGE_KEY = 'placementPortfolioData';

// --- Utility function to clear profile data cache ---
function clearProfileCache() {
  localStorage.removeItem('cachedProfileData');
  localStorage.removeItem('profileCacheTimestamp');
}

// --- Load Data into Form ---
function loadData(formId) {
  const storedData = localStorage.getItem(STORAGE_KEY);
  if (!storedData) return;

  const data = JSON.parse(storedData);
  const form = document.getElementById(formId);
  if (!form) return;

  form.querySelectorAll('input, select, textarea').forEach(el => {
    if (data[el.name] !== undefined && data[el.name] !== null) {
      if (el.type === 'radio') {
        if (el.value === data[el.name]) el.checked = true;
      } else {
        el.value = data[el.name];
      }
    }
  });
}

// --- Load Data Selectively (only if empty) ---
function loadDataSelectively(formId) {
  const storedData = localStorage.getItem(STORAGE_KEY);
  if (!storedData) return;

  const data = JSON.parse(storedData);
  const form = document.getElementById(formId);
  if (!form) return;

  form.querySelectorAll('input, select, textarea').forEach(el => {
    if (data[el.name] !== undefined && data[el.name] !== null && data[el.name] !== '' && data[el.name] !== '-') {
      if (el.type === 'radio') {
        const isAnyChecked = form.querySelector(`input[name="${el.name}"]:checked`);
        if (!isAnyChecked && el.value === data[el.name]) el.checked = true;
      } else {
        if (!el.value || el.value === '' || el.value === '-') el.value = data[el.name];
      }
    }
  });

  // --- Load photo from localStorage ---
  if (data.photo) {
    const preview = document.getElementById('photoPreview');
    if (preview) {
      preview.src = data.photo;
      preview.style.display = 'inline-block';
    }
  }
}

// --- Save Current Form Data ---
function saveCurrentData(formId) {
  const form = document.getElementById(formId);
  if (!form) return {};

  const formData = new FormData(form);
  const data = {};

  for (const [key, value] of formData.entries()) {
    data[key] = value === '' ? null : value;
  }

  const photoInput = form.querySelector('#photo');
  const existingData = JSON.parse(localStorage.getItem(STORAGE_KEY) || '{}');
  const newData = { ...existingData, ...data };

  if (photoInput && photoInput.files && photoInput.files[0]) {
    const reader = new FileReader();
    reader.onload = () => {
      newData['photo'] = reader.result; // Save base64 string
      localStorage.setItem(STORAGE_KEY, JSON.stringify(newData));
    };
    reader.readAsDataURL(photoInput.files[0]);
  } else {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(newData));
  }

  return newData;
}

// --- Save and Navigate ---
function saveAndNavigate(formId, nextPage) {
  saveCurrentData(formId);
  window.location.href = nextPage;
}

// --- Save and Navigate Back ---
function saveAndNavigateBack(formId, prevPage) {
  saveCurrentData(formId);
  window.location.href = prevPage;
}

// --- Navigation Functions ---
window.saveAndNavigateToAcademic = () => saveAndNavigate('personalForm', 'acadamic.html');
window.saveAndNavigateToSkills = () => saveAndNavigate('academicForm', 'skills.html');
window.saveAndNavigateToPersonal = () => saveAndNavigate('academicForm', 'personal.html');
window.saveAndNavigateToAcademicFromSkills = () => saveAndNavigate('skillsForm', 'acadamic.html');

// --- Initialize Page ---
window.initializePage = (formId) => loadData(formId);

// --- Final Submission to Backend ---
window.finalSubmission = async () => {
  saveCurrentData('personalForm');
  saveCurrentData('academicForm');
  saveCurrentData('skillsForm');

  const finalData = JSON.parse(localStorage.getItem(STORAGE_KEY) || '{}');

  const formattedData = {};
  for (const key in finalData) {
    let value = finalData[key];
    if (value === 'null' || value === null || value === '') value = '-';
    formattedData[key] = value;
  }

  const token = localStorage.getItem('authToken');
  if (!token) {
    alert("You must login first!");
    return;
  }

  try {
    const response = await fetch("https://placement-profiling-system-production.up.railway.app/profile/complete", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${token}`
      },
      body: JSON.stringify(formattedData)
    });

    if (response.ok) {
      alert("Portfolio submitted successfully!");
      localStorage.removeItem(STORAGE_KEY);
      clearProfileCache();
      window.location.href = "profile.html";
    } else {
      const error = await response.json();
      alert("Error submitting portfolio: " + (error.message || response.statusText));
    }
  } catch (err) {
    alert("Network or server error occurred!");
  }
};

// --- Copy Final Data ---
window.copyFinalData = () => {
  const finalData = localStorage.getItem(STORAGE_KEY);
  if (!finalData) return;
  const dummy = document.createElement("textarea");
  document.body.appendChild(dummy);
  dummy.value = finalData;
  dummy.select();
  document.execCommand('copy');
  document.body.removeChild(dummy);
  alert("Data copied to clipboard!");
};

// --- Profile Cache (Global Access) ---
window.clearProfileCache = clearProfileCache;

// --- Get Profile Data (with caching) ---
async function getProfileData() {
  const cachedProfile = localStorage.getItem('cachedProfileData');
  const cacheTimestamp = localStorage.getItem('profileCacheTimestamp');
  const cacheExpiry = 5 * 60 * 1000; // 5 minutes

  if (cachedProfile && cacheTimestamp) {
    const age = Date.now() - parseInt(cacheTimestamp);
    if (age < cacheExpiry) return JSON.parse(cachedProfile);
  }

  const token = localStorage.getItem('authToken');
  if (!token) return null;

  try {
    const response = await fetch('https://placement-profiling-system-production.up.railway.app/profile/complete', {
      method: 'GET',
      headers: { 'Authorization': 'Bearer ' + token }
    });

    if (!response.ok) return null;

    const data = await response.json();
    const profile = data.profile;

    if (profile) {
      localStorage.setItem('cachedProfileData', JSON.stringify(profile));
      localStorage.setItem('profileCacheTimestamp', Date.now().toString());
    }

    return profile;
  } catch {
    return null;
  }
}

// --- Populate Personal Form ---
function populatePersonalForm(profile) {
  if (!profile) return;

  const fields = [
    'name', 'roll_no', 'date_of_birth', 'adhaar_no', 'mobile_number',
    'alt_mobile_number', 'personal_email', 'linkedin_url', 'residence_type',
    'address', 'city', 'pincode', 'father_name', 'father_mobile',
    'father_occupation', 'father_company_details', 'father_email',
    'mother_name', 'mother_mobile', 'mother_occupation', 'mother_email',
    'company_aim', 'target_package'
  ];

  fields.forEach(f => {
    const el = document.getElementById(f);
    if (el && profile[f]) el.value = profile[f];
  });

  // Show photo if exists
  if (profile.photo) {
    const preview = document.getElementById('photoPreview');
    if (preview) {
      preview.src = profile.photo;
      preview.style.display = 'inline-block';
    }
  }
}
document.addEventListener('DOMContentLoaded', () => {
    // Get all necessary elements
    const fileInput = document.getElementById('photoUpload');
    const photoPreview = document.getElementById('photoPreview');
    const previewBtn = document.getElementById('previewBtn');
    const deleteBtn = document.getElementById('deleteBtn');
    const fileStatusDiv = document.getElementById('fileStatus');
    const fileNameDisplay = document.getElementById('fileNameDisplay');
    const viewBtn = document.getElementById('viewBtn');

    // --- Core Functionality: Handle file selection ---
    fileInput.addEventListener('change', (event) => {
        const file = event.target.files[0];
        
        if (file) {
            const reader = new FileReader();

            // Set up the reader to display the image
            reader.onload = (e) => {
                photoPreview.src = e.target.result;
                // photoPreview.style.display = 'block'; // Keep hidden until 'View' is clicked

                // Update UI state for an uploaded file
                previewBtn.style.display = 'none'; // Hide Preview button once a file is selected (now handled by View)
                deleteBtn.style.display = 'block';
                
                // Show file name and View button
                fileInput.style.display = 'none'; // Hide the file input
                fileStatusDiv.style.display = 'flex';
                fileNameDisplay.textContent = file.name;
            };

            reader.readAsDataURL(file);
        } else {
            // If the user cancels the file selection
            handleDelete();
        }
    });

    // --- Preview/View Button Logic ---
    // In a real application, "Preview" and "View" often do the same thing.
    // Here, we use 'View' on the uploaded file and ensure the preview is visible.
    viewBtn.addEventListener('click', () => {
        // Toggle the visibility of the preview image
        if (photoPreview.style.display === 'block') {
            photoPreview.style.display = 'none';
            viewBtn.textContent = 'View';
        } else {
            photoPreview.style.display = 'block';
            viewBtn.textContent = 'Hide';
        }
    });

    // --- Delete Button Logic ---
    deleteBtn.addEventListener('click', handleDelete);

    function handleDelete() {
        // 1. Reset the file input (crucial for re-uploading the same file)
        fileInput.value = ''; 
        
        // 2. Clear the preview image
        photoPreview.src = '';
        photoPreview.style.display = 'none';
        
        // 3. Reset UI state
        previewBtn.style.display = 'none';
        deleteBtn.style.display = 'none';
        fileStatusDiv.style.display = 'none';
        fileInput.style.display = 'block';
        viewBtn.textContent = 'View'; // Reset the view button text
    }

});

// --- Populate Academic Form ---
function populateAcademicForm(profile) {
  if (!profile) return;

  const fields = [
    'tenth_percentage', 'twelth_percentage', 'cgpa_sem1', 'cgpa_sem2',
    'cgpa_sem3', 'cgpa_sem4', 'cgpa_overall', 'current_backlogs'
  ];

  fields.forEach(f => {
    const el = document.getElementById(f);
    if (el && profile[f]) el.value = profile[f];
  });

  if (profile.has_backlog_history) {
    document.querySelectorAll('input[name="has_backlog_history"]').forEach(r => (r.checked = false));
    const rb = document.querySelector(`input[name="has_backlog_history"][value="${profile.has_backlog_history}"]`);
    if (rb) rb.checked = true;
  }
}

// --- Populate Skills Form ---
function populateSkillsForm(profile) {
  if (!profile) return;

  const regularFields = [
    'certifications', 'internships', 'workshops', 'awards',
    'hackathons_attended', 'extracurriculars', 'club_participation',
    'future_path', 'strength', 'weakness', 'remarks'
  ];

  const radioFields = [
    'skill_c', 'skill_cpp', 'skill_java', 'skill_python', 'skill_node_js',
    'skill_php', 'skill_web_dev', 'skill_flutter', 'skill_sql', 'skill_no_sql',
    'concept_data_structures', 'concept_algos', 'concept_dbms', 'concept_oops',
    'concept_os', 'concept_networks', 'concept_problem_solving',
    'tool_git', 'tool_linux', 'tool_cloud', 'tool_hacker_rank', 'tool_hacker_earth',
    'skill_aptitude', 'skill_reasoning', 'communication_skills'
  ];

  regularFields.forEach(f => {
    const el = document.getElementById(f);
    if (el && profile[f]) el.value = profile[f];
  });

  radioFields.forEach(f => {
    const val = profile[f];
    if (val) {
      document.querySelectorAll(`input[name="${f}"]`).forEach(r => (r.checked = false));
      const rb = document.querySelector(`input[name="${f}"][value="${val}"]`);
      if (rb) rb.checked = true;
    }
  });
}

// --- Page Initialization with Profile Data ---
async function initializePageWithProfile(formId) {
  const profile = await getProfileData();

  if (profile) {
    switch (formId) {
      case 'personalForm':
        populatePersonalForm(profile);
        break;
      case 'academicForm':
        populateAcademicForm(profile);
        break;
      case 'skillsForm':
        populateSkillsForm(profile);
        break;
    }
  }

  loadDataSelectively(formId);
}

// --- Preview Photo ---
function previewPhoto(event) {
  const preview = document.getElementById('photoPreview');
  const file = event.target.files[0];
  if (file) {
    preview.src = URL.createObjectURL(file);
    preview.style.display = 'inline-block';
  } else {
    preview.src = '';
    preview.style.display = 'none';
  }
}

// --- Make Functions Global ---
window.getProfileData = getProfileData;
window.initializePageWithProfile = initializePageWithProfile;
window.previewPhoto = previewPhoto;
