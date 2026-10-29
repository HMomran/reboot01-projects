# ATM Management System — Project Guide

## What it is

An educational, menu-driven bank account manager written in C. It runs in a terminal and stores users and account records in local text files. It simulates account operations; it has no connection to a bank or payment service.

## How it works

1. `main` opens the login/registration menu. Registration checks for an existing username, assigns the next user ID, and appends a user record. Login compares the supplied username and password against `data/users.txt`; terminal echo is disabled while entering the login password.
2. After login, `mainMenu` dispatches the selected option to an account or transaction function. Operations locate owned accounts using the logged-in username and the requested account number.
3. Creating an account collects a date, account number, country, phone number, positive starting balance, and account type. It checks for an account-number duplicate for that user and appends the record.
4. Details and list views scan `records.txt` and display matching accounts. Details include interest estimates based on account type.
5. Updates, deposits, withdrawals, ownership changes, and deletion write the records to `data/temp.txt`, modifying or omitting the selected record, then replace `records.txt`. This preserves other records in the normal success path.
6. Successful actions offer a return to the main menu. The program uses function calls between menus rather than a central iterative event loop.

## Code map

| File | Responsibility |
| --- | --- |
| `atm-system/src/main.c` | Initial menu, login flow, and operation dispatch |
| `atm-system/src/header.h` | User/record/date structs, limits, interest constants, and function declarations |
| `atm-system/src/auth.c` | Registration, login, and user-file lookup |
| `atm-system/src/account.c` | Record parsing/writing, listing, details, and interest display |
| `atm-system/src/operations.c` | Account creation, updates, and deletion |
| `atm-system/src/transaction.c` | Deposits, withdrawals, and ownership transfers |
| `atm-system/src/utils.c` | Numeric input, date checks, menu helpers, and interest calculation |
| `atm-system/Makefile` | Compiles six C files and links the `atm` executable |
| `atm-system/data/` | Persistent user and account text files |

## Build and run

Use macOS or Linux with Make and a C compiler such as Clang or GCC. Login uses POSIX `termios`, so this source is not a native Windows console implementation.

From `C/atm-management-system/`:

```bash
cd atm-system
make -B
./atm
```

`make -B` forces a rebuild rather than relying on the bundled executable. Always run from `atm-system/`: paths such as `./data/users.txt` are relative to the working directory. Use an interactive terminal because password entry changes terminal settings.

An equivalent direct build is:

```bash
cc -o atm src/main.c src/auth.c src/utils.c src/account.c src/operations.c src/transaction.c
./atm
```

`make clean` removes object files; it does not remove account data or the executable.

## Menu options

| Option | Behavior |
| --- | --- |
| 1 | Create a new account |
| 2 | Update country, phone, balance, or account type; zero keeps the current value |
| 3 | Inspect one account and its interest estimate |
| 4 | List accounts owned by the logged-in user |
| 5 | Deposit or withdraw money |
| 6 | Delete an account after confirmation |
| 7 | Transfer an entire account to an existing registered user |
| 8 | Exit |

Ownership transfer changes the account owner; it does not send a chosen amount to another account. Transactions require a positive amount and reject withdrawals larger than the balance. Deposits and withdrawals are blocked for `fixed01`, `fixed02`, and `fixed03` accounts.

## Account types and interest

| Type | Calculation displayed by the code |
| --- | --- |
| `saving` | Balance × 7% / 12, displayed as monthly interest |
| `current` | No interest |
| `fixed01` | Balance × 4%, displayed at deposit date plus one year |
| `fixed02` | Balance × 5% × 2, displayed at deposit date plus two years |
| `fixed03` | Balance × 8% × 3, displayed at deposit date plus three years |

These are estimates from constants in the source, using simple interest. There is no scheduled interest accrual or automatic balance crediting.

## Stored data

`data/users.txt` uses whitespace-separated fields:

```text
user_id username password
```

`data/records.txt` uses:

```text
record_id user_id username account_number month/day/year country phone balance account_type
```

The stored balance is formatted to two decimal places. Names and countries are single tokens because parsing uses whitespace separation. Record IDs are assigned within the user's records; account-number uniqueness is checked per owner. Keep the supplied text files well formed: the record reader checks for end-of-file rather than verifying every field was parsed.

## Example walkthrough

1. Register a new demo user, then log in with it.
2. Choose option 1 and create a saving account with a valid date, a new account number, and an initial balance of 100.
3. Use option 3 to inspect its balance and interest estimate.
4. Use option 5 to deposit 25, then withdraw 10. The expected balance is 115.
5. Exit, restart, and log in again to confirm the balance is persisted.
6. To explore ownership transfer, register a second demo user and transfer the account to that username. Log in as the recipient and list their accounts.

## Current limitations

Passwords are stored in plaintext, and `success` also prints the logged-in user's password. Use demonstration credentials. Account updates can directly change balances and account types, including fixed-deposit records. There is no transaction ledger, encryption, or concurrent-process file locking. Replacing files uses `remove` followed by `rename` without robust recovery, so this is a learning implementation rather than a durable banking system. Some string reads lack bounds, and repeated menu navigation nests function calls. Phone numbers are stored as integers, which loses leading zeroes and limits their range.

## Verification scope

The copy was checked against the original files. All six C source files compiled and linked successfully using `cc`, with the verification executable written outside the project. An interactive session and transaction changes were not performed during this documentation task.
