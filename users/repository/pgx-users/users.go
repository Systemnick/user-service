package pgx_users

import (
	"context"
	"time"

	"github.com/Systemnick/user-service/domain"
	"github.com/Systemnick/user-service/users"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx"
	"github.com/pkg/errors"
)

type pgxUsersRepository struct {
	schemaName string
	tableName  string
	timeout    time.Duration
	connection *pgx.Conn
}

func New(connString, schemaName, tableName string, timeout time.Duration) (users.Repository, error) {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	connection, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, errors.Wrap(err, "pgx.Connect")
	}

	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	err = connection.Ping(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "connection.Ping")
	}

	r := &pgxUsersRepository{
		schemaName: schemaName,
		tableName:  tableName,
		timeout:    timeout,
		connection: connection,
	}

	return r, nil
}

func (r pgxUsersRepository) Store(ctx context.Context, user *domain.User) error {
	var id pgtype.UUID

	s := `
		INSERT INTO users(id, login, password, email, phone, creation_time)
		VALUES($1, $2, $3, $4, $5, $6)
	`
	_, err := r.connection.Prepare(ctx, "new_user", s)
	if err != nil {
		return errors.Wrap(err, "failed to prepare query")
	}

	err = id.DecodeText(nil, []byte(user.ID))
	if err != nil {
		return errors.Wrap(err, "failed to convert UUID")
	}

	_, err = r.connection.Exec(ctx, "new_user", id, user.Login, user.Password, user.Email, user.Phone, user.CreationTime)
	if err != nil {
		return errors.Wrap(err, "failed to store user")
	}

	return nil
}

func (r pgxUsersRepository) Fetch(ctx context.Context, login, password string) (*domain.User, error) {
	var u domain.User

	u.Login = login
	u.Password = password

	query := "SELECT id FROM users WHERE login=$1 AND password=$2"

	err := r.connection.QueryRow(ctx, query, login, password).Scan(&u.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch user")
	}

	return &u, nil
}
