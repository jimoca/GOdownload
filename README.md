## Clone

```
$ git clone https://github.com/jimoca/GOdownload.git
$ cd GOdownload
```
## Build

```
$ go build main.go
```

## or just run

```
$ go run main.go [-r] -path <path>
```

## CLI usage

```
Usage: main [-r] -path <path> 

Options:
-r      Use reader or not when reading files
-path string
        Path to the files folder in this project
```


## API

```
GET [host]:[port]/<filename_in_files>
```
