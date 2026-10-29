objects = src/main.o src/auth.o src/utils.o src/account.o src/operations.o src/transaction.o

atm : $(objects)
	cc -o atm $(objects)

main.o : src/header.h
auth.o : src/header.h
utils.o : src/header.h
account.o : src/header.h
operations.o : src/header.h
transaction.o : src/header.h

clean :
	rm -f $(objects)
