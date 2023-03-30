package shm

// #include "shm.h"
import "C"
import (
	"encoding/binary"
	"os"
	"unsafe"
)

type VsysFlags int

// reference: https://man7.org/linux/man-pages/man2/shmget.2.html
/*
   The return value of the ftok function, or IPC_PRIVATE, is typically used to generate a unique identifier for a shared memory segment.
   When using IPC_PRIVATE, it is recommended that the two processes that will share the memory segment are related in some way, such as being parent and child processes.
   Otherwise, it may be difficult for another process to obtain the shared memory identifier (return value) generated by the current process.
*/
const (
	StatusIpcNone                = 0               // the constant with the value of 0 that represents no shared memory creation flag
	StatusIpcCreate    VsysFlags = C.IPC_CREAT     // the constant with the value of 512 that represents the flag for creating a newWithReturnId shared memory segment, defined as IPC_CREAT in C language
	StatusIpcExclusive           = C.IPC_EXCL      // the constant with the value of 1024 that represents the flag for creating a newWithReturnId shared memory segment exclusively
	StatusHugePages              = C.SHM_HUGETLB   // the constant with the value of 2048 that represents the flag for requesting shared memory allocation using huge pages.
	StatusNoReserve              = C.SHM_NORESERVE // the constant with the value of C.SHM_NORESERVE that represents the flag for creating a shared memory segment without reserving swap space.
)

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	MajorVersion uint16 = 1
	MinorVersion uint16 = 2
	PatchVersion uint16 = 3
)

// default value for shm
const (
	defautlShmFlag       = StatusIpcCreate | StatusIpcExclusive
	defaultShmPermission = 0600
	defaultMaxKeyValue   = 1 << 10
	defualtMinShmSize    = 2 + 2 + 2 + 8 + 8 + 8 + 4 + 4 + 8 + 4
)

const (
	ErrFailToRetrieveShmSize       = Error("failed to retrieve shm size")
	ErrNegativeOrZeroShmKey        = Error("shm key should not be negative or zero")
	ErrEndOfFile                   = Error("end of file")
	ErrIllegalKey                  = Error("use Illegal key in shm")
	ErrDataDevided                 = Error("data is divided into parts due to lack of space")
	ErrShmIdNotSet                 = Error("shm id not set")
	ErrShmEmptyPoint               = Error("shm point is empty")
	ErrShmAlreadyExist             = Error("shm already exist")
	ErrShmNotExist                 = Error("shm not exist")
	ErrShmFetchInfo                = Error("fetch shm info failed")
	ErrExceedDefaultMaxKeyValue    = Error("exceed default max key value")
	ErrNegativeOrZeroSize          = Error("shm size should not be negative or zero")
	ErrInitializeMajorVersionValue = Error("initialization of major version value failed")
	ErrInitializeMinorVersionValue = Error("initialization of minor version value failed")
	ErrInitializePatchVersionValue = Error("initialization of patch version value failed")
	ErrInitializeKeyValue          = Error("initialization of key value failed")
	ErrInitializeIdValue           = Error("initialization of id value failed")
	ErrInitializeSizeValue         = Error("initialization of size value failed")
	ErrInitializeParameterValue    = Error("initialization of parameter value failed")
	ErrInitializeFlagValue         = Error("initialization of flag value failed")
	ErrInitializeOffsetValue       = Error("initialization of offset value failed")
	ErrInitializeTypeValue         = Error("initialization of type value failed")
)

// VsegmentMap : map key to id
var VsegmentMap map[int64]int64

func init() {
	VsegmentMap = make(map[int64]int64, defaultMaxKeyValue)
}

// Vsegment is a native representation of a SysV shared memory segment
type Vsegment struct {
	key    int64
	id     int64
	size   int64
	offset int64
}

type Vopts struct {
	// These values are user-defined
	Key  int64
	Size int64
	// These values are automatically determined
	flag      VsysFlags
	parameter os.FileMode
}

type Vinfo struct {
	Major     uint16
	Minor     uint16
	Patch     uint16
	Key       int64
	Id        int64
	Size      int64
	Flag      int32
	Parameter [4]int8
	Offset    int64
	Type      int32
}

