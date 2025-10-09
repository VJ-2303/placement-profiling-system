package models

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

// Extended models for the complete schema
type StudentDetails struct {
	StudentID             int64     `json:"student_id"`
	DateOfBirth           time.Time `json:"date_of_birth"`
	MobileNumber          string    `json:"mobile_number"`
	AlternateMobileNumber string    `json:"alternate_mobile_number"`
	PersonalEmail         string    `json:"personal_email"`
	LinkedinProfile       string    `json:"linkedin_profile"`
	Address               string    `json:"address"`
	City                  string    `json:"city"`
	Pincode               string    `json:"pincode"`
	AdhaarNo              string    `json:"adhaar_no"`
	ResidenceType         string    `json:"residence_type"`
	Strength              string    `json:"strength"`
	Weakness              string    `json:"weakness"`
	Remarks               string    `json:"remarks"`
}

type StudentParents struct {
	StudentID            int64  `json:"student_id"`
	FatherName           string `json:"father_name"`
	FatherMobile         string `json:"father_mobile"`
	FatherOccupation     string `json:"father_occupation"`
	FatherCompanyDetails string `json:"father_company_details"`
	FatherEmail          string `json:"father_email"`
	MotherName           string `json:"mother_name"`
	MotherMobile         string `json:"mother_mobile"`
	MotherOccupation     string `json:"mother_occupation"`
	MotherEmail          string `json:"mother_email"`
}

type StudentAcademics struct {
	StudentID         int64   `json:"student_id"`
	TenthPercentage   string  `json:"tenth_percentage"`
	TwelthPercentage  string  `json:"twelth_percentage"`
	CgpaSem1          *string `json:"cgpa_sem1"`
	CgpaSem2          *string `json:"cgpa_sem2"`
	CgpaSem3          *string `json:"cgpa_sem3"`
	CgpaSem4          *string `json:"cgpa_sem4"`
	CgpaOverall       string  `json:"cgpa_overall"`
	CurrentBacklogs   string  `json:"current_backlogs"`
	HasBacklogHistory string  `json:"has_backlog_history"`
}

type StudentAspirations struct {
	StudentID           int64  `json:"student_id"`
	CompanyAim          string `json:"company_aim"`
	TargetPackage       string `json:"target_package"`
	Certifications      string `json:"certifications"`
	Awards              string `json:"awards"`
	Workshops           string `json:"workshops"`
	Internships         string `json:"internships"`
	HackathonsAttended  string `json:"hackathons_attended"`
	Extracurriculars    string `json:"extracurriculars"`
	ClubParticipation   string `json:"club_participation"`
	FuturePath          string `json:"future_path"`
	CommunicationSkills string `json:"communication_skills"`
}

