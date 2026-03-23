package repositories

import (
	"context"
	"errors"
	"strconv"

	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/database"
	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/logger"
	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/types"
	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/utils"
)

func GetAllChallenges(ctx context.Context, limit, pageSize string) ([]types.ChallengesType, int64, int64, error) {
	pool := database.GetPool()
	ctx, cancel := utils.WithTimeout(ctx)
	defer cancel()

	var count int64
	err := pool.QueryRow(ctx,
		"select count(*) from preview_challenges_view").Scan(&count)
	if err != nil {
		logger.Log.Error("Error counting challenges", "error", err)
		return nil, 0, 0, err
	}

	ps, err := strconv.ParseInt(pageSize, 10, 64)
	if err != nil || ps <= 0 {
		logger.Log.Warn("Invalid pageSize, defaulting to 10", "input", pageSize, "error", err)
		ps = 10
	}

	page, err := strconv.ParseInt(limit, 10, 64)
	if err != nil || page <= 0 {
		logger.Log.Warn("Invalid limit, defaulting to 1", "input", limit, "error", err)
		page = 1
	}

	totalPages := (count + ps - 1) / ps
	if page > totalPages && totalPages > 0 {
		page = totalPages
	}

	offset := (page - 1) * ps

	rows, err := pool.Query(ctx,
		`select id, created_at, title, description, type_id, status_id, type, status, marks
			from preview_challenges_view
			where status_id = 2
			order by id desc
			limit $1 offset $2`,
		ps, offset)
	if err != nil {
		logger.Log.Error("Query Error", "error", err)
		return nil, 0, 0, err
	}
	defer rows.Close()

	challenges := []types.ChallengesType{}
	for rows.Next() {
		var challenge types.ChallengesType
		if err := rows.Scan(
			&challenge.ID,
			&challenge.CreatedAt,
			&challenge.Title,
			&challenge.Description,
			&challenge.TypeID,
			&challenge.StatusID,
			&challenge.Type,
			&challenge.Status,
			&challenge.Marks,
		); err != nil {
			logger.Log.Error("Scan Error", "error", err)
			return nil, 0, 0, err
		}
		challenges = append(challenges, challenge)
	}

	if err := rows.Err(); err != nil {
		logger.Log.Error("Rows Error", "error", err)
		return nil, 0, 0, err
	}

	logger.Log.Info("Fetched challenges", "count", len(challenges), "page", page, "totalPages", totalPages)
	return challenges, page, totalPages, nil
}

func GetChallengeType(ctx context.Context, id string) (int, error) {
	pool := database.GetPool()
	ctx, cancel := utils.WithTimeout(ctx)
	defer cancel()
	var challengeType int
	err := pool.QueryRow(ctx,
		"select type_id from challenges where id = $1", id).Scan(&challengeType)
	if err != nil {
		logger.Log.Error("GetChallengeType error", "id", id, "error", err)
	}
	return challengeType, err
}

func GetMarksForChallenge(ctx context.Context, challengeId int) (int, error) {
	pool := database.GetPool()
	ctx, cancel := utils.WithTimeout(ctx)
	defer cancel()

	var marks int
	err := pool.QueryRow(ctx,
		`SELECT marks FROM challenges WHERE id = $1`, challengeId).Scan(&marks)
	if err != nil {
		logger.Log.Error("GetMarksForChallenge: query error", "challenge_id", challengeId, "error", err)
		return 0, err
	}
	return marks, nil
}

func GetChallengeIDBySubmissionID(ctx context.Context, submissionId string) (int, error) {
	pool := database.GetPool()
	ctx, cancel := utils.WithTimeout(ctx)
	defer cancel()
	var challengeId int
	err := pool.QueryRow(ctx, "SELECT challenge_id FROM dsa_submissions WHERE submission_id = $1", submissionId).Scan(&challengeId)
	return challengeId, err
}

func GetUserIDBySubmissionID(ctx context.Context, submissionId string) (string, error) {
	pool := database.GetPool()
	ctx, cancel := utils.WithTimeout(ctx)
	defer cancel()

	var userID string
	err := pool.QueryRow(ctx, "SELECT user_id FROM dsa_submissions WHERE submission_id = $1", submissionId).Scan(&userID)
	return userID, err
}

