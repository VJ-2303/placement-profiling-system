
document.addEventListener('DOMContentLoaded', async () => {
  const loading = document.getElementById('loading');
  const profileContainer = document.getElementById('profile-container');
  
  try {
    
    const profile = await getProfileData();
    
    if (!profile) {
      throw new Error('No profile data found');
    }
    
   
    displayProfile(profile);
    
    
    loading.style.display = 'none';
    profileContainer.style.display = 'block';
    
  } catch (error) {
    
    loading.innerHTML = `
      <div class="error">
        <h3>❌ Error Loading Profile</h3>
        <p>Please make sure you're logged in and try again.</p>
        <button class="btn" onclick="window.location.href='index.html'">Go to Login</button>
      </div>
    `;
  }
});


function displayProfile(profile) {
  const container = document.getElementById('profile-container');
  
  const nameInitials = profile.name 
    ? profile.name.split(' ').map(n => n[0]).join('').toUpperCase() 
    : 'NA';
  
 
  container.innerHTML = `
    <div class="profile-card">
      <div class="profile-header">
        <div class="profile-avatar">${nameInitials}</div>
        <div class="profile-name">${profile.name || 'Name Not Set'}</div>
        <div class="profile-roll">Roll No: ${profile.roll_no || 'Not Set'}</div>
      </div>
      
      <div class="profile-content">
        ${createPersonalSection(profile)}
        ${createAcademicSection(profile)}
        ${createSkillsSection(profile)}
        ${createAdditionalSection(profile)}
      </div>
    </div>
  `;
  
  initializeSectionToggle();
}


function initializeSectionToggle() {
  document.querySelectorAll('.section-header').forEach(header => {
    header.addEventListener('click', () => {
      const section = header.parentElement;
      section.classList.toggle('active');
    });
  });
  
  
  const firstSection = document.querySelector('.section');
  if (firstSection) {
    firstSection.classList.add('active');
  }
}


function createPersonalSection(profile) {
  return `
    <div class="section active">
      <div class="section-header">
        <h3>👤 Personal Information</h3>
        <span class="section-toggle">▼</span>
      </div>
      <div class="section-content">
        <div class="info-grid">
          <div class="info-item">
            <div class="info-label">Full Name</div>
            <div class="info-value">${profile.name || 'Not provided'}</div>
          </div>
          <div class="info-item">
            <div class="info-label">Roll Number</div>
            <div class="info-value">${profile.roll_no || 'Not provided'}</div>
          </div>
          <div class="info-item">
            <div class="info-label">Date of Birth</div>
            <div class="info-value">${formatDate(profile.date_of_birth) || 'Not provided'}</div>
          </div>
          <div class="info-item">
            <div class="info-label">Mobile Number</div>
            <div class="info-value">${profile.mobile_number || 'Not provided'}</div>
          </div>
          <div class="info-item">
            <div class="info-label">Alternate Mobile</div>
            <div class="info-value">${profile.alt_mobile_number || 'Not provided'}</div>
          </div>
          <div class="info-item">
            <div class="info-label">Personal Email</div>
            <div class="info-value">${profile.personal_email || 'Not provided'}</div>
          </div>
          <div class="info-item">
            <div class="info-label">Official Email</div>
            <div class="info-value">${profile.official_email || 'Not provided'}</div>
          </div>
          <div class="info-item">
            <div class="info-label">LinkedIn Profile</div>
            <div class="info-value">${formatLinkedIn(profile.linkedin_url)}</div>
          </div>
          <div class="info-item">
            <div class="info-label">Residence Type</div>
            <div class="info-value">${profile.residence_type || 'Not provided'}</div>
          </div>
          <div class="info-item">
            <div class="info-label">Address</div>
            <div class="info-value">${profile.address || 'Not provided'}</div>
          </div>
          <div class="info-item">
            <div class="info-label">City</div>
            <div class="info-value">${profile.city || 'Not provided'}</div>
          </div>
          <div class="info-item">
            <div class="info-label">Pincode</div>
            <div class="info-value">${profile.pincode || 'Not provided'}</div>
          </div>
        </div>
        
        <h4 class="family-heading">👨‍👩‍👧‍👦 Family Information</h4>
        <div class="info-grid">
          <div class="info-item">
            <div class="info-label">Father's Name</div>
            <div class="info-value">${profile.father_name || 'Not provided'}</div>
          </div>
          <div class="info-item">
            <div class="info-label">Father's Mobile</div>
            <div class="info-value">${profile.father_mobile || 'Not provided'}</div>
          </div>
          <div class="info-item">
            <div class="info-label">Father's Occupation</div>
            <div class="info-value">${profile.father_occupation || 'Not provided'}</div>
          </div>
          <div class="info-item">
            <div class="info-label">Father's Company</div>
            <div class="info-value">${profile.father_company_details || 'Not provided'}</div>
          </div>
          <div class="info-item">
            <div class="info-label">Father's Email</div>
            <div class="info-value">${profile.father_email || 'Not provided'}</div>
          </div>
          <div class="info-item">
            <div class="info-label">Mother's Name</div>
            <div class="info-value">${profile.mother_name || 'Not provided'}</div>
          </div>
          <div class="info-item">
            <div class="info-label">Mother's Mobile</div>
            <div class="info-value">${profile.mother_mobile || 'Not provided'}</div>
          </div>
          <div class="info-item">
            <div class="info-label">Mother's Occupation</div>
            <div class="info-value">${profile.mother_occupation || 'Not provided'}</div>
          </div>
          <div class="info-item">
            <div class="info-label">Mother's Email</div>
            <div class="info-value">${profile.mother_email || 'Not provided'}</div>
          </div>
        </div>
      </div>
    </div>
  `;
}


