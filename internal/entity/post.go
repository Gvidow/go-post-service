package entity

type Post struct {
	ID           int          `json:"id"`
	Author       string       `json:"author"`
	Title        string       `json:"title"`
	Content      string       `json:"content"`
	AllowComment bool         `json:"allowComment"`
	CreatedAt    int          `json:"created_at"`
	Comments     *FeedComment `json:"comments"`
}
