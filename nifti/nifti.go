package nifti

import (
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"reflect"
	"strings"
)

type Nifti1Header struct {
	SizeofHdr    int32    /*!< MUST be 348           */ /* int32 sizeof_hdr;      */
	DataType     [10]byte /*!< ++UNUSED++            */ /* char data_type[10];  */
	DbName       [18]byte /*!< ++UNUSED++            */ /* char db_name[18];    */
	Extents      int32    /*!< ++UNUSED++            */ /* int32 extents;         */
	SessionError int16    /*!< ++UNUSED++            */ /* short session_error; */
	Regular      byte     /*!< ++UNUSED++            */ /* char regular;        */
	DimInfo      byte     /*!< MRI slice ordering.   */ /* char hkey_un0;       */

	/*--- was image_dimension substruct ---*/
	Dim      [8]int16 /*!< Data array dimensions.*/ /* short dim[8];        */
	IntentP1 float32  /*!< 1st intent parameter. */ /* short unused8;       */
	/* short unused9;       */
	IntentP2 float32 /*!< 2nd intent parameter. */ /* short unused10;      */
	/* short unused11;      */
	IntentP3 float32 /*!< 3rd intent parameter. */ /* short unused12;      */
	/* short unused13;      */
	IntentCode    int16      /*!< NIFTI_INTENT_* code.  */ /* short unused14;      */
	Datatype      int16      /*!< Defines data type!    */ /* short datatype;      */
	Bitpix        int16      /*!< Number bits/voxel.    */ /* short bitpix;        */
	SliceStart    int16      /*!< First slice index.    */ /* short dim_un0;       */
	Pixdim        [8]float32 /*!< Grid spacings.        */ /* float32 pixdim[8];     */
	VoxOffset     float32    /*!< Offset into .nii file */ /* float32 vox_offset;    */
	SclSlope      float32    /*!< Data scaling: slope.  */ /* float32 funused1;      */
	SclInter      float32    /*!< Data scaling: offset. */ /* float32 funused2;      */
	SliceEnd      int16      /*!< Last slice index.     */ /* float32 funused3;      */
	SliceCode     byte       /*!< Slice timing order.   */
	XyztUnits     byte       /*!< Units of pixdim[1..4] */
	CalMax        float32    /*!< Max display intensity */ /* float32 cal_max;       */
	CalMin        float32    /*!< Min display intensity */ /* float32 cal_min;       */
	SliceDuration float32    /*!< Time for 1 slice.     */ /* float32 compressed;    */
	Toffset       float32    /*!< Time axis shift.      */ /* float32 verified;      */
	Glmax         int32      /*!< ++UNUSED++            */ /* int32 glmax;           */
	Glmin         int32      /*!< ++UNUSED++            */ /* int32 glmin;           */

	/*--- was data_history substruct ---*/
	Descrip [80]byte /*!< any text you like.    */ /* char descrip[80];    */
	AuxFile [24]byte /*!< auxiliary filename.   */ /* char aux_file[24];   */

	QformCode int16 /*!< NIFTI_XFORM_* code.   */ /*-- all ANALYZE 7.5 ---*/
	SformCode int16 /*!< NIFTI_XFORM_* code.   */ /*   fields below here  */
	/*   are replaced       */
	QuaternB float32 /*!< Quaternion b param.   */
	QuaternC float32 /*!< Quaternion c param.   */
	QuaternD float32 /*!< Quaternion d param.   */
	QoffsetX float32 /*!< Quaternion x shift.   */
	QoffsetY float32 /*!< Quaternion y shift.   */
	QoffsetZ float32 /*!< Quaternion z shift.   */

	SrowX [4]float32 /*!< 1st row affine transform.   */
	SrowY [4]float32 /*!< 2nd row affine transform.   */
	SrowZ [4]float32 /*!< 3rd row affine transform.   */

	IntentName [16]byte /*!< 'name' or meaning of data.  */

	Magic [4]byte /*!< MUST be "ni1\0" or "n+1\0". */

} /**** 348 bytes total ****/

