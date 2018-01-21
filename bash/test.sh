#!/bin/bash

testFile="file.hex"
md5sum="$(md5sum ${testFile} | cut -f1 -d\ )"
name="$(basename ${testFile})"
tempDir=$(mktemp -d)


seed='{
    "type":"{{ type }}",
    "timestamp":"{{ timestamp }}",
    "files":[]
}'

dataSeed='{"data": [], "name": "'${name}'", "md5sum": "'${md5sum}'"}'
echo ${dataSeed} | jq -Src . | tee ${tempDir}/tempSeed.json
echo ${seed} | jq -Src . | tee ${tempDir}/seed.json

while read line;do
    jq ' .data += ["'${line}'"]' ${tempDir}/tempSeed.json > ${tempDir}/growing.json
    cp -a ${tempDir}/growing.json ${tempDir}/tempSeed.json
    echo -n '.'
done < <(cat ${testFile} | base64 )

jq '.files += ["'$(cat ${tempDir}/tempSeed.json | jq -Src . )'"]' > "${tempDir}/toSend.json"

cat "${tempDir}/toSend.json" | jq .
