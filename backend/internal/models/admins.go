package models

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Admin struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	OfficialEmail string    `json:"email"`
	CreatedAt     time.Time `json:"created_at"`
}

type AdminModel struct {
	DB *sql.DB
}

func (m AdminModel) Insert(admin *Admin) error {

	query := `
		INSERT INTO admins (name, email)
		VALUES ($1,$2)
		RETURNING id, created_at
	`
	args := []any{admin.Name, admin.OfficialEmail}

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(
		&admin.ID, &admin.CreatedAt,
	)
}

func (m AdminModel) GetByEmail(email string) (*Admin, error) {

	if email == "" {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, name , email, created_at from admins where email = $1
	`

	var admin Admin

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&admin.ID,
		&admin.Name,
		&admin.OfficialEmail,
		&admin.CreatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &admin, nil
}

func (m AdminModel) GetByID(id int64) (*Admin, error) {

	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, name , email, created_at from admins where id = $1
	`

	var admin Admin

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&admin.ID,
		&admin.Name,
		&admin.OfficialEmail,
		&admin.CreatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &admin, nil
}

func (m AdminModel) GetFullProfileByRollNo(rollNo string) (*FlatProfileResponse, error) {
	query := `
		SELECT
			COALESCE(s.roll_no, ''),s.id, s.name, s.official_email,
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
		WHERE s.roll_no = $1`

	var profile FlatProfileResponse

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, rollNo).Scan(
		&profile.RollNo, &profile.Id, &profile.Name, &profile.OfficialEmail,
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
		WHERE ss.student_id = $1`, profile.Id)
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
