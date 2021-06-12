Semantic Versioning for Golang
==============================

[![GoDoc](https://godoc.org/blitznote.com/src/semver?status.png)](https://godoc.org/blitznote.com/src/semver)

A library for parsing and processing of *Versions* and *Ranges* in:

* [Semantic Versioning](http://semver.org/) (semver) v2.0.0 notation
  * used by npmjs.org, pypi.org…
* Gentoo's ebuild format
* The fastest implementation, and the one that'll actually parse all semver variants correctly and without errors.
* Sorting is in O(n).

Does not rely on *regular expressions* neither does it use package *reflection*.

```bash
$ sudo /bin/bash -c 'for g in /sys/bus/cpu/drivers/processor/cpu[0-9]*/cpufreq/scaling_governor; do echo performance >$g; done'
$ sed -i -e 's@ignore@3rdparty@g' foreign_test.go
$ go mod tidy
$ go test -tags 3rdparty -run=XXX -benchmem -bench=.

Benchmark_Compare-24                           1.392 ns/op 0 B/op   0 allocs/op
Benchmark_NewVersion-24                       32.02 ns/op  0 B/op   0 allocs/op
BenchmarkSemverNewRange-24                    86.83 ns/op  0 B/op   0 allocs/op
Benchmark_SortPtr-24                     2768859 ns/op

BenchmarkLibraryTwo_Compare-24                 5.019 ns/op 0 B/op   0 allocs/op
BenchmarkLibraryTwo_Make-24                  299.5 ns/op  75 B/op   2 allocs/op
BenchmarkLibraryTwo_ParseRange-24           1357 ns/op   456 B/op  13 allocs/op
BenchmarkLibraryTwo_Sort-24             12771299 ns/op

BenchmarkLibraryOne_Compare-24              1442 ns/op   480 B/op  17 allocs/op
BenchmarkLibraryOne_NewVersion-24           1516 ns/op   535 B/op   6 allocs/op
BenchmarkLibraryOne_NewConstraint-24        7024 ns/op  2092 B/op  18 allocs/op
BenchmarkLibraryOne_SortPtr-24         668274885 ns/op

# AMD Epyc 7401P, Linux 5.12.10, Go 1.16.5
# - LibraryOne v1.3.0 sometimes segfaults
# - LibraryTwo v4 errors on 19.4% of the given versions
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

Also check its [go.dev](https://pkg.go.dev/blitznote.com/src/semver/v3?tab=overview) listing
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

Contribute
----------

Pull requests are welcome.

For anything written in Assembly, please contribute your implementation for one
architecture only at first. We'll work with this and once it's in, follow up
with more if you like.

Please add your name and email address to a file *AUTHORS* and/or *CONTRIBUTORS*.  
Thanks!
