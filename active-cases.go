package main

import (
    "encoding/csv"
    "fmt"
    "image"
	"image/color"
    "image/jpeg"
    "io"
    "log"
    "math"
    "os"
    "path/filepath"
    "sort"
    "strconv"
    "sync"
    "github.com/tesseract241/paletteGenerator"
)

const (
    imageHeight             = 3000
    numberOfRegions         = 21
    scale                   = 50
    horizontalOffset        = 100
    verticalOffset          = 50
    availableImageHeight    = imageHeight - 2*verticalOffset
    lineSizeDefault         = 3
    pointSizeDefault        = 7
)

func renderAxes(im *image.RGBA, maxInfected int, length int, color color.Color) {
    verticalScale := float64(availableImageHeight)/float64(maxInfected)
    lineSize    := int(math.Max(lineSizeDefault, float64(lineSizeDefault)/(4.*verticalScale)))
    for i:=verticalOffset;i<imageHeight - verticalOffset;i++ {
        for j:=-lineSize;j<=lineSize;j++ {
            im.Set(horizontalOffset + j, i, color)
        }
    }
    for i:=horizontalOffset;i<length*scale + horizontalOffset;i++ {
        for j:=-lineSize;j<=lineSize;j++ {
            im.Set(i, availableImageHeight + verticalOffset + j, color)
        }
    }

}

func renderGraph(im *image.RGBA, graph []int, maxInfected int, color color.Color) {
    verticalScale := float64(availableImageHeight)/float64(maxInfected)
    pointSize   := int(math.Max(pointSizeDefault, float64(pointSizeDefault)/(7. * verticalScale)))
    for i := range graph {
        for k:=-pointSize;k<pointSize;k++ {
            for l:=-pointSize;l<pointSize;l++ {
                im.Set(i*scale + k + horizontalOffset, availableImageHeight + verticalOffset - int(float64(graph[i])*verticalScale) - l, color)
            }
        }
        if i<len(graph)-1 {
            var m float64 = (float64(graph[i+1] - graph[i])*verticalScale)/float64(scale)
            for k:=0; k<scale; k++ {
        im.Set(i*scale+k + horizontalOffset, imageHeight - verticalOffset - int(float64(graph[i])*verticalScale) - int(m*float64(k)), color)
            }
        }
    }
}


