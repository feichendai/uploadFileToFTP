package saveJson

import (
	"testing"
)

func TestGetJsonInfo(t *testing.T) {

	_, err := GetJsonInfo()
	if (err != nil)   {
		t.Errorf("GetJsonInfo() error = %v", err )
		return
	}
}
