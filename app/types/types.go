package types

type RegisterUserRequestType struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}

type ChallengesType struct {
	ID          int    `json:"id"`
	CreatedAt   string `json:"created_at"`
	Title       string `json:"title"`
	Description string `json:"description"`
	TypeID      int    `json:"type_id"`
	StatusID    int    `json:"status_id"`
	Type        string `json:"type"`
	Status      string `json:"status"`
	Marks       int    `json:"marks"`
}

type AddChallengeRequestType struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	TypeID      int    `json:"type_id"`
	StatusID    int    `json:"status_id"`
	Marks       int    `json:"marks"`
}

type DSAChallengesType struct {
	ID           int     `json:"id"`
	CreatedAt    string  `json:"created_at"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	TypeID       int     `json:"type_id"`
	StatusID     int     `json:"status_id"`
	Type         string  `json:"type"`
	Status       string  `json:"status"`
	SampleInput  string  `json:"sample_input"`
	SampleOutput string  `json:"sample_output"`
	Note         *string `json:"note"`
	Marks        int     `json:"marks"`
}

type AddDSAChallengeRequestType struct {
	Title        string            `json:"title"`
	Description  string            `json:"description"`
	TypeID       int               `json:"type_id"`
	StatusID     int               `json:"status_id"`
	Marks        int               `json:"marks"`
	SampleInput  string            `json:"sample_input"`
	SampleOutput string            `json:"sample_output"`
	Note         *string           `json:"note"`
	TestCases    []DSATestCaseType `json:"test_cases"`
}

type DSATestCaseType struct {
	TestInput  string `json:"test_input"`
	TestOutput string `json:"test_output"`
}

type TestDSAChallengeRequestType struct {
	ChallengeID    int    `json:"challenge_id"`
	LanguageID     int    `json:"language_id"`
	SourceCode     string `json:"source_code"`
	Stdin          string `json:"stdin"`
	ExpectedOutput string `json:"expected_output"`
}

type TestDSAChallengeResponse struct {
	Stdout        string `json:"stdout"`
	Stderr        string `json:"stderr"`
	Token         string `json:"token"`
	CompileOutput string `json:"compile_output"`
	Message       string `json:"message"`
	Status        struct {
		StatusID          int    `json:"id"`
		StatusDescription string `json:"description"`
	} `json:"status"`
}

type SubmitDSAChallengeRequestType struct {
	ChallengeID int    `json:"challenge_id"`
	LanguageID  int    `json:"language_id"`
	SourceCode  string `json:"source_code"`
}

type DSAChallengeTestCase struct {
	ID          int    `json:"id"`
	CreatedAt   string `json:"created_at"`
	ChallengeID int    `json:"challenge_id"`
	TestInput   string `json:"test_input"`
	TestOutput  string `json:"test_output"`
}

type Judge0SubmissionRequest struct {
	LanguageID     int    `json:"language_id"`
	SourceCode     string `json:"source_code"`
	CallbackURL    string `json:"callback_url"`
	Stdin          string `json:"stdin"`
	ExpectedOutput string `json:"expected_output"`
}

type Judge0BatchSubmissionRequest struct {
	Submissions []Judge0SubmissionRequest `json:"submissions"`
}

type LeaderboardUserType struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	XP     int    `json:"xp"`
}

type DSASubmissionType struct {
	ID               int    `json:"id"`
	CreatedAt        string `json:"created_at"`
	SubmissionID     string `json:"submission_id"`
	ChallengeID      int    `json:"challenge_id"`
	UserID           string `json:"user_id"`
	TestCount        int    `json:"test_count"`
	PassCount        int    `json:"pass_count"`
	FailCount        int    `json:"fail_count"`
	EvaluationStatus int    `json:"evaluation_status"`
}
