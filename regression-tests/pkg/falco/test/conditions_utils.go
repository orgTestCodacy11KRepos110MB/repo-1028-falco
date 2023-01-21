package test

import (
	"bufio"
	"io"
	"time"
)

func skewedDuration(d time.Duration) time.Duration {
	return time.Duration(float64(d) * 1.10)
}

func readLineByLine(r io.Reader, f func(string) error) error {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		if err := f(scanner.Text()); err != nil {
			return err
		}
	}
	return scanner.Err()
}
