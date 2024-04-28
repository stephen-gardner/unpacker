# Unpacker
This package is a Go port of the Dean Edward's unpacker from [js-beautify](https://github.com/beautifier/js-beautify/blob/main/python/jsbeautifier/unpackers/packer.py). More unpackers may be added to this package when they are needed in another scraping project.
## Installation
```sh
$ go get -u github.com/stephen-gardner/unpacker
```
## Example
```go
source := "eval(function(p,a,c,k,e,r){e=String;if(!''.replace(/^/,String)){while(c--)r[c]=k[c]||c;k=[function(e){return r[e]}];e=function(){return'\\\\w+'};c=1};while(c--)if(k[c])p=p.replace(new RegExp('\\\\b'+e(c)+'\\\\b','g'),k[c]);return p}('0 2=1',62,3,'var||a'.split('|'),0,{}))"
deu, valid := unpacker.NewDEUnpacker(source)
if valid {
    res, err := deu.Unpack()
    if err != nil {
        log.Fatal("While unpacking js:", err)
    }
    // Output: var a=1
    log.Println(res)
}
```