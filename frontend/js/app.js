// app.js
// Single script file containing all application logic (storage, load, and navigation)

const STORAGE_KEY = 'placementPortfolioData';

// --- Utility Functions: Data Persistence and Retrieval ---

/**
 * Loads data from localStorage and populates the form fields.
 * @param {string} formId - The ID of the form element.
 */
function loadData(formId) {
    const storedData = localStorage.getItem(STORAGE_KEY);
    if (storedData) {
        const data = JSON.parse(storedData);
        const form = document.getElementById(formId);
        
        // Iterate over all form elements and populate them
        form.querySelectorAll('input, select, textarea').forEach(element => {
            if (data[element.name] !== undefined && data[element.name] !== null) {
                if (element.type === 'radio') {
                    if (element.value === data[element.name]) {
                        element.checked = true;
                    }
                } else {
                    element.value = data[element.name];
                }
            }
        });
    }
}

/**
 * Collects data from the specified form and saves it to localStorage.
 * @param {string} formId - The ID of the form element.
 * @returns {Object} The collected data object.
 */
function saveCurrentData(formId) {
    const form = document.getElementById(formId);
    const formData = new FormData(form);
    const data = {};

    // Collect form data
    for (const [key, value] of formData.entries()) {
        // Sanitize input: replace empty strings with null if needed
        data[key] = value === '' ? null : value; 
    }

    // Merge with existing data in localStorage (if any)
    const existingData = JSON.parse(localStorage.getItem(STORAGE_KEY) || '{}');
    const newData = { ...existingData, ...data };
    
    localStorage.setItem(STORAGE_KEY, JSON.stringify(newData));
    return newData;
}

/**
 * Saves current form data and navigates to the next page.
 * @param {string} formId - The ID of the form element.
 * @param {string} nextPage - The HTML file to navigate to.
 */
function saveAndNavigate(formId, nextPage) {
    saveCurrentData(formId);
    window.location.href = nextPage;
}

// --- Page-Specific Navigation Functions (Exposed to HTML) ---

// PERSONAL PAGE (personal.html) navigation
window.saveAndNavigateToAcademic = () => {
    saveAndNavigate('personalForm', 'academic.html');
};

// ACADEMIC PAGE (academic.html) navigation
window.saveAndNavigateToSkills = () => {
    saveAndNavigate('academicForm', 'skills.html');
};

window.saveAndNavigateToPersonal = () => {
    saveAndNavigate('academicForm', 'personal.html');
};

// SKILLS PAGE (skills.html) navigation
window.saveAndNavigateToAcademicFromSkills = () => {
    saveAndNavigate('skillsForm', 'academic.html');
};

/**
 * Consolidates all data, formats it, displays it, and handles final submission.
 */
window.finalSubmission = () => {
    // 1. Save any final changes from the current page
    saveCurrentData('skillsForm');
    
    // 2. Retrieve the complete JSON object
    const finalData = JSON.parse(localStorage.getItem(STORAGE_KEY) || '{}');
    
    // 3. Format/Clean up the final object
    const formattedData = {};
    for (const key in finalData) {
        let value = finalData[key];
        
        // Handle explicit null string replacement
        if (value === 'null') {
            value = null;
        } 
        // Replace genuinely empty strings (from non-required fields left blank) with '-'
        else if (typeof value === 'string' && value.trim() === '') {
            value = '-';
        }

        formattedData[key] = value;
    }

    const jsonOutput = JSON.stringify(formattedData, null, 4);

    // 4. Display the result in the dedicated area
    const resultArea = document.getElementById('resultArea');
    const outputContainer = document.getElementById('outputContainer');
    const mainFormContainer = document.getElementById('mainFormContainer');
    const copyButton = document.getElementById('copyButton');
    const copyMessage = document.getElementById('copyMessage');

    if (resultArea && outputContainer && mainFormContainer) {
        resultArea.textContent = jsonOutput;
        outputContainer.classList.remove('hidden');
        mainFormContainer.classList.add('hidden');

        // Logic for copying JSON to clipboard
        copyButton.onclick = () => {
             // Create a dummy textarea to hold the text
            const dummy = document.createElement("textarea");
            document.body.appendChild(dummy);
            dummy.value = jsonOutput;
            dummy.select();
            // Use document.execCommand('copy') for better iframe compatibility
            document.execCommand('copy');
            document.body.removeChild(dummy);
            
            copyMessage.classList.remove('hidden');
            setTimeout(() => {
                copyMessage.classList.add('hidden');
            }, 2000);
        };
        
        // Print to console for visibility
        console.log("--- FINAL CONSOLIDATED JSON OUTPUT ---");
        console.log(formattedData);
        console.log("--------------------------------------");

    } else {
        console.error("DOM elements for output not found.");
    }
};


// --- Initialization ---

// Global function to initialize the correct page on load
window.initializePage = (formId) => {
    loadData(formId);
};