function createAcademicSection(profile) {
  return `
    <div class="section">
      <div class="section-header">
        <h3>🎓 Academic Information</h3>
        <span class="section-toggle">▼</span>
      </div>
      <div class="section-content">
        <div class="info-grid">
          <div class="info-item">
            <div class="info-label">10th Percentage</div>
            <div class="info-value">${formatPercentage(profile.tenth_percentage)}</div>
          </div>
          <div class="info-item">
            <div class="info-label">12th Percentage</div>
            <div class="info-value">${formatPercentage(profile.twelth_percentage)}</div>
          </div>
          <div class="info-item">
            <div class="info-label">Semester 1 CGPA</div>
            <div class="info-value">${profile.cgpa_sem1 || 'Not provided'}</div>
          </div>
          <div class="info-item">
            <div class="info-label">Semester 2 CGPA</div>
            <div class="info-value">${profile.cgpa_sem2 || 'Not provided'}</div>
          </div>
          <div class="info-item">
            <div class="info-label">Semester 3 CGPA</div>
            <div class="info-value">${profile.cgpa_sem3 || 'Not provided'}</div>
          </div>
          <div class="info-item">
            <div class="info-label">Semester 4 CGPA</div>
            <div class="info-value">${profile.cgpa_sem4 || 'Not provided'}</div>
          </div>
          <div class="info-item">
            <div class="info-label">Overall CGPA</div>
            <div class="info-value">${profile.cgpa_overall || 'Not provided'}</div>
          </div>
          <div class="info-item">
            <div class="info-label">Current Backlogs</div>
            <div class="info-value">${profile.current_backlogs || '0'}</div>
          </div>
          <div class="info-item">
            <div class="info-label">Backlog History</div>
            <div class="info-value">${profile.has_backlog_history || 'Not specified'}</div>
          </div>
        </div>
      </div>
    </div>
  `;
}


function createSkillsSection(profile) {
  const programmingSkills = [
    {name: 'C', level: profile.skill_c},
    {name: 'C++', level: profile.skill_cpp},
    {name: 'Java', level: profile.skill_java},
    {name: 'Python', level: profile.skill_python},
    {name: 'Node.js', level: profile.skill_node_js},
    {name: 'PHP', level: profile.skill_php},
    {name: 'Web Development', level: profile.skill_web_dev},
    {name: 'Flutter', level: profile.skill_flutter},
    {name: 'SQL', level: profile.skill_sql},
    {name: 'NoSQL', level: profile.skill_no_sql}
  ].filter(skill => skill.level);
  
  const concepts = [
    {name: 'Data Structures', level: profile.concept_data_structures},
    {name: 'Algorithms', level: profile.concept_algos},
    {name: 'DBMS', level: profile.concept_dbms},
    {name: 'OOPS', level: profile.concept_oops},
    {name: 'Operating Systems', level: profile.concept_os},
    {name: 'Computer Networks', level: profile.concept_networks},
    {name: 'Problem Solving', level: profile.concept_problem_solving}
  ].filter(skill => skill.level);
  
  const tools = [
    {name: 'Git/GitHub', level: profile.tool_git},
    {name: 'Linux/Unix', level: profile.tool_linux},
    {name: 'Cloud Platforms', level: profile.tool_cloud},
    {name: 'HackerRank', level: profile.tool_hacker_rank},
    {name: 'HackerEarth', level: profile.tool_hacker_earth}
  ].filter(skill => skill.level);
  
  const otherSkills = [
    {name: 'Aptitude', level: profile.skill_aptitude},
    {name: 'Reasoning', level: profile.skill_reasoning},
    {name: 'Communication', level: profile.communication_skills}
  ].filter(skill => skill.level);
  
  return `
    <div class="section">
      <div class="section-header">
        <h3>💻 Technical Skills</h3>
        <span class="section-toggle">▼</span>
      </div>
      <div class="section-content">
        <div class="skills-grid">
          ${programmingSkills.length ? createSkillCategory('Programming Languages', programmingSkills) : ''}
          ${concepts.length ? createSkillCategory('Core Concepts', concepts) : ''}
          ${tools.length ? createSkillCategory('Tools & Platforms', tools) : ''}
          ${otherSkills.length ? createSkillCategory('Other Skills', otherSkills) : ''}
        </div>
        ${(!programmingSkills.length && !concepts.length && !tools.length && !otherSkills.length) 
          ? '<p style="text-align: center; color: #6c757d; font-style: italic;">No skills information available</p>' 
          : ''}
      </div>
    </div>
  `;
}


