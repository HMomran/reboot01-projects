#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define MAX_NAME_LEN 50
#define MAX_COUNTRY_LEN 100
#define MAX_ACCOUNT_TYPE_LEN 10
#define MAX_PASSWORD_LEN 50

#define INTEREST_SAVING 0.07
#define INTEREST_FIXED01 0.04
#define INTEREST_FIXED02 0.05
#define INTEREST_FIXED03 0.08

#define MIN_YEAR 1900
#define MAX_YEAR 2100

struct Date
{
    int month, day, year;
};

struct Record
{
    int id;
    int userId;
    char name[MAX_COUNTRY_LEN];
    char country[MAX_COUNTRY_LEN];
    int phone;
    char accountType[MAX_ACCOUNT_TYPE_LEN];
    int accountNbr;
    double amount;
    struct Date deposit;
    struct Date withdraw;
};

struct User
{
    int id;
    char name[MAX_NAME_LEN];
    char password[MAX_PASSWORD_LEN];
};

int isValidDate(int day, int month, int year);
int clearInputBuffer();
int getValidInt(const char *prompt);
double getValidDouble(const char *prompt);
int isValidCountry(const char *country);
int isValidUsername(const char *username);
void getValidCountry(char *country);
int askRetryOrMenu(struct User u, void (*retryFunc)(struct User));
double calculateInterest(const char *accountType, double amount);
void stayOrReturn(int notGood, void f(struct User u), struct User u);
void success(struct User u);

int getAccountFromFile(FILE *ptr, char name[MAX_NAME_LEN], struct Record *r);
void saveAccountToFile(FILE *ptr, struct User u, struct Record r);

void loginMenu(struct User *u);
void registerMenu(char a[MAX_NAME_LEN], char pass[MAX_PASSWORD_LEN]);
const char *getPassword(struct User u);

void createNewAcc(struct User u);
void mainMenu(struct User u);
void checkAllAccounts(struct User u);
void checkAccountDetails(struct User u);
void updateAcc(struct User u);
void removeCurrentAcc(struct User u);
void transferOwner(struct User u);
void makeTransaction(struct User u);

