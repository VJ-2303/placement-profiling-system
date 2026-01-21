-- Fix: Add UNIQUE constraint to student_achievements table
-- This is required for ON CONFLICT upsert operations

-- Add unique constraint if it doesn't exist
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint 
        WHERE conname = 'student_achievements_student_id_key'
    ) THEN
        -- First, delete duplicate rows keeping only the latest one per student
        DELETE FROM student_achievements a
        USING student_achievements b
        WHERE a.student_id = b.student_id
        AND a.id < b.id;
        
        -- Then add the unique constraint
        ALTER TABLE student_achievements ADD CONSTRAINT student_achievements_student_id_key UNIQUE (student_id);
    END IF;
END $$;
