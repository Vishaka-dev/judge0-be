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
