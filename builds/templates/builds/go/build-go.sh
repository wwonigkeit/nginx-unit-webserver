#!/bin/bash

#Check if the build has previously initialised
if [ -f "$WORKINGDIR/$EXECUTABLE" ]; then
        echo "The application has perviously been compiled - just start NGINX unit"
        /usr/local/bin/docker-entrypoint.sh unitd-debug --log /var/log/unitd.log --control 0.0.0.0:8080 --user root --group root
else
	# install build tools to compile the app
	apt update
	apt install --no-install-recommends --no-install-suggests -y build-essential -q
	apt install git -y -q

	# make a copy of the GitHub repo you're about to compile
	mkdir /src/ && cd /src/
	git clone $GITHUB_REPO
	# directory for the GitHub cloned repo
	set -- /src/*/
	CODENAME=$1
	# Get the name of the directory
	APPNAME=`ls`

	# make the working directory
	mkdir -p $WORKINGDIR
	# make nginx/unit package available at $GOPATH to compile the app
	cp -r /usr/share/gocode/src/* /usr/lib/go-1.11/src/
	cd $CODENAME
	/usr/lib/go-1.11/bin/go build -o $WORKINGDIR/$EXECUTABLE

	/usr/local/bin/docker-entrypoint.sh unitd-debug --log /var/log/unitd.log --control 0.0.0.0:8080 --user root --group root
fi
