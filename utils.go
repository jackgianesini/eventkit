package eventkit

import (
	"errors"
	log "github.com/sirupsen/logrus"
)

// wrapLogFields : Wrap log fields
func wrapLogFields(fields log.Fields) log.Fields {
	fields["package"] = packageName
	fields["version"] = packageVersion
	fields["commit"] = packageCommit

	return fields
}

// try : try to execute a function
func try(call func()) (err error) {
	defer func() {
		if v := recover(); v != nil {
			switch value := v.(type) {
			case error:
				err = v.(error)
			case string:
				err = errors.New(value)
			}
			return
		}
	}()
	call()
	return
}
