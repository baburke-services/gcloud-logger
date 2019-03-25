package main

import (
    "encoding/json"
    "log"
    "os/exec"
)

var DEFAULT_COMMAND = []string{"journalctl", "-f"}

func create_reader(cmd []string) chan interface{} {
    objects := make(chan interface{}, 25);
    _cmd := exec.Cmd{};

    if len(cmd) == 0 {
        cmd = DEFAULT_COMMAND;
    }

    path, err := exec.LookPath(cmd[0]);
    if err != nil {
        log.Fatal(err);
    }

    _cmd.Args = cmd;
    _cmd.Path = path;

    output, err := _cmd.StdoutPipe();
    if err != nil {
        log.Fatal(err);
    }

    if err := _cmd.Start(); err != nil {
        log.Fatal(err)
    }

    defer _cmd.Wait();

    go func() {
        decoder := json.NewDecoder(output);
        var v interface{};

        for decoder.More() {
            err := decoder.Decode(&v);
            if err != nil {
                log.Println(err);
                break;
            }
            objects <- v;
        }

        close(objects);
    }();

    return objects;
}
