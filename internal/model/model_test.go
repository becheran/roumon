package model_test

import (
	"strings"
	"testing"

	"github.com/becheran/roumon/internal/model"
	"github.com/stretchr/testify/assert"
)

var trace_1 string = `goroutine 4431 [running]:
runtime/pprof.writeGoroutineStacks(0xe491c0, 0xc0001380e0, 0x0, 0x0)
	/usr/local/go/src/runtime/pprof/pprof.go:693 +0xb8
runtime/pprof.writeGoroutine(0xe491c0, 0xc0001380e0, 0x2, 0x0, 0x0)
	/usr/local/go/src/runtime/pprof/pprof.go:682 +0x4b
runtime/pprof.(*Profile).WriteTo(0x11306a0, 0xe491c0, 0xc0001380e0, 0x2, 0x0, 0x0)
	/usr/local/go/src/runtime/pprof/pprof.go:331 +0xa2
net/http/pprof.handler.ServeHTTP(0xc000fc8731, 0x9, 0xe531a0, 0xc0001380e0, 0xc000f17f00)
	/usr/local/go/src/net/http/pprof/pprof.go:256 +0x3e8
net/http/pprof.Index(0xe531a0, 0xc0001380e0, 0xc000f17f00)
	/usr/local/go/src/net/http/pprof/pprof.go:367 +0x130
net/http.HandlerFunc.ServeHTTP(0xe02378, 0xe531a0, 0xc0001380e0, 0xc000f17f00)
	/usr/local/go/src/net/http/server.go:2042 +0x44
net/http.(*ServeMux).ServeHTTP(0x117f100, 0xe531a0, 0xc0001380e0, 0xc000f17f00)
	/usr/local/go/src/net/http/server.go:2417 +0x1ab
net/http.serverHandler.ServeHTTP(0xc000138000, 0xe531a0, 0xc0001380e0, 0xc000f17f00)
	/usr/local/go/src/net/http/server.go:2843 +0x22b
net/http.(*conn).serve(0xc000fe5f40, 0xe54aa0, 0xc000fbab80)
	/usr/local/go/src/net/http/server.go:1925 +0x1805
created by net/http.(*Server).Serve
	/usr/local/go/src/net/http/server.go:2969 +0x970

goroutine 1 [chan receive, 16 minutes]:
company/foo/bar/SecureTest/internal/Testservice.(*TestService).Start(0xc0001adc70)
	/home/user/dev/TestService/code/testapp/internal/Testservice/Testservice.go:179 +0x3c5
main.main()
	/home/user/dev/TestService/code/testapp/cmd/TestService/main.go:109 +0xcf0

goroutine 3 [select]:
company/foo/bar/SecureTest/internal/mylib.(*filetestStore).createWatcher.func1(0xc0001b0320)
	/home/user/dev/TestService/code/testapp/internal/mylib/testStore.go:485 +0x1be
created by company/foo/bar/SecureTest/internal/mylib.(*filetestStore).createWatcher
	/home/user/dev/TestService/code/testapp/internal/mylib/testStore.go:411 +0x159

goroutine 35 [IO wait]:
internal/poll.runtime_pollWait(0x7fd3bc60de38, 0x72, 0x0)
	/usr/local/go/src/runtime/netpoll.go:220 +0x65
internal/poll.(*pollDesc).wait(0xc00011c618, 0x72, 0x0, 0x0, 0x0)
	/usr/local/go/src/internal/poll/fd_poll_runtime.go:87 +0x9b
internal/poll.(*pollDesc).waitRead(0xc00011c618, 0xffffffffffffff00, 0x0, 0x0)
	/usr/local/go/src/internal/poll/fd_poll_runtime.go:92 +0x45
internal/poll.(*FD).Accept(0xc00011c600, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0)
	/usr/local/go/src/internal/poll/fd_unix.go:394 +0x419
net.(*netFD).accept(0xc00011c600, 0x0, 0x0, 0x0)
	/usr/local/go/src/net/fd_unix.go:172 +0x85
net.(*TCPListener).accept(0xc0000c8660, 0x0, 0x0, 0x0)
	/usr/local/go/src/net/tcpsock_posix.go:139 +0x65
net.(*TCPListener).Accept(0xc0000c8660, 0x0, 0x0, 0x0, 0x0)
	/usr/local/go/src/net/tcpsock.go:261 +0x78
net/http.(*Server).Serve(0xc000138000, 0xe52ee0, 0xc0000c8660, 0x0, 0x0)
	/usr/local/go/src/net/http/server.go:2937 +0x42e
net/http.(*Server).ListenAndServe(0xc000138000, 0x0, 0x0)
	/usr/local/go/src/net/http/server.go:2866 +0x1c5
net/http.ListenAndServe(0x7fff900e1f0e, 0xe, 0x0, 0x0, 0x0, 0x0)
	/usr/local/go/src/net/http/server.go:3120 +0xdb
company/foo/bar/SecureTest/cmd/TestService/foo.Initfoo.func1()
	/home/user/dev/TestService/code/testapp/cmd/TestService/foo/foo_debug.go:22 +0x5d
created by company/foo/bar/SecureTest/cmd/TestService/foo.Initfoo
	/home/user/dev/TestService/code/testapp/cmd/TestService/foo/foo_debug.go:21 +0x72`

