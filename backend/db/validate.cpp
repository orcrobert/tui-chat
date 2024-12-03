#include <sqlite3.h>
#include <stdio.h>
#include <iostream>
#include <string>
#include <filesystem>

using namespace std;

bool validateUser(const string& username, const string& password) {
    sqlite3* db;
    string dbPath = filesystem::absolute("../../build/users.db").string();

    if (sqlite3_open(dbPath.c_str(), &db) != SQLITE_OK) {
        printf("Error opening database: %s\n", sqlite3_errmsg(db));
        return false;
    }

    const char* query = "SELECT COUNT(*) FROM users WHERE username = ? AND password = ?;";
    sqlite3_stmt* stmt;

    if (sqlite3_prepare(db, query, -1, &stmt, nullptr) != SQLITE_OK) {
        printf("Error preparing statement: %s\n", sqlite3_errmsg(db));
        sqlite3_close(db);
        return false;
    }

    sqlite3_bind_text(stmt, 1, username.c_str(), -1, SQLITE_STATIC);
    sqlite3_bind_text(stmt, 2, password.c_str(), -1, SQLITE_STATIC);

    int count = 0;
    if (sqlite3_step(stmt) == SQLITE_ROW) {
        count = sqlite3_column_int(stmt, 0);
    }

    sqlite3_finalize(stmt);
    sqlite3_close(db);

    return count > 0;
}

int main() {
    string username, password;
    cin >> username >> password;

    if (validateUser(username, password)) {
        printf("Valid!\n");
    }
    else {
        printf("Invalid!\n");
    }

    return 0;
}