-- Main table for core student identity.
CREATE TABLE IF NOT EXISTS students (
    id BIGSERIAL PRIMARY KEY,
    roll_no TEXT UNIQUE,
    name TEXT,
    official_email TEXT UNIQUE NOT NULL,
    photo TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    version INT NOT NULL DEFAULT 1
);

-- Personal and contact details.
CREATE TABLE IF NOT EXISTS student_details (
    student_id BIGINT PRIMARY KEY REFERENCES students(id) ON DELETE CASCADE,
    date_of_birth DATE NOT NULL,
    mobile_number TEXT NOT NULL,
    alternate_mobile_number TEXT NOT NULL,
    personal_email TEXT UNIQUE NOT NULL,
    linkedin_profile TEXT NOT NULL,
    address TEXT NOT NULL,
    city TEXT NOT NULL,
    pincode TEXT, -- Made nullable as it can be null in JSON
    adhaar_no TEXT NOT NULL,
    residence_type TEXT NOT NULL,
    strength TEXT DEFAULT '', -- Made optional with default
    weakness TEXT DEFAULT '', -- Made optional with default
    remarks TEXT DEFAULT '' -- Made optional with default
);

-- Parent/guardian information.
CREATE TABLE IF NOT EXISTS student_parents (
    student_id BIGINT PRIMARY KEY REFERENCES students(id) ON DELETE CASCADE,
    father_name TEXT NOT NULL,
    father_mobile TEXT NOT NULL,
    father_occupation TEXT NOT NULL,
    father_company_details TEXT NOT NULL,
    father_email TEXT NOT NULL,
    mother_name TEXT NOT NULL,
    mother_mobile TEXT NOT NULL,
    mother_occupation TEXT NOT NULL,
    mother_email TEXT NOT NULL
);

-- All academic records.
CREATE TABLE IF NOT EXISTS student_academics (
    student_id BIGINT PRIMARY KEY REFERENCES students(id) ON DELETE CASCADE,
    tenth_percentage TEXT NOT NULL,
    twelth_percentage TEXT NOT NULL,
    cgpa_sem1 TEXT,
    cgpa_sem2 TEXT,
    cgpa_sem3 TEXT,
    cgpa_sem4 TEXT,
    cgpa_overall TEXT NOT NULL,
    current_backlogs TEXT NOT NULL,
    has_backlog_history TEXT NOT NULL
);

-- Career goals and extracurriculars.
CREATE TABLE IF NOT EXISTS student_aspirations (
    student_id BIGINT PRIMARY KEY REFERENCES students(id) ON DELETE CASCADE,
    company_aim TEXT NOT NULL,
    target_package TEXT NOT NULL,
    certifications TEXT NOT NULL,
    awards TEXT NOT NULL,
    workshops TEXT NOT NULL,
    internships TEXT NOT NULL,
    hackathons_attended TEXT DEFAULT '', -- Made optional with default
    extracurriculars TEXT DEFAULT '', -- Made optional with default
    club_participation TEXT DEFAULT '', -- Made optional with default
    future_path TEXT DEFAULT '', -- Made optional with default
    communication_skills TEXT NOT NULL
);

-- Master table for all possible skills (you will pre-populate this).
CREATE TABLE IF NOT EXISTS skills (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    category TEXT NOT NULL -- 'Programming', 'Core Concept', 'Tool'
);

-- Links students to skills and their proficiency.
CREATE TABLE IF NOT EXISTS student_skills (
    student_id BIGINT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    skill_id INT NOT NULL REFERENCES skills(id) ON DELETE CASCADE,
    proficiency_level TEXT NOT NULL,
    PRIMARY KEY (student_id, skill_id)
);
