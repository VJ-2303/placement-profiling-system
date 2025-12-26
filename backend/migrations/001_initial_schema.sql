-- Placement Profiling System Database Schema
-- Version: 1.0
-- Date: 2025-12-26

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================
-- ENUM TYPES
-- ============================================

CREATE TYPE placement_status AS ENUM ('not_placed', 'in_process', 'placed', 'higher_studies', 'entrepreneur');
CREATE TYPE proficiency_level AS ENUM ('beginner', 'intermediate', 'advanced', 'expert');
CREATE TYPE skill_category AS ENUM ('programming_language', 'database', 'framework', 'tool', 'concept', 'soft_skill');
CREATE TYPE residence_type AS ENUM ('day_scholar', 'hosteler');
CREATE TYPE gender AS ENUM ('male', 'female', 'other');

-- ============================================
-- CORE TABLES
-- ============================================

-- Batches table (e.g., 2024, 2025, 2026)
CREATE TABLE batches (
    id SERIAL PRIMARY KEY,
    year INTEGER NOT NULL UNIQUE,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Insert default batches
INSERT INTO batches (year) VALUES (2024), (2025), (2026), (2027);

-- Admins table (pre-registered placement coordinators)
CREATE TABLE admins (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(20),
    designation VARCHAR(100) DEFAULT 'Placement Coordinator',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Students table (core student info)
CREATE TABLE students (
    id SERIAL PRIMARY KEY,
    
    -- Basic Info (from OAuth)
    official_email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    
    -- Academic Identity
    roll_no VARCHAR(20) UNIQUE,
    register_no VARCHAR(20) UNIQUE,
    batch_id INTEGER REFERENCES batches(id),
    
    -- Profile
    photo_url TEXT,
    
    -- Status
    is_profile_completed BOOLEAN DEFAULT false,
    is_eligible_for_placement BOOLEAN DEFAULT true,
    placement_status placement_status DEFAULT 'not_placed',
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_login_at TIMESTAMP WITH TIME ZONE,
    
    -- Version for optimistic locking
    version INTEGER DEFAULT 1
);

-- Create index for faster searches
CREATE INDEX idx_students_roll_no ON students(roll_no);
CREATE INDEX idx_students_batch ON students(batch_id);
CREATE INDEX idx_students_placement_status ON students(placement_status);
CREATE INDEX idx_students_email ON students(official_email);

-- ============================================
-- STUDENT DETAILS TABLES
-- ============================================

-- Personal Details
CREATE TABLE student_personal_details (
    id SERIAL PRIMARY KEY,
    student_id INTEGER NOT NULL UNIQUE REFERENCES students(id) ON DELETE CASCADE,
    
    -- Personal Info
    date_of_birth DATE,
    gender gender,
    blood_group VARCHAR(5),
    
    -- Contact
    mobile_number VARCHAR(15),
    alternate_mobile VARCHAR(15),
    personal_email VARCHAR(255),
    linkedin_url TEXT,
    github_url TEXT,
    portfolio_url TEXT,
    
    -- Identity
    aadhaar_number VARCHAR(12),
    
    -- Address
    address TEXT,
    city VARCHAR(100),
    state VARCHAR(100) DEFAULT 'Tamil Nadu',
    pincode VARCHAR(10),
    residence_type residence_type,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Parent/Guardian Details
CREATE TABLE student_family_details (
    id SERIAL PRIMARY KEY,
    student_id INTEGER NOT NULL UNIQUE REFERENCES students(id) ON DELETE CASCADE,
    
    -- Father Details
    father_name VARCHAR(255),
    father_mobile VARCHAR(15),
    father_email VARCHAR(255),
    father_occupation VARCHAR(255),
    father_company VARCHAR(255),
    father_annual_income VARCHAR(50),
    
    -- Mother Details
    mother_name VARCHAR(255),
    mother_mobile VARCHAR(15),
    mother_email VARCHAR(255),
    mother_occupation VARCHAR(255),
    mother_company VARCHAR(255),
    
    -- Guardian (if different)
    guardian_name VARCHAR(255),
    guardian_mobile VARCHAR(15),
    guardian_relation VARCHAR(50),
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Academic Details
CREATE TABLE student_academics (
    id SERIAL PRIMARY KEY,
    student_id INTEGER NOT NULL UNIQUE REFERENCES students(id) ON DELETE CASCADE,
    
    -- School Education
    tenth_percentage DECIMAL(5,2),
    tenth_board VARCHAR(100),
    tenth_year INTEGER,
    tenth_school VARCHAR(255),
    
    twelfth_percentage DECIMAL(5,2),
    twelfth_board VARCHAR(100),
    twelfth_year INTEGER,
    twelfth_school VARCHAR(255),
    
    -- Diploma (if applicable)
    has_diploma BOOLEAN DEFAULT false,
    diploma_percentage DECIMAL(5,2),
    diploma_branch VARCHAR(100),
    diploma_college VARCHAR(255),
    
    -- College CGPA (Semester-wise)
    cgpa_sem1 DECIMAL(4,2),
    cgpa_sem2 DECIMAL(4,2),
    cgpa_sem3 DECIMAL(4,2),
    cgpa_sem4 DECIMAL(4,2),
    cgpa_sem5 DECIMAL(4,2),
    cgpa_sem6 DECIMAL(4,2),
    cgpa_sem7 DECIMAL(4,2),
    cgpa_sem8 DECIMAL(4,2),
    cgpa_overall DECIMAL(4,2),
    
    -- Backlogs
    current_backlogs INTEGER DEFAULT 0,
    history_of_backlogs BOOLEAN DEFAULT false,
    backlog_details TEXT,
    
    -- Gap Year
    has_gap_year BOOLEAN DEFAULT false,
    gap_year_reason TEXT,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- ============================================
-- SKILLS SYSTEM
-- ============================================

-- Master Skills Table
CREATE TABLE skills (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    category skill_category NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    display_order INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Insert default skills
INSERT INTO skills (name, category, display_order) VALUES
-- Programming Languages
('C', 'programming_language', 1),
('C++', 'programming_language', 2),
('Java', 'programming_language', 3),
('Python', 'programming_language', 4),
('JavaScript', 'programming_language', 5),
('TypeScript', 'programming_language', 6),
('Go', 'programming_language', 7),
('Rust', 'programming_language', 8),
('PHP', 'programming_language', 9),
('Ruby', 'programming_language', 10),
('Kotlin', 'programming_language', 11),
('Swift', 'programming_language', 12),

-- Databases
('MySQL', 'database', 1),
('PostgreSQL', 'database', 2),
('MongoDB', 'database', 3),
('Redis', 'database', 4),
('SQLite', 'database', 5),
('Oracle', 'database', 6),
('Firebase', 'database', 7),

-- Frameworks
('React', 'framework', 1),
('Angular', 'framework', 2),
('Vue.js', 'framework', 3),
('Node.js', 'framework', 4),
('Express.js', 'framework', 5),
('Django', 'framework', 6),
('Flask', 'framework', 7),
('Spring Boot', 'framework', 8),
('Flutter', 'framework', 9),
('React Native', 'framework', 10),
('.NET', 'framework', 11),
('Laravel', 'framework', 12),
('Next.js', 'framework', 13),
('FastAPI', 'framework', 14),
('Tailwind CSS', 'framework', 15),
('Bootstrap', 'framework', 16),

-- Tools
('Git/GitHub', 'tool', 1),
('Docker', 'tool', 2),
('Kubernetes', 'tool', 3),
('Linux', 'tool', 4),
('AWS', 'tool', 5),
('Azure', 'tool', 6),
('GCP', 'tool', 7),
('Jenkins', 'tool', 8),
('Jira', 'tool', 9),
('Postman', 'tool', 10),
('VS Code', 'tool', 11),
('Figma', 'tool', 12),

-- Concepts
('Data Structures', 'concept', 1),
('Algorithms', 'concept', 2),
('DBMS', 'concept', 3),
('Operating Systems', 'concept', 4),
('Computer Networks', 'concept', 5),
('OOP', 'concept', 6),
('System Design', 'concept', 7),
('Software Engineering', 'concept', 8),
('Machine Learning', 'concept', 9),
('Deep Learning', 'concept', 10),
('Cloud Computing', 'concept', 11),
('Cybersecurity', 'concept', 12),

-- Soft Skills
('Communication', 'soft_skill', 1),
('Problem Solving', 'soft_skill', 2),
('Team Work', 'soft_skill', 3),
('Leadership', 'soft_skill', 4),
('Time Management', 'soft_skill', 5),
('Critical Thinking', 'soft_skill', 6),
('Adaptability', 'soft_skill', 7),
('Presentation', 'soft_skill', 8);

-- Student Skills Junction Table
CREATE TABLE student_skills (
    id SERIAL PRIMARY KEY,
    student_id INTEGER NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    skill_id INTEGER NOT NULL REFERENCES skills(id) ON DELETE CASCADE,
    proficiency proficiency_level NOT NULL,
    years_of_experience DECIMAL(3,1) DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(student_id, skill_id)
);

CREATE INDEX idx_student_skills_student ON student_skills(student_id);
CREATE INDEX idx_student_skills_skill ON student_skills(skill_id);
CREATE INDEX idx_student_skills_proficiency ON student_skills(proficiency);

-- ============================================
-- ACHIEVEMENTS & ACTIVITIES
-- ============================================

CREATE TABLE student_achievements (
    id SERIAL PRIMARY KEY,
    student_id INTEGER NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    
    -- Certifications
    certifications TEXT,
    
    -- Awards & Recognition
    awards TEXT,
    
    -- Workshops & Training
    workshops TEXT,
    
    -- Internships
    internships TEXT,
    
    -- Projects
    projects TEXT,
    
    -- Competitive Programming
    leetcode_profile TEXT,
    hackerrank_profile TEXT,
    codeforces_profile TEXT,
    codechef_profile TEXT,
    leetcode_rating INTEGER,
    problems_solved INTEGER,
    
    -- Hackathons
    hackathons_participated INTEGER DEFAULT 0,
    hackathons_won INTEGER DEFAULT 0,
    hackathon_details TEXT,
    
    -- Extra-curricular
    extracurriculars TEXT,
    club_memberships TEXT,
    sports TEXT,
    volunteer_work TEXT,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- ============================================
-- CAREER ASPIRATIONS
-- ============================================

CREATE TABLE student_aspirations (
    id SERIAL PRIMARY KEY,
    student_id INTEGER NOT NULL UNIQUE REFERENCES students(id) ON DELETE CASCADE,
    
    -- Career Goals
    dream_companies TEXT,
    preferred_roles TEXT,
    preferred_locations TEXT,
    expected_salary VARCHAR(50),
    willing_to_relocate BOOLEAN DEFAULT true,
    
    -- Career Path
    career_objective TEXT,
    short_term_goals TEXT,
    long_term_goals TEXT,
    
    -- Self Assessment
    strengths TEXT,
    weaknesses TEXT,
    
    -- Additional
    hobbies TEXT,
    languages_known TEXT,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- ============================================
-- PLACEMENT SYSTEM
-- ============================================

-- Companies Master Table
CREATE TABLE companies (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    website TEXT,
    industry VARCHAR(100),
    company_type VARCHAR(50), -- 'product', 'service', 'startup', 'mnc'
    description TEXT,
    logo_url TEXT,
    
    -- Contact
    hr_name VARCHAR(255),
    hr_email VARCHAR(255),
    hr_phone VARCHAR(20),
    
    -- Location
    headquarters VARCHAR(255),
    locations TEXT, -- JSON array of locations
    
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Placement Records
CREATE TABLE placements (
    id SERIAL PRIMARY KEY,
    student_id INTEGER NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    company_id INTEGER REFERENCES companies(id),
    
    -- If company not in master
    company_name VARCHAR(255),
    
    -- Offer Details
    job_role VARCHAR(255),
    package_lpa DECIMAL(10,2),
    package_ctc VARCHAR(100),
    joining_date DATE,
    offer_date DATE,
    offer_type VARCHAR(50), -- 'full_time', 'internship', 'ppo'
    
    -- Location
    job_location VARCHAR(255),
    
    -- Status
    is_accepted BOOLEAN DEFAULT true,
    
    -- Verification
    verified_by INTEGER REFERENCES admins(id),
    verified_at TIMESTAMP WITH TIME ZONE,
    offer_letter_url TEXT,
    
    remarks TEXT,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_placements_student ON placements(student_id);
CREATE INDEX idx_placements_company ON placements(company_id);

-- ============================================
-- RESUME MANAGEMENT
-- ============================================

CREATE TABLE student_resumes (
    id SERIAL PRIMARY KEY,
    student_id INTEGER NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    file_name VARCHAR(255) NOT NULL,
    file_url TEXT NOT NULL,
    file_size INTEGER,
    is_primary BOOLEAN DEFAULT false,
    uploaded_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- ============================================
-- ACTIVITY LOG (for audit)
-- ============================================

CREATE TABLE activity_logs (
    id SERIAL PRIMARY KEY,
    user_type VARCHAR(20) NOT NULL, -- 'student' or 'admin'
    user_id INTEGER NOT NULL,
    action VARCHAR(100) NOT NULL,
    entity_type VARCHAR(50),
    entity_id INTEGER,
    details JSONB,
    ip_address INET,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_activity_logs_user ON activity_logs(user_type, user_id);
CREATE INDEX idx_activity_logs_created ON activity_logs(created_at);

-- ============================================
-- VIEWS FOR ANALYTICS
-- ============================================

-- View: Student Full Profile
CREATE OR REPLACE VIEW v_student_full_profile AS
SELECT 
    s.id,
    s.official_email,
    s.name,
    s.roll_no,
    s.register_no,
    s.photo_url,
    s.is_profile_completed,
    s.is_eligible_for_placement,
    s.placement_status,
    s.created_at,
    s.updated_at,
    
    b.year as batch_year,
    
    -- Personal Details
    spd.date_of_birth,
    spd.gender,
    spd.mobile_number,
    spd.personal_email,
    spd.linkedin_url,
    spd.github_url,
    spd.city,
    spd.residence_type,
    
    -- Academics
    sa.tenth_percentage,
    sa.twelfth_percentage,
    sa.cgpa_overall,
    sa.current_backlogs,
    
    -- Placement Info
    p.company_name as placed_company,
    p.package_lpa,
    p.job_role
    
FROM students s
LEFT JOIN batches b ON s.batch_id = b.id
LEFT JOIN student_personal_details spd ON s.id = spd.student_id
LEFT JOIN student_academics sa ON s.id = sa.student_id
LEFT JOIN placements p ON s.id = p.student_id AND p.is_accepted = true;

-- View: Placement Statistics
CREATE OR REPLACE VIEW v_placement_stats AS
SELECT 
    b.year as batch_year,
    COUNT(DISTINCT s.id) as total_students,
    COUNT(DISTINCT s.id) FILTER (WHERE s.is_profile_completed = true) as profiles_completed,
    COUNT(DISTINCT s.id) FILTER (WHERE s.placement_status = 'placed') as students_placed,
    COUNT(DISTINCT s.id) FILTER (WHERE s.placement_status = 'not_placed') as students_not_placed,
    COUNT(DISTINCT s.id) FILTER (WHERE s.placement_status = 'in_process') as students_in_process,
    COUNT(DISTINCT s.id) FILTER (WHERE s.placement_status = 'higher_studies') as higher_studies,
    ROUND(AVG(p.package_lpa) FILTER (WHERE p.is_accepted = true), 2) as avg_package,
    MAX(p.package_lpa) FILTER (WHERE p.is_accepted = true) as max_package,
    MIN(p.package_lpa) FILTER (WHERE p.is_accepted = true AND p.package_lpa > 0) as min_package
FROM students s
LEFT JOIN batches b ON s.batch_id = b.id
LEFT JOIN placements p ON s.id = p.student_id
GROUP BY b.year;

-- ============================================
-- FUNCTIONS
-- ============================================

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply trigger to all tables with updated_at
CREATE TRIGGER update_students_updated_at BEFORE UPDATE ON students
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_admins_updated_at BEFORE UPDATE ON admins
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_student_personal_updated_at BEFORE UPDATE ON student_personal_details
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_student_family_updated_at BEFORE UPDATE ON student_family_details
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_student_academics_updated_at BEFORE UPDATE ON student_academics
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_student_skills_updated_at BEFORE UPDATE ON student_skills
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_student_achievements_updated_at BEFORE UPDATE ON student_achievements
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_student_aspirations_updated_at BEFORE UPDATE ON student_aspirations
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_companies_updated_at BEFORE UPDATE ON companies
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_placements_updated_at BEFORE UPDATE ON placements
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================
-- INSERT DEFAULT ADMIN (Update with actual email)
-- ============================================

INSERT INTO admins (name, email, designation) VALUES
('Placement Coordinator', 'placement@kct.ac.in', 'Placement Officer'),
('HOD CSE', 'hod.cse@kct.ac.in', 'Head of Department');

-- ============================================
-- GRANT PERMISSIONS (if using separate user)
-- ============================================

-- GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO your_app_user;
-- GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO your_app_user;
