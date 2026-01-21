package models

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"
)

// ============================================
// ENUMS
// ============================================

type PlacementStatus string

const (
	PlacementStatusNotPlaced     PlacementStatus = "not_placed"
	PlacementStatusInProcess     PlacementStatus = "in_process"
	PlacementStatusPlaced        PlacementStatus = "placed"
	PlacementStatusHigherStudies PlacementStatus = "higher_studies"
	PlacementStatusEntrepreneur  PlacementStatus = "entrepreneur"
)

type ProficiencyLevel string

const (
	ProficiencyBeginner     ProficiencyLevel = "beginner"
	ProficiencyIntermediate ProficiencyLevel = "intermediate"
	ProficiencyAdvanced     ProficiencyLevel = "advanced"
	ProficiencyExpert       ProficiencyLevel = "expert"
)

type SkillCategory string

const (
	SkillCategoryProgramming SkillCategory = "programming_language"
	SkillCategoryDatabase    SkillCategory = "database"
	SkillCategoryFramework   SkillCategory = "framework"
	SkillCategoryTool        SkillCategory = "tool"
	SkillCategoryConcept     SkillCategory = "concept"
	SkillCategorySoftSkill   SkillCategory = "soft_skill"
)

// ============================================
// STUDENT STRUCTS
// ============================================

type Student struct {
	ID                     int64           `json:"id"`
	OfficialEmail          string          `json:"official_email"`
	Name                   string          `json:"name"`
	RollNo                 *string         `json:"roll_no"`
	RegisterNo             *string         `json:"register_no"`
	BatchID                *int            `json:"batch_id"`
	BatchYear              *int            `json:"batch_year,omitempty"`
	PhotoURL               *string         `json:"photo_url"`
	IsProfileCompleted     bool            `json:"is_profile_completed"`
	IsEligibleForPlacement bool            `json:"is_eligible_for_placement"`
	PlacementStatus        PlacementStatus `json:"placement_status"`
	CreatedAt              time.Time       `json:"created_at"`
	UpdatedAt              time.Time       `json:"updated_at"`
	LastLoginAt            *time.Time      `json:"last_login_at,omitempty"`
	Version                int             `json:"version"`
}

type StudentPersonalDetails struct {
	ID              int64   `json:"id"`
	StudentID       int64   `json:"student_id"`
	DateOfBirth     *string `json:"date_of_birth"`
	Gender          *string `json:"gender"`
	BloodGroup      *string `json:"blood_group"`
	MobileNumber    *string `json:"mobile_number"`
	AlternateMobile *string `json:"alternate_mobile"`
	PersonalEmail   *string `json:"personal_email"`
	LinkedinURL     *string `json:"linkedin_url"`
	GithubURL       *string `json:"github_url"`
	PortfolioURL    *string `json:"portfolio_url"`
	AadhaarNumber   *string `json:"aadhaar_number"`
	Address         *string `json:"address"`
	City            *string `json:"city"`
	State           *string `json:"state"`
	Pincode         *string `json:"pincode"`
	ResidenceType   *string `json:"residence_type"`
}

type StudentFamilyDetails struct {
	ID                 int64   `json:"id"`
	StudentID          int64   `json:"student_id"`
	FatherName         *string `json:"father_name"`
	FatherMobile       *string `json:"father_mobile"`
	FatherEmail        *string `json:"father_email"`
	FatherOccupation   *string `json:"father_occupation"`
	FatherCompany      *string `json:"father_company"`
	FatherAnnualIncome *string `json:"father_annual_income"`
	MotherName         *string `json:"mother_name"`
	MotherMobile       *string `json:"mother_mobile"`
	MotherEmail        *string `json:"mother_email"`
	MotherOccupation   *string `json:"mother_occupation"`
	MotherCompany      *string `json:"mother_company"`
	GuardianName       *string `json:"guardian_name"`
	GuardianMobile     *string `json:"guardian_mobile"`
	GuardianRelation   *string `json:"guardian_relation"`
}

