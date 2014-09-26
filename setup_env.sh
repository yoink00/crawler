#!/bin/sh
####
# This script sets up some required Go development
# files.
####
##
# Check for the presence of 'go tool cover'
##
go tool cover > /dev/null 2>&1
RET=$?
if [ $RET -ne 0 -a $RET -ne 1 ]; then
	echo "Go coverage tool does not appear to be installed."
	echo "Attempting to install it..."
	# This needs to run as root because it goes and installs
	# things in naughty places. I'd rather it didn't but hey...
	# I do my development in a BSD jail so I don't care.
	sudo -E go get -u code.google.com/p/go.tools/cmd/cover
	# Restore the ownership of code.google.com
	sudo -E /usr/sbin/chown -R $USER:$GROUP ./src/code.google.com
fi
go get -u github.com/nsf/gocode
go get -u github.com/smartystreets/goconvey
go get -u github.com/puerkitobio/goquery
