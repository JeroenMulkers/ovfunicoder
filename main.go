package main

import (
	"fmt"
	"github.com/mumax/3/data"
	"github.com/mumax/3/draw"
	"github.com/mumax/3/oommf"
	"os"
)

var (
	// maximum width of unicode figure
	// TODO: ideally this should depend on the terminal width
	MaxNx = 64

	X, Y, Z = 0, 1, 2
)

func main() {

	ovfFilenames := os.Args[1:]
	if len(ovfFilenames) < 1 {
		fmt.Println("You need to specify the path of at least one ovf file")
		return
	}

	for _, ovfFilename := range ovfFilenames {
		field, _, err := oommf.ReadFile(ovfFilename)
		if err != nil {
			fmt.Println(err, "\n")
			continue
		}
		fmt.Println(ovfFilename)
		ShowLayer(field, 0)
		fmt.Print("\n")
	}
}

// Showlayer prints a unicode figure of a slice in the stdout. This is done
// by mimicing a grid of pixels with upper half block unicode characters ▀
// and ANSI escape sequences to colorize the upper part (foreground color)
// and the lower part (background color)
func ShowLayer(rawField *data.Slice, layer int) {

	// resample so that the figure is not too wide
	field := Resample(rawField)

	// loop over the unicode character cells
	// each character cell contains two pixels on top of each other
	for i := 0; i < (field.Size()[Y]+1)/2; i++ {
		for j := 0; j < field.Size()[X]; j++ {

			// reset colors
			asciEscape := "\u001b[0m"

			// set foreground color (the upper half of the cell)
			ix := j
			iy := 2 * i
			mx := float32(field.Get(X, ix, iy, layer))
			my := float32(field.Get(Y, ix, iy, layer))
			mz := float32(field.Get(Z, ix, iy, layer))
			fc := draw.HSLMap(mx, my, mz)
			asciEscape += fmt.Sprintf("\u001b[38;2;%d;%d;%dm", fc.R, fc.G, fc.B)

			// set background color (the lower half of the cell)
			ix = j
			iy = 2*i + 1
			if iy < field.Size()[1] {
				mx = float32(field.Get(X, ix, iy, layer))
				my = float32(field.Get(Y, ix, iy, layer))
				mz = float32(field.Get(Z, ix, iy, layer))
				bc := draw.HSLMap(mx, my, mz)
				asciEscape += fmt.Sprintf("\u001b[48;2;%d;%d;%dm", bc.R, bc.G, bc.B)
			}

			// TODO: avoid printing each character seperatly (use buffers)
			fmt.Print(asciEscape + "▀")
		}
		fmt.Print("\u001b[0m\n") // reset colors before going to a new line
	}
}

func Resample(slice *data.Slice) *data.Slice {
	N := slice.Size() // old size
	if N[X] <= MaxNx {
		return slice
	}

	Nnew := N                         // new size
	Nnew[X] = MaxNx                   // make sure that Nx <= MaxNx
	Nnew[Y] = (N[Y] * Nnew[X]) / N[X] // and that the aspect ratio is more or less conserved
	if Nnew[Y] < 1 {                  // and that there is at least one row of cells
		Nnew[Y] = 1
	}
	return data.Resample(slice, Nnew)
}
