# Bcov

[![Docker](https://github.com/tikz/bcov/actions/workflows/docker-publish.yml/badge.svg)](https://github.com/tikz/bcov/actions/workflows/docker-publish.yml)
[![Go](https://github.com/tikz/bcov/actions/workflows/go.yml/badge.svg)](https://github.com/tikz/bcov/actions/workflows/go.yml)
[![Node.js CI](https://github.com/tikz/bcov/actions/workflows/node.js.yml/badge.svg)](https://github.com/tikz/bcov/actions/workflows/node.js.yml)

Effortlessly query coverage and read depths of genes and exonic variants across dozens of exome probe kits and hundreds of gigabytes of original BAM raw data.

---

Backend database preloading and serving of the HTTP API and React frontend is done through the same single binary with zero dependencies.

https://pkg.go.dev/github.com/tikz/bcov/cov

https://pkg.go.dev/github.com/tikz/bcov/api

https://pkg.go.dev/github.com/tikz/bcov/db

## Local build

### Requirements

[Go >=1.18](https://go.dev/doc/install) and [npm](https://nodejs.org/en/download/)

```
sudo apt install nodejs golang
```

### Build

```
git clone https://github.com/tikz/bcov.git
cd bcov

make
```

---

```
./bcov

Bcov
Version:  development                Commit:  unknown

Usage:
  -bam string
        Load a BAM file into the database
  -bams string
        Load multiple BAM files passing a CSV file that in each line contains <bam path>,<capture kit name>
  -delete-bam string
        Delete BAM file from the database with SHA256
  -fetch-exons
        Fetch exons from Ensembl and store in DB
  -fetch-variants
        Fetch SNPs from ClinVar and store in DB
  -help
        Display help
  -kit string
        Capture kit name
  -region string
        Print per position depth of a given range, expressed as <chromosome>:<start>-<end>
  -test-db
        Test database connection
  -web
        Run web server
```

## Database setup

### 1. Load all GRCh38 exons from Ensembl

`./bcov -fetch-exons`

Populates the local database fetching from [Ensembl core schema](https://www.ensembl.org/info/docs/api/core/core_schema.html) all exons in genes that have a valid gene symbol name defined by HGNC and a MANE Select canonical transcript.

![exons](https://i.imgur.com/Jpr8k6s.gif)

### 2. Load all dbSNP variants in exonic regions

`./bcov -fetch-variants`

Populates the local database using the latest `variant_summary.txt` from the [NCBI public FTP server](https://ftp.ncbi.nlm.nih.gov/pub/clinvar/tab_delimited/variant_summary.txt.gz).

![variants](https://i.imgur.com/Kn34LXE.gif)

### 3. Load a BAM file

`./bcov -bam <path> -kit <kitName>`

Loads a BAM file and assigns it under the given kit name.

Each BAM file is meant to be loaded only once.
A file is deemed unique (by practical purposes) and enforced using the SHA256 sum.

**Make sure that the BAM files you intend to load are aligned to GRCh38, otherwise _coordinate mappings will be wrong_**.

Expected performance of traversing a single BAM file is about **1 GB/min** using the SQLite engine, a consumer level NVMe SSD, and Ensembl/ClinVar datasets as of 2022. Resulting size in database is about ~1/10 of the original file.

![bam](https://i.imgur.com/C2demhB.gif)

### \*. Load multiple BAM files

`./bcov -bams <csvPath>`

Load multiple BAM files passing a CSV file in which each row should contain `<path>:<kitName>`

```
bams/P0971.bam, SureSelect Human All Exon V6 (100X)
bams/P0978.bam, SureSelect Human All Exon V7 (200X)
```

### 4. Run web server

`./bcov -web`

## Live

[https://bcov.qb.fcen.uba.ar](https://bcov.qb.fcen.uba.ar)
