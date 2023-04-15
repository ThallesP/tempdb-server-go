package databases

import (
	"crypto/rand"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type DatabasesStorage struct {
	db *sql.DB
}

func NewDatabasesStorage(db *sql.DB) *DatabasesStorage {
	return &DatabasesStorage{db: db}
}

type CreateUserResponse struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (d *DatabasesStorage) Create() (databaseName string, err error) {
	databaseName = uuid.New().String()
	_, err = d.db.Exec(fmt.Sprintf(`CREATE DATABASE "%s"`, databaseName))

	if err != nil {
		return "", err
	}

	_, err = d.db.Exec(fmt.Sprintf(`REVOKE connect ON DATABASE "%s" FROM PUBLIC`, databaseName))

	return databaseName, err
}

func (d *DatabasesStorage) Delete(databaseName string) error {
	_, err := d.db.Exec(fmt.Sprintf(`SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = '%s'`, databaseName))
	
	if err != nil {
		return err
	}

	_, err = d.db.Exec(fmt.Sprintf(`DROP DATABASE "%s"`, databaseName))

	if err != nil {
		return err
	}

	return nil
}

func (d *DatabasesStorage) CreateUser(databaseName string) (createUserResponse *CreateUserResponse, err error) {
	userBuf := make([]byte, 8)
	passwordBuf := make([]byte, 8)

	_, err = rand.Read(userBuf)

	if err != nil {
		return nil, err
	}

	_, err = rand.Read(passwordBuf)

	if err != nil {
		return nil, err
	}

	username := fmt.Sprintf("%x", userBuf)
	password := fmt.Sprintf("%x", passwordBuf)

	_, err = d.db.Exec(fmt.Sprintf(`CREATE USER "%s" WITH PASSWORD '%s'`, username, password))

	if err != nil {
		return nil, err
	}

	_, err = d.db.Exec(fmt.Sprintf(`GRANT connect ON DATABASE "%s" TO "%s"`, databaseName, username))

	if err != nil {
		return nil, err
	}

	return &CreateUserResponse{Username: username, Password: password}, nil
}