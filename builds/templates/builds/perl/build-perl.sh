#!/bin/bash


### NOT COMPLETED
  
#Check if the build has previously initialised
if [ -f "$WORKINGDIR/$SCRIPT" ]; then
        echo "The application has perviously been copied - just start NGINX unit"
        /usr/local/bin/docker-entrypoint.sh unitd-debug --log /var/log/unitd.log --control 0.0.0.0:8080 --user root --group root
else
        # add NGINX Unit and Node.js repos
        # install git
        apt -o APT::Sandbox::User=root update
        apt -o APT::Sandbox::User=root install -y git
        # final cleanup
        apt remove -y build-essential
        apt clean && apt autoclean && apt autoremove --purge -y

        # make a copy of the GitHub repo
        mkdir /src/ && cd /src/
        git clone $GITHUB_REPO
        # directory for the GitHub cloned repo
        set -- /src/*/
        CODENAME=$1
        # Get the name of the directory
        APPNAME=`ls`
        cd $APPNAME

        # Make the working directory and copy the script to the correct location
        mkdir -p $WORKINGDIR && cd $WORKINGDIR
        # Copy the psgi script from the GitHub repo to the correct directory
        cp /$CODENAME/* $WORKINGDIR/

        /usr/local/bin/docker-entrypoint.sh unitd-debug --log /var/log/unitd.log --control 0.0.0.0:8080 --user root --group root
fi
