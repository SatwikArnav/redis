package store

var Set map[string]string

var Hset map[string]map[string]string
func init() {
	Set = make(map[string]string)
	Hset = make(map[string]map[string]string)
}
