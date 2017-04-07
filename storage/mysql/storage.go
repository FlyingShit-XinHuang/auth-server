package mysql

import (
	omysql "github.com/felipeweb/osin-mysql"

	_ "github.com/go-sql-driver/mysql"

	"whispir/auth-server/storage"

	"database/sql"
	"fmt"
	"github.com/RangelReale/osin"
	"github.com/ansel1/merry"
	"log"
	"strings"
	"whispir/auth-server/pkg/api/v1alpha1"
)

// TODO: wrap OAuth2Storage interface with an ORM package (e.g. https://github.com/jinzhu/gorm)

var schemas = []string{`CREATE TABLE IF NOT EXISTS {prefix}users (
	id 		int(11) NOT NULL PRIMARY KEY AUTO_INCREMENT,
	name		varchar(255) NOT NULL,
	password	varchar(255) NOT NULL,
	UNIQUE INDEX users_index (name)
)`,
}

const tabPrefix = "whispir_"

type mysqlStorage struct {
	db *sql.DB
	*omysql.Storage
}

func NewStorageOrDie(user, password, host string, port int, dbname string) storage.OAuth2Storage {
	storage, err := NewStorage(user, password, host, port, dbname)
	if nil != err {
		panic(err)
	}

	return storage
}

func NewStorage(user, password, host string, port int, dbname string) (storage.OAuth2Storage, error) {
	handle, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", user, password, host, port, dbname))
	if nil != err {
		return nil, fmt.Errorf("Failed to open mysql connections: %s", err)
	}

	store := &mysqlStorage{
		db:      handle,
		Storage: omysql.New(handle, tabPrefix),
	}

	if err := store.createSchemas(); nil != err {
		return nil, fmt.Errorf("Failed to create schemas: %s", err)
	}
	return store, nil
}

func (s *mysqlStorage) createSchemas() error {
	if err := s.CreateSchemas(); nil != err {
		return err
	}
	for k, schema := range schemas {
		schema := strings.Replace(schema, "{prefix}", tabPrefix, 4)
		if _, err := s.db.Exec(schema); err != nil {
			log.Printf("Error creating schema %d: %s", k, schema)
			return err
		}
	}
	return nil
}

func (s *mysqlStorage) CreateUser(user *v1alpha1.User) error {
	if _, err := s.db.Exec(
		fmt.Sprintf("INSERT INTO %susers (name, password) VALUES (?, ?)", tabPrefix),
		user.Name, user.Password); err != nil {

		return merry.Wrap(err)
	}
	return nil
}

func (s *mysqlStorage) CreateClient(client *v1alpha1.Client) error {
	return s.Storage.CreateClient(&osin.DefaultClient{
		Id:          client.Id,
		Secret:      client.Secret,
		RedirectUri: client.RedirectURL,
		UserData:    client.Name,
	})
}

func (s *mysqlStorage) GetUserByNameAndPassword(name, password string) (*v1alpha1.User, error) {
	row := s.db.QueryRow(
		fmt.Sprintf("SELECT * FROM %susers WHERE name=? and password=?", tabPrefix), name, password)
	var user v1alpha1.User

	if err := row.Scan(&user.Id, &user.Name, &user.Password); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, merry.Wrap(err)
	}
	return &user, nil
}
