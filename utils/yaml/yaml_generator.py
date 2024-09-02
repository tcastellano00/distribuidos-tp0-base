import sys

def get_client_container_definition(client_id):
    return  f"""
  client{client_id}:
    container_name: client{client_id}
    image: client:latest
    entrypoint: /client
    environment:
        - CLI_ID={client_id}
        - CLI_LOG_LEVEL=DEBUG
    networks:
        - testing_net
    depends_on:
        - server
    volumes:
        - ./client/config.yaml:/config.yaml
"""


def generate_docker_compose_yaml(file_name, clients_number):
    file_name_template = "./utils/yaml/docker-compose-dev-template.yaml"

    with open(file_name_template, "r") as file:
        content = file.read()
    
    clients_definitions = ""
    for client_id in range(0, int(clients_number)):
        clients_definitions += get_client_container_definition(client_id+1)

    content = content.replace("{{clients}}", clients_definitions)

    with open(file_name, "w") as file:
        file.write(content)


if __name__ == "__main__":
    if len(sys.argv) != 3:
        sys.exit(1)

    file_name = sys.argv[1]
    clients_number = sys.argv[2]

    generate_docker_compose_yaml(file_name, clients_number)