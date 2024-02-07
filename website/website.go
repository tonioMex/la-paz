package website

type Website struct {
	ID   int64
	Name string `gorm:"uniqueIndex;not null"`
	URL  string `gorm:"not null"`
	Rank int64  `gorm:"not null"`
}

func (Website) TableName() string {
	return "websites"
}
