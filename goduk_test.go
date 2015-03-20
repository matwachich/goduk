package goduk

import "testing"

func TestSimple(t *testing.T) {
	ctx := CreateHeapDefault()
	defer ctx.Destroy()

	ret := ctx.PevalString(
`
function fib(n) {
    if (n == 0) { return 0; }
    if (n == 1) { return 1; }
    return fib(n-1) + fib(n-2);
}

function test() {
    var res = [];
    for (i = 0; i < 20; i++) {
        res.push(fib(i));
    }
    print(res.join(' '));
	return "السلام";
}

test();
` , "test.js")

	if ret == 0 {
		println("OK - return =", ctx.GetString(-1))
	} else {
		println("Failed - err =", ctx.GetString(-1))
	}
}
