# ATM Management System

A terminal application written in C for managing simulated bank accounts. Users can register, log in, create accounts, inspect balances and interest estimates, update details, deposit or withdraw funds, delete accounts, and transfer account ownership.

See the [project guide](PROJECT_GUIDE.md) for how the code works, its file formats, build instructions, and a walkthrough.

## Quick start

Requires a C compiler, Make, and a macOS/Linux terminal.

```bash
cd atm-system
make -B
./atm
```

Run from `atm-system/` so the relative `data/` paths resolve. The program changes the copied data files; the original repository remains untouched.