func TestParseTrace(t *testing.T) {
	routines, err := model.ParseStackFrame(strings.NewReader(trace_1))
	assert.Nil(t, err)
	assert.Len(t, routines, 4)
	// Routine 0
	r0 := routines[0]
	assert.Equal(t, int64(4431), r0.ID)
	assert.Equal(t, "running", r0.Status)
	assert.Equal(t, int64(0), r0.WaitSinceMin)
	//created by net/http.(*Server).Serve
	//	/usr/local/go/src/net/http/server.go:2969 +0x970
	assert.Equal(t, "/usr/local/go/src/net/http/server.go", r0.CratedBy.File)
	assert.Equal(t, int32(2969), r0.CratedBy.Line)
	assert.Equal(t, 0x970, *r0.CratedBy.Position)
	assert.Equal(t, "net/http.(*Server).Serve", r0.CratedBy.FuncName)
	assert.Equal(t, "runtime/pprof.writeGoroutineStacks(0xe491c0, 0xc0001380e0, 0x0, 0x0)", r0.StackTrace[0].FuncName)
	assert.False(t, r0.LockedToThread)

	// Routine 1
	r1 := routines[1]
	assert.Equal(t, int64(1), r1.ID)
	assert.Equal(t, int64(16), r1.WaitSinceMin)
	assert.Equal(t, "chan receive", r1.Status)
	assert.Nil(t, r1.CratedBy)
	assert.False(t, r1.LockedToThread)

	// Routine 1
	r2 := routines[2]
	assert.Equal(t, int64(3), r2.ID)
	assert.False(t, r2.LockedToThread)

	// Routine 3
	r3 := routines[3]
	assert.Equal(t, int64(35), r3.ID)
	assert.Equal(t, "IO wait", r3.Status)
	assert.Equal(t, "company/foo/bar/SecureTest/cmd/TestService/foo.Initfoo", r3.CratedBy.FuncName)
	assert.Equal(t, "/home/user/dev/TestService/code/testapp/cmd/TestService/foo/foo_debug.go", r3.CratedBy.File)
	assert.Equal(t, int32(21), r3.CratedBy.Line)
	assert.Equal(t, 0x72, *r3.CratedBy.Position)
	r3Stack := r3.StackTrace[len(r3.StackTrace)-1]
	assert.Equal(t, "/home/user/dev/TestService/code/testapp/cmd/TestService/foo/foo_debug.go", r3Stack.File)
	assert.Equal(t, "company/foo/bar/SecureTest/cmd/TestService/foo.Initfoo.func1()", r3Stack.FuncName)
	assert.Equal(t, 0x5d, *r3Stack.Position)
	assert.False(t, r3.LockedToThread)
}

func Benchmark_ParseTrace(b *testing.B) {
	for n := 0; n < b.N; n++ {
		model.ParseStackFrame(strings.NewReader(trace_1))
	}
}

var trace_2 = `goroutine 268 [runnable, locked to thread]:
syscall.Syscall9(0x7ff9af9b0500, 0x7, 0x1f4, 0xc0000902d8, 0x1, 0xc0000902c8, 0xc000090348, 0xc000090298, 0x0, 0x0, ...)
	C:/Program Files/Go/src/runtime/syscall_windows.go:356 +0xf2
syscall.WSARecv(0x1f4, 0xc0000902d8, 0x1, 0xc0000902c8, 0xc000090348, 0xc000090298, 0x0, 0x0, 0x0)
	C:/Program Files/Go/src/syscall/zsyscall_windows.go:1264 +0x12c
created by net/http.(*connReader).startBackgroundRead
	C:/Program Files/Go/src/net/http/server.go:688 +0xdb`

func TestParseLockedToThread(t *testing.T) {
	routines, err := model.ParseStackFrame(strings.NewReader(trace_2))
	assert.Nil(t, err)
	assert.Len(t, routines, 1)
	// Routine 0
	r0 := routines[0]
	assert.Equal(t, int64(268), r0.ID)
	assert.Equal(t, "runnable", r0.Status)
	assert.Equal(t, int64(0), r0.WaitSinceMin)
	assert.True(t, r0.LockedToThread)
}

func Benchmark_ParseHeader(b *testing.B) {
	for n := 0; n < b.N; n++ {
		model.ParseHeader("goroutine 268 [runnable, locked to thread]:")
	}
}

func Test_ParseHeader_Invalid(t *testing.T) {
	_, err := model.ParseHeader("")
	assert.NotNil(t, err)
	_, err = model.ParseHeader("gogogogogoroutines")
	assert.NotNil(t, err)
	_, err = model.ParseHeader("goroutine0fd")
	assert.NotNil(t, err)
}

func Test_ParseHeader_Valid(t *testing.T) {
	result, err := model.ParseHeader("goroutine 268 [runnable, locked to thread]:")
	assert.Nil(t, err)
	assert.Equal(t, int64(268), result.ID)
	assert.Equal(t, "runnable", result.Status)
	assert.Equal(t, true, result.LockedToThread)

	result, err = model.ParseHeader("goroutine 1 [chan receive, 16 minutes]:")
	assert.Nil(t, err)
	assert.Equal(t, int64(1), result.ID)
	assert.Equal(t, "chan receive", result.Status)
	assert.Equal(t, int64(16), result.WaitSinceMin)
	assert.Equal(t, false, result.LockedToThread)

	result, err = model.ParseHeader("goroutine 1 [chan receive, 16 minutes, locked to thread]:")
	assert.Nil(t, err)
	assert.Equal(t, int64(1), result.ID)
	assert.Equal(t, "chan receive", result.Status)
	assert.Equal(t, int64(16), result.WaitSinceMin)
	assert.Equal(t, true, result.LockedToThread)
}
