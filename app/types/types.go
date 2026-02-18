package types

type ChallengesPreviewType struct {
	ID          int    `json:"id"`
	CreatedAt   string `json:"created_at"`
	Title       string `json:"title"`
	Description string `json:"description"`
	TypeID      int    `json:"type_id"`
	StatusID    int    `json:"status_id"`
	Type        string `json:"type"`
	Status      string `json:"status"`
}
