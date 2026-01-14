#!/usr/bin/env bash
#cd ..
#npm run build
#cd main
#VERSION=0.1
cpath=`pwd`
echo $cpath
PROJECT_PATH=${cpath%src*} #从右向左截取第一个 src 后的字符串
echo ${PROJECT_PATH}


export GOPATH=$GOPATH:${PROJECT_PATH}

alias SET_SOCKS5_PROXY='polipo socksParentProxy=127.0.0.1:1086 proxyAddress=0.0.0.0'

alias SET_HTTP_HTTPS_PROXY='env http_proxy=http://127.0.0.1:8123 https_proxy=http://127.0.0.1:8123'