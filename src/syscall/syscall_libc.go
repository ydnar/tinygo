//go:build darwin || nintendoswitch || wasip1 || wasip2

package syscall

import (
	"unsafe"
)

type sliceHeader struct {
	buf *byte
	len uintptr
	cap uintptr
}

func Close(fd int) (err error) {
	if libc_close(int32(fd)) < 0 {
		err = getErrno()
	}
	return
}

func Dup(fd int) (fd2 int, err error) {
	fd2 = int(libc_dup(int32(fd)))
	if fd2 < 0 {
		err = getErrno()
	}
	return
}

func Write(fd int, p []byte) (n int, err error) {
	buf, count := splitSlice(p)
	n = libc_write(int32(fd), buf, uint(count))
	if n < 0 {
		err = getErrno()
	}
	return
}

func Read(fd int, p []byte) (n int, err error) {
	buf, count := splitSlice(p)
	n = libc_read(int32(fd), buf, uint(count))
	if n < 0 {
		err = getErrno()
	}
	return
}

func Pread(fd int, p []byte, offset int64) (n int, err error) {
	buf, count := splitSlice(p)
	n = libc_pread(int32(fd), buf, uint(count), offset)
	if n < 0 {
		err = getErrno()
	}
	return
}

func Pwrite(fd int, p []byte, offset int64) (n int, err error) {
	buf, count := splitSlice(p)
	n = libc_pwrite(int32(fd), buf, uint(count), offset)
	if n < 0 {
		err = getErrno()
	}
	return
}

func Seek(fd int, offset int64, whence int) (newoffset int64, err error) {
	newoffset = libc_lseek(int32(fd), offset, whence)
	if newoffset < 0 {
		err = getErrno()
	}
	return
}

func Open(path string, flag int, mode uint32) (fd int, err error) {
	data := cstring(path)
	fd = int(libc_open(&data[0], int32(flag), mode))
	if fd < 0 {
		err = getErrno()
	}
	return
}

func Fsync(fd int) (err error) {
	if libc_fsync(int32(fd)) < 0 {
		err = getErrno()
	}
	return
}

func Readlink(path string, p []byte) (n int, err error) {
	data := cstring(path)
	buf, count := splitSlice(p)
	n = libc_readlink(&data[0], buf, uint(count))
	if n < 0 {
		err = getErrno()
	}
	return
}

func Chdir(path string) (err error) {
	data := cstring(path)
	fail := int(libc_chdir(&data[0]))
	if fail < 0 {
		err = getErrno()
	}
	return
}

func Mkdir(path string, mode uint32) (err error) {
	data := cstring(path)
	fail := int(libc_mkdir(&data[0], mode))
	if fail < 0 {
		err = getErrno()
	}
	return
}

func Rmdir(path string) (err error) {
	data := cstring(path)
	fail := int(libc_rmdir(&data[0]))
	if fail < 0 {
		err = getErrno()
	}
	return
}

func Rename(from, to string) (err error) {
	fromdata := cstring(from)
	todata := cstring(to)
	fail := int(libc_rename(&fromdata[0], &todata[0]))
	if fail < 0 {
		err = getErrno()
	}
	return
}

func Link(oldname, newname string) (err error) {
	fromdata := cstring(oldname)
	todata := cstring(newname)
	fail := int(libc_link(&fromdata[0], &todata[0]))
	if fail < 0 {
		err = getErrno()
	}
	return
}

func Symlink(from, to string) (err error) {
	fromdata := cstring(from)
	todata := cstring(to)
	fail := int(libc_symlink(&fromdata[0], &todata[0]))
	if fail < 0 {
		err = getErrno()
	}
	return
}

func Unlink(path string) (err error) {
	data := cstring(path)
	fail := int(libc_unlink(&data[0]))
	if fail < 0 {
		err = getErrno()
	}
	return
}

func Chown(path string, uid, gid int) (err error) {
	data := cstring(path)
	fail := int(libc_chown(&data[0], uid, gid))
	if fail < 0 {
		err = getErrno()
	}
	return
}

func Fork() (err error) {
	fail := int(libc_fork())
	if fail < 0 {
		err = getErrno()
	}
	return
}

