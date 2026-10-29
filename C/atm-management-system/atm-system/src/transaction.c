#include "header.h"

extern const char *RECORDS;

void transferOwner(struct User u)
{
    int accountNum;
    char newOwnerName[50];
    struct Record r;
    char userName[50];
    int found = 0;
    int newOwnerExists = 0;
    int newOwnerId = -1;
    int confirm;
    
    system("clear");
    printf("\t\t====== Transfer Account Ownership ======\n\n");
    printf("[Enter 0 to return to main menu]\n\n");
    
    printf("Enter account number to transfer: ");
    accountNum = getValidInt("");
    if (accountNum == 0) {
        mainMenu(u);
        return;
    }
    
    printf("Enter new owner's username: ");
    scanf("%49s", newOwnerName);
    
    FILE *userFile = fopen("./data/users.txt", "r");
    if (userFile) {
        int userId;
        char userName2[50], userPass[50];
        while (fscanf(userFile, "%d %49s %49s", &userId, userName2, userPass) == 3) {
            if (strcmp(userName2, newOwnerName) == 0) {
                newOwnerExists = 1;
                newOwnerId = userId;
                break;
            }
        }
        fclose(userFile);
    }
    
    if (!newOwnerExists) {
        printf("\n✖ New owner '%s' does not exist!\n", newOwnerName);
        askRetryOrMenu(u, transferOwner);
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
            printf("Current Owner: %s\n", userName);
            printf("Country: %s\n", r.country);
            printf("Amount: $%.2f\n", r.amount);
            
            printf("\nTransfer this account to '%s'? (1=Yes, 0=No): ", newOwnerName);
            confirm = getValidInt("");
            
            if (confirm == 1) {
                strcpy(userName, newOwnerName);
                r.userId = newOwnerId;
                printf("\n✔ Ownership transferred successfully!\n");
            } else {
                printf("\n✖ Transfer cancelled!\n");
                found = 0;
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
    } else {
        remove("./data/temp.txt");
        printf("\n✖ Account not found!\n");
        askRetryOrMenu(u, transferOwner);
        return;
    }
    
    success(u);
}

void makeTransaction(struct User u)
{
    int accountNum;
    struct Record r;
    char userName[50];
    int found = 0;
    int transType;
    double amount;
    
    system("clear");
    printf("\t\t====== Make Transaction ======\n\n");
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
    
    found = 0;
    while (getAccountFromFile(pf, userName, &r))
    {
        if (strcmp(userName, u.name) == 0 && r.accountNbr == accountNum)
        {
            found = 1;
            
            printf("\n--- Account Found ---\n");
            printf("Account Number: %d\n", r.accountNbr);
            printf("Current Balance: $%.2f\n", r.amount);
            printf("Account Type: %s\n", r.accountType);
            
            if (strcmp(r.accountType, "fixed01") == 0 || 
                strcmp(r.accountType, "fixed02") == 0 || 
                strcmp(r.accountType, "fixed03") == 0) {
                printf("\n✖ Transactions are not allowed on %s accounts!\n", r.accountType);
                fclose(pf);
                askRetryOrMenu(u, makeTransaction);
                return;
            }
            break;
        }
    }
    fclose(pf);
    
    if (!found) {
        printf("\n✖ Account not found!\n");
        askRetryOrMenu(u, makeTransaction);
        return;
    }
    
    printf("\nSelect transaction type:\n");
    printf("[1] Deposit\n");
    printf("[2] Withdraw\n");
    printf("Enter your choice: ");
    transType = getValidInt("");
    
    if (transType != 1 && transType != 2) {
        printf("\n✖ Invalid transaction type!\n");
        askRetryOrMenu(u, makeTransaction);
        return;
    }
    
    printf("\nEnter amount: $");
    amount = getValidDouble("");
    
    if (amount <= 0) {
        printf("\n✖ Invalid amount!\n");
        askRetryOrMenu(u, makeTransaction);
        return;
    }
    
    pf = fopen(RECORDS, "r");
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
            if (transType == 1) {
                r.amount += amount;
                printf("\n✔ Deposited $%.2f\n", amount);
                printf("New Balance: $%.2f\n", r.amount);
            } else {
                if (r.amount < amount) {
                    printf("\n✖ Insufficient funds! Current balance: $%.2f\n", r.amount);
                    fclose(pf);
                    fclose(tempFile);
                    remove("./data/temp.txt");
                    askRetryOrMenu(u, makeTransaction);
                    return;
                } else {
                    r.amount -= amount;
                    printf("\n✔ Withdrew $%.2f\n", amount);
                    printf("New Balance: $%.2f\n", r.amount);
                }
            }
        }
        
        fprintf(tempFile, "%d %d %s %d %d/%d/%d %s %d %.2lf %s\n\n",
                r.id, r.userId, userName, r.accountNbr,
                r.deposit.month, r.deposit.day, r.deposit.year,
                r.country, r.phone, r.amount, r.accountType);
    }
    
    fclose(pf);
    fclose(tempFile);
    
    remove(RECORDS);
    rename("./data/temp.txt", RECORDS);
    
    success(u);
}
