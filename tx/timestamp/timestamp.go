package timestamp

import (
	"bytes"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
)

// if unix time's nano seconds is 0, RFC3339Nano Format tails nano parts.
const RFC3339NanoFixed = "2006-01-02T15:04:05.000000000Z"

type Timestamp timestamp.Timestamp

func (t *Timestamp) String() string {
	return t.UTC().Format(RFC3339NanoFixed)
}

// UTC returns time.Time with the location set to UTC
func (t *Timestamp) UTC() time.Time {
	return time.Unix(t.Seconds, int64(t.Nanos)).UTC()
}

// MarshalJSON marshals Timestamp as RFC3339NanoFixed format
func (t *Timestamp) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer([]byte{'"'})
	if _, err := buf.WriteString(t.String()); err != nil {
		return nil, err
	}
	if err := buf.WriteByte('"'); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UnmarshalJSON unmarshals RFC3339NanoFixed format bytes to Timestamp
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	str := string(data)
	if str == "null" {
		return nil
	}
	txtime, err := time.Parse(RFC3339NanoFixed, strings.Trim(str, `"`))
	if err != nil {
		return err
	}
	// t.Seconds = txtime.Unix()
	t = &Timestamp{Seconds: txtime.Unix(), Nanos: int32(txtime.Nanosecond())}
	return nil
}
