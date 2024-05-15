package domain

import (
	"context"
	"database/sql"
	db "github.com/Mire0726/unibox/backend/infrastructure/mysql"
	"github.com/Mire0726/unibox/backend/domain/model"
)

type UserRepository interface {
    // GetUserByID は指定されたIDを持つユーザーを取得します。
    GetUserByID(ctx context.Context, id string) (*domain.User, error)

    // CreateUser は新しいユーザーをデータベースに保存します。
    CreateUser(ctx context.Context, user *domain.User) error

    // UpdateUser は既存のユーザー情報を更新します。
    UpdateUser(ctx context.Context, user *domain.User) error

    // DeleteUser は指定されたIDを持つユーザーをデータベースから削除します。
    DeleteUser(ctx context.Context, id string) error
}

type SQLUserRepository struct {
    db *sql.DB
}

func NewSQLUserRepository(db *sql.DB) UserRepository {
    return &SQLUserRepository{db: db}
}