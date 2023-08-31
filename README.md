# Config [![GoDoc](https://godoc.org/github.com/rusriver/config?status.png)](https://godoc.org/github.com/rusriver/config)

Package config provides convenient access methods to configuration
stored as JSON or YAML.

This is a fork of [olebedev/config](https://github.com/olebedev/config),
which in turn is a fork of [original version](https://github.com/moraes/config).
It has incompatibilities with the original(s), and quite an extended functionality.
In fact, it has diverged a lot, and continues to.

Can be used not just as a config, but as a general data model storage, with
thread-safe high-performance mutability. It has limitations, but if you know
what this is about, you'll see them yourself. After all, it's open source.

TODO NEXT:

- Serk support
- provide extensive examples of use, suitable for copy-paste
- repair old original tests
- refactor it to the MAV model

## The v2 improvements:

U*() functions are removed. Instead, get-type functions behave this way:

- If Err() is set, it sets error if any;
- If Ok() is set, it sets true on success, false on error;
- If U() is set, it doesn't panic on error, if no Err() or Ok() were set;
- If neither set, it panics on error (to be intercepted as exception);
- In either case, the default value is returned;
- ErrOk() resets error to nil and ok to true, if either is set;
- ErrPtr and OkPtr are accessible directly as well;

Path now is specified always in P(), and type parsing happens as separate function,
and path is specified as a []string. E.g., instead of

```
    .UDuration("dot.path")
```

you now have to write:

```
    .U().P("dot", "path").Duration()
```

The API model was simplified, calling accessors with empty path is no longer needed,
the GetNestedConfig() was removed.

Different mapping rules for env variables - now all dashes are removed. For example,
if you have a path "a.s-d.f", previously the env variable A_S-D_F would be looked for,
now it will be A_SD_F. Obviously, both "sd" and "s-d" would map to the same thing,
but it's not a problem if know about it.

New idiom to load config, with automatic file type or data format detection:

```
    var err error
    conf := (&config.InitContext{}).FromFile("filename.yaml").Err(&err).Load()  // detected by suffix

    err = nil
    conf2 := (&config.InitContext{}).FromBytes([]byte(`text here`)).Err(&err).Load() // tries all known formats
```

Added LoadWithParenting().

## Thread-safety

There are three M.O. to use it:

1) Load once, then use concurrently, but keep it immutable. Most often used case.

2) Use Set() method, but then it's not thread-safe.

3) Thread-safe mutability support, using high-performance RCU, and Source object idiom.

Example of initializing the Source object:

```
        k.Source = config.NewSource(func (opts *config.NewSource_Options) {
            opts.Config = conf
            opts.Context = ctx
            // + other opts, optionally
        })
```

Example of using the config from the Source:

```
        conf := k.Source.Config
        // use conf as usual, but update it once in a while
```

In variant 2, the Set() method is just a by-pass to the NonThreadSafe_Set().

In 2 and 3, the usage of Set() by user is identical.

## Explicit synchronization of a completion of a batch of concurrent Set()