// FlatProfileResponse represents the complete flattened profile data
type FlatProfileResponse struct {
	Id                    int64   `json:"id"`
	RollNo                string  `json:"roll_no"`
	Name                  string  `json:"name"`
	OfficialEmail         string  `json:"official_email"`
	DateOfBirth           string  `json:"date_of_birth"`
	MobileNumber          string  `json:"mobile_number"`
	AltMobileNumber       string  `json:"alt_mobile_number"`
	PersonalEmail         string  `json:"personal_email"`
	LinkedInUrl           string  `json:"linkedin_url"`
	FatherName            string  `json:"father_name"`
	FatherMobile          string  `json:"father_mobile"`
	FatherOccupation      string  `json:"father_occupation"`
	FatherCompanyDetails  string  `json:"father_company_details"`
	FatherEmail           string  `json:"father_email"`
	MotherName            string  `json:"mother_name"`
	MotherMobile          string  `json:"mother_mobile"`
	MotherOccupation      string  `json:"mother_occupation"`
	MotherEmail           string  `json:"mother_email"`
	ResidenceType         string  `json:"residence_type"`
	Address               string  `json:"address"`
	City                  string  `json:"city"`
	Pincode               string  `json:"pincode"`
	AdhaarNo              string  `json:"adhaar_no"`
	Photo                 string  `json:"photo"`
	CompanyAim            string  `json:"company_aim"`
	TargetPackage         string  `json:"target_package"`
	TenthPercentage       string  `json:"tenth_percentage"`
	TwelthPercentage      string  `json:"twelth_percentage"`
	CgpaSem1              *string `json:"cgpa_sem1"`
	CgpaSem2              *string `json:"cgpa_sem2"`
	CgpaSem3              *string `json:"cgpa_sem3"`
	CgpaSem4              *string `json:"cgpa_sem4"`
	CgpaOverall           string  `json:"cgpa_overall"`
	CurrentBacklogs       string  `json:"current_backlogs"`
	HasBacklogHistory     string  `json:"has_backlog_history"`
	Certifications        string  `json:"certifications"`
	Awards                string  `json:"awards"`
	Workshops             string  `json:"workshops"`
	Internships           string  `json:"internships"`
	SkillC                string  `json:"skill_c"`
	SkillCpp              string  `json:"skill_cpp"`
	SkillJava             string  `json:"skill_java"`
	SkillPython           string  `json:"skill_python"`
	SkillNodeJs           string  `json:"skill_node_js"`
	SkillSql              string  `json:"skill_sql"`
	SkillNoSql            string  `json:"skill_no_sql"`
	SkillWebDev           string  `json:"skill_web_dev"`
	SkillPhp              string  `json:"skill_php"`
	SkillFlutter          string  `json:"skill_flutter"`
	SkillAptitude         string  `json:"skill_aptitude"`
	SkillReasoning        string  `json:"skill_reasoning"`
	ConceptDataStructures string  `json:"concept_data_structures"`
	ConceptDbms           string  `json:"concept_dbms"`
	ConceptOops           string  `json:"concept_oops"`
	ConceptProblemSolving string  `json:"concept_problem_solving"`
	ConceptNetworks       string  `json:"concept_networks"`
	ConceptOs             string  `json:"concept_os"`
	ConceptAlgos          string  `json:"concept_algos"`
	ToolGit               string  `json:"tool_git"`
	ToolLinux             string  `json:"tool_linux"`
	ToolCloud             string  `json:"tool_cloud"`
	ToolCompCoding        string  `json:"tool_comp_coding"`
	ToolHackerRank        string  `json:"tool_hacker_rank"`
	ToolHackerEarth       string  `json:"tool_hacker_earth"`
	CommunicationSkills   string  `json:"communication_skills"`
	HackathonsAttended    string  `json:"hackathons_attended"`
	Extracurriculars      string  `json:"extracurriculars"`
	ClubParticipation     string  `json:"club_participation"`
	FuturePath            string  `json:"future_path"`
	Strength              string  `json:"strength"`
	Weakness              string  `json:"weakness"`
	Remarks               string  `json:"remarks"`
}

// Model instances for the new tables
type StudentDetailsModel struct {
	DB *sql.DB
}

func (m StudentDetailsModel) Insert(tx *sql.Tx, details *StudentDetails) error {
	query := `
		INSERT INTO student_details (
			student_id, date_of_birth, mobile_number, alternate_mobile_number,
			personal_email, linkedin_profile, address, city, pincode, adhaar_no,
			residence_type, strength, weakness, remarks
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		ON CONFLICT (student_id) DO UPDATE SET
			date_of_birth = EXCLUDED.date_of_birth,
			mobile_number = EXCLUDED.mobile_number,
			alternate_mobile_number = EXCLUDED.alternate_mobile_number,
			personal_email = EXCLUDED.personal_email,
			linkedin_profile = EXCLUDED.linkedin_profile,
			address = EXCLUDED.address,
			city = EXCLUDED.city,
			pincode = EXCLUDED.pincode,
			adhaar_no = EXCLUDED.adhaar_no,
			residence_type = EXCLUDED.residence_type,
			strength = EXCLUDED.strength,
			weakness = EXCLUDED.weakness,
			remarks = EXCLUDED.remarks`

	args := []any{
		details.StudentID, details.DateOfBirth, details.MobileNumber,
		details.AlternateMobileNumber, details.PersonalEmail, details.LinkedinProfile,
		details.Address, details.City, details.Pincode, details.AdhaarNo,
		details.ResidenceType, details.Strength, details.Weakness, details.Remarks,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, args...)
	return err
}

