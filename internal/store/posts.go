package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64
	Content   string `json:"content"`
	Title     string `json:"title"`
	UserID    int64
	Tags      []string `json:"tags"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string
	Version   int        `json:"version"`
	Comments  []*Comment `json:"comments"`
	User      User       `json:"user"`
}

type PostWithMetaData struct {
	Post
	CommentCount int `json:"comments_count"`
}

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) GetUserFeed(ctx context.Context, userID int64) ([]PostWithMetaData, error) {

	query := `
        SELECT 
            p.id,
            p.user_id,
            p.title,
            p.content,
            p.created_at,
            p.version,
            p.tags,
            u.username,
            COUNT(c.id) AS comments_count
        FROM posts p
        LEFT JOIN comments c ON c.post_id = p.id
        JOIN followers f ON f.follower_id = p.user_id
        JOIN users u ON u.id = p.user_id
        WHERE f.follower_id = $1 OR p.user_id = $1
        GROUP BY p.id, u.username
        ORDER BY p.created_at DESC
    `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feed []PostWithMetaData

	for rows.Next() {
		var p PostWithMetaData

		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Title,
			&p.Content,
			&p.CreatedAt,
			&p.Version,
			pq.Array(&p.Tags),
			&p.User.UserName,
			&p.CommentCount,
		)
		if err != nil {
			return nil, err
		}

		feed = append(feed, p)
	}

	return feed, nil
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `INSERT INTO posts(content,title,user_id,tags)
	values ($1, $2, $3,$4) RETURNING id, created_at, updated_at`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Content, post.Title, post.UserID, pq.Array(post.Tags),
	).Scan(
		&post.ID, &post.CreatedAt, &post.UpdatedAt,
	)
	if err != nil {
		return err

	}
	return nil
}

func (s *PostStore) GetById(ctx context.Context, id int64) (*Post, error) {
	query := `SELECT id,user_id,title,content,created_at,updated_at,version,tags FROM posts WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var post Post
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Version,
		pq.Array(&post.Tags),
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err

		}

	}

	return &post, nil
}
func (s *PostStore) Delete(ctx context.Context, postID int64) error {
	query := `DELETE FROM posts WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, postID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}
func (s *PostStore) Update(ctx context.Context, post *Post) error {
	query := `UPDATE posts SET title = $1 , content = $2 , version = version +1
	WHERE id = $3 AND version = $4
	RETURNING version`

	err := s.db.QueryRowContext(ctx, query, post.Title, post.Content, post.ID, post.Version).Scan(&post.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}
	return nil

}
