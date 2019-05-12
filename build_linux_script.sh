#!/bin/sh -xv

rm -rf princess_build princess
GOOS=linux go build
mkdir princess_build
cp -R statics princess_build
cp -R templates princess_build
mv Princess princess_build
mv princess_build princess
tar czf princess.tgz princess
rm -rf princess
