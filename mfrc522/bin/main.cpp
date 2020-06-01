#include <iostream>
#include <mfrc522/mfrc522.h>

#ifdef WIN32
#include <windows.h>
#else
#include <csignal>
#include <unistd.h>
#endif

void delay(int ms)
{
#ifdef WIN32
    Sleep(ms);
#else
    usleep(ms * 1000);
#endif
}

volatile bool quit = false;

void handleSignal(int) {
    quit = true;
}

int main()
{
    Mfrc522::Device rfidReader;
    signal(SIGINT, handleSignal);

    rfidReader.pcdInit();

    while (true) {
        if (quit) {
            std::cout << "Terminating" << std::endl;
            break;
        }

        // Look for a card
        if (!rfidReader.piccIsNewCardPresent())
            continue;

        if (!rfidReader.piccReadCardSerial())
            continue;

        // Print UID
        for (uint8_t i = 0; i < rfidReader.getUid().size; ++i) {
            if (rfidReader.getUid().uidByte[i] < 0x10) {
                printf(" 0");
                printf("%X", rfidReader.getUid().uidByte[i]);
            } else {
                printf(" ");
                printf("%X", rfidReader.getUid().uidByte[i]);
            }
        }
        printf("\n");
        delay(1000);
    }
    return 0;
}