func Execve(pathname string, argv []string, envv []string) (err error) {
	argv0 := cstring(pathname)

	// transform argv and envv into the format expected by execve
	argv1 := make([]*byte, len(argv)+1)
	for i, arg := range argv {
		argv1[i] = &cstring(arg)[0]
	}
	argv1[len(argv)] = nil

	env1 := make([]*byte, len(envv)+1)
	for i, env := range envv {
		env1[i] = &cstring(env)[0]
	}
	env1[len(envv)] = nil

	fail := int(libc_execve(&argv0[0], &argv1[0], &env1[0]))
	if fail < 0 {
		err = getErrno()
	}
	return
}

func Truncate(path string, length int64) (err error) {
	data := cstring(path)
	fail := int(libc_truncate(&data[0], length))
	if fail < 0 {
		err = getErrno()
	}
	return
}

func Faccessat(dirfd int, path string, mode uint32, flags int) (err error)

func Kill(pid int, sig Signal) (err error) {
	return ENOSYS // TODO
}

type SysProcAttr struct{}

// TODO
type WaitStatus uint32

func (w WaitStatus) Exited() bool       { return false }
func (w WaitStatus) ExitStatus() int    { return 0 }
func (w WaitStatus) Signaled() bool     { return false }
func (w WaitStatus) Signal() Signal     { return 0 }
func (w WaitStatus) CoreDump() bool     { return false }
func (w WaitStatus) Stopped() bool      { return false }
func (w WaitStatus) Continued() bool    { return false }
func (w WaitStatus) StopSignal() Signal { return 0 }
func (w WaitStatus) TrapCause() int     { return 0 }

func Getenv(key string) (value string, found bool) {
	data := cstring(key)
	raw := libc_getenv(&data[0])
	if raw == nil {
		return "", false
	}

	ptr := uintptr(unsafe.Pointer(raw))
	for size := uintptr(0); ; size++ {
		v := *(*byte)(unsafe.Pointer(ptr))
		if v == 0 {
			src := *(*[]byte)(unsafe.Pointer(&sliceHeader{buf: raw, len: size, cap: size}))
			return string(src), true
		}
		ptr += unsafe.Sizeof(byte(0))
	}
}

func Setenv(key, val string) (err error) {
	if len(key) == 0 {
		return EINVAL
	}
	for i := 0; i < len(key); i++ {
		if key[i] == '=' || key[i] == 0 {
			return EINVAL
		}
	}
	for i := 0; i < len(val); i++ {
		if val[i] == 0 {
			return EINVAL
		}
	}
	runtimeSetenv(key, val)
	return
}

func Unsetenv(key string) (err error) {
	runtimeUnsetenv(key)
	return
}

func Clearenv() {
	for _, s := range Environ() {
		for j := 0; j < len(s); j++ {
			if s[j] == '=' {
				Unsetenv(s[0:j])
				break
			}
		}
	}
}

func Mmap(fd int, offset int64, length int, prot int, flags int) (data []byte, err error) {
	addr := libc_mmap(nil, uintptr(length), int32(prot), int32(flags), int32(fd), uintptr(offset))
	if addr == unsafe.Pointer(^uintptr(0)) {
		return nil, getErrno()
	}
	return (*[1 << 30]byte)(addr)[:length:length], nil
}

func Munmap(b []byte) (err error) {
	errCode := libc_munmap(unsafe.Pointer(&b[0]), uintptr(len(b)))
	if errCode != 0 {
		err = getErrno()
	}
	return err
}

func Mprotect(b []byte, prot int) (err error) {
	errCode := libc_mprotect(unsafe.Pointer(&b[0]), uintptr(len(b)), int32(prot))
	if errCode != 0 {
		err = getErrno()
	}
	return
}

// BytePtrFromString returns a pointer to a NUL-terminated array of
// bytes containing the text of s. If s contains a NUL byte at any
// location, it returns (nil, EINVAL).
func BytePtrFromString(s string) (*byte, error) {
	for i := 0; i < len(s); i++ {
		if s[i] == 0 {
			return nil, EINVAL
		}
	}
	return &cstring(s)[0], nil
}

// cstring converts a Go string to a C string.
func cstring(s string) []byte {
	data := make([]byte, len(s)+1)
	copy(data, s)
	// final byte should be zero from the initial allocation
	return data
}

func splitSlice(p []byte) (buf *byte, len uintptr) {
	slice := (*sliceHeader)(unsafe.Pointer(&p))
	return slice.buf, slice.len
}

// These two functions are provided by the runtime.
func runtimeSetenv(key, value string)
func runtimeUnsetenv(key string)
