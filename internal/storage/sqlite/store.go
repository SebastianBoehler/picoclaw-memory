package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"time"

	_ "modernc.org/sqlite"

	"picoclaw-memory/internal/memory"
)

type Store struct {
	db *sql.DB
}

func Open(path string) (*Store, error) {
	dsn := fmt.Sprintf(
		"file:%s?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=foreign_keys(1)",
		url.PathEscape(path),
	)

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	store := &Store{db: db}
	if err := store.initSchema(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return store, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) Save(ctx context.Context, item memory.Memory) error {
	const query = `
		INSERT INTO memories (
			id,
			container_tag,
			content,
			source,
			created_at,
			expires_at
		) VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.ExecContext(
		ctx,
		query,
		item.ID,
		item.ContainerTag,
		item.Content,
		item.Source,
		item.CreatedAt.UTC().Format(time.RFC3339Nano),
		nullableTime(item.ExpiresAt),
	)
	if err != nil {
		return fmt.Errorf("insert memory: %w", err)
	}

	return nil
}

func (s *Store) Recall(ctx context.Context, req memory.RecallRequest) ([]memory.RecallResult, error) {
	const query = `
		SELECT
			id,
			container_tag,
			content,
			source,
			created_at,
			expires_at,
			CASE
				WHEN lower(content) = lower(?) THEN 100
				WHEN instr(lower(content), lower(?)) > 0 THEN 60
				ELSE 0
			END AS score
		FROM memories
		WHERE container_tag = ?
			AND (expires_at IS NULL OR expires_at > ?)
			AND instr(lower(content), lower(?)) > 0
		ORDER BY score DESC, created_at DESC
		LIMIT ?
	`

	rows, err := s.db.QueryContext(
		ctx,
		query,
		req.Query,
		req.Query,
		req.ContainerTag,
		time.Now().UTC().Format(time.RFC3339Nano),
		req.Query,
		req.Limit,
	)
	if err != nil {
		return nil, fmt.Errorf("recall memories: %w", err)
	}
	defer rows.Close()

	results := make([]memory.RecallResult, 0, req.Limit)
	for rows.Next() {
		item, score, err := scanRecallRow(rows)
		if err != nil {
			return nil, err
		}

		results = append(results, memory.RecallResult{
			Memory: item,
			Score:  score,
			Reason: "lexical_match",
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate recall rows: %w", err)
	}

	return results, nil
}

func (s *Store) ListRecent(ctx context.Context, containerTag string, limit int) ([]memory.Memory, error) {
	const query = `
		SELECT
			id,
			container_tag,
			content,
			source,
			created_at,
			expires_at
		FROM memories
		WHERE container_tag = ?
			AND (expires_at IS NULL OR expires_at > ?)
		ORDER BY created_at DESC
		LIMIT ?
	`

	rows, err := s.db.QueryContext(ctx, query, containerTag, time.Now().UTC().Format(time.RFC3339Nano), limit)
	if err != nil {
		return nil, fmt.Errorf("list memories: %w", err)
	}
	defer rows.Close()

	items := make([]memory.Memory, 0, limit)
	for rows.Next() {
		item, err := scanMemoryRow(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate memory rows: %w", err)
	}

	return items, nil
}

func (s *Store) initSchema() error {
	const schema = `
		CREATE TABLE IF NOT EXISTS memories (
			id TEXT PRIMARY KEY,
			container_tag TEXT NOT NULL,
			content TEXT NOT NULL,
			source TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL,
			expires_at TEXT
		);

		CREATE INDEX IF NOT EXISTS idx_memories_container_created_at
			ON memories (container_tag, created_at DESC);
	`

	if _, err := s.db.Exec(schema); err != nil {
		return fmt.Errorf("init sqlite schema: %w", err)
	}

	return nil
}

func scanRecallRow(scanner interface {
	Scan(dest ...any) error
}) (memory.Memory, int64, error) {
	var (
		item      memory.Memory
		createdAt string
		expiresAt sql.NullString
		score     int64
	)

	if err := scanner.Scan(
		&item.ID,
		&item.ContainerTag,
		&item.Content,
		&item.Source,
		&createdAt,
		&expiresAt,
		&score,
	); err != nil {
		return memory.Memory{}, 0, fmt.Errorf("scan recall row: %w", err)
	}

	parsed, err := time.Parse(time.RFC3339Nano, createdAt)
	if err != nil {
		return memory.Memory{}, 0, fmt.Errorf("parse created_at: %w", err)
	}
	item.CreatedAt = parsed

	if expiresAt.Valid {
		parsedExpiry, err := time.Parse(time.RFC3339Nano, expiresAt.String)
		if err != nil {
			return memory.Memory{}, 0, fmt.Errorf("parse expires_at: %w", err)
		}
		item.ExpiresAt = &parsedExpiry
	}

	return item, score, nil
}

func scanMemoryRow(scanner interface {
	Scan(dest ...any) error
}) (memory.Memory, error) {
	var (
		item      memory.Memory
		createdAt string
		expiresAt sql.NullString
	)

	if err := scanner.Scan(
		&item.ID,
		&item.ContainerTag,
		&item.Content,
		&item.Source,
		&createdAt,
		&expiresAt,
	); err != nil {
		return memory.Memory{}, fmt.Errorf("scan memory row: %w", err)
	}

	parsed, err := time.Parse(time.RFC3339Nano, createdAt)
	if err != nil {
		return memory.Memory{}, fmt.Errorf("parse created_at: %w", err)
	}
	item.CreatedAt = parsed

	if expiresAt.Valid {
		parsedExpiry, err := time.Parse(time.RFC3339Nano, expiresAt.String)
		if err != nil {
			return memory.Memory{}, fmt.Errorf("parse expires_at: %w", err)
		}
		item.ExpiresAt = &parsedExpiry
	}

	return item, nil
}

func nullableTime(value *time.Time) any {
	if value == nil {
		return nil
	}
	return value.UTC().Format(time.RFC3339Nano)
}
