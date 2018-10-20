# CoreDNS Governance

## Principles

The CoreDNS community adheres to the following principles:

- Open: CoreDNS is open source, advertised on [our website](https://coredns.io/community).
- Welcoming and respectful: See [Code of
  Conduct](https://github.com/coredns/coredns/blob/master/CODE-OF-CONDUCT.md).
- Transparent and accessible: Work and collaboration are done in public.
- Merit: Ideas and contributions are accepted according to their technical merit and alignment with
  project objectives, scope, and design principles.

## Project Lead

The CoreDNS project has a project lead.

A project lead in CoreDNS is
a single person that has a final say in any decision concerning the CoreDNS project.

The term of the project lead is two years, with no term limit restriction.

The project lead is elected by CoreDNS maintainers
according to an individual's technical merit to CoreDNS project.

The current "project lead" is identified in the [top level OWNERS
file](https://github.com/coredns/coredns/blob/master/OWNERS) with the string
`project lead` and the term behind the name.


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

## Changes in Project Lead

Changes in project lead or term is initiated by opening a github PR.

Anyone from CoreDNS community can vote on the PR with either +1 or -1.

Only the following votes are binding:
1) Any maintainer that has been listed in the top-level [OWNERS](OWNERS) file before the PR is opened.
2) Any maintainer from an organization may cast the vote for that organization. However, no organization
should have more binding votes than 1/5 of the total number of maintainers defined in 1).

The PR should be kept open for no less than 4 weeks. The PR can only be merged after the end of the
last project lead's term, with more +1 than -1 in the binding votes.

When there are conflicting PRs about changes in project lead, the PR with the most binding +1 votes is merged.

The project lead can volunteer to step down.

## Decision making process

Decisions are build on consensus between maintainers.
Proposals and ideas can either be submitted for agreement via a github issue or PR,
or by sending an email to `maintainer@coredns.io`.

In general, we prefer that technical issues and maintainer membership are amicably worked out between the persons involved.
If a dispute cannot be decided independently, the project lead has the final say to decide an issue.

Decision making process should be transparent and accessible to adhere to
the principles of CoreDNS project.

All proposals, ideas, and decisions, by maintainers or the project lead,
should either be part of a github issue or PR, or be sent to `maintainer@coredns.io`.

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
