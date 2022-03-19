package urls

type UrlMapping struct {
	Key     string `bson:"key"`
	URL     string `bson:"url"`
	Counter uint32 `bson:"counter"`
}
