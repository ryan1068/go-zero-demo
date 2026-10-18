package logic

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"go-zero-demo/media-api/internal/model"
	"go-zero-demo/media-api/internal/storage"
	"go-zero-demo/pkg/auth"
	"go-zero-demo/pkg/errcode"
)

const findMediaSQL = "(?i)select .* from media_file where file_name = \\? and deleted_at is null"
const softDeleteMediaSQL = "(?i)update media_file set deleted_at"

func mediaRows(uploaderID int64) *sqlmock.Rows {
	return sqlmock.NewRows([]string{
		"id", "file_name", "original_name", "mime", "size", "uploader_id", "created_at", "deleted_at",
	}).AddRow(1, "a.png", "a.png", "image/png", 16, uploaderID, time.Now(), nil)
}

func TestDeleteFile_InvalidName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create sqlmock: %v", err)
	}
	defer db.Close()

	dir := t.TempDir()
	if writeErr := os.WriteFile(filepath.Join(dir, "a.png"), pngBytes, 0o644); writeErr != nil {
		t.Fatalf("write file: %v", writeErr)
	}

	ctx := auth.WithUserRole(context.Background(), 1, "admin", auth.RoleAdmin)
	logic := NewDeleteFileLogic(ctx, storage.NewLocal(dir), model.NewMediaFileModel(sqlx.NewSqlConnFromDB(db)))
	resp, err := logic.Delete("../a.png")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Code != errcode.BadRequest {
		t.Fatalf("expected code %d, got %d (%s)", errcode.BadRequest, resp.Code, resp.Msg)
	}
	if _, statErr := os.Stat(filepath.Join(dir, "a.png")); statErr != nil {
		t.Fatalf("file should still exist: %v", statErr)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestDeleteFile_Permissions(t *testing.T) {
	cases := []struct {
		name         string
		ctx          context.Context
		fileExists   bool
		rows         *sqlmock.Rows
		wantCode     int32
		wantDelete   bool
		wantFileGone bool
	}{
		{
			name:         "owner deletes",
			ctx:          auth.WithUserRole(context.Background(), 2, "buyer", 0),
			fileExists:   true,
			rows:         mediaRows(2),
			wantCode:     errcode.Success,
			wantDelete:   true,
			wantFileGone: true,
		},
		{
			name:         "admin deletes another user's file",
			ctx:          auth.WithUserRole(context.Background(), 1, "admin", auth.RoleAdmin),
			fileExists:   true,
			rows:         mediaRows(2),
			wantCode:     errcode.Success,
			wantDelete:   true,
			wantFileGone: true,
		},
		{
			name:         "other user is forbidden",
			ctx:          auth.WithUserRole(context.Background(), 3, "other", 0),
			fileExists:   true,
			rows:         mediaRows(2),
			wantCode:     errcode.Forbidden,
			wantDelete:   false,
			wantFileGone: false,
		},
		{
			name:         "missing file",
			ctx:          auth.WithUserRole(context.Background(), 2, "buyer", 0),
			fileExists:   false,
			rows:         sqlmock.NewRows([]string{"id"}),
			wantCode:     errcode.NotFound,
			wantDelete:   false,
			wantFileGone: true,
		},
		{
			name:         "unauthenticated",
			ctx:          context.Background(),
			fileExists:   true,
			rows:         nil,
			wantCode:     errcode.Unauthorized,
			wantDelete:   false,
			wantFileGone: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("create sqlmock: %v", err)
			}
			defer db.Close()

			dir := t.TempDir()
			if tc.fileExists {
				if writeErr := os.WriteFile(filepath.Join(dir, "a.png"), pngBytes, 0o644); writeErr != nil {
					t.Fatalf("write file: %v", writeErr)
				}
			}

			if tc.rows != nil {
				mock.ExpectQuery(findMediaSQL).WithArgs("a.png").WillReturnRows(tc.rows)
			}
			if tc.wantDelete {
				mock.ExpectExec(softDeleteMediaSQL).WithArgs("a.png").WillReturnResult(sqlmock.NewResult(0, 1))
			}

			logic := NewDeleteFileLogic(tc.ctx, storage.NewLocal(dir), model.NewMediaFileModel(sqlx.NewSqlConnFromDB(db)))
			resp, err := logic.Delete("a.png")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if resp.Code != tc.wantCode {
				t.Fatalf("expected code %d, got %d (%s)", tc.wantCode, resp.Code, resp.Msg)
			}

			_, statErr := os.Stat(filepath.Join(dir, "a.png"))
			fileGone := os.IsNotExist(statErr)
			if fileGone != tc.wantFileGone {
				t.Fatalf("expected file gone=%v, got %v (err=%v)", tc.wantFileGone, fileGone, statErr)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet sql expectations: %v", err)
			}
		})
	}
}
