package structs

type News2 struct {
	Id      int64  `db:"Id"`
	Title   string `db:"Title"`
	Content string `db:"Content"`
}

type NewsResult struct {
	Id         int64   `db:"Id"`
	Title      string  `db:"Title"`
	Content    string  `db:"Content"`
	Categories []int64 `db:"Categories"`
}
type PostNews struct {
	Id         int64   `json:"id"`
	Title      string  `json:"title"`
	Content    string  `json:"content"`
	Categories []int64 `json:"categories"`
}

type NewsCategories struct {
	NewsId     int64 `db:"NewsId"`
	Categoryid int64 `db:"CategoryId"`
}
