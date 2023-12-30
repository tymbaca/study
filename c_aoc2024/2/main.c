#include <stdio.h>
#include <ctype.h>

int getNumFromLine(char *line) {
    char num1 = '\0';
    char num2 = '\0';
    for (int i = 0; i < 1024; i++) {
        char ch = line[i];

        if (ch == '\n') {
            break;
        }

        if (isdigit(ch)) {
            if (num1 == '\0') {
                num1 = ch;
            }
            num2 = ch;
        }
    }
    int n1 = num1 - '0';
    int n2 = num2 - '0';
    int result = n1*10 + n2; // 2 and 6 -> 20+6 -> 26
    return result;
}



int main() {
    char name[] = "Привет";
    printf("%s\n", name);
    printf("%zu\n", sizeof(name));
    char n = name[1];
    printf("%c\n", n);
    /* FILE *fptr = fopen("input.txt", "r"); */
    
    /* if (fptr == NULL) { */
    /*     printf("Cannot open the file\n"); */
    /*     return 1; */
    /* } */

    /* // Read line (with max size of 1024 bytes (chars)) */
    /* const int size = 1024; */
    /* char line[size]; */
    /* int sum = 0; */
    /* while (fgets(line, size, fptr)) { */
    /*     sum += getNumFromLine(line); */
    /* } */
    /* printf("%d\n", sum); */
}
