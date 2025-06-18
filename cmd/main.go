package main

import (
	"log/slog"
)

func main() {

}

func Handle(err error) {
	if err != nil {
		slog.Error(err.Error())
	}
}
