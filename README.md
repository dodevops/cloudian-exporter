# Cloudian metrics exporter

This Prometheus exporter provides a metrics endpoint with metrics from a [Cloudian](https://cloudian.com/) installation.

**Note**: This is not affiliated in any way with Cloudian and Cloudian does not offer support for this.

## Configuration

The exporter requires the following configuration using environment variables:

* CLOUDIAN_URL: Cloudian CMC REST api base URL (e.g. https://cmc.company.com:19443)
* CLOUDIAN_USERNAME: REST api username (usually sysadmin)
* CLOUDIAN_PASSWORD: REST api password
* EXPORTER_REFRESH: Interval (in minutes) in which the metric values are updated [5]
* EXPORTER_LOGLEVEL: Loglevel to use [info]
* EXPORTER_LISTEN: Port to listen to [8080]

## Usage

Aside from building the go binary and running it, a container is available as well. Run it using e.g.:

   docker run --rm -it -P -e CLOUDIAN_URL=... -e CLOUDIAN_USERNAME=... -e CLOUDIAN_PASSWORD=... ghcr.io/dodevops/cloudian-exporter:main

## Exported metrics

* cloudian_bucket_size: Size of a bucket in Cloudian in bytes
  * Labels available
    * group_id: Group in Cloudian
    * user_id: User holding the bucket
    * bucket: Bucket name