func main() {
    var length, maxInfected int
    var maxInfectedByRegion [numberOfRegions]int
    files, err := filepath.Glob("../data/dati-regioni/*[0-9].csv")
    if err != nil {
	    log.Fatal(err)
    }
    sort.Strings(files)
    length = len(files)
    maxInfected = 0
    graphs := make([][]int, numberOfRegions)
    regionsNames := make([]string, 0)
    g, err := os.Open(files[0])
    if err != nil {
        log.Fatal(err)
    }
    csvReader := csv.NewReader(g)
    for j:=0; ;j++{
        record, err := csvReader.Read()
	    if err == io.EOF {
			    break
		}
        if err != nil {
			log.Fatal(err)
        }
	    if j>0 {
            regionsNames = append(regionsNames, record[3])
        }
    }

    for i:=0 ;i<numberOfRegions; i++ {
        graphs[i] = make([]int, length)
    }

    for i, name := range files {
        f, err := os.Open(name)
        if err != nil {
	        log.Fatal(err)
        }
        csvReader := csv.NewReader(f)
        for j:=0; ;j++{
            record, err := csvReader.Read()
		    if err == io.EOF {
			    break
		    }
		    if err != nil {
			    log.Fatal(err)
            }
            if j>0 {
                dummy, e := strconv.Atoi(record[10])
                if e != nil {
                    log.Printf("While reading file %v failed to parse the number of infected (%v) with error %v\n", name, record[9], e)
                }
                if dummy>maxInfected {
                    maxInfected = dummy
                }
                if dummy>maxInfectedByRegion[j-1] {
                    maxInfectedByRegion[j-1] = dummy
                }
                graphs[j-1][i] = dummy
                //fmt.Printf("From file %v got value %v for region %s\n", name, dummy, record[3])
            }
        }
    f.Close()
    }
//    verticalScale := float32(availableImageHeight)/float32(maxInfected)
//    lineSize    := int(float32(lineSizeDefault)/(2.*verticalScale))
//    pointSize   := int(float32(pointSizeDefault)/(3.*verticalScale))
//    colors, err := paletteGenerator.GeneratePalette(numberOfRegions)
//    if err != nil {
//        //TODO Check if we need to release any resource
//        log.Fatal(err)
//    }
//    r := image.Rect(0,0, length * scale + horizontalOffset*2, imageHeight)
//    outPic := image.NewRGBA(r)
//    green := color.RGBA{0,255,0,255}
//    for i:=verticalOffset;i<imageHeight - verticalOffset;i++ {
//        for j:=-lineSize;j<=lineSize;j++ {
//            outPic.Set(horizontalOffset + j, i, green)
//        }
//    }
//    for i:=horizontalOffset;i<length*scale + horizontalOffset;i++ {
//        for j:=-lineSize;j<=lineSize;j++ {
//            outPic.Set(i, availableImageHeight + verticalOffset + j, green)
//        }
//    }
//    for j:=0;j<numberOfRegions;j++ {
//        for i := range(graphs) {
//            for k:=-pointSize;k<pointSize;k++ {
//                for l:=-pointSize;l<pointSize;l++ {
//                    outPic.Set(i*scale + k + horizontalOffset, availableImageHeight + verticalOffset - int(float32(graphs[i][j])*verticalScale) - l, colors[j])
//                }
//            }
//            if i<len(graphs)-1 {
//                var m float32 = (float32(graphs[i+1][j] - graphs[i][j])*verticalScale)/float32(scale)
//                for k:=0;k<scale;k++ {
//                    outPic.Set(i*scale+k + horizontalOffset, imageHeight - verticalOffset - int(float32(graphs[i][j])*verticalScale) - int(m*float32(k)), colors[j])
//                }
//            }
//        }
//    }
//    f, err = os.Create("output.jpg")
//    if err != nil {
//        log.Fatal(err)
//    }
//    jpeg.Encode(f, outPic, nil)
//
//
//
//
//
//
    r := image.Rect(0, 0, length * scale + horizontalOffset*2, imageHeight)
    outPic := image.NewRGBA(r)
    green := color.RGBA{0, 255, 0, 255}
    colors, err := paletteGenerator.GeneratePalette(numberOfRegions)
    var mux sync.Mutex
    c := make(chan int, len(graphs))
    for i, graph := range graphs {
        go func(i int, graph []int, c chan int) {
            outPic := image.NewRGBA(r)
            renderAxes(outPic, maxInfectedByRegion[i], length, green)
            renderGraph(outPic, graph, maxInfectedByRegion[i], colors[i])
            f, err := os.Create("covid-output/active-cases-" + regionsNames[i] + ".jpg")
            if err != nil {
                log.Fatal(err)
            }
            jpeg.Encode(f, outPic, nil)
            mux.Lock()
            fmt.Printf("Generated the image for region %s\n", regionsNames[i])
            mux.Unlock()
            c <- 1
        }(i, graph, c)
    }

    outPic = image.NewRGBA(r)
    renderAxes(outPic, maxInfected, length, green)
    for i, graph := range graphs {
        renderGraph(outPic, graph, maxInfected, colors[i])
        f, err := os.Create("covid-output/active-cases.jpg")
        if err != nil {
            log.Fatal(err)
        }
        jpeg.Encode(f, outPic, nil)
        mux.Lock()
        fmt.Printf("Added region %s to the cumulative graph image\n", regionsNames[i])
        mux.Unlock()
    }
    for _ = range graphs {
        <- c
    }
}
