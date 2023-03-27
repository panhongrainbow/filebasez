package shm

// #include "shm.h"
import "C"
import "os"

type SharedMemoryFlags int

// reference: https://man7.org/linux/man-pages/man2/shmget.2.html
/*
   The return value of the ftok function, or IPC_PRIVATE, is typically used to generate a unique identifier for a shared memory segment.
   When using IPC_PRIVATE, it is recommended that the two processes that will share the memory segment are related in some way, such as being parent and child processes.
   Otherwise, it may be difficult for another process to obtain the shared memory identifier (return value) generated by the current process.
*/
const (
	StatusIpcNone                        = 0               // the constant with the value of 0 that represents no shared memory creation flag
	StatusIpcCreate    SharedMemoryFlags = C.IPC_CREAT     // the constant with the value of 512 that represents the flag for creating a new shared memory segment, defined as IPC_CREAT in C language
	StatusIpcExclusive                   = C.IPC_EXCL      // the constant with the value of 1024 that represents the flag for creating a new shared memory segment exclusively
	StatusHugePages                      = C.SHM_HUGETLB   // the constant with the value of 2048 that represents the flag for requesting shared memory allocation using huge pages.
	StatusNoReserve                      = C.SHM_NORESERVE // the constant with the value of C.SHM_NORESERVE that represents the flag for creating a shared memory segment without reserving swap space.
)

type Error string

func (e Error) Error() string {
	return string(e)
}

// default value for shm
const (
	defautlShmFlag       = StatusIpcCreate | StatusIpcExclusive
	defaultShmPermission = 0600
)

const (
	ErrFailToRetrieveShmSize = Error("failed to retrieve shm size")
	ErrNegativeShmKey        = Error("shm key should not be negative")
)

// Vsegment is a native representation of a SysV shared memory segment
type Vsegment struct {
	Key    int64
	Id     int64
	Size   int64
	offset int64
}

type Vopts struct {
	key       int64
	size      int64
	flag      SharedMemoryFlags
	parameter os.FileMode
}

/*
New creates a new shared memory segment based on the given "Vopts" options.
It returns a pointer to a "Vsegment" variable and an error.
*/
func New(opts Vopts) (segment *Vsegment, err error) {
	switch {
	/*
		New creates a new shared memory segment based on the given "Vopts" options
		It returns a pointer to a "Vsegment" variable and an error
	*/
	case opts.key > 0:
		if opts.size > 0 {
			/*
				If the "size" field in the "opts" argument is also greater than 0,
				the function sets the "StatusIpcCreate" and "StatusIpcExclusive" flags in the "opts" argument
			*/
			opts.flag = StatusIpcCreate | StatusIpcExclusive
		} else {
			/*
				If the "size" field in the "opts" argument is 0,
				the function sets the "StatusIpcNone" flag in the "opts" argument
			*/
			opts.flag = StatusIpcNone
		}
		opts.parameter = 0600
		segment, err = createShmWithKey(opts)
		return
	case opts.key < 0:
		/*
			If the "key" field in the "opts" argument is less than 0,
			the function sets the "err" variable to an "ErrNegativeShmKey" error
		*/
		err = ErrNegativeShmKey
		return
	case opts.key == 0:
		/*
			If the "key" field in the "opts" argument is equal to 0,
			the function calls "createShm" to create a shared memory segment
		*/
		segment, err = createShm(opts)
		return
	default:
		// If none of the above conditions are met, the function does nothing and returns
	}
	return
}

// createShm to create a new shared memory segment with given size
func createShm(opts Vopts) (segment *Vsegment, err error) {
	// Declare variables to store shared memory ID and size
	var shmId C.int
	var shmSize C.ulong

	// Open shared memory segment with given size and default flags and permissions
	shmId, err = C.sysv_shm_open(C.int(opts.size), C.int(defautlShmFlag), C.int(defaultShmPermission))
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
			Id:   int64(shmId),
			Size: int64(shmSize),
		}
	}

	// Return the segment and err values
	return
}

// createShm to create a new shared memory segment with given size by using the key
func createShmWithKey(opts Vopts) (segment *Vsegment, err error) {
	// Declare variables to store shared memory ID and size
	var shmId C.int
	var shmSize C.ulong

	// Open shared memory segment with given size and default flags and permissions by using the key
	shmId, err = C.sysv_shm_open_with_key(C.int(opts.key), C.int(opts.size), C.int(opts.flag), C.int(opts.parameter))
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
			Key:  opts.key,
			Id:   int64(shmId),
			Size: int64(shmSize),
		}
	}

	// Return the segment and err values
	return
}
