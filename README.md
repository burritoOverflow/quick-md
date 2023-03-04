Convert an input directory of markdown files to html, using [Blackfriday](https://github.com/russross/blackfriday/tree/v2).

Build:

```
$ go build
```

Run:

```
$ ./quickmd
```

Optional flags `in-dir` and `out-dir` for the input (markdown) directory and the ouput (generated html) directory. These default to `inputs` and `dist`, respectively.
