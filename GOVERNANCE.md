# CoreDNS Governance

## Principles

The CoreDNS community adheres to the following principles:

- Open: CoreDNS is open source, advertised on [our website](https://coredns.io/community)
- Welcoming and respectful: See [Code of
  Conduct](https://github.com/coredns/coredns/blob/master/CODE-OF-CONDUCT.md).
- Transparent and accessible: Work and collaboration are done in public.
- Merit: Ideas and contributions are accepted according to their technical merit and alignment with
  project objectives, scope, and design principles.

## Benevolent Dictator for Life

The CoreDNS project has a benevolent dictator for life.

A [Benevolent dictator for life](https://en.wikipedia.org/wiki/Benevolent_dictator_for_life) is
a single person that has a final say in any decision concerning the CoreDNS project.

The benevolent dictator for life can volunteer to step down and appoint a new "dictator for life".

The current "benevolent dictator for life" is identified in the [top level OWNERS
file](https://github.com/coredns/coredns/OWNERS) with the string `benevolent dictator for life`
behind the name.

## Expectations from Maintainers

Every one carries water....

Making a community work requires input/effort from everyone. Maintainers should actively
participate in Pull Request reviews. Maintainers are expected to respond to assigned Pull Requests
in a *reasonable* time frame, either providing insights, or assign the Pull Requests to other
maintainers.

Every Maintainer is listed in the top-level [OWNERS](https://github.com/coredns/coredns/blob/master/OWNERS)
file, with their Github handle and a possibly obfuscated email address. Everyone in the
`approvers` list is a Maintainer.

A Maintainer is also listed in a plugin specific OWNER file.

A Maintainer should be a member of `maintainer@coredns.io`, although this is not a hard requirement.

## Becoming a Maintainers

On successful merge of a significant pull request, any current maintainer can reach
to the author behind the pull request and ask them if they are willing to become a CoreDNS
maintainer.

## Changes in Maintainership

If a Maintainer feels she/he can not fulfill the "Expectations from Maintainers", they are free to
step down.

The CoreDNS organization will never forcefully remove a current Maintainer, unless a maintainer
fails to meet the principles of CoreDNS community,
or adhere to the [Code of Conduct](https://github.com/coredns/coredns/blob/master/CODE-OF-CONDUCT.md).

## Github Project Administration

The __coredns__ GitHub project maintainers team reflects the list of Maintainers.

## Other Projects

The CoreDNS organization is open to receive new sub-projects under its umbrella. To accept a project
into the __CoreDNS__ organization, it has to meet the following criteria:

- Must be licensed under the terms of the Apache License v2.0.
- Must be related to one or more scopes of the CoreDNS ecosystem:
  - CoreDNS project artifacts (website, deployments, CI, etc ...)
  - External plugins
  - Other DNS related processing
- Must be supported by a Maintainer.

The submission process starts as a Pull Request or Issue on the
[coredns/coredns](https://github.com/coredns/coredns) repository with the required information
mentioned above. Once a project is accepted, it's considered a __CNCF sub-project under the umbrella
of CoreDNS__

## Code of Conduct

The [CoreDNS Code of Conduct](https://github.com/coredns/coredns/blob/master/CODE-OF-CONDUCT.md) is aligned with the CNCF Code of Conduct.

## Credits

Sections of this documents have been borrowed from [Fluentd](https://github.com/fluent/fluentd/blob/master/GOVERNANCE.md) and [Envoy](https://github.com/envoyproxy/envoy/blob/master/GOVERNANCE.md) projects.