function createSkillCategory(title, skills) {
  return `
    <div class="skill-category">
      <h4>${title}</h4>
      ${skills.map(skill => `
        <div class="skill-item">
          <span class="skill-name">${skill.name}</span>
          <span class="skill-level ${skill.level.toLowerCase()}">${skill.level}</span>
        </div>
      `).join('')}
    </div>
  `;
}

function createAdditionalSection(profile) {
  return `
    <div class="section">
      <div class="section-header">
        <h3>🎯 Career & Additional Information</h3>
        <span class="section-toggle">▼</span>
      </div>
      <div class="section-content">
        <div class="info-grid">
          <div class="info-item">
            <div class="info-label">Target Company/Sector</div>
            <div class="info-value">${profile.company_aim || 'Not specified'}</div>
          </div>
          <div class="info-item">
            <div class="info-label">Target Package</div>
            <div class="info-value">${profile.target_package || 'Not specified'}</div>
          </div>
          <div class="info-item">
            <div class="info-label">Certifications</div>
            <div class="info-value">${profile.certifications || 'None listed'}</div>
          </div>
          <div class="info-item">
            <div class="info-label">Internships</div>
            <div class="info-value">${profile.internships || 'None listed'}</div>
          </div>
          <div class="info-item">
            <div class="info-label">Workshops Attended</div>
            <div class="info-value">${profile.workshops || 'None listed'}</div>
          </div>
          <div class="info-item">
            <div class="info-label">Awards & Achievements</div>
            <div class="info-value">${profile.awards || 'None listed'}</div>
          </div>
          <div class="info-item">
            <div class="info-label">Hackathons Participated</div>
            <div class="info-value">${profile.hackathons_attended || 'None listed'}</div>
          </div>
          <div class="info-item">
            <div class="info-label">Extracurricular Activities</div>
            <div class="info-value">${profile.extracurriculars || 'None listed'}</div>
          </div>
          <div class="info-item">
            <div class="info-label">Club Participation</div>
            <div class="info-value">${profile.club_participation || 'None listed'}</div>
          </div>
          <div class="info-item">
            <div class="info-label">Future Career Path</div>
            <div class="info-value">${profile.future_path || 'Not specified'}</div>
          </div>
          <div class="info-item">
            <div class="info-label">Strengths</div>
            <div class="info-value">${profile.strength || 'Not specified'}</div>
          </div>
          <div class="info-item">
            <div class="info-label">Areas for Improvement</div>
            <div class="info-value">${profile.weakness || 'Not specified'}</div>
          </div>
          <div class="info-item">
            <div class="info-label">Additional Remarks</div>
            <div class="info-value">${profile.remarks || 'None provided'}</div>
          </div>
        </div>
      </div>
    </div>
  `;
}

function formatDate(dateString) {
  if (!dateString) return null;
  try {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-IN', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  } catch (error) {
    return dateString; 
  }
}

function formatPercentage(value) {
  if (!value) return 'Not provided';
  return value.toString().includes('%') ? value : `${value}%`;
}

function formatLinkedIn(url) {
  if (!url) return 'Not provided';
  return `<a href="${url}" target="_blank" rel="noopener noreferrer">View Profile</a>`;
}
