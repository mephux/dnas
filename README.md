# DNAS - Domain Name Analytics System
[![Build Status](https://drone.io/github.com/mephux/dnas/status.png)](https://drone.io/github.com/mephux/dnas/latest)

Logs all DNS questions and answers for searching and metrics. DNA answers are stored as
a bloom filter for better performance (less questions asked to the embeded key/value store).

DNAS supports logging to an embeded Bolt (https://github.com/boltdb/bolt) key / value store 
(-d database.db) or to a flat file as json (-w filename.txt).

The plan is to continue to scale DNAS for powerful searches and metrics. 
i.e malware blah.exe sent data to blah.org what ips did that resolve to at that time.

## Install

  * Using Go Get

    * Note: You will need libpcap-dev before you build DNAS.
    * `go get github.com/mephux/dnas`

  * Using Git & Go Build
  
    * Note: You will need libpcap-dev before you build DNAS.
    * `git clone https://github.com/mephux/dnas.git`
    * `cd dnas`
    * `make`

  * Using Vagrant & Docker

    * `vagrant up`

## OUTPUT

  `Example: sudo dnas -i en0 -H -u mephux`

  * Output Support:

    * sqlite3
    * Mysql
    * Postgres
    * Json

## Usage

```
  DNAS (0.1.0) - Domain Name Analytics System

  Usage: dnas [options]

  Options:
    -i, --interface=eth0          Interface to monitor
    -p, --port=53                 DNS port (53)
    -D, --daemon                  Run DNAS in daemon mode
    -w, --write=FILE              Write JSON output to log file
    -u, --user=USER               Drop privileges to this user
    -H, --hexdump                 Show hexdump of DNS packet
        --mysql                   Enable Mysql Output Support
        --postgres                Enable Postgres Output Support
        --sqlite3                 Enable Sqlite3 Output Support
        --db-user=root            Database User (root)
        --db-password=PASSWORD    Database Password
        --db-database=dnas        Database Database (dnas)
        --db-host=127.0.0.1       Database Host
        --db-port=3306            Database Port
        --db-path=~/.dnas.db      Path to Database on disk. (sqlite3 only)
    -q, --quiet                   Suppress DNAS output
    -v, --version                 Show version information

  Help Options:
    -h, --help                    Show this help message
```

## Self-Promotion

Like DNAS? Follow the repository on
[GitHub](https://github.com/mephux/dnas) and if
you would like to stalk me, follow [mephux](http://dweb.io/) on
[Twitter](http://twitter.com/mephux) and
[GitHub](https://github.com/mephux).

# MIT LICENSE

The MIT License (MIT) - [LICENSE](https://github.com/mephux/dnas/blob/master/LICENSE)
