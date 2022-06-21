// Wordle! Solver
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/minipub/wordle/internal"
)

func main() {
	conn, err := net.Dial("tcp", "0.0.0.0:8080")
	if err != nil {
		os.Exit(1)
	}

	r := bufio.NewReaderSize(conn, 1)
	w := bufio.NewWriterSize(conn, 1)

	var iWord [5]byte // only words
	copy(iWord[:], "great")

	for {
		rs := readLoop(r, func() {
			r.Reset(conn)
		})

		// print server response
		// fmt.Printf("{{ %s }}", rs)
		fmt.Printf("%s", rs)

		if IsTheEnd(rs) {
			break
		}

		var pos [5]int
		if !IsTheStart(rs) {
			var vs []byte
			for i, j, k := 0, 0, 0; ; {
				v := rs[i]

				if j > 0 && j%internal.ColoredByteNum == 0 {
					// fmt.Printf("vs: %+v\n", vs)

					if bytes.HasPrefix(vs, []byte(internal.ColorRed)) {
						pos[k] = internal.Miss
					} else if bytes.HasPrefix(vs, []byte(internal.ColorYellow)) {
						pos[k] = internal.Appear
					} else if bytes.HasPrefix(vs, []byte(internal.ColorGreen)) {
						pos[k] = internal.Hit
					}
					k++

					// if len(pos) == 5 {
					if k == 5 {
						break
					} else {
						vs = make([]byte, 0)
						i += (internal.ColorResetByteNum + 1)
						j = 0
					}
				} else {
					vs = append(vs, v)
					i++
					j++
				}
			}
		}

		// fmt.Printf("pos: %+v\n", pos)

		if bytes.HasSuffix(rs, []byte(internal.Prompt)) {
			// fmt.Fprintln(os.Stderr, "The second guess")
			time.Sleep(time.Second)
			// solve word
			if !IsTheStart(rs) {
				iWord = internal.SolveWord(pos, iWord)
			}
			// print client request
			fmt.Printf("%s\n", iWord)
			w.Write(iWord[:])
		}

		// fmt.Println("1111111111")

	}

}

func IsTheStart(b []byte) bool {
	return bytes.HasPrefix(b, []byte(internal.PreText))
}

func IsTheEnd(b []byte) bool {
	return bytes.HasSuffix(b, []byte(internal.ByeText))
}

// read next input or the end
func readLoop(r io.Reader, f func()) (rs []byte) {
	var keepRead bool
	var n int
	var err error
	var b [512]byte

	defer f()

	for {
		if keepRead {
			fmt.Printf("keepRead: %t\n", keepRead)
			n, err = r.Read(b[n:])
		} else {
			n, err = r.Read(b[:])
		}
		if err != nil {
			fmt.Printf("readLoop err: %+v\n", err)
			os.Exit(2)
		}

		rs = b[0:n]
		fmt.Fprintf(os.Stderr, "resp b: {{ %+v }}, {{ %s }}\n", rs, rs)

		if IsTheEnd(rs) {
			break
		} else if !bytes.HasSuffix(rs, []byte(internal.Prompt)) {
			// continue to read if Prompt not direct after PreText or Colored Response
			fmt.Fprintln(os.Stderr, "Prompt not afterwards!")
			keepRead = true
			continue
		} else {
			break
		}
		// fmt.Println("2222222222")

		// if IsTheEnd(rs) {
		// 	break
		// } else if IsTheStart(rs) ||
		// 	bytes.HasPrefix(rs, []byte(internal.ColorRed)) ||
		// 	bytes.HasPrefix(rs, []byte(internal.ColorYellow)) ||
		// 	bytes.HasPrefix(rs, []byte(internal.ColorGreen)) {

		// 	// continue to read if Prompt not direct after PreText or Colored Response
		// 	if !bytes.HasSuffix(rs, []byte(internal.Prompt)) {
		// 		fmt.Fprintln(os.Stderr, "Prompt not afterwards!")
		// 		keepRead = true
		// 		continue
		// 	} else {
		// 		break
		// 	}
		// }
		// // fmt.Println("2222222222")
	}

	return
}
