#include <termios.h>
#include "header.h"

char *USERS = "./data/users.txt";

void loginMenu(struct User *u)
{
    struct termios oflags, nflags;
    char inputName[MAX_NAME_LEN], inputPass[MAX_PASSWORD_LEN];
    char fileUser[MAX_NAME_LEN], filePass[MAX_PASSWORD_LEN];
    int fileId;
    int found = 0;

    system("clear");
    printf("\n\n\n\t\t\t\t   Bank Management System\n\t\t\t\t\t User Login:\n");
    printf("[Enter '0' at username to go back to main menu]\n\n");
    printf("Enter username: ");
    if (scanf("%49s", inputName) != 1) {
        printf("Error reading username!\n");
        clearInputBuffer();
        return;
    }

    if (strcmp(inputName, "0") == 0) {
        u->id = -2;
        return;
    }

    tcgetattr(fileno(stdin), &oflags);
    nflags = oflags;
    nflags.c_lflag &= ~ECHO;
    nflags.c_lflag |= ECHONL;
    if (tcsetattr(fileno(stdin), TCSANOW, &nflags) != 0)
    {
        perror("tcsetattr");
        return;
    }
    printf("Enter password: ");
    if (scanf("%49s", inputPass) != 1) {
        printf("Error reading password!\n");
        clearInputBuffer();
        tcsetattr(fileno(stdin), TCSANOW, &oflags);
        return;
    }

    if (tcsetattr(fileno(stdin), TCSANOW, &oflags) != 0)
    {
        perror("tcsetattr");
        return;
    }

    FILE *fp = fopen("./data/users.txt", "r");
    if (!fp) {
        printf("Error opening users file!\n");
        return;
    }
    while (fscanf(fp, "%d %49s %49s", &fileId, fileUser, filePass) == 3) {
        if (strcmp(fileUser, inputName) == 0 && strcmp(filePass, inputPass) == 0) {
            found = 1;
            break;
        }
    }
    fclose(fp);

    if (found) {
        printf("\n✔ Login successful! Welcome, %s.\n", inputName);
        u->id = fileId;
        strcpy(u->name, fileUser);
        strcpy(u->password, filePass);
    } else {
        printf("\n✖ Invalid username or password.\n");
        u->id = -1;
    }
};

const char *getPassword(struct User u)
{
    FILE *fp;
    int fileId;
    char fileName[50], filePass[50];

    if ((fp = fopen("./data/users.txt", "r")) == NULL)
    {
        printf("Error! opening file");
        exit(1);
    }

    while (fscanf(fp, "%d %49s %49s", &fileId, fileName, filePass) == 3)
    {
        if (strcmp(fileName, u.name) == 0)
        {
            fclose(fp);
            char *buff = strdup(filePass);
            return buff;
        }
    }

    fclose(fp);
    return "no user found";
}

void registerMenu(char a[MAX_NAME_LEN], char pass[MAX_PASSWORD_LEN])
{
    FILE *user;
    struct User userChecker;
    int maxId;
    int currentId;
    int found;

    system("clear");
    printf("\n\n\n\t\t\t\t   Bank Management System\n\t\t\t\t\t User Registration:\n");
    printf("\n[Enter '0' at username to go back to main menu]\n");

    while (1) {
        printf("\nEnter username: ");
        if (scanf("%49s", a) != 1) {
            printf("Error reading username!\n");
            clearInputBuffer();
            continue;
        }

        if (strcmp(a, "0") == 0) {
            printf("\nReturning to main menu...\n");
            strcpy(a, "");
            strcpy(pass, "");
            return;
        }

        if (!isValidUsername(a)) {
            printf("\n✖ Invalid username! Username must contain at least one letter.\n");
            continue;
        }

        maxId = -1;
        found = 0;

        if ((user = fopen("./data/users.txt", "r")) != NULL)
        {
            while (fscanf(user, "%d %49s %49s", &currentId, userChecker.name, userChecker.password) != EOF)
            {
                if (strcmp(userChecker.name, a) == 0)
                {
                    found = 1;
                }
                if (currentId > maxId)
                {
                    maxId = currentId;
                }
            }
            fclose(user);
        }

        if (found)
        {
            printf("\n✖ Username '%s' already exists.\n", a);
            printf("\nChoose:\n");
            printf("[1] Try again\n");
            printf("[0] Return to main menu\n");
            printf("Enter choice: ");
            int choice;
            if (scanf("%d", &choice) != 1) {
                clearInputBuffer();
                continue;
            }
            clearInputBuffer();
            if (choice == 0) {
                strcpy(a, "");
                strcpy(pass, "");
                return;
            }
            system("clear");
            printf("\n\n\n\t\t\t\t   Bank Management System\n\t\t\t\t\t User Registration:\n");
            printf("\n[Enter '0' at username to go back to main menu]\n");
            continue;
        }

        break;
    }

    printf("Enter password: ");
    if (scanf("%49s", pass) != 1) {
        printf("Error reading password!\n");
        return;
    }

    if ((user = fopen("./data/users.txt", "a")) == NULL)
    {
        perror("Error opening users file for append");
        return;
    }

    int newId = maxId + 1;
    fprintf(user, "%d %s %s\n", newId, a, pass);
    fclose(user);

    printf("\n✔ Registration successful! Your ID is: %d\n", newId);
}