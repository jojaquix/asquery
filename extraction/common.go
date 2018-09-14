package extraction

//InfoIterable ... is a generic closable interface for iterating over Infos
type InfoIterable interface {
	Next() (*InfoIterable, error)
	ForEach(func(*InfoIterable) error) error
	Close()
}
