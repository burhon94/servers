package pkg

import (
	"bufio"
)

func ReadLine(reader *bufio.Reader) (line string, err error) {
	return reader.ReadString('\n')
}

func WriteLine(string string, writer *bufio.Writer) (err error) {
	_, err = writer.WriteString(string + "\n")
	if err != nil {
		return
	}
	err = writer.Flush()
	if err != nil {
		return
	}
	return
}