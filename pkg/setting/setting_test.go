package setting

import (
	"testing"
)

func TestSetup(t *testing.T) {
	Setup("../../conf/app.ini")
	PrintSetting()
}
