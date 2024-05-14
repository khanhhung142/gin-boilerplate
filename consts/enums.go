package consts

type NoSQLCollection string

const (
	MongoDBCollectionUsers NoSQLCollection = "users"
)

func (m NoSQLCollection) String() string {
	return string(m)
}
