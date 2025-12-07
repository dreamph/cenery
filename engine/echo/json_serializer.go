package echo

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/goccy/go-json"
	"github.com/labstack/echo/v4"
)

// fastJSONSerializer uses goccy/go-json for faster encode/decode.
type fastJSONSerializer struct{}

func (fastJSONSerializer) Serialize(c echo.Context, i any, indent string) error {
	var (
		b   []byte
		err error
	)
	if indent != "" {
		b, err = json.MarshalIndent(i, "", indent)
	} else {
		b, err = json.Marshal(i)
	}
	if err != nil {
		return err
	}
	_, err = c.Response().Writer.Write(b)
	return err
}

func (fastJSONSerializer) Deserialize(c echo.Context, i any) error {
	req := c.Request()
	if req.Body == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Request body can't be empty")
	}

	data, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return nil
	}

	req.Body = io.NopCloser(bytes.NewBuffer(data))
	return json.Unmarshal(data, i)
}

// fastBinder replaces JSON decoding with goccy/go-json and falls back to the default binder.
type fastBinder struct {
	echo.DefaultBinder
}

func (b *fastBinder) Bind(i interface{}, c echo.Context) error {
	req := c.Request()
	ct := req.Header.Get(echo.HeaderContentType)
	if strings.HasPrefix(ct, echo.MIMEApplicationJSON) {
		if req.Body == nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Request body can't be empty")
		}

		data, err := io.ReadAll(req.Body)
		if err != nil {
			return err
		}
		if len(data) == 0 {
			return nil
		}

		req.Body = io.NopCloser(bytes.NewBuffer(data))
		if err := json.Unmarshal(data, i); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return nil
	}

	return b.DefaultBinder.Bind(i, c)
}
