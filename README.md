# KDBX CSV to Pass

Takes your KDBX exported CSV file and calls the pass insert command - inserting with the following format %name%_%username%.

The app takes one argument that is the path to the kdbx csv file.

```sh
go run ./main.go %kdbx-csv-path%

```

## Note

This was written for fun to try out golang.
