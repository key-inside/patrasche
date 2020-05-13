# Patrasche

## Update the version string in source codes
```sh
% make config
```

## Configuration

> Command line flags
```sh
% testdapp inspect --patrasche.block 1004
```

> Environment variables
```sh
% export TESTDAPP_PATRASCHE_BLOCK=1004
% testdapp inspect
```

> Configuration file
```sh
% testdapp inspect --patrasche.config=./config.yaml
```

> ARN
* SSM (Parameter Store)
```sh
% testdapp inspect --patrasche.config.region=ap-northeast-2 --patrasche.config.parameter=test/patrasche/config.yaml
% testdapp inspect --patrasche.config=arn:aws:ssm:region:account-id:parameter/test/patrasche/config.yaml
```
* Secrets Manager
```sh
% testdapp inspect --patrasche.config=arn:aws:secretsmanager:region:account-id:secret:test/patrasche/config.yaml
or via SSM
% testdapp inspect --patrasche.config.region=ap-northeast-2 --patrasche.config.parameter=/aws/reference/secretsmanager/test/patrasche/config.yaml
% testdapp inspect --patrasche.config=arn:aws:ssm:region:account-id:/aws/reference/secretsmanager/test/patrasche/config.yaml
```

* The extention of the ARN, .yaml or .json, is used for the content-type. If it's not set, the content-type is JSON.

### and ...

* See sample-config/config.yaml
