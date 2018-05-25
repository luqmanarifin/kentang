package service

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/luqmanarifin/kentang/model"
)

type MySQL struct {
	db *sql.DB
}

// Option holds all necessary options for database.
type MySQLOption struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
	Charset  string
}

// NewMySQL returns a pointer of MySQL instance and error.
func NewMySQL(opt MySQLOption) (*MySQL, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", opt.User, opt.Password, opt.Host, opt.Port, opt.Database, opt.Charset))
	if err != nil {
		return &MySQL{}, err
	}

	return &MySQL{db: db}, nil
}

func (m *MySQL) CreateDictionary(d *model.Dictionary) error {
	_, err := m.db.Exec("INSERT INTO dictionaries(source, keyword, description, timestamp) VALUES(?, ?, ?, ?)",
		d.Source, d.Keyword, d.Description, time.Now())
	return err
}

func (m *MySQL) RemoveDictionary(d *model.Dictionary) error {
	_, err := m.db.Exec("DELETE FROM dictionaries WHERE source=? AND keyword=?",
		d.Source, d.Keyword)
	return err
}

func (m *MySQL) GetDictionary(id int) (model.Dictionary, error) {
	var d model.Dictionary

	err := m.db.QueryRow("SELECT id, source, keyword, description FROM dictionaries WHERE id = ?", id).Scan(&d.ID, &d.Source, &d.Keyword, &d.Description)
	if err != nil {
		return model.Dictionary{}, err
	}

	return d, nil
}

func (m *MySQL) GetDictionaryByKeyword(source, keyword string) (model.Dictionary, error) {
	var d model.Dictionary

	err := m.db.QueryRow("SELECT id, source, keyword, description FROM dictionaries WHERE source = ? AND keyword = ?", source, keyword).Scan(&d.ID, &d.Source, &d.Keyword, &d.Description)
	if err != nil {
		return model.Dictionary{}, err
	}

	return d, nil
}

func (m *MySQL) GetAllDictionaries(source string) ([]model.Dictionary, error) {
	var ds []model.Dictionary

	rows, err := m.db.Query(`
			SELECT id, source, keyword, description
			FROM dictionaries
			WHERE source = ?
	`, source)
	if err != nil {
		return ds, err
	}

	defer rows.Close()
	for rows.Next() {
		var d model.Dictionary

		if err = rows.Scan(&d.ID, &d.Source, &d.Keyword, &d.Description); err != nil {
			return ds, err
		}

		ds = append(ds, d)
	}

	return ds, nil
}

func (m *MySQL) CreateEntry(entry *model.Entry) error {
	_, err := m.db.Exec("INSERT INTO entries(source, keyword, timestamp) VALUES(?, ?, ?)",
		entry.Source, entry.Keyword, time.Now())
	return err
}

func (m *MySQL) RemoveEntryByKeyword(source, keyword string) error {
	_, err := m.db.Exec("DELETE FROM entries WHERE source=? AND keyword=?",
		source, keyword)
	return err
}

func (m *MySQL) RemoveEntryBySource(source string) error {
	_, err := m.db.Exec("DELETE FROM entries WHERE source=?",
		source)
	return err
}

func (m *MySQL) GetAllEntries(source string) ([]model.Entry, error) {
	var es []model.Entry

	rows, err := m.db.Query(`
			SELECT id, source, keyword
			FROM entries
			WHERE source = ?
	`, source)
	if err != nil {
		return es, err
	}

	defer rows.Close()
	for rows.Next() {
		var e model.Entry

		if err = rows.Scan(&e.ID, &e.Source, &e.Keyword); err != nil {
			return es, err
		}

		es = append(es, e)
	}

	return es, nil
}

func (m *MySQL) getEntriesByDay(source string, day int) ([]model.Entry, error) {
	var es []model.Entry

	rows, err := m.db.Query(`
			SELECT id, source, keyword, timestamp
			FROM entries
			WHERE source = ?
			HAVING DATEDIFF(?, timestamp) <= ?
	`, source, time.Now(), day)
	if err != nil {
		return es, err
	}

	defer rows.Close()
	for rows.Next() {
		var e model.Entry
		var dummy string

		if err = rows.Scan(&e.ID, &e.Source, &e.Keyword, &dummy); err != nil {
			return es, err
		}

		es = append(es, e)
	}

	return es, nil
}

func (m *MySQL) GetMonthEntries(source string) ([]model.Entry, error) {
	return m.getEntriesByDay(source, 30)
}

func (m *MySQL) GetWeekEntries(source string) ([]model.Entry, error) {
	return m.getEntriesByDay(source, 7)
}

func (m *MySQL) GetDayEntries(source string) ([]model.Entry, error) {
	return m.getEntriesByDay(source, 1)
}
