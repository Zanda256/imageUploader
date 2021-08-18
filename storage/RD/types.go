package storage

type Img struct {
	ID        string `db:"id"`
	ShortDesc string `db:"description"`
	Region    string `db:"region"`
	Location  string `db:"location"`
	Content   []byte `db:"content"`
	Size      int    `db:"size"`
	Name      string `db:"name"`
	CreatedAt string `db:"added"`
}
