package main

import (
	"crypto/sha256"
	"fmt"
	"os"

	myErrs "github.com/sakamoto-max/diablo/internal/pkg/myerrors"
)

func GetShaIndex(data []byte) [32]byte {
	h := sha256.Sum256(data)
	return h
}

func ReadFile(filePath string) ([]byte, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		// todo : handle different file errors
		return nil, myErrs.Wrap(fmt.Errorf("failed to read file of path %v : %w", filePath, err), myErrs.Internal)
	}

	return bytes, nil
}

func main() {

	// str := `diablo/internal/config/config.go`

	// // d := strings.Split(str, `\`)
	// fmt.Println(str)

	// err := filepath.WalkDir("../../../../diablo",
	// 	func(path string, d fs.DirEntry, err error) error {
	// 		if err != nil {
	// 			return fmt.Errorf("error from call back func :%w", err)
	// 		}

	// 		t := d.Type()

	// 		var typ string
	// 		if t == fs.ModeDir {
	// 			typ = "dir"
	// 		} else {
	// 			typ = "file"
	// 		}

	// 		path = strings.TrimPrefix(path, "../../../../")
	// 		path = strings.TrimPrefix(path, `..\..\..\..\`)

	// 		loc := domain.Dir{Path: path, Type: typ}

	// 		fmt.Println(loc)

	// 		return nil
	// 	})
	// if err != nil {
	// 	panic(fmt.Errorf("failed to read the directrory : %w", err))
	// }

}

// {diablo dir}
// {.env file}
// {cmd dir}
// {cmd\server dir}
// {cmd\server\main.go file}
// {go.mod file}
// {go.sum file}
// {internal dir}
// {internal\app dir}
// {internal\app\app.go file}
// {internal\config dir}
// {internal\config\config.go file}
// {internal\database dir}
// {internal\database\migrations dir}
// {internal\database\migrations\001_users.sql file}
// {internal\database\migrator.go file}
// {internal\database\pool.go file}
// {internal\domain dir}
// {internal\domain\dir.go file}
// {internal\domain\user.go file}
// {internal\dto dir}
// {internal\dto\client.go file}
// {internal\dto\user.go file}
// {internal\dto\validations.go file}
// {internal\env dir}
// {internal\env\env.go file}
// {internal\handlers dir}
// {internal\handlers\handlers.go file}
// {internal\handlers\users.go file}
// {internal\middleware dir}
// {internal\middleware\auth.go file}
// {internal\middleware\middleawares.go file}
// {internal\pkg dir}
// {internal\pkg\myerrors dir}
// {internal\pkg\myerrors\grpc.go file}
// {internal\pkg\myerrors\http.go file}
// {internal\pkg\myerrors\wraper.go file}
// {internal\pkg\password dir}
// {internal\pkg\password\password.go file}
// {internal\pkg\sha dir}
// {internal\pkg\sha\main.go file}
// {internal\pkg\token dir}
// {internal\pkg\token\jwt.go file}
// {internal\repository dir}
// {internal\repository\repository.go file}
// {internal\repository\user.go file}
// {internal\router dir}
// {internal\router\router.go file}
// {internal\services dir}
// {internal\services\services.go file}
// {internal\services\users.go file}
// {makefile file}
// {mvp.txt file}
// {tern.conf file}
// {todo.txt file}