func AddMarksToLeaderboard(ctx context.Context, userID string, marks int) error {
	if marks <= 0 {
		return nil
	}

	pool := database.GetPool()
	ctx, cancel := utils.WithTimeout(ctx)
	defer cancel()

	cmdTag, err := pool.Exec(ctx,
		"UPDATE leaderboard SET marks = marks + $2 WHERE user_id = $1",
		userID,
		marks,
	)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		_, err = pool.Exec(ctx,
			"INSERT INTO leaderboard (user_id, marks) VALUES ($1, $2)",
			userID,
			marks,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// DSA Challenges

func GetDSAChallenge(ctx context.Context, id string) (types.DSAChallengesType, error) {
	pool := database.GetPool()
	ctx, cancel := utils.WithTimeout(ctx)
	defer cancel()
	var challenge types.DSAChallengesType
	err := pool.QueryRow(ctx,
		`select id, created_at, title, description,marks, type_id, status_id, type, status,
			sample_input, sample_output, note
			from get_dsa_challenges_view where id = $1 and status_id = 2`, id).Scan(
		&challenge.ID,
		&challenge.CreatedAt,
		&challenge.Title,
		&challenge.Description,
		&challenge.Marks,
		&challenge.TypeID,
		&challenge.StatusID,
		&challenge.Type,
		&challenge.Status,
		&challenge.SampleInput,
		&challenge.SampleOutput,
		&challenge.Note,
	)
	if err != nil {
		logger.Log.Error("GetDSAChallenge error", "id", id, "error", err)
	}
	return challenge, err
}

func AddDSAChallenge(ctx context.Context, challenge types.AddDSAChallengeRequestType) (int, error) {
	pool := database.GetPool()
	ctx, cancel := utils.WithTimeout(ctx)
	defer cancel()

	tx, err := pool.Begin(ctx)
	if err != nil {
		logger.Log.Error("AddDSAChallenge: begin transaction error", "error", err)
		return 0, err
	}
	defer tx.Rollback(ctx)

	var challengeID int
	err = tx.QueryRow(ctx,
		`INSERT INTO challenges (title, description, type_id, status_id, marks)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id`,
		challenge.Title,
		challenge.Description,
		challenge.TypeID,
		challenge.StatusID,
		challenge.Marks,
	).Scan(&challengeID)
	if err != nil {
		logger.Log.Error("AddDSAChallenge: insert challenge error", "error", err)
		return 0, err
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO dsa_challenges (challenge_id, sample_input, sample_output, note)
			VALUES ($1, $2, $3, $4)`,
		challengeID,
		challenge.SampleInput,
		challenge.SampleOutput,
		challenge.Note,
	)
	if err != nil {
		logger.Log.Error("AddDSAChallenge: insert dsa_challenge error", "error", err)
		return 0, err
	}

	for _, testCase := range challenge.TestCases {
		_, err = tx.Exec(ctx,
			`INSERT INTO dsa_test_cases (challenge_id, test_input, test_output)
				VALUES ($1, $2, $3)`,
			challengeID,
			testCase.TestInput,
			testCase.TestOutput,
		)
		if err != nil {
			logger.Log.Error("AddDSAChallenge: insert dsa_test_case error", "challenge_id", challengeID, "error", err)
			return 0, err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		logger.Log.Error("AddDSAChallenge: commit error", "error", err)
		return 0, err
	}

	logger.Log.Info("DSA Challenge added", "challenge_id", challengeID)
	return challengeID, nil
}

func GetDSAChallengeTestCases(ctx context.Context, challengeID int) ([]types.DSAChallengeTestCase, error) {
	pool := database.GetPool()
	ctx, cancel := utils.WithTimeout(ctx)
	defer cancel()

	rows, err := pool.Query(ctx,
		`SELECT id, challenge_id, test_input, test_output FROM dsa_test_cases WHERE challenge_id = $1 ORDER BY id ASC`,
		challengeID,
	)
	if err != nil {
		logger.Log.Error("GetDSAChallengeTestCases: query error", "challenge_id", challengeID, "error", err)
		return nil, err
	}
	defer rows.Close()

	var testCases []types.DSAChallengeTestCase
	for rows.Next() {
		var tc types.DSAChallengeTestCase
		if err := rows.Scan(
			&tc.ID,
			&tc.ChallengeID,
			&tc.TestInput,
			&tc.TestOutput,
		); err != nil {
			logger.Log.Error("GetDSAChallengeTestCases: scan error", "error", err)
			return nil, err
		}
		testCases = append(testCases, tc)
	}
	if err := rows.Err(); err != nil {
		logger.Log.Error("GetDSAChallengeTestCases: rows error", "error", err)
		return nil, err
	}
	logger.Log.Info("Fetched DSA challenge test cases", "challenge_id", challengeID, "count", len(testCases))
	return testCases, nil
}

func GetDSATestCaseCount(ctx context.Context, challengeID int) (int64, error) {
	pool := database.GetPool()
	ctx, cancel := utils.WithTimeout(ctx)
	defer cancel()
	var count int64
	err := pool.QueryRow(ctx,
		"SELECT COUNT(*) FROM dsa_test_cases WHERE challenge_id = $1", challengeID).Scan(&count)
	return count, err
}

func AddDSASubmission(ctx context.Context, submissionId string, challengeId int, userId string, testCount int) (bool, error) {
	pool := database.GetPool()
	ctx, cancel := utils.WithTimeout(ctx)
	defer cancel()
	_, err := pool.Exec(ctx,
		`INSERT INTO dsa_submissions (submission_id, challenge_id, user_id,test_count) VALUES ($1, $2, $3,$4)`,
		submissionId, challengeId, userId, testCount,
	)

	if err != nil {
		logger.Log.Error("AddDSASubmission: insert error", "submission_id", submissionId, "error", err)
		return false, err
	}
	return true, nil
}

func GetDSASubmissionCount(ctx context.Context, submissionId string) (int64, error) {
	pool := database.GetPool()
	ctx, cancel := utils.WithTimeout(ctx)
	defer cancel()
	var count int64
	err := pool.QueryRow(ctx, "SELECT COUNT(*) FROM dsa_submission_results WHERE submission_id = $1", submissionId).Scan(&count)
	return count, err
}

func HasUserPassedDSAChallenge(ctx context.Context, userId string, challengeId int) (bool, error) {
	pool := database.GetPool()
	ctx, cancel := utils.WithTimeout(ctx)
	defer cancel()

	var count int64
	err := pool.QueryRow(ctx,
		`SELECT COUNT(*)
		 FROM dsa_submissions
		 WHERE user_id = $1
		   AND challenge_id = $2
		   AND evaluation_status = 2`,
		userId,
		challengeId,
	).Scan(&count)
	if err != nil {
		logger.Log.Error("HasUserPassedDSAChallenge: query error", "user_id", userId, "challenge_id", challengeId, "error", err)
		return false, err
	}

	return count > 0, nil
}

func GetPassDSASubmissionCount(ctx context.Context, submissionId string) (int64, error) {
	pool := database.GetPool()
	ctx, cancel := utils.WithTimeout(ctx)
	defer cancel()
	var count int64
	err := pool.QueryRow(ctx, "SELECT COUNT(*) FROM dsa_submission_results WHERE submission_id = $1 AND status = 2", submissionId).Scan(&count)
	return count, err
}

func AddDSASubmissionResult(ctx context.Context, submissionId string, payload types.TestDSAChallengeResponse) (bool, error) {
	pool := database.GetPool()
	ctx, cancel := utils.WithTimeout(ctx)
	defer cancel()

	if payload.Token == "" {
		err := errors.New("missing callback token")
		logger.Log.Error("AddDSASubmissionResult: invalid payload", "submission_id", submissionId, "error", err)
		return false, err
	}

	statusID := payload.Status.StatusID
	status := 2
	switch {
	case statusID < 3:
		status = 1
	case statusID > 3:
		status = 3
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		logger.Log.Error("AddDSASubmissionResult: begin transaction error", "submission_id", submissionId, "error", err)
		return false, err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`INSERT INTO dsa_submission_results (submission_id, status, token)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (token)
		 DO UPDATE SET
		 	submission_id = EXCLUDED.submission_id,
		 	status = EXCLUDED.status`,
		submissionId,
		status,
		payload.Token,
	)
	if err != nil {
		logger.Log.Error("AddDSASubmissionResult: upsert error", "submission_id", submissionId, "token", payload.Token, "error", err)
		return false, err
	}

	if err = tx.Commit(ctx); err != nil {
		logger.Log.Error("AddDSASubmissionResult: commit error", "submission_id", submissionId, "error", err)
		return false, err
	}

	return true, nil
}

func UpdateDSASubmission(ctx context.Context, submissionId string, payload types.TestDSAChallengeResponse) (bool, error) {
	pool := database.GetPool()
	ctx, cancel := utils.WithTimeout(ctx)
	defer cancel()

	if payload.Token == "" {
		err := errors.New("missing callback token")
		logger.Log.Error("AddDSASubmissionResult: invalid payload", "submission_id", submissionId, "error", err)
		return false, err
	}

	challengeId, err := GetChallengeIDBySubmissionID(context.Background(), submissionId)
	if err != nil {
		logger.Log.Error("UpdateDSASubmission: failed to get challenge ID", "submission_id", submissionId, "error", err)
		return false, err
	}

	testCount, err := GetDSATestCaseCount(context.Background(), challengeId)
	if err != nil {
		logger.Log.Error("UpdateDSASubmission: failed to get test case count", "submission_id", submissionId, "error", err)
		return false, err
	}
	submissionCount, err := GetDSASubmissionCount(context.Background(), submissionId)
	if err != nil {
		logger.Log.Error("UpdateDSASubmission: failed to get submission count", "submission_id", submissionId, "error", err)
		return false, err
	}

	if submissionCount < testCount {
		return true, nil
	}

	passCount, err := GetPassDSASubmissionCount(context.Background(), submissionId)
	if err != nil {
		logger.Log.Error("UpdateDSASubmission: failed to count successful results", "submission_id", submissionId, "error", err)
		return false, err
	}

	failCount := submissionCount - passCount
	evaluationStatus := 3
	if passCount >= testCount {
		evaluationStatus = 2
	}

	_, err = pool.Exec(ctx,
		`UPDATE dsa_submissions
		 SET pass_count = $1,
		 	fail_count = $2,
		 	evaluation_status = $3
		 WHERE submission_id = $4`,
		passCount,
		failCount,
		evaluationStatus,
		submissionId,
	)
	if err != nil {
		logger.Log.Error("UpdateDSASubmission: failed to update submission aggregates", "submission_id", submissionId, "error", err)
		return false, err
	}

	return true, nil

}
