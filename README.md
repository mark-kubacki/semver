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

BenchmarkHashicorpNewVersion-24          1000000  1805 ns/op   544 B/op   5 allocs/op
BenchmarkBlangMake-24                    3000000   516 ns/op    96 B/op   3 allocs/op
BenchmarkSemverNewVersion-24            20000000    60.6 ns/op   0 B/op   0 allocs/op ←

BenchmarkHashicorpNewConstraint-24        200000  7033 ns/op  2096 B/op  18 allocs/op
BenchmarkBlangParseRange-24              1000000  1600 ns/op   480 B/op  13 allocs/op
BenchmarkSemverNewRange-24              10000000   158 ns/op     0 B/op   0 allocs/op ←

BenchmarkHashicorpCompare-24             1000000  1075 ns/op   384 B/op  10 allocs/op
BenchmarkBlangCompare-24               100000000    18.1 ns/op   0 B/op   0 allocs/op
BenchmarkSemverCompare-24              100000000    15.0 ns/op   0 B/op   0 allocs/op ←

```

Licensed under a [BSD-style license](LICENSE).

Usage
-----
```bash
$ go get -v -d github.com/wmark/semver

or, better yet,

$ dep ensure --add github.com/wmark/semver@^2
```

```go
import "github.com/wmark/semver"

v1, err := semver.NewVersion([]byte("1.2.3-beta"))
v2, err := semver.NewVersion([]byte("2.0.0-alpha20140805.456-rc3+build1800"))
v1.Less(v2) // true

r1, err := NewRange("~1.2")
r1.Contains(v1)      // true
r1.IsSatisfiedBy(v1) // false (rejects pre-releases: alphas, betas…)
```

Also check the [GoDocs](http://godoc.org/github.com/wmark/semver)
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
