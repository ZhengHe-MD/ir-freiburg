package postinglist

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type PostingList struct {
	cap             int
	num             int
	docIDList       []int64
	scoreList       []int64
	skipPointerList []int
}

func NewPostingList() *PostingList {
	return &PostingList{}
}

func (m *PostingList) ReadFromFile(filename string) (err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		err = errors.New("file less than one line")
		return
	}
	firstLine := scanner.Text()
	n, err := strconv.ParseInt(firstLine, 10, 64)
	if err != nil {
		return
	}
	m.Reserve(int(n))

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) != 2 {
			err = errors.New("invalid line format")
			return
		}

		var id, score int64
		if id, err = strconv.ParseInt(parts[0], 10, 64); err != nil {
			return
		}
		if score, err = strconv.ParseInt(parts[1], 10, 64); err != nil {
			return
		}

		m.AddPosting(id, score)
	}
	return
}

func (m *PostingList) ReadFromFileWithSentinel(filename string) (err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		err = errors.New("file less than one line")
		return
	}
	firstLine := scanner.Text()
	n, err := strconv.ParseInt(firstLine, 10, 64)
	if err != nil {
		return
	}
	// for sentinel
	m.Reserve(int(n + 1))

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) != 2 {
			err = errors.New("invalid line format")
			return
		}

		var id, score int64
		if id, err = strconv.ParseInt(parts[0], 10, 64); err != nil {
			return
		}
		if score, err = strconv.ParseInt(parts[1], 10, 64); err != nil {
			return
		}

		m.AddPosting(id, score)
	}

	m.AddPosting(math.MaxInt64, 0)
	return
}

func (m *PostingList) ReadFromFileWithSkipPointer(filename string, numSkipPointers int) (err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		err = errors.New("file less than one line")
		return
	}
	firstLine := scanner.Text()
	n, err := strconv.ParseInt(firstLine, 10, 64)
	if err != nil {
		return
	}
	m.Reserve(int(n))

	gap := int(n) / numSkipPointers
	for scanner.Scan() {
		line := scanner.Text()
		if m.num%gap == 0 && (m.num+gap < int(n)) {
			m.AddSkipPointer(m.num + gap)
		} else {
			m.AddSkipPointer(0)
		}

		parts := strings.Fields(line)
		if len(parts) != 2 {
			err = errors.New("invalid line format")
			return
		}

		var id, score int64
		if id, err = strconv.ParseInt(parts[0], 10, 64); err != nil {
			return
		}
		if score, err = strconv.ParseInt(parts[1], 10, 64); err != nil {
			return
		}

		m.AddPosting(id, score)
	}

	return
}

func (m *PostingList) AddPosting(id, score int64) {
	m.docIDList[m.num] = id
	m.scoreList[m.num] = score
	m.num += 1
}

func (m *PostingList) AddSkipPointer(pos int) {
	m.skipPointerList[m.num] = pos
}

func (m *PostingList) Reserve(n int) {
	m.docIDList = make([]int64, n)
	m.scoreList = make([]int64, n)
	m.skipPointerList = make([]int, n)
	m.cap = n
	m.num = 0
}

func (m *PostingList) Size() int {
	return m.num
}

func (m *PostingList) GetId(idx int) int64 {
	return m.docIDList[idx]
}

func (m *PostingList) GetScore(idx int) int64 {
	return m.scoreList[idx]
}

func (m *PostingList) GetSkipPointer(idx int) (int, bool) {
	return m.skipPointerList[idx], m.skipPointerList[idx] != 0
}

func (m *PostingList) String() string {
	sb := strings.Builder{}
	sb.WriteString("[")
	for i := 0; i < m.num; i++ {
		sb.WriteString(fmt.Sprintf("(%d, %d)", m.GetId(i), m.GetScore(i)))
		if i < m.num-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("]")
	return sb.String()
}

func (m *PostingList) HasSentinel() bool {
	return m.docIDList[m.num-1] == math.MaxInt64
}
