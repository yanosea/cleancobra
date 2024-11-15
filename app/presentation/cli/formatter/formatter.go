package formatter

type Formatter interface {
	Format(result interface{}) (string, error)
}
