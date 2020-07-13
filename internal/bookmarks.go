package internal

type Bookmark struct {
	Version     string            `json:"version"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	URL         string            `json:"url"`
	Tags        map[string]string `json:"tags"`
	Mode        string            `json:"mode"`
}

type BookmarkRepository interface {
	Find() ([]*Bookmark, error)
	Store(bookmark *Bookmark) (string, error)
	Delete(version string) (bool, error)
}
