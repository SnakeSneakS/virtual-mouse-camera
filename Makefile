# USAGE: make help

.DEFAULT_GOAL := help

ADDR=127.0.0.1
PORT=8080
CAMERA_INDEX=1

.PHONY: run-hand-tracker 
run-hand-tracker: ## run hand-tracker ## make run-hand-tracker ADDR=localhost PORT=8080 CAMERA_INDEX=1
	@echo handling socket on $(ADDR):$(PORT) using camera $(CAMERA_INDEX)
	cd hand-tracker && python main.py $(ADDR) $(PORT) $(CAMERA_INDEX)

.PHONY: run-mouse-controller
run-mouse-controller: ## run mouse-controller ## make run-mouse-controller ADDR=localhost PORT=8080 
	@echo handling socket on $(ADDR):$(PORT) 
	cd mouse-controller && go run main.go $(ADDR) $(PORT)

#.PHONY: run-all
#run:  ## run all ## make run-all ADDR=localhost PORT=8080
#	@echo handling socket on $(ADDR):$(PORT) 
#	cd hand-tracker && python main.py $(ADDR) $(PORT) & \
#	cd mouse-controller && go run main.go $(ADDR) $(PORT)

# thanks to https://ktrysmt.github.io/blog/write-useful-help-command-by-shell/
.PHONY: help
help: ## show help ## make help 
	@echo "--- Makefile Help ---"
	@echo ""
	@echo "Usage: make SUB_COMMAND argument_name=argument_value"
	@echo ""
	@echo "Command list:"
	@printf "\033[36m%-30s\033[0m %-30s %s\n" "[Sub command]" "[Description]" "[Example]"
	@grep -E '^[/a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | perl -pe 's%^([/a-zA-Z_-]+):.*?(##)%$$1 $$2%' | awk -F " *?## *?" '{printf "\033[36m%-30s\033[0m %-30s %-30s\n", $$1, $$2, $$3}'