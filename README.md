# Config [![GoDoc](https://godoc.org/github.com/rusriver/config/v2?status.png)](https://godoc.org/github.com/rusriver/config/v2)

Package config provides convenient access methods to configuration stored as JSON or YAML.

Originally this was a fork of [olebedev/config](https://github.com/olebedev/config), which in turn was
a fork of [moraes/config](https://github.com/moraes/config). Since then it grew and developed quite a lot,
got quite a lot of new functionality and power, it also has diverged a lot, and continues to. The same time
it is still the development of the same way of doing things as in the prior works.

Can be used not just as a config, but as a general data model storage, with thread-safe high-performance
mutability. It has limitations, but if you know what this is about, you'll see them yourself. After all,
it's open source, if in doubt - just look at the code.

---
20240900 Announcement: Further development will be in the /v3 folder, and will not be backward compatible.
This is production-grade version, but for new projects consider using /v3 instead.

---

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

```go
    .UDuration("dot.path")
```

you now have to write:

```go
    .U().P("dot", "path").Duration()
```

The API model was simplified, calling accessors with empty path is no longer needed,
the GetNestedConfig() was removed.

Different mapping rules for env variables - now all dashes are removed. For example,
if you have a path "a.s-d.f", previously the env variable A_S-D_F would be looked for,
now it will be A_SD_F. Obviously, both "sd" and "s-d" would map to the same thing,
but it's not a problem if know about it. Update: all punctuations except "_" are removed,
see the ExtendByEnvs_WithPrefix(), and also added the more powerful alternative the
ExtendByEnvsV2_WithPrefix().

New idiom to load config, with automatic file type or data format detection:

```go
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

```go
        k.Source = config.NewSource(func (opts *config.NewSource_Options) {
            opts.Config = conf
            opts.Context = ctx
            // + other opts, optionally
        })
```

Example of using the config from the Source:

```go
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

```go
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

```go
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

```go
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

```go
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

## Preserving default values, if the path is missing

Sometimes we have this situation:

```go
    someVariable = "default value"
    someVariable = conf.DotP("path.path").String()
```

Now we want to keep the someVariable with default value, if the value is missing in the config. Of course, we can use `if` with Ok(), but there's simpler idiom for that:

```go
    someVariable = "default value"
    someVariable = conf.DotP("path.path").String(func() string { return someVariable })
```

If the value is missing, the same default value will be returned and assigned back to it. This is convenient idiom.

## Prototype inheritance

We implement two kinds of prototype inheritance - file-level (parenting) and structural ($isa).

File-level ineritance is called "parenting", because you specify file parents for files.

The structural prototype inheritance, or the $isa inheritance, is similar to the parenting,
except the parenting applies to files, and the $isa applies to the structs inside the already
loaded data model, after the whole thing is already in memory.

Technically, both the parenting and the $isa are the instances of prototype inheritance,
just applied to different things at different stages and different levels.

### Parenting

In any configuration file you can specify, at the root level, the field "parent: string", or
"parents: [string, ...]", and specify parent config files, paths to them.

How does this work:

    1) Load the file;
    2) Look for fields "parent" or "parents", make a list of parent files;
    3) For each one - recursively repeat this algorithm, from step 1;
    4) Apply files on top of each other, in the order they were specified, and the current file apply on top of them all;

In other words: first go deep down the recursion, opening parent files; on returning up, apply each current file
on top of the parent one.

"Apply" means that when the fields are in conflict between the file we apply and the file we apply it to,
the fields we apply on top get the priority. I.e., the children has priority on parents.

There's an example of this mechanism in common use - the CSS in Internet browsers. They work by the same
principle: the styles defined at lower levels in hierarchy take the priority and redefine the styles
of higher levels of hierarchy. The difference is that here the hierarchy is specified explicitly, multiple
inheritance is possible, and the processing begins from the element (the file) located at the very bottom
of the hierarchy.

What it's useful for:

    1) It partically substitutes the functionality of "includes" - you can inherit several files, and they'll be "included" in the main one;

    2) It can be used to specify default values, and have multitude of specialization configs, which are applied on top of the default one;

This functionality is implemented in the `LoadWithParenting()`.

The resulting file get the field `parents-inherited` set at root level, which contains an array of paths
to all the files inherited. This is useful to see what file took part in it, quickly.

### The $isa

It means the "is-a". Implemented by the `TheIsa()`.

Suppose you have this structure:

```yaml
        objects:
            object-01:
                $isa: objects.object-02
                a: 2
                b: 3
            object-02:
                a: 0
                c: 4
```

After applying the `TheIsa()`, we'll get this YAML equivalent:

```yaml
        objects:
            object-01:
                a: 2
                b: 3
                c: 4
            object-02:
                a: 0
                c: 4
```

First we get the parent object, and apply the current one on top of that. That is, in case of intersections,
the current one takes the priority.

In this case, concerning the object `objects.object-01`:

1. we took the `objects.object-02`, and made a copy of it;
2. the "a" got a new value;
3. added the "b";
4. the "c" was left as it was, inherited from the prototype;

### The $isa Advanced

Handling the recursion - we do:

- recursive traversal of the tree;
- hash map to protect against loops;
- deep copy the isa struct, replace the current, remove the isa key from current, then apply current on the top;
- after this, re-parse the current structure again, to make sure there's no isa inherited left;

Multiple inheritance - you can specify several parent objects, and they all will be applied in order (just like
files), and the current one will be on top:

```yaml
    object-01:
        $isa: [objects.o1, defaults.x1]
```

Handling recursion - the re-parse happens when applying each consequtive layer.

Using relative paths:

```yaml
    objects:
        defaults:
            a: 1
            b: 2
        obj-01:
            $isa: ^.defaults
            name: obj-01
        obj-02:
            $isa: ^.defaults
            name: obj-02
        obj-03:
            $isa: ^.defaults
            name: obj-03
```

## Integrity checks

It supports integrity checks of source config files:

```go
	var err error
	conf := (&config.InitContext{}).
		FromFile("conf-test-files/config.yaml",
			"W7vJif0qw764-gXERIZ2HyQ0Rg0yvX5FnRc1USNAynI=",
			"ldtT3WFuCGSdXyGr5jUvibkNVLJzmlS_ajJzip95jEc=",
			"Kc_snY875G-uzwdVXGRpV-h8o7AodUgF_MfAugsx2PA=",
			"yEsexBCPF8HP86euYrrpDIrK7JHLrrVMhBUJFvYqWsE=",
			"Jo0cHrhVyEwtjIIApxQ0i_fr5UqsOOcE9Y6tWQlKGoM=",
		).
		Err(&err).
		LoadWithParenting()

	fmt.Println("CONFIG HASHES:", conf.InitContext.SourceHashesActual)
```

This main purpose is to tie test config files to tests, but can also be used for
other purposes as well.

## DDash language support

- set
- string-interpolate
- append (thread-unsafe (to be improved))

See tests for usage examples.
