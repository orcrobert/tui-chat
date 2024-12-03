#include <sqlite3.h>
#include <iostream>
#include <string>
#include <stdio.h>

using namespace std;

void initializeDatabase() {
    sqlite3* db;
    char* errMsg = nullptr;

    int rc = sqlite3_open("users.db", &db);
    if (rc) {
        printf("Cannot open database: %s\n", sqlite3_errmsg(db));
        return;
    }

    const char* createQuery = "CREATE TABLE IF NOT EXISTS users (username TEXT NOT NULL PRIMARY KEY, password TEXT NOT NULL)";

    rc = sqlite3_exec(db, createQuery, nullptr, nullptr, &errMsg);
    if (rc != SQLITE_OK) {
        printf("SQL error: %s\n", errMsg);
        sqlite3_free(errMsg);
    }
    else {
        printf("Database and table initialized successfully!\n");
    }

    sqlite3_close(db);
}

void addUser(const string& username, const string& password) {
    sqlite3* db;
    sqlite3_open("users.db", &db);

    const char* addQuery = "INSERT INTO users (username, password) VALUES (?, ?);";
    sqlite3_stmt* stmt;

    sqlite3_prepare(db, addQuery, -1, &stmt, nullptr);
    sqlite3_bind_text(stmt, 1, username.c_str(), -1, SQLITE_STATIC);
    sqlite3_bind_text(stmt, 2, password.c_str(), -1, SQLITE_STATIC);

    if (sqlite3_step(stmt) == SQLITE_DONE) {
        printf("User added successfully!\n");
    }
    else {
        printf("Failed to add user: %s\n", sqlite3_errmsg(db));
    }

    sqlite3_finalize(stmt);
    sqlite3_close(db);
}

int main() {
    initializeDatabase();

    // add some users for testing
    addUser("root", "toor");
    addUser("admin", "admin");
    addUser("robert", "caca");

    return 0;
}