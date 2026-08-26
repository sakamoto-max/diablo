.PHONY: diablo watcher

diablo:
	@ cd cmd && cd server && go run main.go
watcher:
	@ cd internal && cd pkg && cd sha && go run main.go
