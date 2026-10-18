package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type MediaFile struct {
	Id           int64        `db:"id"`
	FileName     string       `db:"file_name"`
	OriginalName string       `db:"original_name"`
	Mime         string       `db:"mime"`
	Size         int64        `db:"size"`
	UploaderId   int64        `db:"uploader_id"`
	CreatedAt    time.Time    `db:"created_at"`
	DeletedAt    sql.NullTime `db:"deleted_at"`
}

type MediaFileModel interface {
	Insert(ctx context.Context, data *MediaFile) (int64, error)
	FindActiveByName(ctx context.Context, name string) (*MediaFile, error)
	SoftDelete(ctx context.Context, name string) (int64, error)
}

type customMediaFileModel struct {
	conn sqlx.SqlConn
}

func NewMediaFileModel(conn sqlx.SqlConn) MediaFileModel {
	return &customMediaFileModel{conn: conn}
}

const (
	insertSQL      = `INSERT INTO media_file (file_name, original_name, mime, size, uploader_id) VALUES (?, ?, ?, ?, ?)`
	findActiveSQL  = `SELECT id, file_name, original_name, mime, size, uploader_id, created_at, deleted_at FROM media_file WHERE file_name = ? AND deleted_at IS NULL LIMIT 1`
	softDeleteSQL  = `UPDATE media_file SET deleted_at = NOW() WHERE file_name = ? AND deleted_at IS NULL`
)

func (m *customMediaFileModel) Insert(ctx context.Context, data *MediaFile) (int64, error) {
	res, err := m.conn.ExecCtx(ctx, insertSQL, data.FileName, data.OriginalName, data.Mime, data.Size, data.UploaderId)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (m *customMediaFileModel) FindActiveByName(ctx context.Context, name string) (*MediaFile, error) {
	var data MediaFile
	if err := m.conn.QueryRowCtx(ctx, &data, findActiveSQL, name); err != nil {
		return nil, err
	}
	return &data, nil
}

func (m *customMediaFileModel) SoftDelete(ctx context.Context, name string) (int64, error) {
	res, err := m.conn.ExecCtx(ctx, softDeleteSQL, name)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
