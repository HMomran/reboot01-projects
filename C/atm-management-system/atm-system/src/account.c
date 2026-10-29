#include "header.h"

const char *RECORDS = "./data/records.txt";

int getAccountFromFile(FILE *ptr, char name[MAX_NAME_LEN], struct Record *r)
{
    return fscanf(ptr, "%d %d %s %d %d/%d/%d %s %d %lf %s",
                  &r->id,
		  &r->userId,
		  name,
                  &r->accountNbr,
                  &r->deposit.month,
                  &r->deposit.day,
                  &r->deposit.year,
                  r->country,
                  &r->phone,
                  &r->amount,
                  r->accountType) != EOF;
}

void saveAccountToFile(FILE *ptr, struct User u, struct Record r)
{
    fprintf(ptr, "%d %d %s %d %d/%d/%d %s %d %.2lf %s\n\n",
            r.id,
            u.id,
            u.name,
            r.accountNbr,
            r.deposit.month,
            r.deposit.day,
            r.deposit.year,
            r.country,
            r.phone,
            r.amount,
            r.accountType);
}

void checkAllAccounts(struct User u)
{
    char userName[MAX_COUNTRY_LEN];
    struct Record r;

    FILE *pf = fopen(RECORDS, "r");

    system("clear");
    printf("\t\t====== All accounts from user, %s =====\n\n", u.name);
    while (getAccountFromFile(pf, userName, &r))
    {
        if (strcmp(userName, u.name) == 0)
        {
            printf("_____________________\n");
            printf("\nAccount number:%d\nDeposit Date:%d/%d/%d \ncountry:%s \nPhone number:%d \nAmount deposited: $%.2f \nType Of Account:%s\n",
                   r.accountNbr,
                   r.deposit.day,
                   r.deposit.month,
                   r.deposit.year,
                   r.country,
                   r.phone,
                   r.amount,
                   r.accountType);
        }
    }
    fclose(pf);
    success(u);
}

void checkAccountDetails(struct User u)
{
    int accountNum;
    struct Record r;
    char userName[MAX_NAME_LEN];
    int found = 0;
    
    system("clear");
    printf("\t\t====== Account Details ======\n\n");
    printf("[Enter 0 to return to main menu]\n\n");
    
    printf("Enter account number: ");
    accountNum = getValidInt("");
    if (accountNum == 0) {
        mainMenu(u);
        return;
    }
    
    FILE *pf = fopen(RECORDS, "r");
    if (!pf) {
        printf("Error opening records file!\n");
        return;
    }
    
    while (getAccountFromFile(pf, userName, &r))
    {
        if (strcmp(userName, u.name) == 0 && r.accountNbr == accountNum)
        {
            found = 1;
            printf("\n--- Account Details ---\n");
            printf("Account Number: %d\n", r.accountNbr);
            printf("Date Created: %d/%d/%d\n", r.deposit.month, r.deposit.day, r.deposit.year);
            printf("Country: %s\n", r.country);
            printf("Phone: %d\n", r.phone);
            printf("Balance: $%.2f\n", r.amount);
            printf("Account Type: %s\n", r.accountType);
            
            double interest = calculateInterest(r.accountType, r.amount);
            if (interest > 0) {
                if (strcmp(r.accountType, "fixed01") == 0) {
                    double yearlyInterest = r.amount * INTEREST_FIXED01;
                    printf("\n-> You will get $%.2f as interest on %d/%d/%d\n", 
                           yearlyInterest, r.deposit.month, r.deposit.day, r.deposit.year + 1);
                } else if (strcmp(r.accountType, "fixed02") == 0) {
                    double yearlyInterest = r.amount * INTEREST_FIXED02 * 2;
                    printf("\n-> You will get $%.2f as interest on %d/%d/%d\n", 
                           yearlyInterest, r.deposit.month, r.deposit.day, r.deposit.year + 2);
                } else if (strcmp(r.accountType, "fixed03") == 0) {
                    double yearlyInterest = r.amount * INTEREST_FIXED03 * 3;
                    printf("\n-> You will get $%.2f as interest on %d/%d/%d\n", 
                           yearlyInterest, r.deposit.month, r.deposit.day, r.deposit.year + 3);
                } else {
                    printf("\n-> You will get $%.2f as interest on day %d of every month\n", 
                           interest, r.deposit.day);
                }
            } else if (strcmp(r.accountType, "current") == 0) {
                printf("\n-> You will not get interests because the account is of type current\n");
            }
            
            break;
        }
    }
    
    fclose(pf);
    
    if (!found) {
        printf("\n✖ Account not found!\n");
        askRetryOrMenu(u, checkAccountDetails);
        return;
    }
    
    success(u);
}
