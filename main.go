package main

import (
	"image"
	"image/color"
	"image/jpeg"
	_ "image/jpeg"
	"log"
	"os"

	"./nifti"
)

func main() {
	var x nifti.Nifti1Image
	x.LoadImage("MNI152.nii.gz", true)
	sliceTest := x.GetSlice(40, 0)
	// fmt.Println(len(sliceTest), len(sliceTest[0])) //start from 0
	// fmt.Println(x.GetTimeSeries(50, 50, 50))
	img := image.NewGray16(image.Rect(0, 0, len(sliceTest), len(sliceTest[0])))
	// fmt.Println(sliceTest)
	for i, row := range sliceTest {
		for j, col := range row {
			img.SetGray16(i, j, color.Gray16{Y: uint16(col)})
		}
	}
	f, err := os.Create("test.jpg")
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	err = jpeg.Encode(f, img, &jpeg.Options{Quality: 100})
	if err != nil {
		log.Fatal(err)
	}

}
