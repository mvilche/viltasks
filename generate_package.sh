#!/bin/bash -e
rm -rf package viltasks.tar.gz
revel build . package/
cd package/
ln -sf src/viltasks/conf conf
mkdir -p database
tar zcf viltasks.tar.gz * && mv viltasks.tar.gz ../
cd -
rm -rf package/* && mkdir -p package/
mv viltasks.tar.gz package/
echo "PACKAGE GENERATED IN package/ FOLDER!!!"