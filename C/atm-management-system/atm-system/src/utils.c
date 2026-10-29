#include "header.h"
#include <ctype.h>

int clearInputBuffer()
{
    int c;
    while ((c = getchar()) != '\n' && c != EOF);
    return 1;
}

int getValidInt(const char *prompt)
{
    int value;
    while (1) {
        printf("%s", prompt);
        if (scanf("%d", &value) == 1) {
            clearInputBuffer();
            return value;
        } else {
            printf("\n✖ Invalid input! Please enter a number.\n");
            clearInputBuffer();
        }
    }
}

double getValidDouble(const char *prompt)
{
    double value;
    while (1) {
        printf("%s", prompt);
        if (scanf("%lf", &value) == 1) {
            clearInputBuffer();
            return value;
        } else {
            printf("\n✖ Invalid input! Please enter a number.\n");
            clearInputBuffer();
        }
    }
}

void getValidCountry(char *country)
{
    int hasLetter = 0;
    while (1) {
        printf("\nEnter the country: ");
        scanf("%99s", country);
        clearInputBuffer();
        
        hasLetter = 0;
        for (int i = 0; country[i] != '\0'; i++) {
            if (isalpha(country[i])) {
                hasLetter = 1;
                break;
            }
        }
        
        if (hasLetter) {
            return;
        } else {
            printf("✖ Invalid country name! Must contain letters.\n");
        }
    }
}

int askRetryOrMenu(struct User u, void (*retryFunc)(struct User))
{
    int choice;
    printf("\nChoose an option:\n");
    printf("[1] Try again\n");
    printf("[0] Return to main menu\n");
    printf("Your choice: ");
    
    if (scanf("%d", &choice) != 1) {
        clearInputBuffer();
        mainMenu(u);
        return 0;
    }
    clearInputBuffer();
    
    if (choice == 1) {
        retryFunc(u);
        return 1;
    } else {
        mainMenu(u);
        return 0;
    }
}

int isValidCountry(const char *country)
{
    int i;
    int hasLetter = 0;
    
    if (strlen(country) == 0)
        return 0;
    
    for (i = 0; country[i] != '\0'; i++) {
        if (isalpha(country[i])) {
            hasLetter = 1;
            break;
        }
    }
    
    return hasLetter;
}

int isValidUsername(const char *username)
{
    int i;
    int hasLetter = 0;
    
    if (strlen(username) == 0)
        return 0;
    
    for (i = 0; username[i] != '\0'; i++) {
        if (isalpha(username[i])) {
            hasLetter = 1;
            break;
        }
    }
    
    return hasLetter;
}

int isValidDate(int day, int month, int year)
{
    if (year < MIN_YEAR || year > MAX_YEAR) return 0;
    if (month < 1 || month > 12) return 0;
    int daysInMonth[] = {31,28,31,30,31,30,31,31,30,31,30,31};
    if (month == 2) {
        int isLeap = (year % 4 == 0 && (year % 100 != 0 || year % 400 == 0));
        if (isLeap) daysInMonth[1] = 29;
    }
    if (day < 1 || day > daysInMonth[month-1]) return 0;
    return 1;
}

double calculateInterest(const char *accountType, double amount)
{
    if (strcmp(accountType, "saving") == 0) 
        return amount * INTEREST_SAVING / 12;
    if (strcmp(accountType, "fixed01") == 0) 
        return amount * INTEREST_FIXED01 / 12;
    if (strcmp(accountType, "fixed02") == 0) 
        return amount * INTEREST_FIXED02 / 12;
    if (strcmp(accountType, "fixed03") == 0) 
        return amount * INTEREST_FIXED03 / 12;
    return 0.0;
}

void stayOrReturn(int notGood, void f(struct User u), struct User u)
{
    int option;
    if (notGood == 0)
    {
        system("clear");
        printf("\n✖ Record not found!!\n");
    invalid:
        printf("\nChoose an option:\n");
        printf("[1] Try again\n");
        printf("[0] Return to main menu\n");
        printf("Your choice: ");
        scanf("%d", &option);
        if (option == 1)
            f(u);
        else if (option == 0)
            mainMenu(u);
        else
        {
            printf("Insert a valid operation!\n");
            goto invalid;
        }
    }
    else
    {
        printf("\nEnter 1 to go to the main menu and 0 to try again: ");
        scanf("%d", &option);
    }
    if (option == 1)
    {
        system("clear");
        mainMenu(u);
    }
    else
    {
        system("clear");
        f(u);
    }
}

void success(struct User u)
{
    printf("User Details:\n");
    printf("  ID      : %d\n", u.id);
    printf("  Name    : %s\n", u.name);
    printf("  Password: %s\n", u.password);
    int option;
    printf("\n✔ Success!\n\n");
invalid:
    printf("Press 1 to return to main menu: ");
    scanf("%d", &option);
    system("clear");
    if (option == 1)
    {
        mainMenu(u);
    }
    else
    {
        printf("Insert a valid operation!\n");
        goto invalid;
    }
}