func (header *Nifti1Header) LoadHeader(filepath string) {
	reader, err := gzipOpen(filepath)
	defer reader.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = binary.Read(reader, binary.LittleEndian, header)
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Printf("%+v\n", header)
}

type Nifti1Image struct { /*!< Image storage struct **/
	ndim     int32    /*!< last dimension greater than 1 (1..7) */
	nx       int32    /*!< dimensions of grid array             */
	ny       int32    /*!< dimensions of grid array             */
	nz       int32    /*!< dimensions of grid array             */
	nt       int32    /*!< dimensions of grid array             */
	nu       int32    /*!< dimensions of grid array             */
	nv       int32    /*!< dimensions of grid array             */
	nw       int32    /*!< dimensions of grid array             */
	dim      [8]int32 /*!< dim[0]=ndim, dim[1]=nx, etc.         */
	nvox     uint     /*!< number of voxels = nx*ny*nz*...*nw   */
	nbyper   int32    /*!< bytes per voxel, matches datatype    */
	datatype int32    /*!< type of data in voxels: DT_* code    */

	dx     float32    /*!< grid spacings      */
	dy     float32    /*!< grid spacings      */
	dz     float32    /*!< grid spacings      */
	dt     float32    /*!< grid spacings      */
	du     float32    /*!< grid spacings      */
	dv     float32    /*!< grid spacings      */
	dw     float32    /*!< grid spacings      */
	pixdim [8]float32 /*!< pixdim[1]=dx, etc. */

	sclSlope float32 /*!< scaling parameter - slope        */
	sclInter float32 /*!< scaling parameter - intercept    */

	calMin float32 /*!< calibration parameter, minimum   */
	calMax float32 /*!< calibration parameter, maximum   */

	qformCode int32 /*!< codes for (x,y,z) space meaning  */
	sformCode int32 /*!< codes for (x,y,z) space meaning  */

	freqDim  int32 /*!< indexes (1,2,3, or 0) for MRI    */
	phaseDim int32 /*!< directions in dim[]/pixdim[]     */
	sliceDim int32 /*!< directions in dim[]/pixdim[]     */

	sliceCode     int32   /*!< code for slice timing pattern    */
	sliceStart    int32   /*!< index for start of slices        */
	sliceEnd      int32   /*!< index for end of slices          */
	sliceDuration float32 /*!< time between individual slices   */

	/*! quaternion transform parameters
	  [when writing a dataset, these are used for qform, NOT qto_xyz]   */
	quaternB, quaternC, quaternD,
	qoffsetX, qoffsetY, qoffsetZ,
	qfac float32

	qtoXyz [4][4]float32 /*!< qform: transform (i,j,k) to (x,y,z) */
	qtoIjk [4][4]float32 /*!< qform: transform (x,y,z) to (i,j,k) */

	stoXyz [4][4]float32 /*!< sform: transform (i,j,k) to (x,y,z) */
	stoIjk [4][4]float32 /*!< sform: transform (x,y,z) to (i,j,k) */

	toffset float32 /*!< time coordinate offset */

	xyzUnits  int32 /*!< dx,dy,dz units: NIFTI_UNITS_* code  */
	timeUnits int32 /*!< dt       units: NIFTI_UNITS_* code  */

	niftiType int32 /*!< 0==ANALYZE, 1==NIFTI-1 (1 file),
	  2==NIFTI-1 (2 files),
	  3==NIFTI-ASCII (1 file) */
	intentCode int32    /*!< statistic type (or something)       */
	intentP1   float32  /*!< intent parameters                   */
	intentP2   float32  /*!< intent parameters                   */
	intentP3   float32  /*!< intent parameters                   */
	intentName [16]byte /*!< optional description of intent data */

	descrip [80]byte /*!< optional text to describe dataset   */
	auxFile [24]byte /*!< auxiliary filename                  */

	fname       string      /*!< header filename (.hdr or .nii)         */
	iname       string      /*!< image filename  (.img or .nii)         */
	inameOffset int32       /*!< offset into iname where data starts    */
	swapsize    int32       /*!< swap unit in image data (might be 0)   */
	byteorder   int32       /*!< byte order on disk (MSB_ or LSB_FIRST) */
	data        []byte      /*!< pointer to data: nbyper*nvox bytes     */
	volumeN     int         //defined by me, volume vox num
	byte2floatF interface{} //defined by me, byte2floatF
	float2byteF interface{}
	header      Nifti1Header //defined by me, it might be a good idea store the img header in the image structure

	numExt int32 /*!< number of extensions in ext_list       */
	// nifti1_extension       *ext_list        /*!< array of extension structs (with data) */
	// analyze_75_orient_code analyze75_orient /*!< for old analyze files, orient */
}

