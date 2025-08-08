package main
var Handler map[string]func(*bufio.Reader) (any, error)