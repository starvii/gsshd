// refer to https://github.com/jpillora/sshd-lite/blob/master/main.go
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"github.com/jpillora/sshd-lite/server"
	"os/exec"
	"path/filepath"
)

var VERSION string = "0.0.1" //set via ldflags

var help = `
	Usage: gsshd [options] <auth>
	Version: ` + VERSION + `
	Options:
	  --host, listening interface (defaults to all)
	  --port -p, listening port (defaults to 22, then fallsback to 2200)
	  --shell, the type of to use shell for remote sessions (defaults to bash)
	  --keyfile, a filepath to an private key (for example, an 'id_rsa' file)
	  --keyseed, a string to use to seed key generation
	  --version, display version
	  -v, verbose logs
	<auth> must be set to one of:
	  1. a username and password string separated by a colon ("user:pass")
	  2. a path to an ssh authorized keys file ("~/.ssh/authorized_keys")
	  3. "none" to disable client authentication :WARNING: very insecure
	Notes:
	  * if no keyfile and no keyseed are set, a random RSA2048 key is used
	  * once authenticated, clients will have access to a shell of the
	  current user. sshd-lite does not lookup system users.
	  * sshd-lite only supports remotes shells. tunnelling and command
	  execution are not currently supported.
	Read more: https://github.com/jpillora/sshd-lite
`

func parse_parameters() sshd.Config {
	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, help)
		os.Exit(1)
	}

	//init config from flags
	c := &sshd.Config{}
	flag.StringVar(&c.Host, "host", "0.0.0.0", "")
	flag.StringVar(&c.Port, "p", "", "")
	flag.StringVar(&c.Port, "port", "", "")
	flag.StringVar(&c.Shell, "shell", "", "")
	flag.StringVar(&c.KeyFile, "keyfile", "", "")
	flag.StringVar(&c.KeySeed, "keyseed", "", "")
	flag.BoolVar(&c.LogVerbose, "v", false, "")

	//help/version
	h1f := flag.Bool("h", false, "")
	h2f := flag.Bool("help", false, "")
	vf := flag.Bool("version", false, "")
	flag.Parse()

	if *vf {
		_, _ = fmt.Fprintf(os.Stderr, VERSION)
		os.Exit(0)
	}
	if *h1f || *h2f {
		flag.Usage()
	}

	args := flag.Args()
	if len(args) == 1 {
		c.AuthType = args[0]
	} else {
		c.AuthType = "0:Mechrev0"
	}
	return *c
}

func run_in_background() {
	if os.Getppid() != 1 {
		exe, _ := filepath.Abs(os.Args[0])
		cmd := exec.Command(exe, os.Args[1:]...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		_ = cmd.Start()
		log.Printf("Start background process: %v\n", cmd.Process.Pid)
		os.Exit(cmd.Process.Pid)
	}
}

func main() {
	c := parse_parameters()

	run_in_background()
	s, err := sshd.NewServer(&c)
	if err != nil {
		log.Fatal(err)
	}
	err = s.Start()
	if err != nil {
		log.Fatal(err)
	}
}

