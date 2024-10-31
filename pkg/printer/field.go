package printer

type DisplayField struct {
	DisplayName string
	ColumnTag   string
	Value       func() string
}
