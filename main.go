package main

import (
	"flag"
	"log"
	"os"
	"source.baburke.net/baburke-services/gcloud-logger/glogger"
)

type Args struct {
	follow bool
}

func parse_args(argv []string) *Args {
	var args Args

	set := flag.NewFlagSet(argv[0], flag.ExitOnError)
	set.BoolVar(&args.follow, "follow", false, "follow journald?")
	set.Parse(argv[1:])

	log.Printf("follow? %v", args.follow)

	return &args
}

func main() {
	args := parse_args(os.Args)
	cursor, err := glogger.ReadCurrentCursor()
	if err == glogger.ERROR_NO_CURSOR {
		log.Print("cursor not found; proceeding")
	} else if err != nil {
		panic(err)
	}

	journald, err := glogger.NewJournaldReader(cursor, args.follow)
	if err != nil {
		panic(err)
	}
	defer journald.Close()

	log.Printf("read cursor %q", cursor)

	reader := glogger.NewLogReader(journald.Reader)
	log_channel := reader.StartStream(0)
	post_cursor, err := glogger.CursorProcessor(log_channel)
	if err != nil {
		panic(err)
	}

	glog, gclient, err := glogger.NewGLogger()
	if err != nil {
		panic(err)
	}
	defer gclient.Close()
	defer glog.Flush()

	done, err := glogger.StartLogForwarder(glog, post_cursor)
	<-done
}

// vim: noexpandtab
