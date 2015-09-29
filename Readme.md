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
go run main.go --write-files

Will write files when done..

Working on: Novack Cafe
39 / 39 [=====================================================] 100.00 % 11/s 3s
Done getting nutrition info for: Novack Cafe

Working on: 53 Commons
775 / 775 [================================================] 100.00 % 10/s 1m10s
Done getting nutrition info for: 53 Commons

Working on: Courtyard Cafe
1148 / 1148 [==============================================] 100.00 % 11/s 1m40s
Done getting nutrition info for: Courtyard Cafe
Entire Scrape took 2m57.80872732s
```

## TODO
* Add database once we figure out the schema
* Maybe add an HTTP API and Frontend progress monitor
* Add ability to get a range of dates
* hash and skip items we've already scraped to speed up the scrape
