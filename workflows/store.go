package workflows

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/IAmRiteshKoushik/gravemind/cmd"
	_ "modernc.org/sqlite"
)

var validWorkflowFields = map[string]map[string]bool{
	"achievement_workflow": {
		"job_update_hash_set": true,
		"job_check_new_badge": true,
		"processed":           true,
	},
	"bounty_workflow": {
		"job_update_redis":    true,
		"job_check_top_three": true,
		"job_update_postgres": true,
		"processed":           true,
	},
	"pull_request_workflow": {
		"job_update_pr_count":      true,
		"job_check_for_pr_badge":   true,
		"job_check_top_three":      true,
		"job_update_pr_language":   true,
		"job_check_for_issubian":   true,
		"job_update_enamoured_set": true,
		"job_check_enamoured_set":  true,
		"processed":                true,
	},
}

var localDb *sql.DB

// Setup SQLite for local persistance
func InitSQLite() error {
	db, err := sql.Open("sqlite", "gravemind.db")
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err = db.PingContext(context.Background()); err != nil {
		db.Close()
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Pool configuration
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(0)

	schema := `
CREATE TABLE IF NOT EXISTS achievement_workflow (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  stream_id TEXT NOT NULL,
	job_update_hash_set BOOLEAN DEFAULT FALSE,
	job_check_new_badge BOOLEAN DEFAULT FALSE,
  processed BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS bounty_workflow (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	stream_id TEXT NOT NULL,
	job_update_redis BOOLEAN DEFAULT FALSE,
	job_check_top_three BOOLEAN DEFAULT FALSE,
	job_update_postgres BOOLEAN DEFAULT FALSE,
	processed BOOLEAN DEFAULT VALUE
);

CREATE TABLE IF NOT EXISTS pull_request_workflow (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  stream_id TEXT NOT NULL,
	job_update_pr_count BOOLEAN DEFAULT FALSE,
	job_check_for_pr_badge BOOLEAN DEFAULT FALSE,
	job_check_top_three BOOLEAN DEFAULT FALSE,
	job_update_pr_language BOOLEAN DEFAULT FALSE,
	job_check_for_issubian BOOLEAN DEFAULT FALSE,
	job_update_enamoured_set BOOLEAN DEFAULT FALSE,
	job_check_enamoured_set BOOLEAN DEFAULT FALSE,
  processed BOOLEAN DEFAULT FALSE
);`

	_, err = db.ExecContext(context.Background(), schema)
	if err != nil {
		db.Close()
		return fmt.Errorf("failed to create tables: %w", err)
	}

	localDb = db
	return nil
}

func CloseDb() {
	if localDb != nil {
		if err := localDb.Close(); err != nil {
			cmd.Log.Error("Error closing database connection: %v", err)
		} else {
			cmd.Log.Info("Database connection closed.")
		}
	}
}

func entry(tableName, streamId string) {
	if localDb == nil {
		cmd.Log.Error("Database is not initialized.", fmt.Errorf("Call InitSQLite first."))
		return
	}
	validTables := map[string]bool{
		"achievement_workflow":  true,
		"bounty_workflow":       true,
		"pull_request_workflow": true,
	}

	if !validTables[tableName] {
		cmd.Log.Error(fmt.Sprintf("Invalid table name provided: %s", tableName), fmt.Errorf("Table does not exist"))
		return
	}

	query := fmt.Sprintf(`INSERT INTO %s (stream_id) VALUES (?)`, tableName)
	_, err := localDb.ExecContext(context.Background(), query, streamId)
	if err != nil {
		cmd.Log.Error(fmt.Sprintf("Failed to insert entry into %s for stream_id %s", tableName, streamId), err)
	} else {
		cmd.Log.Info(fmt.Sprintf("Entry for stream_id %s inserted into %s successfully.", streamId, tableName))
	}
}

// --- Wrapper Functions for Entries ---
func AchivementEntry(entryId string) {
	entry("achievement_workflow", entryId)
}

func PullRequestEntry(entryId string) {
	entry("pull_request_workflow", entryId)
}

func BountyEntry(entryId string) {
	entry("bounty_workflow", entryId)
}

func complete(tableName, fieldName, streamId string) {
	if localDb == nil {
		cmd.Log.Error("Database is not initialized. Call InitSQLite first.",
			fmt.Errorf("db not found"))
		return
	}

	// Validate table name and table field name using a map
	fields, tableExists := validWorkflowFields[tableName]
	if !tableExists {
		cmd.Log.Error(fmt.Sprintf("Invalid table name provided for completion: %s",
			tableName), fmt.Errorf("table does not exist"))
		return
	}
	if !fields[fieldName] {
		cmd.Log.Error(fmt.Sprintf("Invalid field name '%s' for table '%s'.",
			fieldName, tableName), fmt.Errorf("field does not exist"))
		return
	}

	query := fmt.Sprintf(`UPDATE %s SET %s = TRUE WHERE stream_id = ?`, tableName,
		fieldName)
	result, err := localDb.ExecContext(context.Background(), query, streamId)
	if err != nil {
		cmd.Log.Error(fmt.Sprintf("Failed to update %s.%s for stream_id %s",
			tableName, fieldName, streamId), err)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		cmd.Log.Error(fmt.Sprintf("Failed to get rows affected for update %s.%s for stream_id %s",
			tableName, fieldName, streamId), err)
		return
	}

	if rowsAffected == 0 {
		cmd.Log.Warn(fmt.Sprintf("No rows updated for %s.%s with stream_id %s. Entry might not exist.",
			tableName, fieldName, streamId))
	} else {
		cmd.Log.Info(fmt.Sprintf("Successfully updated %s.%s for stream_id %s.",
			tableName, fieldName, streamId))
	}
}

// --- Wrapper Functions for Achievement Workflow ---
func JobUpdateHashSet(streamId string) {
	complete("achievement_workflow", "job_update_hash_set", streamId)
}

func JobCheckNewBadge(streamId string) {
	complete("achievement_workflow", "job_check_new_badge", streamId)
}

func ProcessedAchievement(streamId string) {
	complete("achievement_workflow", "processed", streamId)
}

// --- Wrapper Functions for Bounty Workflow ---
func JobUpdateRedis(streamId string) {
	complete("bounty_workflow", "job_update_redis", streamId)
}

func JobCheckTopThreeBounty(streamId string) {
	complete("bounty_workflow", "job_check_top_three", streamId)
}

func JobUpdatePostgres(streamId string) {
	complete("bounty_workflow", "job_update_postgres", streamId)
}

func ProcessedBounty(streamId string) {
	complete("bounty_workflow", "processed", streamId)
}

// --- Wrapper Functions for Pull Request Workflow ---
func JobUpdatePrCount(streamId string) {
	complete("pull_request_workflow", "job_update_pr_count", streamId)
}

func JobCheckForPrBadge(streamId string) {
	complete("pull_request_workflow", "job_check_for_pr_badge", streamId)
}

func JobCheckTopThreePR(streamId string) {
	complete("pull_request_workflow", "job_check_top_three", streamId)
}

func JobUpdatePrLanguage(streamId string) {
	complete("pull_request_workflow", "job_update_pr_language", streamId)
}

func JobCheckForIssubian(streamId string) {
	complete("pull_request_workflow", "job_check_for_issubian", streamId)
}

func JobUpdateEnamouredSet(streamId string) {
	complete("pull_request_workflow", "job_update_enamoured_set", streamId)
}

func JobCheckEnamouredSet(streamId string) {
	complete("pull_request_workflow", "job_check_enamoured_set", streamId)
}

func ProcessedPullRequest(streamId string) {
	complete("pull_request_workflow", "processed", streamId)
}

func ReadLastEntry(table string) (string, error) {
	if localDb == nil {
		return "", fmt.Errorf("database is not initialized. Call InitSQLite first")
	}

	var tableName string
	switch table {
	case "bounty":
		tableName = "bounty_workflow"
	case "pull_request":
		tableName = "pull_request_workflow"
	case "achievement":
		tableName = "achievement_workflow"
	default:
		return "", fmt.Errorf("failed to read entry: unknown table identifier '%s'", table)
	}

	var streamID string
	query := fmt.Sprintf(`SELECT stream_id FROM %s ORDER BY id DESC LIMIT 1`, tableName)

	row := localDb.QueryRowContext(context.Background(), query)
	err := row.Scan(&streamID)

	if err != nil {
		if err == sql.ErrNoRows {
			// This means the table is empty
			cmd.Log.Info(fmt.Sprintf("No entries found in table '%s'.", tableName))
			return "", err
		}
		cmd.Log.Error(fmt.Sprintf("Failed to read last entry from '%s'", tableName), err)
		return "", err
	}

	cmd.Log.Info(fmt.Sprintf("Successfully read last stream_id '%s' from table '%s'.", streamID, tableName))
	return streamID, nil
}
