const STORAGE_KEY = 'placementPortfolioData';

// --- Utility function to clear profile data cache ---
function clearProfileCache() {
    localStorage.removeItem('cachedProfileData');
    localStorage.removeItem('profileCacheTimestamp');
}

// --- Load Data into Form ---
function loadData(formId) {
    const storedData = localStorage.getItem(STORAGE_KEY);
    if (storedData) {
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
}

// --- Load Data Selectively (only for empty fields) ---
function loadDataSelectively(formId) {
    const storedData = localStorage.getItem(STORAGE_KEY);
    if (storedData) {
        const data = JSON.parse(storedData);
        const form = document.getElementById(formId);
        if (!form) return;
        
        form.querySelectorAll('input, select, textarea').forEach(el => {
            // Only load localStorage data if field is empty and localStorage has valid data
            if (data[el.name] !== undefined && data[el.name] !== null && data[el.name] !== '' && data[el.name] !== '-') {
                if (el.type === 'radio') {
                    // For radio buttons, only set if no radio is currently selected
                    const isAnyChecked = form.querySelector(`input[name="${el.name}"]:checked`);
                    if (!isAnyChecked && el.value === data[el.name]) {
                        el.checked = true;
                    }
                } else {
                    // For other inputs, only set if current value is empty
                    if (!el.value || el.value === '' || el.value === '-') {
                        el.value = data[el.name];
                    }
                }
            }
        });
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
    const existingData = JSON.parse(localStorage.getItem(STORAGE_KEY) || '{}');
    const newData = { ...existingData, ...data };
    localStorage.setItem(STORAGE_KEY, JSON.stringify(newData));
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

// --- Initialize Page (backward compatibility) ---
window.initializePage = (formId) => loadData(formId);

// --- Final Submission to Backend ---
window.finalSubmission = async () => {
    // Save all forms before submission
    saveCurrentData('personalForm');
    saveCurrentData('academicForm');
    saveCurrentData('skillsForm');

    const finalData = JSON.parse(localStorage.getItem(STORAGE_KEY) || '{}');

    // Replace empty/null values with "-"
    const formattedData = {};
    for (const key in finalData) {
        let value = finalData[key];
        if (value === 'null' || value === null || value === '') value = '-';
        formattedData[key] = value;
    }

    // Get auth token (this is stored separately in localStorage, not merged into finalData)
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
                "Authorization": `Bearer ${token}`   // token goes in headers, not in body
            },
            body: JSON.stringify(formattedData)
        });

        if (response.ok) {
            alert("Portfolio submitted successfully!");
            console.log("--- SUBMITTED DATA ---", formattedData);
            localStorage.removeItem(STORAGE_KEY); // clear saved data
            clearProfileCache(); // clear cached profile data so it gets fresh data next time
            window.location.href = "profile.html"; // redirect after success
        } else {
            const error = await response.json();
            alert("Error submitting portfolio: " + (error.message || response.statusText));
        }
    } catch (err) {
        console.error(err);
        alert("Network or server error occurred!");
    }
};

// --- Optional: Copy JSON for debugging ---
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

// --- Make clearProfileCache available globally ---
window.clearProfileCache = clearProfileCache;

// --- Centralized function to get profile data (cached or fresh) ---
async function getProfileData() {
    // Check if we have cached profile data
    const cachedProfile = localStorage.getItem('cachedProfileData');
    const cacheTimestamp = localStorage.getItem('profileCacheTimestamp');
    const cacheExpiry = 5 * 60 * 1000; // 5 minutes cache
    
    // Use cached data if it exists and is not expired
    if (cachedProfile && cacheTimestamp) {
        const age = Date.now() - parseInt(cacheTimestamp);
        if (age < cacheExpiry) {
            return JSON.parse(cachedProfile);
        }
    }
    
    // Fetch fresh data from API
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
        
        // Cache the profile data
        if (profile) {
            localStorage.setItem('cachedProfileData', JSON.stringify(profile));
            localStorage.setItem('profileCacheTimestamp', Date.now().toString());
        }
        
        return profile;
    } catch (error) {
        console.log('Error fetching profile data:', error);
        return null;
    }
}

