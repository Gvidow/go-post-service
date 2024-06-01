package entity

type FeedComment struct {
	Comments []*Comment `json:"comments"`
	Cursor   int        `json:"cursor"`
}

type FeedPost struct {
	Posts  []*Post `json:"posts"`
	Cursor int     `json:"cursor"`
}
