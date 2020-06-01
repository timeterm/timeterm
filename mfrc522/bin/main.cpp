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

void handleInterrupt(int) {
    quit = true;
}

int main()
{
    Mfrc522::Device mfrc;
    signal(SIGTERM, handleInterrupt);

    mfrc.pcdInit();

    while (true) {
        // Look for a card
        if (!mfrc.piccIsNewCardPresent())
            continue;

        if (!mfrc.piccReadCardSerial())
            continue;

        // Print UID
        for (uint8_t i = 0; i < mfrc.getUid().size; ++i) {
            if (mfrc.getUid().uidByte[i] < 0x10) {
                printf(" 0");
                printf("%X", mfrc.getUid().uidByte[i]);
            } else {
                printf(" ");
                printf("%X", mfrc.getUid().uidByte[i]);
            }
        }
        printf("\n");
        delay(1000);

        if (quit) {
            std::cout << "Terminating" << std::endl;
            break;
        }
    }
    return 0;
}
