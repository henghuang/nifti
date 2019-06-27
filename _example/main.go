package main

import (
	"fmt"

	"github.com/henghuang/nifti"
)

func main() {
	var test nifti.Nifti1Image
	test.LoadImage("./MNI152.nii.gz", true)   //open a nifti image
	fmt.Println(test.GetAt(50, 50, 50, 0))    //get value at x,y,z,t; index start at 0
	test.SetAt(50, 50, 50, 0, 100)            //set value 100 at x,y,z,t; index start at 0
	test.Save("save.nii")                     //save nifti image and compress it
	newImg := nifti.NewImg(100, 100, 100, 20) //create an empyt image
	fmt.Println(newImg.GetDims())             //print dimensions
}
