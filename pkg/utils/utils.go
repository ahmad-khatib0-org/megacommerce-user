package utils

type Utils struct{}

type UtilsArgs struct{}

func NewUtils(ua *UtilsArgs) (*Utils, error) {
	u := &Utils{}
	return u, nil
}
