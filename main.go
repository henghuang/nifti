package main

import (
	_ "image/jpeg"

	"./nifti"
)

func main() {
	var x nifti.Nifti1Image
	x.LoadImage("MNI152.nii.gz", true)
	y := nifti.NewImg(100, 100, 100, 10)
	y.SetAt(50, 50, 50, 0, 1.1)
	y.Save("test.nii")
}
