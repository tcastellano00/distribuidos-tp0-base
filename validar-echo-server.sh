
MESSAGE="test" 
SERVER_CONTAINER_NAME="server"
SERVER_CONTAINER_PORT=12345
NETWORK_NAME="tp0_testing_net"

RESPONSE=$(docker run --rm --network "$NETWORK_NAME" busybox:latest sh -c "echo $MESSAGE | nc -w 1 $SERVER_CONTAINER_NAME $SERVER_CONTAINER_PORT")

if [ "$RESPONSE" = "$MESSAGE" ]; then
    echo "action: test_echo_server | result: success"
else
    echo "action: test_echo_server | result: fail"
fi
