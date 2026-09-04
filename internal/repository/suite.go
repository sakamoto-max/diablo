package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	// "github.com/jackc/pgx/v5"

	"github.com/sakamoto-max/diablo/internal/domain"
	"github.com/sakamoto-max/diablo/internal/dto"
)

type Suite struct {
	db *sql.DB
}

type SuiteIface interface {
	New(ctx context.Context, suite *dto.NewSuiteReq) error
	Sync(ctx context.Context, data *dto.EventsReq) error
	Get(ctx context.Context, input dto.Suite) (dto.Suite, error)
	GetSuitesAndEvents(ctx context.Context, userIp string) ([]domain.LastSyncedData, error)
	GetFileContents(ctx context.Context, suiteId string, path string) ([]byte, error)
}

func (s *Suite) New(ctx context.Context, suite *dto.NewSuiteReq) error {

	trnx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to create transaction")
	}

	defer trnx.Rollback()

	for _, suite := range suite.Suites {
		query := `
			INSERT INTO SUITES(NAME, LAST_CHANGED)
			VALUES(
				@suiteName,
				@lastChanged
			)

			RETURNING ID
	 	`

		log.Println("suite Name is : %v", suite.Name)

		var suiteId int

		

		err = trnx.QueryRowContext(ctx, query,
			sql.Named("suiteName", suite.Name),
			sql.Named("lastChanged", time.Now()),
		).Scan(&suiteId)
		if err != nil {
			return fmt.Errorf("failed to insert suite : %w", err)
		}

		fmt.Println("suite id", suiteId)


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

			_, err := trnx.ExecContext(ctx, query,
				sql.Named("suiteId", suiteId),
				sql.Named("path", file.Path),
				sql.Named("data", file.Contents),
				sql.Named("fileType", file.FileType),
			)
			if err != nil {
				return fmt.Errorf("failed to insert file of path %v : %w", file.Path, err)
			}
		}
	}

	err = trnx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction : %w", err)
	}

	return nil
}

func (s *Suite) Sync(ctx context.Context, data *dto.EventsReq) error {

	trnx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to create transaction")
	}

	defer trnx.Rollback()

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

		err = trnx.QueryRowContext(ctx, query,
			sql.Named("suiteName", event.SuiteName),
		).Scan(&suiteId)
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

		_, err = trnx.ExecContext(ctx, query,
			sql.Named("suiteId", suiteId),
			sql.Named("path", event.Path),
			sql.Named("data", event.Contents),
			sql.Named("fileType", event.FileType),
			sql.Named("newPath", event.RenamedTo),
			sql.Named("oldPath", event.Path),
		)
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

		_, err := trnx.ExecContext(ctx, query,
			sql.Named("suiteId", suiteId),
			sql.Named("eventName", event.Event),
			sql.Named("path", event.Path),
		)

		if err != nil {
			return fmt.Errorf("failed to insert event : %w", err)
		}
	}

	err = trnx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction : %w", err)
	}

	return nil
}

func (s *Suite) Get(ctx context.Context, input dto.Suite) (dto.Suite, error) {

	fmt.Println(input.Name)

	trnx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return dto.Suite{}, fmt.Errorf("failed to create transaction")
	}

	defer trnx.Rollback()

	query := `
		SELECT
			ID
		FROM 
			SUITES
		WHERE
			NAME = @suiteName
	`

	var suiteId string

	err = trnx.QueryRowContext(ctx, query,
		sql.Named("suiteName", input.Name),
	).Scan(&suiteId)
	if err != nil {
		return dto.Suite{}, fmt.Errorf("failed to get suite id : %w", err)
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

	rows, err := trnx.QueryContext(ctx, query,
		sql.Named("suiteId", suiteId),
	)
	if err != nil {
		return dto.Suite{}, fmt.Errorf("failed to get rows : %w", err)
	}

	var allFiles []dto.Event

	var path string
	var data []byte
	var fileType string

	for rows.Next() {
		err := rows.Scan(&path, &data, &fileType)
		if err != nil {
			return dto.Suite{}, fmt.Errorf("failed to scan rows : %w", err)
		}

		allFiles = append(allFiles, dto.Event{
			SuiteName: input.Name,
			Path:      path,
			Contents:  data,
			FileType:  fileType,
		})
	}

	if len(allFiles) == 0 {
		return dto.Suite{}, nil
	}

	// todo : uncomment this later

	// query = `
	// 	INSERT INTO USER_IPS(
	// 		SUITE_ID,
	// 		IP,
	// 		LAST_SYNCED
	// 	)

	// 	VALUES(
	// 		@suiteId,
	// 		@ip,
	// 		@lastSynced
	// 	)
	
	// `

	// _, err = trnx.ExecContext(ctx, query,
	// 	sql.Named("suiteId", suiteId),
	// 	sql.Named("ip", input.Ip),
	// 	sql.Named("lastSynced", time.Now()),
	// )
	// if err != nil {
	// 	return dto.Suite{}, fmt.Errorf("failed to insert user ip : %w", err)
	// }

	err = trnx.Commit()
	if err != nil {
		return dto.Suite{}, fmt.Errorf("failed to commit transaction : %w", err)
	}

	return dto.Suite{
		Name:  input.Name,
		Files: allFiles,
	}, nil
}

func (s *Suite) GetSuitesAndEvents(ctx context.Context, userIp string) ([]domain.LastSyncedData, error) {
	trnx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin the transaction : %w", err)
	}

	defer trnx.Rollback()

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

	rows, err := trnx.QueryContext(ctx, query, sql.Named("ip", userIp))
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

		eventRows, err := trnx.QueryContext(ctx, query,
			sql.Named("suiteId", suite.SuiteId),
			sql.Named("lastSynced", suite.LastSynced),
		)

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

	err = trnx.Commit()
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

	err := s.db.QueryRowContext(ctx, query,
		sql.Named("suiteId", suiteId),
		sql.Named("path", path),
	).Scan(&data)
	if err != nil {
		return nil, fmt.Errorf("failed to get file contents : %w", err)
	}

	return data, nil
}
