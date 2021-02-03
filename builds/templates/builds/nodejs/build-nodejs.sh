#!/bin/bash
  
#Check if the build has previously initialised
if [ -f "$WORKINGDIR/$EXECUTABLE" ]; then
        echo "The application has perviously been compiled - just start NGINX unit"
        /usr/local/bin/docker-entrypoint.sh unitd-debug --log /var/log/unitd.log --control 0.0.0.0:8080 --user root --group root
else
        # add NGINX Unit and Node.js repos
        apt -o APT::Sandbox::User=root update
        apt -o APT::Sandbox::User=root install -y apt-transport-https gnupg1
        curl -sL https://nginx.org/keys/nginx_signing.key | apt-key add -
        echo "deb https://packages.nginx.org/unit/debian/ buster unit" > /etc/apt/sources.list.d/unit.list
        echo "deb-src https://packages.nginx.org/unit/debian/ buster unit" >> /etc/apt/sources.list.d/unit.list
        curl https://deb.nodesource.com/setup_12.x | bash -
        # install build chain
        apt -o APT::Sandbox::User=root update
        apt -o APT::Sandbox::User=root install -y build-essential nodejs unit-dev git
        # add dependencies locally
        npm install -g --unsafe-perm unit-http
        # final cleanup
        apt remove -y build-essential unit-dev apt-transport-https gnupg1
        apt clean && apt autoclean && apt autoremove --purge -y
        rm -rf /var/lib/apt/lists/* /etc/apt/sources.list.d/*.list

        # make a copy of the GitHub repo
        mkdir /$WORKINGDIR/ && mkdir /src/ && cd /src/
        git clone $GITHUB_REPO
        # directory for the GitHub cloned repo
        set -- /src/*/
        CODENAME=$1
        # Get the name of the directory
        APPNAME=`ls`

        # make app.js executable; link unit-http locally
        cd /$CODENAME && mv -f * /$WORKINGDIR/
        cd /$WORKINGDIR/ && chmod +x $EXECUTABLE
        pwd
        ls -la
        npm link unit-http && npm install express && npm install yargs
        pwd
        ls -la

        /usr/local/bin/docker-entrypoint.sh unitd-debug --log /var/log/unitd.log --control 0.0.0.0:8080 --user root --group root
fi
