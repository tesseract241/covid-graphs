
package main

import (
    "log"
    "path/filepath"
    "fmt"
    "io"
    "os"
    "sort"
    "encoding/csv"
    "strconv"
    "image"
    "image/color"
    "image/jpeg"
)

const (
    scale = 50
    horizontalOffset    = 100
    verticalOffset      = 50
    pointSize           = 5
    lineSize            = 2
    wantedProvinceCode  = "0" //Your province code here
)
func main() {
    var length, maxInfected int
    files, err := filepath.Glob("../data/dati-province/*[0-9].csv") 
    if err != nil {
        log.Fatal(err)
    }
    sort.Strings(files)
    length = len(files)
    maxInfected = 0
    points := make([]image.Point, 0, length) 
    for i, name := range files {
        f, err := os.Open(name)
        if err != nil {
            log.Fatal(err)
        }
        csvReader := csv.NewReader(f)
        for {
            record, err := csvReader.Read()
            if err == io.EOF {
                break
            }
            if err != nil {
                log.Fatal(err)
            }
            if record[6]== wantedProvinceName {
                dummy, e := strconv.Atoi(record[9])
                if e != nil {
                    log.Printf("While reading file %v failed to parse the number of infected (%v) with error %v\n", name, record[9], e)
                }
                if dummy>maxInfected {
                    maxInfected = dummy
                }
                points = append(points, image.Pt(i, dummy))
                fmt.Printf("From file %v got value %v\n", name, dummy)
            }
        }
    f.Close()
    }
    r := image.Rect(0,0, length * scale + horizontalOffset*2, maxInfected + verticalOffset*2)
    outPic := image.NewRGBA(r)
    pointColour := color.RGBA{255,0,0,255}
    lineColour := color.RGBA{128,0,0,255}
    axesColour := color.RGBA{191,191,191,255}
    //TODO Give the image some horizontal and vertical padding, draw the axes and mark their divisions
    for i:=verticalOffset;i<maxInfected + verticalOffset;i++ {
        for j:=-lineSize;j<=lineSize;j++ {
            outPic.Set(horizontalOffset + j, i, axesColour)
        }
    }
    for i:=horizontalOffset;i<length*scale + horizontalOffset;i++ {
        for j:=-lineSize;j<=lineSize;j++ {
            outPic.Set(i, j + maxInfected + verticalOffset, axesColour)
        }
    }
    for i, p := range(points) {
        for j:=-pointSize;j<=pointSize;j++ {
            for k:=-pointSize;k<=pointSize;k++ {
                outPic.Set(p.X*scale + j + horizontalOffset, maxInfected + verticalOffset - p.Y - k , pointColour)
            }
        }
        if i<len(points)-1 {
            var m float32 = float32(points[i+1].Y - points[i].Y)/float32(scale)
            for j:=0;j<scale;j++ {
                outPic.Set(p.X*scale+j + horizontalOffset,  maxInfected + verticalOffset - p.Y - int(m*float32(j)), lineColour)
            }
        }
    }
    f, err := os.Create("output.jpg")
    if err != nil {
        log.Fatal(err)
    }
    jpeg.Encode(f, outPic, nil)
}
