#include "header.h"
#include <ctype.h>

extern const char *RECORDS;

const char *validTypes[] = {
    "saving", "current", "fixed01", "fixed02", "fixed03"
};
int numTypes = 5;

void createNewAcc(struct User u)
{
    struct Record r;
    struct Record cr;
    char userName[50];
    FILE *pf = fopen(RECORDS, "a+");

noAccount:
    system("clear");
    printf("\t\t\t===== New record =====\n");
    printf("[Enter 0 to return to main menu]\n");

    printf("\nEnter today's date (mm/dd/yyyy): ");
    while (1) {
        char dateInput[50];
        if (scanf("%s", dateInput) == 1) {
            if (strcmp(dateInput, "0") == 0) {
                clearInputBuffer();
                fclose(pf);
                mainMenu(u);
                return;
            }
            if (sscanf(dateInput, "%d/%d/%d", &r.deposit.month, &r.deposit.day, &r.deposit.year) == 3) {
                clearInputBuffer();
                if (r.deposit.month == 0 || r.deposit.day == 0 || r.deposit.year == 0) {
                    fclose(pf);
                    mainMenu(u);
                    return;
                }
                if (isValidDate(r.deposit.day, r.deposit.month, r.deposit.year)) {
                    break;
                } else {
                    printf("Invalid date! Please enter a valid date (mm/dd/yyyy): ");
                }
            } else {
                clearInputBuffer();
                printf("Invalid format! Please use mm/dd/yyyy (or 0 to go back): ");
            }
        } else {
            clearInputBuffer();
            printf("Invalid input! Please use mm/dd/yyyy (or 0 to go back): ");
        }
    }
    
    printf("\nEnter the account number:");
    r.accountNbr = getValidInt("");
    if (r.accountNbr == 0) {
        fclose(pf);
        mainMenu(u);
        return;
    }
    
    rewind(pf);
    while (getAccountFromFile(pf, userName, &cr))
    {
        if (strcmp(userName, u.name) == 0 && cr.accountNbr == r.accountNbr)
        {
            printf("\n✖ You already have an account with number %d!\n", r.accountNbr);
            printf("Choose an option:\n");
            printf("[1] Try another account number\n");
            printf("[0] Return to main menu\n");
            printf("Your choice: ");
            
            int choice;
            scanf("%d", &choice);
            
            if (choice == 1) {
                fclose(pf);
                pf = fopen(RECORDS, "a+");
                goto noAccount;
            } else {
                fclose(pf);
                mainMenu(u);
                return;
            }
        }
    }
    
    printf("\nEnter the country:");
    scanf("%s", r.country);
    
    while (!isValidCountry(r.country)) {
        printf("✖ Invalid country name! Country must contain letters.\n");
        printf("Enter the country:");
        scanf("%s", r.country);
    }
    
    printf("\nEnter the phone number:");
    r.phone = getValidInt("");
    printf("\nEnter amount to deposit: $");
    r.amount = getValidDouble("");
    
    if (r.amount <= 0) {
        printf("\n✖ Amount must be positive!\n");
        fclose(pf);
        askRetryOrMenu(u, createNewAcc);
        return;
    }
    
    printf("\nChoose the type of account:\n");
    printf("\t[1] Saving\n");
    printf("\t[2] Current\n");
    printf("\t[3] Fixed deposit (1 year)\n");
    printf("\t[4] Fixed deposit (2 years)\n");
    printf("\t[5] Fixed deposit (3 years)\n");
    printf("\nEnter your choice (1-5): ");
    
    int typeChoice;
    while (1) {
        typeChoice = getValidInt("");
        if (typeChoice >= 1 && typeChoice <= 5) {
            strcpy(r.accountType, validTypes[typeChoice - 1]);
            break;
        } else {
            printf("Invalid choice! Please enter 1-5: ");
        }
    }
    
    int maxUserAccountId = 0;
    rewind(pf);
    struct Record temp;
    char tempName[50];
    while (getAccountFromFile(pf, tempName, &temp)) {
        if (strcmp(tempName, u.name) == 0 && temp.id > maxUserAccountId)
            maxUserAccountId = temp.id;
    }
    r.id = maxUserAccountId + 1;

    saveAccountToFile(pf, u, r);

    fclose(pf);
    success(u);
}

