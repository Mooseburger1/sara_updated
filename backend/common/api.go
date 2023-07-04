package common

import "strconv"

// SaraServer is ...
type SaraServer interface {
	StartServer()
	ShutdownServer()
}

// str2Int32 is a package private helper function
// for type conversion
func Str2Int32(val string) (int32, error) {
	i, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}

	return int32(i), nil
}
