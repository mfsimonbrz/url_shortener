package data

import (
	"database/sql"
	"url_shortener/internals/models"
)

type EntryData struct {
	db *sql.DB
}

func NewEntryData(db *sql.DB) *EntryData {
	return &EntryData{db: db}
}

func (d *EntryData) InitDB() error {
	sql := "CREATE TABLE IF NOT EXISTS public.entries (id INTEGER PRIMARY KEY NOT NULL GENERATED ALWAYS AS IDENTITY, url TEXT, short_url TEXT, entry_date DATE)"

	_, err := d.db.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

func (d *EntryData) CreateUrlEntry(entry *models.Entry) error {
	sql := "INSERT INTO public.entries (url, short_url, entry_date) VALUES ($1, $2, now())"
	_, err := d.db.Exec(sql, entry.Url, entry.ShortUrl)
	if err != nil {
		return err
	}

	return nil
}

func (d *EntryData) GetUrlEntryByShortUrl(url string) (*models.Entry, error) {
	sql := "SELECT * FROM public.entries WHERE short_url = $1"
	var entry models.Entry

	row := d.db.QueryRow(sql, url)
	err := row.Scan(&entry.ID, &entry.Url, &entry.ShortUrl, &entry.Date)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func (d *EntryData) GetUrlEntryByFullUrl(url string) (*models.Entry, error) {
	sql := "SELECT * FROM public.entries WHERE url = $1"
	var entry models.Entry

	row := d.db.QueryRow(sql, url)
	err := row.Scan(&entry.ID, &entry.Url, &entry.ShortUrl, &entry.Date)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}
