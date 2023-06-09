# Config [![GoDoc](https://godoc.org/github.com/rusriver/config?status.png)](https://godoc.org/github.com/rusriver/config)

Package config provides convenient access methods to configuration
stored as JSON or YAML.

This is a fork of [olebedev/config](https://github.com/olebedev/config),
which in turn is a fork of [original version](https://github.com/moraes/config).

It has incompatibilities with the original(s), and quite an extended functionality.

TODO:

- test ExtendBy_v2() (not tested yet, at all)
- Serk support
- provide extensive examples of use, suitable for copy-paste
- write tests, and make them work

## The v2 improvements:

U*() functions are removed. Instead, get-type functions behave this way:

- If E() is set, it sets error if any;
- If Ok() is set, it sets true on success, false on error;
- If U() is set, it doesn't panic on error, if no E() or Ok() were set;
- If neither set, it panics on error (to be intercepted as exception);
- In either case, the default value is returned;

Path now is specified always in P(), and type parsing happens as separate function. E.g.,
instead of

```
    .UDuration("dot.path")
```

you now have to write:

```
    .U().P("dot.path").Duration()
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


