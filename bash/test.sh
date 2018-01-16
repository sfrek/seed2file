#!/bin/bash


baseDir="/home/freko/src/go/src/github.com/sfrek/seed2file/bash"
testFile="${baseDir}/../test/file.hex"
md5sum="$(md5sum ${testFile} | cut -f1 -d\ )"
name="$(basename ${testFile})"
tempDir=$(mktemp -d)

dataSeed='{"data": [], "name": "'${name}'", "md5sum": "'${md5sum}'"}'
echo ${dataSeed} | jq -Src . | tee ${tempDir}/tempSeed.json

while read line;do
    jq ' .data += ["'${line}'"]' ${tempDir}/tempSeed.json > ${tempDir}/growing.json
    cp -a ${tempDir}/growing.json ${tempDir}/tempSeed.json
    echo -n '.'
done < <(cat ${testFile} | base64 )

cat ${tempDir}/tempSeed.json | jq -Src .
