package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"runtime/pprof"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/oscisn93/tetris/tetris"
)

var cpuProfile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
  flag.Parse()

  if *cpuProfile != "" {
    file, err := os.Create(*cpuProfile)

    if err != nil {
      log.Fatal(err)
    }

    writer := bufio.NewWriter(file)

    if err := pprof.StartCPUProfile(writer); err != nil {
      log.Fatal(err)
    }

    defer func() {
      if err := writer.Flush(); err != nil {
        log.Fatal(err)
      }
    }()

    defer pprof.StopCPUProfile()
  }
  
	ebiten.SetWindowSize(tetris.ScreenWidth, tetris.ScreenHeight)
	ebiten.SetWindowTitle("Vimtris (Vim Keybindings for Tetris)")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(&tetris.Game{}); err != nil {
		log.Fatal(err)
	}
}
