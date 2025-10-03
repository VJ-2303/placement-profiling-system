const STORAGE_KEY = 'placementPortfolioData';

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

// --- Initialize Page ---
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
