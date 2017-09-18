package devices

// DbInserter interface is used to model database ORMs that are able to perform SQL INSERT statements.
type DbInserter interface {
	Insert(list ...interface{}) error
}
