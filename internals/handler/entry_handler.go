package handler

import (
	"net/url"
	"url_shortener/internals/cache"
	"url_shortener/internals/data"
	"url_shortener/internals/models"
	"url_shortener/internals/utils"
)

type EntryHandler struct {
	entryData    *data.EntryData
	cacheHandler *cache.CacheHandler
}

func NewEntryHandler(entryData *data.EntryData, cacheHandler *cache.CacheHandler) *EntryHandler {
	return &EntryHandler{entryData: entryData, cacheHandler: cacheHandler}
}

func (h *EntryHandler) AddUrlEntry(fullUrl string) (*models.Entry, error) {
	if _, err := url.Parse(fullUrl); err != nil {
		return nil, err
	}

	retrievedEntry, err := h.entryData.GetUrlEntryByFullUrl(fullUrl)
	if err != nil {
		entry := &models.Entry{Url: fullUrl, ShortUrl: utils.GenerateRandomString()}
		h.entryData.CreateUrlEntry(entry)
		h.cacheHandler.Add(entry)

		return entry, nil
	}

	return retrievedEntry, nil
}

func (h *EntryHandler) RetrieveUrl(shortUrl string) (string, error) {
	result, err := h.cacheHandler.Get(shortUrl)
	if err != nil {
		entry, err := h.entryData.GetUrlEntryByShortUrl(shortUrl)
		if err != nil {
			return "", err
		}

		err = h.cacheHandler.Add(entry)
		if err != nil {
			return "", err
		}

		return entry.Url, nil
	} else {
		return result, nil
	}
}
