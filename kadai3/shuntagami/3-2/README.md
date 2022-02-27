# How to use

## build

```
$ cd cmd
$ go build ./godl.go
```

## help

```
$ ./godl -h
Usage of ./godl:
  -buffer-size int
    	The buffer size to copy from http response body (default 32768)
  -f string
    	Output file name
  -n int
    	Concurrency level (default 1)
  -resume
    	Resume the download
  -u string
    	* Download url
```

## Download a file

```
./dl -u {URL} -n {CONCURRENCY_Level}
./dl -u https://sample-img.lb-product.com/wp-content/themes/hitchcock/images/100MB.png -n 10
```

## Interupt/Pause the download

Ctrl+c

## Resume the download

Use --resume

```
./dl -u https://sample-img.lb-product.com/wp-content/themes/hitchcock/images/100MB.png -n 10 --resume
```
