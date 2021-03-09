# image-migrate

A tiny tool for pull images, tagging them with a registry and pushing them to the registry.

## Usage

1. Connect to the Internet
1. Navigate to the location of the YAML file you wish to handle.
1. Run `im -r -skipPush --registry registry.mydomain.com .` . This will pull the images and tag them locally.
1. Now connect to the corporate network and run.
1. Run `im -r -u --skipPull --registry registry.mydomain.com .` . This will push the images to the registry and update the yaml files.