// --- Populate Personal Form Fields ---
function populatePersonalForm(profile) {
    if (!profile) return;
    
    // Personal information fields
    const personalFields = [
        'name', 'roll_no', 'date_of_birth', 'adhaar_no', 'mobile_number',
        'alt_mobile_number', 'personal_email', 'linkedin_url', 'residence_type',
        'address', 'city', 'pincode'
    ];
    
    personalFields.forEach(fieldName => {
        const element = document.getElementById(fieldName);
        if (element) {
            // Always use database value if it exists, otherwise keep current value
            element.value = profile[fieldName] || element.value || '';
        }
    });

    // Father's details
    const fatherFields = [
        'father_name', 'father_mobile', 'father_occupation', 
        'father_company_details', 'father_email'
    ];
    
    fatherFields.forEach(fieldName => {
        const element = document.getElementById(fieldName);
        if (element) {
            element.value = profile[fieldName] || element.value || '';
        }
    });

    // Mother's details
    const motherFields = [
        'mother_name', 'mother_mobile', 'mother_occupation', 'mother_email'
    ];
    
    motherFields.forEach(fieldName => {
        const element = document.getElementById(fieldName);
        if (element) {
            element.value = profile[fieldName] || element.value || '';
        }
    });
}

// --- Populate Academic Form Fields ---
function populateAcademicForm(profile) {
    if (!profile) return;
    
    // Academic fields
    const academicFields = [
        'tenth_percentage', 'twelth_percentage', 'cgpa_sem1', 'cgpa_sem2',
        'cgpa_sem3', 'cgpa_sem4', 'cgpa_overall', 'current_backlogs'
    ];
    
    academicFields.forEach(fieldName => {
        const element = document.getElementById(fieldName);
        if (element) {
            // Always use database value if it exists, otherwise keep current value
            element.value = profile[fieldName] || element.value || '';
        }
    });
    
    // Handle radio buttons for backlog history
    if (profile.has_backlog_history) {
        // Clear any existing selection first
        document.querySelectorAll('input[name="has_backlog_history"]').forEach(radio => {
            radio.checked = false;
        });
        // Set the correct value
        const radioButton = document.querySelector(`input[name="has_backlog_history"][value="${profile.has_backlog_history}"]`);
        if (radioButton) radioButton.checked = true;
    }
}

// --- Populate Skills Form Fields ---
function populateSkillsForm(profile) {
    if (!profile) return;
    
    // Career aspirations
    const aspirationFields = ['company_aim', 'target_package'];
    aspirationFields.forEach(fieldName => {
        const element = document.getElementById(fieldName);
        if (element) {
            element.value = profile[fieldName] || element.value || '';
        }
    });
    
    // Certifications & experience
    const experienceFields = ['certifications', 'internships', 'workshops', 'awards'];
    experienceFields.forEach(fieldName => {
        const element = document.getElementById(fieldName);
        if (element) {
            element.value = profile[fieldName] || element.value || '';
        }
    });
    
    // All skill dropdowns
    const skillFields = [
        'skill_c', 'skill_cpp', 'skill_java', 'skill_python', 'skill_node_js',
        'skill_php', 'skill_web_dev', 'skill_flutter', 'skill_sql', 'skill_no_sql',
        'concept_data_structures', 'concept_algos', 'concept_dbms', 'concept_oops',
        'concept_os', 'concept_networks', 'tool_git', 'tool_linux', 'tool_cloud',
        'concept_problem_solving', 'tool_hacker_rank', 'tool_hacker_earth',
        'skill_aptitude', 'skill_reasoning', 'communication_skills'
    ];
    
    skillFields.forEach(fieldName => {
        const element = document.getElementById(fieldName);
        if (element) {
            // For dropdowns, use database value if it exists, otherwise keep current selection
            if (profile[fieldName]) {
                element.value = profile[fieldName];
            }
        }
    });
    
    // Communication & extracurricular fields (optional fields)
    const optionalFields = [
        'hackathons_attended', 'extracurriculars', 'club_participation',
        'future_path', 'strength', 'weakness', 'remarks'
    ];
    
    optionalFields.forEach(fieldName => {
        const element = document.getElementById(fieldName);
        if (element) {
            element.value = profile[fieldName] || element.value || '';
        }
    });
}

// --- Enhanced page initialization with profile data ---
async function initializePageWithProfile(formId) {
    console.log('Initializing page with profile for:', formId);
    
    // First fetch profile data from database
    const profile = await getProfileData();
    console.log('Fetched profile data:', profile);
    
    if (profile) {
        // Populate with database data first (priority)
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
        console.log('Populated form with database data');
    }
    
    // Then load localStorage data only for fields that are still empty
    // This preserves any unsaved changes while prioritizing database data
    loadDataSelectively(formId);
}

// --- Make functions available globally ---
window.getProfileData = getProfileData;
window.initializePageWithProfile = initializePageWithProfile;
