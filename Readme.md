Dartmouth Nutrition Scraper
===========================
## Introduction
This program is part of the Dartmouth Nutrition Scraper
project for CS98 at Dartmouth College.  

The purpose of this scraper is to grab all the nutritional
information from the [Dartmouth Nutrition Website](http://nutrition.dartmouth.edu:8088/)

## Terminology
1. venue : Dining Hall. eg: FoCo
2. sid : Dartmouth's API uses this to indentify each venue
3. recipe : Dartmouth's name for each of the available items at each venue

## Known Issues
The progress bars for some reason do not go up to 100%.
However the data all seems to get filled in so this may be a non issue. Still should look into it some more.

## Requirements
To build this program you must have Go installed. A binary is included but has only been tested on OSX.

## Installation
```
git clone https://github.com/jesusrmoreno/Dartmouth-Nutrition-Scraper scraper

cd scraper
go run main.go --help
go run main.go // will not create output files
go run main.go --write-files // will create output files can also use --wf

// Or

go build
./nutrition-scraper --help
./nutrition-scraper // will not create output files
./nutrition-scraper --write-files // will create output files can also use --wf
```

## Output
```
Starting scrape at:  2015-09-28 10:47:23.354492421 -0400 EDT
Starting NOVACK at 2015-09-28 10:47:23.379001018 -0400 EDT
Getting recipes for  NOVACK at 2015-09-28 10:47:23.405734455 -0400 EDT
37 / 37 [=====================================================] 100.00 % 14/s 2s
Done getting recipes for: NOVACK
Starting DDS at 2015-09-28 10:47:26.114652518 -0400 EDT
Getting recipes for  DDS at 2015-09-28 10:47:26.824192169 -0400 EDT
751 / 751 [==================================================] 100.00 % 12/s 58s
Done getting recipes for: DDS
Starting CYC at 2015-09-28 10:48:25.935099207 -0400 EDT
Getting recipes for  CYC at 2015-09-28 10:48:26.614847761 -0400 EDT
1144 / 1144 [==============================================] 100.00 % 12/s 1m28s
Done getting recipes for: CYC
Entire Scrape took 2m34.0107052s

```

## TODO
* Add database once we figure out the schema
* Maybe add an HTTP API and Frontend progress monitor
* Add ability to get a range of dates
* hash and skip items we've already scraped to speed up the scrape
