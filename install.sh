#!/usr/bin/env bash

name=`gogogen`
archi=`arch`
if [ "$archi" == "x86_64" ];then
  archi="amd64"
elif [ "$archi" == "i386" ];then
  archi="arm64"
fi

os=`uname | tr '[A-Z]' '[a-z]'`

gopath=`go env var GOPATH | grep "/"`
if [ "$gopath" == "" ];then
  gopath=`go env var GOROOT | grep "/"`
fi

package=`curl -s https://api.github.com/repos/vine-io/${name}/releases/latest | grep browser_download_url | grep ${os} | cut -d'"' -f4 | grep "${name}-${os}-${archi}"`

echo "install package: ${package}"
wget ${package} -O /tmp/${name}.tar.gz && tar -xvf /tmp/${name}.tar.gz -C /tmp/

mv /tmp/${os}-${archi}/* $gopath/bin

rm -fr /tmp/${os}-${archi}
rm -fr /tmp/${name}.tar.gz
