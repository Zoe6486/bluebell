#windows 用不了makefile,用cmd/migrate/main.go代替了
# DB_URL = mysql://root:root1234@tcp(127.0.0.1:23306)/db_bluebell?multiStatements=true
# MIGRATIONS = ./migrations

# migrate-up:
# 	migrate -path $(MIGRATIONS) -database "$(DB_URL)" up

# migrate-down:
# 	migrate -path $(MIGRATIONS) -database "$(DB_URL)" down 1

# migrate-force:
# 	migrate -path $(MIGRATIONS) -database "$(DB_URL)" force $(v)

# # 敲这个就行了make migrate-up
