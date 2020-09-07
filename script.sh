#!/bin/bash
ROOT_FOLDER="$HOME/go/github.com/tesseract241/"
OUTPUT_FOLDER="$HOME/Nextcloud/Photos/"
cd "$ROOT_FOLDER" || (echo "Folder $ROOT_FOLDER does not exist" && exit)
(cd data && git pull)

if cd covid-active-cases ; then
    ./app
    if cp -r covid-output "$OUTPUT_FOLDER"; then
        echo "Correctly copied covid-active-cases output to $OUTPUT_FOLDER"
    else
        echo "Could not copy covid-active-cases output to $OUTPUT_FOLDER"
    fi
else
    echo "Folder $OUTPUT_FOLDER does not exist"
fi

if cd ../covid-total-cases-NA; then
    ./app
    if cp output.jpg "${OUTPUT_FOLDER}"covid-output/covid-total-cases-NA.jpg; then
        echo "Correctly copied covid-total-cases-NA output to ${OUTPUT_FOLDER}covid-output"
    else
        echo "Could not copy covid-total-cases-NA output to ${OUTPUT_FOLDER}covid-output"
    fi
else
    echo "Folder $ROOT_FOLDER covid-total-cases-NA does not exist"
fi
