package repositories

import (
	"context"
	"errors"
	"strconv"

	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/database"
	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/logger"
	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/types"
	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/utils"
	"github.com/jackc/pgx/v5"
)

func GetLeaderboard(ctx context.Context, page, pageSize string) ([]types.LeaderboardUserType, int64, int64, error) {
	pool := database.GetPool()
	ctx, cancel := utils.WithTimeout(ctx)
	defer cancel()

	var count int64
	err := pool.QueryRow(ctx,
		"SELECT count(*) FROM leaderboard").Scan(&count)
	if err != nil {
		logger.Log.Error("Error counting leaderboard entries", "error", err)
		return nil, 0, 0, err
	}

	ps, err := strconv.ParseInt(pageSize, 10, 64)
	if err != nil || ps <= 0 {
		logger.Log.Warn("Invalid pageSize, defaulting to 10", "input", pageSize, "error", err)
		ps = 10
	}

	p, err := strconv.ParseInt(page, 10, 64)
	if err != nil || p <= 0 {
		logger.Log.Warn("Invalid page, defaulting to 1", "input", page, "error", err)
		p = 1
	}

	totalPages := (count + ps - 1) / ps
	if p > totalPages && totalPages > 0 {
		p = totalPages
	}

	offset := (p - 1) * ps

	rows, err := pool.Query(ctx,
		`SELECT l.user_id, u.name, l.marks
			FROM leaderboard l
			JOIN users u ON u.user_id = l.user_id
			ORDER BY l.marks DESC, l.user_id ASC
			LIMIT $1 OFFSET $2`,
		ps, offset)
	if err != nil {
		logger.Log.Error("Query Error", "error", err)
		return nil, 0, 0, err
	}
	defer rows.Close()

	users := []types.LeaderboardUserType{}
	for rows.Next() {
		var user types.LeaderboardUserType
		if err := rows.Scan(
			&user.UserID,
			&user.Name,
			&user.XP,
		); err != nil {
			logger.Log.Error("Scan Error", "error", err)
			return nil, 0, 0, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		logger.Log.Error("Rows Error", "error", err)
		return nil, 0, 0, err
	}

	logger.Log.Info("Fetched leaderboard", "count", len(users), "page", p, "totalPages", totalPages)
	return users, p, totalPages, nil
}

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
	if err != nil {
		logger.Log.Error("GetChallengeIDBySubmissionID: query error", "submission_id", submissionId, "error", err)
		return 0, err
	}
	return challengeId, err
}

func GetUserIDBySubmissionID(ctx context.Context, submissionId string) (string, error) {
	pool := database.GetPool()
	ctx, cancel := utils.WithTimeout(ctx)
	defer cancel()

	var userID string
	err := pool.QueryRow(ctx, "SELECT user_id FROM dsa_submissions WHERE submission_id = $1", submissionId).Scan(&userID)
	if err != nil {
		logger.Log.Error("GetUserIDBySubmissionID: query error", "submission_id", submissionId, "error", err)
		return "", err
	}
	return userID, err
}

func AddMarksToLeaderboard(ctx context.Context, userID string, marks int) error {
	if marks <= 0 {
		logger.Log.Info("AddMarksToLeaderboard: skipped non-positive marks", "user_id", userID, "marks", marks)
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
		logger.Log.Error("AddMarksToLeaderboard: update error", "user_id", userID, "marks", marks, "error", err)
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		_, err = pool.Exec(ctx,
			"INSERT INTO leaderboard (user_id, marks) VALUES ($1, $2)",
			userID,
			marks,
		)
		if err != nil {
			logger.Log.Error("AddMarksToLeaderboard: insert error", "user_id", userID, "marks", marks, "error", err)
			return err
		}
		logger.Log.Info("AddMarksToLeaderboard: leaderboard row created", "user_id", userID, "marks_added", marks)
		return nil
	}

	logger.Log.Info("AddMarksToLeaderboard: leaderboard marks updated", "user_id", userID, "marks_added", marks)

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
	if err != nil {
		logger.Log.Error("GetDSASubmissionCount: query error", "submission_id", submissionId, "error", err)
		return 0, err
	}
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
	if err != nil {
		logger.Log.Error("GetPassDSASubmissionCount: query error", "submission_id", submissionId, "error", err)
		return 0, err
	}
	return count, err
}

func AddDSASubmissionResult(ctx context.Context, submissionId string, payload types.TestDSAChallengeResponse) (bool, error) {
	pool := database.GetPool()
	ctx, cancel := utils.WithTimeout(ctx)
	defer cancel()
	logger.Log.Info("AddDSASubmissionResult: processing callback", "submission_id", submissionId, "token", payload.Token, "judge0_status_id", payload.Status.StatusID)

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

	logger.Log.Info("AddDSASubmissionResult: callback persisted", "submission_id", submissionId, "token", payload.Token, "mapped_status", status)

	return true, nil
}

