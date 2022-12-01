#!/bin/sh

# until [ "$(curl -w '%{response_code}' --no-keepalive -o /dev/null --connect-timeout 1 http://user-service:8080/manager/text/list)" == "404" ];
# do echo --- sleeping for 1 second;
# sleep 1;
# done

until [ "$(curl -w '%{response_code}' --no-keepalive -o /dev/null --connect-timeout 1 http://user-service:8080/manager/text/list)" -eq "404" ]
do 
    echo "waiting for user-service to be up... sleeping for 1 second";
    sleep 2;
done
echo
echo "user-service UP!"
echo

# until [ "$(curl -w '%{response_code}' --no-keepalive -o /dev/null --connect-timeout 1 http://notification-service:8080/manager/text/list)" -eq "404" ]
# do 
#     echo "waiting for notification-service to be up... sleeping for 1 second";
#     sleep 2;
# done
# echo
# echo "notification-service UP!"
# echo

until [ "$(curl -w '%{response_code}' --no-keepalive -o /dev/null --connect-timeout 1 http://item-service:8088/manager/text/list)" -eq "404" ]
do 
    echo "waiting for item-service to be up... sleeping for 1 second";
    sleep 2;
done
echo
echo "item-service UP!"
echo

# run main
echo
echo "java services done building..."
echo 

cd /go/src/auctions-service-debug/main
go run . sql

