#!/bin/sh
#
# This is a very basic script that will pull down a list of content type to extensions and then output it into a Go map.  
#
TYPE_FILE="types.txt"
GO_FILE="../types.go"
wget -qO- http://svn.apache.org/repos/asf/httpd/httpd/trunk/docs/conf/mime.types | egrep -v ^# | awk '{ for (i=2; i<=NF; i++) {print $i" "$1}}' | sort > $TYPE_FILE

rm ../types.go

echo  "package main\n" >> $GO_FILE
echo "var contentTypes = map[string]string{" >> $GO_FILE 

while read p; do
  array=($p)
  echo "  " \"${array[0]}\" : \"${array[1]}\", >> $GO_FILE
done <$TYPE_FILE
echo "}" >> $GO_FILE