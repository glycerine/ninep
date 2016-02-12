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
	var user ninep.User
	var err error
	var c *clnt.Clnt
	var file *clnt.File
	var d []*ninep.Dir

	flag.Parse()
	user = ninep.OsUsers.Uid2User(os.Geteuid())
	clnt.DefaultDebuglevel = *debuglevel
	c, err = clnt.Mount("tcp", *addr, "", uint32(*msize), user)
	if err != nil {
		log.Fatal(err)
	}

	lsarg := "/"
	if flag.NArg() == 1 {
		lsarg = flag.Arg(0)
	} else if flag.NArg() > 1 {
		log.Fatal("error: only one argument expected")
	}

	file, err = c.FOpen(lsarg, ninep.OREAD)
	if err != nil {
		log.Fatal(err)
	}

	for {
		d, err = file.Readdir(0)
		if d == nil || len(d) == 0 || err != nil {
			break
		}

		for i := 0; i < len(d); i++ {
			os.Stdout.WriteString(d[i].Name + "\n")
		}
	}

	file.Close()
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}

	return
}