//filepath ,read data?
//header.VoxOffset, header.Bitpix,header.Dim are very important
func (img *Nifti1Image) LoadImage(filepath string, rdata bool) {
	var header Nifti1Header
	header.LoadHeader(filepath)
	img.header = header

	// set dimensions of data array
	img.ndim, img.dim[0] = int32(header.Dim[0]), int32(header.Dim[0])
	img.nx, img.dim[1] = int32(header.Dim[1]), int32(header.Dim[1])
	img.ny, img.dim[2] = int32(header.Dim[2]), int32(header.Dim[2])
	img.nz, img.dim[3] = int32(header.Dim[3]), int32(header.Dim[3])
	img.nt, img.dim[4] = int32(header.Dim[4]), int32(header.Dim[4])
	img.nu, img.dim[5] = int32(header.Dim[5]), int32(header.Dim[5])
	img.nv, img.dim[6] = int32(header.Dim[6]), int32(header.Dim[6])
	img.nw, img.dim[7] = int32(header.Dim[7]), int32(header.Dim[7])
	img.nvox = 1
	for i := int16(1); i <= header.Dim[0]; i++ {
		img.nvox *= uint(header.Dim[i])
	}
	img.volumeN = int(img.dim[1] * img.dim[2] * img.dim[3])

	if header.Bitpix == 0 {
		fmt.Println("empty header.Bitpix")
		fmt.Println(header)
		return
	}

	//init byte2float32 function
	img.nbyper = int32(header.Bitpix) / 8

	// fmt.Println(header)
	//setting function to convert float2byte or byte2float
	if img.nbyper == 2 {
		img.byte2floatF = func(b []byte) float32 {
			v := binary.LittleEndian.Uint16(b)
			return float32(v)
		}
		img.float2byteF = func(buff []byte, x float32) {
			binary.LittleEndian.PutUint16(buff, uint16(x))
		}

	} else if img.nbyper == 4 {
		img.byte2floatF = func(b []byte) float32 {
			v := binary.LittleEndian.Uint32(b)
			return math.Float32frombits(v)
		}
		img.float2byteF = func(buff []byte, x float32) {
			v := math.Float32bits(x)
			binary.LittleEndian.PutUint32(buff, v)
		}

	} else if img.nbyper == 8 {
		img.byte2floatF = func(b []byte) float32 {
			v := binary.LittleEndian.Uint64(b)
			return float32(math.Float64frombits(v))
		}
		img.float2byteF = func(buff []byte, x float32) {
			v := math.Float64bits(float64(x))
			binary.LittleEndian.PutUint64(buff, v)
		}
	} else {
		fmt.Println("input nbyper:", img.nbyper)
		panic("(img *Nifti1Image) byte2float, only support 16 32 and 64 bit")
	}
	fmt.Println(img.nbyper)
	// set the grid spacings
	img.dx, img.pixdim[1] = header.Pixdim[1], header.Pixdim[1]
	img.dy, img.pixdim[2] = header.Pixdim[2], header.Pixdim[2]
	img.dz, img.pixdim[3] = header.Pixdim[3], header.Pixdim[3]
	img.dt, img.pixdim[4] = header.Pixdim[4], header.Pixdim[4]
	img.du, img.pixdim[5] = header.Pixdim[5], header.Pixdim[5]
	img.dv, img.pixdim[6] = header.Pixdim[6], header.Pixdim[6]
	img.dw, img.pixdim[7] = header.Pixdim[7], header.Pixdim[7]

	if rdata {
		reader, err := gzipOpen(filepath)
		defer reader.Close()
		if err != nil {
			fmt.Println("open data error")
			return
		}
		data, err := ioutil.ReadAll(reader)
		if err != nil {
			fmt.Println("read data error")
			return
		}
		img.data = data[uint(header.VoxOffset):len(data)]
		// fmt.Println(len(img.data), len(data), uint(header.VoxOffset))
	}
}

