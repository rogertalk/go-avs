// Package multipart2 implements MIME multipart parsing for open streams.
// The built-in mime/multipart package hangs until close, which doesn't work
// for longlived streams (see https://github.com/golang/go/issues/15431).
package multipart2

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/textproto"
)

const (
	sExpectingPart = iota
	sInsidePart
	sAfterPart
	sDone
)

const peekBufferSize = 1024

type Part struct {
	Header            textproto.MIMEHeader
	reader            io.Reader
	partReader        *partReader
	disposition       string
	dispositionParams map[string]string
}

func (p *Part) Close() error {
	if p.partReader == nil {
		return nil
	}
	return p.partReader.Close()
}

func (p *Part) FormName() string {
	if p.dispositionParams == nil {
		p.parseContentDisposition()
	}
	if p.disposition != "form-data" {
		return ""
	}
	return p.dispositionParams["name"]
}

func (p *Part) FileName() string {
	if p.dispositionParams == nil {
		p.parseContentDisposition()
	}
	return p.dispositionParams["filename"]
}

func (p *Part) parseContentDisposition() {
	v := p.Header.Get("Content-Disposition")
	var err error
	p.disposition, p.dispositionParams, err = mime.ParseMediaType(v)
	if err != nil {
		p.dispositionParams = make(map[string]string)
	}
}

func newPart(r *Reader) *Part {
	p := new(Part)
	if r.currentPart != nil {
		panic("newPart: expected reader.currentPart to be nil")
	}
	r.currentPart = p
	p.partReader = &partReader{r, p, false}
	rd := bufio.NewReader(p.partReader)
	if header, _ := textproto.NewReader(rd).ReadMIMEHeader(); header != nil {
		p.Header = header
	} else {
		p.Header = textproto.MIMEHeader{}
	}
	p.reader = rd
	return p
}

func (p *Part) Read(d []byte) (n int, err error) {
	return p.reader.Read(d)
}

type partReader struct {
	r          *Reader
	p          *Part
	needsTopUp bool
}

func (pr *partReader) Close() error {
	if pr.r != nil {
		io.Copy(ioutil.Discard, pr)
	}
	return nil
}

func (pr *partReader) Read(d []byte) (n int, err error) {
	if pr.r == nil {
		return 0, io.EOF
	}
	if pr.r.currentPart != pr.p {
		panic("partReader: reading from wrong part")
	}
	if pr.r.w <= pr.r.r || pr.needsTopUp {
		// Need to read more data from the underlying reader.
		err = pr.r.topUp()
		pr.needsTopUp = false
	}
	// Attempt to find a boundary in the available data.
	x := pr.r.r
	for {
		idx := x + bytes.IndexByte(pr.r.buf[x:pr.r.w], pr.r.nlDashBoundary[0])
		if idx < x || idx >= x+len(d) {
			// No boundary beginning possible in the requested data.
			break
		}
		x = idx
		for {
			if pr.r.buf[x] != pr.r.nlDashBoundary[x-idx] {
				// Boundary mismatch.
				break
			} else if x == idx-1+len(pr.r.nlDashBoundary) {
				// Potential boundary match. Check newline/end dashes.
				x += len(pr.r.nl)
				if x < pr.r.w {
					rest := pr.r.buf[x-len(pr.r.nl)+1 : x+1]
					if bytes.Equal(rest, pr.r.nl) {
						pr.r.state = sAfterPart
					} else {
						if len(pr.r.nl) == 1 {
							// We need one more byte.
							x++
						}
						if x < pr.r.w {
							if bytes.Equal(pr.r.buf[x-1:x+1], pr.r.dash) {
								pr.r.state = sDone
							} else {
								// Boundary mismatch.
								break
							}
						}
					}
					if x < pr.r.w {
						pr.r.partsRead++
						pr.r.currentPart = nil
						pr.p.partReader = nil
						n = copy(d, pr.r.buf[pr.r.r:idx])
						pr.r.r += n
						pr.r, pr.p = nil, nil
						return
					}
				}
			}
			x++
			if x >= pr.r.w {
				// Not enough data to determine boundary.
				n = copy(d, pr.r.buf[pr.r.r:idx])
				pr.r.r += n
				pr.needsTopUp = true
				if err == io.EOF {
					err = io.ErrUnexpectedEOF
				}
				return
			}
		}
	}
	// No boundary found; read as much as possible.
	n = copy(d, pr.r.buf[pr.r.r:pr.r.w])
	pr.r.r += n
	if err == io.EOF {
		err = io.ErrUnexpectedEOF
	}
	return
}

type Reader struct {
	reader         io.Reader
	buf            []byte
	partsRead      int
	currentPart    *Part
	state          int
	r, w           int
	dash           []byte
	dashBoundary   []byte
	nl             []byte
	nlDashBoundary []byte
}

func NewReader(r io.Reader, boundary string) *Reader {
	b := []byte("\r\n--" + boundary)
	return &Reader{
		reader:         r,
		buf:            make([]byte, peekBufferSize),
		dash:           b[2:4],
		dashBoundary:   b[2:],
		nl:             b[:2],
		nlDashBoundary: b,
	}
}

func (r *Reader) NextPart() (*Part, error) {
	for {
		if r.state == sDone {
			return nil, io.EOF
		}
		if r.state == sInsidePart {
			r.currentPart.Close()
			continue
		}
		line, err := r.readSlice('\n')
		if err != nil {
			return nil, err
		}
		if r.state == sAfterPart {
			if !bytes.Equal(line, r.nl) {
				return nil, fmt.Errorf("expected newline, got %#v", string(line))
			}
			r.state = sExpectingPart
		}
		if r.state != sExpectingPart {
			return nil, fmt.Errorf("expected state to be %d, was %d", sExpectingPart, r.state)
		}
		if bytes.HasPrefix(line, r.dashBoundary) {
			rest := line[len(r.dashBoundary):]
			if r.partsRead == 0 {
				// Switch to newline mode if first boundary ends with \n instead of \r\n.
				if len(rest) == 1 && rest[0] == '\n' {
					r.nl = r.nl[1:]
					r.nlDashBoundary = r.nlDashBoundary[1:]
				}
			}
			if bytes.Equal(rest, r.nl) {
				r.state = sInsidePart
				return newPart(r), nil
			}
		}
	}
}

func (r *Reader) readSlice(delim byte) (line []byte, err error) {
	for {
		idx := bytes.IndexByte(r.buf[r.r:r.w], delim)
		if idx == -1 {
			if err := r.topUp(); err != nil {
				return nil, err
			}
			continue
		}
		line = r.buf[r.r : r.r+idx+1]
		r.r += idx + 1
		return
	}
}

func (r *Reader) topUp() error {
	if r.r > 0 {
		copy(r.buf, r.buf[r.r:r.w])
		r.w -= r.r
		r.r = 0
	}
	if r.w >= len(r.buf) {
		return fmt.Errorf("can't top up full buffer")
	}
	n, err := r.reader.Read(r.buf[r.w:])
	r.w += n
	return err
}
