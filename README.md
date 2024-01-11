# Node ID

Command-line application and package written in Golang. Helps build identifying
data for Orcfax components.

## Usage

Run, `node-id -h` for more information its usage. Basic identity creation
requires calling the app with a URL pointing to an Orcfax validator websocket.

E.g.

```sh
./node-id -ws ws://<validator-websocket>/ws/node
```

## Based on IP Info

We're using IP Info. For more information on the tooling provided by IP Info
checkout the following:

* [Library source code][ipinfo-1].
* [Library documentation][ipinfo-2].
* [Short YCombinator article about IPInfo][ipinfo-3].

[ipinfo-1]: https://github.com/ipinfo/go/
[ipinfo-2]: https://pkg.go.dev/github.com/ipinfo/go/v2/ipinfo
[ipinfo-3]: https://news.ycombinator.com/item?id=37509114

IP Info is augmented to help generate a unique node identity for Orcfax
collector nodes.

A node identify object looks as follows:

```json
{
   "node_id": "bc178d43-6e11-421e-8d8e-976ed4754035",
   "location": {
      "ip": "0.0.0.0",
      "city": "City",
      "region": "Region",
      "country": "CO",
      "loc": "0.0000,0.0000",
      "org": "VM hosting Global",
      "postal": "88212",
      "timezone": "Coordinated Universal Time (UTC)",
      "readme": "https://ipinfo.io/"
   },
   "initialization": "2023-12-05T14:46:30Z",
   "validator_web_socket": "ws://validator-websocket"
}
```

And is converted into provenance data in the Orcfax COOP metadata.

## Provenance example

Orcfax validator records wrap information about collector nodes, e.g. a node
in the federated model's validation metadata.

```json
  "contributor": {
    "@type": "Organization",
    "name": "AS14061 DigitalOcean, LLC",
    "locationCreated": {
      "address": {
        "@type": "PostalAddress",
        "addressLocality": "London",
        "addressRegion": "England, GB,",
        "geo": "51.5085,-0.1257"
      },
      "additionalType": {
        "@type": "PropertyValue",
        "name": "ip address",
        "value": "167.71.137.7",
        "valueReference": "https://ipinfo.io/"
      }
    }
  },
```

For more context you can explore records in the Orcfax explorer.

* Explorer example: [urn:orcfax:5699f21a-4162-42e4-8409-cc50e6e9f992][explorer-1]

[explorer-1]: https://explorer.orcfax.io/5699f21a-4162-42e4-8409-cc50e6e9f992

## Makefile

The `makefile` provides a number of helpers. Run `make` to see those.

### Linting

Linting is a special case, you don't necessarily want it to fail when
you want an overview of all errors. To run all linters:

```sh
make --ignore-errors lint
```

vs. the following which will fail on errors. This isn't necessarily a big deal
as you will want to resolve those anyway.

```sh
make lint
```

### Goreleaser

The `makefile` wraps helpful `goreleaser` commands which help to test the
output of the code in release-like conditions.

The best way to install goreleaser is via Go with Go installed:

```sh
go install github.com/goreleaser/goreleaser@latest
```

More install options are available, see: [goreleaser website][releaser-1]

[releaser-1]: https://goreleaser.com/install/#nur

The quickest way to test locally is to run:

```sh
make build-local-snapshot
```

#### Semantic versioning

Goreleaser can be pedantic about how semantic versioning looks in git tags,
especially for release-candidates which are useful leading up to a release.

Valid semantic versioning looks as follows:

```text
vMM.mm.pp-rc.n
```

Where `-rc.n` are the components required for release candidates to be properly
identified and built.

> NB. If you are unsure at all about how to tag and version releases,take a
> look at the tag history, e.g. using:
>
> ```sh
> git tag -ln20
> ```
>
> Where the tags will be listed with their repspective tag messages up to
> 20 characters.

## Signing

Signing is enabled in the `goreleaser` component but the signing process is
currently a WIP. Currently checksums and binary objects will be signed with
the default gpg in your keyring.

### Configuring a GPG id

In `.goreleaser.yml` if the signing process needs to be modified to use a
specific key, then the `signs` section needs to be modified. Under `args`
change the command list to add your key-id or email address, e.g. from:

```sh
args: ["--output", "${signature}", "--detach-sign", "${artifact}"]
```

<!--markdownlint-disable -->

to:

```sh
args: ["-u", "<key-id or email_addr>", "--output", "${signature}", "--detach-sign", "${artifact}"]
```

<!--markdownlint-enable -->

### Verifying a signature

Signatures can be verified after creation using the options in the `makefile`.
They can also be verified manually with:

```sh
gpg --verify <signature-file>
```

If the `<signature-file>` is called `my-file.txt.sig` then GPG will assume the
original to verify the signature against is called `my-file.txt` returning:

```sh
$ gpg --verify my-file.txt.sig
gpg: assuming signed data in 'my-file.txt.sig'
gpg: Signature made Tue 05 Dec 2023 03:31:52 PM CET
gpg:                using RSA key E7274118998A052A32D77FC157B8D1DB7C7C611F
gpg:                issuer "developer@orcfax.io"
gpg: Good signature from "Anon Developer <developer@orcfax.io>" [ultimate]
```

or if the signature fails:

```sh
$ gpg --verify my-file.txt.sig
gpg: assuming signed data in 'my-file.txt'
gpg: Signature made Tue 05 Dec 2023 03:31:52 PM CET
gpg:                using RSA key E7274118998A052A32D77FC157B8D1DB7C7C611F
gpg:                issuer "developer@orcfax.io"
gpg: BAD signature from "Anon Developer <developer@orcfax.io>" [ultimate]
```

## Building private repos

While repositories such as this may eventually be open sourced, the development
may initially occur behind closed doors. This might, for example, reduce the
burden felt around API design or versioning, to give two examples.

Repositories can be imported into other golang projects even when private when
users have access to said repos, e.g. via ssh.

In fact, ssh makes this much easier in today's GitHub ecosystem where ssh is
now the default mechanism for access.

Git, however, still needs configuring with the following values set in the
environment (shown as a set of shell commands for convenience).

```sh
# Enable go to access ssh via git in git config.
git config --global url.git@github.com:.insteadOf https://github.com/
# View git changes.
cat ~/.gitconfig
# Add organizations or username as comma-separated-values to .bashrc.
export GOPRIVATE=github.com/orcfax/node-id/*
# Apply .bashrc changes immediately.
source ~/.bashrc
```
