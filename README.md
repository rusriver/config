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

- If E() is set, it sets error if any;
- If Ok() is set, it sets true on success, false on error;
- If U() is set, it doesn't panic on error, if no E() or Ok() were set;
- If neither set, it panics on error (to be intercepted as exception);
- In either case, the default value is returned;

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
    conf := config.InitContext{}.FromFile("filename.yaml").E(&err).Load()  // detected by suffix

    conf2 := config.InitContext{}.FromBytes([]byte(`text here`)).E(&err).Load() // tries all known formats
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

