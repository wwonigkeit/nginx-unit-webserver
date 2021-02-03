#!/bin/bash
  
#Check if the build has previously initialised
if [ -f "$WEBAPP" ]; then
        echo "The application has perviously been compiled - just start NGINX unit"
        /usr/local/bin/docker-entrypoint.sh unitd-debug --log /var/log/unitd.log --control 0.0.0.0:8080 --user root --group root
else
        # add NGINX Unit and Node.js repos
        # install git
        apt -o APT::Sandbox::User=root update
        apt -o APT::Sandbox::User=root install -y git openjdk-11-jdk
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

        # get the war file directory structure from the WEBAPP environment variable
        mkdir $(dirname "${WEBAPP}")
        # create the WAR file
        jar -cvf $WEBAPP *

        /usr/local/bin/docker-entrypoint.sh unitd-debug --log /var/log/unitd.log --control 0.0.0.0:8080 --user root --group root
fi
