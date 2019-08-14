# Docker/OCI image relocation

[![GoDoc](https://godoc.org/github.com/pivotal/image-relocation?status.svg)](https://godoc.org/github.com/pivotal/image-relocation)
[![Go Report Card](https://goreportcard.com/badge/pivotal/image-relocation)](https://goreportcard.com/report/pivotal/image-relocation)
[![Build Status](https://dev.azure.com/projectriff/pivotal-image-relocation/_apis/build/status/pivotal.image-relocation?branchName=master)](https://dev.azure.com/projectriff/pivotal-image-relocation/_build/latest?definitionId=11&branchName=master)
[![codecov](https://codecov.io/gh/pivotal/image-relocation/branch/master/graph/badge.svg)](https://codecov.io/gh/pivotal/image-relocation)

This repository contains a Go module for relocating Docker and OCI images.

## What is image relocation?
_Relocating_ an image means copying it to another repository, possibly in a private registry.

Using a separate registry has some advantages:
* It provides complete control over when the image is updated or deleted:
    * This provides isolation from unwanted updates or deletion of the original image.
    * If the image becomes stale, for instance when it has known vulnerabilities, it can be deleted.
* The registry can be hosted on a private network for security or other reasons.

A highly desirable property of image relocation is that the image digest of the relocated image is the same as that of the original images.
This gives the user confidence that the relocated image consists of the same bits as the original image.

## Relocating image names
An image name consists of a domain name (with optional port) and a path. The image name may also contain a tag and/or a digest.
The domain name determines the network location of a registry.
The path consists of one or more components separated by forward slashes.
The first component is sometimes, by convention for certain registries, a user name providing access control to the image.

Let’s look at some examples:
* The image name `docker.io/istio/proxyv2` refers to an image with user name `istio` residing in the docker hub registry at `docker.io`.
* The image name `projectriff/builder:v1` is short-hand for `docker.io/projectriff/builder:v1` which refers to an image with user name `projectriff` also residing at `docker.io`.
* The image name `gcr.io/cf-elafros/knative-releases/github.com/knative/serving/cmd/autoscaler@sha256:deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef` refers to an image with user name `cf-elafros` residing at `gcr.io`.

When an image is relocated to a registry, the domain name is set to that of the registry.
Relocation takes a _repository prefix_ which is used to prefix the relocated image names.

The path of a relocated image may:
* Include the original user name for readability
* Be “flattened” to accommodate registries which do not support hierarchical paths with more than two components
* End with a hash of the image name (to avoid collisions)
* Preserve any tag in the original image name
* Preserve any digest in the original image name.

For instance, when relocated to a repository prefix `example.com/user`, the above image names might become something like this:
* `example.com/user/istio-proxyv2-f93a2cacc6cafa0474a2d6990a4dd1a0`
* `example.com/user/projectriff-builder-a4a25a99d48adad8310050be267a10ce:v1`
* `example.com/user/cf-elafros-knative-releases-github.com-knative-serving-cmd-autoscaler-c74d62dc488234d6d1aaa38808898140@sha256:deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef`

The hash added to the end of the relocated image path should not depend on any tag and/or digest in
the original image name. This ensures a one-to-one mapping between repositories. In other words, if:

    x maps to y

where `x` and `y` are image names without tags or digests, then

    x:t maps to y:t (for all tags t)

and

    x@d maps to y@d (for all digests d).

## Packages provided

The Go packages provided by this repository include:
 * some rich types representing image names and digests
 * a "path mapping" utility for relocating image names
 * a registry package for copying images between repositories and between repositories and an [OCI image layout](https://github.com/opencontainers/image-spec/blob/master/image-layout.md) on disk.

For details, please refer to the [package documentation](https://godoc.org/github.com/pivotal/image-relocation).

### Docker daemon

This repository reads images directly from their repositories and does not attempt to read images
from the Docker daemon. This is primarily because the daemon doesn't guarantee to provide the 
same digest of an image as when the image has been pushed to a repository.

## Command line interface

A CLI, `irel`, is provided for manual use and experimentation. Issue `make irel` to build it.

