package goduk

import (
	"time"
	"testing"
	"fmt"
)

func duk_myError(ctx *Context) int {
	ctx.Error(0, "filename", 10, ctx.SafeToString(0))
	return 0
}

func duk_Sleep(ctx *Context) int {
	time.Sleep(time.Duration(ctx.ToInt(0)) * time.Millisecond)
	return 0
}

func execute(script string) {
	fmt.Println("Start:", script)

	ctx := CreateHeapDefault()
	defer ctx.Destroy()

	ctx.PushGlobalObject()
	ctx.PutGoFunc(-1, GoFunc{"Sleep", duk_Sleep, 1})
	ctx.PutGoFunc(-1, GoFunc{"myError", duk_myError, 1})
	ctx.Pop()

	if ctx.PcompileString(0, script, "filename") != 0 {
		fmt.Println("Pcompile error:", ctx.SafeToString(-1))
	} else {
		if ctx.Pcall(0) != 0 {
			fmt.Println("Pcall error:", ctx.SafeToString(-1))
		} else {
			fmt.Println("End")
		}
	}
}

func TestErrors(t *testing.T) {
	execute(`print("Will call error..."); Sleep(1000); throw "My error texte"; myError("myError() message")`)
}