type dsaSubmissionMeta struct {
	userID           string
	challengeID      int
	testCount        int64
	evaluationStatus int
}

func getDSASubmissionMetaForUpdate(ctx context.Context, tx pgx.Tx, submissionId string) (dsaSubmissionMeta, error) {
	var meta dsaSubmissionMeta
	err := tx.QueryRow(ctx,
		`SELECT user_id, challenge_id, test_count, evaluation_status
		 FROM dsa_submissions
		 WHERE submission_id = $1
		 FOR UPDATE`,
		submissionId,
	).Scan(&meta.userID, &meta.challengeID, &meta.testCount, &meta.evaluationStatus)
	if err != nil {
		return dsaSubmissionMeta{}, err
	}

	return meta, nil
}

func lockUserChallengeEvaluation(ctx context.Context, tx pgx.Tx, userID string, challengeID int) error {
	_, err := tx.Exec(ctx,
		`SELECT pg_advisory_xact_lock(hashtext($1 || ':' || $2::text)::bigint)`,
		userID,
		challengeID,
	)
	return err
}

func getDSASubmissionResultCounts(ctx context.Context, tx pgx.Tx, submissionId string) (int64, int64, error) {
	var submissionCount int64
	var passCount int64
	err := tx.QueryRow(ctx,
		`SELECT
			COUNT(*) AS total_count,
			COUNT(*) FILTER (WHERE status = 2) AS pass_count
		 FROM dsa_submission_results
		 WHERE submission_id = $1`,
		submissionId,
	).Scan(&submissionCount, &passCount)
	if err != nil {
		return 0, 0, err
	}

	return submissionCount, passCount, nil
}

func hasUserPassedDSAChallengeTx(ctx context.Context, tx pgx.Tx, userID string, challengeID int) (bool, error) {
	var hasPassedBefore bool
	err := tx.QueryRow(ctx,
		`SELECT EXISTS (
			SELECT 1
			FROM dsa_submissions
			WHERE user_id = $1
			  AND challenge_id = $2
			  AND evaluation_status = 2
		)`,
		userID,
		challengeID,
	).Scan(&hasPassedBefore)
	if err != nil {
		return false, err
	}

	return hasPassedBefore, nil
}

func awardLeaderboardMarksTx(ctx context.Context, tx pgx.Tx, userID string, challengeID int) error {
	var marks int
	err := tx.QueryRow(ctx,
		"SELECT marks FROM challenges WHERE id = $1",
		challengeID,
	).Scan(&marks)
	if err != nil {
		return err
	}

	if marks <= 0 {
		return nil
	}

	cmdTag, err := tx.Exec(ctx,
		"UPDATE leaderboard SET marks = marks + $2 WHERE user_id = $1",
		userID,
		marks,
	)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() > 0 {
		return nil
	}

	_, err = tx.Exec(ctx,
		"INSERT INTO leaderboard (user_id, marks) VALUES ($1, $2)",
		userID,
		marks,
	)
	return err
}

