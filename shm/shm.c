#include "shm.h"

/*
    When writing processes, sometimes you may want certain values or data structures to be passed or modified between different processes.
    However, since data structures are stored in virtual memory, memory blocks between different processes are independent and
    cannot access memory locations of other processes.
    To achieve the functionality of accessing the same memory block between different processes, shared memory must be used.
*/

/*
    sysv_shm_open is called to create a shared memory segment.
    It takes three arguments:
        - size:  the size of the shared memory segment
        - flags: the flags used to create the shared memory segment, such as IPC_CREAT or IPC_EXCL
        - perm:  the permissions of the shared memory segment, represented as an octal value
*/
int sysv_shm_open(int size, int flags, int perm) {
    int shm_id;
    /*
        If the size argument is non-zero,
        the function creates a new shared memory segment with the specified size and permissions
    */
    if(size) {
        // Unless otherwise specified, segment is owner-read/write (no exec)
        if(!perm){
            perm = 0600;
        }
        return shmget(IPC_PRIVATE, size, flags|perm);
    } else {
       return shmget(IPC_PRIVATE, size, 0);
    }
}

// sysv_shm_open_with_key creates or opens a shared memory segment with the given key, size, and permissions.
int sysv_shm_open_with_key(int key,int size, int flags, int perm) {
    int shm_id;

    // If the size parameter is non-zero, it creates a new shared memory segment with the given key, size, and permissions
    if(size) {
        // Unless otherwise specified, segment is owner-read/write (no exec)
        if(!perm){
            perm = 0600;
        }
        return shmget((key_t)key, size, flags|perm);
    } else {
        return shmget(key, size, 0);
    }
}

/*
    sysv_shm_write writes data to a shared memory segment identified by an ID, offset, and length,
    after attaching to it and detaching from it
    It takes three arguments:
        - shm_id:  the integer that represents the ID of the shared memory segment that we want to write to
        - input:   the pointer to the memory location that contains the data that we want to write to the shared memory segment
        - len:     the integer that represents the length of the data that we want to write to the shared memory segment
        - offset:  the integer that represents the offset from the beginning of the shared memory segment where we want to start writing the data
*/
int sysv_shm_write(int shm_id, void* input, int len, int offset) {
    // Attach to the given segment to get its memory address
    char* addr = sysv_shm_attach(shm_id);

    // Check if attaching to the segment was successful
    if(addr == (char*)(-1)){
        return -1;
    }

    // Copy len bytes from input into the shared memory segment starting from offset
    memcpy(addr+offset, input, len);

    // Detach from the shared memory segment
    sysv_shm_detach(addr);

    // Return 0 to indicate success
    return 0;
}

/*
   sysv_shm_attach takes one parameter shm_id,
   which is an integer that represents the ID of the shared memory segment that we want to attach to.
*/
void *sysv_shm_attach(int shm_id) {
    //
    return shmat(shm_id, NULL, 0);
}

/*
    sysv_shm_detach takes one parameter addr,
    which is a pointer to the memory address of the shared memory segment that we want to detach from.
*/
int sysv_shm_detach(void *addr) {
    return shmdt(addr);
}

/*
    sysv_shm_read reads data from a shared memory segment identified by shm_id and stores it in output.
    It takes three arguments:
        - shm_id: The integer that represents the ID of the shared memory segment that we want to read from
        - output: The pointer to the memory address where the data read from the shared memory segment
        - len:    The integer that represents the number of bytes to read from the shared memory segment
        - offset: The integer that represents the offset (in bytes) from the start of the shared memory segment where the read operation will begin
*/
int sysv_shm_read(int shm_id, void* output, int len, int offset) {
    // Attach to the given segment to get its memory address
    char* addr = sysv_shm_attach(shm_id);

    // Check if attachment was successful
    if(addr == (char*)(-1)){
        return -1;
    }

    // Copy len bytes from addr into output
    memcpy(output, addr+offset, len);

    // Detach from the shared memory segment
    sysv_shm_detach(addr);

    // Return 0 to indicate successful read
    return 0;
}

/*
    This code locks a System V shared memory segment with the given ID
    using the shmctl function with the SHM_LOCK command.
*/
int sysv_shm_lock(int shm_id) {
    return shmctl(shm_id, SHM_LOCK, NULL);
}

/*
    sysv_shm_unlock unlocks a System V shared memory segment with the given ID using the shmctl function with the SHM_UNLOCK command.
*/
int sysv_shm_unlock(int shm_id) {
    return shmctl(shm_id, SHM_UNLOCK, NULL);
}

// sysv_shm_close removes a System V shared memory segment with the given ID using the shmctl function with the IPC_RMID command.
int sysv_shm_close(int shm_id) {
    return shmctl(shm_id, IPC_RMID, NULL);
}

/*
    sysv_shm_get_size  retrieves the size of a System V shared memory segment identified by shm_id using the shmctl function with the IPC_STAT command.
    - shm_id: the integer argument shm_id, which represents the ID of a System V shared memory segment.
*/
size_t sysv_shm_get_size(int shm_id) {
    struct shmid_ds shm;

    /*
        In this case, the function returns the size of the shared memory segment in bytes,
        which is stored in the shm_segsz field of the shmid_ds structure.
    */
    if(shmctl(shm_id, IPC_STAT, &shm) >= 0) {
        return shm.shm_segsz;
    }else{
        return -1;
    }
}
