Go binding for duktape scripting language. Nearly everything is supported (including foreign functions).

Limitations:
- No lightfuncs
- Unable to use magic number for foreign functions (they are used to resolve the called function)
- Duktape's error handling is incompatible with Go (setjmp, longjmp). You can not throw errors from Go code (Context.Error, Context.RequireXXX) so those must not be used

Needs more testing, but the basics seems to work well!
