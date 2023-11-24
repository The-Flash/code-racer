apk add libseccomp-dev # necessary for nosocket to work
apk add util-linux # to get prlimit

prlimit
# check that the previous command worked or exists
if [ $? -ne 0 ]; then
    exit
fi

echo "READY" >> /READY