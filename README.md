# prov

`prov` is a simple provenance system for files produced by command line tools that write to stdout.

The main use case is figuring out when and how a particular file was created.

`prov` is actually two commands: `prov` and `whence`.

## Example Usage

```
% prov cat test.in > test.out
```

> hours, days, months later...

```
% cat test.out
He who has a shady past knows that nice guys finish last.
```

> hmm, that's odd, where did this file come from?

```
% whence test.out
Hash: e5dea09392dd886ca63531aaa00571dc07554bb6
Time: 2015-12-12 16:49:09.747598265 -0800 PST
User: banksean
Directory: /Users/banksean/src/prov
Command: cat test.in
```

> Oh, I remember that now.

## Details

The design is very simple, basically a process wrapper that logs some stuff:

### How prov works
- Take everything that comes after `prov` on the command line, and run it in a subprocess. 
- Read the subprocess's stdout stream into a buffer, add it to the SHA1, and also write it to our own process's stdout.
- Subprocess's stderr just goes to our stderr.
- Add a new line to `~/.prov` with the SHA1 of the file, plus some additional information about how and when the file was created.

### How whence works
- Buffer the file in question, calculate its SHA1.
- Check `~/.prov` for that SHA1 and print anything it finds.

Note that in the example, *even if `test.out` had been moved to a different location* after the fact, you could still run whence and get the same information.

## Ideas for future improvements

- Move the data currently kept in `~/.prov` to a central location shared by multiple users.  (Just realized this is very similar to @gojomo's now-defunct https://bitzi.com/ project :)
- Try to determine if any of the subprocess's input files *also* have provenance information, and store references to each input file's SHA1. Keep a family tree of the files.
- FUSE LOL

## Why

I work with lots of scratch data files (nothing you'd keep in source control), transforming them in various ways.  Sometimes I lose track of how I created a particular file or if it was really filtered from that other file the way I currently believe it was.  

I'm aware of more full-featured implementations of this idea but I just wanted something dirt simple I can use with few external dependencies.
