package main

import (
	"fmt"
	"log"
	"os"
	"regexp"

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

type Options struct {
	Interface    string `short:"i" long:"interface" description:"Interface to monitor" value-name:"eth0"`
	Port         int    `short:"p" long:"port" description:"DNS port" default:"53" value-name:"53"`
	Database     string `short:"d" long:"database" description:"Database file path" value-name:"FILE"`
	Filter       string `short:"F" long:"filter" description:"Filter by question" default:"" value-name:"*.com"`
	Daemon       bool   `short:"D" long:"daemon" description:"Run DNAS in daemon mode"`
	Write        string `short:"w" long:"write" description:"Write JSON output to log file" value-name:"FILE"`
	User         string `short:"u" long:"user" description:"Drop privileges to this user" value-name:"USER"`
	Hexdump      bool   `short:"H" long:"hexdump" description:"Show hexdump of DNS packet"`
	FindQuestion string `short:"q" long:"find-question" description:"Search for DNS record by question" value-name:"STRING"`
	FindAnswer   string `short:"a" long:"find-answer" description:"Search for DNS records by answer data" value-name:"STRING"`
	List         bool   `short:"l" long:"list" description:"List all seen DNS questions"`
	Version      bool   `short:"v" long:"version" description:"Show version information"`
}

func Usage(p *flags.Parser) {
	fmt.Printf("\n  %s (%s) - %s\n",
		NAME,
		VERSION,
		DESCRIPTION,
	)

	p.WriteHelp(os.Stdout)
	fmt.Printf("\n")
	os.Exit(1)
}

func Version() {
	fmt.Printf("%s - %s - Version: %s\n",
		NAME,
		DESCRIPTION,
		VERSION,
	)

	os.Exit(1)
}

func CLIRun(f func(options *Options)) {

	options := &Options{}

	var parser = flags.NewParser(options, flags.Default)

	if _, err := parser.Parse(); err != nil {
		Usage(parser)
	}

	if options.Version {
		Version()
	}

	if options.Database == "" {
		options.Database = DATABASE
	}

	if options.List {
		listAllQuestions(options.Database)
		os.Exit(1)
	} else if options.FindQuestion != "" {
		r, err := regexp.Compile(options.FindQuestion)

		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}

		findByQuestion(options.Database, r)
		os.Exit(1)
	} else if options.FindAnswer != "" {
		findByAnswer(options.Database, []byte(options.FindAnswer))
		os.Exit(1)
	}

	if options.Interface == "" {
		Usage(parser)
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
