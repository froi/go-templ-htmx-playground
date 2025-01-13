
GREEN='\033[0;32m'
NC='\033[0m'


.PHONY: templ
templ:
	@echo "${GREEN}Generating templ code...${NC}"
	templ generate -watch -proxy=http://dev.localhost:8080 ./ui/

.PHONY: tailwindcss-dev
tailwindcss-dev:
	@echo "${GREEN}Generating and watching Tailwind CSS code...${NC}"
	npx tailwindcss -c ./ui/tailwind.config.js -i ./ui/main.css -o ./ui/static/main.css --watch
.PHONY: tailwindcss
tailwindcss:
	@echo "${GREEN}Generating Tailwind CSS code...${NC}"
	npx tailwindcss -c ./ui/tailwind.config.js -i ./ui/main.css -o ./ui/static/main.css --minify

.PHONY: air
air:
	@air -c .air.toml

.PHONY: sqlc
sqlc:
	@echo "${GREEN}Generating SQLC code...${NC}"
	sqlc generate -f ./sqlc.yaml


.PHONY: clean
clean:
	@echo "Killing Air process"
	ps | rg "ai[r]" | awk '{print $1}' | xargs kill

codegen:
	@echo "${GREEN}Generating templ code...${NC}"
	templ generate ./ui/
	# @echo "${GREEN}Generating SQLC code...${NC}"
	# sqlc generate -f ./sqlc.yaml
	@echo "${GREEN}Generating Tailwind CSS code...${NC}"
	npx tailwindcss -c ./ui/tailwind.config.js -i ./ui/main.css -o ./ui/static/main.css --minify

migrate_up:
	@echo "${GREEN}Migrating up...${NC}"
	go run cmd/migrate/main.go upOne

build_ui: codegen
	@echo "${GREEN}Building UI...${NC}"
	go build -o bin/ui/ui cmd/ui/main.go
	cp -r ./ui/static ./bin/ui/

run_ui: codegen
	@echo "${GREEN}Running UI...${NC}"
	go run cmd/ui/main.go
