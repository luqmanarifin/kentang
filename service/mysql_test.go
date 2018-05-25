package service

import (
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/luqmanarifin/kentang/model"
)

func getConnection(t *testing.T) (*MySQL, error) {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatal("Error loading .env file")
	}

	opt := MySQLOption{
		User:     os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		Host:     os.Getenv("MYSQL_HOST"),
		Port:     os.Getenv("MYSQL_PORT"),
		Database: os.Getenv("MYSQL_DATABASE"),
		Charset:  os.Getenv("MYSQL_CHARSET"),
	}
	m, err := NewMySQL(opt)
	if err != nil {
		t.Fatalf("%s", err.Error())
		return &MySQL{}, err
	}
	return m, nil
}

func TestCreateDictionary(t *testing.T) {
	m, err := getConnection(t)
	if err != nil {
		t.Fatal("error connecting")
	}
	dict := &model.Dictionary{
		Source:      "source",
		Keyword:     "a",
		Description: "b",
	}
	err = m.CreateDictionary(dict)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
}

func TestRemoveDictionary(t *testing.T) {
	m, err := getConnection(t)
	if err != nil {
		t.Fatal("error connecting")
	}
	dict := &model.Dictionary{
		Source:      "source",
		Keyword:     "a",
		Description: "b",
	}
	err = m.RemoveDictionary(dict)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
}

func TestGetDictionary(t *testing.T) {
	m, err := getConnection(t)
	if err != nil {
		t.Fatal("error connecting")
	}
	dict, err := m.GetDictionary(1)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	t.Logf("%+v", dict)
	// fmt.Printf("%d, %s, %s, %s\n", dict.ID, dict.Source, dict.Keyword, dict.Description)
}

func TestGetDictionaryByKeyword(t *testing.T) {
	m, err := getConnection(t)
	if err != nil {
		t.Fatal("error connecting")
	}
	dict, err := m.GetDictionaryByKeyword("source", "bird")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	t.Logf("%+v", dict)
	// fmt.Printf("%d, %s, %s, %s\n", dict.ID, dict.Source, dict.Keyword, dict.Description)
}

func TestGetAllDictionaries(t *testing.T) {
	m, err := getConnection(t)
	if err != nil {
		t.Fatal("error connecting")
	}
	dict, err := m.GetAllDictionaries("source")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	t.Logf("%+v", dict)
	// fmt.Printf("%d, %s, %s, %s\n", dict.ID, dict.Source, dict.Keyword, dict.Description)

}

func TestCreateEntry(t *testing.T) {
	m, err := getConnection(t)
	if err != nil {
		t.Fatal("error connecting")
	}
	entry := &model.Entry{
		Source:  "source",
		Keyword: "asuasu",
	}
	err = m.CreateEntry(entry)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
}

func TestRemoveEntryByKeyword(t *testing.T) {
	m, err := getConnection(t)
	if err != nil {
		t.Fatal("error connecting")
	}
	entry := &model.Entry{
		Source:  "source",
		Keyword: "asuasu",
	}
	err = m.RemoveEntryByKeyword(entry.Source, entry.Keyword)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
}

func TestRemoveEntryBySource(t *testing.T) {
	m, err := getConnection(t)
	if err != nil {
		t.Fatal("error connecting")
	}
	entry := &model.Entry{
		Source:  "luqman",
		Keyword: "a",
	}
	err = m.RemoveEntryBySource(entry.Source)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
}

func TestGetAllEntries(t *testing.T) {
	m, err := getConnection(t)
	if err != nil {
		t.Fatal("error connecting")
	}
	entries, err := m.GetAllEntries("source")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	for _, e := range entries {
		t.Logf("%+v", e)
	}
}

func TestGetMonthEntries(t *testing.T) {
	m, err := getConnection(t)
	if err != nil {
		t.Fatal("error connecting")
	}
	entries, err := m.GetMonthEntries("source")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	for _, e := range entries {
		t.Logf("%+v", e)
	}
}
