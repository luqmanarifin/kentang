package service

import (
	"database/sql"
	"fmt"
)

type MySQL struct {
	db *sql.DB
}

// Option holds all necessary options for database.
type Option struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
	Charset  string
}

// NewMySQL returns a pointer of MySQL instance and error.
func NewMySQL(opt Option) (*MySQL, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", opt.User, opt.Password, opt.Host, opt.Port, opt.Database, opt.Charset))
	if err != nil {
		return &MySQL{}, err
	}

	return &MySQL{db: db}, nil
}

func (m *MySQL) AddDictionary(d *Dictionary) error {

}

func (m *MySQL) RemoveDictionary(d *Dictionary) error {

}

func (m *MySQL) GetDictionary(source, keyword string) (Dictionary, error) {

}

func (m *MySQL) GetAllDictionaries(source string) ([]Dictionary, error) {

}

func (m *MySQL) GetAllEntries(source string) ([]Entry, error) {

}
