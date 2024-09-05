//go:build darwin || nintendoswitch || wasip1

package syscall

import "unsafe"

// int getpagesize(void);
//
//export getpagesize
func libc_getpagesize() int

// char *getcwd(char *buf, size_t size)
//
//export getcwd
func libc_getcwd(buf *byte, size uint) *byte

// int stat(const char *path, struct stat * buf);
//
//export stat
func libc_stat(pathname *byte, ptr unsafe.Pointer) int32

// int fstat(int fd, struct stat * buf);
//
//export fstat
func libc_fstat(fd int32, ptr unsafe.Pointer) int32

// int lstat(const char *path, struct stat * buf);
//
//export lstat
func libc_lstat(pathname *byte, ptr unsafe.Pointer) int32

// int open(const char *pathname, int flags, mode_t mode);
//
//export open
func libc_open(pathname *byte, flags int32, mode uint32) int32

// DIR *fdopendir(int);
//
//export fdopendir
func libc_fdopendir(fd int32) unsafe.Pointer

// int fdclosedir(DIR *);
//
//export fdclosedir
func libc_fdclosedir(unsafe.Pointer) int32

// struct dirent *readdir(DIR *);
//
//export readdir
func libc_readdir(unsafe.Pointer) *Dirent

//export strlen
func libc_strlen(ptr unsafe.Pointer) uintptr

// ssize_t write(int fd, const void *buf, size_t count)
//
//export write
func libc_write(fd int32, buf *byte, count uint) int

// char *getenv(const char *name);
//
//export getenv
func libc_getenv(name *byte) *byte

// ssize_t read(int fd, void *buf, size_t count);
//
//export read
func libc_read(fd int32, buf *byte, count uint) int

// ssize_t pread(int fd, void *buf, size_t count, off_t offset);
//
//export pread
func libc_pread(fd int32, buf *byte, count uint, offset int64) int

// ssize_t pwrite(int fd, void *buf, size_t count, off_t offset);
//
//export pwrite
func libc_pwrite(fd int32, buf *byte, count uint, offset int64) int

// ssize_t lseek(int fd, off_t offset, int whence);
//
//export lseek
func libc_lseek(fd int32, offset int64, whence int) int64

// int close(int fd)
//
//export close
func libc_close(fd int32) int32

// int dup(int fd)
//
//export dup
func libc_dup(fd int32) int32

// void *mmap(void *addr, size_t length, int prot, int flags, int fd, off_t offset);
//
//export mmap
func libc_mmap(addr unsafe.Pointer, length uintptr, prot, flags, fd int32, offset uintptr) unsafe.Pointer

// int munmap(void *addr, size_t length);
//
//export munmap
func libc_munmap(addr unsafe.Pointer, length uintptr) int32

// int mprotect(void *addr, size_t len, int prot);
//
//export mprotect
func libc_mprotect(addr unsafe.Pointer, len uintptr, prot int32) int32

// int chdir(const char *pathname, mode_t mode);
//
//export chdir
func libc_chdir(pathname *byte) int32

// int chmod(const char *pathname, mode_t mode);
//
//export chmod
func libc_chmod(pathname *byte, mode uint32) int32

// int chown(const char *pathname, uid_t owner, gid_t group);
//
//export chown
func libc_chown(pathname *byte, owner, group int) int32

// int mkdir(const char *pathname, mode_t mode);
//
//export mkdir
func libc_mkdir(pathname *byte, mode uint32) int32

// int rmdir(const char *pathname);
//
//export rmdir
func libc_rmdir(pathname *byte) int32

// int rename(const char *from, *to);
//
//export rename
func libc_rename(from, to *byte) int32

// int symlink(const char *from, *to);
//
//export symlink
func libc_symlink(from, to *byte) int32

// int link(const char *oldname, *newname);
//
//export link
func libc_link(oldname, newname *byte) int32

// int fsync(int fd);
//
//export fsync
func libc_fsync(fd int32) int32

// ssize_t readlink(const char *path, void *buf, size_t count);
//
//export readlink
func libc_readlink(path *byte, buf *byte, count uint) int

// int unlink(const char *pathname);
//
//export unlink
func libc_unlink(pathname *byte) int32

// pid_t fork(void);
//
//export fork
func libc_fork() int32

// int execve(const char *filename, char *const argv[], char *const envp[]);
//
//export execve
func libc_execve(filename *byte, argv **byte, envp **byte) int

// int truncate(const char *path, off_t length);
//
//export truncate
func libc_truncate(path *byte, length int64) int32
