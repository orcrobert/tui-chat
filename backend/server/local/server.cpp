#include <iostream>
#include <string>
#include <vector>
#include <thread>
#include <mutex>
#include <cstring>
#include <sys/socket.h>
#include <arpa/inet.h>
#include <unistd.h>

#define PORT 8080
#define MAX_CLIENTS 20

std::vector<int> client_sockets;
std::vector<std::string> client_usernames;
std::mutex lock;

void broadcast_message(const std::string& message, int sender_socket) {
    std::lock_guard<std::mutex> guard(lock);
    for (size_t i = 0; i < client_sockets.size(); ++i) {
        if (client_sockets[i] != sender_socket) {
            send(client_sockets[i], message.c_str(), message.size() + 1, 0);
        }
    }
}

void handle_client(int client_socket) {
    char buffer[1024];
    int read_size;

    // Receive username
    read_size = recv(client_socket, buffer, sizeof(buffer), 0);
    std::string username(buffer);
    if (read_size > 0) {
        buffer[read_size] = '\0';

        {
            std::lock_guard<std::mutex> guard(lock);
            client_usernames.push_back(username);
        }

        std::cout << username << " connected" << std::endl;

        // Notify others about the new user
        std::string join_message = username + " has joined the chat.";
        broadcast_message(join_message, client_socket);
    }

    while ((read_size = recv(client_socket, buffer, sizeof(buffer), 0)) > 0) {
        buffer[read_size] = '\0';
        std::string received_message(buffer);

        std::string formatted_message = username + ": " + received_message;
        std::cout << formatted_message << std::endl;
        broadcast_message(formatted_message, client_socket);
    }

    // Client disconnected
    {
        std::lock_guard<std::mutex> guard(lock);
        for (size_t i = 0; i < client_sockets.size(); ++i) {
            if (client_sockets[i] == client_socket) {
                std::string leave_message = client_usernames[i] + " has left the chat.";
                broadcast_message(leave_message, client_socket);

                client_usernames.erase(client_usernames.begin() + i);
                client_sockets.erase(client_sockets.begin() + i);
                break;
            }
        }
    }

    close(client_socket);
    std::cout << "Client disconnected!" << std::endl;
}

int main() {
    int server_socket, client_socket, c;
    struct sockaddr_in server, client;

    server_socket = socket(AF_INET, SOCK_STREAM, 0);
    if (server_socket < 0) {
        std::cerr << "Socket could not be created!" << std::endl;
        return 1;
    }

    server.sin_family = AF_INET;
    server.sin_addr.s_addr = INADDR_ANY;
    server.sin_port = htons(PORT);

    if (bind(server_socket, (struct sockaddr*)&server, sizeof(server)) < 0) {
        perror("Bind failed!");
        return 1;
    }

    listen(server_socket, MAX_CLIENTS);
    std::cout << "Waiting for connections..." << std::endl;

    c = sizeof(struct sockaddr_in);

    while ((client_socket = accept(server_socket, (struct sockaddr*)&client, (socklen_t*)&c))) {
        {
            std::lock_guard<std::mutex> guard(lock);
            client_sockets.push_back(client_socket);
        }

        std::thread client_thread(handle_client, client_socket);
        client_thread.detach(); // Let the thread run independently
    }

    if (client_socket < 0) {
        perror("Accept failed!");
        return 1;
    }

    close(server_socket);
    return 0;
}
