package error_formats

import (
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/silvergama/efficientAPI/utils/errorutils"
)

func ParseError(err error) errorutils.MessageErr {
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		if strings.Contains(err.Error(), "no rows in result set") {
			return errorutils.NewNotFoundError("no record matching gived id")
		}
		return errorutils.NewInternalServerError(fmt.Sprintf("error when trying to save message %s", err.Error()))
	}
	switch sqlErr.Number {
	case 1062:
		return errorutils.NewInternalServerError("title already token")
	}
	return errorutils.NewInternalServerError(fmt.Sprintf("error when processing request: %s", err.Error()))
}