type StudentAcademics struct {
	ID                int64    `json:"id"`
	StudentID         int64    `json:"student_id"`
	TenthPercentage   *float64 `json:"tenth_percentage"`
	TenthBoard        *string  `json:"tenth_board"`
	TenthYear         *int     `json:"tenth_year"`
	TenthSchool       *string  `json:"tenth_school"`
	TwelfthPercentage *float64 `json:"twelfth_percentage"`
	TwelfthBoard      *string  `json:"twelfth_board"`
	TwelfthYear       *int     `json:"twelfth_year"`
	TwelfthSchool     *string  `json:"twelfth_school"`
	HasDiploma        bool     `json:"has_diploma"`
	DiplomaPercentage *float64 `json:"diploma_percentage"`
	DiplomaBranch     *string  `json:"diploma_branch"`
	DiplomaCollege    *string  `json:"diploma_college"`
	CGPASem1          *float64 `json:"cgpa_sem1"`
	CGPASem2          *float64 `json:"cgpa_sem2"`
	CGPASem3          *float64 `json:"cgpa_sem3"`
	CGPASem4          *float64 `json:"cgpa_sem4"`
	CGPASem5          *float64 `json:"cgpa_sem5"`
	CGPASem6          *float64 `json:"cgpa_sem6"`
	CGPASem7          *float64 `json:"cgpa_sem7"`
	CGPASem8          *float64 `json:"cgpa_sem8"`
	CGPAOverall       *float64 `json:"cgpa_overall"`
	CurrentBacklogs   int      `json:"current_backlogs"`
	HistoryOfBacklogs bool     `json:"history_of_backlogs"`
	BacklogDetails    *string  `json:"backlog_details"`
	HasGapYear        bool     `json:"has_gap_year"`
	GapYearReason     *string  `json:"gap_year_reason"`
}

type StudentAchievements struct {
	ID                     int64   `json:"id"`
	StudentID              int64   `json:"student_id"`
	Certifications         *string `json:"certifications"`
	Awards                 *string `json:"awards"`
	Workshops              *string `json:"workshops"`
	Internships            *string `json:"internships"`
	Projects               *string `json:"projects"`
	LeetcodeProfile        *string `json:"leetcode_profile"`
	HackerrankProfile      *string `json:"hackerrank_profile"`
	CodeforcesProfile      *string `json:"codeforces_profile"`
	CodechefProfile        *string `json:"codechef_profile"`
	LeetcodeRating         *int    `json:"leetcode_rating"`
	ProblemsSolved         *int    `json:"problems_solved"`
	HackathonsParticipated int     `json:"hackathons_participated"`
	HackathonsWon          int     `json:"hackathons_won"`
	HackathonDetails       *string `json:"hackathon_details"`
	Extracurriculars       *string `json:"extracurriculars"`
	ClubMemberships        *string `json:"club_memberships"`
	Sports                 *string `json:"sports"`
	VolunteerWork          *string `json:"volunteer_work"`
}

type StudentAspirations struct {
	ID                 int64   `json:"id"`
	StudentID          int64   `json:"student_id"`
	DreamCompanies     *string `json:"dream_companies"`
	PreferredRoles     *string `json:"preferred_roles"`
	PreferredLocations *string `json:"preferred_locations"`
	ExpectedSalary     *string `json:"expected_salary"`
	WillingToRelocate  bool    `json:"willing_to_relocate"`
	CareerObjective    *string `json:"career_objective"`
	ShortTermGoals     *string `json:"short_term_goals"`
	LongTermGoals      *string `json:"long_term_goals"`
	Strengths          *string `json:"strengths"`
	Weaknesses         *string `json:"weaknesses"`
	Hobbies            *string `json:"hobbies"`
	LanguagesKnown     *string `json:"languages_known"`
}

type StudentSkill struct {
	ID                int64            `json:"id"`
	StudentID         int64            `json:"student_id"`
	SkillID           int              `json:"skill_id"`
	SkillName         string           `json:"skill_name,omitempty"`
	SkillCategory     SkillCategory    `json:"skill_category,omitempty"`
	Proficiency       ProficiencyLevel `json:"proficiency"`
	YearsOfExperience float64          `json:"years_of_experience"`
}

// Full student profile combining all details
type StudentFullProfile struct {
	Student      Student                 `json:"student"`
	Personal     *StudentPersonalDetails `json:"personal"`
	Family       *StudentFamilyDetails   `json:"family"`
	Academics    *StudentAcademics       `json:"academics"`
	Achievements *StudentAchievements    `json:"achievements"`
	Aspirations  *StudentAspirations     `json:"aspirations"`
	Skills       []StudentSkill          `json:"skills"`
	Placement    *PlacementRecord        `json:"placement,omitempty"`
}

// ============================================
// STUDENT MODEL
// ============================================

type StudentModel struct {
	DB *sql.DB
}

// Insert creates a new student (from OAuth)
func (m StudentModel) Insert(student *Student) error {
	query := `
		INSERT INTO students (official_email, name)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at, version`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, student.OfficialEmail, student.Name).Scan(
		&student.ID,
		&student.CreatedAt,
		&student.UpdatedAt,
		&student.Version,
	)
}

