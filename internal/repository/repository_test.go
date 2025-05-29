package repository

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ecoderat/dispatch-go/internal/model"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	assert.NoError(t, err)
	cleanup := func() { db.Close() }
	return gormDB, mock, cleanup
}

func TestMessageRepository_GetAll(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()
	logger := &logrus.Logger{}
	repo := NewMessageRepository(db, logger)

	query := `SELECT * FROM "message" WHERE "message"."deleted_at" IS NULL`

	rows := sqlmock.NewRows([]string{"id", "recipient", "content", "status"}).
		AddRow(1, "+123", "hi", "pending")
	mock.ExpectQuery(query).WillReturnRows(rows)

	msgs, err := repo.GetAll(context.Background())
	assert.NoError(t, err)
	assert.Len(t, msgs, 1)
	assert.Equal(t, "+123", msgs[0].Recipient)
	assert.Equal(t, "hi", msgs[0].Content)
	assert.Equal(t, "pending", string(msgs[0].Status))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_GetAll_WithStatus(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()
	logger := &logrus.Logger{}
	repo := NewMessageRepository(db, logger)

	query := `SELECT * FROM "message" WHERE status IN ($1,$2) AND "message"."deleted_at" IS NULL`

	rows := sqlmock.NewRows([]string{"id", "recipient", "content", "status"}).
		AddRow(1, "+123", "hi", "pending").
		AddRow(2, "+456", "hello", "failed")
	mock.ExpectQuery(query).WithArgs("pending", "failed").WillReturnRows(rows)

	msgs, err := repo.GetAll(context.Background(), model.StatusPending, model.StatusFailed)
	assert.NoError(t, err)
	assert.Len(t, msgs, 2)
	assert.Equal(t, "+123", msgs[0].Recipient)
	assert.Equal(t, "hi", msgs[0].Content)
	assert.Equal(t, "pending", string(msgs[0].Status))
	assert.Equal(t, "+456", msgs[1].Recipient)
	assert.Equal(t, "hello", msgs[1].Content)
	assert.Equal(t, "failed", string(msgs[1].Status))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_GetAll_Fails(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()
	logger := &logrus.Logger{}
	repo := NewMessageRepository(db, logger)

	query := `SELECT * FROM "message" WHERE "message"."deleted_at" IS NULL`

	mock.ExpectQuery(query).WillReturnError(assert.AnError)

	msgs, err := repo.GetAll(context.Background())
	assert.Error(t, err)
	assert.Nil(t, msgs)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_Create(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()
	logger := &logrus.Logger{}
	repo := NewMessageRepository(db, logger)

	query := `INSERT INTO "message" ("recipient","content","status","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`

	msg := model.Message{Recipient: "+123", Content: "hi", Status: "pending"}
	mock.ExpectBegin()
	mock.ExpectQuery(query).WithArgs(
		msg.Recipient,
		msg.Content,
		string(msg.Status),
		sqlmock.AnyArg(), // created_at
		sqlmock.AnyArg(), // updated_at
		sqlmock.AnyArg(), // deleted_at
	).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.Create(context.Background(), msg)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_Update(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()
	logger := &logrus.Logger{}
	repo := NewMessageRepository(db, logger)

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "message" SET "status"=$1,"updated_at"=$2 WHERE id = $3 AND "message"."deleted_at" IS NULL`).
		WithArgs("sent", sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Update(context.Background(), 1, "sent")
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_Delete(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()
	logger := &logrus.Logger{}
	repo := NewMessageRepository(db, logger)

	query := `UPDATE "message" SET "deleted_at"=$1 WHERE id = $2 AND "message"."deleted_at" IS NULL`
	messageID := 1

	mock.ExpectBegin()
	mock.ExpectExec(query).WithArgs(sqlmock.AnyArg(), messageID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete(context.Background(), messageID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