// >>>>> >>>>> >>>>> [Basic Functions]

/*
newWithReturnId creates a shm shared memory segment based on the given "Vopts" options.
It returns a pointer to a "Vsegment" variable and an error.
*/
func newWithReturnId(opts Vopts) (segment *Vsegment, err error) {
	switch {
	/*
		newWithReturnId creates a new shared memory segment based on the given "Vopts" options
		It returns a pointer to a "Vsegment" variable and an error
	*/
	case opts.Key > 0:
		if opts.Size > 0 {
			/*
				If the "size" field in the "opts" argument is also greater than 0,
				the function sets the "StatusIpcCreate" and "StatusIpcExclusive" flags in the "opts" argument
			*/
			// Set the flag options for creating the shared memory segment with IPC_CREATE and IPC_EXCL
			opts.flag = StatusIpcCreate | StatusIpcExclusive
			// Set the access permissions for the shared memory segment
			opts.parameter = 0600
			// Attempt to create the shared memory segment using the specified options
			segment, err = createShmWithKey(opts)
		} else if opts.Size <= 0 {
			/*
				If the "size" field in the "opts" argument is less than 0,
				the function sets the "err" variable to an "ErrNegativeOrZeroSize" error
			*/
			err = ErrNegativeOrZeroSize
			/*
				shm cannot run properly when the size is 0 or negative, so the entire block of code should be removed.
				If the "size" field in the "opts" argument is 0,
				the function sets the "StatusIpcNone" flag in the "opts" argument
				opts.flag = StatusIpcNone
				opts.parameter = 0600
				segment, err = createShmWithKey(opts)
			*/
		}
	case opts.Key <= 0:
		/*
			If the "key" field in the "opts" argument is less than or equal to 0,
			the function sets the "err" variable to an "ErrNegativeOrZeroShmKey" error
		*/
		err = ErrNegativeOrZeroShmKey
		// <----- no opts.Key == 0 ----->
		/*
			case opts.key == 0:
			I've done a lot of research, and it's not advised to use a key value of 0.
			[reference](https://hackmd.io/@sysprog/linux-shared-memory)
			If the "key" field in the "opts" argument is equal to 0,
			the function calls "createShm" to create a shared memory segment.
			segment, err = createShm(opts)
		*/
	default:
		/*
			If the "key" field in the "opts" argument is equal to 0,
			the function sets the "err" variable to an "ErrIllegalKey" error
		*/
		err = ErrIllegalKey
	}
	return
}

// createShm to create a new shared memory segment with given size
func createShm(opts Vopts) (segment *Vsegment, err error) {
	// Declare variables to store shared memory ID and size
	var shmId C.int
	var shmSize C.ulong

	// Open shared memory segment with given size and default flags and permissions
	shmId, err = C.sysv_shm_open(C.int(opts.Size), C.int(defautlShmFlag), C.int(defaultShmPermission))
	if err == nil {
		// Retrieve the size of the shared memory segment
		shmSize, err = C.sysv_shm_get_size(shmId)

		// Return error if failed to retrieve the size
		if err != nil {
			err = ErrFailToRetrieveShmSize
			return
		}

		// Create a new Vsegment struct to represent the shared memory segment
		segment = &Vsegment{
			id:   int64(shmId),
			size: int64(shmSize),
		}
	}

	// Return the segment and err values
	return
}

// createShmWithKey to create a new shared memory segment with given size by using the key
func createShmWithKey(opts Vopts) (segment *Vsegment, err error) {
	// Declare variables to store shared memory ID and size
	var shmId C.int
	var shmSize C.ulong

	// Open shared memory segment with given size and default flags and permissions by using the key
	shmId, err = C.sysv_shm_open_with_key(C.int(opts.Key), C.int(opts.Size), C.int(opts.flag), C.int(opts.parameter))
	if err == nil {
		// Retrieve the size of the shared memory segment by using the key
		shmSize, err = C.sysv_shm_get_size(shmId)

		// Return error if failed to retrieve the size
		if err != nil {
			err = ErrFailToRetrieveShmSize
			return
		}

		// Create a new Vsegment struct to represent the shared memory segment by using the key
		segment = &Vsegment{
			key:  opts.Key,
			id:   int64(shmId),
			size: int64(shmSize),
		}
	}

	// Return the segment and err values
	return
}

