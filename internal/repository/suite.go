package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sakamoto-max/diablo/internal/domain"
	"github.com/sakamoto-max/diablo/internal/dto"
)

type Suite struct {
	db *pgxpool.Pool
}

type SuiteIface interface {
	New(ctx context.Context, suite *dto.NewSuiteReq) error
	Sync(ctx context.Context, data *dto.EventsReq) error
	Get(ctx context.Context, input dto.Event) (*dto.Suite, error)
	GetSuitesAndEvents(ctx context.Context, userIp string) ([]domain.LastSyncedData, error)
	GetFileContents(ctx context.Context, suiteId string, path string) ([]byte, error)
}

func (s *Suite) New(ctx context.Context, suite *dto.NewSuiteReq) error {

	trnx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to create transaction")
	}

	defer trnx.Rollback(ctx)

	for _, suite := range suite.Suites {
		query := `
			INSERT INTO SUITES(NAME, LAST_CHANGED)
			VALUES(
				@suiteName,
				@lastChanged
			)
			
			RETURNING ID
		`

		var suiteUUId string

		err = trnx.QueryRow(ctx, query, pgx.NamedArgs{
			"suiteName":   suite.Name,
			"lastChanged": time.Now(),
		}).Scan(&suiteUUId)
		if err != nil {
			return fmt.Errorf("failed to insert suite : %w", err)
		}

		for _, file := range suite.Files {
			query = `
				INSERT INTO FILES(
					SUITE_ID,
					PATH,
					DATA,
					FILE_TYPE
				)
				VALUES(
					@suiteId,
					@path,
					@data,
					@fileType
				)
			`

			_, err = trnx.Exec(ctx, query, pgx.NamedArgs{
				"suiteId":  suiteUUId,
				"path":     file.Path,
				"data":     file.Contents,
				"fileType": file.FileType,
			})
			if err != nil {
				return fmt.Errorf("failed to insert file of path %v : %w", file.Path, err)
			}
		}
	}

	err = trnx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction : %w", err)
	}

	return nil
}

func (s *Suite) Sync(ctx context.Context, data *dto.EventsReq) error {

	trnx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to create transaction")
	}

	defer trnx.Rollback(ctx)

	for _, event := range data.Events {
		query := `
			SELECT
				ID
			FROM
				SUITES
			WHERE
				NAME = @suiteName
		`
		var suiteId string

		err = trnx.QueryRow(ctx, query, pgx.NamedArgs{
			"suiteName": event.SuiteName,
		}).Scan(&suiteId)
		if err != nil {
			return fmt.Errorf("failed to get suite id : %w", err)
		}
		switch event.Event {
		case "written":

			query = `
					UPDATE
						FILES
					SET 
						DATA = @data
					WHERE
						SUITE_ID = @suiteId
						AND PATH = @path
				`

		case "created":

			query = `
					INSERT INTO FILES(
						SUITE_ID,
						PATH,
						DATA,
						FILE_TYPE
					)
					VALUES(
						@suiteId,
						@path,
						@data,
						@fileType
					)
				`

		case "deleted":

			query = `
					DELETE FROM FILES
					WHERE
						SUITE_ID = @suiteId
						AND PATH = @path
				`
		case "renamed":

			query = `
					UPDATE FILES
					SET PATH = @newPath
					WHERE
						SUITE_ID = @suiteId
						AND PATH = @oldPath
				`
		}

		_, err = trnx.Exec(ctx, query, pgx.NamedArgs{
			"suiteId":  suiteId,
			"path":     event.Path,
			"data":     event.Contents,
			"fileType": event.FileType,
			"newPath":  event.RenamedTo,
			"oldPath":  event.Path,
		})
		if err != nil {
			return fmt.Errorf("failed to update file : %w", err)
		}

		query = `
			INSERT INTO EVENTS(
				SUITE_ID,
				NAME,
				PATH
			)

			VALUES(
				@suiteId,
				@eventName,
				@path
			)
		`

		_, err := trnx.Exec(ctx, query, pgx.NamedArgs{
			"suiteId":   suiteId,
			"eventName": event.Event,
			"path":      event.Path,
		})

		if err != nil {
			return fmt.Errorf("failed to insert event : %w", err)
		}
	}

	err = trnx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction : %w", err)
	}

	return nil
}

