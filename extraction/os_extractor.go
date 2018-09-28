package extraction

//OsVersionInfo ...
type OsVersionInfo struct {
	Name    string
	Version string
}

type OsVersionInfoIterable interface {
	Next() (*OsVersionInfo, error)
	ForEach(func(*OsVersionInfo) error) error
}

type OsExtractor interface {
	OsVersionInfoIterable
	getOsVersion() (*OsVersionInfo, error)
}
