package repositories

import (
	"context"

	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/database"
	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/types"
)

func GetAllChallenges() ([]types.ChallengesPreviewType, error) {
	pool := database.GetPool()

	query := `SELECT * FROM preview_challenges_view`

	rows, err := pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var challenges []types.ChallengesPreviewType

	for rows.Next() {
		var c types.ChallengesPreviewType
		if err := rows.Scan(
			&c.ID,
			&c.CreatedAt,
			&c.Title,
			&c.Description,
			&c.TypeID,
			&c.StatusID,
			&c.Type,
			&c.Status,
		); err != nil {
			return nil, err
		}
		challenges = append(challenges, c)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return challenges, nil
}
