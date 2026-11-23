// Ukaz dom5 prikaže uporabo paketa vreme z ukaznim vmesnikom.
//
// Uporaba:
//
//	dom5 [--minTemp -50] [--maxTemp 50] [--stMeritev 3]
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/nemke357x/dom5/vreme"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "dom5",
		Usage: "obdelava meritev vremenskih postaj",
		Flags: []cli.Flag{
			&cli.Float64Flag{Name: "minTemp", Value: -50, Usage: "najnižja dovoljena temperatura meritve (°C)"},
			&cli.Float64Flag{Name: "maxTemp", Value: 50, Usage: "najvišja dovoljena temperatura meritve (°C)"},
			&cli.IntFlag{Name: "stMeritev", Value: 3, Usage: "najmanjše število meritev za izračun statistike"},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return zazeni(cmd.Float64("minTemp"), cmd.Float64("maxTemp"), int(cmd.Int("stMeritev")))
		},
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func zazeni(minTemp, maxTemp float64, stMeritev int) error {
	postaje := make(map[string]vreme.Postaja)
	vreme.DodajPostajo(postaje, "LJ01", "Ljubljana", "Ljubljana")
	vreme.DodajPostajo(postaje, "MB01", "Maribor", "Maribor")
	vreme.DodajPostajo(postaje, "CE01", "Celje", "Celje")
	vreme.DodajPostajo(postaje, "KP01", "Koper", "Koper")

	meritve := []struct {
		oznaka string
		temp   float64
	}{
		{"LJ01", 40.5}, {"LJ01", 20.2}, {"LJ01", 30.1}, {"LJ01", 15.5},
		{"MB01", 15.8}, {"MB01", 15.4}, {"MB01", 10.0}, {"MB01", 60.0},
		{"CE01", 3.5}, {"CE01", 2.2}, {"CE01", 4.0},
		{"KP01", 15.0}, {"KP01", 15.5},
		{"TJ01", 20.0},
	}
	for _, m := range meritve {
		if err := vreme.DodajMeritev(postaje, m.oznaka, m.temp, minTemp, maxTemp); err != nil {
			fmt.Printf("%s (%.1f °C): %v\n", m.oznaka, m.temp, err)
		}
	}

	vreme.IzpisMeritev(postaje)
	fmt.Println("POROČILO:")
	vreme.IzpisPorocila(postaje, stMeritev)
	return nil
}
