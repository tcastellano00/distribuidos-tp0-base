
MESSAGE="test" 
SERVER_CONTAINER_NAME="server"

CONTAINER_IP=$(docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' "$SERVER_CONTAINER_NAME")
CONTAINER_PORT=12345

RESPONSE=$(echo "$MESSAGE" | nc $CONTAINER_IP $CONTAINER_PORT)

if [ "$RESPONSE" = "$MESSAGE" ];then
    echo "action: test_echo_server | result: success"
else
    echo "action: test_echo_server | result: fail"
fi
