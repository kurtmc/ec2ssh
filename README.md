ec2ssh
======

## Installation

```
go get github.com/kurtmc/ec2ssh && go install github.com/kurtmc/ec2ssh
```

## Usage

```
ec2ssh my-app*dev.company.com
root@my-app-02c5f6bda5a87c5f1.dev.company.com:~$

ec2ssh my-app*dev.company.com 'free -h'
my-app-077c3b56b553ae19e.company.com: total        used        free      shared  buff/cache   available
Mem:           3.9G        1.1G        2.0G         20M        834M        2.7G
Swap:            0B          0B          0B
my-app-02c5f6bda5a87c5f1.company.com: total        used        free      shared  buff/cache   available
Mem:           3.9G        1.1G        2.0G         15M        825M        2.7G
Swap:            0B          0B          0B
```
