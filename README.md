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

## Bundles and relocation mappings

From an image relocation point of view, a _bundle_ is a software package which declares the images it uses,
typically using image references.

A _thin_ bundle's image references refer to repositories on the internet.
When the thin bundle is installed and run, the images are pulled from their repositories.

A _thick_ bundle's image references refer to binary images packaged with the bundle, typically in an archive file.
The images must somehow be loaded from the bundle before they can be used.

The images of a bundle can be relocated to a registry in which case a _relocation mapping_ maps
each original image reference to its relocated counterpart. The keys of the relocation mapping are the images
declared by the bundle. The relocation mapping needs to be applied to the bundle so that when the bundle
is installed and run, it will pull its images from the registry.

A thick bundle needs to be relocated before its images can be pulled.
A thin bundle _may_ be relocated, although this is not usually necessary.

Note: the terminology in this section originated in the [CNAB standard](https://cnab.io/).

## Example scenarios

The following scenarios, adapted from the
[relocation guide](https://github.com/deislabs/duffle/blob/master/docs/guides/relocation-guide.md) of the CNAB reference
implementation, describe relocation of thin and thick bundles.

### Thin bundle relocation

The [Acme Corporation](https://en.wikipedia.org/wiki/Acme_Corporation) needs to install some "forge" software packaged as a thin bundle (`forge.json`).
Acme is used to things going wrong, so they have evolved sophisticated processes to protect their systems.
In particular, all their production software must be loaded from Acme-internal repositories.
This protects them from outages when an external repository goes down.
It also gives them complete control over what software they run in production.

So Acme needs to pull the images referenced by `forge.json` from external repositories and store them in an Acme-internal registry.
This will be done in a DMZ with access to the internet and write access to the internal registry.

Suppose their internal registry is hosted at `registry.internal.acme.com` and they have created a user `smith` to manage the forge software.
They can relocate the images to their registry using a repository prefix `registry.internal.acme.com/smith`.
They can now install the bundle and use the relocation mapping to reconfigure the bundle to use
the relocated image names instead of the original image names.

When the bundle runs, the images are pulled from the internal registry.

### Thick bundle relocation

Gringotts Wizarding Bank (GWB) needs to install some software into a new coin sorting machine.
For GWB, security is paramount. Like Acme, all their production software must be loaded from internal repositories.
However, GWB regard a networked DMZ as too insecure. Their data center has no connection to the external internet.

Software is delivered to GWB encoded in Base64 and etched on large stones which are then rolled by hand into the
GWB data center, scanned, and decoded. The stones are stored for future security audits.

GWB obtains the new software as a thick bundle (`sort.tgz`) and relocates it to their private registry 
using a repository prefix of `registry.gold.gwb.dia/griphook`. 
This loads the images from `sort.tgz` into the private registry. Relocating from a thick bundle does not need
access to the original image repositories (which would prevent it from running inside the GWB data center).  

They can now install the bundle and use the relocation mapping to reconfigure the bundle to use
the relocated image names instead of the original image names.

Again when the bundle runs, the images are pulled from the internal registry.
Since relocation need not modify the original bundle or produce a new bundle, GWB can use the original stones
in security audits.

## Packages provided

The Go packages provided by this repository include:
 * some rich types representing image names and digests
 * a "path mapping" utility for relocating image names
 * a registry package for:
   * obtaining the digest of an image
   * copying images between repositories
   * copying images between repositories and an [OCI image layout](https://github.com/opencontainers/image-spec/blob/master/image-layout.md) on disk, e.g. to implement thick bundles.

For details, please refer to the [package documentation](https://godoc.org/github.com/pivotal/image-relocation).

### Docker daemon

This repository reads images directly from their repositories and does not attempt to read images
from the Docker daemon. This is primarily because the daemon doesn't guarantee to provide the 
same digest of an image as when the image has been pushed to a repository.

## Command line interface

A CLI, `irel`, is provided for manual use and experimentation. Issue `make irel` to build it.

## Where is this repository used?

This repository was originally factored out of the [Pivotal Function Service](https://pivotal.io/platform/pivotal-function-service)
which provided a command line interface for relocating the images in its distributions.

[duffle](https://github.com/deislabs/duffle), the [CNAB](https://cnab.io/) reference implementation, uses this repository to relocate bundles.

The riff project also [experimented](https://github.com/projectriff/cnab-k8s-installer-base) with using this
repository to create CNAB bundles which could relocate themselves (before duffle could do relocation).

## Alternatives

If this repository isn't quite what you're looking for, try:
* the underlying library: [ggcr](https://github.com/google/go-containerregistry)
* [kbld](https://github.com/k14s/kbld) (also based on ggcr)
* [k8s container image promoter](https://github.com/kubernetes-sigs/k8s-container-image-promoter) (also based on ggcr)