type StudentParentsModel struct {
	DB *sql.DB
}

func (m StudentParentsModel) Insert(tx *sql.Tx, parents *StudentParents) error {
	query := `
		INSERT INTO student_parents (
			student_id, father_name, father_mobile, father_occupation,
			father_company_details, father_email, mother_name, mother_mobile,
			mother_occupation, mother_email
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (student_id) DO UPDATE SET
			father_name = EXCLUDED.father_name,
			father_mobile = EXCLUDED.father_mobile,
			father_occupation = EXCLUDED.father_occupation,
			father_company_details = EXCLUDED.father_company_details,
			father_email = EXCLUDED.father_email,
			mother_name = EXCLUDED.mother_name,
			mother_mobile = EXCLUDED.mother_mobile,
			mother_occupation = EXCLUDED.mother_occupation,
			mother_email = EXCLUDED.mother_email`

	args := []any{
		parents.StudentID, parents.FatherName, parents.FatherMobile,
		parents.FatherOccupation, parents.FatherCompanyDetails, parents.FatherEmail,
		parents.MotherName, parents.MotherMobile, parents.MotherOccupation, parents.MotherEmail,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, args...)
	return err
}

type StudentAcademicsModel struct {
	DB *sql.DB
}

func (m StudentAcademicsModel) Insert(tx *sql.Tx, academics *StudentAcademics) error {
	query := `
		INSERT INTO student_academics (
			student_id, tenth_percentage, twelth_percentage, cgpa_sem1,
			cgpa_sem2, cgpa_sem3, cgpa_sem4, cgpa_overall, current_backlogs,
			has_backlog_history
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (student_id) DO UPDATE SET
			tenth_percentage = EXCLUDED.tenth_percentage,
			twelth_percentage = EXCLUDED.twelth_percentage,
			cgpa_sem1 = EXCLUDED.cgpa_sem1,
			cgpa_sem2 = EXCLUDED.cgpa_sem2,
			cgpa_sem3 = EXCLUDED.cgpa_sem3,
			cgpa_sem4 = EXCLUDED.cgpa_sem4,
			cgpa_overall = EXCLUDED.cgpa_overall,
			current_backlogs = EXCLUDED.current_backlogs,
			has_backlog_history = EXCLUDED.has_backlog_history`

	args := []any{
		academics.StudentID, academics.TenthPercentage, academics.TwelthPercentage,
		academics.CgpaSem1, academics.CgpaSem2, academics.CgpaSem3, academics.CgpaSem4,
		academics.CgpaOverall, academics.CurrentBacklogs, academics.HasBacklogHistory,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, args...)
	return err
}

type StudentAspirationsModel struct {
	DB *sql.DB
}

