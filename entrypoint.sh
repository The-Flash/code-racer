#!/bin/bash
cp /bin/nosocket /opt/code-racer/nosocket
/bin/code-racer -f ${MANIFEST_PATH} -m ${MNTFS} -r ${RUNNERS_PATH} -n ${NOSOCKET}