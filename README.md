Semantic Versioning for Golang
==============================

[![GoDoc](https://godoc.org/github.com/wmark/semver?status.png)](https://godoc.org/github.com/wmark/semver)

A library for parsing and processing of *Versions* and *Ranges* in:

* [Semantic Versioning](http://semver.org/) (semver) v2.0.0 notation
  * used by npmjs.org, pypi.org…
* Gentoo's ebuild format

Does not rely on *regular expressions* neither does it use package *reflection*.

```bash
$ sed -i -e 's@ignore@3rdparty@g' *_test.go
$ go test -tags 3rdparty -run=XXX -benchmem -bench=.

BenchmarkHashicorpNewVersion-24          1000000  1800 ns/op   551 B/op   5 allocs/op
BenchmarkBlangMake-24                    3000000   790 ns/op   320 B/op   5 allocs/op
BenchmarkSemverNewVersion-24            30000000    48.5 ns/op   0 B/op   0 allocs/op ←

BenchmarkHashicorpNewConstraint-24        200000  6350 ns/op  2096 B/op  18 allocs/op
BenchmarkBlangParseRange-24              1000000  1440 ns/op   480 B/op  13 allocs/op
BenchmarkSemverNewRange-24              10000000   120 ns/op     0 B/op   0 allocs/op ←

BenchmarkHashicorpCompare-24             1000000  1005 ns/op   395 B/op  12 allocs/op
BenchmarkBlangCompare-24               100000000    20.2 ns/op   0 B/op   0 allocs/op
BenchmarkSemverCompare-24              100000000    14.6 ns/op   0 B/op   0 allocs/op ←

```

Licensed under a [BSD-style license](LICENSE).

Usage
-----

Using _go modules_ you'd just:

```go
import "blitznote.com/src/semver/v3"
```

… or, with older versions of _Go_ leave out the version suffix `/v…` and:

```bash
$ dep ensure --add blitznote.com/src/semver@^3
```

After which you can use the module as usual, like this:

```go
v1 := semver.MustParse("1.2.3-beta")
v2 := semver.MustParse("2.0.0-alpha20140805.456-rc3+build1800")
v1.Less(v2) // true

r1, _ := NewRange("~1.2")
r1.Contains(v1)      // true
r1.IsSatisfiedBy(v1) // false (pre-releases don't satisfy)
```

Also check the [GoDocs](https://godoc.org/blitznote.com/src/semver)
and [Gentoo Linux Ebuild File Format](http://devmanual.gentoo.org/ebuild-writing/file-format/),
[Gentoo's notation of dependencies](http://devmanual.gentoo.org/general-concepts/dependencies/).

Please Note
-----------

It is, ordered from lowest to highest:

    alpha < beta < pre < rc < (no release type/»common«) < r (revision) < p

Therefore it is:

    Version("1.0.0-pre1") < Version("1.0.0") < Version("1.0.0-p1")

### Limitations

Version 2 no longer supports dot-tag notation.
That is, `1.8.rc2` will be rejected, valid are `1.8rc2` and `1.8-rc2`.

Usage Note
----------

Most *NodeJS* authors write **~1.2.3** where **>=1.2.3** would fit better.
*~1.2.3* is ```Range(">=1.2.3 <1.3.0")``` and excludes versions such as *1.4.0*,
which almost always work.

Contribute
----------

Pull requests are welcome.  
Please add your name and email address to a file *AUTHORS* and/or *CONTRIBUTORS*.  
Thanks!
