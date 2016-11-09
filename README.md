# terraform-markdown-extractor
Extract json resource description from terraform repository markdown

## Usage:
*./extractor -dir=./_test/r/ -out=out.json*
* -dir: relative root path to terraform markdown like [AWS provider markdown](https://github.com/hashicorp/terraform/tree/master/website/source/docs/providers/aws), directories will be searched recursively
* -out: relative path to resulting file (in JSON)