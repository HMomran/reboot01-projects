#include "header.h"

void mainMenu(struct User u)
{
    int option;
    system("clear");
    printf("\n\n\t\t======= ATM =======\n\n");
    printf("\n\t\t-->> Feel free to choose one of the options below <<--\n");
    printf("\n\t\t[1]- Create a new account\n");
    printf("\n\t\t[2]- Update account information\n");
    printf("\n\t\t[3]- Check accounts\n");
    printf("\n\t\t[4]- Check list of owned account\n");
    printf("\n\t\t[5]- Make Transaction\n");
    printf("\n\t\t[6]- Remove existing account\n");
    printf("\n\t\t[7]- Transfer ownership\n");
    printf("\n\t\t[8]- Exit\n");
    option = getValidInt("");

    switch (option)
    {
    case 1:
        createNewAcc(u);
        break;
    case 2:
        updateAcc(u);
        break;
    case 3:
        checkAccountDetails(u);
        break;
    case 4:
        checkAllAccounts(u);
        break;
    case 5:
        makeTransaction(u);
        break;
    case 6:
        removeCurrentAcc(u);
        break;
    case 7:
        transferOwner(u);
        break;
    case 8:
        exit(1);
        break;
    default:
        printf("Invalid operation!\n");
    }
};

void initMenu(struct User *u)
{
    int r = 0;
    int option;
    system("clear");
    printf("\n\n\t\t======= ATM =======\n");
    printf("\n\t\t-->> Feel free to login / register :\n");
    printf("\n\t\t[1]- login\n");
    printf("\n\t\t[2]- register\n");
    printf("\n\t\t[3]- exit\n");
    while (!r)
    {
        option = getValidInt("");
        switch (option)
        {
        case 1:
            while (!r) {
                loginMenu(u);
                if (u->id == -2) {
                    system("clear");
                    printf("\n\n\t\t======= ATM =======\n");
                    printf("\n\t\t-->> Feel free to login / register :\n");
                    printf("\n\t\t[1]- login\n");
                    printf("\n\t\t[2]- register\n");
                    printf("\n\t\t[3]- exit\n");
                    break;
                }
                if (u->id != -1 && strcmp(u->password, getPassword(*u)) == 0)
                {
                    printf("\n\nPassword Match!");
                    r = 1;
                    break;
                }
                else
                {
                    printf("\n✖ Login failed!\n");
                    printf("\nChoose:\n");
                    printf("[1] Try again\n");
                    printf("[0] Return to main menu\n");
                    printf("Enter choice: ");
                    int choice = getValidInt("");
                    if (choice == 0) {
                        system("clear");
                        printf("\n\n\t\t======= ATM =======\n");
                        printf("\n\t\t-->> Feel free to login / register :\n");
                        printf("\n\t\t[1]- login\n");
                        printf("\n\t\t[2]- register\n");
                        printf("\n\t\t[3]- exit\n");
                        break;
                    }
                    system("clear");
                    printf("\n\n\t\t\t\t   Bank Management System\n\t\t\t\t\t User Login:\n");
                    printf("[Enter '0' at username to go back to main menu]\n\n");
                }
            }
            break;
        case 2:
            registerMenu(u->name,u->password);
            if (strlen(u->name) == 0) {
                system("clear");
                printf("\n\n\t\t======= ATM =======\n");
                printf("\n\t\t-->> Feel free to login / register :\n");
                printf("\n\t\t[1]- login\n");
                printf("\n\t\t[2]- register\n");
                printf("\n\t\t[3]- exit\n");
            } else {
                printf("\nPress Enter to return to login menu...");
                getchar();
                system("clear");
                printf("\n\n\t\t======= ATM =======\n");
                printf("\n\t\t-->> Feel free to login / register :\n");
                printf("\n\t\t[1]- login\n");
                printf("\n\t\t[2]- register\n");
                printf("\n\t\t[3]- exit\n");
            }
            break;
        case 3:
            exit(1);
            break;
        default:
            printf("Insert a valid operation!\n");
        }
    }
};

int main()
{
    struct User u;
    
    initMenu(&u);
    mainMenu(u);
    return 0;
}

