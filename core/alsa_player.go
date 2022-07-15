package core

import (
	"os/exec"
	"strconv"
	"syscall"
)

type AlsaPlayer struct {
	playCmd *exec.Cmd
}

func NewAlsaPlayer() AlsaPlayer {
	player := AlsaPlayer{}
	return player
}

func (player *AlsaPlayer) Play(audioFile string, duration int) error {
	player.playCmd = exec.Command("/usr/bin/aplay", audioFile, "-d", strconv.Itoa(duration))
	e := player.playCmd.Run()
	player.playCmd = nil
	return e
}

func (player *AlsaPlayer) Stop() {
	if player.playCmd != nil {
		player.playCmd.Process.Signal(syscall.SIGTERM)
	}
}
