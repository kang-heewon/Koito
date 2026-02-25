package catalog

type MbzMappingError struct {
	Message string
	Entity  string // "artist", "album", "track"
	MbzID   string
}

func (e *MbzMappingError) Error() string {
	return e.Message
}
