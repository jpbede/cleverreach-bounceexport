# cleverreach-bounceexport

Little go tool to export bounces as CSV from CleverReach.
It uses the CleverReach REST API v3.

# How to use
Create a OAuth app at https://eu.cleverreach.com/admin/account_rest.php

Then call the tool as following:

```cleverreach-bounceexport --oauth_id <your oauth client id> --oauth_secret <your oauth client secret>```

The tool always exports, by default, the bounces of the account to which the OAuth app belongs. 

#### Debug
If you want a log with a bit more infos, just append ```--debug```

#### Agencies
If you are a agency, using the whitelabel feature and have the right scope you can also export bounces from your sub-accounts 
with the following parameter ```--client_id <desired client id>```

Keep in mind this requires a extra scope, given by CleverReach.

# Questions/Ideas/Anything else
Just open a issue :)