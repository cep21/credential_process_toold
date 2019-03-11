# credential_process_toold
Convert `credential_process` inside aws credentials file to an older style creds file.

## Install

`go get -u github.com/cep21/credential_process_toold`

With go1.11 or later,

`GO111MODULE=on go get -u github.com/cep21/credential_process_toold@v1.0.0`

## Why

Some older programs do not understand the `credential_process` field inside `~/.aws/credentials` files.
To get around this, you can read the file and regenerate it with credentials from an external process in a
format that older cli programs will understand.

## Example

Let's say you have a ~/.aws/credentials file like this.

```ini
[company-hatmaker]
region                = us-west-2
credential_process    = some_process --account-id 12345 --role Admin

[company-shoemaker]
region                = us-west-2
credential_process    = some_process --account-id 67891 --role Admin
```

You could then run the **credential_process_toold** binary on this file and pipe the output to another file

```bash
credential_process_toold generate > /tmp/another_file
```

If you were to cat this file, you will see generated credentials from the credentials_process.

```ini
[company-hatmaker]
region                = us-west-2
credential_process    = some_process --account-id 12345 --role Admin
aws_access_key_id     = DEADBEAFDEADBEAFDEADBEAF
aws_secret_access_key = ABCDEFG/ABCDEFGABCDEFG0123+
aws_session_token     = ABCDEFGABCDEFGABCDEFGABCDEFG+DEADBEAF+333211111


[company-shoemaker]
region                = us-west-2
credential_process    = some_process --account-id 67891 --role Admin
aws_access_key_id     = DEADBEAFDEADBEAFDEADBEAF
aws_secret_access_key = ABCDEFG/ABCDEFGABCDEFG0123+
aws_session_token     = ABCDEFGABCDEFGABCDEFGABCDEFG+DEADBEAF+333211111
```

You can now explicitly set the credentials file to what you generated and run older AWS-CLI style commands.

```bash
AWS_SHARED_CREDENTIALS_FILE=/tmp/another_file terraform init
```