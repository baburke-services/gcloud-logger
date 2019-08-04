package main

import (
	"flag"
	"fmt"
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

	log.Printf("read cursor %q", cursor)

	reader := glogger.NewLogReader(journald.Reader)
	log_channel := reader.StartStream(0)
	post_cursor, err := glogger.CursorProcessor(log_channel)
	if err != nil {
		panic(err)
	}

	for entry := range post_cursor {
		fmt.Printf("%s: %s\n", entry.LevelName, entry.Message)
	}

	if err := journald.Close(); err != nil {
		panic(err)
	}
}

// vim: noexpandtab
