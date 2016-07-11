/* "Benchmarking" a table sort under Linux.
   # R.Wobst, @(#) Dec 27 2011, 20:02:29
   A table of ROWxCOL size is filled with random 31 bit integers and
   then sorted by row no. N where N runs from 0 to COL.
   The time is measured and also the memory usage, by grepping all lines
   from /dev/proc/status which start with "Vm".

   At the end, the same is done in Python to compare.

   The tests were executed on a AMD-X2 PC (3GHz, 64bit, dual core)

   Timings over few runs:
   0.8 sec for Go program
   1.3 sec for Python script

   Memory usage:
   40 MB peak for Go program
   62 MB peak for Python script
*/

package main

import (
        "fmt"
        "os"
        "rand"
        "sort"
        "strings"
        "time"
)

// the global var: sort column

var N int

const ROW int = 10000
const COL int = 100

// the list elements
type Row []int32
type Table []Row

func (lst Table) Len() int           { return len(lst) }
func (lst Table) Swap(i, j int)      { lst[j], lst[i] = lst[i], lst[j] }
func (lst Table) Less(i, j int) bool { return lst[i][N] < lst[j][N] }

// get memory usage
func prtmem() {
    procdev := fmt.Sprintf("/proc/%d/status", os.Getpid())
    fd, _ := os.Open(procdev)
    defer fd.Close()

    b := make([]byte, 4096)
    fd.Read(b)
    for _, lin := range strings.Split(string(b), "\n") {
        if "Vm" == lin[:2] {
            fmt.Printf("%s\n", lin)
        }
    }
}


func main() {
        tb := Table(make([]Row, ROW))

        // fill the table
        for i := 0; i < ROW; i++ {
                r := make(Row, COL)
                for j := 0; j < COL; j++ {
                        r[j] = rand.Int31()
                }
                tb[i] = r
        }

        t0 := time.Nanoseconds()

        // do the sort, one sort for each column
        for N = 0; N < COL; N++ {
                sort.Sort(tb)
        }

        t1 := time.Nanoseconds()
        fmt.Printf("%.2f seconds\n", float32(t1-t0)/1.e+9)

        // print memory usage
        prtmem()
}

/* The Python code which does the same.

#!/usr/bin/env python
# -*- coding: utf-8 -*-

import random, time, os

COL = 100
ROW = 10000
UPPER = (1<<31) - 1

Table = []

# fill table

for row in xrange(ROW):
    Row = []
    for col in xrange(COL):
        Row.append(random.randint(0, UPPER))
    Table.append(Row)

#sort table

t0 = time.time()

for scol in xrange(COL):
    Table.sort(key = lambda x: x[scol])

t1 = time.time()

print "%.2f seconds" % (t1-t0)

# get memory usage
fd = open('/proc/%d/status' % os.getpid(), 'r')
for l in fd.readlines():
    if l.strip().startswith('Vm'):
        print l,
fd.close()
*/