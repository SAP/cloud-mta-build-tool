#!/bin/bash
#bumping version automatically
echo 0.0.9 | 
awk '
BEGIN {
    FS="."                  # set the separators
    OFS=""
}
{
    sub(/^v/,"")            # remove the v
    $1=$1                   # rebuild record to remove periods
    $0=sprintf("%03d",$0+1) # zeropad the number after adding 1 to it
}
END {
    FS=""                   # reset the separators 
    OFS="."
    $0=$0                   # this to keep the value in $0
    $1=$1                   # rebuild to get the periods back
    print "v" $0            # output
}'

