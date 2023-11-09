package cio

import (
	"io"
	"sync"
)

func Pipe(src io.ReadWriteCloser, dst io.ReadWriteCloser) (int64, int64) {
	var sent, received int64
	// var isHTTP bool
	var wg sync.WaitGroup
	var o sync.Once
	close := func() {
		src.Close()
		dst.Close()
	}
	wg.Add(2)

	// go func() {
	// 	// Read the initial data from src to determine if it's HTTP
	// 	buf := make([]byte, 10240) // Adjust the buffer size as needed
	// 	n, err := src.Read(buf)
	// 	if err == nil {
	// 		data := string(buf[:n])
	// 		if strings.HasPrefix(data, "GET /") {
	// 			fmt.Println("isHTTPisHTTP", isHTTP)
	// 			fmt.Println("datadata", data)
	// 			isHTTP = true
	// 		} else {
	// 			fmt.Println("isHTTPisHTTP", 111)
	// 		}
	// 	}
	// }()

	go func() {
		received, _ = io.Copy(src, dst)
		// fmt.Printf("Received from %v\n", src)
		o.Do(close)
		wg.Done()
	}()
	go func() {
		sent, _ = io.Copy(dst, src)
		// fmt.Printf("Sent to %v\n", dst)
		o.Do(close)
		wg.Done()
	}()
	wg.Wait()
	return sent, received
}

// const vis = false

// type pipeVisPrinter struct {
// 	name string
// }

// func (p pipeVisPrinter) Write(b []byte) (int, error) {
// 	log.Printf(">>> %s: %x", p.name, b)
// 	return len(b), nil
// }

// func pipeVis(name string, r io.Reader) io.Reader {
// 	if vis {
// 		return io.TeeReader(r, pipeVisPrinter{name})
// 	}
// 	return r
// }
