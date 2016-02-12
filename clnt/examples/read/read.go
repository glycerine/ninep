package main

import (
	"github.com/rminnich/ninep"
	"github.com/rminnich/ninep/clnt"
	"flag"
	"io"
	"log"
	"os"
)

var debuglevel = flag.Int("d", 0, "debuglevel")
var addr = flag.String("addr", "127.0.0.1:5640", "network address")
var msize = flag.Uint("m", 8192, "Msize for 9p")

func main() {
	var n int
	var user ninep.User
	var err error
	var c *clnt.Clnt
	var file *clnt.File
	var buf []byte

	flag.Parse()
	user = ninep.OsUsers.Uid2User(os.Geteuid())
	clnt.DefaultDebuglevel = *debuglevel
	c, err = clnt.Mount("tcp", *addr, "", uint32(*msize), user)
	if err != nil {
		goto error
	}

	if flag.NArg() != 1 {
		log.Println("invalid arguments")
		return
	}

	file, err = c.FOpen(flag.Arg(0), ninep.OREAD)
	if err != nil {
		goto error
	}

	buf = make([]byte, 8192)
	for {
		n, err = file.Read(buf)
		if n == 0 {
			break
		}

		os.Stdout.Write(buf[0:n])
	}

	file.Close()

	if err != nil && err != io.EOF {
		goto error
	}

	return

error:
	log.Println("Error", err)
}