func (m StudentAspirationsModel) Insert(tx *sql.Tx, aspirations *StudentAspirations) error {
	query := `
		INSERT INTO student_aspirations (
			student_id, company_aim, target_package, certifications, awards,
			workshops, internships, hackathons_attended, extracurriculars,
			club_participation, future_path, communication_skills
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT (student_id) DO UPDATE SET
			company_aim = EXCLUDED.company_aim,
			target_package = EXCLUDED.target_package,
			certifications = EXCLUDED.certifications,
			awards = EXCLUDED.awards,
			workshops = EXCLUDED.workshops,
			internships = EXCLUDED.internships,
			hackathons_attended = EXCLUDED.hackathons_attended,
			extracurriculars = EXCLUDED.extracurriculars,
			club_participation = EXCLUDED.club_participation,
			future_path = EXCLUDED.future_path,
			communication_skills = EXCLUDED.communication_skills`

	args := []any{
		aspirations.StudentID, aspirations.CompanyAim, aspirations.TargetPackage,
		aspirations.Certifications, aspirations.Awards, aspirations.Workshops,
		aspirations.Internships, aspirations.HackathonsAttended, aspirations.Extracurriculars,
		aspirations.ClubParticipation, aspirations.FuturePath, aspirations.CommunicationSkills,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, args...)
	return err
}

type SkillsModel struct {
	DB *sql.DB
}

func (m SkillsModel) GetAllAsMap(tx *sql.Tx) (map[string]int, error) {
	query := `SELECT id, name FROM skills`

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	skillMap := make(map[string]int)
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, err
		}
		skillMap[name] = id
	}

	return skillMap, rows.Err()
}

type StudentSkillsModel struct {
	DB *sql.DB
}

func (m StudentSkillsModel) Insert(tx *sql.Tx, studentID int64, skillID int, proficiency string) error {
	query := `
		INSERT INTO student_skills (student_id, skill_id, proficiency_level)
		VALUES ($1, $2, $3)
		ON CONFLICT (student_id, skill_id) DO UPDATE SET
			proficiency_level = EXCLUDED.proficiency_level`

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, studentID, skillID, proficiency)
	return err
}