// GetByEmail retrieves a student by email
func (m StudentModel) GetByEmail(email string) (*Student, error) {
	query := `
		SELECT s.id, s.official_email, s.name, s.roll_no, s.register_no, 
		       s.batch_id, b.year, s.photo_url, s.is_profile_completed, 
		       s.is_eligible_for_placement, s.placement_status,
		       s.created_at, s.updated_at, s.last_login_at, s.version
		FROM students s
		LEFT JOIN batches b ON s.batch_id = b.id
		WHERE s.official_email = $1`

	var student Student
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&student.ID, &student.OfficialEmail, &student.Name, &student.RollNo,
		&student.RegisterNo, &student.BatchID, &student.BatchYear, &student.PhotoURL,
		&student.IsProfileCompleted, &student.IsEligibleForPlacement,
		&student.PlacementStatus, &student.CreatedAt, &student.UpdatedAt,
		&student.LastLoginAt, &student.Version,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &student, nil
}

// GetByID retrieves a student by ID
func (m StudentModel) GetByID(id int64) (*Student, error) {
	query := `
		SELECT s.id, s.official_email, s.name, s.roll_no, s.register_no, 
		       s.batch_id, b.year, s.photo_url, s.is_profile_completed, 
		       s.is_eligible_for_placement, s.placement_status,
		       s.created_at, s.updated_at, s.last_login_at, s.version
		FROM students s
		LEFT JOIN batches b ON s.batch_id = b.id
		WHERE s.id = $1`

	var student Student
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&student.ID, &student.OfficialEmail, &student.Name, &student.RollNo,
		&student.RegisterNo, &student.BatchID, &student.BatchYear, &student.PhotoURL,
		&student.IsProfileCompleted, &student.IsEligibleForPlacement,
		&student.PlacementStatus, &student.CreatedAt, &student.UpdatedAt,
		&student.LastLoginAt, &student.Version,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &student, nil
}

// UpdateBasicInfo updates student's basic info
func (m StudentModel) UpdateBasicInfo(student *Student) error {
	query := `
		UPDATE students 
		SET name = $1, roll_no = $2, register_no = $3, batch_id = $4, 
		    photo_url = $5, version = version + 1
		WHERE id = $6 AND version = $7
		RETURNING version, updated_at`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query,
		student.Name, student.RollNo, student.RegisterNo, student.BatchID,
		student.PhotoURL, student.ID, student.Version,
	).Scan(&student.Version, &student.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrEditConflict
		}
		return err
	}

	return nil
}

// SetProfileCompleted marks profile as completed
func (m StudentModel) SetProfileCompleted(tx *sql.Tx, studentID int64) error {
	query := `UPDATE students SET is_profile_completed = true WHERE id = $1`
	_, err := tx.Exec(query, studentID)
	return err
}

// UpdateLastLogin updates last login timestamp
func (m StudentModel) UpdateLastLogin(id int64) error {
	query := `UPDATE students SET last_login_at = NOW() WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

// GetBatchIDByYear retrieves batch ID by year
func (m StudentModel) GetBatchIDByYear(year int) (int, error) {
	query := `SELECT id FROM batches WHERE year = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var id int
	err := m.DB.QueryRowContext(ctx, query, year).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrRecordNotFound
		}
		return 0, err
	}
	return id, nil
}

// ============================================
// PERSONAL DETAILS
// ============================================

func (m StudentModel) UpsertPersonalDetails(tx *sql.Tx, details *StudentPersonalDetails) error {
	query := `
		INSERT INTO student_personal_details (
			student_id, date_of_birth, gender, blood_group, mobile_number,
			alternate_mobile, personal_email, linkedin_url, github_url, portfolio_url,
			aadhaar_number, address, city, state, pincode, residence_type
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		ON CONFLICT (student_id) DO UPDATE SET
			date_of_birth = EXCLUDED.date_of_birth,
			gender = EXCLUDED.gender,
			blood_group = EXCLUDED.blood_group,
			mobile_number = EXCLUDED.mobile_number,
			alternate_mobile = EXCLUDED.alternate_mobile,
			personal_email = EXCLUDED.personal_email,
			linkedin_url = EXCLUDED.linkedin_url,
			github_url = EXCLUDED.github_url,
			portfolio_url = EXCLUDED.portfolio_url,
			aadhaar_number = EXCLUDED.aadhaar_number,
			address = EXCLUDED.address,
			city = EXCLUDED.city,
			state = EXCLUDED.state,
			pincode = EXCLUDED.pincode,
			residence_type = EXCLUDED.residence_type`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := tx.ExecContext(ctx, query,
		details.StudentID, details.DateOfBirth, details.Gender, details.BloodGroup,
		details.MobileNumber, details.AlternateMobile, details.PersonalEmail,
		details.LinkedinURL, details.GithubURL, details.PortfolioURL,
		details.AadhaarNumber, details.Address, details.City, details.State,
		details.Pincode, details.ResidenceType,
	)
	return err
}

