package types

type ChallengesType struct {
	ID          int    `json:"id"`
	CreatedAt   string `json:"created_at"`
	Title       string `json:"title"`
	Description string `json:"description"`
	TypeID      int    `json:"type_id"`
	StatusID    int    `json:"status_id"`
	Type        string `json:"type"`
	Status      string `json:"status"`
}

type AddChallengeRequestType struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	TypeID      int    `json:"type_id"`
	StatusID    int    `json:"status_id"`
}

type DSAChallengesType struct {
	ID           int    `json:"id"`
	CreatedAt    string `json:"created_at"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	TypeID       int    `json:"type_id"`
	StatusID     int    `json:"status_id"`
	Type         string `json:"type"`
	Status       string `json:"status"`
	SampleInput  string `json:"sample_input"`
	SampleOutput string `json:"sample_output"`
	Note         string `json:"note"`
}

type AddDSAChallengeRequestType struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	TypeID       int    `json:"type_id"`
	StatusID     int    `json:"status_id"`
	SampleInput  string `json:"sample_input"`
	SampleOutput string `json:"sample_output"`
	Note         string `json:"note"`
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
