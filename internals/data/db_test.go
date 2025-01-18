package data

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"url_shortener/internals/models"

	_ "github.com/lib/pq"
)

func setup() (*sql.DB, error) {
	dbHost := os.Getenv("testDbHost")
	dbPort := os.Getenv("testDbPort")
	dbUser := os.Getenv("testDbUser")
	dbPass := os.Getenv("testDbPass")
	dbName := os.Getenv("testDbName")

	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName))
	if err != nil {
		return nil, err
	}

	return db, nil
}

func shutdown(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM public.entries")
	if err != nil {
		return err
	}

	defer db.Close()

	return nil
}

func TestCreateUrlEntry(t *testing.T) {
	db, err := setup()
	if err != nil {
		t.Error(err)
	}

	defer shutdown(db)

	entryData := NewEntryData(db)
	entryData.InitDB()

	entry := &models.Entry{Url: "https://g1.globo.com/"}
	err = entryData.CreateUrlEntry(entry)
	if err != nil {
		t.Error(err)
	}
	expected := "https://g1.globo.com/"
	row := db.QueryRow("SELECT url FROM public.entries LIMIT 1")
	var entryFromDb models.Entry
	row.Scan(&entryFromDb.Url)
	got := entryFromDb.Url

	if got != expected {
		t.Errorf("got %q, expected %q", got, expected)
	}
}

func TestGetUrlEntryByShortUrl(t *testing.T) {
	db, err := setup()
	defer shutdown(db)
	if err != nil {
		t.Error(err)
	}

	_, err = db.Exec("INSERT INTO public.entries (url, short_url) VALUES ($1, $2)", "https://g1.globo.com/", "abc1234")
	if err != nil {
		t.Error(err)
	}

	expected := "https://g1.globo.com/"
	row := db.QueryRow("SELECT url FROM public.entries WHERE short_url = 'abc1234'")
	var entryFromDb models.Entry
	row.Scan(&entryFromDb.Url)
	got := entryFromDb.Url

	if got != expected {
		t.Errorf("got %q, expected %q", got, expected)
	}
}

func TestGetUrlEntryByFullUrl(t *testing.T) {
	db, err := setup()
	defer shutdown(db)
	if err != nil {
		t.Error(err)
	}

	_, err = db.Exec("INSERT INTO public.entries (url, short_url) VALUES ($1, $2)", "https://g1.globo.com/", "abc1234")
	if err != nil {
		t.Error(err)
	}

	expected := "abc1234"
	row := db.QueryRow("SELECT short_url FROM public.entries WHERE url = 'https://g1.globo.com/'")
	var entryFromDb models.Entry
	row.Scan(&entryFromDb.ShortUrl)
	got := entryFromDb.ShortUrl

	if got != expected {
		t.Errorf("got %q, expected %q", got, expected)
	}
}