void updateAcc(struct User u)
{
    int accountNum;
    struct Record r;
    char userName[50];
    int found = 0;
    
    system("clear");
    printf("\t\t====== Update Account Information ======\n\n");
    printf("[Enter 0 to return to main menu]\n\n");
    
    printf("Enter account number to update: ");
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
    
    FILE *tempFile = fopen("./data/temp.txt", "w");
    if (!tempFile) {
        printf("Error creating temp file!\n");
        fclose(pf);
        return;
    }
    
    while (getAccountFromFile(pf, userName, &r))
    {
        if (strcmp(userName, u.name) == 0 && r.accountNbr == accountNum)
        {
            found = 1;
            printf("\n--- Current Account Information ---\n");
            printf("Country: %s\n", r.country);
            printf("Phone: %d\n", r.phone);
            printf("Amount: $%.2f\n", r.amount);
            printf("Account Type: %s\n", r.accountType);
            
            printf("\n--- Enter New Information (Enter 0 to keep current) ---\n");
            printf("Enter new country (or 0 to keep current): ");
            char newCountry[100];
            scanf("%s", newCountry);
            if (strcmp(newCountry, "0") != 0) {
                while (!isValidCountry(newCountry)) {
                    printf("✖ Invalid country name! Country must contain letters.\n");
                    printf("Enter new country (or 0 to keep current): ");
                    scanf("%s", newCountry);
                    if (strcmp(newCountry, "0") == 0)
                        break;
                }
                if (strcmp(newCountry, "0") != 0)
                    strcpy(r.country, newCountry);
            }
            
            printf("Enter new phone number (or 0 to keep current): ");
            int newPhone = getValidInt("");
            if (newPhone != 0)
                r.phone = newPhone;
            
            printf("Enter new amount (or 0 to keep current): ");
            double newAmount = getValidDouble("");
            if (newAmount != 0) {
                if (newAmount < 0) {
                    printf("✖ Amount cannot be negative! Keeping current.\n");
                } else {
                    r.amount = newAmount;
                }
            }
            
            printf("\nEnter new account type (or 0 to keep current):\n");
            printf("  [1] Saving\n");
            printf("  [2] Current\n");
            printf("  [3] Fixed deposit (1 year)\n");
            printf("  [4] Fixed deposit (2 years)\n");
            printf("  [5] Fixed deposit (3 years)\n");
            
            int typeChoice;
            while (1) {
                printf("Enter your choice (1-5, or 0 to keep): ");
                typeChoice = getValidInt("");
                if (typeChoice >= 0 && typeChoice <= 5) {
                    break;
                } else {
                    printf("✖ Invalid choice! Please enter a number between 0-5.\n");
                }
            }
            
            if (typeChoice >= 1 && typeChoice <= 5) {
                strcpy(r.accountType, validTypes[typeChoice - 1]);
            }
        }
        
        fprintf(tempFile, "%d %d %s %d %d/%d/%d %s %d %.2lf %s\n\n",
                r.id, r.userId, userName, r.accountNbr,
                r.deposit.month, r.deposit.day, r.deposit.year,
                r.country, r.phone, r.amount, r.accountType);
    }
    
    fclose(pf);
    fclose(tempFile);
    
    if (found) {
        remove(RECORDS);
        rename("./data/temp.txt", RECORDS);
        printf("\n✔ Account updated successfully!\n");
    } else {
        remove("./data/temp.txt");
        printf("\n✖ Account not found!\n");
        askRetryOrMenu(u, updateAcc);
        return;
    }
    
    success(u);
}

void removeCurrentAcc(struct User u)
{
    int accountNum;
    struct Record r;
    char userName[50];
    int found = 0;
    
    system("clear");
    printf("\t\t====== Remove Account ======\n\n");
    printf("[Enter 0 to return to main menu]\n\n");
    
    printf("Enter account number to remove: ");
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
    
    FILE *tempFile = fopen("./data/temp.txt", "w");
    if (!tempFile) {
        printf("Error creating temp file!\n");
        fclose(pf);
        return;
    }
    
    while (getAccountFromFile(pf, userName, &r))
    {
        if (strcmp(userName, u.name) == 0 && r.accountNbr == accountNum)
        {
            found = 1;
            printf("\n--- Account Found ---\n");
            printf("Account Number: %d\n", r.accountNbr);
            printf("Country: %s\n", r.country);
            printf("Phone: %d\n", r.phone);
            printf("Amount: $%.2f\n", r.amount);
            printf("Account Type: %s\n", r.accountType);
            
            printf("\nAre you sure you want to delete this account? (1=Yes, 0=No): ");
            int confirm = getValidInt("");
            
            if (confirm != 1) {
                fprintf(tempFile, "%d %d %s %d %d/%d/%d %s %d %.2lf %s\n\n",
                        r.id, r.userId, userName, r.accountNbr,
                        r.deposit.month, r.deposit.day, r.deposit.year,
                        r.country, r.phone, r.amount, r.accountType);
                found = 0;
            }
        }
        else
        {
            fprintf(tempFile, "%d %d %s %d %d/%d/%d %s %d %.2lf %s\n\n",
                    r.id, r.userId, userName, r.accountNbr,
                    r.deposit.month, r.deposit.day, r.deposit.year,
                    r.country, r.phone, r.amount, r.accountType);
        }
    }
    
    fclose(pf);
    fclose(tempFile);
    
    if (found) {
        remove(RECORDS);
        rename("./data/temp.txt", RECORDS);
        printf("\n✔ Account deleted successfully!\n");
    } else {
        remove("./data/temp.txt");
        printf("\n✖ Account not found or deletion cancelled!\n");
        askRetryOrMenu(u, removeCurrentAcc);
        return;
    }
    
    success(u);
}
