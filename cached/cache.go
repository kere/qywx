package cached

import "github.com/kere/gno"

var (
	expiresVal = 0
)

// Expires int
func Expires() int {
	if expiresVal > 0 {
		return expiresVal
	}

	expiresVal = gno.GetConfig().GetConf("data").DefaultInt("data_expires", 72000)
	return expiresVal
}
