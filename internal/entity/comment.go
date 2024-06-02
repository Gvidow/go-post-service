package entity

type Comment struct {
	ID        int    `json:"id"`
	Author    string `json:"author"`
	Content   string `json:"content"`
	Parent    int    `json:"parent"`
	Depth     int    `json:"depth"`
	CreatedAt int    `json:"created_at"`
}

type NotifyComment struct {
	Err     error
	Comment *Comment
}
