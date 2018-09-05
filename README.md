# Project for Interview

## Running the project

Build the docker image for the application and its code:

```
docker build -t interview_app .
```

Then it is ready to run:

```
docker-compose up
```

This will read data specified and will write the output to:

```
data/output/summary.txt
```

### Benchmarks

Timing is run for each file. The size of the file is not reported with
the timing to ingest. Results for 3 parsers (configurable) for a file
of size 2.3M is 15 seconds. For a small file such as 2.5K it takes
~26ms.

Timing is done for the total process. For 1003 files size 1-5K (~13MB
total) with 3 file readers that have 3 parsers each takes ~10 seconds.

The total process for the 1004 files, total size ~15MB takes ~25.5
seconds to ingest data and ~83ms to report the information.