func (s *Suite) Get(ctx context.Context, input dto.Event) (*dto.Suite, error) {

	trnx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction")
	}

	defer trnx.Rollback(ctx)

	query := `
		SELECT
			ID
		FROM 
			SUITES
		WHERE
			NAME = @suiteName
	`

	var suiteId string

	err = trnx.QueryRow(ctx, query, pgx.NamedArgs{
		"suiteName": input.SuiteName,
	}).Scan(&suiteId)
	if err != nil {
		return nil, fmt.Errorf("failed to get suite id : %w", err)
	}

	query = `
		SELECT
			PATH,
			DATA,
			FILE_TYPE
		FROM 
			FILES
		WHERE 
			SUITE_ID = @suiteId
	`

	rows, err := trnx.Query(ctx, query, pgx.NamedArgs{
		"suiteId": suiteId,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get rows : %w", err)
	}

	var allFiles []dto.Event

	var path string
	var data []byte
	var fileType string

	for rows.Next() {
		err := rows.Scan(&path, &data, &fileType)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rows : %w", err)
		}

		allFiles = append(allFiles, dto.Event{
			Path:     path,
			Contents: data,
			FileType: fileType,
		})
	}

	if len(allFiles) == 0 {
		// todo
	}

	query = `
		INSERT INTO USER_IPS(
			SUITE_ID,
			IP,
			LAST_SYNCED
		)

		VALUES(
			@suiteId,
			@ip,
			@lastSynced
		)
	
	`

	_, err = trnx.Exec(ctx, query, pgx.NamedArgs{
		"suiteId":    suiteId,
		"ip":         input.Ip,
		"lastSynced": time.Now(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to insert user ip : %w", err)
	}

	err = trnx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction : %w", err)
	}

	return &dto.Suite{
		Name:  input.SuiteName,
		Files: allFiles,
	}, nil
}

func (s *Suite) GetSuitesAndEvents(ctx context.Context, userIp string) ([]domain.LastSyncedData, error) {
	trnx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin the transaction : %w", err)
	}

	defer trnx.Rollback(ctx)

	// get all the suites of the user
	query := `
		SELECT 
			USER_IPS.SUITE_ID,
			USER_IPS.LAST_SYNCED,
			SUITES.NAME
		FROM 
			USER_IPS
		INNER JOIN 
			SUITES
		ON 
			USER_IPS.SUITE_ID = SUITES.ID
		WHERE 
			IP = @ip
	`

	rows, err := trnx.Query(ctx, query, pgx.NamedArgs{
		"ip": userIp,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get rows : %w", err)
	}

	var allSuites []domain.LastSyncedData

	var suiteId string
	var lastSynced time.Time
	var suiteName string

	for rows.Next() {
		err := rows.Scan(&suiteId, &lastSynced, &suiteName)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rows : %w", err)
		}

		allSuites = append(allSuites, domain.LastSyncedData{
			SuiteId:    suiteId,
			LastSynced: lastSynced,
			SuiteName:  suiteName,
		})
	}

	for _, suite := range allSuites {

		query = `
			SELECT
				NAME,
				PATH,
				RENAMED_TO
			FROM 
				EVENTS
			WHERE 
				SUITE_ID = @suiteId
				AND CHANGED_AT > @lastSynced
		`

		eventRows, err := trnx.Query(ctx, query, pgx.NamedArgs{
			"suiteId":    suite.SuiteId,
			"lastSynced": suite.LastSynced,
		})

		if err != nil {
			return nil, fmt.Errorf("failed to get event rows : %w", err)
		}

		var eventName string
		var path string
		var renamedTo string

		for eventRows.Next() {
			err := rows.Scan(&eventName, &path, &renamedTo)
			if err != nil {
				return nil, fmt.Errorf("failed to scan event rows : %w", err)
			}

			suite.Events = append(suite.Events, dto.Event{
				SuiteName: suite.SuiteName,
				Path:      path,
				Event:     eventName,
				RenamedTo: renamedTo,
			})

		}
	}

	err = trnx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction : %w", err)
	}

	return allSuites, nil
}

func (s *Suite) GetFileContents(ctx context.Context, suiteId string, path string) ([]byte, error) {
	query := `
		SELECT
			DATA
		FROM 
			FILES
		WHERE 
			SUITE_ID = @suiteId
			AND PATH = @path
	`

	var data []byte

	err := s.db.QueryRow(ctx, query, pgx.NamedArgs{
		"suiteId": suiteId,
		"path":    path,
	}).Scan(&data)
	if err != nil {
		return nil, fmt.Errorf("failed to get file contents : %w", err)
	}

	return data, nil
}
