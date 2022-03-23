# Setup Droplet Image

To build SRS droplet image for [DigitalOcean Marketplace](https://marketplace.digitalocean.com/).

## Usage

First of all, please create a [token](https://cloud.digitalocean.com/account/api/tokens) and setup the env 
`DIGITALOCEAN_TOKEN`, for example:

```bash
export DIGITALOCEAN_TOKEN=xxx
```

Then [install Packer](https://www.packer.io/intro/getting-started/install.html):

```bash
brew tap hashicorp/tap
brew install hashicorp/tap/packer
```

Finally, start to build SRS image by:

```bash
packer build srs.json
```

or from root:

```bash
cd scripts/setup-droplet
PACKER_LOG=1 packer build srs.json
```

Please check the [snapshot](https://cloud.digitalocean.com/images/snapshots/droplets).

