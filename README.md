# collector-clair

[![codecov](https://codecov.io/gh/rode/collector-clair/branch/main/graph/badge.svg)](https://codecov.io/gh/rode/collector-clair)

A rode collector for clair scans

## Running standalone clair scanner
1. Download the clair scanner https://github.com/arminc/clair-scanner/releases
2. Get the DB started up as well as the scanner:
    ```
    docker run -p 5432:5432 -d --name db arminc/clair-db:latest
    docker run -p 6060:6060 --link db:postgres -d --name clair arminc/clair-local-scan:latest```
3. Execute the scan
    ```
    ./clair-scanner --report output.json --ip $(ipconfig getifaddr en0) alpine:3.5
    ```
    You must specify your ip address on the network. Using localhost/127.0.0.1/0.0.0.0 doesn't seem to work
4. Analyze the output. Here is a sample json file:
    ```
    {
        "image": "abc4d437894f",
        "unapproved": [
            "CVE-2020-28928"
        ],
        "vulnerabilities": [
            {
                "featurename": "musl",
                "featureversion": "1.1.24-r8",
                "vulnerability": "CVE-2020-28928",
                "namespace": "alpine:v3.12",
                "description": "",
                "link": "https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2020-28928",
                "severity": "Low",
                "fixedby": "1.1.24-r10"
            }
        ]
    }
    ```