// writeWithId writes data to a shared memory segment and checks if the data to be written exceeds the available space in the segment.
// If the data to write is too large, it reduces the length to the remaining available space and sets a previous error variable.
func (receive *Vsegment) writeWithId(data []byte) (wroteLength int64, err error) {
	// Check if the Vsegment pointer is nil
	if receive == nil {
		err = ErrShmEmptyPoint
	}

	// Check if the shared memory ID has been set
	if receive.id == 0 {
		err = ErrShmIdNotSet
		return
	}

	// Check if the current offset exceeds the size of the shared memory segment
	if receive.offset >= receive.size {
		err = ErrEndOfFile
		return
	}

	// Declare a variable to store an error returned by previous operations
	var previousErr error

	// Determine the length of the data to write
	wroteLength = int64(len(data))
	if (wroteLength + receive.offset) > receive.size {
		/*
			If the data to write exceeds the available space in the shared memory segment,
			reduce the wroteLength to the remaining available space and set the previous error variable to ErrDataDevided
		*/
		wroteLength = receive.size - receive.offset
		previousErr = ErrDataDevided
	}

	// Write the data to the shared memory segment
	_, err = C.sysv_shm_write(C.int(receive.id), unsafe.Pointer(&data[0]), C.int(wroteLength), C.int(receive.offset))
	if err != nil {
		// If an error occurs during the write operation, set the wroteLength to 0 and return
		wroteLength = 0
		return
	}

	// If the write operation is successful, update the offset and check if there was a previous error
	if err == nil {
		receive.offset += wroteLength
		if previousErr != nil {
			err = previousErr
		}
	}

	// Return the wroteLength and any error that occurred during the operation
	return
}

// readWithId is a function that reads data from a shared memory segment.
// It includes error checking and buffer allocation before reading data from shared memory.
func (receive *Vsegment) readWithId(data []byte) (readLength int64, err error) {
	// Check if the Vsegment pointer is nil
	if receive == nil {
		err = ErrShmEmptyPoint
	}

	// Check if the shared memory ID has been set
	if receive.id == 0 {
		err = ErrShmIdNotSet
		return
	}

	// Check if the current offset exceeds the size of the shared memory segment
	if receive.offset >= receive.size {
		readLength = 0
		err = ErrEndOfFile
		return
	}

	// Calculate the actual number of bytes that can be read from the shared memory segment
	length := int64(len(data))
	if (length + receive.offset) > receive.size {
		length = receive.size - receive.offset
	}

	// Allocate a buffer for reading data from shared memory
	buffer := C.malloc(C.size_t(length))
	defer C.free(buffer)

	// Read data from shared memory
	_, err = C.sysv_shm_read(C.int(receive.id), buffer, C.int(length), C.int(receive.offset))
	if err != nil {
		readLength = 0
		return
	}

	// Copy the read data to the output buffer
	count := copy(data, C.GoBytes(buffer, C.int(length)))
	if count > 0 {
		receive.offset += int64(count)
		readLength = int64(count)
		return
	}

	// Return the number of bytes read and the error
	return
}

// deleteWithId defines a Go function that calls an external C function to close a shared memory region.
func (receive *Vsegment) deleteWithId() (err error) {

	// Check if the Vsegment pointer is nil
	if receive == nil {
		err = ErrShmEmptyPoint
		return
	}

	// Check if the shared memory ID has been set
	if receive.id == 0 {
		err = ErrShmIdNotSet
		return
	}

	// Call an external C function `sysv_shm_close` to close the specified shared memory region.
	// The `receive.id` argument is passed as an integer and is cast to a C integer using `C.int`.
	_, err = C.sysv_shm_close(C.int(receive.id))
	// Return any error that occurred during the call to `sysv_shm_close`.
	return
}

// >>>>> >>>>> >>>>> [Extension Function]