func (m StudentModel) GetPersonalDetails(studentID int64) (*StudentPersonalDetails, error) {
	query := `
		SELECT id, student_id, date_of_birth::text, gender, blood_group, mobile_number,
		       alternate_mobile, personal_email, linkedin_url, github_url, portfolio_url,
		       aadhaar_number, address, city, state, pincode, residence_type
		FROM student_personal_details
		WHERE student_id = $1`

	var d StudentPersonalDetails
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, studentID).Scan(
		&d.ID, &d.StudentID, &d.DateOfBirth, &d.Gender, &d.BloodGroup,
		&d.MobileNumber, &d.AlternateMobile, &d.PersonalEmail, &d.LinkedinURL,
		&d.GithubURL, &d.PortfolioURL, &d.AadhaarNumber, &d.Address, &d.City,
		&d.State, &d.Pincode, &d.ResidenceType,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Not an error, just no data yet
		}
		return nil, err
	}

	return &d, nil
}

// ============================================
// FAMILY DETAILS
// ============================================

func (m StudentModel) UpsertFamilyDetails(tx *sql.Tx, details *StudentFamilyDetails) error {
	query := `
		INSERT INTO student_family_details (
			student_id, father_name, father_mobile, father_email, father_occupation,
			father_company, father_annual_income, mother_name, mother_mobile, mother_email,
			mother_occupation, mother_company, guardian_name, guardian_mobile, guardian_relation
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		ON CONFLICT (student_id) DO UPDATE SET
			father_name = EXCLUDED.father_name,
			father_mobile = EXCLUDED.father_mobile,
			father_email = EXCLUDED.father_email,
			father_occupation = EXCLUDED.father_occupation,
			father_company = EXCLUDED.father_company,
			father_annual_income = EXCLUDED.father_annual_income,
			mother_name = EXCLUDED.mother_name,
			mother_mobile = EXCLUDED.mother_mobile,
			mother_email = EXCLUDED.mother_email,
			mother_occupation = EXCLUDED.mother_occupation,
			mother_company = EXCLUDED.mother_company,
			guardian_name = EXCLUDED.guardian_name,
			guardian_mobile = EXCLUDED.guardian_mobile,
			guardian_relation = EXCLUDED.guardian_relation`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := tx.ExecContext(ctx, query,
		details.StudentID, details.FatherName, details.FatherMobile, details.FatherEmail,
		details.FatherOccupation, details.FatherCompany, details.FatherAnnualIncome,
		details.MotherName, details.MotherMobile, details.MotherEmail,
		details.MotherOccupation, details.MotherCompany, details.GuardianName,
		details.GuardianMobile, details.GuardianRelation,
	)
	return err
}

func (m StudentModel) GetFamilyDetails(studentID int64) (*StudentFamilyDetails, error) {
	query := `
		SELECT id, student_id, father_name, father_mobile, father_email, father_occupation,
		       father_company, father_annual_income, mother_name, mother_mobile, mother_email,
		       mother_occupation, mother_company, guardian_name, guardian_mobile, guardian_relation
		FROM student_family_details
		WHERE student_id = $1`

	var d StudentFamilyDetails
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, studentID).Scan(
		&d.ID, &d.StudentID, &d.FatherName, &d.FatherMobile, &d.FatherEmail,
		&d.FatherOccupation, &d.FatherCompany, &d.FatherAnnualIncome,
		&d.MotherName, &d.MotherMobile, &d.MotherEmail, &d.MotherOccupation,
		&d.MotherCompany, &d.GuardianName, &d.GuardianMobile, &d.GuardianRelation,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &d, nil
}

// ============================================
// ACADEMICS
// ============================================