func updateDSASubmissionAggregatesTx(ctx context.Context, tx pgx.Tx, submissionId string, passCount, failCount int64, evaluationStatus int) error {
	_, err := tx.Exec(ctx,
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
	return err
}

type dsaSubmissionEvaluation struct {
	meta             dsaSubmissionMeta
	submissionCount  int64
	passCount        int64
	failCount        int64
	evaluationStatus int
	shouldFinalize   bool
}

func commitDSASubmissionTx(ctx context.Context, tx pgx.Tx, submissionId string) error {
	if err := tx.Commit(ctx); err != nil {
		logger.Log.Error("UpdateDSASubmission: commit error", "submission_id", submissionId, "error", err)
		return err
	}
	return nil
}

func prepareDSASubmissionEvaluationTx(ctx context.Context, tx pgx.Tx, submissionId string) (dsaSubmissionEvaluation, error) {
	meta, err := getDSASubmissionMetaForUpdate(ctx, tx, submissionId)
	if err != nil {
		return dsaSubmissionEvaluation{}, err
	}

	err = lockUserChallengeEvaluation(ctx, tx, meta.userID, meta.challengeID)
	if err != nil {
		return dsaSubmissionEvaluation{}, err
	}

	if meta.evaluationStatus != 1 {
		logger.Log.Info("UpdateDSASubmission: already finalized, skipping", "submission_id", submissionId, "evaluation_status", meta.evaluationStatus)
		return dsaSubmissionEvaluation{meta: meta, shouldFinalize: false}, nil
	}

	submissionCount, passCount, err := getDSASubmissionResultCounts(ctx, tx, submissionId)
	if err != nil {
		return dsaSubmissionEvaluation{}, err
	}

	logger.Log.Info("UpdateDSASubmission: submission progress", "submission_id", submissionId, "received_results", submissionCount, "expected_results", meta.testCount)
	if submissionCount < meta.testCount {
		logger.Log.Info("UpdateDSASubmission: waiting for more callbacks", "submission_id", submissionId, "remaining", meta.testCount-submissionCount)
		return dsaSubmissionEvaluation{
			meta:            meta,
			submissionCount: submissionCount,
			passCount:       passCount,
			shouldFinalize:  false,
		}, nil
	}

	failCount := submissionCount - passCount
	evaluationStatus := 3
	if passCount == submissionCount {
		evaluationStatus = 2
	}

	logger.Log.Info("UpdateDSASubmission: computed aggregate", "submission_id", submissionId, "pass_count", passCount, "fail_count", failCount, "evaluation_status", evaluationStatus)

	return dsaSubmissionEvaluation{
		meta:             meta,
		submissionCount:  submissionCount,
		passCount:        passCount,
		failCount:        failCount,
		evaluationStatus: evaluationStatus,
		shouldFinalize:   true,
	}, nil
}

func maybeAwardSubmissionMarksTx(ctx context.Context, tx pgx.Tx, submissionId string, eval dsaSubmissionEvaluation) error {
	if eval.evaluationStatus != 2 {
		return nil
	}

	hasPassedBefore, err := hasUserPassedDSAChallengeTx(ctx, tx, eval.meta.userID, eval.meta.challengeID)
	if err != nil {
		logger.Log.Error("UpdateDSASubmission: failed to check prior pass", "submission_id", submissionId, "user_id", eval.meta.userID, "challenge_id", eval.meta.challengeID, "error", err)
		return err
	}

	if hasPassedBefore {
		logger.Log.Info("UpdateDSASubmission: user already passed challenge, skipping marks", "submission_id", submissionId, "user_id", eval.meta.userID, "challenge_id", eval.meta.challengeID)
		return nil
	}

	err = awardLeaderboardMarksTx(ctx, tx, eval.meta.userID, eval.meta.challengeID)
	if err != nil {
		logger.Log.Error("UpdateDSASubmission: failed to add leaderboard marks", "submission_id", submissionId, "user_id", eval.meta.userID, "challenge_id", eval.meta.challengeID, "error", err)
		return err
	}

	logger.Log.Info("UpdateDSASubmission: leaderboard marks awarded", "submission_id", submissionId, "user_id", eval.meta.userID, "challenge_id", eval.meta.challengeID)
	return nil
}

func UpdateDSASubmission(ctx context.Context, submissionId string, payload types.TestDSAChallengeResponse) (bool, error) {
	pool := database.GetPool()
	ctx, cancel := utils.WithTimeout(ctx)
	defer cancel()
	logger.Log.Info("UpdateDSASubmission: evaluating submission", "submission_id", submissionId, "token", payload.Token)

	if payload.Token == "" {
		err := errors.New("missing callback token")
		logger.Log.Error("UpdateDSASubmission: invalid payload", "submission_id", submissionId, "error", err)
		return false, err
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		logger.Log.Error("UpdateDSASubmission: begin transaction error", "submission_id", submissionId, "error", err)
		return false, err
	}
	defer tx.Rollback(ctx)

	eval, err := prepareDSASubmissionEvaluationTx(ctx, tx, submissionId)
	if err != nil {
		logger.Log.Error("UpdateDSASubmission: failed to prepare submission evaluation", "submission_id", submissionId, "error", err)
		return false, err
	}

	if !eval.shouldFinalize {
		err = commitDSASubmissionTx(ctx, tx, submissionId)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	err = maybeAwardSubmissionMarksTx(ctx, tx, submissionId, eval)
	if err != nil {
		return false, err
	}

	err = updateDSASubmissionAggregatesTx(ctx, tx, submissionId, eval.passCount, eval.failCount, eval.evaluationStatus)
	if err != nil {
		logger.Log.Error("UpdateDSASubmission: failed to update submission aggregates", "submission_id", submissionId, "error", err)
		return false, err
	}

	err = commitDSASubmissionTx(ctx, tx, submissionId)
	if err != nil {
		return false, err
	}

	logger.Log.Info("UpdateDSASubmission: submission finalized", "submission_id", submissionId, "pass_count", eval.passCount, "fail_count", eval.failCount, "evaluation_status", eval.evaluationStatus)

	return true, nil

}
