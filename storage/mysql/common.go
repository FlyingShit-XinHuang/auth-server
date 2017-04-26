package mysql

import (
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
)`, `CREATE TABLE IF NOT EXISTS {prefix}client (
	id           varchar(255) NOT NULL PRIMARY KEY,
	secret 		 varchar(255) NOT NULL,
	extra 		 varchar(255) NOT NULL,
	redirect_uri varchar(255) NOT NULL
)`,
}

const tabPrefix = "whispir_"

var notFoundError = merry.New("Not found")

type mysqlStorage struct {
	db *sql.DB
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
		db: handle,
	}

	if err := store.createSchemas(); nil != err {
		return nil, fmt.Errorf("Failed to create schemas: %s", err)
	}
	return store, nil
}

func (s *mysqlStorage) createSchemas() error {
	for k, schema := range schemas {
		schema := strings.Replace(schema, "{prefix}", tabPrefix, 4)
		if _, err := s.db.Exec(schema); err != nil {
			log.Printf("Error creating schema %d: %s", k, schema)
			return err
		}
	}
	return nil
}

func (s *mysqlStorage) Clone() osin.Storage {
	return s
}

// Close the resources the Storage potentially holds (using Clone for example)
func (s *mysqlStorage) Close() {
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
	if _, err := s.db.Exec(
		fmt.Sprintf("INSERT INTO %sclient (id, secret, redirect_uri, extra) VALUES (?, ?, ?, ?)", tabPrefix),
		client.Id, client.Secret, client.RedirectURL, client.Name); err != nil {

		return merry.Wrap(err)
	}
	return nil
}

func (s *mysqlStorage) GetClient(id string) (osin.Client, error) {
	row := s.db.QueryRow(fmt.Sprintf("SELECT id, secret, redirect_uri, extra FROM %sclient WHERE id=?", tabPrefix), id)
	var c osin.DefaultClient
	var extra string

	if err := row.Scan(&c.Id, &c.Secret, &c.RedirectUri, &extra); err == sql.ErrNoRows {
		return nil, notFoundError
	} else if err != nil {
		return nil, merry.Wrap(err)
	}
	c.UserData = extra
	return &c, nil
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

func (s *mysqlStorage) GetUserById(id int) (*v1alpha1.User, error) {
	row := s.db.QueryRow(
		fmt.Sprintf("SELECT * FROM %susers WHERE id=?", tabPrefix), id)
	var user v1alpha1.User

	if err := row.Scan(&user.Id, &user.Name, &user.Password); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, merry.Wrap(err)
	}
	return &user, nil
}
