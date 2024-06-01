package entity

type Comment struct {
	ID        int          `json:"id"`
	Author    string       `json:"author"`
	Content   string       `json:"content"`
	CreatedAt int          `json:"created_at"`
	Replies   *FeedComment `json:"replies"`
}
