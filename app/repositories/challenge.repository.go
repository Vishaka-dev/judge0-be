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
		"select count(*) from challenges").Scan(&count)

	if err != nil {
		log.Println("Error:", err)
		return nil, 0, 0, err
	}

	ps, err := strconv.ParseInt(pageSize, 10, 64)
	if err != nil || ps <= 0 {
		return nil, 0, 0, err
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

	log.Println("Count:", count)
	log.Println("Total Pages:", totalPages)
	log.Println("Current Page:", page)

	rows, err := pool.Query(context.Background(),
		"select id, created_at, title, description, type_id, status_id from challenges order by id desc limit $1 offset $2",
		ps, offset)
	if err != nil {
		log.Println("Query Error:", err)
		return nil, 0, 0, err
	}
	defer rows.Close()

	challenges := []types.ChallengesPreviewType{}
	for rows.Next() {
		var challenge types.ChallengesPreviewType
		err := rows.Scan(&challenge.ID, &challenge.CreatedAt, &challenge.Title, &challenge.Description, &challenge.TypeID, &challenge.StatusID)
		if err != nil {
			log.Println("Scan Error:", err)
			continue
		}
		challenges = append(challenges, challenge)
	}

	return challenges, page, totalPages, nil
}
