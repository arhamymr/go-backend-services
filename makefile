# Load environment variables from .env file
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# Define migration command
MIGRATE_CMD=migrate -path ./migrations -database $(POSTGRES_URL)

# Default target
.PHONY: all
all: migrate-up

# Run migrations up
.PHONY: migrate-up
migrate-up:
	$(MIGRATE_CMD) up

# Run migrations down
.PHONY: migrate-down
migrate-down:
	$(MIGRATE_CMD) down

# Create a new migration
.PHONY: create-migration
create-migration:
	@read -p "Enter migration name: " name; \
  migrate create -ext sql -dir ./migrations -seq $$name

# Force a specific version
.PHONY: migrate-force
migrate-force:
	@read -p "Enter version to force: " version; \
  $(MIGRATE_CMD) force $$version

# Drop everything
.PHONY: migrate-drop
migrate-drop:
	$(MIGRATE_CMD) drop


# Run the application
.PHONY: run
run:
	sh run.sh