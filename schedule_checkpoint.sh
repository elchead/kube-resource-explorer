# echo "./servicetool -application:itso -file:input-milp-noRuntime.zip -runtime:0 -jobsize:1 | grep id: > jobid.txt" | at now + 2 minutes
echo "WORKER=${WORKER} ./checkpoint.sh" | at now + 15 hours
