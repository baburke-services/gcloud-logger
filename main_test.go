package main

import (
    "testing"
)

func TestParseArgs(t *testing.T) {
    t.Run("NoFollow", func(t *testing.T) {
        args := parse_args([]string{"exec"})
        if args.follow != false {
            t.Error("follow is not false")
        }
    })
    t.Run("Follow", func(t *testing.T) {
        args := parse_args([]string{"exec", "-follow"})
        if args.follow != true {
            t.Error("follow is not true")
        }
    })
}
