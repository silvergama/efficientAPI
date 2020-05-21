package error_formats

import (
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/silvergama/efficient-api/utils/error_utils"
)

func ParseError(err error) error_utils.MessageErr {
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		if strings.Contains(err.Error(), "no rows in result set") {
			return error_utils.NewNotFoundError("no record matching gived id")
		}
		return error_utils.NewInternalServerError(fmt.Sprintf("error when trying to save message %s", err.Error()))
	}
	switch sqlErr.Number {
	case 1062:
		return error_utils.NewInternalServerError("title already token")
	}
	return error_utils.NewInternalServerError(fmt.Sprintf("error when processing request: %s", err.Error()))
}
