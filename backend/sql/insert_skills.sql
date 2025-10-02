-- Insert all the skills that will be used in the system
INSERT INTO skills (name, category) VALUES
-- Programming Languages
('C', 'Programming'),
('C++', 'Programming'),
('JAVA', 'Programming'),
('PYTHON', 'Programming'),
('Node.js', 'Programming'),
('SQL Database', 'Programming'),
('NoSQL Database', 'Programming'),
('Web Developement', 'Programming'),
('PHP', 'Programming'),
('Mobile App development-flutter', 'Programming'),
('Aptitude level', 'Programming'),
('logical and verbal Reasoning', 'Programming'),

-- Core Concepts
('DataStructure', 'Core Concept'),
('DBMS', 'Core Concept'),
('OOPS', 'Core Concept'),
('Problem Solving/Coding Tests', 'Core Concept'),
('Computer Networks', 'Core Concept'),
('Operating System', 'Core Concept'),
('Design and Analysis of Algorithm', 'Core Concept'),

-- Tools
('Git/Github', 'Tool'),
('Linux/Unix', 'Tool'),
('Cloud Basics (AWS/Azure/GCP)', 'Tool'),
('Competitive Coding (Codeforces/LeetCode/Hackerrank)', 'Tool'),
('Hacker Rank', 'Tool'),
('Hacker Earth', 'Tool')

ON CONFLICT (name) DO NOTHING;