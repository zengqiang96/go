// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build amd64 386

package runtime

import (
	"runtime/internal/sys"
	"unsafe"
)

// adjust Gobuf as if it executed a call to fn with context ctxt
// and then did an immediate gosave.
func gostartcall(buf *gobuf, fn, ctxt unsafe.Pointer) {
	sp := buf.sp //newg的栈顶，目前newg栈上只有fn函数的参数，sp指向的是fn的第一参数
	if sys.RegSize > sys.PtrSize {
		sp -= sys.PtrSize
		*(*uintptr)(unsafe.Pointer(sp)) = 0
	}
	sp -= sys.PtrSize                        //为返回地址预留空间，
	*(*uintptr)(unsafe.Pointer(sp)) = buf.pc //在栈上放入goexit+1的地址
	buf.sp = sp                              //重新设置newg的栈顶寄存器
	//这里才真正让newg的ip寄存器指向fn函数，注意，这里只是在设置newg的一些信息，newg还未执行，
	//等到newg被调度起来运行时，调度器会把buf.pc放入cpu的IP寄存器，
	//从而使newg得以在cpu上真正的运行起来
	buf.pc = uintptr(fn)
	buf.ctxt = ctxt
}
