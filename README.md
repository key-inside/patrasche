# Patrasche

## Update the version string in source codes
```sh
% make config
```

## Configuration

> Command line flags
```sh
% patrasche inspect --block 1004
```

> Environment variables
```sh
% export PATRASCHE_BLOCK=1004
% patrasche inspect
```

> Configuration file
```sh
% partrasche inspect --config=./config.yaml
```

> ARN
* SSM (Parameter Store)
```sh
% partrasche inspect --config.region=ap-northeast-2 --config.parameter=test/patrasche/config.yaml
% partrasche inspect --config=arn:aws:ssm:region:account-id:parameter/test/patrasche/config.yaml
```
* Secrets Manager
```sh
% partrasche inspect --config=arn:aws:secretsmanager:region:account-id:secret:test/patrasche/config.yaml
or via SSM
% partrasche inspect --config.region=ap-northeast-2 --config.parameter=/aws/reference/secretsmanager/test/patrasche/config.yaml
% partrasche inspect --config=arn:aws:ssm:region:account-id:/aws/reference/secretsmanager/test/patrasche/config.yaml
```

* The extention of the ARN, .yaml or .json, is used for the content-type. If it's not set, the content-type is JSON.

### and ...

* --network, --identity, --keep and extened options can be ARN
* The config file can contain ARN values
* See sample-config/config.yaml
