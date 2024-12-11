package helper

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func GetClaimsFromToken(c echo.Context) (jwt.MapClaims, error) {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		fmt.Println(ok)
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Error Fetching Token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Token Claims Error")
	}

	return claims, nil
}

// CustomDate is a custom type to represent the date in DD-MM-YYYY format
type CustomDate struct {
	time.Time
}

// UnmarshalJSON is a custom unmarshaler for the CustomDate type
func (d *CustomDate) UnmarshalJSON(b []byte) error {
	// Define the format we expect
	const layout = "02-01-2006"
	// Parse the date using the layout
	parsedTime, err := time.Parse(`"`+layout+`"`, string(b))
	if err != nil {
		return err
	}
	d.Time = parsedTime
	return nil
}

// MarshalJSON is a custom marshaler for the CustomDate type
func (d CustomDate) MarshalJSON() ([]byte, error) {
	// Format the time in DD-MM-YYYY
	return json.Marshal(d.Format("02-01-2006"))
}

// Value implements the Valuer interface for CustomDate
// This method is used when writing data to the database
func (d CustomDate) Value() (driver.Value, error) {
	// Return the date as a string in the correct format
	return d.Format("02-01-2006"), nil
}

// Scan implements the Scanner interface for CustomDate
// This method is used when reading data from the database
func (d *CustomDate) Scan(value interface{}) error {
	// Handle NULL values
	if value == nil {
		d.Time = time.Time{}
		return nil
	}

	// Expect the value to be a string in DD-MM-YYYY format
	strValue, ok := value.(string)
	if !ok {
		return errors.New("invalid type for CustomDate")
	}

	// Parse the date string into the CustomDate type
	const layout = "02-01-2006"
	parsedTime, err := time.Parse(layout, strValue)
	if err != nil {
		return err
	}
	d.Time = parsedTime
	return nil
}
