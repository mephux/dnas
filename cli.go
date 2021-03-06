package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"runtime"

	"code.google.com/p/gopass"

	"github.com/sevlyar/go-daemon"
	"github.com/visionmedia/go-flags"

	"os/user"
	"strconv"
	"syscall"
)

func chuser(username string) (uid, gid int) {
	usr, err := user.Lookup(username)

	if err != nil {
		fmt.Printf("failed to find user %q: %s", username, err)
	}

	uid, err = strconv.Atoi(usr.Uid)

	if err != nil {
		fmt.Printf("bad user ID %q: %s", usr.Uid, err)
	}

	gid, err = strconv.Atoi(usr.Gid)

	if err != nil {
		fmt.Printf("bad group ID %q: %s", usr.Gid, err)
	}

	if err := syscall.Setgid(gid); err != nil {
		fmt.Printf("setgid(%d): %s", gid, err)
	}

	if err := syscall.Setuid(uid); err != nil {
		fmt.Printf("setuid(%d): %s", uid, err)
	}

	return uid, gid
}

// Options cli command options
type Options struct {
	Interface string `short:"i" long:"interface" description:"Interface to monitor" value-name:"eth0"`
	Port      int    `short:"p" long:"port" description:"DNS port" default:"53" value-name:"53"`
	Daemon    bool   `short:"D" long:"daemon" description:"Run DNAS in daemon mode"`
	Write     string `short:"w" long:"write" description:"Write JSON output to log file" value-name:"FILE"`
	User      string `short:"u" long:"user" description:"Drop privileges to this user" value-name:"USER"`
	Hexdump   bool   `short:"H" long:"hexdump" description:"Show hexdump of DNS packet"`

	Mysql    bool `long:"mysql" description:"Enable Mysql Output Support"`
	Postgres bool `long:"postgres" description:"Enable Postgres Output Support"`
	Sqlite3  bool `long:"sqlite3" description:"Enable Sqlite3 Output Support"`

	DbUser         string `long:"db-user" description:"Database User" value-name:"root" default:"root"`
	DbPassword     string `long:"db-password" description:"Database Password" value-name:"PASSWORD"`
	DbDatabase     string `long:"db-database" description:"Database Database" value-name:"dnas" default:"dnas"`
	DbHost         string `long:"db-host" description:"Database Host" value-name:"localhost"`
	DbPort         string `long:"db-port" description:"Database Port" value-name:"3306"`
	DbPath         string `long:"db-path" description:"Path to Database on disk. (sqlite3 only)" default:"./dnas.db"`
	DbSsl          bool   `long:"db-ssl" description:"Enable TLS / SSL encrypted connection to the database. (mysql/postgres only)" value-name:"false"`
	DbFlush        bool   `long:"db-flush" description:"Flush all data from the database and start fresh"`
	DbSkipVerify   bool   `long:"db-skip-verify" description:"Allow Self-signed or invalid certificate (mysql/postgres only)" value-name:"false"`
	DatabaseOutput bool   `long:"db-verbose" description:"Show database logs in STDOUT"`

	Quiet   bool `short:"q" long:"quiet" description:"Suppress DNAS output"`
	Version bool `short:"v" long:"version" description:"Show version information"`

	// Other
	InterfaceData *net.Interface
	Hostname      string
	Ip            string
	Client        *Client
}

func printUsage(p *flags.Parser) {
	fmt.Printf("\n  %s (%s) - %s\n",
		Name,
		Version,
		Description,
	)

	p.WriteHelp(os.Stdout)
	fmt.Printf("\n")
	os.Exit(1)
}

func printVersion() {
	fmt.Printf("%s - %s - Version: %s\n",
		Name,
		Description,
		Version,
	)

	os.Exit(1)
}

// CLIRun start DNAS and process all command-line options
func CLIRun(f func(options *Options)) {

	runtime.GOMAXPROCS(runtime.NumCPU())

	options := &Options{}

	var parser = flags.NewParser(options, flags.Default)

	options.Hostname, _ = os.Hostname()

	addrs, _ := net.InterfaceAddrs()

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				options.Ip = ipnet.IP.String()
			}
		}
	}

	if _, err := parser.Parse(); err != nil {
		printUsage(parser)
	}

	if options.Version {
		printVersion()
	}

	if options.Mysql {

		if options.Postgres || options.Sqlite3 {
			fmt.Println("DNAS Error: Only one database output plugin can be selected.")
			os.Exit(1)
		}

		if options.DbHost != "" {

			options.DbHost = "tcp(" + options.DbHost + ":"

			if options.DbPort == "" {
				options.DbPort = "3306"
			}

			options.DbHost = options.DbHost + options.DbPort + ")"
		}

	} else if options.Postgres {

		if options.Sqlite3 {
			fmt.Println("DNAS Error: Only one database output plugin can be selected.")
			os.Exit(1)
		}

		if options.DbPort == "" {
			options.DbPort = "5432"
		}

		if options.DbHost == "" {
			options.DbHost = "127.0.0.1"
		}
	}

	if options.Mysql || options.Postgres {
		if options.DbPassword == "" {
			password, err := gopass.GetPass("Database Password: ")

			if err != nil {
				panic(err)
			}

			options.DbPassword = password
		} else if options.DbPassword == "none" {
			options.DbPassword = ""
		}
	}

	if options.DbUser == "" {
		options.DbUser = "root"
	}

	if options.DbDatabase == "" {
		options.DbDatabase = "dnas"
	}

	if options.Interface == "" {
		printUsage(parser)
	} else {

		iface, err := net.InterfaceByName(options.Interface)

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		options.InterfaceData = iface
	}

	if options.Daemon {

		cntxt := &daemon.Context{
			PidFileName: "dnas.pid",
			PidFilePerm: 0644,
			LogFileName: "dnas.log",
			LogFilePerm: 0640,
			WorkDir:     "./",
			Umask:       027,
		}

		d, err := cntxt.Reborn()

		if err != nil {
			log.Fatalln(err)
		}

		if d != nil {
			return
		}

		defer cntxt.Release()

		go f(options)

		err = daemon.ServeSignals()

		if err != nil {
			log.Println("Error:", err)
		}
	} else {
		f(options)
	}
}
