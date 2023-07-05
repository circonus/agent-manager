# Installer

1. install package for OS
1. call `circonus-am --register=token` to register the Agent Manager instance
    1. the registration token is sent as a header to the registration endpoint along with some json with claims (e.g. `{"sub":"hostname"}`)
    1. if all is good, a jwt will be returned and saved in `etc/.id/token`
    1. token is used with all other requests to the api
1. if successful, start the circonus-am service