func (m StudentModel) UpsertAcademics(tx *sql.Tx, a *StudentAcademics) error {
	query := `
		INSERT INTO student_academics (
			student_id, tenth_percentage, tenth_board, tenth_year, tenth_school,
			twelfth_percentage, twelfth_board, twelfth_year, twelfth_school,
			has_diploma, diploma_percentage, diploma_branch, diploma_college,
			cgpa_sem1, cgpa_sem2, cgpa_sem3, cgpa_sem4, cgpa_sem5, cgpa_sem6, cgpa_sem7, cgpa_sem8,
			cgpa_overall, current_backlogs, history_of_backlogs, backlog_details,
			has_gap_year, gap_year_reason
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27)
		ON CONFLICT (student_id) DO UPDATE SET
			tenth_percentage = EXCLUDED.tenth_percentage,
			tenth_board = EXCLUDED.tenth_board,
			tenth_year = EXCLUDED.tenth_year,
			tenth_school = EXCLUDED.tenth_school,
			twelfth_percentage = EXCLUDED.twelfth_percentage,
			twelfth_board = EXCLUDED.twelfth_board,
			twelfth_year = EXCLUDED.twelfth_year,
			twelfth_school = EXCLUDED.twelfth_school,
			has_diploma = EXCLUDED.has_diploma,
			diploma_percentage = EXCLUDED.diploma_percentage,
			diploma_branch = EXCLUDED.diploma_branch,
			diploma_college = EXCLUDED.diploma_college,
			cgpa_sem1 = EXCLUDED.cgpa_sem1,
			cgpa_sem2 = EXCLUDED.cgpa_sem2,
			cgpa_sem3 = EXCLUDED.cgpa_sem3,
			cgpa_sem4 = EXCLUDED.cgpa_sem4,
			cgpa_sem5 = EXCLUDED.cgpa_sem5,
			cgpa_sem6 = EXCLUDED.cgpa_sem6,
			cgpa_sem7 = EXCLUDED.cgpa_sem7,
			cgpa_sem8 = EXCLUDED.cgpa_sem8,
			cgpa_overall = EXCLUDED.cgpa_overall,
			current_backlogs = EXCLUDED.current_backlogs,
			history_of_backlogs = EXCLUDED.history_of_backlogs,
			backlog_details = EXCLUDED.backlog_details,
			has_gap_year = EXCLUDED.has_gap_year,
			gap_year_reason = EXCLUDED.gap_year_reason`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := tx.ExecContext(ctx, query,
		a.StudentID, a.TenthPercentage, a.TenthBoard, a.TenthYear, a.TenthSchool,
		a.TwelfthPercentage, a.TwelfthBoard, a.TwelfthYear, a.TwelfthSchool,
		a.HasDiploma, a.DiplomaPercentage, a.DiplomaBranch, a.DiplomaCollege,
		a.CGPASem1, a.CGPASem2, a.CGPASem3, a.CGPASem4, a.CGPASem5, a.CGPASem6, a.CGPASem7, a.CGPASem8,
		a.CGPAOverall, a.CurrentBacklogs, a.HistoryOfBacklogs, a.BacklogDetails,
		a.HasGapYear, a.GapYearReason,
	)
	return err
}

func (m StudentModel) GetAcademics(studentID int64) (*StudentAcademics, error) {
	query := `
		SELECT id, student_id, tenth_percentage, tenth_board, tenth_year, tenth_school,
		       twelfth_percentage, twelfth_board, twelfth_year, twelfth_school,
		       has_diploma, diploma_percentage, diploma_branch, diploma_college,
		       cgpa_sem1, cgpa_sem2, cgpa_sem3, cgpa_sem4, cgpa_sem5, cgpa_sem6, cgpa_sem7, cgpa_sem8,
		       cgpa_overall, current_backlogs, history_of_backlogs, backlog_details,
		       has_gap_year, gap_year_reason
		FROM student_academics
		WHERE student_id = $1`

	var a StudentAcademics
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, studentID).Scan(
		&a.ID, &a.StudentID, &a.TenthPercentage, &a.TenthBoard, &a.TenthYear, &a.TenthSchool,
		&a.TwelfthPercentage, &a.TwelfthBoard, &a.TwelfthYear, &a.TwelfthSchool,
		&a.HasDiploma, &a.DiplomaPercentage, &a.DiplomaBranch, &a.DiplomaCollege,
		&a.CGPASem1, &a.CGPASem2, &a.CGPASem3, &a.CGPASem4, &a.CGPASem5, &a.CGPASem6, &a.CGPASem7, &a.CGPASem8,
		&a.CGPAOverall, &a.CurrentBacklogs, &a.HistoryOfBacklogs, &a.BacklogDetails,
		&a.HasGapYear, &a.GapYearReason,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &a, nil
}

// ============================================
// ACHIEVEMENTS
// ============================================

