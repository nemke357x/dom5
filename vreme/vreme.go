// Package vreme omogoča hranjenje in obdelavo meritev temperature
// vremenskih postaj.
//
// Postaje se hranijo v slovarju tipa map[string]Postaja, kjer je ključ
// oznaka postaje (npr. "LJ01"). Ker so polja strukture Postaja neizvožena,
// se nove postaje ustvarijo s funkcijo DodajPostajo:
//
//	postaje := make(map[string]vreme.Postaja)
//	vreme.DodajPostajo(postaje, "LJ01", "Ljubljana", "Ljubljana")
//	vreme.DodajMeritev(postaje, "LJ01", 21.5, -50, 50)
//	vreme.IzpisPorocila(postaje, 3)
package vreme

import (
	"errors"
	"fmt"
)

// Napake, ki jih vrača funkcija DodajMeritev.
var (
	// ErrIzvenObmocja pomeni, da je temperatura izven dovoljenega območja.
	ErrIzvenObmocja = errors.New("temperatura je izven dovoljenega območja")
	// ErrNiPostaje pomeni, da postaja z dano oznako ne obstaja.
	ErrNiPostaje = errors.New("postaja ne obstaja")
)

// Postaja predstavlja vremensko postajo z imenom, krajem in
// zaporedjem izmerjenih temperatur v stopinjah Celzija.
type Postaja struct {
	ime     string
	kraj    string
	meritev []float64
}

// DodajPostajo v slovar postaje doda novo postajo z oznako oznaka,
// imenom ime in krajem kraj. Nova postaja nima meritev.
// Če postaja z isto oznako že obstaja, jo nadomesti.
func DodajPostajo(postaje map[string]Postaja, oznaka, ime, kraj string) {
	postaje[oznaka] = Postaja{ime: ime, kraj: kraj}
}

// DodajMeritev postaji z oznako oznaka doda meritev temperature.
// Meritev je sprejeta le, če velja minTemp <= temperatura <= maxTemp;
// sicer vrne ErrIzvenObmocja. Če postaja ne obstaja, vrne ErrNiPostaje.
func DodajMeritev(postaje map[string]Postaja, oznaka string, temperatura, minTemp, maxTemp float64) error {
	if temperatura < minTemp || temperatura > maxTemp {
		return ErrIzvenObmocja
	}
	postaja, ok := postaje[oznaka]
	if !ok {
		return ErrNiPostaje
	}
	postaja.meritev = append(postaja.meritev, temperatura)
	postaje[oznaka] = postaja
	return nil
}

// statistika vrne najnižjo, najvišjo in povprečno temperaturo postaje
// z oznako oznaka. Vrednost ok je false, če postaja ne obstaja ali
// ima manj kot stMeritev meritev.
func statistika(postaje map[string]Postaja, oznaka string, stMeritev int) (najnizja, najvisja, povprecje float64, ok bool) {
	postaja, obstaja := postaje[oznaka]
	if !obstaja || len(postaja.meritev) == 0 || len(postaja.meritev) < stMeritev {
		return 0, 0, 0, false
	}
	var vsota float64
	for i, temp := range postaja.meritev {
		if i == 0 || temp < najnizja {
			najnizja = temp
		}
		if i == 0 || temp > najvisja {
			najvisja = temp
		}
		vsota += temp
	}
	povprecje = vsota / float64(len(postaja.meritev))
	return najnizja, najvisja, povprecje, true
}

// IzpisMeritev na standardni izhod izpiše vse postaje in njihove meritve.
func IzpisMeritev(postaje map[string]Postaja) {
	fmt.Println("MERITVE:")
	for oznaka, postaja := range postaje {
		fmt.Printf("%s - %s (%s): %v\n", oznaka, postaja.ime, postaja.kraj, postaja.meritev)
	}
}

// IzpisPorocila za vsako postajo izpiše najnižjo, najvišjo in povprečno
// temperaturo ter oceno vremena (Vroče, Toplo, Hladno ali Mrzlo).
// Postaje z manj kot stMeritev meritvami so označene s "premalo meritev".
func IzpisPorocila(postaje map[string]Postaja, stMeritev int) {
	for oznaka, postaja := range postaje {
		najnizja, najvisja, povprecje, ok := statistika(postaje, oznaka, stMeritev)
		if !ok {
			fmt.Printf("%s: premalo meritev\n", postaja.ime)
			continue
		}
		var ocena string
		switch {
		case povprecje >= 25:
			ocena = "Vroče"
		case povprecje >= 15:
			ocena = "Toplo"
		case povprecje >= 5:
			ocena = "Hladno"
		default:
			ocena = "Mrzlo"
		}
		fmt.Printf("%s: min %.1f °C, max %.1f °C, povprečje %.1f °C -> %s\n",
			postaja.ime, najnizja, najvisja, povprecje, ocena)
	}
}