/*
NewShm defines a function called "NewShm" which creates a new segment map with specified options.
It writes values to the map and returns an error if any of the writes fail.
The values written include version information, a key, ID, size, and a parameter value.
*/
func NewShm(opts Vopts) (err error) {
	// Check if the value of opts.Key exceeds the default maximum allowed value
	if opts.Key > defaultMaxKeyValue {
		err = ErrExceedDefaultMaxKeyValue
		return
	}

	// Check if the segment already exists in VsegmentMap
	if _, ok := VsegmentMap[opts.Key]; ok {
		err = ErrShmAlreadyExist
		return
	}

	// Create the shared memory segment
	sg, err := newWithReturnId(opts)
	if err != nil {
		return
	}

	// Store the segment ID in VsegmentMap
	VsegmentMap[opts.Key] = sg.id

	// Write major version information to the shared memory segment
	_, err = sg.writeWithId([]byte{
		byte(MajorVersion),
		byte(MajorVersion >> 8),
	})
	if err != nil {
		err = ErrInitializeMajorVersionValue
		return
	}

	// Write minor version information to the shared memory segment
	_, err = sg.writeWithId([]byte{
		byte(MinorVersion),
		byte(MinorVersion >> 8),
	})
	if err != nil {
		err = ErrInitializeMinorVersionValue
		return
	}

	// Write patch version information to the shared memory segment
	_, err = sg.writeWithId([]byte{
		byte(PatchVersion),
		byte(PatchVersion >> 8),
	})
	if err != nil {
		err = ErrInitializePatchVersionValue
		return
	}

	// Write key information to the shared memory segment
	_, err = sg.writeWithId([]byte{
		byte(sg.key),
		byte(sg.key >> 8),
		byte(sg.key >> 16),
		byte(sg.key >> 24),
		byte(sg.key >> 32),
		byte(sg.key >> 40),
		byte(sg.key >> 48),
		byte(sg.key >> 56),
	})
	if err != nil {
		err = ErrInitializeKeyValue
		return
	}

	// Write id information to the shared memory segment
	_, err = sg.writeWithId([]byte{
		byte(sg.id),
		byte(sg.id >> 8),
		byte(sg.id >> 16),
		byte(sg.id >> 24),
		byte(sg.id >> 32),
		byte(sg.id >> 40),
		byte(sg.id >> 48),
		byte(sg.id >> 56),
	})
	if err != nil {
		err = ErrInitializeIdValue
		return
	}

	// Write size information to the shared memory segment
	_, err = sg.writeWithId([]byte{
		byte(sg.size),
		byte(sg.size >> 8),
		byte(sg.size >> 16),
		byte(sg.size >> 24),
		byte(sg.size >> 32),
		byte(sg.size >> 40),
		byte(sg.size >> 48),
		byte(sg.size >> 56),
	})
	if err != nil {
		err = ErrInitializeSizeValue
		return
	}

	// Write flag information to the shared memory segment
	opts.flag = StatusIpcCreate | StatusIpcExclusive
	_, err = sg.writeWithId([]byte{
		byte(opts.flag),
		byte(opts.flag >> 8),
		byte(opts.flag >> 16),
		byte(opts.flag >> 24),
	})
	if err != nil {
		err = ErrInitializeFlagValue
		return
	}

	// Write parameter information to the shared memory segment
	_, err = sg.writeWithId([]byte{0, 6, 0, 0})
	if err != nil {
		err = ErrInitializeParameterValue
		return
	}

	// Write offset information to the shared memory segment
	_, err = sg.writeWithId([]byte{
		byte(sg.offset + 8 + 4),
		byte((sg.offset + 8 + 4) >> 8),
		byte((sg.offset + 8 + 4) >> 16),
		byte((sg.offset + 8 + 4) >> 24),
		byte((sg.offset + 8 + 4) >> 32),
		byte((sg.offset + 8 + 4) >> 40),
		byte((sg.offset + 8 + 4) >> 48),
		byte((sg.offset + 8 + 4) >> 56),
	})
	if err != nil {
		err = ErrInitializeOffsetValue
		return
	}

	// Write type information to the shared memory segment
	/*_, err = sg.writeWithId([]byte{
		byte(),
		byte( >> 8),
		byte( >> 16),
		byte( >> 24),
	})
	if err != nil {
		err = ErrInitializeTypeValue
		return
	}*/

	// Return the error values
	return
}

