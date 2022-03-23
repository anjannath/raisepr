### How to run

To test `raisepr` locally, you need [ngrok](http://ngrok.com) to expose the http service that'll handle the callback request from Github

```bash
$ ngrok http 9999 # the http serivce listens on this port
```
1. Grab the URL from o/p of the above command. 
2. Create a new OAuth app in Github. Put the URL from previous step as the call back URL for the Github OAuth app.
3. Grab the **Client ID** and **Client Secret** from above step.
3. Build `raisepr`
```bash
$ make # this will create a raisepr binary on pwd
```
4. Run `raisepr` using the following command
```bash
# Replace the values of the env variables with the correct ones as obtained from above steps
$ REPO_NAME=<name_of_repo> REPO_OWNER=<anjannath> CLIENT_ID=<client_id_from_github> SECRET=<shhh> ./raisepr
```