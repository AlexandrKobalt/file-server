package duration

import (
	"encoding/json"
	"time"
)

type Minutes struct {
	time.Duration
}

func (d *Minutes) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Minutes())
}

func (d *Minutes) UnmarshalJSON(b []byte) error {
	var minutes int64
	if err := json.Unmarshal(b, &minutes); err != nil {
		return err
	}

	d.Duration = time.Duration(minutes * int64(time.Minute))

	return nil
}
