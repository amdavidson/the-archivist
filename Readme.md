# The Archivist 

A simple service to archive various web services.

## Running The Archivist

The Archivist has a fairly typical command structure, the most basic command is in this format

    archivist ServiceName backup

## Service Instructions

### Github
    
To backup Github repositories, pull a [personal access token](https://github.com/settings/tokens) and add it 
to your `.the-archivist.yaml` config file.

    github: 
        user: my-username
        token: 238u2n3f98va8enf238as9d7hoi23dshaf

Then you can run `archivist github backup` and all your personal repositories will be stored in your `data` directory.



