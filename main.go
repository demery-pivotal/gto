package main
import (
	"fmt"
	"bufio"
	"io"
	"path"
	"os"
)

func main() {
	file, err := os.Open("output.bin")
	handleError(err)
	defer file.Close()

	reader := bufio.NewReader(file)
	serializer := KryoSerializer{reader}

	for ;; {
		event, err := NewEvent(serializer)
		handleError(err)
		WriteEvent(event)
	}
}

func handleError(e error) {
	if e == nil {
		return
	}
	fmt.Fprintln(os.Stderr, e)
	os.Exit(1)
}

func WriteEvent(e Event) {
	classDir := fmt.Sprintf("class-%05d", e.ClassId)
	methodDir := fmt.Sprintf("method-%05d", e.MethodId)
	dirPath := path.Join("results", classDir, methodDir)
	err := os.MkdirAll(dirPath, 0750)
	handleError(err)

	filePath := path.Join(dirPath, e.File)
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	handleError(err)
	defer f.Close()

	_, err = f.WriteString(e.Text)
	handleError(err)
}

type Event struct {
	ClassId int
	MethodId int
	File string
	Text string
}

func  NewEvent(s KryoSerializer) (Event, error) {
	isStdout, err := s.readBool()
	if err != nil {
		return Event{}, err
	}
	var file string
	if isStdout {
		file = "stdout.log"
	} else {
		file = "stderr.log"
	}
	classId, err := s.readInt()
	if err != nil {
		return Event{}, err
	}
	methodId, err := s.readInt()
	if err != nil {
		return Event{}, err
	}
	text, err := s.readText()
	if err != nil {
		return Event{}, err
	}
	return Event{
		ClassId: classId,
		MethodId: methodId,
		File: file,
		Text: text,
	}, nil
}

type KryoSerializer struct {
	reader *bufio.Reader
}

func (s KryoSerializer) readBool() (bool, error) {
	b, err := s.reader.ReadByte()
	return b != 0, err
}

func (s KryoSerializer) readInt() (int, error) {
	var shift = 0
	var value int
	var done bool


	for ; !done; {
		b, err := s.reader.ReadByte()
		if err != nil {
			return 0, err
		}
		done = b & 0x80 == 0
		byteValue := (int) (b & 0x7f)
		value |= byteValue << shift
		shift += 7
	}
	return value, nil
}

func (s KryoSerializer) readText() (string, error) {
	textLen, err := s.readInt()
	if err != nil {
		return "", err
	}
	buf := make([]byte, textLen)
	_, err = io.ReadFull(s.reader, buf)
	if err != nil {
		return "", err
	}
	return string(buf), err
}
