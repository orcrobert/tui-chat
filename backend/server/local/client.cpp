#include <iostream>
#include <string>
#include <thread>
#include <unistd.h>
#include <sys/socket.h>
#include <arpa/inet.h>

#define PORT 8080
#define MAX_CLIENTS 20

using namespace std;

int main(int argc, char *argv[]) {
    // Create socket
    int server_socket;
    struct sockaddr_in server;
    server_socket = socket(AF_INET, SOCK_STREAM, 0);

    if (server_socket < 0) {
        cout << "Error creating socket!" << endl;
        return 1;
    }

    server.sin_family = AF_INET;
    server.sin_addr.s_addr = inet_addr("127.0.0.1");
    server.sin_port = htons(PORT);

    // Connect to the server
    if (connect(server_socket, (struct sockaddr *)&server, sizeof(server)) < 0) {
        cout << "Connection failed!" << endl;
        return 1;
    }

    string username;
    cout << "Enter username: ";
    getline(cin, username);

    // Send username to server
    send(server_socket, username.c_str(), username.length(), 0);

    // Start listening for input messages and sending them to the server
    thread recv_thread([&]() {
        char buffer[1024];
        while (true) {
            int bytes_received = recv(server_socket, buffer, sizeof(buffer), 0);
            if (bytes_received > 0) {
                buffer[bytes_received] = '\0';
                cout << buffer << endl;
            }
        }
    });

    // Send messages from stdin to the server
    string message;
    while (true) {
        getline(cin, message);
        if (message.empty()) continue;
        send(server_socket, message.c_str(), message.length(), 0);
    }

    close(server_socket);
    return 0;
}
