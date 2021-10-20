This is a small tool for validating glibc versions on ELF binaries.

Examples:

```
$ go build ./cmd/glibc-check

$ ./glibc-check list-versions ./glibc-check
2.2.5
2.3.2

$ ./glibc-check max ./glibc-check
2.3.2

$ ./glibc-check min ./glibc-check
2.2.5

$ ./glibc-check assert-all 'major == 2 && minor >= 2' ./glibc-check
$ echo $?
0

$ ./glibc-check assert-all 'major == 2 && minor != 2' ./glibc-check
condition did not hold for versions: 2.2.5
$ echo $?
1
```
