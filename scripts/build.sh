#!/bin/bash

#sleep 1

BuildOutPath="$PWD"
username="duobuilduser"
emailaddress="duobuilduser@duosoftware.com"
password="DuoS12345"

mkdir "bin" 
mkdir "pkg" 
mkdir "src" 

export PATH=$PATH:/usr/local/go/bin;
export GOPATH=$PWD;
export PATH=$PATH:$GOPATH/bin;
cd bin
rm *
cd ../
cd src

git config --global user.name $username
git config --global user.email $emailaddress
echo ""
echo "BEGIN REPO PULL v6engine-deps"
mkdir "depo"
cd depo
git clone https://github.com/DuoSoftware/v6engine-deps
cd ../



cd depo/v6engine-deps
git pull
cp * -r $GOPATH/src
cd ../
cd ../
echo "END REPO PULL v6engine-deps"
echo ""
echo "BEGIN REPO PULL v6engine"
git clone -b development https://github.com/DuoSoftware/v6engine
mv v6engine duov6.com
