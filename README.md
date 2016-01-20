# UN locode API

in-memory translation API for the UN locode, written in golang

## datasource

taken from http://www.unece.org/fileadmin/DAM/cefact/locode/loc152csv.zip

```
$ curl -O http://www.unece.org/fileadmin/DAM/cefact/locode/loc152csv.zip
$ unzip loc152csv.zip
$ cat *UNLOCODE*Part*.csv > combined.csv
```
