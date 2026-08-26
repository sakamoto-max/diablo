package services

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sakamoto-max/diablo/internal/domain"
	"github.com/sakamoto-max/diablo/internal/dto"
	"github.com/sakamoto-max/diablo/internal/repository"
)

type Synchronizer struct {
	db *repository.Db
}

const (
	InitialFileSystem = `C:\Users\clikithh\Desktop\GO_CODE\PROJECTS\diablod_target_data`
)

// func (s *Synchronizer) New(ctx context.Context, input *dto.AllFiles) error {

// 	log.Println("req received in service")

// 	for _, file := range input.Files {

// 		parentPath, fileName := getParentAndFileName(file.Path)

// 		err := makeDir(parentPath)
// 		if err != nil {
// 			return fmt.Errorf("failed to make dir : %w", err)
// 		}

// 		parentPath = filepath.Join(InitialFileSystem, parentPath)
// 		// make child
// 		filePath := filepath.Join(parentPath, fileName)
// 		osFile, err := os.OpenFile(filePath, os.O_CREATE, os.ModePerm)
// 		if err != nil {
// 			return fmt.Errorf("failed to create file %v : %w", file.Path, err)
// 		}

// 		// write to the file
// 		_, err = osFile.Write(file.Data)
// 		if err != nil {
// 			return fmt.Errorf("failed to write to file %v : %w", file.Path, err)
// 		}

// 		osFile.Close()
// 	}
// 	log.Println("req over in service")

// 	return nil
// }

func (s *Synchronizer) New(ctx context.Context, input *dto.NewSuiteReq) error {

	err := s.db.Suite.New(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to create a new suite : %w", err)
	}

	return nil
}

func (s *Synchronizer) Ping(ctx context.Context, input dto.UserIp) ([]dto.Event, error) {
	suitesAndEvents, err := s.db.Suite.GetSuitesAndEvents(ctx, input.IP)
	if err != nil {
		return nil, fmt.Errorf("failed to get suites and events : %w", err)
	}

	suitesAndEvents = Ledge(suitesAndEvents)

	var allEvents []dto.Event

	for _, suite := range suitesAndEvents {
		for _, event := range suite.Events {

			if event.Event == "written" {
				contents, err := s.db.Suite.GetFileContents(ctx, suite.SuiteId, event.Path)
				if err != nil {
					return nil, fmt.Errorf("failed to get file contents : %w", err)
				}

				event.Contents = contents
			}

			allEvents = append(allEvents, dto.Event{
				SuiteName: suite.SuiteName,
				Path:      event.Path,
				Event:     event.Event,
				Contents:  event.Contents,
				FileType:  event.FileType,
				RenamedTo: event.RenamedTo,
			})
		}

	}

	return allEvents, nil
}

func Ledge(suitesAndEvents []domain.LastSyncedData) []domain.LastSyncedData {

	for i := range suitesAndEvents {

		allEvents := make(map[string]dto.Event)

		for _, event := range suitesAndEvents[i].Events {

			if event.Event == "deleted" {
				_, ok := allEvents[fmt.Sprintf("%v-%v", event.Path, event.Event)]
				if ok {
					delete(allEvents, fmt.Sprintf("%v-%v", event.Path, "created"))
					delete(allEvents, fmt.Sprintf("%v-%v", event.Path, "written"))
				} else {
					allEvents[fmt.Sprintf("%v-%v", event.Path, event.Event)] = event
				}
			}

			allEvents[fmt.Sprintf("%v-%v", event.Path, event.Event)] = event
		}

		var newEvents []dto.Event

		for _, event := range allEvents {
			newEvents = append(newEvents, event)
		}

		suitesAndEvents[i].Events = newEvents
	}

	return suitesAndEvents
}

// main.txt create
// main.txt write
// main.txt delete
// user  create
// repository create
// repository.go create
// repository.go write
// suite.go create
// suite.go delete
// main.txt create

// func (s *Synchronizer) Sync(ctx context.Context, input *dto.FileSystem) error {

// 	switch input.Event {
// 	case "written":

// 		err := writeToFile(input.Path, input.Contents)
// 		if err != nil {
// 			return fmt.Errorf("failed to write to file %v : %w", input.Path, err)
// 		}

// 	case "created":

// 		switch input.IsDir {
// 		case true:
// 			err := makeDir(input.Path)
// 			if err != nil {
// 				return fmt.Errorf("failed to create dir %v : %w", input.Path, err)
// 			}
// 		case false:
// 			err := makeFile(input.Path)
// 			if err != nil {
// 				return fmt.Errorf("failed to create file %v : %w", input.Path, err)
// 			}
// 		}

// 	case "deleted":

// 		err := deleteFile(input.Path)
// 		if err != nil {
// 			return fmt.Errorf("failed to delete file %v : %w", input.Path, err)
// 		}

// 	case "renamed":
// 		err := renameFile(input.Path, input.RenamedTo)
// 		if err != nil {
// 			return fmt.Errorf("failed to rename file %v : %w", input.Path, err)
// 		}
// 	}

// 	return nil
// }

func (s *Synchronizer) Sync(ctx context.Context, input *dto.EventsReq) error {
	return s.db.Suite.Sync(ctx, input)
}

func (s *Synchronizer) GetSuite(ctx context.Context, input dto.Event) (*dto.Suite, error) {
	return s.db.Suite.Get(ctx, input)
}

func getParentAndFileName(path string) (string, string) {
	allFiles := strings.Split(path, `\`)

	parent := strings.Join(allFiles[:len(allFiles)-1], `\`)
	fileName := allFiles[len(allFiles)-1]

	return parent, fileName
}

func makeDir(path string) error {

	path = filepath.Join(InitialFileSystem, path)

	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create file %v : %w", path, err)
	}

	return nil
}

func makeFile(path string) error {
	parentPath, fileName := getParentAndFileName(path)

	// make parent
	parentPath = filepath.Join(InitialFileSystem, parentPath)
	err := os.MkdirAll(parentPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create file %v : %w", path, err)
	}

	// make child
	filePath := filepath.Join(parentPath, fileName)
	osFile, err := os.OpenFile(filePath, os.O_CREATE, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create file %v : %w", path, err)
	}

	osFile.Close()

	return nil
}

func writeToFile(path string, contents []byte) error {
	path = filepath.Join(InitialFileSystem, path)

	err := os.WriteFile(path, contents, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write to file %v : %w", path, err)
	}

	return nil
}

func deleteFile(path string) error {

	path = filepath.Join(InitialFileSystem, path)
	err := os.Remove(path)
	if err != nil {
		return fmt.Errorf("failed to delete file %v : %w", path, err)
	}

	return nil
}

func renameFile(path string, newPath string) error {
	oldPath := filepath.Join(InitialFileSystem, path)
	newPath = filepath.Join(InitialFileSystem, newPath)

	err := os.Rename(oldPath, newPath)
	if err != nil {
		return fmt.Errorf("failed to rename file %v : %w", path, err)
	}

	return nil
}
