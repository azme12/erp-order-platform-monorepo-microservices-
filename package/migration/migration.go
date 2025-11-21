package migration

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/lib/pq"
)

type Migration struct {
	Version int
	Name    string
	UpSQL   string
	DownSQL string
}

func RunMigrations(ctx context.Context, db *sql.DB, migrationsDir string) error {
	migrations, err := loadMigrations(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	if err := createMigrationsTable(ctx, db); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	applied, err := getAppliedMigrations(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	for _, migration := range migrations {
		if _, exists := applied[migration.Version]; exists {
			continue
		}

		if err := applyMigration(ctx, db, migration); err != nil {
			return fmt.Errorf("failed to apply migration %d: %w", migration.Version, err)
		}
	}

	return nil
}

func loadMigrations(dir string) ([]Migration, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	migrations := make(map[int]*Migration)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()
		if !strings.HasSuffix(name, ".up.sql") && !strings.HasSuffix(name, ".down.sql") {
			continue
		}

		var version int
		var direction string
		if strings.HasSuffix(name, ".up.sql") {
			direction = "up"
			_, err := fmt.Sscanf(name, "%d_", &version)
			if err != nil {
				continue
			}
		} else if strings.HasSuffix(name, ".down.sql") {
			direction = "down"
			_, err := fmt.Sscanf(name, "%d_", &version)
			if err != nil {
				continue
			}
		}

		if migrations[version] == nil {
			migrations[version] = &Migration{
				Version: version,
				Name:    name,
			}
		}

		content, err := ioutil.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return nil, err
		}

		if direction == "up" {
			migrations[version].UpSQL = string(content)
			migrations[version].Name = strings.TrimSuffix(name, ".up.sql")
		} else {
			migrations[version].DownSQL = string(content)
		}
	}

	result := make([]Migration, 0, len(migrations))
	for _, m := range migrations {
		if m.UpSQL == "" {
			continue
		}
		result = append(result, *m)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Version < result[j].Version
	})

	return result, nil
}

func createMigrationsTable(ctx context.Context, db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := db.ExecContext(ctx, query)
	return err
}

func getAppliedMigrations(ctx context.Context, db *sql.DB) (map[int]bool, error) {
	rows, err := db.QueryContext(ctx, "SELECT version FROM schema_migrations ORDER BY version")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[int]bool)
	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}

	return applied, rows.Err()
}

func applyMigration(ctx context.Context, db *sql.DB, migration Migration) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, migration.UpSQL); err != nil {
		return fmt.Errorf("failed to execute migration SQL: %w", err)
	}

	if _, err := tx.ExecContext(ctx, "INSERT INTO schema_migrations (version) VALUES ($1)", migration.Version); err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	return tx.Commit()
}

func RunMigrationsFromDir(dbURL string, migrationsDir string) error {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return RunMigrations(ctx, db, migrationsDir)
}