// Add method to get full profile with all data
func (m StudentModel) GetFullProfile(id int64) (*FlatProfileResponse, error) {
	query := `
		SELECT
			COALESCE(s.roll_no, ''), s.name, s.official_email,
			COALESCE(sd.date_of_birth::text, ''), COALESCE(sd.mobile_number, ''),
			COALESCE(sd.alternate_mobile_number, ''), COALESCE(sd.personal_email, ''),
			COALESCE(sd.linkedin_profile, ''), COALESCE(sd.address, ''), COALESCE(sd.city, ''),
			COALESCE(sd.pincode, ''), COALESCE(sd.adhaar_no, ''), COALESCE(sd.residence_type, ''),
			COALESCE(sd.strength, ''), COALESCE(sd.weakness, ''), COALESCE(sd.remarks, ''),
			COALESCE(sp.father_name, ''), COALESCE(sp.father_mobile, ''), COALESCE(sp.father_occupation, ''),
			COALESCE(sp.father_company_details, ''), COALESCE(sp.father_email, ''), COALESCE(sp.mother_name, ''),
			COALESCE(sp.mother_mobile, ''), COALESCE(sp.mother_occupation, ''), COALESCE(sp.mother_email, ''),
			COALESCE(sa.tenth_percentage, ''), COALESCE(sa.twelth_percentage, ''), sa.cgpa_sem1, sa.cgpa_sem2,
			sa.cgpa_sem3, sa.cgpa_sem4, COALESCE(sa.cgpa_overall, ''), COALESCE(sa.current_backlogs, ''),
			COALESCE(sa.has_backlog_history, ''),
			COALESCE(sas.company_aim, ''), COALESCE(sas.target_package, ''), COALESCE(sas.certifications, ''),
			COALESCE(sas.awards, ''), COALESCE(sas.workshops, ''), COALESCE(sas.internships, ''),
			COALESCE(sas.hackathons_attended, ''), COALESCE(sas.extracurriculars, ''), COALESCE(sas.club_participation, ''),
			COALESCE(sas.future_path, ''), COALESCE(sas.communication_skills, '')
		FROM students s
		LEFT JOIN student_details sd ON s.id = sd.student_id
		LEFT JOIN student_parents sp ON s.id = sp.student_id
		LEFT JOIN student_academics sa ON s.id = sa.student_id
		LEFT JOIN student_aspirations sas ON s.id = sas.student_id
		WHERE s.id = $1`

	var profile FlatProfileResponse

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&profile.RollNo, &profile.Name, &profile.OfficialEmail,
		&profile.DateOfBirth, &profile.MobileNumber, &profile.AltMobileNumber,
		&profile.PersonalEmail, &profile.LinkedInUrl, &profile.Address, &profile.City,
		&profile.Pincode, &profile.AdhaarNo, &profile.ResidenceType, &profile.Strength,
		&profile.Weakness, &profile.Remarks,
		&profile.FatherName, &profile.FatherMobile, &profile.FatherOccupation,
		&profile.FatherCompanyDetails, &profile.FatherEmail, &profile.MotherName,
		&profile.MotherMobile, &profile.MotherOccupation, &profile.MotherEmail,
		&profile.TenthPercentage, &profile.TwelthPercentage, &profile.CgpaSem1,
		&profile.CgpaSem2, &profile.CgpaSem3, &profile.CgpaSem4, &profile.CgpaOverall,
		&profile.CurrentBacklogs, &profile.HasBacklogHistory,
		&profile.CompanyAim, &profile.TargetPackage, &profile.Certifications,
		&profile.Awards, &profile.Workshops, &profile.Internships,
		&profile.HackathonsAttended, &profile.Extracurriculars, &profile.ClubParticipation,
		&profile.FuturePath, &profile.CommunicationSkills,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	// Fetch skills separately
	skillRows, err := m.DB.QueryContext(ctx, `
		SELECT s.name, ss.proficiency_level
		FROM skills s
		JOIN student_skills ss ON s.id = ss.skill_id
		WHERE ss.student_id = $1`, id)
	if err != nil {
		return nil, err
	}
	defer skillRows.Close()

	for skillRows.Next() {
		var skillName, proficiency string
		if err := skillRows.Scan(&skillName, &proficiency); err != nil {
			return nil, err
		}

		// Map skills to the flat structure
		switch skillName {
		case "C":
			profile.SkillC = proficiency
		case "C++":
			profile.SkillCpp = proficiency
		case "JAVA":
			profile.SkillJava = proficiency
		case "PYTHON":
			profile.SkillPython = proficiency
		case "Node.js":
			profile.SkillNodeJs = proficiency
		case "SQL Database":
			profile.SkillSql = proficiency
		case "NoSQL Database":
			profile.SkillNoSql = proficiency
		case "Web Developement":
			profile.SkillWebDev = proficiency
		case "PHP":
			profile.SkillPhp = proficiency
		case "Mobile App development-flutter":
			profile.SkillFlutter = proficiency
		case "Aptitude level":
			profile.SkillAptitude = proficiency
		case "logical and verbal Reasoning":
			profile.SkillReasoning = proficiency
		case "DataStructure":
			profile.ConceptDataStructures = proficiency
		case "DBMS":
			profile.ConceptDbms = proficiency
		case "OOPS":
			profile.ConceptOops = proficiency
		case "Problem Solving/Coding Tests":
			profile.ConceptProblemSolving = proficiency
		case "Computer Networks":
			profile.ConceptNetworks = proficiency
		case "Operating System":
			profile.ConceptOs = proficiency
		case "Design and Analysis of Algorithm":
			profile.ConceptAlgos = proficiency
		case "Git/Github":
			profile.ToolGit = proficiency
		case "Linux/Unix":
			profile.ToolLinux = proficiency
		case "Cloud Basics (AWS/Azure/GCP)":
			profile.ToolCloud = proficiency
		case "Competitive Coding (Codeforces/LeetCode/Hackerrank)":
			profile.ToolCompCoding = proficiency
		case "Hacker Rank":
			profile.ToolHackerRank = proficiency
		case "Hacker Earth":
			profile.ToolHackerEarth = proficiency
		}
	}

	return &profile, nil
}
