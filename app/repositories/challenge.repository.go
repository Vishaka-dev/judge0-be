package repositories

import (
	"context"
	"log"
	"strconv"

	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/database"
	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/types"
)

func GetAllChallenges(limit, pageSize string) ([]types.ChallengesPreviewType, int64, int64, error) {
	pool := database.GetPool()

	var count int64
	err := pool.QueryRow(context.Background(),
		"select count(*) from preview_challenges_view").Scan(&count)
	if err != nil {
		log.Println("Error counting challenges:", err)
		return nil, 0, 0, err
	}

	ps, err := strconv.ParseInt(pageSize, 10, 64)
	if err != nil || ps <= 0 {
		ps = 10
	}

	page, err := strconv.ParseInt(limit, 10, 64)
	if err != nil || page <= 0 {
		page = 1
	}

	totalPages := (count + ps - 1) / ps
	if page > totalPages && totalPages > 0 {
		page = totalPages
	}

	offset := (page - 1) * ps

	rows, err := pool.Query(context.Background(),
		`select id, created_at, title, description, type_id, status_id, type, status
		 from preview_challenges_view
		 order by id desc
		 limit $1 offset $2`,
		ps, offset)
	if err != nil {
		log.Println("Query Error:", err)
		return nil, 0, 0, err
	}
	defer rows.Close()

	challenges := []types.ChallengesPreviewType{}
	for rows.Next() {
		var challenge types.ChallengesPreviewType
		if err := rows.Scan(
			&challenge.ID,
			&challenge.CreatedAt,
			&challenge.Title,
			&challenge.Description,
			&challenge.TypeID,
			&challenge.StatusID,
			&challenge.Type,
			&challenge.Status,
		); err != nil {
			log.Println("Scan Error:", err)
			continue
		}
		challenges = append(challenges, challenge)
	}

	if err := rows.Err(); err != nil {
		log.Println("Rows Error:", err)
		return nil, 0, 0, err
	}

	return challenges, page, totalPages, nil
}

func GetChallengeType(id string) (string, error) {
	pool := database.GetPool()
	var challengeType string
	err := pool.QueryRow(context.Background(),
		"select type_id  from challenges where id = $1", id).Scan(&challengeType)
	return challengeType, err
}

func GetDSAChallenge(id string) (types.DSAChallengesType, error) {
	pool := database.GetPool()
	var challenge types.DSAChallengesType
	err := pool.QueryRow(context.Background(),
		`select id, created_at, title, description, type_id, status_id, type, status,
		 sample_input, sample_output, note
		 from get_dsa_challenges_view where id = $1`, id).Scan(
		&challenge.ID,
		&challenge.CreatedAt,
		&challenge.Title,
		&challenge.Description,
		&challenge.TypeID,
		&challenge.StatusID,
		&challenge.Type,
		&challenge.Status,
		&challenge.SampleInput,
		&challenge.SampleOutput,
		&challenge.Note,
	)
	return challenge, err
}