// x,y,z,t,start at zero
func (img *Nifti1Image) GetAt(x, y, z, t int) float32 {

	tIndex := int32(t) * img.nx * img.ny * img.nz
	zIndex := img.nx * img.ny * int32(z)
	yIndex := img.nx * int32(y)
	xIndex := int32(x)
	index := int32(tIndex) + zIndex + yIndex + xIndex
	return img.byte2float(img.data[index*img.nbyper : (index+1)*img.nbyper]) //shift byte
}

func (img *Nifti1Image) SetAt(x, y, z, t int, elem float32) {

	tIndex := int32(t) * img.nx * img.ny * img.nz
	zIndex := img.nx * img.ny * int32(z)
	yIndex := img.nx * int32(y)
	xIndex := int32(x)
	index := int32(tIndex) + zIndex + yIndex + xIndex
	copy(img.data[index*img.nbyper:(index+1)*img.nbyper], img.float2byte(elem))
}

func (img *Nifti1Image) GetTimeSeries(x, y, z int) []float32 {
	var timeSeries []float32
	timeNum := len(img.data) / int(img.volumeN*int(img.nbyper))
	for i := 0; i < timeNum; i++ {
		timeSeries = append(timeSeries, img.GetAt(x, y, z, i))
	}
	return timeSeries
}

func (img *Nifti1Image) GetSlice(z, t int) [][]float32 {
	sliceX := int(img.nx)
	sliceY := int(img.ny)
	slice := make([][]float32, sliceX)
	for i := range slice {
		slice[i] = make([]float32, sliceY)
	}
	for xi := 0; xi < sliceX; xi++ {
		for yi := 0; yi < sliceY; yi++ {
			slice[xi][yi] = img.GetAt(xi, yi, z, t)
			// fmt.Println(xi, yi, z, t)
		}
	}
	return slice
}

//Return [x,y,z,t]
func (img *Nifti1Image) GetDims() [4]int {
	var dim [4]int
	for i, v := range img.dim[1:5] {
		dim[i] = int(v)
	}
	return dim
}

//convert byte to float32,init in  LoadImage
func (img *Nifti1Image) byte2float(data []byte) float32 {
	v := reflect.ValueOf(img.byte2floatF).Call([]reflect.Value{reflect.ValueOf(data)})[0].Interface()
	return v.(float32)
}

//convert float32 to byte,init in  LoadImage
func (img *Nifti1Image) float2byte(data float32) []byte {
	buff := make([]byte, img.nbyper)
	reflect.ValueOf(img.float2byteF).Call([]reflect.Value{reflect.ValueOf(buff), reflect.ValueOf(data)})
	return buff
}

func gzipOpen(filepath string) (io.ReadCloser, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	filepathS := strings.Split(filepath, ".")
	if filepathS[len(filepathS)-1] == "gz" {
		gzipReader, err := gzip.NewReader(f)
		if err != nil {
			return nil, err
		}
		return gzipReader, nil
	}
	return f, nil
}
