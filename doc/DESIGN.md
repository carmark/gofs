Introduction
------------

## Architecture
The general architecture  of the GoFS has three main components:
* the third-party cloud storage for the file data
* the metadata cluster service for managing the metadata and to support synchronization
* the client that implements most of the GoFS functionality, and corresponds to the file system client mounted at the user machine.
![Alt text](./gofs.png)

## Details
### Metadata server
The metadata service resorts to the coordination service to store file and directory metadata, together with information required for enforcing access control. Each file system object is represented by a metadata tuple containing: 
* The namespace (multi-tenant's identifier)
* The object name
* The type (file, directory or link)
* The parent object (in the hierarchical file namespace)
* The object metadata (size, date of creation, owner, ACLs, etc.)
* An opaque identifier referencing the file in the storage service

These two last fields represent the id and the hash stored in the consistency anchor. Metadata tuples are accessed through a set of operations offered by the local metadata service, which are then translated into different calls to the coordination service.

To avoid single point of failure, GoFS does not use master-slave mode to keep the service up, but uses [raft](https://raft.github.io/) protocol to do the consistency during the multiple nodes.

Metadata of the file, directory and link is stored in [BoltDB](https://github.com/boltdb/bolt),  and cached the active metadata information in memory.  Since we want to serve the public service, many users may have huge numbers of files, we can not cache all data in memory, such as GFS.  A proposed cache idea is a dynamic subtree partitioning strategy for distributing metadata across a cluster of metadata servers introduced in [Dynamic Metadata Management for Petabyte-scale File Systems](http://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.78.3205&rep=rep1&type=pdf).

We plan to provide the same [APIs of HDFS](https://hadoop.apache.org/docs/r1.0.4/webhdfs.html) to make it more compatible, more detailed API information will be introduced later.

### Client

On the client side, a client software is needed to do the FUSE interface implement and interact with remote APIs.  A simple implement will be did with [bazil's FUSE library](https://github.com/bazil/fuse). 

The key crux of the client side implement is the cache, it will decide the user experience.  The idea is to retain recently accessed file's metadata in the cache, so that repeated accesses to the same information can be handled locally, without additional network traffic. A client machine is sometimes faced with the problem of deciding whether a locally cached copy of data is consistent with the master copy (and hence can be used). If the client machine determines that its cached data are out of date, it must cache an up-to-date copy of the data before allowing further accesses.

There are two solutions for cache consistency:
* The client tries to check with metadata server at fixed time intervals
* When the server detects a potential inconsistency, it sends the update to each client.

Furthermore, a lots of papers introduced their research results, such as 
[Adaptation of Distributed File System to VDI Storage by Client-Side Cache](http://www.jcomputers.us/vol11/jcp1101-02.pdf)
[Knockoff: Cheap versions in the cloud](https://www.usenix.org/sites/default/files/fast17_full_proceedings.pdf#page=84)

### Third-party OSS
We plan to support several public cloud object storage service, such as aliyun oss, aws s3 and so on.  Most of the OSS has the s3 compatible API, so it is more easily to develop.

We use [minio object service](https://github.com/minio/minio) which is an open source object storage server compatible with Amazon S3 APIs to do the test.  You may also setup the `minio` and `GoFS` service in your private systems.