func (m StudentModel) UpsertAchievements(tx *sql.Tx, a *StudentAchievements) error {
	query := `
		INSERT INTO student_achievements (
			student_id, certifications, awards, workshops, internships, projects,
			leetcode_profile, hackerrank_profile, codeforces_profile, codechef_profile,
			leetcode_rating, problems_solved, hackathons_participated, hackathons_won,
			hackathon_details, extracurriculars, club_memberships, sports, volunteer_work
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
		ON CONFLICT (student_id) DO UPDATE SET
			certifications = EXCLUDED.certifications,
			awards = EXCLUDED.awards,
			workshops = EXCLUDED.workshops,
			internships = EXCLUDED.internships,
			projects = EXCLUDED.projects,
			leetcode_profile = EXCLUDED.leetcode_profile,
			hackerrank_profile = EXCLUDED.hackerrank_profile,
			codeforces_profile = EXCLUDED.codeforces_profile,
			codechef_profile = EXCLUDED.codechef_profile,
			leetcode_rating = EXCLUDED.leetcode_rating,
			problems_solved = EXCLUDED.problems_solved,
			hackathons_participated = EXCLUDED.hackathons_participated,
			hackathons_won = EXCLUDED.hackathons_won,
			hackathon_details = EXCLUDED.hackathon_details,
			extracurriculars = EXCLUDED.extracurriculars,
			club_memberships = EXCLUDED.club_memberships,
			sports = EXCLUDED.sports,
			volunteer_work = EXCLUDED.volunteer_work`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := tx.ExecContext(ctx, query,
		a.StudentID, a.Certifications, a.Awards, a.Workshops, a.Internships, a.Projects,
		a.LeetcodeProfile, a.HackerrankProfile, a.CodeforcesProfile, a.CodechefProfile,
		a.LeetcodeRating, a.ProblemsSolved, a.HackathonsParticipated, a.HackathonsWon,
		a.HackathonDetails, a.Extracurriculars, a.ClubMemberships, a.Sports, a.VolunteerWork,
	)
	return err
}

func (m StudentModel) GetAchievements(studentID int64) (*StudentAchievements, error) {
	query := `
		SELECT id, student_id, certifications, awards, workshops, internships, projects,
		       leetcode_profile, hackerrank_profile, codeforces_profile, codechef_profile,
		       leetcode_rating, problems_solved, hackathons_participated, hackathons_won,
		       hackathon_details, extracurriculars, club_memberships, sports, volunteer_work
		FROM student_achievements
		WHERE student_id = $1`

	var a StudentAchievements
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, studentID).Scan(
		&a.ID, &a.StudentID, &a.Certifications, &a.Awards, &a.Workshops, &a.Internships, &a.Projects,
		&a.LeetcodeProfile, &a.HackerrankProfile, &a.CodeforcesProfile, &a.CodechefProfile,
		&a.LeetcodeRating, &a.ProblemsSolved, &a.HackathonsParticipated, &a.HackathonsWon,
		&a.HackathonDetails, &a.Extracurriculars, &a.ClubMemberships, &a.Sports, &a.VolunteerWork,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &a, nil
}

// ============================================
// ASPIRATIONS
// ============================================

func (m StudentModel) UpsertAspirations(tx *sql.Tx, a *StudentAspirations) error {
	query := `
		INSERT INTO student_aspirations (
			student_id, dream_companies, preferred_roles, preferred_locations, expected_salary,
			willing_to_relocate, career_objective, short_term_goals, long_term_goals,
			strengths, weaknesses, hobbies, languages_known
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		ON CONFLICT (student_id) DO UPDATE SET
			dream_companies = EXCLUDED.dream_companies,
			preferred_roles = EXCLUDED.preferred_roles,
			preferred_locations = EXCLUDED.preferred_locations,
			expected_salary = EXCLUDED.expected_salary,
			willing_to_relocate = EXCLUDED.willing_to_relocate,
			career_objective = EXCLUDED.career_objective,
			short_term_goals = EXCLUDED.short_term_goals,
			long_term_goals = EXCLUDED.long_term_goals,
			strengths = EXCLUDED.strengths,
			weaknesses = EXCLUDED.weaknesses,
			hobbies = EXCLUDED.hobbies,
			languages_known = EXCLUDED.languages_known`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := tx.ExecContext(ctx, query,
		a.StudentID, a.DreamCompanies, a.PreferredRoles, a.PreferredLocations, a.ExpectedSalary,
		a.WillingToRelocate, a.CareerObjective, a.ShortTermGoals, a.LongTermGoals,
		a.Strengths, a.Weaknesses, a.Hobbies, a.LanguagesKnown,
	)
	return err
}