The Set() command is totally async (or how you'd expect to wait for each command completion?)
Instead, see next paragraph.
        
What is you do want to flush you commands immediatelly, and make sure they were indeed
executed, so you can safely get an updated Config from Source?

For this, the ChFlushSignal has ChDown, by which you can get notified if flush indeed
completed.

Because there can be several concurrent flush customers, we need to make ChFlushSignal
buffered, e.g. 10% of ChCmd.

So, if a user wants to make sure its commands took effect, it does this:

```
	// sync the config
	chDown := make(chan struct{})
	configSource.ChFlushSignal <- &config.MsgFlushSignal{ChDown: chDown}
	<-chDown
	
	// here we're certain the commands were flushed
```

Be careful: if you send explicit flush signal, with ChDown != nil, and never read from it,
you'll hang the whole write-back updater goroutine.

## Clarification on use of err, ok, and ErrOk()

Repeated use of err, and or misuse of ok, and forgetting to use ErrOk(). No method in this
library does explicitly sets ok=true, or err=nil, a user must do this itself. For example:

```
    // This code is totally wrong

    var err error
    var ok bool  // not set to true

    ... = conf.Ok(&ok).Err(&err).P("a", "s").MapConfig()    // expr-1
    if err != nil {}    // correct
    if !ok {}           // wrong, because ok is anyway false

    someFunc(conf)

    ... = conf.Ok(&ok).Err(&err).P("a", "s").MapConfig()  // expr-2
    ... = conf.P("a", "s").MapConfig()    // expr-3, identical to expr-2, because err and ok are left set in conf
    if err != nil {}    // wrong, will react to the error from expr-1 or expr-4
    if !ok {}           // wrong, because ok is anyway false, also will react to !ok from expr-1 or expr-4

    ... = conf.ErrOk().Ok(&ok).Err(&err).P("a", "s").MapConfig()  // expr-5, WRONG

func someFunc(conf *config.Config) {
    ... = conf.Duration() // expr-4
}
```

So how to use it right? Several rules:

1) Always write `ok := true`, instead of `var ok bool`;
2) Before new expression, if you re-use the err, write `err = nil`;
3) Can call `.ErrOk().` in the beginning of expression, to reset err=nil and ok=true explicitly;
4) In an epression where Err() or Ok() are set, the ErrOk() must be __after__ them;
5) In long expression, only __first__ failing method sets err and ok. This is because the only first one is relevant, all subsequent ones would fail anyway with not useful error values, therefore we are interested in only the first error in an expression.

Here's the same code, re-written correctly:

```
    var err error
    ok := true

    ... = conf.Ok(&ok).Err(&err).P("a", "s").MapConfig()    // expr-1
    if err != nil {}    // correct
    if !ok {}           // correct

    someFunc(conf)

    err = nil
    ok = true
    ... = conf.Ok(&ok).Err(&err).P("a", "s").MapConfig()  // expr-2
    if err != nil {}    // correct
    if !ok {}           // correct

    ... = conf.ErrOk().         // using ErrOk() is less LOC
        P("a", "s").MapConfig() // expr-3, identical to expr-2, because err and ok are left set in conf
    if err != nil {}    // correct
    if !ok {}           // correct

    ... = conf.Ok(&ok).Err(&err).ErrOk().P("a", "s").MapConfig()  // expr-5, correct

func someFunc(conf *config.Config) {
    ... = conf.Duration() // expr-4
}
```

## Default value callbacks

An optional callback function can be supplied to each value func, which will
be called in case the expression failed, and yield the default value. If this happens,
the function is called, and whatever it returns would be the result. The idea is that it's
not just a constant default value, but that instead there can be some logic, which makes a
dependency on other values (e.g. "default is false if ca-cert is set, true if ca-cert
isn't set"), those values being captured in a closure.

Example:

```
        caUseSystem = php.ErrOk().P("ca-use-system").Bool(func() {
            return !(len(caCert) > 0)   // caCert is captured in a closure
        })
```

Of course, you can achieve the same effect just setting Ok(), and then checking
it in an "if". Or you can check Err(). Yes you can. But this approach is somehow more
clean sometimes.

Clarification on how it detects that an expression failed:

1) If Err() or Ok() were set, we look at them;
2) Otherwise, we look at the .ExpressionFailure, if it is set:
   * =0    okay, nothing failed
   * =1    failed
   * =2    failed, one callback was already executed, new callbacks will cause panic
3) The .ExpressionFailure is being reset by ErrOk(), or explicitly;

This means, that if you don't have Err() or Ok() set, please use ErrOk()
in the beginning of your expressions.

To protect against missing ErrOk(), the logic guards against repeated calls
of default functions in subsequent expressions, without prior ErrOk(), when
there's no Err() or Ok() set.

