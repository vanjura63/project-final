package dbcomment

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Data base.
type DB struct {
	db *pgxpool.Pool
}

type Comment struct {
	Id       int    `json:"id"`
	NewsId   int    `json:"newsId"`
	Text     string `json:"text"`
	ParentId int    `json:"parentId"`
	PubTime  int    `json:"pubTime"`
	Profane  bool   `json:"-"`
}

// Constructor, accepts a DB connection string.
func New(constr string) (*DB, error) {
	db, err := pgxpool.Connect(context.Background(), constr)
	if err != nil {
		return nil, err
	}
	s := DB{
		db: db,
	}
	return &s, nil
}

//Добавление коментария к новости
func (db *DB) AddComment(comment Comment) error {

	_, err := db.db.Exec(context.Background(), `
		INSERT INTO comment(news_id, text_comment, parent_id, pub_time, profane)
		VALUES ($1, $2, $3, $4, $5)`,
		comment.NewsId,
		comment.Text,
		comment.ParentId,
		comment.PubTime,
		comment.Profane,
	)

	return err
}

// Коментарии к новости.
func (db *DB) Comment(id int) ([]Comment, error) {

	rows, err := db.db.Query(context.Background(), `
	SELECT id, news_id, text_comment, parent_id, pub_time, profane FROM comment 
	WHERE news_id=$1
	ORDER BY pub_time DESC
		`,
		id,
	)
	if err != nil {
		return nil, err
	}
	var comments []Comment
	for rows.Next() {
		var p Comment
		err = rows.Scan(
			&p.Id,
			&p.NewsId,
			&p.Text,
			&p.ParentId,
			&p.PubTime,
			&p.Profane,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, p)
	}
	return comments, rows.Err()
}
