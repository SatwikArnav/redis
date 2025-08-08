package main

import "bufio"

var Handler map[string]func(*bufio.Reader) (any, error)