/*
InfoShm retrieves information about a segment map identified by a key.
It reads data from the map and returns a Vinfo struct containing the major, minor, and patch versions, key, ID, size, flag, and offset.
An error is returned if the map does not exist or if the data could not be read.
*/
func InfoShm(key int64) (vinfo Vinfo, err error) {
	// Check if the value of opts.Key exceeds the default maximum allowed value
	if key > defaultMaxKeyValue {
		err = ErrExceedDefaultMaxKeyValue
		return
	}

	// Check if the given key exists in the VsegmentMap
	id, ok := VsegmentMap[key]
	if !ok {
		err = ErrShmNotExist
		return
	}

	// Create a byte slice to store the raw information and read the information from the shared memory segment into it
	vg := new(Vsegment)
	vg.key = key
	vg.id = id
	vg.offset = 0
	vg.size = defualtMinShmSize

	// Create a byte slice to store the raw information and read the information from the shared memory segment into it
	rawInfo := make([]byte, defualtMinShmSize)
	var count int64
	count, err = vg.readWithId(rawInfo)
	if count != defualtMinShmSize {
		err = ErrShmFetchInfo
		return
	}

	// Use binary.LittleEndian to extract information from the raw byte slice and assign it to the corresponding fields in the Vinfo struct
	vinfo.Major = binary.LittleEndian.Uint16(rawInfo[0:2])           // Extract the Major version number
	vinfo.Minor = binary.LittleEndian.Uint16(rawInfo[2:4])           // Extract the Minor version number
	vinfo.Patch = binary.LittleEndian.Uint16(rawInfo[4:6])           // Extract the Patch version number
	vinfo.Key = int64(binary.LittleEndian.Uint64(rawInfo[6:14]))     // Extract the Key value
	vinfo.Id = int64(binary.LittleEndian.Uint64(rawInfo[14:22]))     // Extract the Id value
	vinfo.Size = int64(binary.LittleEndian.Uint64(rawInfo[22:30]))   // Extract the Size value
	vinfo.Flag = int32(binary.LittleEndian.Uint32(rawInfo[30:34]))   // Extract the flag value
	vinfo.Parameter[0] = int8(rawInfo[34:35][0])                     // Extract the Parameter value
	vinfo.Parameter[1] = int8(rawInfo[35:36][0])                     // Extract the Parameter value
	vinfo.Parameter[2] = int8(rawInfo[36:37][0])                     // Extract the Parameter value
	vinfo.Parameter[3] = int8(rawInfo[37:38][0])                     // Extract the Parameter value
	vinfo.Offset = int64(binary.LittleEndian.Uint64(rawInfo[38:46])) // Extract the Offset value
	vinfo.Type = int32(binary.LittleEndian.Uint32(rawInfo[46:50]))   // Extract the type value

	// Return the extracted Vinfo struct
	return
}

func WriteOffset(key, offset int64) (vinfo Vinfo, err error) {
	// Check if the value of opts.Key exceeds the default maximum allowed value
	if key > defaultMaxKeyValue {
		err = ErrExceedDefaultMaxKeyValue
		return
	}

	// Check if the given key exists in the VsegmentMap
	id, ok := VsegmentMap[key]
	if !ok {
		err = ErrShmNotExist
		return
	}

	// Create a byte slice to store the raw information and read the information from the shared memory segment into it
	vg := new(Vsegment)
	vg.key = key
	vg.id = id
	vg.offset = 2 + 2 + 2 + 8 + 8 + 8 + 4 + 4

	vg.size = defualtMinShmSize

	// Write offset information to the shared memory segment
	_, err = vg.writeWithId([]byte{
		byte(offset),
		byte(offset >> 8),
		byte(vg.offset >> 16),
		byte(vg.offset >> 24),
		byte(vg.offset >> 32),
		byte(vg.offset >> 40),
		byte(vg.offset >> 48),
		byte(vg.offset >> 56),
	})
	if err != nil {
		err = ErrInitializeOffsetValue
		return
	}

	return
}