func (m StudentModel) GetAspirations(studentID int64) (*StudentAspirations, error) {
	query := `
		SELECT id, student_id, dream_companies, preferred_roles, preferred_locations, expected_salary,
		       willing_to_relocate, career_objective, short_term_goals, long_term_goals,
		       strengths, weaknesses, hobbies, languages_known
		FROM student_aspirations
		WHERE student_id = $1`

	var a StudentAspirations
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, studentID).Scan(
		&a.ID, &a.StudentID, &a.DreamCompanies, &a.PreferredRoles, &a.PreferredLocations, &a.ExpectedSalary,
		&a.WillingToRelocate, &a.CareerObjective, &a.ShortTermGoals, &a.LongTermGoals,
		&a.Strengths, &a.Weaknesses, &a.Hobbies, &a.LanguagesKnown,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &a, nil
}

// ============================================
// STUDENT SKILLS
// ============================================

func (m StudentModel) UpsertSkills(tx *sql.Tx, studentID int64, skills []StudentSkill) error {
	// First delete existing skills
	_, err := tx.Exec("DELETE FROM student_skills WHERE student_id = $1", studentID)
	if err != nil {
		return err
	}

	// Insert new skills
	for _, skill := range skills {
		query := `
			INSERT INTO student_skills (student_id, skill_id, proficiency, years_of_experience)
			VALUES ($1, $2, $3, $4)`

		_, err := tx.Exec(query, studentID, skill.SkillID, skill.Proficiency, skill.YearsOfExperience)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m StudentModel) GetSkills(studentID int64) ([]StudentSkill, error) {
	query := `
		SELECT ss.id, ss.student_id, ss.skill_id, s.name, s.category, ss.proficiency, ss.years_of_experience
		FROM student_skills ss
		JOIN skills s ON ss.skill_id = s.id
		WHERE ss.student_id = $1
		ORDER BY s.category, s.display_order`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []StudentSkill
	for rows.Next() {
		var s StudentSkill
		if err := rows.Scan(&s.ID, &s.StudentID, &s.SkillID, &s.SkillName, &s.SkillCategory, &s.Proficiency, &s.YearsOfExperience); err != nil {
			return nil, err
		}
		skills = append(skills, s)
	}

	return skills, rows.Err()
}

// ============================================
// FULL PROFILE
// ============================================

func (m StudentModel) GetFullProfile(studentID int64) (*StudentFullProfile, error) {
	student, err := m.GetByID(studentID)
	if err != nil {
		return nil, err
	}

	profile := &StudentFullProfile{Student: *student}

	profile.Personal, _ = m.GetPersonalDetails(studentID)
	profile.Family, _ = m.GetFamilyDetails(studentID)
	profile.Academics, _ = m.GetAcademics(studentID)
	profile.Achievements, _ = m.GetAchievements(studentID)
	profile.Aspirations, _ = m.GetAspirations(studentID)
	profile.Skills, _ = m.GetSkills(studentID)

	return profile, nil
}

// ============================================
// SEARCH & FILTER (for Admin)
// ============================================

type StudentFilter struct {
	Search          string
	BatchYear       *int
	PlacementStatus *PlacementStatus
	MinCGPA         *float64
	MaxCGPA         *float64
	HasBacklogs     *bool
	SkillIDs        []int
	Page            int
	PageSize        int
}

type StudentListItem struct {
	ID                 int64           `json:"id"`
	Name               string          `json:"name"`
	OfficialEmail      string          `json:"official_email"`
	RollNo             *string         `json:"roll_no"`
	BatchYear          *int            `json:"batch_year"`
	PhotoURL           *string         `json:"photo_url"`
	IsProfileCompleted bool            `json:"is_profile_completed"`
	PlacementStatus    PlacementStatus `json:"placement_status"`
	CGPAOverall        *float64        `json:"cgpa_overall"`
	MobileNumber       *string         `json:"mobile_number"`
	PlacedCompany      *string         `json:"placed_company,omitempty"`
	PackageLPA         *float64        `json:"package_lpa,omitempty"`
}

type StudentListResult struct {
	Students   []StudentListItem `json:"students"`
	Total      int               `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalPages int               `json:"total_pages"`
}

func (m StudentModel) List(filter StudentFilter) (*StudentListResult, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 || filter.PageSize > 100 {
		filter.PageSize = 20
	}

	var conditions []string
	var args []interface{}
	argNum := 1

	// Build WHERE clause
	if filter.Search != "" {
		conditions = append(conditions, "(s.name ILIKE $"+string(rune('0'+argNum))+" OR s.roll_no ILIKE $"+string(rune('0'+argNum))+" OR s.official_email ILIKE $"+string(rune('0'+argNum))+")")
		args = append(args, "%"+filter.Search+"%")
		argNum++
	}

	if filter.BatchYear != nil {
		conditions = append(conditions, "b.year = $"+string(rune('0'+argNum)))
		args = append(args, *filter.BatchYear)
		argNum++
	}

	if filter.PlacementStatus != nil {
		conditions = append(conditions, "s.placement_status = $"+string(rune('0'+argNum)))
		args = append(args, *filter.PlacementStatus)
		argNum++
	}

	if filter.MinCGPA != nil {
		conditions = append(conditions, "sa.cgpa_overall >= $"+string(rune('0'+argNum)))
		args = append(args, *filter.MinCGPA)
		argNum++
	}

	if filter.MaxCGPA != nil {
		conditions = append(conditions, "sa.cgpa_overall <= $"+string(rune('0'+argNum)))
		args = append(args, *filter.MaxCGPA)
		argNum++
	}

	if filter.HasBacklogs != nil {
		if *filter.HasBacklogs {
			conditions = append(conditions, "sa.current_backlogs > 0")
		} else {
			conditions = append(conditions, "(sa.current_backlogs = 0 OR sa.current_backlogs IS NULL)")
		}
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total
	countQuery := `
		SELECT COUNT(DISTINCT s.id)
		FROM students s
		LEFT JOIN batches b ON s.batch_id = b.id
		LEFT JOIN student_academics sa ON s.id = sa.student_id
		` + whereClause

	var total int
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, err
	}

	// Get paginated results
	offset := (filter.Page - 1) * filter.PageSize
	args = append(args, filter.PageSize, offset)

	query := `
		SELECT DISTINCT s.id, s.name, s.official_email, s.roll_no, b.year, s.photo_url,
		       s.is_profile_completed, s.placement_status, sa.cgpa_overall, spd.mobile_number,
		       p.company_name, p.package_lpa
		FROM students s
		LEFT JOIN batches b ON s.batch_id = b.id
		LEFT JOIN student_academics sa ON s.id = sa.student_id
		LEFT JOIN student_personal_details spd ON s.id = spd.student_id
		LEFT JOIN placements p ON s.id = p.student_id AND p.is_accepted = true
		` + whereClause + `
		ORDER BY s.name ASC
		LIMIT $` + string(rune('0'+argNum)) + ` OFFSET $` + string(rune('0'+argNum+1))

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []StudentListItem
	for rows.Next() {
		var s StudentListItem
		if err := rows.Scan(
			&s.ID, &s.Name, &s.OfficialEmail, &s.RollNo, &s.BatchYear, &s.PhotoURL,
			&s.IsProfileCompleted, &s.PlacementStatus, &s.CGPAOverall, &s.MobileNumber,
			&s.PlacedCompany, &s.PackageLPA,
		); err != nil {
			return nil, err
		}
		students = append(students, s)
	}

	totalPages := (total + filter.PageSize - 1) / filter.PageSize

	return &StudentListResult{
		Students:   students,
		Total:      total,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetByRollNo retrieves a student by roll number
func (m StudentModel) GetByRollNo(rollNo string) (*Student, error) {
	query := `
		SELECT s.id, s.official_email, s.name, s.roll_no, s.register_no, 
		       s.batch_id, b.year, s.photo_url, s.is_profile_completed, 
		       s.is_eligible_for_placement, s.placement_status,
		       s.created_at, s.updated_at, s.last_login_at, s.version
		FROM students s
		LEFT JOIN batches b ON s.batch_id = b.id
		WHERE s.roll_no = $1`

	var student Student
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, rollNo).Scan(
		&student.ID, &student.OfficialEmail, &student.Name, &student.RollNo,
		&student.RegisterNo, &student.BatchID, &student.BatchYear, &student.PhotoURL,
		&student.IsProfileCompleted, &student.IsEligibleForPlacement,
		&student.PlacementStatus, &student.CreatedAt, &student.UpdatedAt,
		&student.LastLoginAt, &student.Version,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &student, nil
}

// UpdatePlacementStatus updates a student's placement status
func (m StudentModel) UpdatePlacementStatus(studentID int64, status PlacementStatus) error {
	query := `UPDATE students SET placement_status = $1 WHERE id = $2`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, status, studentID)
	return err
}
