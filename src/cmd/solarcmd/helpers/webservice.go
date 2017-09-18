package helpers

import (
	"github.com/go-gorp/gorp"
)

// WebService is used to sent measurement records stored temporarily in the local db to the remote web service.
type WebService struct {
	dbMap *gorp.DbMap
}

// NewWebService is used to create a new WebService object.
func NewWebService() (ws *WebService, err error) {
	ws = &WebService{}
	err = ws.init()
	return
}

func (ws *WebService) init() (err error) {
	return
}

// UpdateWebService retrieves data from the db and pushes it to the web service
func (ws *WebService) UpdateWebService() (err error) {
	return
}
