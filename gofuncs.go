package goduk

/*
#cgo CFLAGS: -std=c99 -O2 -Os -fomit-frame-pointer -fstrict-aliasing -DDUK_OPT_NO_ES6_OBJECT_SETPROTOTYPEOF -DDUK_OPT_NO_ES6_OBJECT_PROTO_PROPERTY -DDUK_OPT_NO_ES6_PROXY -DDUK_OPT_NO_AUGMENT_ERRORS -DDUK_OPT_NO_TRACEBACKS

#include <stdlib.h>
#include "duktape.h"

extern int funcs_mainForeignFunc(duk_context* ctx);

static duk_idx_t __duk_push_main_c_function(duk_context* ctx, duk_idx_t nargs) {
	return duk_push_c_function(ctx, funcs_mainForeignFunc, nargs);
}
*/
import "C"
import "unsafe"

/* ------ */
/* public */
/* ------ */

type GoFunc struct {
	Name  string
	F     func(*Context)int
	Nargs int
}

// NEVER call these function from inside a foreign function body
func (this Context) PutGoFunc(obj_index int, f GoFunc) {
	obj_index = this.NormalizeIndex(obj_index)

	// put the function in context's funcs map entry
	magic := funcs_ctxAddFunc(unsafe.Pointer(this.ctx), f.F)

	// push funcs_mainForeignFunc to stack
	C.__duk_push_main_c_function(unsafe.Pointer(this.ctx), C.duk_idx_t(f.Nargs))

	// set new function's magic
	C.duk_set_magic(unsafe.Pointer(this.ctx), C.duk_idx_t(-1), C.duk_int_t(magic))

	// create CString of function name
	cS_name := C.CString(f.Name)
	defer C.free(unsafe.Pointer(cS_name))

	// put function
	C.duk_put_prop_string(unsafe.Pointer(this.ctx), C.duk_idx_t(obj_index), cS_name)
}

func (this Context) PutGoFuncs(obj_index int, f []GoFunc) {
	for _, val := range f {
		this.PutGoFunc(obj_index, val)
	}
}

/* ------- */
/* private */
/* ------- */

// maps between a context and a foreign functions slice
var funcsMap map[unsafe.Pointer][]func(*Context)int


func funcs_ctxAddFunc(ctx unsafe.Pointer, f func(*Context)int) int {
	if funcsMap == nil {
		funcsMap = make(map[unsafe.Pointer][]func(*Context)int, 0)
	}

	_, ok := funcsMap[ctx]
	if !ok {
		funcsMap[ctx] = make([]func(*Context)int, 0)
	}

	funcsMap[ctx] = append(funcsMap[ctx], f)
	if len(funcsMap[ctx]) > 0 {
		return len(funcsMap[ctx]) - 1
	} else {
		return -1
	}
}

func funcs_ctxGetFunc(ctx unsafe.Pointer, magic_key int) func(*Context)int {
	slice, ok := funcsMap[ctx]
	if ok {
		if magic_key >= 0 && magic_key < len(slice) {
			return slice[magic_key]
		}
	}
	return nil
}

func funcs_ctxDel(ctx unsafe.Pointer) {
	delete(funcsMap, ctx)
}


//export funcs_mainForeignFunc
func funcs_mainForeignFunc(pctx unsafe.Pointer) C.int {
	ctx := new(Context)
	ctx.ctx = (*C.struct_duk_context)(pctx)

	// get magic value of the current function
	ctx.PushCurrentFunction()
	magic := int(C.duk_get_magic(unsafe.Pointer(pctx), C.duk_idx_t(-1)))
	ctx.Pop()

	// retreive go func from the funcsMap
	f := funcs_ctxGetFunc(pctx, magic)
	if f != nil {
		return C.int(f(ctx))
	}

	return RET_UNIMPLEMENTED_ERROR
}
